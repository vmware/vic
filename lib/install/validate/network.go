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
	"path"
	"strings"

	"context"

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
)

func (v *Validator) getEndpoint(ctx context.Context, conf *config.VirtualContainerHostConfigSpec, network data.NetworkConfig, epName, contNetName string, def bool, ns []net.IP) (*executor.NetworkEndpoint, error) {
	defer trace.End(trace.Begin(""))
	var gw net.IPNet
	var dest []net.IPNet
	var staticIP *net.IPNet

	if !network.Empty() {
		log.Debugf("Setting static IP for %q on port group %q", contNetName, network.Name)
		gw = network.Gateway
		dest = network.Destinations
		staticIP = &network.IP
	}

	moid, err := v.networkHelper(ctx, network.Name)
	if err != nil {
		return nil, err
	}

	e := &executor.NetworkEndpoint{
		Common: executor.Common{
			Name: epName,
		},
		Network: executor.ContainerNetwork{
			Common: executor.Common{
				Name: contNetName,
				ID:   moid,
			},
			Default:      def,
			Destinations: dest,
			Gateway:      gw,
			Nameservers:  ns,
		},
		IP: staticIP,
	}
	if staticIP != nil {
		e.Static = true
	}

	return e, nil
}

func (v *Validator) checkNetworkConflict(bridgeNetName, otherNetName, otherNetType string) {
	if bridgeNetName == otherNetName {
		v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - %s also uses %q", otherNetType, bridgeNetName))
	}
}

// portGroupConfig gets the input config for all networks
// for use in checking that the config is valid
func (v *Validator) portGroupConfig(input *data.Data, ips map[string][]data.NetworkConfig) {
	defer trace.End(trace.Begin(""))

	if input.ManagementNetwork.Name != "" {
		if !input.ManagementNetwork.Empty() {
			ips[input.ManagementNetwork.Name] = append(ips[input.ManagementNetwork.Name], input.ManagementNetwork)
		}
	}
	if input.ClientNetwork.Name != "" {
		if !input.ClientNetwork.Empty() {
			ips[input.ClientNetwork.Name] = append(ips[input.ClientNetwork.Name], input.ClientNetwork)
		}
	}
	if input.PublicNetwork.Name != "" {
		if !input.PublicNetwork.Empty() {
			ips[input.PublicNetwork.Name] = append(ips[input.PublicNetwork.Name], input.PublicNetwork)
		}
	}
}

// checkPortGroups checks that network config is valid
// enforce that networks that share a port group with public are configured via the pubic args
// prevent assigning > 1 static IP to the same port group
// warn if assigning addresses in the same subnet to > 1 port group
func (v *Validator) checkPortGroups(input *data.Data, ips map[string][]data.NetworkConfig) error {
	defer trace.End(trace.Begin(""))

	networks := make(map[string]string)

	shared := false
	// check for networks that share port group with public
	for nn, n := range map[string]*data.NetworkConfig{
		"client":     &input.ClientNetwork,
		"management": &input.ManagementNetwork,
	} {
		if n.Name == input.PublicNetwork.Name && !n.Empty() {
			log.Errorf("%s network shares port group with public network, but has static IP configuration", nn)
			log.Errorf("To resolve this, configure static IP for public network and assign %s network to same port group", nn)
			log.Error("The static IP will be automatically configured for networks sharing the port group")
			shared = true
		}
	}

	if shared {
		return fmt.Errorf("Static IP on network sharing port group with public network - Configuration ONLY allowed through public network options")
	}

	for pg, config := range ips {
		if len(config) > 1 {
			var msgIPs []string
			for _, v := range config {
				msgIPs = append(msgIPs, v.IP.IP.String())
			}
			log.Errorf("Port group %q is configured for networks with more than one static IP: %s", pg, msgIPs)
			log.Error("All VCH networks on the same port group must have the same IP address")
			log.Error("To resolve this, configure static IP for one network and assign other networks to same port group")
			log.Error("The static IP will be automatically configured for networks sharing the port group")
			return fmt.Errorf("Incorrect static IP configuration for networks on port group %q", pg)
		}

		// check if same subnet assigned to multiple portgroups - this can cause routing problems
		_, net, _ := net.ParseCIDR(config[0].IP.String())
		netAddr := net.String()
		if networks[netAddr] != "" {
			log.Warnf("Unsupported static IP configuration: Same subnet %q is assigned to multiple port groups %q and %q", netAddr, networks[netAddr], pg)
		} else {
			networks[netAddr] = pg
		}
	}

	return nil
}

