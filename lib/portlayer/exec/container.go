// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/iolog"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/sys"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"
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

	containerLogName = "output.log"

	vmNotSuspendedKey = "msg.suspend.powerOff.notsuspended"
	vmPoweringOffKey  = "msg.rpc.error.poweringoff"
)

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

func IsConcurrentAccessError(err error) bool {
	_, ok := err.(ConcurrentAccessError)
	return ok
}

type DevicesInUseError struct {
	Devices []string
}

func (e DevicesInUseError) Error() string {
	return fmt.Sprintf("device %s in use", strings.Join(e.Devices, ","))
}

// Container is used to return data about a container during inspection calls
// It is a copy rather than a live reflection and does not require locking
type ContainerInfo struct {
	containerBase

	state State

	// Size of the leaf (unused)
	VMUnsharedDisk int64
}

// Container is used for an entry in the container cache - this is a "live" representation
// of containers in the infrastructure.
// DANGEROUS USAGE CONSTRAINTS:
//   None of the containerBase fields should be partially updated - consider them immutable once they're
//   part of a cache entry
//   i.e. Do not make changes in containerBase.ExecConfig - only swap, under lock, the pointer for a
//   completely new ExecConfig.
//   This constraint allows us to avoid deep copying those structs every time a container is inspected
type Container struct {
	m sync.Mutex

	ContainerInfo

	logFollowers []io.Closer

	newStateEvents map[State]chan struct{}
}

// newContainer constructs a Container suitable for adding to the cache
// it's state is set from the Runtime.PowerState field, or StateCreated if that is not
// viable
// This copies (shallow) the containerBase that's provided
func newContainer(base *containerBase) *Container {
	c := &Container{
		ContainerInfo: ContainerInfo{
			containerBase: *base,
			state:         StateCreated,
		},
		newStateEvents: make(map[State]chan struct{}),
	}

	// if this is a creation path, then Runtime will be nil
	if base.Runtime != nil {
		// set state
		switch base.Runtime.PowerState {
		case types.VirtualMachinePowerStatePoweredOn:
			// the containerVM is poweredOn, so set state to starting
			// then check to see if a start was successful
			c.state = StateStarting
			// If any sessions successfully started then set to running
			for _, s := range base.ExecConfig.Sessions {
				if s.Started != "" {
					c.state = StateRunning
					break
				}
			}
		case types.VirtualMachinePowerStatePoweredOff:
			// check if any of the sessions was started
			for _, s := range base.ExecConfig.Sessions {
				if s.Started != "" {
					c.state = StateStopped
					break
				}
			}
		case types.VirtualMachinePowerStateSuspended:
			c.state = StateSuspended
			log.Warnf("container VM %s: invalid power state %s", base.vm.Reference(), base.Runtime.PowerState)
		}
	}

	return c
}

func GetContainer(ctx context.Context, id uid.UID) *Handle {
	// get from the cache
	container := Containers.Container(id.String())
	if container != nil {
		return container.NewHandle(ctx)
	}

	return nil
}

func (c *ContainerInfo) String() string {
	return c.ExecConfig.ID
}

// State returns the state at the time the ContainerInfo object was created
func (c *ContainerInfo) State() State {
	return c.state
}

// Info returns a copy of the public container configuration that
// is consistent and copied under lock
func (c *Container) Info() *ContainerInfo {
	c.m.Lock()
	defer c.m.Unlock()

	info := c.ContainerInfo
	return &info
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

// transitionState changes the container state to finalState if the current state is initialState
// and returns an error otherwise.
func (c *Container) transitionState(initialState, finalState State) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.state == initialState {
		c.state = finalState
		log.Debugf("Set container %s state: %s", c.ExecConfig.ID, finalState)
		return nil
	}

	return fmt.Errorf("container state is %s and was not changed to %s", c.state, finalState)
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
	// Call property collector to fill the data
	if c.vm != nil {
		// FIXME: this should be calling the cache to decide if a refresh is needed
		if err := c.Refresh(ctx); err != nil {
			log.Errorf("refreshing container %s failed: %s", c.ExecConfig.ID, err)
			return nil // nil indicates error
		}
	}

	// return a handle that represents zero changes over the current configuration
	// for this container
	return newHandle(c)
}

// Refresh updates config and runtime info, holding a lock only while swapping
// the new data for the old
func (c *Container) Refresh(ctx context.Context) error {
	c.m.Lock()
	defer c.m.Unlock()

	if err := c.refresh(ctx); err != nil {
		return err
	}

	// sync power state (see issue 4872)
	switch c.containerBase.Runtime.PowerState {
	case types.VirtualMachinePowerStatePoweredOn:
		c.state = StateRunning
	case types.VirtualMachinePowerStatePoweredOff:
		c.state = StateStopped
	}

	return nil
}

func (c *Container) refresh(ctx context.Context) error {
	return c.containerBase.refresh(ctx)
}

