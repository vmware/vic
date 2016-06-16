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
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/cmd/vic-machine/data"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/session"

	"golang.org/x/net/context"
)

type Validator struct {
	TargetPath            string
	DatacenterName        string
	ClusterPath           string
	ResourcePoolPath      string
	ImageStorePath        string
	ExternalNetworkPath   string
	BridgeNetworkPath     string
	BridgeNetworkName     string
	ManagementNetworkPath string
	ManagementNetworkName string

	Session *session.Session
	Context context.Context
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(input *data.Data) (*metadata.VirtualContainerHostConfigSpec, error) {

	vchConfig := &metadata.VirtualContainerHostConfigSpec{}
	vchConfig.ApplianceSize.CPU.Limit = int64(input.NumCPUs)
	vchConfig.ApplianceSize.Memory.Limit = int64(input.MemoryMB)
	vchConfig.Name = input.DisplayName

	v.ParseComputeResourcePath(input.ComputeResourcePath) // Ignore error so we can suggest values

	v.TargetPath = input.URL.String()
	vchConfig.Target = v.TargetPath
	vchConfig.Insecure = input.Insecure

	if v.DatacenterName != "" {
		v.ImageStorePath = fmt.Sprintf("/%s/datastore/%s", v.DatacenterName, input.ImageDatastoreName)
		v.ExternalNetworkPath = fmt.Sprintf("/%s/network/%s", v.DatacenterName, input.ExternalNetworkName)
		v.BridgeNetworkPath = fmt.Sprintf("/%s/network/%s", v.DatacenterName, input.BridgeNetworkName)
		v.ManagementNetworkPath = fmt.Sprintf("/%s/network/%s", v.DatacenterName, input.ManagementNetworkName)
		vchConfig.DatacenterName = v.DatacenterName
	}

	v.BridgeNetworkName = input.BridgeNetworkName
	v.ManagementNetworkName = input.ManagementNetworkName
	vchConfig.ImageStoreName = input.ImageDatastoreName

	if v.ClusterPath != "" {
		vchConfig.ClusterPath = v.ClusterPath
	}

	if err := v.validateConfiguration(input, vchConfig); err != nil {
		return nil, err
	}
	return vchConfig, nil
}

func (v *Validator) ParseComputeResourcePath(rp string) error {
	resources := strings.Split(rp, "/")
	if len(resources) < 2 || resources[1] == "" {
		err := errors.Errorf("Could not determine datacenter from specified --compute-resource path: %s", rp)
		return err
	}
	v.DatacenterName = resources[1]
	v.ClusterPath = strings.Split(rp, "/Resources")[0]

	if v.ClusterPath == "" {
		err := errors.Errorf("Could not determine cluster from specified --compute-resource path: %s", rp)
		return err
	}
	v.ResourcePoolPath = rp
	return nil
}

func (v *Validator) GetResourcePool(input *data.Data) (*object.ResourcePool, error) {
	if err := v.ParseComputeResourcePath(input.ComputeResourcePath); err != nil {
		return nil, err
	}
	if _, err := v.getConnection(input); err != nil {
		return nil, err
	}

	if err := v.getResources(v.Context); err != nil {
		log.Errorf("Failed to get resources:\n%s", err)
		return nil, err
	}

	return v.Session.Pool, nil
}

func (v *Validator) getConnection(input *data.Data) (bool, error) {
	var err error
	v.Context = context.TODO()
	sessionconfig := &session.Config{
		Service:        input.URL.String(),
		Insecure:       input.Insecure,
		DatacenterPath: v.DatacenterName,
		ClusterPath:    v.ClusterPath,
		DatastorePath:  v.ImageStorePath,
		PoolPath:       v.ResourcePoolPath,
	}

	v.Session = session.NewSession(sessionconfig)

	if _, err = v.Session.Connect(v.Context); err != nil {
		return false, errors.Errorf("%s\nFailed to connect. Verify --target, --user, and --password", err.Error())
	}

	return true, nil
}
func (v *Validator) validateConfiguration(input *data.Data, vchConfig *metadata.VirtualContainerHostConfigSpec) error {
	log.Infof("Validating supplied configuration")
	var err error

	connected, err := v.getConnection(input) // Continue to validate if connected
	if connected == false {
		return err
	}

	if err = v.getResources(v.Context); err != nil {
		return err
	}
	// find the host(s) attached to given storage
	if _, err = v.Session.Datastore.AttachedClusterHosts(v.Context, v.Session.Cluster); err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		return err
	}

	if err = v.createBridgeNetwork(); err != nil && !input.Force {
		return errors.Errorf("Creating bridge network failed with %s", err)
	}

	if err = v.setNetworks(vchConfig); err != nil {
		return errors.Errorf("Find networks failed with %s", err)
	}
	vchConfig.Cluster = v.Session.Cluster.Reference()
	vchConfig.Datacenter = v.Session.Datacenter.Reference()
	vchConfig.ComputeResources = append(vchConfig.ComputeResources, v.Session.Pool.Reference())
	vchConfig.ImageStore = v.ImageStorePath

	//TODO: Add more configuration validation
	return nil
}

func (v *Validator) suggestComputeResource() {
	log.Info("Suggesting valid values for --compute-resource")

	numPools := v.suggestResourcePools(v.ClusterPath)
	if numPools == 0 {
		// ClusterPath not valid, suggest ALL the things!
		log.Info("Failed to find resource pool in the provided path. Showing all resource pools.")
		numPools = v.suggestAllResourcePools()
	}
	if numPools == 0 {
		log.Info("No resource pools found")
	}
}

func (v *Validator) suggestAllResourcePools() int {
	return v.suggestResourcePools("*")
}

