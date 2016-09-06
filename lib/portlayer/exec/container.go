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
	"fmt"
	"io"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/sys"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

type State int

const (
	StateStarting = iota + 1
	StateRunning
	StateStopping
	StateStopped
	StateSuspending
	StateSuspended
	StateCreated
	StateRemoving
	StateRemoved
	StateUnknown

	propertyCollectorTimeout = 3 * time.Minute
)

type Container struct {
	sync.Mutex

	ExecConfig *executor.ExecutorConfig
	State      State

	VMUnsharedDisk int64

	vm *vm.VirtualMachine
}

func NewContainer(id uid.UID) *Handle {
	con := &Container{
		ExecConfig: &executor.ExecutorConfig{},
		State:      StateStopped,
	}
	con.ExecConfig.ID = id.String()
	return con.newHandle()
}

func GetContainer(id uid.UID) *Handle {
	// get from the cache
	container := containers.Container(id.String())
	if container != nil {
		return container.newHandle()
	}
	return nil

}

func (s State) String() string {
	switch s {
	case StateCreated:
		return "Created"
	case StateStarting:
		return "Starting"
	case StateRunning:
		return "Running"
	case StateRemoving:
		return "Removing"
	case StateRemoved:
		return "Removed"
	case StateStopping:
		return "Stopping"
	case StateStopped:
		return "Stopped"

	}
	return ""
}

func (c *Container) newHandle() *Handle {
	return newHandle(c)
}

func (c *Container) Commit(ctx context.Context, sess *session.Session, h *Handle, waitTime *int32) error {
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
			res, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
				return VCHConfig.VirtualApp.CreateChildVM_Task(ctx, *h.Spec.Spec(), nil)
			})

			c.State = StateCreated
		} else {
			// Find the Virtual Machine folder that we use
			var folders *object.DatacenterFolders
			folders, err = sess.Datacenter.Folders(ctx)
			if err != nil {
				log.Errorf("Could not get folders")
				return err
			}
			parent := folders.VmFolder

			// Create the vm
			res, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
				return parent.CreateVM(ctx, *h.Spec.Spec(), VCHConfig.ResourcePool, nil)
			})

			c.State = StateCreated
		}

		if err != nil {
			log.Errorf("Something failed. Spec was %+v", *h.Spec.Spec())
			return err
		}

		c.vm = vm.NewVirtualMachine(ctx, sess, res.Result.(types.ManagedObjectReference))

		// clear the spec as we've acted on it
		h.Spec = nil
	}

	// if we're stopping the VM, do so before the reconfigure to preserve the extraconfig
	if h.State != nil && *h.State == StateStopped {
		// stop the container
		if err := h.Container.stop(ctx, waitTime); err != nil {
			return err
		}

		c.State = *h.State
	}

	if h.Spec != nil {
		// FIXME: add check that the VM is powered off - it should be, but this will destroy the
		// extraconfig if it's not.

		s := h.Spec.Spec()
		_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
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
	// get existing state and set to starting
	// if there's a failure we'll revert to existing
	existingState := c.State
	c.State = StateStarting

	// Power on
	_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return c.vm.PowerOn(ctx)
	})
	if err != nil {
		c.State = existingState
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
		c.State = existingState
		return fmt.Errorf("unable to wait for process launch status: %s", err.Error())
	}

	if detail != "true" {
		c.State = existingState
		return errors.New(detail)
	}

	return nil
}

func (c *Container) waitForPowerState(ctx context.Context, max time.Duration, state types.VirtualMachinePowerState) (bool, error) {
	timeout, cancel := context.WithTimeout(ctx, max)
	defer cancel()

	err := c.vm.WaitForPowerState(timeout, state)
	if err != nil {
		return timeout.Err() == err, err
	}

	return false, nil
}

func (c *Container) shutdown(ctx context.Context, waitTime *int32) error {
	wait := 10 * time.Second // default
	if waitTime != nil && *waitTime > 0 {
		wait = time.Duration(*waitTime) * time.Second
	}

	stop := []string{c.ExecConfig.StopSignal, string(ssh.SIGKILL)}
	if stop[0] == "" {
		stop[0] = string(ssh.SIGTERM)
	}

	for _, sig := range stop {
		msg := fmt.Sprintf("sending kill -%s %s", sig, c.ExecConfig.ID)
		log.Info(msg)

		err := c.startGuestProgram(ctx, "kill", sig)
		if err != nil {
			return fmt.Errorf("%s: %s", msg, err)
		}

		log.Infof("waiting %s for %s to power off", wait, c.ExecConfig.ID)
		timeout, err := c.waitForPowerState(ctx, wait, types.VirtualMachinePowerStatePoweredOff)
		if err == nil {
			return nil // VM has powered off
		}

		if !timeout {
			return err // error other than timeout
		}

		log.Warnf("timeout (%s) waiting for %s to power off via SIG%s", wait, c.ExecConfig.ID, sig)
	}

	return fmt.Errorf("failed to shutdown %s via kill signals %s", c.ExecConfig.ID, stop)
}