// RefreshFromHandle updates config and runtime info, holding a lock only while swapping
// the new data for the old
func (c *Container) RefreshFromHandle(ctx context.Context, h *Handle) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.Config != nil && (h.Config == nil || h.Config.ChangeVersion != c.Config.ChangeVersion) {
		log.Warnf("container and handle ChangeVersions do not match: %s != %s", c.Config.ChangeVersion, h.Config.ChangeVersion)
		return
	}

	// power off doesn't necessarily cause a change version increment and bug1898149 occasionally impacts power on
	if c.Runtime != nil && (h.Runtime == nil || h.Runtime.PowerState != c.Runtime.PowerState) {
		log.Warnf("container and handle PowerStates do not match: %s != %s", c.Runtime.PowerState, h.Runtime.PowerState)
		return
	}

	// copy over the new state
	c.containerBase = h.containerBase
	if c.Config != nil {
		log.Debugf("Update: updated change version from handle: %s", c.Config.ChangeVersion)
	}
}

// Start starts a container vm with the given params
func (c *Container) start(ctx context.Context) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}
	// Set state to Starting
	c.SetState(StateStarting)

	err := c.containerBase.start(ctx)
	if err != nil {
		// change state to stopped because start task failed
		c.SetState(StateStopped)

		// check if locked disk error
		devices := disk.LockedDisks(err)
		if len(devices) > 0 {
			for i := range devices {
				// get device id from datastore file path
				// FIXME: find a reasonable way to get device ID from datastore path in exec
				devices[i] = strings.TrimSuffix(path.Base(devices[i]), ".vmdk")
			}
			return DevicesInUseError{devices}
		}
		return err
	}

	// wait task to set started field to something
	ctx, cancel := context.WithTimeout(ctx, constants.PropertyCollectorTimeout)
	defer cancel()

	err = c.waitForSession(ctx, c.ExecConfig.ID)
	if err != nil {
		// leave this in state starting - if it powers off then the event
		// will cause transition to StateStopped which is likely our original state
		// if the container was just taking a very long time it'll eventually
		// become responsive.

		// TODO: mechanism to trigger reinspection of long term transitional states
		return err
	}

	// Transition the state to Running only if it's Starting.
	// The current state is already Stopped if the container's process has exited or
	// a poweredoff event has been processed.
	if err = c.transitionState(StateStarting, StateRunning); err != nil {
		log.Debug(err)
	}

	return nil
}

func (c *Container) stop(ctx context.Context, waitTime *int32) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	defer c.onStop()

	// get existing state and set to stopping
	// if there's a failure we'll revert to existing
	finalState := c.SetState(StateStopping)

	err := c.containerBase.stop(ctx, waitTime)
	if err != nil {
		// we've got no idea what state the container is in at this point
		// running is an _optimistic_ statement
		// If the current state is Stopping, revert it to the old state.
		if stateErr := c.transitionState(StateStopping, finalState); stateErr != nil {
			log.Debug(stateErr)
		}

		return err
	}

	// Transition the state to Stopped only if it's Stopping.
	if err = c.transitionState(StateStopping, StateStopped); err != nil {
		log.Debug(err)
	}

	return nil
}

func (c *Container) Signal(ctx context.Context, num int64) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}

	if num == int64(syscall.SIGKILL) {
		return c.containerBase.kill(ctx)
	}

	return c.startGuestProgram(ctx, "kill", fmt.Sprintf("%d", num))
}

func (c *Container) onStop() {
	lf := c.logFollowers
	c.logFollowers = nil

	log.Debugf("Container(%s) closing %d log followers", c.ExecConfig.ID, len(lf))
	for _, l := range lf {
		// #nosec: Errors unhandled.
		_ = l.Close()
	}
}

func (c *Container) LogReader(ctx context.Context, tail int, follow bool, since int64) (io.ReadCloser, error) {
	defer trace.End(trace.Begin(c.ExecConfig.ID))
	c.m.Lock()
	defer c.m.Unlock()

	if c.vm == nil {
		return nil, fmt.Errorf("vm not set")
	}

	url, err := c.vm.VMPathNameAsURL(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s/%s", url.Path, containerLogName)

	var via string

	if c.state == StateRunning && c.vm.IsVC() {
		// #nosec: Errors unhandled.
		hosts, _ := c.vm.Datastore.AttachedHosts(ctx)
		if len(hosts) > 1 {
			// In this case, we need download from the VM host as it owns the file lock
			// #nosec: Errors unhandled.
			h, _ := c.vm.HostSystem(ctx)
			if h != nil {
				ctx = c.vm.Datastore.HostContext(ctx, h)
				via = fmt.Sprintf(" via %s", h.Reference())
			}
		}
	}

	log.Infof("pulling %s%s", name, via)

	file, err := c.vm.Datastore.Open(ctx, name)
	if err != nil {
		return nil, err
	}

	if since > 0 {
		err = file.TailFunc(tail, func(line int, message string) bool {
			if tail <= line && tail != -1 {
				return false
			}

			buf := bytes.NewBufferString(message)

			entry, err := iolog.ParseLogEntry(buf)
			if err != nil {
				log.Errorf("Error parsing log entry: %s", err.Error())
				return false
			}

			if entry.Timestamp.Unix() <= since {
				return false
			}

			return true
		})
	} else if tail >= 0 {
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
	url, err := c.vm.VMPathNameAsURL(ctx)
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

	ds, err := sess.Finder.Datastore(ctx, url.Host)
	if err != nil {
		return err
	}

	// enable Destroy
	c.vm.EnableDestroy(ctx)

	//removes the vm from vsphere, but detaches the disks first
	_, err = c.vm.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
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
		case *types.ConcurrentAccess:
			// We are getting ConcurrentAccess errors from DeleteExceptDisks - even though we don't set ChangeVersion in that path
			// We are ignoring the error because in reality the operation finishes successfully.
			log.Warnf("DeleteExceptDisks failed with ConcurrentAccess error. Ignoring it.")
		default:
			log.Debugf("Fault while attempting to destroy vm: %#v", f.Fault())
			c.updateState(existingState)
			return err
		}
	}

	// remove from datastore
	fm := ds.NewFileManager(sess.Datacenter, true)

	if err = fm.Delete(ctx, url.Path); err != nil {
		// at this phase error doesn't matter. Just log it.
		log.Debugf("Failed to delete %s, %s", url, err)
	}

	//remove container from cache
	Containers.Remove(c.ExecConfig.ID)
	publishContainerEvent(c.ExecConfig.ID, time.Now(), events.ContainerRemoved)

	return nil
}

