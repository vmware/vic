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

package backends

//****
// system.go
//
// Rules for code to be in here:
// 1. No remote or swagger calls.  Move those code to system_portlayer.go
// 2. Always return docker engine-api compatible errors.
//		- Do NOT return fmt.Errorf()
//		- Do NOT return errors.New()
//		- DO USE the aliased docker error package 'derr'
//		- It is OK to return errors returned from functions in system_portlayer.go

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	urlfetcher "github.com/vmware/vic/pkg/fetcher"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"

	"github.com/docker/docker/pkg/platform"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/go-units"
)

type System struct {
	systemProxy VicSystemProxy
}

const (
	systemStatusMhz    = " VCH CPU limit"
	systemStatusMemory = " VCH memory limit"
	systemStatusCPUUsageMhz    = " VCH CPU usage"
	systemStatusMemUsage = " VCH memory usage"
	systemOS           = " VMware OS"
	systemOSVersion    = " VMware OS version"
	systemProductName  = " VMware Product"
	volumeStoresID     = "VolumeStores"
	loginTimeout       = 20 * time.Second
)

func NewSystemBackend() *System {
	return &System{
		systemProxy: &SystemProxy{},
	}
}

func (s *System) SystemInfo() (*types.Info, error) {
	defer trace.End(trace.Begin("SystemInfo"))
	client := PortLayerClient()

	// Retrieve container status from port layer
	running, paused, stopped, err := s.systemProxy.ContainerCount()
	if err != nil {
		log.Infof("System.SytemInfo unable to get global status on containers: %s", err.Error())
	}

	// Build up the struct that the Remote API and CLI wants
	info := &types.Info{
		Driver:             PortLayerName(),
		IndexServerAddress: IndexServerAddress,
		ServerVersion:      ProductVersion(),
		ID:                 ProductName(),
		Containers:         running + paused + stopped,
		ContainersRunning:  running,
		ContainersPaused:   paused,
		ContainersStopped:  stopped,
		Images:             getImageCount(),
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
		Name:          VchConfig().Name,
		KernelVersion: "",
		Architecture:  platform.Architecture, //stubbed

		// NOTE: These values have no meaning for VIC.  We default them to true to
		// prevent the CLI from displaying warning messages.
		CPUCfsPeriod:      true,
		CPUCfsQuota:       true,
		CPUShares:         true,
		CPUSet:            true,
		OomKillDisable:    true,
		MemoryLimit:       true,
		SwapLimit:         true,
		KernelMemory:      true,
		IPv4Forwarding:    true,
		BridgeNfIptables:  true,
		BridgeNfIP6tables: true,
		HTTPProxy:         "",
		HTTPSProxy:        "",
		NoProxy:           "",
	}

	// Add in network info from the VCH via guestinfo
	for _, network := range VchConfig().ContainerNetworks {
		info.Plugins.Network = append(info.Plugins.Network, network.Name)
	}

	info.SystemStatus = make([][2]string, 0)

	// Add in volume label from the VCH via guestinfo
	volumeStoreString, err := FetchVolumeStores(client)
	if err != nil {
		log.Infof("Unable to get the volume store list from the portlayer : %s", err.Error())
	} else {
		customInfo := [2]string{volumeStoresID, volumeStoreString}
		info.SystemStatus = append(info.SystemStatus, customInfo)
	}

	if s.systemProxy.PingPortlayer() {
		status := [2]string{PortLayerName(), "RUNNING"}
		info.SystemStatus = append(info.SystemStatus, status)
	} else {
		status := [2]string{PortLayerName(), "STOPPED"}
		info.SystemStatus = append(info.SystemStatus, status)
	}

	// Add in vch information
	vchInfo, err := s.systemProxy.VCHInfo()
	if err != nil || vchInfo == nil {
		log.Infof("System.SystemInfo unable to get vch info from port layer: %s", err.Error())
	} else {
		if vchInfo.CPUMhz > 0 {
			info.NCPU = int(vchInfo.CPUMhz)

			customInfo := [2]string{systemStatusMhz, fmt.Sprintf("%d Mhz", info.NCPU)}
			info.SystemStatus = append(info.SystemStatus, customInfo)
		}
		if vchInfo.Memory > 0 {
			info.MemTotal = vchInfo.Memory * 1024 * 1024 // Get Mebibytes

			customInfo := [2]string{systemStatusMemory, units.BytesSize(float64(info.MemTotal))}
			info.SystemStatus = append(info.SystemStatus, customInfo)
		}
		if vchInfo.CPUUsage >= 0 {
			customInfo := [2]string{systemStatusCPUUsageMhz, fmt.Sprintf("%d Mhz", int(vchInfo.CPUUsage))}
			info.SystemStatus = append(info.SystemStatus, customInfo)
		}
		if vchInfo.MemUsage >= 0 {
			customInfo := [2]string{systemStatusMemUsage, units.BytesSize(float64(vchInfo.MemUsage))}
			info.SystemStatus = append(info.SystemStatus, customInfo)
		}
		if vchInfo.HostProductName != "" {
			customInfo := [2]string{systemProductName, vchInfo.HostProductName}
			info.SystemStatus = append(info.SystemStatus, customInfo)
		}
		if vchInfo.HostOS != "" {
			info.OperatingSystem = vchInfo.HostOS
			info.OSType = vchInfo.HostOS //Value for OS and OS Type the same from vmomi

			customInfo := [2]string{systemOS, vchInfo.HostOS}
			info.SystemStatus = append(info.SystemStatus, customInfo)
		}
		if vchInfo.HostOSVersion != "" {
			customInfo := [2]string{systemOSVersion, vchInfo.HostOSVersion}
			info.SystemStatus = append(info.SystemStatus, customInfo)
		}
	}

	return info, nil
}

