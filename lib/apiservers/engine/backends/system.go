// Copyright 2016 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vicbackends

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/misc"
	"github.com/vmware/vic/pkg/trace"

	"github.com/docker/docker/pkg/parsers/kernel"
	"github.com/docker/docker/pkg/platform"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
)

type System struct {
}

type ContainerStatus struct {
	count      int
	numRunning int
	numStopped int
	numPaused  int
}

const (
	etcReleaseFile    = "/etc/os-release"
	usrlibReleaseFile = "/usr/lib/os-release"
	unknown           = "<unknown"
)

func (s *System) SystemInfo() (*types.Info, error) {
	defer trace.End(trace.Begin("SystemInfo"))

	// Use docker pkgs to get some system data
	kernelVersion := unknown
	if kv, err := kernel.GetKernelVersion(); err != nil {
		log.Warnf("Could not get kernel version: %v", err)
	} else {
		kernelVersion = kv.String()
	}

	operatingSystem := getOperatingSystem()

	meminfo, err := system.ReadMemInfo()
	if err != nil {
		log.Errorf("Could not read system memory info: %v", err)
	}

	// Check if portlayer server is up
	plClient := PortLayerClient()

	systemStatus := make([][2]string, 1)
	systemStatus[0][0] = PortLayerName()
	if pingPortlayer(plClient) {
		systemStatus[0][1] = "RUNNING"
	} else {
		systemStatus[0][1] = "STOPPED"
	}

	// Retrieve number of images from storage port layer
	numImages := getImageCount(plClient)
	if err != nil {
		log.Infof("System.SytemInfo unable to get image count: %s.", err.Error())
	}

	// Retieve container status from port layer
	containerStatus, err := getContainerStatus(plClient)
	if err != nil {
		log.Infof("System.SytemInfo unable to get global status on containers: ", err.Error())
	}

	// Build up the struct that the Remote API and CLI wants
	info := &types.Info{
		Driver:             PortLayerName(),
		IndexServerAddress: IndexServerAddress,
		ServerVersion:      ProductVersion(),
		ID:                 ProductName(),
		Containers:         containerStatus.count,
		ContainersRunning:  containerStatus.numRunning,
		ContainersPaused:   containerStatus.numPaused,
		ContainersStopped:  containerStatus.numStopped,
		Images:             numImages,
		SystemStatus:       systemStatus,
		Debug:              VchConfig().Diagnostics.DebugLevel > 0,
		NGoroutines:        runtime.NumGoroutine(),
		SystemTime:         time.Now().Format(time.RFC3339Nano),
		ExecutionDriver:    PortLayerName(),
		LoggingDriver:      "",
		CgroupDriver:       "",
		DockerRootDir:      "",
		ClusterStore:       "",
		ClusterAdvertise:   "",

		// FIXME: Get this info once we have event listening service
		//	NEventsListener    int

		// These are system related.  Some refer to cgroup info.  Others are
		// retrieved from the port layer and are information about the resource
		// pool.
		// FIXME: update this once we get the resource pool info
		Name:            VchConfig().Name,
		KernelVersion:   kernelVersion,         //stubbed
		OperatingSystem: operatingSystem,       //stubbed
		OSType:          platform.OSType,       //stubbed
		Architecture:    platform.Architecture, //stubbed
		NCPU:            runtime.NumCPU(),      //stubbed
		MemTotal:        meminfo.MemTotal,      //stubbed
		CPUCfsPeriod:    false,
		CPUCfsQuota:     false,
		CPUShares:       false,
		CPUSet:          false,
		OomKillDisable:  true,
		//	MemoryLimit        bool
		//	SwapLimit          bool
		KernelMemory: false,
		//	IPv4Forwarding     bool
		//	BridgeNfIptables   bool
		//	BridgeNfIP6tables  bool `json:"BridgeNfIp6tables"`
		HTTPProxy:  "",
		HTTPSProxy: "",
		NoProxy:    "",
	}

	// Add in network info from the VCH via guestinfo
	for _, network := range VchConfig().ContainerNetworks {
		info.Plugins.Network = append(info.Plugins.Network, network.Name)
	}

	// Add in volume info from the VCH via guestinfo
	for _, location := range VchConfig().VolumeLocations {
		info.Plugins.Volume = append(info.Plugins.Volume, location.String())
	}

	return info, nil
}

func (s *System) SystemVersion() types.Version {
	APIVersion := "1.22"
	Arch := runtime.GOARCH
	// FIXME: fill with real build time
	BuildTime := "-"
	Experimental := true
	// FIXME: fill with real commit id
	GitCommit := "-"
	GoVersion := runtime.Version()
	// FIXME: fill with real kernel version
	KernelVersion := "-"
	Os := runtime.GOOS
	Version := "0.0.1"

	// go runtime panics without this so keep this here
	// until we find a repro case and report it to upstream
	_ = Arch

	version := types.Version{
		APIVersion:    APIVersion,
		Arch:          Arch,
		BuildTime:     BuildTime,
		Experimental:  Experimental,
		GitCommit:     GitCommit,
		GoVersion:     GoVersion,
		KernelVersion: KernelVersion,
		Os:            Os,
		Version:       Version,
	}

	return version
}

func (s *System) SubscribeToEvents(since, sinceNano int64, ef filters.Args) ([]events.Message, chan interface{}) {
	return make([]events.Message, 0, 0), make(chan interface{})
}

func (s *System) UnsubscribeFromEvents(chan interface{}) {

}

func (s *System) AuthenticateToRegistry(ctx context.Context, authConfig *types.AuthConfig) (string, string, error) {
	return "", "", fmt.Errorf("%s does not implement System.AuthenticateToRegistry", ProductName())
}

// Utility functions

// Returns the name of the OS on the appliance vm.  Very likely Photon OS.  This
// function takes the place of Docker's function that relies on other packages.
// Since VIC will likely run on either Photon OS or devbox (ubuntu), we can
// make some assumptions that avoid having to pull in other packages.
func getOperatingSystem() string {
	releaseFile, err := os.Open(etcReleaseFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return "<unknown>"
		}
		releaseFile, err = os.Open(usrlibReleaseFile)
		if err != nil {
			return "<unknown>"
		}
	}

	var prettyName string
	re := regexp.MustCompile("\"")
	scanner := bufio.NewScanner(releaseFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			parts := strings.Split(line, "=")
			prettyName = re.ReplaceAllString(parts[1], "")

			return prettyName
		}
	}

	return "Linux"
}

func pingPortlayer(plClient *client.PortLayer) bool {
	if plClient != nil {
		pingParams := misc.NewPingParams()
		_, err := plClient.Misc.Ping(pingParams)
		if err != nil {
			log.Info("Ping to portlayer failed")
			return false
		}
		return true
	}

	log.Errorf("Portlayer client is invalid")
	return false
}

func getImageCount(plClient *client.PortLayer) int {

	images := cache.ImageCache().GetImages()
	return len(images)
}

func getContainerStatus(plClient *client.PortLayer) (ContainerStatus, error) {
	return ContainerStatus{0, 0, 0, 0}, nil
}
