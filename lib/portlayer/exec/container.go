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
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/sys"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

type State int

const (
	StateUnknown State = iota
	StateStarting
	StateRunning
	StateStopping
	StateStopped
	StateSuspending
	StateSuspended
	StateCreated
	StateCreating
	StateRemoving
	StateRemoved

	propertyCollectorTimeout = 3 * time.Minute
	containerLogName         = "output.log"

	vmNotSuspendedKey = "msg.suspend.powerOff.notsuspended"
)

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

// ConcurrentAccessError is returned when concurrent calls tries to modify same object
type ConcurrentAccessError struct {
	err error
}

func (r ConcurrentAccessError) Error() string {
	return r.err.Error()
}

type Container struct {
	m sync.Mutex

	// Current values
	ExecConfig *executor.ExecutorConfig
	state      State

	// Size of the leaf (unused)
	VMUnsharedDisk int64

	vm *vm.VirtualMachine

	logFollowers []io.Closer

	// Current state
	Config  *types.VirtualMachineConfigInfo
	Runtime *types.VirtualMachineRuntimeInfo

	newStateEvents map[State]chan struct{}
}

func NewContainer(id uid.UID) *Handle {
	con := &Container{
		ExecConfig:     &executor.ExecutorConfig{},
		state:          StateCreating,
		newStateEvents: make(map[State]chan struct{}),
	}
	con.ExecConfig.ID = id.String()
	return newHandle(con, con.state)
}

func GetContainer(ctx context.Context, id uid.UID) *Handle {
	// get from the cache
	container := Containers.Container(id.String())
	if container != nil {
		return container.NewHandle(ctx)
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
	case StateUnknown:
		return "Unknown"
	}
	return ""
}

// CurrentState returns current state.
func (c *Container) CurrentState() State {
	c.m.Lock()
	defer c.m.Unlock()
	return c.state
}

// SetState changes container state.
func (c *Container) SetState(s State) State {
	c.m.Lock()
	defer c.m.Unlock()
	return c.updateState(s)
}

func (c *Container) updateState(s State) State {
	log.Debugf("Setting container %s state: %s", c.ExecConfig.ID, s)
	prevState := c.state
	if s != c.state {
		c.state = s
		if ch, ok := c.newStateEvents[s]; ok {
			delete(c.newStateEvents, s)
			close(ch)
		}
	}
	return prevState
}

var closedEventChannel = func() <-chan struct{} {
	a := make(chan struct{})
	close(a)
	return a
}()

// WaitForState subscribes a caller to an event returning
// a channel that will be closed when an expected state is set.
// If expected state is already set the caller will receive a closed channel immediately.
func (c *Container) WaitForState(s State) <-chan struct{} {
	c.m.Lock()
	defer c.m.Unlock()

	if s == c.state {
		return closedEventChannel
	}

	if ch, ok := c.newStateEvents[s]; ok {
		return ch
	}

	eventChan := make(chan struct{})
	c.newStateEvents[s] = eventChan
	return eventChan
}

func (c *Container) NewHandle(ctx context.Context) *Handle {
	c.m.Lock()
	defer c.m.Unlock()

	// Call property collector to fill the data
	if c.vm != nil {
		if err := c.refresh(ctx); err != nil {
			log.Errorf("refreshing container %s failed: %s", c.ExecConfig.ID, err)
			return nil // nil indicates error
		}
	}

	return newHandle(c, c.state)
}

// Refresh calls the property collector to get config and runtime info and Guest RPC for ExtraConfig
func (c *Container) Refresh(ctx context.Context) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	c.m.Lock()
	defer c.m.Unlock()

	return c.refresh(ctx)
}

func (c *Container) refresh(ctx context.Context) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	var o mo.VirtualMachine

	// make sure we have vm
	if c.vm == nil {
		return fmt.Errorf("There is no backing VirtualMachine %#v", c)
	}
	if err := c.vm.Properties(ctx, c.vm.Reference(), []string{"config", "runtime"}, &o); err != nil {
		return err
	}

	c.Config = o.Config
	c.Runtime = &o.Runtime

	// Get the ExtraConfig
	extraconfig.Decode(vmomi.OptionValueSource(o.Config.ExtraConfig), c.ExecConfig)

	return nil
}