func (v *Validator) suggestResourcePools(path string) int {
	var numPools int
	clusters, _ := v.Session.Finder.ComputeResourceList(v.Context, path)
	if clusters != nil {
		for _, c := range clusters {
			numPools += v.listResourcePools(c)
		}
	}
	return numPools
}

func (v *Validator) listResourcePools(c *object.ComputeResource) int {
	children := fmt.Sprintf("%s/Resources/*", c.InventoryPath)
	log.Infof("  %s/Resources", c.InventoryPath) // Suggest `/Resources`
	pools, _ := v.Session.Finder.ResourcePoolList(v.Context, children)
	if len(pools) > 0 {
		for _, p := range pools {
			log.Infof("    %s", p.InventoryPath)
		}
	}
	return len(pools)
}

func (v *Validator) getResources(ctx context.Context) error {
	var errs []string
	var err error

	finder := v.Session.Finder

	log.Debug("Check vSphere resources ...")
	if v.DatacenterName != "" {
		v.Session.Datacenter, err = finder.DatacenterOrDefault(ctx, v.DatacenterName)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failure finding dc (%s): %s", v.DatacenterName, err.Error()))
		} else {
			finder.SetDatacenter(v.Session.Datacenter)
			log.Debugf("Found dc: %s", v.DatacenterName)
		}
	}

	if v.ClusterPath != "" {
		v.Session.Cluster, err = finder.ComputeResourceOrDefault(ctx, v.ClusterPath)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failure finding cluster (%s): %s", v.ClusterPath, err.Error()))
		} else {
			log.Debugf("Found cluster: %s", v.ClusterPath)
		}
	}

	if v.ImageStorePath != "" {
		v.Session.Datastore, err = finder.DatastoreOrDefault(ctx, v.ImageStorePath)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failure finding ds (%s): %s", v.ImageStorePath, err.Error()))
		} else {
			log.Debugf("Found ds: %s", v.ImageStorePath)
		}
	}

	v.Session.Host, err = finder.HostSystemOrDefault(ctx, v.Session.HostPath)
	if err != nil {
		if _, ok := err.(*find.DefaultMultipleFoundError); !ok || !v.Session.IsVC() {
			errs = append(errs, fmt.Sprintf("Failure finding host (%s): %s", v.Session.HostPath, err.Error()))
		}
	} else {
		log.Debugf("Found host: %s", v.Session.HostPath)
	}

	v.Session.Pool, err = finder.ResourcePoolOrDefault(ctx, v.ResourcePoolPath)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Failure finding pool (%s): %s", v.ResourcePoolPath, err.Error()))
		if v.Session.Datacenter == nil {
			log.Info("Invalid datacenter. Unable to suggest valid --compute-resource")
		} else {
			v.suggestComputeResource()
		}
	} else {
		log.Debugf("Found pool: %s", v.ResourcePoolPath)
	}

	if len(errs) > 0 {
		log.Debugf("Error count populating vSphere cache: (%d)", len(errs))
		return errors.New(strings.Join(errs, "\n"))
	}
	log.Debug("vSphere resources populated...")
	return nil
}

func (v *Validator) getNetworkPath(net object.NetworkReference) (string, string, error) {
	switch t := net.(type) {
	case *object.DistributedVirtualPortgroup:
		return t.InventoryPath, t.Name(), nil
	case *object.Network:
		return t.InventoryPath, t.Name(), nil
	case *object.DistributedVirtualSwitch:
		return "", "", errors.Errorf("Distributed Virtual Switch is not acceptable, please change to Distributed Virtual Port Group")
	default:
		return "", "", errors.Errorf("Unknown network card type: %s", reflect.TypeOf(t))
	}
}

func (v *Validator) setNetworks(vchConfig *metadata.VirtualContainerHostConfigSpec) error {
	var path, name string
	vchConfig.Networks = make(map[string]*metadata.NetworkInfo)

	// bridge network
	network, err := v.Session.Finder.NetworkOrDefault(v.Context, v.BridgeNetworkPath)
	if err != nil {
		err = errors.Errorf("Failed to get bridge network: %s", err)
		return err
	}
	path, name, err = v.getNetworkPath(network)
	if err != nil {
		err = errors.Errorf("Failed to get bridge network path: %s", err)
		return err
	}
	v.BridgeNetworkName = name
	vchConfig.BridgeNetwork = name
	v.BridgeNetworkPath = path
	vchConfig.Networks["bridge"] = &metadata.NetworkInfo{
		PortGroup:     network,
		PortGroupName: name,
		PortGroupRef:  network.Reference(),
		InventoryPath: path,
	}

	// client network
	network, err = v.Session.Finder.NetworkOrDefault(v.Context, v.ExternalNetworkPath)
	if err != nil {
		err = errors.Errorf("Failed to get external network: %s", err)
		return err
	}
	path, name, err = v.getNetworkPath(network)
	if err != nil {
		err = errors.Errorf("Failed to get client network path: %s", err)
		return err
	}
	v.ExternalNetworkPath = path
	vchConfig.Networks["client"] = &metadata.NetworkInfo{
		PortGroup:     network,
		PortGroupName: name,
		PortGroupRef:  network.Reference(),
		InventoryPath: path,
	}

	// management network
	if v.ManagementNetworkName != "" {
		network, err = v.Session.Finder.Network(v.Context, v.ManagementNetworkPath)
		if err != nil {
			err = errors.Errorf("Failed to get management network: %s", err)
			return err
		}
		path, name, err = v.getNetworkPath(network)
		if err != nil {
			err = errors.Errorf("Failed to get management network path: %s", err)
			return err
		}
		vchConfig.Networks["management"] = &metadata.NetworkInfo{
			PortGroup:     network,
			PortGroupName: name,
			PortGroupRef:  network.Reference(),
			InventoryPath: path,
		}
	}
	return nil
}