// configureSharedPortGroups sets VCH static IP for networks that share a
// portgroup with another network that has a configured static IP
func (v *Validator) configureSharedPortGroups(input *data.Data, ips map[string][]data.NetworkConfig) error {
	defer trace.End(trace.Begin(""))

	// find other networks using same portgroup and copy the NetworkConfig to them
	for name, config := range ips {
		if len(config) != 1 {
			return fmt.Errorf("Failed to configure static IP for additional networks using port group %q", name)
		}
		log.Infof("Configuring static IP for additional networks using port group %q", name)
		if input.ClientNetwork.Name == name && input.ClientNetwork.Empty() {
			input.ClientNetwork = config[0]
		}
		if input.PublicNetwork.Name == name && input.PublicNetwork.Empty() {
			input.PublicNetwork = config[0]
		}
		if input.ManagementNetwork.Name == name && input.ManagementNetwork.Empty() {
			input.ManagementNetwork = config[0]
		}
	}

	return nil
}

func (v *Validator) network(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	var e *executor.NetworkEndpoint
	var err error

	// set default portgroup if user input not provided
	if input.ClientNetwork.Name == "" {
		input.ClientNetwork.Name = input.PublicNetwork.Name
	}
	if input.ManagementNetwork.Name == "" {
		input.ManagementNetwork.Name = input.ClientNetwork.Name
	}

	i := make(map[string][]data.NetworkConfig) // user configured IPs for portgroup
	v.portGroupConfig(input, i)

	err = v.checkPortGroups(input, i)
	v.NoteIssue(err)

	err = v.configureSharedPortGroups(input, i)
	v.NoteIssue(err)

	// client and management networks need to have at least one
	// routing destination if gateway was specified
	for nn, n := range map[string]*data.NetworkConfig{
		"client":     &input.ClientNetwork,
		"management": &input.ManagementNetwork,
	} {
		if n.Name == input.PublicNetwork.Name {
			// no Destinations required if sharing with PublicNetwork
			continue
		}
		if !ip.IsUnspecifiedIP(n.Gateway.IP) && len(n.Destinations) == 0 {
			v.NoteIssue(fmt.Errorf("%s network gateway specified without at least one routing destination", nn))
		}
	}

	// if static ip is specified for public network, gateway must be specified
	if !ip.IsUnspecifiedIP(input.PublicNetwork.IP.IP) && ip.IsUnspecifiedIP(input.PublicNetwork.Gateway.IP) {
		v.NoteIssue(errors.New("public network must have both static IP and gateway specified"))
	}

	// public network should not have any routing destinations specified
	// if a gateway was specified
	if !ip.IsUnspecifiedIP(input.PublicNetwork.Gateway.IP) && len(input.PublicNetwork.Destinations) > 0 {
		v.NoteIssue(errors.New("public network has the default route and must not have any routing destinations specified for gateway"))
	}

	// check if static IP on all networks and no user provided DNS servers
	specifiedDNS := len(input.DNS) > 0
	usingDHCP := ip.IsUnspecifiedIP(input.ClientNetwork.IP.IP) || ip.IsUnspecifiedIP(input.PublicNetwork.IP.IP) || ip.IsUnspecifiedIP(input.ManagementNetwork.IP.IP)

	if !usingDHCP && !specifiedDNS { // Set default DNS servers
		log.Debug("Setting default DNS servers 8.8.8.8 and 8.8.4.4")
		input.DNS = []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("8.8.4.4")}
	}

	// Public net
	// public network is default for appliance
	e, err = v.getEndpoint(ctx, conf, input.PublicNetwork, "public", "public", true, input.DNS)
	if err != nil {
		v.NoteIssue(fmt.Errorf("Error checking network for --public-network: %s", err))
		v.suggestNetwork("--public-network", true)
	}
	// Bridge network should be different than all other networks
	v.checkNetworkConflict(input.BridgeNetworkName, input.PublicNetwork.Name, "public")
	conf.AddNetwork(e)

	// Client net - defaults to connect to same portgroup as public
	e, err = v.getEndpoint(ctx, conf, input.ClientNetwork, "client", "client", false, input.DNS)
	if err != nil {
		v.NoteIssue(fmt.Errorf("Error checking network for --client-network: %s", err))
		v.suggestNetwork("--client-network", true)
	}
	v.checkNetworkConflict(input.BridgeNetworkName, input.ClientNetwork.Name, "client")
	conf.AddNetwork(e)

	// Management net - defaults to connect to the same portgroup as client
	e, err = v.getEndpoint(ctx, conf, input.ManagementNetwork, "", "management", false, input.DNS)
	if err != nil {
		v.NoteIssue(fmt.Errorf("Error checking network for --management-network: %s", err))
		v.suggestNetwork("--management-network", true)
	}
	v.checkNetworkConflict(input.BridgeNetworkName, input.ManagementNetwork.Name, "management")
	conf.AddNetwork(e)

	// Bridge net -
	//   vCenter: must exist and must be a DPG
	//   ESX: doesn't need to exist - we will create with default value
	//
	// for now we're hardcoded to "bridge" for the container host name
	conf.BridgeNetwork = "bridge"
	endpointMoref, err := v.dpgHelper(ctx, input.BridgeNetworkName)

	var bridgeID, netMoid string
	if err != nil {
		bridgeID = ""
		netMoid = ""
	} else {
		bridgeID = endpointMoref.String()
		netMoid = endpointMoref.String()
	}

	checkBridgeVDS := true
	if err != nil {
		if _, ok := err.(*find.NotFoundError); !ok || v.IsVC() {
			v.NoteIssue(fmt.Errorf("An existing distributed port group must be specified for bridge network on vCenter: %s", err))
			v.suggestNetwork("--bridge-network", false)
			checkBridgeVDS = false // prevent duplicate error output
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
		Static: true,
		IP:     &net.IPNet{IP: net.IPv4zero}, // static but managed externally
		Network: executor.ContainerNetwork{
			Common: executor.Common{
				Name: "bridge",
				ID:   netMoid,
			},
			Type: "bridge",
		},
	}
	// we need to have the bridge network identified as an available container network
	conf.AddContainerNetwork(&bridgeNet.Network)
	// we also need to have the appliance attached to the bridge network to allow
	// port forwarding
	conf.AddNetwork(bridgeNet)

	// make sure that the bridge IP pool is large enough for bridge networks
	err = v.checkBridgeIPRange(input.BridgeIPRange)
	if err != nil {
		v.NoteIssue(err)
	}
	conf.BridgeIPRange = input.BridgeIPRange

	log.Debug("Network configuration:")
	for net, val := range conf.ExecutorConfig.Networks {
		log.Debugf("\tNetwork: %s NetworkEndpoint: %v", net, val)
	}

	err = v.checkVDSMembership(ctx, endpointMoref, input.BridgeNetworkName)
	if err != nil && checkBridgeVDS {
		v.NoteIssue(fmt.Errorf("Unable to check hosts in vDS for %q: %s", input.BridgeNetworkName, err))
	}

	// add mapped networks (from --container-network)
	//   these should be a distributed port groups in vCenter
	suggestedMapped := false // only suggest mapped nets once
	for name, net := range input.MappedNetworks {
		checkMappedVDS := true
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
		if err != nil {
			v.NoteIssue(fmt.Errorf("Error adding container network %q: %s", name, err))
			checkMappedVDS = false
			if !suggestedMapped {
				v.suggestNetwork("--container-network", true)
				suggestedMapped = true
			}
		}
		mappedNet := &executor.ContainerNetwork{
			Common: executor.Common{
				Name: name,
				ID:   moref.String(),
			},
			Type:        "external",
			Gateway:     gw,
			Nameservers: dns,
			Pools:       pools,
		}
		if input.BridgeNetworkName == net {
			v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - %q also mapped as container network %q", input.BridgeNetworkName, name))
		}

		err = v.checkVDSMembership(ctx, moref, net)
		if err != nil && checkMappedVDS {
			v.NoteIssue(fmt.Errorf("Unable to check hosts in vDS for %q: %s", net, err))
		}

		conf.AddContainerNetwork(mappedNet)
	}
	v.nicNumbers(conf)

	conf.AsymmetricRouting = input.AsymmetricRouting
}