// layout for build time as per constants defined in https://golang.org/src/time/format.go
const buildTimeLayout = "2006/01/02@15:04:05"

func (s *System) SystemVersion() types.Version {
	APIVersion := "1.22"
	Arch := runtime.GOARCH

	BuildTime := version.BuildDate
	if t, err := time.Parse(buildTimeLayout, BuildTime); err == nil {
		// match time format from docker version's output
		BuildTime = t.Format(time.ANSIC)
	}

	Experimental := true
	GitCommit := version.GitCommit
	GoVersion := runtime.Version()
	// FIXME: fill with real kernel version
	KernelVersion := "-"
	Os := runtime.GOOS
	Version := version.Version
	if Version != "" && Version[0] == 'v' {
		// match version format from docker version's output
		Version = Version[1:]
	}

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
	defer trace.End(trace.Begin(""))

	fetcher := urlfetcher.NewURLFetcher(urlfetcher.Options{
		Timeout:  loginTimeout,
		Username: authConfig.Username,
		Password: authConfig.Password,
		RootCAs:  RegistryCertPool,
	})

	// Only look at V2 registries
	registryAddress := authConfig.ServerAddress
	if !strings.Contains(authConfig.ServerAddress, "/v2") {
		registryAddress = registryAddress + "/v2/"
	}

	registryURL, err := url.Parse(registryAddress)
	if err != nil {
		msg := fmt.Sprintf("Bad login address: %s", registryAddress)
		log.Errorf(msg)
		return msg, "", err
	}

	// Check if requested registry is in our list of allowed insecure registries
	var insecureOk bool
	insecureRegistries := InsecureRegistries()
	for _, registry := range insecureRegistries {
		if registry == registryURL.Host {
			insecureOk = true
			break
		}
	}

	dologin := func(scheme string) (string, error) {
		registryURL.Scheme = scheme

		var authURL *url.URL

		// Attempt to get the Auth URL from HEAD operation to the registry
		hdr, err := fetcher.Head(registryURL)
		if err == nil && fetcher.IsStatusUnauthorized() {
			authURL, err = fetcher.ExtractOAuthURL(hdr.Get("www-authenticate"), nil)
		}
		if err != nil {
			log.Errorf("Looking up OAuth URL failed: %s", err)
			return "", err
		}

		log.Debugf("logging onto %s", authURL.String())

		// Just check if we get a token back.
		token, err := fetcher.FetchAuthToken(authURL)
		if err != nil || token.Token == "" {
			log.Errorf("Fetch auth token failed: %s", err)
			return "", err
		}

		return token.Token, nil
	}

	_, err = dologin("https")
	if err != nil && insecureOk {
		_, err = dologin("http")
	}

	if err != nil {
		return "", "", err
	}

	// We don't return the token.  The config.json will store token if we return
	// it, but the regular docker daemon doesn't seem to return it  either.
	return "Login Succeeded", "", nil
}

// Utility functions

func getImageCount() int {
	images := cache.ImageCache().GetImages()
	return len(images)
}

func FetchVolumeStores(client *client.PortLayer) (string, error) {

	res, err := client.Storage.VolumeStoresList(storage.NewVolumeStoresListParamsWithContext(ctx))
	if err != nil {
		return "", err
	}

	return strings.Join(res.Payload.Stores, " "), nil
}
