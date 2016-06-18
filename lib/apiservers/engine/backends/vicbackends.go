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
	"net"

	httptransport "github.com/go-swagger/go-swagger/httpkit/client"
	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
)

const (
	Imagec = "imagec"
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
	// update the image cache at startup
	if err := imageCache.Update(portLayerClient); err != nil {
		return fmt.Errorf("Error refreshing image cache: %s", err)
	}

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
