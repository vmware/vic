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
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/lib/portlayer/event/collector/vsphere"
	"github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var containerEvents []events.Event

func TestEventedState(t *testing.T) {
	// poweredOn event
	event := events.ContainerPoweredOn
	assert.EqualValues(t, StateStarting, eventedState(event, StateStarting))
	assert.EqualValues(t, StateRunning, eventedState(event, StateRunning))
	assert.EqualValues(t, StateRunning, eventedState(event, StateStopped))
	assert.EqualValues(t, StateRunning, eventedState(event, StateSuspended))

	// powerOff event
	event = events.ContainerPoweredOff
	assert.EqualValues(t, StateStopping, eventedState(event, StateStopping))
	assert.EqualValues(t, StateStopped, eventedState(event, StateStopped))
	assert.EqualValues(t, StateStopped, eventedState(event, StateRunning))

	// suspended event
	event = events.ContainerSuspended
	assert.EqualValues(t, StateSuspending, eventedState(event, StateSuspending))
	assert.EqualValues(t, StateSuspended, eventedState(event, StateSuspended))
	assert.EqualValues(t, StateSuspended, eventedState(event, StateRunning))

	// removed event
	event = events.ContainerRemoved
	assert.EqualValues(t, StateRemoved, eventedState(event, StateRunning))
	assert.EqualValues(t, StateRemoved, eventedState(event, StateStopped))
	assert.EqualValues(t, StateRemoving, eventedState(event, StateRemoving))
}

func TestPublishContainerEvent(t *testing.T) {

	NewContainerCache()
	containerEvents = make([]events.Event, 0)
	Config = Configuration{}

	mgr := event.NewEventManager()
	Config.EventManager = mgr
	mgr.Subscribe(events.NewEventType(events.ContainerEvent{}).Topic(), "testing", containerCallback)

	// create new running container and place in cache
	id := "123439"
	container := newTestContainer(id)
	addTestVM(container)
	container.SetState(StateRunning)
	Containers.Put(container)

	publishContainerEvent(id, time.Now().UTC(), events.ContainerPoweredOff)
	time.Sleep(time.Millisecond * 30)

	assert.Equal(t, 1, len(containerEvents))
	assert.Equal(t, id, containerEvents[0].Reference())
	assert.Equal(t, events.ContainerPoweredOff, containerEvents[0].String())
}

func containerCallback(ee events.Event) {
	containerEvents = append(containerEvents, ee)
}

func TestRegisteredEvent(t *testing.T) {
	//      target := fmt.Sprintf("https://%s:%s@%s", "administrator@vsphere.local", "Alfred!23", "bart.eng.vmware.com")

	//      ctx := context.Background()

	//      sessionconfig := &session.Config{
	//              Service:        target,
	//              Insecure:       true,
	//              DatacenterPath: "stressdatacenter",
	//              ClusterPath:    "/stressdatacenter/host/cls",
	//              DatastorePath:  "/stressdatacenter/datastore/vsanDatastore",
	//      }
	log.SetLevel(log.DebugLevel)
	target := fmt.Sprintf("https://%s:%s@%s", "root", "welc0ME", "10.17.109.138")

	ctx := context.Background()

	sessionconfig := &session.Config{
		Service:        target,
		Insecure:       true,
		DatacenterPath: "ha-datacenter",
		ClusterPath:    "/ha-datacenter/host/localhost.eng.vmware.com",
		DatastorePath:  "/ha-datacenter/datastore/datastore1",
		PoolPath:       "/ha-datacenter/host/localhost.eng.vmware.com/Resources/virtual-container-host",
	}

	installSession, err := session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		t.Errorf("Failed to create session: %s", err)
	}

	_, err = installSession.Connect(ctx)
	if err != nil {
		t.Errorf("Failed to connect session: %s", err)
	}
	if _, err = installSession.Populate(ctx); err != nil {
		t.Errorf("Failed to get resources: %s", err)
	}

	// we want to monitor the cluster, so create a vSphere Event Collector
	// The cluster managed object will either be a proper vSphere Cluster or
	// a specific host when standalone mode
	ec := vsphere.NewCollector(installSession.Vim25(), installSession.Cluster.Reference().String())

	// start the collection of vsphere events
	err = ec.Start()
	if err != nil {
		err = fmt.Errorf("%s failed to start: %s", ec.Name(), err)
		t.Error(err)
		return
	}

	eventSession = installSession
	Config.ResourcePool = installSession.Pool
	// create the event manager &  register the existing collector
	EventManager := event.NewEventManager(ec)

	// subscribe the exec layer to the event stream for Vm events
	EventManager.Subscribe(events.NewEventType(vsphere.VMEvent{}).Topic(), "exec", eventCallback)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {

		select {
		case <-ticker.C:
			t.Log("ticker tick")
		}
	}
	t.Log("test finished")
}
