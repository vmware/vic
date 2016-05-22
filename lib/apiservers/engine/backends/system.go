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
	"fmt"
	"runtime"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
)

type System struct {
	ProductName string
}

func (s *System) SystemInfo() (*types.Info, error) {
	Driver := "Portlayer Storage"
	IndexServerAddress := "https://index.docker.io/v1/"
	ServerVersion := "0.0.1"
	Name := "VIC"

	info := &types.Info{
		Driver:             Driver,
		IndexServerAddress: IndexServerAddress,
		ServerVersion:      ServerVersion,
		Name:               Name,
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

func (s *System) AuthenticateToRegistry(authConfig *types.AuthConfig) (string, error) {
	return "", fmt.Errorf("%s does not implement System.AuthenticateToRegistry", s.ProductName)
}
