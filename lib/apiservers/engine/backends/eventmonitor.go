// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package backends

//**** eventmonitor.go
//
// Handles monitoring of events from the portlayer.  Events that are applicable to
// Docker events are then translated and published to the Docker event subscribers.
// NOTE:  This does not handle all Docker events.  In fact, most docker events are
// passively handled by API calls in the backend routers, with no feedback from
// the portlayer.

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	eventtypes "github.com/docker/docker/api/types/events"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/events"
	plevents "github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/pkg/trace"
)

const (
	containerDieEvent     = "die"
	containerDestroyEvent = "destroy"
	containerStopEvent    = "stop"
	containerStartEvent   = "start"
	containerCreateEvent  = "create"
	containerRestartEvent = "restart"
	containerAttachEvent  = "attach"
	containerDetachEvent  = "detach"
	containerKillEvent    = "kill"
	containerResizeEvent  = "resize"
)

// for unit testing purposes
type eventproxy interface {
	StreamEvents(ctx context.Context, out io.Writer) error
}

type eventpublisher interface {
	PublishEvent(event plevents.BaseEvent)
}

type PlEventProxy struct {
}

type DockerEventPublisher struct {
}

type PortlayerEventMonitor struct {
	stop      chan struct{}
	proxy     eventproxy
	publisher eventpublisher
}

// StreamEvents() handles all swagger interaction to the Portlayer's event manager
//
// Input:
//	context and a io.Writer
func (ep PlEventProxy) StreamEvents(ctx context.Context, out io.Writer) error {
	defer trace.End(trace.Begin(""))

	plClient := PortLayerClient()
	if plClient == nil {
		return InternalServerError("eventproxy.StreamEvents failed to get a portlayer client")
	}

	params := events.NewGetEventsParamsWithContext(ctx)
	if _, err := plClient.Events.GetEvents(params, out); err != nil {
		switch err := err.(type) {
		case *events.GetEventsInternalServerError:
			return InternalServerError("Server error from the events port layer")
		default:
			//Check for EOF.  Since the connection, transport, and data handling are
			//encapsulated inside of Swagger, we can only detect EOF by checking the
			//error string
			if strings.Contains(err.Error(), swaggerSubstringEOF) {
				return nil
			}
			return InternalServerError(fmt.Sprintf("Unknown error from the interaction port layer: %s", err))
		}
	}

	return nil
}

func NewPortlayerEventMonitor(proxy eventproxy, publisher eventpublisher) *PortlayerEventMonitor {
	return &PortlayerEventMonitor{proxy: proxy, publisher: publisher}
}

// Start() starts the portlayer event monitoring
func (m *PortlayerEventMonitor) Start() error {
	defer trace.End(trace.Begin(""))

	if m.stop != nil {
		return fmt.Errorf("Portlayer event monitor: Already started")
	}

	m.stop = make(chan struct{})
	go m.monitor()

	return nil
}

// Stop() stops the portlayer event monitoring
func (m *PortlayerEventMonitor) Stop() {
	defer trace.End(trace.Begin(""))

	if m.stop != nil {
		close(m.stop)
	}
}

// monitor() establishes a streaming connection to the portlayer's event
// endpoint, decodes the results, translate it to Docker events if needed,
// and publishes the event to Docker event subscribers.
func (m *PortlayerEventMonitor) monitor() error {
	defer trace.End(trace.Begin(""))

	var wg sync.WaitGroup
	errors := make(chan error, 2)

	reader, writer := io.Pipe()
	ctx, cancel := context.WithCancel(context.TODO())
	// Start streaming events
	wg.Add(1)
	go func() {
		var err error

		defer wg.Done()

		if err = m.proxy.StreamEvents(ctx, writer); err != nil {
			if ctx.Err() != context.Canceled {
				log.Warnf("Event streaming from portlayer returned: %#v", err)
			}
		}
		if ctx.Err() == context.Canceled {
			log.Infof("Event streaming from portlayer was cancelled")
			return
		}
		errors <- err

		writer.Close()
		reader.Close()
	}()

	// Start decoding event stream json
	wg.Add(1)
	go func() {
		var err error
		var event plevents.BaseEvent

		defer wg.Done()

		decoder := json.NewDecoder(reader)
		for decoder.More() {
			if err = decoder.Decode(&event); err == nil {
				m.publisher.PublishEvent(event)
			}
		}
		errors <- err

		reader.Close()
		writer.Close()
	}()

	// Create a channel signaling when the waitgroup finishes
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(errors)
		close(done)
	}()

	select {
	case <-done:
		for err := range errors {
			if err != nil {
				log.Warnf("Exiting Events Monitor: %#v", err)
				return err
			}
		}
	case <-m.stop:
		cancel()
		writer.Close()
		reader.Close()
	}

	return nil
}

// publishEvent() translate a portlayer event into a Docker event if the event is for a
// known container and publishes them to Docker event subscribers.
func (p DockerEventPublisher) PublishEvent(event plevents.BaseEvent) {
	defer trace.End(trace.Begin(""))

	vc := cache.ContainerCache().GetContainer(event.Ref)
	if vc == nil {
		log.Warnf("Portlayer event for container %s but not found in cache", event.Ref)
		return
	}

	var attrs map[string]string
	// TODO: move to a container.OnEvent() so that container drives the necessary changes
	// based on event activity
	switch event.Event {
	case plevents.ContainerStopped,
		plevents.ContainerPoweredOff:
		// since we are going to make a call to the portLayer lets execute this in a
		// go routine
		go func() {
			attrs := make(map[string]string)
			// get the containerEngine
			code, err := NewContainerBackend().containerProxy.exitCode(vc)
			if err != nil {
				// log the error, but continue
				log.Errorf("unable to get exitCode for die event: %s", err)
			}
			// if the docker client is unable to convert the code to an int then
			// then the client will return 125
			attrs["exitCode"] = code
			log.Infof("Sending die event for container %s - code: %s", vc.ContainerID, code)
			actor := CreateContainerEventActorWithAttributes(vc, attrs)
			EventService().Log(containerDieEvent, eventtypes.ContainerEventType, actor)

			// auto-remove if required
			if vc.HostConfig.AutoRemove {
				config := &types.ContainerRmConfig{
					ForceRemove:  true,
					RemoveVolume: true,
				}

				err := NewContainerBackend().ContainerRm(vc.Name, config)
				if err != nil {
					log.Warnf("failed to auto remove container: %s", err.Error())
				}
			}

		}()
	case plevents.ContainerRemoved:
		attrs = make(map[string]string)
		//pop the destroy event...
		actor := CreateContainerEventActorWithAttributes(vc, attrs)
		EventService().Log(containerDestroyEvent, eventtypes.ContainerEventType, actor)
		if err := UnmapPorts(vc.HostConfig); err != nil {
			log.Warn(err)
		}
		// remove from the container cache...
		cache.ContainerCache().DeleteContainer(vc.ContainerID)
	default:
		// let everything else slide on by...
	}

}
