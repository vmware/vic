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

package simulator

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"github.com/vmware/vic/pkg/vsphere/simulator/vc"
	"golang.org/x/net/context"
)

// Model is used to populate a Model with an initial set of managed entities.
// This is a simple helper for tests running against a simulator, to populate an inventory
// with commonly used models.
type Model struct {
	Service *Service

	ServiceContent types.ServiceContent
	RootFolder     mo.Folder

	// Datacenter specifies the number of Datacenter entities to create
	Datacenter int

	// Portgroup specifies the number of DistributedVirtualPortgroup entities to create per Datacenter
	Portgroup int

	// Host specifies the number of standalone HostSystems entities to create per Datacenter
	Host int

	// Cluster specifies the number of ClusterComputeResource entities to create per Datacenter
	Cluster int

	// ClusterHost specifies the number of HostSystems entities to create within a Cluster
	ClusterHost int

	// Pool specifies the number of ResourcePool entities to create per Cluster
	Pool int

	// Datastore specifies the number of Datastore entities to create
	// Each Datastore will have temporary local file storage and will be mounted
	// on every HostSystem created by the ModelConfig
	Datastore int

	// Machine specifies the number of VirtualMachine entities to create per ResourcePool
	Machine int

	dirs []string
}

// ESX is the default Model for a standalone ESX instance
func ESX() *Model {
	return &Model{
		ServiceContent: esx.ServiceContent,
		RootFolder:     esx.RootFolder,
		Datastore:      1,
		Machine:        2,
	}
}

// VPX is the default Model for a vCenter instance
func VPX() *Model {
	return &Model{
		ServiceContent: vc.ServiceContent,
		RootFolder:     vc.RootFolder,
		Datacenter:     1,
		Host:           1,
		Cluster:        1,
		ClusterHost:    3,
		Datastore:      1,
		Machine:        2,
	}
}

func (*Model) fmtName(prefix string, num int) string {
	return fmt.Sprintf("%s%d", prefix, num)
}

// Create populates the Model with the given ModelConfig
func (m *Model) Create() error {
	m.Service = New(NewServiceInstance(m.ServiceContent, m.RootFolder))

	ctx := context.Background()
	client := m.Service.client
	root := object.NewRootFolder(client)

	// After all hosts are created, this var is used to mount the host datastores.
	var hosts []*object.HostSystem
	// We need to defer VM creation until after the datastores are created.
	var vms []func() error

	// addHost adds a cluster host or a stanalone host.
	addHost := func(name string, f func(types.HostConnectSpec) (*object.Task, error)) (*object.HostSystem, error) {
		spec := types.HostConnectSpec{
			HostName: name,
		}

		task, err := f(spec)
		if err != nil {
			return nil, err
		}

		info, err := task.WaitForResult(context.Background(), nil)
		if err != nil {
			return nil, err
		}

		host := object.NewHostSystem(client, info.Result.(types.ManagedObjectReference))
		hosts = append(hosts, host)

		return host, nil
	}

	// addMachine returns a func to create a VM.
	addMachine := func(prefix string, host *object.HostSystem, pool *object.ResourcePool, folders *object.DatacenterFolders) {
		f := func() error {
			for i := 0; i < m.Machine; i++ {
				name := m.fmtName(prefix+"_VM", i)

				config := types.VirtualMachineConfigSpec{
					Name:    name,
					GuestId: string(types.VirtualMachineGuestOsIdentifierOtherGuest),
					Files: &types.VirtualMachineFileInfo{
						VmPathName: "[LocalDS_0]",
					},
				}

				if pool == nil {
					pool, _ = host.ResourcePool(ctx)
				}

				task, err := folders.VmFolder.CreateVM(ctx, config, pool, host)
				if err != nil {
					return err
				}

				err = task.Wait(ctx)
				if err != nil {
					return err
				}
			}

			return nil
		}

		vms = append(vms, f)
	}

	for ndc := 0; ndc < m.Datacenter; ndc++ {
		dcName := m.fmtName("DC", ndc)

		dc, err := root.CreateDatacenter(ctx, dcName)
		if err != nil {
			return err
		}

		folders, err := dc.Folders(ctx)
		if err != nil {
			return err
		}

		for nhost := 0; nhost < m.Host; nhost++ {
			name := m.fmtName(dcName+"_H", nhost)

			host, err := addHost(name, func(spec types.HostConnectSpec) (*object.Task, error) {
				return folders.HostFolder.AddStandaloneHost(ctx, spec, true, nil, nil)
			})
			if err != nil {
				return err
			}

			addMachine(name, host, nil, folders)
		}

		for ncluster := 0; ncluster < m.Cluster; ncluster++ {
			clusterName := m.fmtName(dcName+"_C", ncluster)

			cluster, err := folders.HostFolder.CreateCluster(ctx, clusterName, types.ClusterConfigSpecEx{})
			if err != nil {
				return err
			}

			// TODO: create DistributedVirtualPortgroup for npg := 0; npg < m.Portgroup; npg++

			// TODO: create ResourcePool for npool := 0; npool < m.Pool; pool++

			for nhost := 0; nhost < m.ClusterHost; nhost++ {
				name := m.fmtName(clusterName+"_H", nhost)

				_, err := addHost(name, func(spec types.HostConnectSpec) (*object.Task, error) {
					return cluster.AddHost(ctx, spec, true, nil, nil)
				})
				if err != nil {
					return err
				}
			}

			pool, err := cluster.ResourcePool(ctx)
			if err != nil {
				return err
			}
			addMachine(clusterName+"_RP0", nil, pool, folders)
		}
	}

	if m.ServiceContent.RootFolder == esx.RootFolder.Reference() {
		// ESX model
		host := object.NewHostSystem(client, esx.HostSystem.Reference())
		hosts = append(hosts, host)

		dc := object.NewDatacenter(client, esx.Datacenter.Reference())
		folders, err := dc.Folders(ctx)
		if err != nil {
			return err
		}

		addMachine(host.Reference().Value, host, nil, folders)
	}

	for i := 0; i < m.Datastore; i++ {
		err := m.createLocalDatastore(m.fmtName("LocalDS_", i), hosts)
		if err != nil {
			return err
		}
	}

	for _, createVM := range vms {
		err := createVM()
		if err != nil {
			return err
		}
	}

	return nil
}

var tempDir = func() (string, error) {
	return ioutil.TempDir("", "govcsim-")
}

func (m *Model) createLocalDatastore(name string, hosts []*object.HostSystem) error {
	ctx := context.Background()
	dir, err := tempDir()
	if err != nil {
		return err
	}

	m.dirs = append(m.dirs, dir)

	for _, host := range hosts {
		dss, err := host.ConfigManager().DatastoreSystem(ctx)
		if err != nil {
			return err
		}

		_, err = dss.CreateLocalDatastore(ctx, name, dir)
		if err != nil {
			return err
		}
	}

	return nil
}

// Remove cleans up items created by the Model, such as local datastore directories
func (m *Model) Remove() {
	for _, dir := range m.dirs {
		_ = os.RemoveAll(dir)
	}
}