func (c *Container) stop(ctx context.Context, waitTime *int32) error {
	defer trace.End(trace.Begin("Container.stop"))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}
	// get existing state and set to stopping
	// if there's a failure we'll revert to existing
	existingState := c.State
	c.State = StateStopping

	err := c.shutdown(ctx, waitTime)
	if err == nil {
		return nil
	}

	log.Warnf("stopping %s via hard power off due to: %s", c.ExecConfig.ID, err)

	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return c.vm.PowerOff(ctx)
	})

	if err != nil {
		// It is possible the VM has finally shutdown in between, ignore the error in that case
		if terr, ok := err.(task.Error); ok {
			if serr, ok := terr.Fault().(*types.InvalidPowerState); ok {
				if serr.ExistingState == types.VirtualMachinePowerStatePoweredOff {
					log.Warnf("power off %s task skipped (state was already %s)", c.ExecConfig.ID, serr.ExistingState)
					return nil
				}
			}
		}
		c.State = existingState
	}
	return err
}

func (c *Container) startGuestProgram(ctx context.Context, name string, args string) error {
	o := guest.NewOperationsManager(c.vm.Client.Client, c.vm.Reference())
	m, err := o.ProcessManager(ctx)
	if err != nil {
		return err
	}

	spec := types.GuestProgramSpec{
		ProgramPath: name,
		Arguments:   args,
	}

	auth := types.NamePasswordAuthentication{
		Username: c.ExecConfig.ID,
	}

	_, err = m.StartProgram(ctx, &auth, &spec)

	return err
}

func (c *Container) Signal(ctx context.Context, num int64) error {
	defer trace.End(trace.Begin("Container.Signal"))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}

	return c.startGuestProgram(ctx, "kill", fmt.Sprintf("%d", num))
}

func (c *Container) LogReader(ctx context.Context) (io.ReadCloser, error) {
	defer trace.End(trace.Begin("Container.LogReader"))

	if c.vm == nil {
		return nil, fmt.Errorf("vm not set")
	}

	p := soap.DefaultDownload

	url, err := c.vm.DSPath(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s/%s.log", url.Path, c.ExecConfig.ID)

	log.Infof("pulling %s", name)

	r, _, err := c.vm.Datastore.Download(ctx, name, &p)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// NotFoundError is returned when a types.ManagedObjectNotFound is returned from a vmomi call
type NotFoundError struct {
	err error
}

func (r NotFoundError) Error() string {
	return "VM has either been deleted or has not been fully created"
}

// RemovePowerError is returned when attempting to remove a containerVM that is powered on
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
		return NotFoundError{}
	}

	// check state first
	if c.State == StateRunning {
		return RemovePowerError{fmt.Errorf("Container is powered on")}
	}

	// get existing state and set to removing
	// if there's a failure we'll revert to existing
	existingState := c.State
	c.State = StateRemoving

	// get the folder the VM is in
	url, err := c.vm.DSPath(ctx)
	if err != nil {

		// handle the out-of-band removal case
		if soap.IsSoapFault(err) {
			fault := soap.ToSoapFault(err).VimFault()
			if _, ok := fault.(types.ManagedObjectNotFound); ok {
				containers.Remove(c.ExecConfig.ID)
				return NotFoundError{}
			}
		}

		log.Errorf("Failed to get datastore path for %s: %s", c.ExecConfig.ID, err)
		c.State = existingState
		return err
	}
	// FIXME: was expecting to find a utility function to convert to/from datastore/url given
	// how widely it's used but couldn't - will ask around.
	dsPath := fmt.Sprintf("[%s] %s", url.Host, url.Path)

	//removes the vm from vsphere, but detaches the disks first
	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return c.vm.DeleteExceptDisks(ctx)
	})
	if err != nil {
		c.State = existingState
		return err
	}

	// remove from datastore
	fm := object.NewFileManager(c.vm.Client.Client)

	if _, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return fm.DeleteDatastoreFile(ctx, dsPath, sess.Datacenter)
	}); err != nil {
		c.State = existingState
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

// return specific container to caller
func ContainerInfo(containerID string) *Container {
	return containers.Container(containerID)
}

func Containers(all bool) []*Container {
	return containers.Containers(all)
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

func instanceUUID(id string) (string, error) {
	// generate VM instance uuid, which will be used to query back VM
	u, err := sys.UUID()
	if err != nil {
		return "", err
	}
	namespace, err := uuid.Parse(u)
	if err != nil {
		return "", errors.Errorf("unable to parse VCH uuid: %s", err)
	}
	return uuid.NewSHA1(namespace, []byte(id)).String(), nil
}

// find the childVM for this resource pool by name
func childVM(ctx context.Context, sess *session.Session, name string) (*vm.VirtualMachine, error) {
	// Search container back through instance UUID
	uuid, err := instanceUUID(name)
	if err != nil {
		detail := fmt.Sprintf("unable to get instance UUID: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	searchIndex := object.NewSearchIndex(sess.Client.Client)
	child, err := searchIndex.FindByUuid(ctx, sess.Datacenter, uuid, true, nil)

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
func convertInfraContainers(vms []mo.VirtualMachine, all bool) []*Container {
	var containerVMs []*Container

	for i := range vms {
		// poweredOn or all states
		if !all && vms[i].Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff {
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

		// set state
		if vms[i].Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
			container.State = StateRunning
		} else {
			// look in the container cache and check state
			// if it's created we'll take that as it's been created, but
			// not started
			cached := containers.Container(container.ExecConfig.ID)
			if cached != nil && cached.State == StateCreated {
				container.State = StateCreated
			} else {
				container.State = StateStopped
			}
		}
		if vms[i].Summary.Storage != nil {
			container.VMUnsharedDisk = vms[i].Summary.Storage.Unshared
		}

		containerVMs = append(containerVMs, container)

	}

	return containerVMs
}
