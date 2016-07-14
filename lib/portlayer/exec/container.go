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

package exec

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
	"golang.org/x/net/context"
)

var containers map[ID]*Container
var containersLock sync.Mutex

func init() {
	containers = make(map[ID]*Container)
}

type State int

const (
	StateRunning = iota + 1
	StateStopped

	propertyCollectorTimeout = 3 * time.Minute
)

type Container struct {
	sync.Mutex

	ExecConfig *metadata.ExecutorConfig
	State      State
	// friendly description of state
	Status string

	VMUnsharedDisk int64

	vm *vm.VirtualMachine
}

func NewContainer(id ID) *Handle {
	con := &Container{
		ExecConfig: &metadata.ExecutorConfig{},
		State:      StateStopped,
	}
	con.ExecConfig.ID = id.String()

	containersLock.Lock()
	containers[id] = con
	containersLock.Unlock()

	return con.newHandle()
}

func GetContainer(id ID) *Handle {
	containersLock.Lock()
	defer containersLock.Unlock()

	if con, ok := containers[id]; ok {
		return con.newHandle()
	}

	return nil
}

func (c *Container) newHandle() *Handle {
	return newHandle(c)
}

func (c *Container) cacheExecConfig(ec *metadata.ExecutorConfig) {
	c.Lock()
	defer c.Unlock()

	c.ExecConfig = ec
}

func (c *Container) Commit(ctx context.Context, sess *session.Session, h *Handle) error {
	defer trace.End(trace.Begin("Committing handle"))

	c.Lock()
	defer c.Unlock()

	if c.vm == nil {
		// the only permissible operation is to create a VM
		if h.Spec == nil {
			return fmt.Errorf("only create operations can be committed without an existing VM")
		}

		// Find the Virtual Machine folder that we use
		folders, err := sess.Datacenter.Folders(ctx)
		if err != nil {
			return err
		}
		parent := folders.VmFolder

		// Create the vm
		res, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
			return parent.CreateVM(ctx, *h.Spec.Spec(), Config.ResourcePool, nil)
		})
		if err != nil {
			return err
		}

		c.vm = vm.NewVirtualMachine(ctx, sess, res.Result.(types.ManagedObjectReference))

		// clear the spec as we've acted on it
		h.Spec = nil
	}

	// if we're stopping the VM, do so before the reconfigure to preserve the extraconfig
	if h.State != nil && *h.State == StateStopped {
		// stop the container
		if err := h.Container.stop(ctx); err != nil {
			return err
		}

		c.State = *h.State
	}

	if h.Spec != nil {
		// FIXME: add check that the VM is powered off - it should be, but this will destroy the
		// extraconfig if it's not.

		s := h.Spec.Spec()
		_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
			return c.vm.Reconfigure(ctx, *s)
		})

		if err != nil {
			return err
		}
	}

	if h.State != nil && *h.State == StateRunning {
		// start the container
		if err := h.Container.start(ctx); err != nil {
			return err
		}

		c.State = *h.State
	}

	c.ExecConfig = &h.ExecConfig

	return nil
}

// Start starts a container vm with the given params
func (c *Container) start(ctx context.Context) error {
	defer trace.End(trace.Begin("Container.start"))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}

	// Power on
	_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return c.vm.PowerOn(ctx)
	})
	if err != nil {
		return err
	}

	// guestinfo key that we want to wait for
	key := fmt.Sprintf("guestinfo..sessions|%s.started", c.ExecConfig.ID)
	var detail string

	// Wait some before giving up...
	ctx, cancel := context.WithTimeout(ctx, propertyCollectorTimeout)
	defer cancel()

	detail, err = c.vm.WaitForKeyInExtraConfig(ctx, key)
	if err != nil {
		return fmt.Errorf("unable to wait for process launch status: %s", err.Error())
	}

	if detail != "true" {
		return errors.New(detail)
	}

	return nil
}

func (c *Container) stop(ctx context.Context) error {
	defer trace.End(trace.Begin("Container.stop"))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}

	//TODO: make the shutdown much cleaner, right now we just pull the plug on the vm.(may need corresponding work in tether.)

	// Power off
	_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return c.vm.PowerOff(ctx)
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Container) Remove(ctx context.Context, sess *session.Session) error {
	c.Lock()
	defer c.Unlock()

	if c.vm == nil {
		return fmt.Errorf("VM has already been removed")
	}

	// get the folder the VM is in
	url, err := c.vm.DSPath(ctx)
	if err != nil {
		log.Errorf("Failed to get datastore path for %s: %s", c.ExecConfig.ID, err)
		return err
	}
	// FIXME: was expecting to find a utility function to convert to/from datastore/url given
	// how widely it's used but couldn't - will ask around.
	dsPath := fmt.Sprintf("[%s] %s", url.Host, url.Path)

	//removes the vm from vsphere, but detaches the disks first
	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return c.vm.DeleteExceptDisks(ctx)
	})
	if err != nil {
		return err
	}

	// remove from datastore
	fm := object.NewFileManager(c.vm.Client.Client)

	if _, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return fm.DeleteDatastoreFile(ctx, dsPath, sess.Datacenter)
	}); err != nil {
		log.Debugf("Failed to delete %s, %s", dsPath, err)
	}

	//removes container from map
	containersLock.Lock()
	delete(containers, ParseID(c.ExecConfig.ID))
	containersLock.Unlock()

	return nil
}