// nicNumbers will check vch appliance nic numbers. currently we don't support more than three nics for issue #1674.
// FIXME: this limitation should be removed after #1674 is fixed
func (v *Validator) nicNumbers(conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))
	nics := make(map[string]bool)
	for _, net := range conf.ExecutorConfig.Networks {
		nics[net.Network.ID] = true
	}
	if len(nics) > 3 {
		v.NoteIssue(fmt.Errorf("Four different networks including bridge network are not allowed at this time. At least two network roles of client network, management network and public network should share same network"))
	}
}

// generateBridgeName returns a name that can be used to create a switch/pg pair on ESX
func (v *Validator) generateBridgeName(ctx, input *data.Data, conf *config.VirtualContainerHostConfigSpec) string {
	defer trace.End(trace.Begin(""))

	return input.DisplayName
}

// checkBridgeIPRange verifies that the bridge network pool is large enough
// port layer currently defaults to /16 for bridge network so pool must be >= /16
func (v *Validator) checkBridgeIPRange(bridgeIPRange *net.IPNet) error {
	if bridgeIPRange == nil {
		return nil
	}
	ones, _ := bridgeIPRange.Mask.Size()
	if ones > 16 {
		return fmt.Errorf("Specified bridge network range is not large enough for the default bridge network size. --bridge-network-range must be /16 or larger network.")
	}
	return nil
}

