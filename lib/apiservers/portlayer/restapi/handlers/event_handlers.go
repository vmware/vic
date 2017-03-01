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

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/events"
	ple "github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
)

// EventHandlerImpl is the receiver for all of the event handler methods
type EventsHandlerImpl struct {
	handlerCtx *HandlerContext
}

// Configure assigns functions to all the exec api handlers
func (handler *EventsHandlerImpl) Configure(api *operations.PortLayerAPI, handlerCtx *HandlerContext) {
	api.EventsGetEventsHandler = events.GetEventsHandlerFunc(handler.GetEventsHandler)
	handler.handlerCtx = handlerCtx
}

// GetEventsHandler provides a stream of events
func (handler *EventsHandlerImpl) GetEventsHandler(params events.GetEventsParams) middleware.Responder {
	defer trace.End(trace.Begin(""))

	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	flusher := NewFlushingReader(r)

	// uid for subscription
	id := uid.New().String()
	sub := fmt.Sprintf("%s-%s", "PLE", id)

	// currently only containerEvents will be streamed
	topic := ple.NewEventType(ple.ContainerEvent{}).Topic()

	// func to clean up the event stream
	onClose := func() {
		exec.Config.EventManager.Unsubscribe(topic, sub)
		closePipe(r, w)
	}

	// subscribe to event stream
	exec.Config.EventManager.Subscribe(topic, sub, func(ie ple.Event) {
		err := enc.Encode(ie)
		if err != nil {
			log.Errorf("Encoding Error: %s", err.Error())
			exec.Config.EventManager.Unsubscribe(topic, sub)
			closePipe(r, w)
		}
	})

	outHandler := &EventOutputHandler{
		outputName:  "events",
		onHTTPClose: onClose,
	}

	return outHandler.WithPayload(flusher)
}

// closePipe is a convenience function for closing the event stream pipe
func closePipe(pipeReader *io.PipeReader, pipeWriter *io.PipeWriter) {
	if pipeReader != nil {
		pipeReader.Close()
	}
	if pipeWriter != nil {
		pipeWriter.Close()
	}
}

// EventOutputHandler is custom return handlers for event stream
type EventOutputHandler struct {
	outputStream *FlushingReader
	outputName   string
	onHTTPClose  func()
}

// NewEventOutputHandler creates EventOutputHandler with default headers values
func NewEventOutputHandler(name string) *EventOutputHandler {
	return &EventOutputHandler{outputName: name}
}

// WithPayload adds the payload to output handler
func (c *EventOutputHandler) WithPayload(payload *FlushingReader) *EventOutputHandler {
	c.outputStream = payload
	return c
}

// WriteResponse to the client
func (c *EventOutputHandler) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusOK)
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
		c.outputStream.AddFlusher(f)
	}

	notify := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		// execute cleanup function
		c.onHTTPClose()
	}()

	_, err := io.Copy(rw, c.outputStream)
	if err != nil {
		log.Debugf("Error copying %s stream: %s", c.outputName, err)
	} else {
		log.Debugf("Finished copying %s stream", c.outputName)
	}
}
