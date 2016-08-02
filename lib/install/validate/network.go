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

package validate

import (
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
	"golang.org/x/net/context"
)

func (v *Validator) network(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// External net
	extMoref, err := v.networkHelper(ctx, input.ExternalNetworkName)
	v.NoteIssue(err)
	conf.AddNetwork(&executor.NetworkEndpoint{
		Common: executor.Common{
			Name: "external",
		},
		Network: executor.ContainerNetwork{
			Common: executor.Common{
				Name: "external",
				ID:   extMoref,
			},
			Default: true, // external network is default for appliance
		},
	})
	// Bridge network should be different with all other networks
	if input.BridgeNetworkName == input.ExternalNetworkName {
		v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - external also uses %q", input.BridgeNetworkName))
	}

	// Client net
	if input.ClientNetworkName == "" {
		input.ClientNetworkName = input.ExternalNetworkName
	}
	clientMoref, err := v.networkHelper(ctx, input.ClientNetworkName)
	v.NoteIssue(err)
	conf.AddNetwork(&executor.NetworkEndpoint{
		Common: executor.Common{
			Name: "client",
		},
		Network: executor.ContainerNetwork{
			Common: executor.Common{
				Name: "client",
				ID:   clientMoref,
			},
		},
	})
	if input.BridgeNetworkName == input.ClientNetworkName {
		v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - client also uses %q", input.BridgeNetworkName))
	}

	// Management net
	if input.ManagementNetworkName == "" {
		input.ManagementNetworkName = input.ClientNetworkName
	}
	managementMoref, err := v.networkHelper(ctx, input.ManagementNetworkName)
	v.NoteIssue(err)
	conf.AddNetwork(&executor.NetworkEndpoint{
		Network: executor.ContainerNetwork{
			Common: executor.Common{
				Name: "management",
				ID:   managementMoref,
			},
		},
	})
	if input.BridgeNetworkName == input.ManagementNetworkName {
		v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - management also uses %q", input.BridgeNetworkName))
	}

	// Bridge net -
	//   vCenter: must exist and must be a DPG
	//   ESX: doesn't need to exist
	//
	// for now we're hardcoded to "bridge" for the container host name
	conf.BridgeNetwork = "bridge"
	endpointMoref, err := v.dpgHelper(ctx, input.BridgeNetworkName)
	var netMoid string
	bridgeID := endpointMoref.String()
	if endpointMoref.Type == "" && endpointMoref.Value == "" {
		netMoid = ""
		bridgeID = ""
	} else {
		netMoid = endpointMoref.String()
	}
	if err != nil {
		if _, ok := err.(*find.NotFoundError); !ok || v.IsVC() {
			v.NoteIssue(fmt.Errorf("An existing distributed port group must be specified for bridge network on vCenter: %s", err))
		}

		// this allows the dispatcher to create the network with corresponding name
		// if BridgeNetworkName doesn't already exist then we set the ContainerNetwork
		// ID to the name, but leaving the NetworkEndpoint moref as ""
		netMoid = input.BridgeNetworkName
	}

	bridgeNet := &executor.NetworkEndpoint{
		Common: executor.Common{
			Name: "bridge",
			ID:   bridgeID,
		},
		Static: &net.IPNet{IP: net.IPv4zero}, // static but managed externally
		Network: executor.ContainerNetwork{
			Common: executor.Common{
				Name: "bridge",
				ID:   netMoid,
			},
		},
	}
	// we  need to have the bridge network identified as an available container
	// network
	conf.AddContainerNetwork(&bridgeNet.Network)
	// we also need to have the appliance attached to the bridge network to allow
	// port forwarding
	conf.AddNetwork(bridgeNet)

	err = v.checkVDSMembership(ctx, endpointMoref, input.BridgeNetworkName)
	if err != nil {
		v.NoteIssue(fmt.Errorf("Unable to check hosts in vDS for %q: %s", input.BridgeNetworkName, err))
	}

	// add mapped networks
	//   these should be a distributed port groups in vCenter
	for name, net := range input.MappedNetworks {
		// "bridge" is reserved
		if name == "bridge" {
			v.NoteIssue(fmt.Errorf("Cannot use reserved name \"bridge\" for container network"))
			continue
		}

		gw := input.MappedNetworksGateways[name]
		pools := input.MappedNetworksIPRanges[name]
		dns := input.MappedNetworksDNS[name]
		if len(pools) != 0 && ip.IsUnspecifiedSubnet(&gw) {
			v.NoteIssue(fmt.Errorf("IP range specified without gateway for container network %q", name))
			continue
		}

		if !ip.IsUnspecifiedSubnet(&gw) && !ip.IsRoutableIP(gw.IP, &gw) {
			v.NoteIssue(fmt.Errorf("Gateway %s is not a routable address", gw.IP))
			continue
		}

		err = nil
		// verify ip ranges are within subnet,
		// and don't overlap with each other
		for i, r := range pools {
			if !gw.Contains(r.FirstIP) || !gw.Contains(r.LastIP) {
				err = fmt.Errorf("IP range %q is not in subnet %q", r, gw)
				break
			}

			for _, r2 := range pools[i+1:] {
				if r2.Overlaps(r) {
					err = fmt.Errorf("Overlapping ip ranges: %q %q", r2, r)
					break
				}
			}

			if err != nil {
				break
			}
		}

		if err != nil {
			v.NoteIssue(err)
			continue
		}

		moref, err := v.dpgHelper(ctx, net)
		v.NoteIssue(err)
		mappedNet := &executor.ContainerNetwork{
			Common: executor.Common{
				Name: name,
				ID:   moref.String(),
			},
			Gateway:     gw,
			Nameservers: dns,
			Pools:       pools,
		}
		if input.BridgeNetworkName == net {
			v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - %q also mapped as container network %q", input.BridgeNetworkName, name))
		}

		err = v.checkVDSMembership(ctx, moref, net)
		if err != nil {
			v.NoteIssue(fmt.Errorf("Unable to check hosts in vDS for %q: %s", net, err))
		}

		conf.AddContainerNetwork(mappedNet)
	}
}