// getNetwork gets a moref based on the network name
func (v *Validator) getNetwork(ctx context.Context, name string) (object.NetworkReference, error) {
	defer trace.End(trace.Begin(name))

	nets, err := v.Session.Finder.NetworkList(ctx, name)
	if err != nil {
		log.Debugf("no such network %q", name)
		// TODO: error message about no such match and how to get a network list
		// we return err directly here so we can check the type
		return nil, err
	}
	if len(nets) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return nil, errors.New("ambiguous network " + name)
	}
	return nets[0], nil
}

// networkHelper gets a moid based on the network name
func (v *Validator) networkHelper(ctx context.Context, name string) (string, error) {
	defer trace.End(trace.Begin(name))

	net, err := v.getNetwork(ctx, name)
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
	defer trace.End(trace.Begin(""))

	for _, h := range dvpHosts {
		if host == h {
			return true
		}
	}
	return false
}

// checkVDSMembership verifes all hosts in the vCenter are connected to the vDS
func (v *Validator) checkVDSMembership(ctx context.Context, network types.ManagedObjectReference, netName string) error {
	defer trace.End(trace.Begin(network.Value))

	var dvp mo.DistributedVirtualPortgroup
	var nonMembers []string

	if !v.IsVC() {
		return nil
	}

	if v.Session.Cluster == nil {
		return errors.New("Invalid cluster. Check --compute-resource")
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

// suggestNetwork suggests all networks
// incStdNets includes standard Networks in addition to DPGs
func (v *Validator) suggestNetwork(flag string, incStdNets bool) {
	defer trace.End(trace.Begin(flag))

	log.Infof("Suggesting valid networks for %s", flag)

	var validNets []string

	nets, err := v.Session.Finder.NetworkList(v.Context, "*")
	if err != nil {
		log.Errorf("Unable to list networks: %s", err)
		return
	}

	if len(nets) == 0 {
		log.Info("No networks found")
		return
	}

	for _, net := range nets {
		switch o := net.(type) {
		case *object.DistributedVirtualPortgroup:
			if v.isNetworkNameValid(o.InventoryPath, flag) {
				// Filter out DVS uplink
				if !v.isDVSUplink(net.Reference()) {
					validNets = append(validNets, path.Base(o.InventoryPath))
				}
			}
		case *object.Network:
			if incStdNets && v.isNetworkNameValid(o.InventoryPath, flag) {
				validNets = append(validNets, path.Base(o.InventoryPath))
			}
		}
	}

	if len(validNets) == 0 {
		log.Info("No valid networks found")
		return
	}

	log.Infof("Suggested values for %s:", flag)
	for _, n := range validNets {
		log.Infof("  %q", n)
	}
}

// isDVSUplink determines if the DVP is an uplink
func (v *Validator) isDVSUplink(ref types.ManagedObjectReference) bool {
	defer trace.End(trace.Begin(ref.Value))

	var dvp mo.DistributedVirtualPortgroup

	r := object.NewDistributedVirtualPortgroup(v.Session.Client.Client, ref)
	if err := r.Properties(v.Context, r.Reference(), []string{"tag"}, &dvp); err != nil {
		log.Errorf("Unable to check tags on %q: %s", ref, err)
		return false
	}
	for _, t := range dvp.Tag {
		if strings.Contains(t.Key, "UPLINKPG") {
			return true
		}
	}
	return false
}

// isNetworkNameValid determines if the network name in inventoryPath is
// not a reserved name
func (v *Validator) isNetworkNameValid(inventoryPath, flag string) bool {
	n := path.Base(inventoryPath)
	if flag != "--bridge-network" && n == "bridge" {
		return false
	}
	return true
}