// Commit executes the requires steps on the handle
func (c *Container) Commit(ctx context.Context, sess *session.Session, h *Handle, waitTime *int32) error {
	defer trace.End(trace.Begin(h.ExecConfig.ID))

	// hold the event that has occurred
	var commitEvent string

	c.m.Lock()
	defer c.m.Unlock()

	// If an event has occurred then put the container in the cache
	// and publish the container event
	defer func() {
		log.Debugf("Commiting container %s status as: %s", c.ExecConfig.ID, commitEvent)
		if commitEvent != "" {
			Containers.Put(c)
			publishContainerEvent(c.ExecConfig.ID, time.Now().UTC(), commitEvent)
		}
	}()

	if c.vm == nil {
		if sess == nil {
			// session must not be nil
			return fmt.Errorf("no session provided for commit operation")
		}

		// the only permissible operation is to create a VM
		if h.Spec == nil {
			return fmt.Errorf("only create operations can be committed without an existing VM")
		}

		var res *types.TaskInfo
		var err error
		if sess.IsVC() && Config.VirtualApp.ResourcePool != nil {
			// Create the vm
			res, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
				return Config.VirtualApp.CreateChildVM_Task(ctx, *h.Spec.Spec(), nil)
			})
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
				return parent.CreateVM(ctx, *h.Spec.Spec(), Config.ResourcePool, nil)
			})
		}

		if err != nil {
			log.Errorf("Something failed. Spec was %+v", *h.Spec.Spec())
			return err
		}

		c.vm = vm.NewVirtualMachine(ctx, sess, res.Result.(types.ManagedObjectReference))
		c.updateState(StateCreated)

		commitEvent = events.ContainerCreated

		// clear the spec as we've acted on it
		h.Spec = nil

		c.ExecConfig = &h.ExecConfig
		// refresh the struct with what propery collector provides
		if err = c.refresh(ctx); err != nil {
			return err
		}
	}

	// if we're stopping the VM, do so before the reconfigure to preserve the extraconfig
	if h.CurrentState() == StateStopped &&
		c.Runtime != nil && c.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
		// stop the container
		if err := h.Container.stop(ctx, waitTime); err != nil {
			return err
		}

		commitEvent = events.ContainerStopped

		// refresh the struct with what propery collector provides
		if err := c.refresh(ctx); err != nil {
			return err
		}
	}

	if h.Spec != nil {
		s := h.Spec.Spec()

		// nilify ExtraConfig if vm is running
		if c.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
			log.Errorf("Nilifying ExtraConfig as we are running")
			s.ExtraConfig = nil
		}

		// set ChangeVersion. This property is useful because it guards against updates that have happened between when the VM’s config is read and when it is applied.
		// Will return "Cannot complete operation due to concurrent modification by another operation.." on failure
		s.ChangeVersion = c.Config.ChangeVersion

		_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
			return c.vm.Reconfigure(ctx, *s)
		})
		if err != nil {
			log.Errorf("Reconfigure failed with %#+v", err)

			// Check whether we get ConcurrentAccess and wrap it if needed
			if f, ok := err.(types.HasFault); ok {
				switch f.Fault().(type) {
				case *types.ConcurrentAccess:
					log.Errorf("We have ConcurrentAccess for version %s", s.ChangeVersion)

					return ConcurrentAccessError{err}
				}
			}
			return err
		}

		c.ExecConfig = &h.ExecConfig

		// refresh the struct with what propery collector provides
		if err = c.refresh(ctx); err != nil {
			return err
		}
	}

	if h.CurrentState() == StateRunning &&
		c.Runtime != nil && c.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff {
		// start the container
		if err := h.Container.start(ctx); err != nil {
			return err
		}

		commitEvent = events.ContainerStarted

		// refresh the struct with what property collector provides
		if err := c.refresh(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Start starts a container vm with the given params
func (c *Container) start(ctx context.Context) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}
	// get existing state and set to starting
	// if there's a failure we'll revert to existing
	finalState := c.updateState(StateStarting)

	defer func() { c.updateState(finalState) }()

	// Power on
	_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return c.vm.PowerOn(ctx)
	})
	if err != nil {
		return err
	}

	// guestinfo key that we want to wait for
	key := fmt.Sprintf("guestinfo.vice..sessions|%s.started", c.ExecConfig.ID)
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

	// this state will be set by defer function.
	finalState = StateRunning
	return nil
}

