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
	"time"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"

	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/lib/portlayer/event/collector/vsphere"
	"github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"

	"github.com/vmware/vic/pkg/vsphere/session"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

func Init(ctx context.Context, f *find.Finder, session *session.Session) error {

	ccount := len(VCHConfig.ComputeResources)
	if ccount != 1 {
		detail := fmt.Sprintf("expected singular compute resource element, found %d", ccount)
		log.Errorf(detail)
		return errors.New(detail)
	}

	cr := VCHConfig.ComputeResources[0]

	r, err := f.ObjectReference(ctx, cr)
	if err != nil {
		detail := fmt.Sprintf("could not get resource pool or virtual app reference from %q: %s", cr.String(), err)
		log.Errorf(detail)
		return err
	}
	switch o := r.(type) {
	case *object.VirtualApp:
		VCHConfig.VirtualApp = o
		VCHConfig.ResourcePool = o.ResourcePool
	case *object.ResourcePool:
		VCHConfig.ResourcePool = o
	default:
		detail := fmt.Sprintf("could not get resource pool or virtual app from reference %q: object type is wrong", cr.String())
		log.Errorf(detail)
		return errors.New(detail)
	}

	// we want to monitor the resource pool, so create a vSphere Event Collector
	ec := vsphere.NewCollector(session.Vim25(), VCHConfig.ResourcePool.Reference().String())

	// start the collection of vsphere events
	err = ec.Start()
	if err != nil {
		detail := fmt.Sprintf("%s failed to start: %s", ec.Name(), err.Error())
		log.Error(detail)
		return err
	}

	// create the event manager &  register the existing collector
	VCHConfig.EventManager = event.NewEventManager(ec)

	// subscribe the exec layer to the event stream for Vm events
	VCHConfig.EventManager.Subscribe(events.NewEventType(vsphere.VmEvent{}).Topic(), "exec", eventCallback)

	// instantiate the container cache now
	NewContainerCache()

	// Grab the AboutInfo about our host environment
	about := session.Vim25().ServiceContent.About
	VCHConfig.VCHMhz = NCPU(ctx)
	VCHConfig.VCHMemoryLimit = MemTotal(ctx)
	VCHConfig.HostOS = about.OsType
	VCHConfig.HostOSVersion = about.Version
	VCHConfig.HostProductName = about.Name
	log.Debugf("Host - OS (%s), version (%s), name (%s)", about.OsType, about.Version, about.Name)
	log.Debugf("VCH limits - %d Mhz, %d MB", VCHConfig.VCHMhz, VCHConfig.VCHMemoryLimit)

	return nil
}

// eventCallback will process events
func eventCallback(ie events.Event) {
	// grab the container from the cache
	container := containers.Container(ie.Reference())
	if container != nil {

		newState := eventedState(ie.String(), container.State)
		// do we have a state change
		if newState != container.State {
			switch newState {
			case StateStopping,
				StateRunning,
				StateStopped,
				StateSuspended:

				log.Debugf("Container(%s) state set to %s via event activity", container.ExecConfig.ID, newState.String())
				container.State = newState
				// container state has changed so we need to update the container attributes
				// we'll do this in a go routine to avoid blocking
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), propertyCollectorTimeout)
					defer cancel()

					_, err := container.Update(ctx, container.vm.Session)
					if err != nil {
						log.Errorf("Event driven container update failed: %s", err.Error())
					}
					// regardless of update success failure publish the container event
					publishContainerEvent(container.ExecConfig.ID, ie.Created(), ie.String())
				}()
			case StateRemoved:
				log.Debugf("Container(%s) %s via event activity", container.ExecConfig.ID, newState.String())
				containers.Remove(container.ExecConfig.ID)
				publishContainerEvent(container.ExecConfig.ID, ie.Created(), ie.String())

			}
		}
	}

	return
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
	if VCHConfig.EventManager == nil || eventType == "" {
		return
	}

	ce := &events.ContainerEvent{
		&events.BaseEvent{
			Ref:         id,
			CreatedTime: created,
			Event:       eventType,
			Detail:      fmt.Sprintf("Container %s %s", id, eventType),
		},
	}

	VCHConfig.EventManager.Publish(ce)

}

func WaitForContainerStop(ctx context.Context, id string) error {
	defer trace.End(trace.Begin(id))

	listen := make(chan interface{})
	defer close(listen)

	watch := func(ce events.Event) {
		event := ce.String()
		if ce.Reference() == id {
			switch event {
			case events.ContainerStopped,
				events.ContainerPoweredOff:
				listen <- event
			}
		}
	}

	sub := fmt.Sprintf("%s:%s(%d)", id, "watcher", &watch)
	topic := events.NewEventType(events.ContainerEvent{}).Topic()
	VCHConfig.EventManager.Subscribe(topic, sub, watch)
	defer VCHConfig.EventManager.Unsubscribe(topic, sub)

	// wait for the event to be pushed on the channel or
	// the context to be complete
	select {
	case <-listen:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("WaitForContainerStop(%s) Error: %s", id, ctx.Err())
	}
}