// Grab the info for the requested container
// TODO:  Possibly change so that handler requests a handle to the
// container and if it's not present then search and return a handle
func ContainerInfo(ctx context.Context, sess *session.Session, containerID ID) (*Container, error) {
	var cvm *Container
	// first  lets see if we have it in the cache
	container := containers[containerID]
	// if we missed search for it...
	if container == nil {
		// search
		vm, err := childVM(ctx, sess, containerID.String())
		if err != nil || vm == nil {
			log.Debugf("ContainerInfo failed to find childVM: %s", err.Error())
			return cvm, fmt.Errorf("Container Not Found: %s", containerID)
		}
		container = &Container{vm: vm}
	}

	// get properties for specific containerVMs
	// we must do this since we don't know that we have a valid state
	// This will be refactored when event streaming hits
	vms, err := populateVMAttributes(ctx, sess, []types.ManagedObjectReference{container.vm.Reference()})
	if err != nil {
		return cvm, err
	}
	// convert the VMs to container objects -- include
	// powered off vms
	all := true
	cc := convertInfraContainers(vms, &all)

	switch len(cc) {
	case 0:
		// we found a vm, but it doesn't appear to be a container VM
		return cvm, fmt.Errorf("%s does not appear to be a container", containerID)
	case 1:
		// we have a winner
		cvm = &cc[0]
	default:
		// we manged to find multiple vms
		return cvm, fmt.Errorf("multiple containers named %s found", containerID)
	}

	return cvm, nil
}

// return a list of container attributes
func List(ctx context.Context, sess *session.Session, all *bool) ([]Container, error) {

	// for now we'll go to the infrastructure
	// future iteration will utilize cache & event stream
	moVMs, err := infraContainers(ctx, sess)
	if err != nil {
		return nil, err
	}
	// convert to container
	containers := convertInfraContainers(moVMs, all)
	return containers, nil
}

// get the containerVMs from infrastructure for this resource pool
func infraContainers(ctx context.Context, sess *session.Session) ([]mo.VirtualMachine, error) {
	var rp mo.ResourcePool

	// popluate the vm property of the vch resource pool
	if err := Config.ResourcePool.Properties(ctx, Config.ResourcePool.Reference(), []string{"vm"}, &rp); err != nil {
		name := Config.ResourcePool.Name()
		log.Errorf("List failed to get %s resource pool child vms: %s", name, err)
		return nil, err
	}
	vms, err := populateVMAttributes(ctx, sess, rp.Vm)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

// find the childVM for this resource pool by name
func childVM(ctx context.Context, sess *session.Session, name string) (*vm.VirtualMachine, error) {
	searchIndex := object.NewSearchIndex(sess.Client.Client)
	child, err := searchIndex.FindChild(ctx, Config.ResourcePool.Reference(), name)
	if err != nil {
		return nil, fmt.Errorf("Unable to find container(%s): %s", name, err.Error())
	}
	if child == nil {
		return nil, fmt.Errorf("Unable to find container %s", name)
	}
	// instantiate the vm object
	return vm.NewVirtualMachine(ctx, sess, child.Reference()), nil
}

// populate the vm attributes for the specified morefs
func populateVMAttributes(ctx context.Context, sess *session.Session, refs []types.ManagedObjectReference) ([]mo.VirtualMachine, error) {
	var vms []mo.VirtualMachine

	// current attributes we care about
	attrib := []string{"config", "runtime.powerState", "summary"}

	// populate the vm properties
	err := sess.Retrieve(ctx, refs, attrib, &vms)
	return vms, err
}

// convert the infra containers to a container object
func convertInfraContainers(vms []mo.VirtualMachine, all *bool) []Container {
	var containerVMs []Container

	for i := range vms {
		// poweredOn or all states
		if !*all && vms[i].Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff {
			// don't want it
			log.Debugf("Skipping poweredOff VM %s", vms[i].Config.Name)
			continue
		}

		container := &Container{ExecConfig: &metadata.ExecutorConfig{}}
		source := extraconfig.OptionValueSource(vms[i].Config.ExtraConfig)
		extraconfig.Decode(source, container.ExecConfig)

		// check extraConfig to see if we have a containerVM -- assumes
		// that ID will always be populated for each containerVM
		if container.ExecConfig == nil || container.ExecConfig.ID == "" {
			log.Debugf("Skipping non-container vm %s", vms[i].Config.Name)
			continue
		}

		// set state & friendly status
		if vms[i].Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
			container.State = StateRunning
			container.Status = "Running"
		} else {
			container.State = StateStopped
			container.Status = "Stopped"
		}
		container.VMUnsharedDisk = vms[i].Summary.Storage.Unshared

		containerVMs = append(containerVMs, *container)

	}

	return containerVMs
}