func (c *Container) waitForPowerState(ctx context.Context, max time.Duration, state types.VirtualMachinePowerState) (bool, error) {
	defer trace.End(trace.Begin(c.ExecConfig.ID))
	timeout, cancel := context.WithTimeout(ctx, max)
	defer cancel()

	err := c.vm.WaitForPowerState(timeout, state)
	if err != nil {
		return timeout.Err() == err, err
	}

	return false, nil
}

func (c *Container) shutdown(ctx context.Context, waitTime *int32) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))
	wait := 10 * time.Second // default
	if waitTime != nil && *waitTime > 0 {
		wait = time.Duration(*waitTime) * time.Second
	}
	cs := c.ExecConfig.Sessions[c.ExecConfig.ID]
	stop := []string{cs.StopSignal, string(ssh.SIGKILL)}
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
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}

	defer c.onStop()

	// get existing state and set to stopping
	// if there's a failure we'll revert to existing

	existingState := c.updateState(StateStopping)

	err := c.shutdown(ctx, waitTime)
	if err == nil {
		c.updateState(StateStopped)
		return nil
	}

	log.Warnf("stopping %s via hard power off due to: %s", c.ExecConfig.ID, err)

	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return c.vm.PowerOff(ctx)
	})

	if err != nil {

		// It is possible the VM has finally shutdown in between, ignore the error in that case
		if terr, ok := err.(task.Error); ok {
			switch terr := terr.Fault().(type) {
			case *types.InvalidPowerState:
				if terr.ExistingState == types.VirtualMachinePowerStatePoweredOff {
					log.Warnf("power off %s task skipped (state was already %s)", c.ExecConfig.ID, terr.ExistingState)
					return nil
				}
				log.Warnf("invalid power state during power off: %s", terr.ExistingState)

			case *types.GenericVmConfigFault:

				// Check if the poweroff task was canceled due to a concurrent guest shutdown
				if len(terr.FaultMessage) > 0 && terr.FaultMessage[0].Key == vmNotSuspendedKey {
					log.Infof("power off %s task skipped due to guest shutdown", c.ExecConfig.ID)
					return nil
				}
				log.Warnf("generic vm config fault during power off: %#v", terr)

			default:
				log.Warnf("hard power off failed due to: %#v", terr)
			}
		}
		c.updateState(existingState)
		return err
	}
	c.updateState(StateStopped)
	return nil
}

func (c *Container) startGuestProgram(ctx context.Context, name string, args string) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))
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
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}

	return c.startGuestProgram(ctx, "kill", fmt.Sprintf("%d", num))
}

func (c *Container) onStop() {
	lf := c.logFollowers
	c.logFollowers = nil

	log.Debugf("Container(%s) closing %d log followers", c.ExecConfig.ID, len(lf))
	for _, l := range lf {
		_ = l.Close()
	}
}

