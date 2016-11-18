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
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/lib/portlayer/event/collector/vsphere"
	"github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/pkg/backoff"
	"github.com/vmware/vic/pkg/uid"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var (
	initializer struct {
		err  error
		once sync.Once
	}
)

// ConcurrentAccessError is returned when concurrent calls tries to modify same object
type retryError struct {
	err error
}

func (r retryError) Error() string {
	return r.err.Error()
}

func Init(ctx context.Context, sess *session.Session, source extraconfig.DataSource, _ extraconfig.DataSink) error {
	initializer.once.Do(func() {
		var err error
		defer func() {
			if err != nil {
				initializer.err = err
			}
		}()
		f := find.NewFinder(sess.Vim25(), false)

		extraconfig.Decode(source, &Config)

		log.Debugf("Decoded VCH config for execution: %#v", Config)
		ccount := len(Config.ComputeResources)
		if ccount != 1 {
			err = fmt.Errorf("expected singular compute resource element, found %d", ccount)
			log.Error(err)
			return
		}

		cr := Config.ComputeResources[0]
		var r object.Reference
		r, err = f.ObjectReference(ctx, cr)
		if err != nil {
			err = fmt.Errorf("could not get resource pool or virtual app reference from %q: %s", cr.String(), err)
			log.Error(err)
			return
		}
		switch o := r.(type) {
		case *object.VirtualApp:
			Config.VirtualApp = o
			Config.ResourcePool = o.ResourcePool
		case *object.ResourcePool:
			Config.ResourcePool = o
		default:
			err = fmt.Errorf("could not get resource pool or virtual app from reference %q: object type is wrong", cr.String())
			log.Error(err)
			return
		}

		// we want to monitor the cluster, so create a vSphere Event Collector
		// The cluster managed object will either be a proper vSphere Cluster or
		// a specific host when standalone mode
		ec := vsphere.NewCollector(sess.Vim25(), sess.Cluster.Reference().String())

		// start the collection of vsphere events
		err = ec.Start()
		if err != nil {
			err = fmt.Errorf("%s failed to start: %s", ec.Name(), err)
			log.Error(err)
			return
		}

		// create the event manager &  register the existing collector
		Config.EventManager = event.NewEventManager(ec)

		// subscribe the exec layer to the event stream for Vm events
		Config.EventManager.Subscribe(events.NewEventType(vsphere.VMEvent{}).Topic(), "exec", eventCallback)
		// subscribe callback to handle vm registered event
		Config.EventManager.Subscribe(events.NewEventType(vsphere.VMEvent{}).Topic(), "registeredVMEvent", func(ie events.Event) {
			registeredVMCallback(sess, ie)
		})
		// subscribe callback to add container stopTime if not set by tether
		Config.EventManager.Subscribe(events.NewEventType(events.ContainerEvent{}).Topic(), "containerTime", containerTimeCallback)
		// instantiate the container cache now
		NewContainerCache()

		// Grab the AboutInfo about our host environment
		about := sess.Vim25().ServiceContent.About
		Config.VCHMhz = NCPU(ctx)
		Config.VCHMemoryLimit = MemTotal(ctx)
		Config.HostOS = about.OsType
		Config.HostOSVersion = about.Version
		Config.HostProductName = about.Name
		log.Debugf("Host - OS (%s), version (%s), name (%s)", about.OsType, about.Version, about.Name)
		log.Debugf("VCH limits - %d Mhz, %d MB", Config.VCHMhz, Config.VCHMemoryLimit)

		// sync container cache
		if err = Containers.sync(ctx, sess); err != nil {
			return
		}
	})
	return initializer.err
}