// generateBridgeName returns a name that can be used to create a switch/pg pair on ESX
func (v *Validator) generateBridgeName(ctx, input *data.Data, conf *config.VirtualContainerHostConfigSpec) string {
	defer trace.End(trace.Begin(""))

	return input.DisplayName
}

func (v *Validator) getNetwork(ctx context.Context, path string) (object.NetworkReference, error) {
	nets, err := v.Session.Finder.NetworkList(ctx, path)
	if err != nil {
		log.Debugf("no such network %q", path)
		// TODO: error message about no such match and how to get a network list
		// we return err directly here so we can check the type
		return nil, err
	}
	if len(nets) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return nil, errors.New("ambiguous network " + path)
	}
	return nets[0], nil
}

func (v *Validator) networkHelper(ctx context.Context, path string) (string, error) {
	defer trace.End(trace.Begin(path))

	net, err := v.getNetwork(ctx, path)
	if err != nil {
		return "", err
	}
	moref := net.Reference()
	return moref.String(), nil
}

func (v *Validator) dpgMorefHelper(ctx context.Context, ref string) (string, error) {
	defer trace.End(trace.Begin(ref))

	moref := new(types.ManagedObjectReference)
	ok := moref.FromString(ref)
	if !ok {
		// TODO: error message about no such match and how to get a network list
		return "", errors.New("could not restore serialized managed object reference: " + ref)
	}

	net, err := v.Session.Finder.ObjectReference(ctx, *moref)
	if err != nil {
		// TODO: error message about no such match and how to get a network list
		return "", errors.New("unable to locate network from moref: " + ref)
	}

	// ensure that the type of the network is a Distributed Port Group if the target is a vCenter
	// if it's not then any network suffices
	if v.IsVC() {
		_, dpg := net.(*object.DistributedVirtualPortgroup)
		if !dpg {
			return "", fmt.Errorf("%q is not a Distributed Port Group", ref)
		}
	}

	return ref, nil
}

func (v *Validator) dpgHelper(ctx context.Context, path string) (types.ManagedObjectReference, error) {
	defer trace.End(trace.Begin(path))

	net, err := v.getNetwork(ctx, path)
	if err != nil {
		return types.ManagedObjectReference{}, err
	}

	// ensure that the type of the network is a Distributed Port Group if the target is a vCenter
	// if it's not then any network suffices
	if v.IsVC() {
		_, dpg := net.(*object.DistributedVirtualPortgroup)
		if !dpg {
			return types.ManagedObjectReference{}, fmt.Errorf("%q is not a Distributed Port Group", path)
		}
	}

	return net.Reference(), nil
}

// inDVP checks if the host is in the distributed virtual portgroup (dvpHosts)
func (v *Validator) inDVP(host types.ManagedObjectReference, dvpHosts []types.ManagedObjectReference) bool {
	for _, h := range dvpHosts {
		if host == h {
			return true
		}
	}
	return false
}

// checkVDSMembership verifes all hosts in the vCenter are connected to the vDS
func (v *Validator) checkVDSMembership(ctx context.Context, network types.ManagedObjectReference, netName string) error {
	var dvp mo.DistributedVirtualPortgroup
	var nonMembers []string

	if !v.IsVC() {
		return nil
	}

	clusterHosts, err := v.Session.Cluster.Hosts(ctx)
	if err != nil {
		return err
	}

	r := object.NewDistributedVirtualPortgroup(v.Session.Client.Client, network)
	if err := r.Properties(ctx, r.Reference(), []string{"name", "host"}, &dvp); err != nil {
		return err
	}

	for _, h := range clusterHosts {
		if !v.inDVP(h.Reference(), dvp.Host) {
			nonMembers = append(nonMembers, h.InventoryPath)
		}
	}

	if len(nonMembers) > 0 {
		log.Errorf("vDS configuration incorrect on %q. All cluster hosts must be in the vDS.", netName)
		log.Errorf("  %q is missing hosts:", netName)
		for _, hs := range nonMembers {
			log.Errorf("    %q", hs)
		}

		errMsg := fmt.Sprintf("All cluster hosts must be in the vDS. %q is missing hosts: %s", netName, nonMembers)
		v.NoteIssue(errors.New(errMsg))
	} else {
		log.Infof("vDS configuration OK on %q", netName)
	}
	return nil
}
