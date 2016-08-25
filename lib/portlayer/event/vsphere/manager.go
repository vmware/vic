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

package vsphere

import (
	"fmt"

	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/pkg/vsphere/session"

	vmwEvents "github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/vim25/types"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type vSphereManager struct {
	vmwManager *vmwEvents.Manager
	callbacks  map[string]func(event.Event, *session.Session)
	blacklist  map[types.ManagedObjectReference]string
	objects    map[string]types.ManagedObjectReference
	session    *session.Session
}

func NewEventManager(s *session.Session) *vSphereManager {
	return &vSphereManager{
		callbacks: make(map[string]func(event.Event, *session.Session)),
		blacklist: make(map[types.ManagedObjectReference]string),
		objects:   make(map[string]types.ManagedObjectReference),
		session:   s,
	}

}

func (mgr *vSphereManager) AddMonitoredObject(ref string) error {
	log.Debugf("EventManager added (%s) to monitor", ref)
	moRef := types.ManagedObjectReference{}
	if !moRef.FromString(ref) {
		return fmt.Errorf("Invalid Managed Object provided for Montioring(%s)", ref)
	}
	mgr.objects[ref] = moRef
	return nil
}

func (mgr *vSphereManager) RemoveMonitoredObject(ref string) {
	delete(mgr.objects, ref)
}

func (mgr *vSphereManager) monitoredObjects() []types.ManagedObjectReference {
	refs := make([]types.ManagedObjectReference, 0, len(mgr.objects))
	for k := range mgr.objects {
		refs = append(refs, mgr.objects[k])
	}
	return refs
}

// Blacklist object from callbacks
func (mgr *vSphereManager) Blacklist(ref string) error {
	moRef := types.ManagedObjectReference{}
	if !moRef.FromString(ref) {
		return fmt.Errorf("Invalid Managed Object provided to Blacklist(%s)", ref)
	}
	mgr.blacklist[moRef] = ref
	return nil
}

// Unblacklist removes object blacklisting and will allow callbacks to occur
func (mgr *vSphereManager) Unblacklist(ref string) {
	moRef := types.ManagedObjectReference{}
	moRef.FromString(ref)
	delete(mgr.blacklist, moRef)
}

// BlacklistCount
func (mgr *vSphereManager) BlacklistCount() int {
	return len(mgr.blacklist)
}

func (mgr *vSphereManager) blacklisted(ref types.ManagedObjectReference) bool {
	if _, ok := mgr.blacklist[ref]; ok {
		log.Debugf("EventManager skipping blacklisted object (%s)", ref.String())
		return true
	}
	return false
}

// Register the callback
func (mgr *vSphereManager) Register(caller string, callback func(event.Event, *session.Session)) {
	log.Debugf("Registering %s with EventManager", caller)
	mgr.callbacks[caller] = callback
}

// Unregister call back
func (mgr *vSphereManager) Unregister(caller string) {
	delete(mgr.callbacks, caller)
}

// Registry eturns map of callbacks
func (mgr *vSphereManager) Registry() map[string]func(event.Event, *session.Session) {
	return mgr.callbacks
}

// RegistryCount returns the callback count
func (mgr *vSphereManager) RegistryCount() int {
	return len(mgr.callbacks)
}

// Stop listening by destroying the manager
// this will be best effort -- as event manager will
// be killed when session is killed
func (mgr *vSphereManager) Stop() {
	_, err := mgr.vmwManager.Destroy(context.Background())
	if err != nil {
		log.Warnf("EventManger failed to destroy the govmomi manager: %s", err.Error())
	}
}

// Start the event listener
func (mgr *vSphereManager) Start() error {
	// array of managed objects
	refs := mgr.monitoredObjects()

	// only continue if objects were added to monitor
	if len(refs) == 0 {
		return fmt.Errorf("EventManager requires at least one ManagedObject for monitoring, currently %d have been added", 0)
	}

	log.Debugf("EventManager starting event listener for %d managed objects...", len(refs))

	// create the govmomi view manager
	mgr.vmwManager = vmwEvents.NewManager(mgr.session.Client.Client)

	// we don't want the event listener to timeout
	ctx := context.Background()

	// vars used for clarity
	// page size of lastest page of events
	// setting to 1 to avoid having to track last event
	// processed -- this may change in future if too chatty

	// TODO: need to create an "event storm" to ensure we can handle
	// a 1 event page
	pageSize := int32(1)
	// bool to follow the stream
	followStream := true
	// don't exceed the govmomi object limit
	force := false

	//TODO: need a proper way to handle failures / status
	go func(page int32, follow bool, ff bool, refs []types.ManagedObjectReference) error {
		// the govmomi event listener can only be configured once per session -- so if it's already listening it
		// will be replaced
		//
		// the manager & listener will be closed with the session
		err := mgr.vmwManager.Events(ctx, refs, 1, followStream, force, func(page []types.BaseEvent) error {
			evented(mgr, page)
			return nil
		})
		// TODO: this will disappear in the ether
		if err != nil {
			log.Debugf("Error configuring event manager %s \n", err.Error())
			return err
		}
		return nil
	}(pageSize, followStream, force, refs)

	return nil
}

// evented will process the event and call registered callbacks
//
// Initial implmentation will only act on certain events -- future implementations
// may provide more flexibilty
func evented(mgr *vSphereManager, page []types.BaseEvent) {

	if mgr.RegistryCount() == 0 {
		log.Debug("EventManager has no callbacks configured...")
	}

	for i := range page {
		var vmwEve event.Event
		// reference to object evented
		var ref types.ManagedObjectReference

		e := page[i].GetEvent()

		// what type of event do we have
		switch page[i].(type) {
		case *types.VmGuestShutdownEvent,
			*types.VmPoweredOnEvent,
			*types.VmPoweredOffEvent,
			*types.VmRemovedEvent,
			*types.VmSuspendedEvent:
			vmwEve = NewVMEvent(page[i])
			ref = e.Vm.Vm
		}

		// if we have event and it's not blacklisted
		if vmwEve != nil && !mgr.blacklisted(ref) {
			action(mgr, vmwEve)
		}
	}

}

// action will execute the callbacks
func action(mgr *vSphereManager, vmwEve event.Event) {
	for caller, callback := range mgr.callbacks {
		log.Debugf("EventManager callling back to %s on %s for Event(%s)", caller, vmwEve.Reference(), vmwEve.String())
		callback(vmwEve, mgr.session)
	}
}