// eventCallback will process events
func eventCallback(ie events.Event) {
	// grab the container from the cache
	container := Containers.Container(ie.Reference())
	if container != nil {

		newState := eventedState(ie.String(), container.CurrentState())
		// do we have a state change
		if newState != container.CurrentState() {
			switch newState {
			case StateStopped:

				log.Debugf("Container(%s) state set to %s via event activity", container.ExecConfig.ID, newState.String())

				container.SetState(newState)
				container.onStop()
				// update container state if stoptime is not set by tether
				ctx, cancel := context.WithTimeout(context.Background(), propertyCollectorTimeout)
				defer cancel()

				if err := setStopTime(ctx, ie.Reference()); err != nil {
					log.Errorf("Event driven container update failed: %s", err)
				}
				// regardless of update success failure publish the container event
				publishContainerEvent(container.ExecConfig.ID, ie.Created(), ie.String())
			case StateStopping,
				StateRunning,
				StateSuspended:

				log.Debugf("Container(%s) state set to %s via event activity", container.ExecConfig.ID, newState.String())

				container.SetState(newState)
				if newState == StateStopped {
					container.onStop()
				}

				// container state has changed so we need to update the container attributes
				// we'll do this in a go routine to avoid blocking
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), propertyCollectorTimeout)
					defer cancel()
					err := container.Refresh(ctx)
					if err != nil {
						log.Errorf("Event driven container update failed: %s", err.Error())
					}
					// regardless of update success failure publish the container event
					publishContainerEvent(container.ExecConfig.ID, ie.Created(), ie.String())
				}()
			case StateRemoved:
				log.Debugf("Container(%s) %s via event activity", container.ExecConfig.ID, newState.String())
				Containers.Remove(container.ExecConfig.ID)
				publishContainerEvent(container.ExecConfig.ID, ie.Created(), ie.String())

			}
		}
	}
	return
}

func waitForContainerStarted(container *Container) bool {
	ctx, cancel := context.WithTimeout(context.Background(), propertyCollectorTimeout)
	defer cancel()

	if container.vm == nil {
		log.Errorf("Container(%s) wait failed for vm is not found", container.ExecConfig.ID)
		return false
	}
	log.Debugf("Waiting for container(%s) started", container.ExecConfig.ID)
	for _, k := range extraconfig.CalculateKeys(container.ExecConfig, "Sessions.*.Detail.StartTime", "") {
		var poweredoff bool
		var detail string
		waitFunc := func(pc []types.PropertyChange) bool {
			for _, c := range pc {
				if c.Op != types.PropertyChangeOpAssign {
					continue
				}

				switch v := c.Val.(type) {
				case types.ArrayOfOptionValue:
					for _, value := range v.OptionValue {
						// check the status of the key and return true if it's been set to non-nil
						if k == value.GetOptionValue().Key {
							detail = value.GetOptionValue().Value.(string)
							if detail != "0" {
								return true
							}
							break // continue the outer loop as we may have a powerState change too
						}
					}
				case types.VirtualMachinePowerState:
					if v != types.VirtualMachinePowerStatePoweredOn {
						// Give up if the vm has powered off
						poweredoff = true
						return true
					}
				}

			}
			return false
		}
		log.Debugf("Waiting for container(%s) key %s", container.ExecConfig.ID, k)

		err := container.vm.WaitForExtraConfig(ctx, waitFunc)
		if poweredoff {
			log.Errorf("Failed to wait container started for vm is powered off")
			return false
		}
		if err != nil {
			log.Errorf("Failed to wait container started: %s", err)
		}
	}
	return true
}

func setStopTime(ctx context.Context, id string) error {
	operation := func() error {
		log.Debugf("Updating container(%s) stopTime", id)
		h := GetContainer(ctx, uid.Parse(id))
		if h == nil {
			log.Errorf("Failed to get handle of container(%s), stop retry", id)
			return nil
		}
		defer h.Close()

		if h.vm == nil {
			log.Errorf("Container(%s) update failed for vm is not found", h.ExecConfig.ID)
			return nil
		}

		var update bool
		if h.ExecConfig == nil {
			return &retryError{fmt.Errorf("Failed to get container(%s) configuration, retry later", h.ExecConfig.ID)}
		}
		if h.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOff {
			log.Debugf("container(%s) power state is changed to %s, stop configuration updating for powered off event", h.ExecConfig.ID, h.Runtime.PowerState)
			return nil
		}
		for _, sc := range h.ExecConfig.Sessions {
			if sc.StopTime == 0 {
				sc.StopTime = time.Now().UTC().Unix()
				sc.StartTime = 0
				log.Debugf("Container(%s) session %s stop time is set to %s", h.ExecConfig.ID, sc.ID, time.Unix(sc.StopTime, 0))
				update = true
			} else {
				log.Debugf("Container(%s) session %s stop time is already set to %s", h.ExecConfig.ID, sc.ID, time.Unix(sc.StopTime, 0))
			}
		}
		if !update {
			log.Debugf("No need to update container(%s) configuration", h.ExecConfig.ID)
			return nil
		}
		log.Debugf("commit container(%s)", h.ExecConfig.ID)
		// set waittime to 0 for vm is already stopped, no need to wait again
		return h.Commit(ctx, h.vm.Session, new(int32))
	}
	retryOnError := func(err error) bool {
		switch err.(type) {
		case ConcurrentAccessError, retryError:
			log.Debugf("Retry container commit")
			return true
		default:
			return false
		}
	}
	return backoff.Retry(operation, retryOnError)
}

