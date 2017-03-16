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
	"fmt"
	"sync"
	"time"

	"context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/lib/portlayer/event/collector/vsphere"
	"github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/pkg/retry"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var (
	initializer struct {
		err  error
		once sync.Once
	}
)

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

		// instantiate the container cache now
		NewContainerCache()

		// create the event manager &  register the existing collector
		Config.EventManager = event.NewEventManager(ec)

		// subscribe to host events
		Config.EventManager.Subscribe(events.NewEventType(vsphere.HostEvent{}).Topic(), "host",
			func(ie events.Event) {
				hostEventCallback(ctx, ie)
			},
		)
		// subscribe the exec layer to the event stream for Vm events
		vmSub := Config.EventManager.Subscribe(events.NewEventType(vsphere.VMEvent{}).Topic(), "exec", func(e events.Event) {
			if c := Containers.Container(e.Reference()); c != nil {
				c.OnEvent(e)
			}
		})

		// Grab the AboutInfo about our host environment
		about := sess.Vim25().ServiceContent.About

		p, err := GetVCHstats(ctx)
		if err != nil {
			log.Errorf("Failed to get VCH stats: %s", err)
		} else {
			Config.VCHMhz = p.Config.CpuAllocation.GetResourceAllocationInfo().Limit
			Config.VCHMemoryLimit = p.Config.MemoryAllocation.GetResourceAllocationInfo().Limit
		}
		Config.HostOS = about.OsType
		Config.HostOSVersion = about.Version
		Config.HostProductName = about.Name
		log.Debugf("Host - OS (%s), version (%s), name (%s)", about.OsType, about.Version, about.Name)
		log.Debugf("VCH limits - %d Mhz, %d MB", Config.VCHMhz, Config.VCHMemoryLimit)

		// sync container cache
		vmSub.Suspend(true)
		defer vmSub.Resume()
		if err = Containers.sync(ctx, sess); err != nil {
			return
		}
	})

	return initializer.err
}

// listens migrated events and connects the file backed serial ports
func hostEventCallback(ctx context.Context, ie events.Event) {
	defer trace.End(trace.Begin(""))

	switch ie.String() {
	// https://www.vmware.com/support/developer/vc-sdk/visdk25pubs/ReferenceGuide/vim.event.EnteringMaintenanceModeEvent.html
	// This event records that a host has begun the process of entering maintenance mode. All virtual machine operations are blocked, except the following:
	// MigrateVM
	// PowerOffVM
	// SuspendVM
	// ShutdownGuest
	// StandbyGuest
	//
	// Because of that limitation the only thing that we can do is shutting down the guest without calling reconfigure
	case events.HostEnteringMaintenanceMode:
		log.Debugf("Received %s event", ie)
		ref := ie.Reference()

		// we are interested with running vms
		state := new(State)
		*state = StateRunning
		for _, v := range Containers.Containers(state) {
			host := v.Runtime.Host.String()
			if host != ref {
				log.Debugf("Skipping %q as it is not on %q", v, ref)
				continue
			}

			log.Debugf("%q is on %q", v, ref)

			// grab the container from the cache
			container := Containers.Container(v.String())
			if container == nil {
				log.Errorf("Container %s not found", v)
				continue
			}

			operation := func() error {
				var err error

				handle := container.NewHandle(ctx)
				if handle == nil {
					err = fmt.Errorf("Handle for %s cannot be created", v)
					log.Error(err)
					return err
				}
				defer handle.Close()

				// this check needs to be after we get a new handle otherwise we could receive stalled data
				needsToBePoweredOff := false

				// get the virtual device list
				devices := object.VirtualDeviceList(v.Config.Hardware.Device)

				// select the virtual serial ports
				serials := devices.SelectByBackingInfo((*types.VirtualSerialPortURIBackingInfo)(nil))
				log.Debugf("Found %d devices with the desired backing", len(serials))

				// iterate over them and set needsToBePoweredOff if necessary
				for _, serial := range serials {
					needsToBePoweredOff = serial.GetVirtualDevice().Connectable.Connected

					log.Debug("Connected: %t", needsToBePoweredOff)
					if needsToBePoweredOff {
						break
					}
				}
				if !needsToBePoweredOff {
					log.Debugf("Skipping %q. Serial is not connected so it will be migrated", v)
					return nil
				}

				handle.SetTargetState(StateStopped)

				// call CommitWithoutSpec which sets spec to nil
				if err = handle.CommitWithoutSpec(ctx, nil, nil); err != nil {
					log.Errorf("Failed to commit handle after getting %s event for container %s: %s", ie, v, err)
					return err
				}
				return nil
			}

			if err := retry.Do(operation, IsConcurrentAccessError); err != nil {
				log.Errorf("Multiple attempts failed to commit handle after getting %s event for container %s: %s", ie, v, err)
			}
		}

	}
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
