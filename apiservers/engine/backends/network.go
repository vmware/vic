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

	"github.com/docker/engine-api/types/network"
	"github.com/docker/libnetwork"
)

type Network struct {
	ProductName string
}

func (n *Network) NetworkControllerEnabled() bool {
	return false
}

func (n *Network) FindNetwork(idName string) (libnetwork.Network, error) {
	return nil, fmt.Errorf("%s does not implement network.FindNetwork", n.ProductName)
}

func (n *Network) GetNetworkByName(idName string) (libnetwork.Network, error) {
	return nil, fmt.Errorf("%s does not implement network.GetNetworkByName", n.ProductName)
}

func (n *Network) GetNetworksByID(partialID string) []libnetwork.Network {
	return make([]libnetwork.Network, 0, 0)
}

func (n *Network) GetAllNetworks() []libnetwork.Network {
	return make([]libnetwork.Network, 0, 0)
}

func (n *Network) CreateNetwork(name, driver string, ipam network.IPAM, options map[string]string, internal bool, enableIPv6 bool) (libnetwork.Network, error) {
	return nil, fmt.Errorf("%s does not implement network.CreateNetwork", n.ProductName)
}

func (n *Network) ConnectContainerToNetwork(containerName, networkName string, endpointConfig *network.EndpointSettings) error {
	return fmt.Errorf("%s does not implement network.ConnectContainerToNetwork", n.ProductName)
}

func (n *Network) DisconnectContainerFromNetwork(containerName string, network libnetwork.Network, force bool) error {
	return fmt.Errorf("%s does not implement network.DisconnectContainerFromNetwork", n.ProductName)
}

func (n *Network) DeleteNetwork(name string) error {
	return fmt.Errorf("%s does not implement network.DeleteNetwork", n.ProductName)
}
