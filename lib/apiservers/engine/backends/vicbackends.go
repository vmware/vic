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
)

const (
	Imagec           = "imagec"
	Retries          = 5
	RetryTimeSeconds = 2
)

var (
	portLayerClient     *client.PortLayer
	portLayerServerAddr string

	imageCache *cache.ImageCache
)

func Init(portLayerAddr string) error {
	_, _, err := net.SplitHostPort(portLayerAddr)
	if err != nil {
		return err
	}

	t := httptransport.New(portLayerAddr, "/", []string{"http"})
	portLayerClient = client.New(t, nil)
	portLayerServerAddr = portLayerAddr

	imageCache = cache.NewImageCache()

	// attempt to update the image cache at startup
	log.Info("Refreshing image cache...")
	go func() {
		for i := 0; i < Retries; i++ {

			// initial pause to wait for the portlayer to come up
			time.Sleep(RetryTimeSeconds * time.Second)

			if err := imageCache.Update(portLayerClient); err == nil {
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

func ImageCache() *cache.ImageCache {
	return imageCache
}