func containerTimeCallback(ie events.Event) {
	container := Containers.Container(ie.Reference())
	if container == nil {
		log.Debugf("Container(%s) is not found", ie.Reference())
	}
	go func() {
		switch ie.String() {
		case events.ContainerPoweredOn:
			// container powered on event received, still need to wait container startTime set by tether, to sync between portlayer and vmx
			// we'll do this in a go routine to avoid blocking
			if !waitForContainerStarted(container) {
				return
			}
			fallthrough
		case events.ContainerStarted:
			// container started event received, refresh diretly
			// we might have timeout to wait container started, so here reset timeout to refresh
			log.Debugf("Container(%s) started, refresh", container.ExecConfig.ID)
			ctx, cancel := context.WithTimeout(context.Background(), propertyCollectorTimeout)
			defer cancel()
			err := container.Refresh(ctx)
			if err != nil {
				log.Errorf("Event driven container update failed: %s", err.Error())
			}
		}
	}()

	return
}

// registeredVMCallback will process registeredVMEvent
func registeredVMCallback(sess *session.Session, ie events.Event) {
	// check container registered event if this container is not found in container cache
	// grab the container from the cache
	container := Containers.Container(ie.Reference())
	if container != nil {
		// if container exists, ingore it
		return
	}
	switch ie.String() {
	case events.ContainerRegistered:
		moref := new(types.ManagedObjectReference)
		if ok := moref.FromString(ie.Reference()); !ok {
			log.Errorf("Failed to get event VM mobref: %s", ie.Reference())
			return
		}
		if !isManagedbyVCH(sess, *moref) {
			return
		}
		log.Debugf("Register container VM %s", moref)
		ctx := context.Background()
		vms, err := populateVMAttributes(ctx, sess, []types.ManagedObjectReference{*moref})
		if err != nil {
			log.Error(err)
			return
		}
		registeredContainers := convertInfraContainers(ctx, sess, vms)
		for i := range registeredContainers {
			Containers.put(registeredContainers[i])
			log.Debugf("Registered container %q", registeredContainers[i].Config.Name)
		}
	}
	return
}

func isManagedbyVCH(sess *session.Session, moref types.ManagedObjectReference) bool {
	var vm mo.VirtualMachine

	// current attributes we care about
	attrib := []string{"resourcePool", "config.name"}

	// populate the vm properties
	ctx := context.Background()
	if err := sess.RetrieveOne(ctx, moref, attrib, &vm); err != nil {
		log.Errorf("Failed to query registered vm object %s: %s", moref.String(), err)
		return false
	}
	if *vm.ResourcePool != Config.ResourcePool.Reference() {
		log.Debugf("container vm %q does not belong to this VCH, ignoring", vm.Config.Name)
		return false
	}
	return true
}

// eventedState will determine the target container
// state based on the current container state and the vsphere event
func eventedState(e string, current State) State {
	switch e {
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

// publishContainerEvent will publish a ContainerEvent to the vic event stream
func publishContainerEvent(id string, created time.Time, eventType string) {
	if Config.EventManager == nil || eventType == "" {
		return
	}

	ce := &events.ContainerEvent{
		BaseEvent: &events.BaseEvent{
			Ref:         id,
			CreatedTime: created,
			Event:       eventType,
			Detail:      fmt.Sprintf("Container %s %s", id, eventType),
		},
	}

	Config.EventManager.Publish(ce)
}
