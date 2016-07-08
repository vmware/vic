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
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	httptransport "github.com/go-swagger/go-swagger/httpkit/client"
	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/metadata"
)

const (
	Imagec             = "imagec"
	PortlayerName      = "Backend Engine"
	IndexServerAddress = "registry-1.docker.io"

	// Retries defines how many times to attempt refreshing the image cache at startup
	Retries = 5

	// RetryTimeSeconds defines how many seconds to wait between retries
	RetryTimeSeconds = 2
)

var (
	portLayerClient     *client.PortLayer
	portLayerServerAddr string
	portLayerName       string
	productName         string
	productVersion      string

	vchConfig *metadata.VirtualContainerHostConfigSpec
)

func Init(portLayerAddr, product string, config *metadata.VirtualContainerHostConfigSpec) error {
	_, _, err := net.SplitHostPort(portLayerAddr)
	if err != nil {
		return err
	}

	vchConfig = config
	productName = product

	if config != nil {
		productVersion = config.Version
		if productVersion == "" {
			portLayerName = product + " Backend Engine"
		} else {
			portLayerName = product + " " + productVersion + " Backend Engine"
		}
	} else {
		portLayerName = product + " Backend Engine"
	}

	t := httptransport.New(portLayerAddr, "/", []string{"http"})
	portLayerClient = client.New(t, nil)
	portLayerServerAddr = portLayerAddr

	// attempt to update the image cache at startup
	log.Info("Refreshing image cache...")
	go func() {
		for i := 0; i < Retries; i++ {

			// initial pause to wait for the portlayer to come up
			time.Sleep(RetryTimeSeconds * time.Second)

			if err := cache.ImageCache().Update(portLayerClient); err == nil {
				log.Info("Image cache updated successfully")
				return
			}
			log.Info("Failed to refresh image cache, retrying...")
		}
		log.Warn("Failed to refresh image cache. Is the portlayer server down?")
	}()

	return nil
}

func PortLayerClient() *client.PortLayer {
	return portLayerClient
}

func PortLayerServer() string {
	return portLayerServerAddr
}

func PortLayerName() string {
	return portLayerName
}

func ProductName() string {
	return productName
}

func ProductVersion() string {
	return productVersion
}

func VchConfig() *metadata.VirtualContainerHostConfigSpec {
	return vchConfig
}
