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
package main

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
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

func (v *Validator) Validate(input *Data) (*metadata.VirtualContainerHostConfigSpec, error) {

	vchConfig := &metadata.VirtualContainerHostConfigSpec{}
	vchConfig.ApplianceSize.CPU.Limit = input.numCPUs
	vchConfig.ApplianceSize.Memory.Limit = input.memoryMB
	vchConfig.Name = input.displayName

	resources := strings.Split(input.computeResourcePath, "/")
	if len(resources) < 2 || resources[1] == "" {
		err := errors.Errorf("Could not determine datacenter from specified -compute path, %s", input.computeResourcePath)
		return nil, err
	}
	v.DatacenterName = resources[1]
	v.ClusterPath = strings.Split(input.computeResourcePath, "/Resources")[0]

	if v.ClusterPath == "" {
		err := errors.Errorf("Could not determine cluster from specified -compute path, %s", input.computeResourcePath)
		return nil, err
	}

	v.ResourcePoolPath = input.computeResourcePath
	v.ImageStorePath = fmt.Sprintf("/%s/datastore/%s", v.DatacenterName, input.imageDatastoreName)

	v.TargetPath = fmt.Sprintf("%s:%s@%s/sdk", input.user, *input.passwd, input.target)
	vchConfig.Target = v.TargetPath
	vchConfig.Insecure = input.insecure

	v.ExternalNetworkPath = fmt.Sprintf("/%s/network/%s", v.DatacenterName, input.externalNetworkName)
	v.BridgeNetworkPath = fmt.Sprintf("/%s/network/%s", v.DatacenterName, input.bridgeNetworkName)
	v.BridgeNetworkName = input.bridgeNetworkName
	if input.managementNetworkName != "" {
		v.ManagementNetworkPath = fmt.Sprintf("/%s/network/%s", v.DatacenterName, input.managementNetworkName)
		v.ManagementNetworkName = input.managementNetworkName
	}
	vchConfig.ImageStoreName = input.imageDatastoreName
	vchConfig.DatacenterName = v.DatacenterName
	vchConfig.ClusterPath = v.ClusterPath

	if err := v.validateConfiguration(input, vchConfig); err != nil {
		return nil, err
	}
	return vchConfig, nil
}

func (v *Validator) validateConfiguration(input *Data, vchConfig *metadata.VirtualContainerHostConfigSpec) error {
	log.Infof("Validating supplied configuration")

	var err error
	v.Context = context.TODO()
	sessionconfig := &session.Config{
		Service:        v.TargetPath,
		Insecure:       input.insecure,
		DatacenterPath: v.DatacenterName,
		ClusterPath:    v.ClusterPath,
		DatastorePath:  v.ImageStorePath,
		PoolPath:       v.ResourcePoolPath,
	}

	v.Session, err = session.NewSession(sessionconfig).Create(v.Context)
	if err != nil {
		log.Errorf("Failed to create session: %s", err)
		return err
	}
	_, err = v.Session.Connect(v.Context)
	if err != nil {
		return err
	}

	if _, err = v.Session.Populate(v.Context); err != nil {
		log.Errorf("Failed to get resources: %s", err)
		return err
	}

	// find the host(s) attached to given storage
	if _, err = v.Session.Datastore.AttachedClusterHosts(v.Context, v.Session.Cluster); err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		return err
	}

	if err = v.createBridgeNetwork(); err != nil && !input.force {
		return errors.Errorf("Creating bridge network failed with %s", err)
	}

	if err = v.setNetworks(vchConfig); err != nil {
		return errors.Errorf("Find networks failed with %s", err)
	}
	vchConfig.ComputeResources = append(vchConfig.ComputeResources, v.Session.Pool)
	vchConfig.ImageStores = append(vchConfig.ImageStores, v.ImageStorePath)
	//TODO: Add more configuration validation
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
			InventoryPath: path,
		}
	}
	return nil
}
