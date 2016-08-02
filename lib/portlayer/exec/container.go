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
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
	"golang.org/x/net/context"
)

func init() {
	NewContainerCache()
}

type State int

const (
	StateRunning = iota + 1
	StateStopped
	StateCreated

	propertyCollectorTimeout = 3 * time.Minute
)

type Container struct {
	sync.Mutex

	ExecConfig *executor.ExecutorConfig
	State      State

	// friendly description of state
	Status string

	VMUnsharedDisk int64

	vm *vm.VirtualMachine
}

func NewContainer(id ID) *Handle {
	con := &Container{
		ExecConfig: &executor.ExecutorConfig{},
		State:      StateStopped,
	}
	con.ExecConfig.ID = id.String()
	return con.newHandle()
}

func GetContainer(id ID) *Handle {
	// get from the cache
	container := containers.Container(id.String())
	if container != nil {
		return container.newHandle()
	}
	return nil

}

func (c *Container) newHandle() *Handle {
	return newHandle(c)
}

// TODO: Is this used anywhere?
func (c *Container) cacheExecConfig(ec *executor.ExecutorConfig) {
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

		var res *types.TaskInfo
		var err error

		if sess.IsVC() && VCHConfig.VirtualApp != nil {
			// Create the vm
			res, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
				return VCHConfig.VirtualApp.CreateChildVM_Task(ctx, *h.Spec.Spec(), nil)
			})
			// set the status to created
			c.State = StateCreated
		} else {
			// Find the Virtual Machine folder that we use
			var folders *object.DatacenterFolders
			folders, err = sess.Datacenter.Folders(ctx)
			if err != nil {
				return err
			}
			parent := folders.VmFolder

			// Create the vm
			res, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
				return parent.CreateVM(ctx, *h.Spec.Spec(), VCHConfig.ResourcePool, nil)
			})

			// set the status to created
			c.State = StateCreated
		}

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

	// add or overwrite the container in the cache
	containers.Put(c)
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

type RemovePowerError struct {
	err error
}

func (r RemovePowerError) Error() string {
	return r.err.Error()
}

func (c *Container) Remove(ctx context.Context, sess *session.Session) error {
	defer trace.End(trace.Begin("Container.Remove"))
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

	state, err := c.vm.PowerState(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get vm power status %q: %s", c.vm.Reference(), err)
	}

	// don't remove the containerVM if it is powered on
	if state == types.VirtualMachinePowerStatePoweredOn {
		return RemovePowerError{fmt.Errorf("Container is powered on")}
	}

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

	//remove container from cache
	containers.Remove(c.ExecConfig.ID)
	return nil
}

func (c *Container) Update(ctx context.Context, sess *session.Session) (*executor.ExecutorConfig, error) {
	defer trace.End(trace.Begin("Container.Update"))
	c.Lock()
	defer c.Unlock()

	if c.vm == nil {
		return nil, fmt.Errorf("container does not have a vm")
	}

	var vm []mo.VirtualMachine

	if err := sess.Retrieve(ctx, []types.ManagedObjectReference{c.vm.Reference()}, []string{"config"}, &vm); err != nil {
		return nil, err
	}

	extraconfig.Decode(vmomi.OptionValueSource(vm[0].Config.ExtraConfig), c.ExecConfig)
	return c.ExecConfig, nil
}

// Grab the info for the requested container
// TODO:  Possibly change so that handler requests a handle to the
// container and if it's not present then search and return a handle
func ContainerInfo(ctx context.Context, sess *session.Session, containerID ID) (*Container, error) {
	// first  lets see if we have it in the cache
	container := containers.Container(containerID.String())
	// if we missed search for it...
	if container == nil {
		// search
		vm, err := childVM(ctx, sess, containerID.String())
		if err != nil || vm == nil {
			log.Debugf("ContainerInfo failed to find childVM: %s", err.Error())
			return nil, fmt.Errorf("Container Not Found: %s", containerID)
		}
		container = &Container{vm: vm}
	}

	// get properties for specific containerVMs
	// we must do this since we don't know that we have a valid state
	// This will be refactored when event streaming hits
	vms, err := populateVMAttributes(ctx, sess, []types.ManagedObjectReference{container.vm.Reference()})
	if err != nil {
		return nil, err
	}
	// convert the VMs to container objects -- include
	// powered off vms
	all := true
	cc := convertInfraContainers(vms, &all)

	switch len(cc) {
	case 0:
		// we found a vm, but it doesn't appear to be a container VM
		return nil, fmt.Errorf("%s does not appear to be a container", containerID)
	case 1:
		// we have a winner
		return cc[0], nil
	default:
		// we manged to find multiple vms
		return nil, fmt.Errorf("multiple containers named %s found", containerID)
	}
}

// return a list of container attributes
func List(ctx context.Context, sess *session.Session, all *bool) ([]*Container, error) {

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
	if err := VCHConfig.ResourcePool.Properties(ctx, VCHConfig.ResourcePool.Reference(), []string{"vm"}, &rp); err != nil {
		name := VCHConfig.ResourcePool.Name()
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
	child, err := searchIndex.FindChild(ctx, VCHConfig.ResourcePool.Reference(), name)
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
func convertInfraContainers(vms []mo.VirtualMachine, all *bool) []*Container {
	var containerVMs []*Container

	for i := range vms {
		// poweredOn or all states
		if !*all && vms[i].Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff {
			// don't want it
			log.Debugf("Skipping poweredOff VM %s", vms[i].Config.Name)
			continue
		}

		container := &Container{ExecConfig: &executor.ExecutorConfig{}}
		source := vmomi.OptionValueSource(vms[i].Config.ExtraConfig)
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
			// look in the container cache and check status
			// if it's created we'll take that as it's been created, but
			// not started
			cached := containers.Container(container.ExecConfig.ID)
			if cached != nil && cached.State == StateCreated {
				container.State = StateCreated
				container.Status = "Created"
			} else {
				container.State = StateStopped
				container.Status = "Stopped"
			}
		}
		container.VMUnsharedDisk = vms[i].Summary.Storage.Unshared

		containerVMs = append(containerVMs, container)

	}

	return containerVMs
}
