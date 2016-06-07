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
	"sync"

	"net/http"

	log "github.com/Sirupsen/logrus"
	derr "github.com/docker/docker/errors"
	apinet "github.com/docker/engine-api/types/network"
	"github.com/docker/libnetwork"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
)

type Network struct {
	ProductName string
}

func (n *Network) NetworkControllerEnabled() bool {
	return false
}

func (n *Network) FindNetwork(idName string) (libnetwork.Network, error) {
	ok, err := PortLayerClient().Scopes.List(scopes.NewListParams().WithIDName(idName))
	if err != nil {
		switch err := err.(type) {
		case *scopes.ListNotFound:
			return nil, derr.NewRequestNotFoundError(fmt.Errorf("network %s not found", idName))

		case *scopes.ListDefault:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

		default:

			return nil, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}
	}

	return &network{cfg: ok.Payload[0]}, nil
}

func (n *Network) GetNetworkByName(idName string) (libnetwork.Network, error) {
	ok, err := PortLayerClient().Scopes.List(scopes.NewListParams().WithIDName(idName))
	if err != nil {
		switch err := err.(type) {
		case *scopes.ListNotFound:
			return nil, nil

		case *scopes.ListDefault:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

		default:

			return nil, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}
	}

	return &network{cfg: ok.Payload[0]}, nil
}

func (n *Network) GetNetworksByID(partialID string) []libnetwork.Network {
	ok, err := PortLayerClient().Scopes.List(scopes.NewListParams().WithIDName(partialID))
	if err != nil {
		return nil
	}

	nets := make([]libnetwork.Network, len(ok.Payload))
	for i, cfg := range ok.Payload {
		nets[i] = &network{cfg: cfg}
	}

	return nets
}

func (n *Network) GetAllNetworks() []libnetwork.Network {
	ok, err := PortLayerClient().Scopes.ListAll(scopes.NewListAllParams())
	if err != nil {
		return nil
	}

	nets := make([]libnetwork.Network, len(ok.Payload))
	for i, cfg := range ok.Payload {
		nets[i] = &network{cfg: cfg}
	}

	return nets
}

func (n *Network) CreateNetwork(name, driver string, ipam apinet.IPAM, options map[string]string, labels map[string]string, internal bool, enableIPv6 bool) (libnetwork.Network, error) {
	if len(ipam.Config) > 1 {
		return nil, fmt.Errorf("at most one ipam config supported")
	}

	var gateway, subnet *string
	var pools []string
	if len(ipam.Config) > 0 {
		if ipam.Config[0].Gateway != "" {
			gateway = new(string)
			*gateway = ipam.Config[0].Gateway
		}

		if ipam.Config[0].Subnet != "" {
			subnet = new(string)
			*subnet = ipam.Config[0].Subnet
		}

		if ipam.Config[0].IPRange != "" {
			pools = append(pools, ipam.Config[0].IPRange)
		}
	}

	if driver == "" {
		driver = "bridge"
	}

	cfg := &models.ScopeConfig{
		Gateway:   gateway,
		Name:      name,
		ScopeType: driver,
		Subnet:    subnet,
		IPAM:      pools,
	}

	created, err := PortLayerClient().Scopes.CreateScope(scopes.NewCreateScopeParams().WithConfig(cfg))
	if err != nil {
		switch err := err.(type) {
		case *scopes.CreateScopeConflict:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf("network %s already exists", name), http.StatusConflict)

		case *scopes.CreateScopeDefault:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

		default:
			return nil, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}
	}

	return &network{cfg: created.Payload}, nil
}

func (n *Network) ConnectContainerToNetwork(containerName, networkName string, endpointConfig *apinet.EndpointSettings) error {
	client := PortLayerClient()
	getRes, err := client.Containers.Get(containers.NewGetParams().WithID(containerName))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf(err.Payload.Message))

		case *containers.GetDefault:
			return derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

		default:
			return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}
	}

	h := getRes.Payload
	nc := &models.NetworkConfig{NetworkName: networkName}
	if endpointConfig != nil && endpointConfig.IPAMConfig != nil && endpointConfig.IPAMConfig.IPv4Address != "" {
		nc.Address = &endpointConfig.IPAMConfig.IPv4Address
	}

	addConRes, err := client.Scopes.AddContainer(scopes.NewAddContainerParams().WithHandle(h).WithNetworkConfig(nc))
	if err != nil {
		switch err := err.(type) {
		case *scopes.AddContainerNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf(err.Payload.Message))

		case *scopes.AddContainerDefault:
			return derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

		default:
			return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}
	}

	defer func() {
		if err != nil {
			if _, err2 := client.Scopes.RemoveContainer(scopes.NewRemoveContainerParams().WithHandle(h).WithScope(nc.NetworkName)); err2 != nil {
				log.Warnf("failed add container to network rollback: %s", err2)
			}
		}
	}()

	h = addConRes.Payload

	// only bind if the container is running
	// get the state of the container
	getStateRes, err := client.Containers.GetState(containers.NewGetStateParams().WithHandle(h))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetStateNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf(err.Payload.Message))

		case *containers.GetStateDefault:
			return derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

		default:
			return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}
	}

	h = getStateRes.Payload.Handle
	if getStateRes.Payload.State == "RUNNING" {
		bindRes, err := client.Scopes.BindContainer(scopes.NewBindContainerParams().WithHandle(h))
		if err != nil {
			switch err := err.(type) {
			case *scopes.BindContainerNotFound:
				return derr.NewRequestNotFoundError(fmt.Errorf(err.Payload.Message))

			case *scopes.BindContainerDefault:
				return derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

			default:
				return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
			}
		}

		defer func() {
			if err == nil {
				return
			}
			if _, err2 := client.Scopes.UnbindContainer(scopes.NewUnbindContainerParams().WithHandle(h)); err2 != nil {
				log.Warnf("failed bind container rollback: %s", err2)
			}
		}()
		h = bindRes.Payload
	}

	// commit handle
	_, err = client.Containers.Commit(containers.NewCommitParams().WithHandle(h))
	if err != nil {
		switch err := err.(type) {
		case *containers.CommitNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf(err.Payload.Message))

		case *containers.CommitDefault:
			return derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

		default:
			return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}
	}

	return nil
}