// eventedState will determine the target container
// state based on the current container state and the vsphere event
func eventedState(e events.Event, current State) State {
	defer trace.End(trace.Begin(fmt.Sprintf("event %s received for id: %d", e.String(), e.EventID())))
	switch e.String() {
	case events.ContainerPoweredOn:
		// are we in the process of starting
		if current != StateStarting {
			return StateRunning
		}
	case events.ContainerPoweredOff:
		// are we in the process of stopping
		if current != StateStopping {
			return StateStopped
		}
	case events.ContainerSuspended:
		// are we in the process of suspending
		if current != StateSuspending {
			return StateSuspended
		}
	case events.ContainerRemoved:
		if current != StateRemoving {
			return StateRemoved
		}
	}
	return current
}

func (c *Container) OnEvent(e events.Event) {
	defer trace.End(trace.Begin(fmt.Sprintf("event %s received for id: %d", e.String(), e.EventID())))
	c.m.Lock()
	defer c.m.Unlock()

	if c.vm == nil {
		return
	}
	newState := eventedState(e, c.state)
	// do we have a state change
	if newState != c.state {
		switch newState {
		case StateStopping,
			StateRunning,
			StateStopped,
			StateSuspended:

			// container state has changed so we need to update the container attributes
			ctx, cancel := context.WithTimeout(context.Background(), constants.PropertyCollectorTimeout)
			defer cancel()

			if err := c.refresh(ctx); err != nil {
				log.Errorf("Event driven container update failed: %s", err)
			}

			c.updateState(newState)
			if newState == StateStopped {
				c.onStop()
			}
			log.Debugf("Container(%s) state set to %s via event activity", c, newState)
		case StateRemoved:
			if c.vm != nil && c.vm.IsFixing() {
				// is fixing vm, which will be registered back soon, so do not remove from containers cache
				log.Debugf("Container(%s) %s is being fixed - %s event ignored", c.ExecConfig.ID, newState)

				// Received remove event triggered by unregister VM operation - leave
				// fixing state now. In a loaded environment, the remove event may be
				// received after vm.fixVM() has returned, at which point the container
				// should still be in fixing state to avoid removing it from the cache below.
				c.vm.LeaveFixingState()
				// since we're leaving the container in cache, just return w/o allowing
				// a container event to be propogated to subscribers
				return
			}
			log.Debugf("Container(%s) %s via event activity", c, newState)
			// if we are here the containerVM has been removed from vSphere, so lets remove it
			// from the portLayer cache
			Containers.Remove(c.ExecConfig.ID)
			c.vm = nil
		default:
			return
		}

		// regardless of state update success or failure publish the container event
		publishContainerEvent(c.ExecConfig.ID, e.Created(), e.String())
		return
	}

	switch e.String() {
	case events.ContainerRelocated:
		// container relocated so we need to update the container attributes
		ctx, cancel := context.WithTimeout(context.Background(), constants.PropertyCollectorTimeout)
		defer cancel()

		err := c.refresh(ctx)
		if err != nil {
			log.Errorf("Event driven container update failed for %s with %s", c, err)
		}
	}
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
		vm := vm.NewVirtualMachine(ctx, sess, v.Reference())
		base := newBase(vm, v.Config, &v.Runtime)
		c := newContainer(base)

		id := uid.Parse(c.ExecConfig.ID)
		if id == uid.NilUID {
			log.Warnf("skipping converting container VM %s: could not parse id", v.Reference())
			continue
		}

		if v.Summary.Storage != nil {
			c.VMUnsharedDisk = v.Summary.Storage.Unshared
		}

		cons = append(cons, c)
	}

	return cons
}