func (c *Container) LogReader(ctx context.Context, tail int, follow bool) (io.ReadCloser, error) {
	defer trace.End(trace.Begin(c.ExecConfig.ID))
	c.m.Lock()
	defer c.m.Unlock()

	if c.vm == nil {
		return nil, fmt.Errorf("vm not set")
	}

	url, err := c.vm.DSPath(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s/%s", url.Path, containerLogName)

	log.Infof("pulling %s", name)

	file, err := c.vm.Datastore.Open(ctx, name)
	if err != nil {
		return nil, err
	}

	if tail >= 0 {
		err = file.Tail(tail)
		if err != nil {
			return nil, err
		}
	}

	if follow && c.state == StateRunning {
		follower := file.Follow(time.Second)

		c.logFollowers = append(c.logFollowers, follower)

		return follower, nil
	}

	return file, nil
}

// Remove removes a containerVM after detaching the disks
func (c *Container) Remove(ctx context.Context, sess *session.Session) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))
	c.m.Lock()
	defer c.m.Unlock()

	if c.vm == nil {
		return NotFoundError{}
	}

	// check state first
	if c.state == StateRunning {
		return RemovePowerError{fmt.Errorf("Container is powered on")}
	}

	// get existing state and set to removing
	// if there's a failure we'll revert to existing
	existingState := c.updateState(StateRemoving)

	// get the folder the VM is in
	url, err := c.vm.DSPath(ctx)
	if err != nil {

		// handle the out-of-band removal case
		if soap.IsSoapFault(err) {
			fault := soap.ToSoapFault(err).VimFault()
			if _, ok := fault.(types.ManagedObjectNotFound); ok {
				Containers.Remove(c.ExecConfig.ID)
				return NotFoundError{}
			}
		}

		log.Errorf("Failed to get datastore path for %s: %s", c.ExecConfig.ID, err)
		c.updateState(existingState)
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
		f, ok := err.(types.HasFault)
		if !ok {
			c.updateState(existingState)
			return err
		}
		switch f.Fault().(type) {
		case *types.InvalidState:
			log.Warnf("container VM is in invalid state, unregistering")
			if err := c.vm.Unregister(ctx); err != nil {
				log.Errorf("Error while attempting to unregister container VM: %s", err)
				return err
			}
		default:
			log.Debugf("Fault while attempting to destroy vm: %#v", f.Fault())
			c.updateState(existingState)
			return err
		}
	}

	// remove from datastore
	fm := object.NewFileManager(c.vm.Client.Client)

	if _, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return fm.DeleteDatastoreFile(ctx, dsPath, sess.Datacenter)
	}); err != nil {
		// at this phase error doesn't matter. Just log it.
		log.Debugf("Failed to delete %s, %s", dsPath, err)
	}

	//remove container from cache
	Containers.Remove(c.ExecConfig.ID)
	return nil
}

// get the containerVMs from infrastructure for this resource pool
func infraContainers(ctx context.Context, sess *session.Session) ([]*Container, error) {
	defer trace.End(trace.Begin(""))
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

	return convertInfraContainers(ctx, sess, vms), nil
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
	defer trace.End(trace.Begin(""))
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
	defer trace.End(trace.Begin(fmt.Sprintf("populating %d refs", len(refs))))
	var vms []mo.VirtualMachine

	// current attributes we care about
	attrib := []string{"config", "runtime.powerState", "summary"}

	// populate the vm properties
	err := sess.Retrieve(ctx, refs, attrib, &vms)
	return vms, err
}

// convert the infra containers to a container object
func convertInfraContainers(ctx context.Context, sess *session.Session, vms []mo.VirtualMachine) []*Container {
	defer trace.End(trace.Begin(fmt.Sprintf("converting %d containers", len(vms))))
	var cons []*Container

	for _, v := range vms {
		source := vmomi.OptionValueSource(v.Config.ExtraConfig)
		c := &Container{
			state:          StateCreated,
			newStateEvents: make(map[State]chan struct{}),
		}
		extraconfig.Decode(source, &c.ExecConfig)
		id := uid.Parse(c.ExecConfig.ID)
		if id == uid.NilUID {
			log.Warnf("skipping converting container VM %s: could not parse id", v.Reference())
			continue
		}

		// set state
		switch v.Runtime.PowerState {
		case types.VirtualMachinePowerStatePoweredOn:
			c.state = StateRunning
		case types.VirtualMachinePowerStatePoweredOff:
			// check if any of the sessions was started
			for _, s := range c.ExecConfig.Sessions {
				if s.Started != "" {
					c.state = StateStopped
					break
				}
			}
		case types.VirtualMachinePowerStateSuspended:
			log.Warnf("skipping converting container VM %s: invalid power state %s", v.Reference(), v.Runtime.PowerState)
			continue
		}

		if v.Summary.Storage != nil {
			c.VMUnsharedDisk = v.Summary.Storage.Unshared
		}

		c.vm = vm.NewVirtualMachine(ctx, sess, v.Reference())
		cons = append(cons, c)
	}

	return cons
}