func (n *Network) DisconnectContainerFromNetwork(containerName string, network libnetwork.Network, force bool) error {
	return fmt.Errorf("%s does not implement network.DisconnectContainerFromNetwork", n.ProductName)
}

func (n *Network) DeleteNetwork(name string) error {
	client := PortLayerClient()

	if _, err := client.Scopes.DeleteScope(scopes.NewDeleteScopeParams().WithIDName(name)); err != nil {
		if _, ok := err.(*scopes.DeleteScopeNotFound); ok {
			return derr.NewRequestNotFoundError(err)
		}

		return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
	}

	return nil
}

// network implements the libnetwork.Network and libnetwork.NetworkInfo interfaces
type network struct {
	sync.Mutex

	cfg *models.ScopeConfig
}

// A user chosen name for this network.
func (n *network) Name() string {
	return n.cfg.Name
}

// A system generated id for this network.
func (n *network) ID() string {
	return *n.cfg.ID
}

// The type of network, which corresponds to its managing driver.
func (n *network) Type() string {
	return n.cfg.ScopeType
}

// Create a new endpoint to this network symbolically identified by the
// specified unique name. The options parameter carry driver specific options.
func (n *network) CreateEndpoint(name string, options ...libnetwork.EndpointOption) (libnetwork.Endpoint, error) {
	return nil, fmt.Errorf("not implemented")
}

// Delete the network.
func (n *network) Delete() error {
	return fmt.Errorf("not implemented")
}

// Endpoints returns the list of Endpoint(s) in this network.
func (n *network) Endpoints() []libnetwork.Endpoint {
	return nil

}

// WalkEndpoints uses the provided function to walk the Endpoints
func (n *network) WalkEndpoints(walker libnetwork.EndpointWalker) {
	return
}

// EndpointByName returns the Endpoint which has the passed name. If not found, the error ErrNoSuchEndpoint is returned.
func (n *network) EndpointByName(name string) (libnetwork.Endpoint, error) {
	return nil, fmt.Errorf("not implemented")

}

// EndpointByID returns the Endpoint which has the passed id. If not found, the error ErrNoSuchEndpoint is returned.
func (n *network) EndpointByID(id string) (libnetwork.Endpoint, error) {
	return nil, fmt.Errorf("not implemented")
}

// Return certain operational data belonging to this network
func (n *network) Info() libnetwork.NetworkInfo {
	return n.Info()
}

func (n *network) IpamConfig() (string, map[string]string, []*libnetwork.IpamConf, []*libnetwork.IpamConf) {
	n.Lock()
	defer n.Unlock()

	confs := make([]*libnetwork.IpamConf, len(n.cfg.IPAM))
	for j, i := range n.cfg.IPAM {
		conf := &libnetwork.IpamConf{
			PreferredPool: *n.cfg.Subnet,
			Gateway:       "",
		}

		if i != *n.cfg.Subnet {
			conf.SubPool = i
		}

		if n.cfg.Gateway != nil {
			conf.Gateway = *n.cfg.Gateway
		}

		confs[j] = conf
	}

	return "", make(map[string]string), confs, nil
}

func (n *network) IpamInfo() ([]*libnetwork.IpamInfo, []*libnetwork.IpamInfo) {
	n.Lock()
	defer n.Unlock()

	var infos []*libnetwork.IpamInfo
	for _, i := range n.cfg.IPAM {
		_, pool, err := net.ParseCIDR(i)
		if err != nil {
			continue
		}

		info := &libnetwork.IpamInfo{
			Meta: make(map[string]string),
		}

		info.Pool = pool
		if n.cfg.Gateway != nil {
			info.Gateway = &net.IPNet{IP: net.ParseIP(*n.cfg.Gateway), Mask: net.CIDRMask(32, 32)}
		}

		info.AuxAddresses = make(map[string]*net.IPNet)
		infos = append(infos, info)
	}

	return infos, nil
}

func (n *network) DriverOptions() map[string]string {
	return make(map[string]string)
}

func (n *network) Scope() string {
	return ""
}

func (n *network) IPv6Enabled() bool {
	return false
}

func (n *network) Internal() bool {
	return false
}
