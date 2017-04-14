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

package handlers

import (
	"context"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/interaction"
	"github.com/vmware/vic/lib/portlayer/attach"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
)

// InteractionHandlersImpl is the receiver for all of the interaction handler methods
type InteractionHandlersImpl struct {
	attachServer *attach.Server
}

const (
	interactionTimeout    time.Duration = 30 * time.Second
	attachStdinInitString               = "v1c#>"

	// in sync with lib/tether/tether_linux.go
	// 115200 bps is 14.4 KB/s so use that
	ioCopyBufferSize = 14 * 1024
)

func (i *InteractionHandlersImpl) Configure(api *operations.PortLayerAPI, _ *HandlerContext) {

	api.InteractionInteractionJoinHandler = interaction.InteractionJoinHandlerFunc(i.JoinHandler)
	api.InteractionInteractionBindHandler = interaction.InteractionBindHandlerFunc(i.BindHandler)
	api.InteractionInteractionUnbindHandler = interaction.InteractionUnbindHandlerFunc(i.UnbindHandler)

	api.InteractionContainerResizeHandler = interaction.ContainerResizeHandlerFunc(i.ContainerResizeHandler)
	api.InteractionContainerSetStdinHandler = interaction.ContainerSetStdinHandlerFunc(i.ContainerSetStdinHandler)
	api.InteractionContainerGetStdoutHandler = interaction.ContainerGetStdoutHandlerFunc(i.ContainerGetStdoutHandler)
	api.InteractionContainerGetStderrHandler = interaction.ContainerGetStderrHandlerFunc(i.ContainerGetStderrHandler)

	api.InteractionContainerCloseStdinHandler = interaction.ContainerCloseStdinHandlerFunc(i.ContainerCloseStdinHandler)

	i.attachServer = attach.NewAttachServer(constants.ManagementHostName, constants.AttachServerPort)

	if err := i.attachServer.Start(false); err != nil {
		log.Fatalf("Attach server unable to start: %s", err)
	}
}

// JoinHandler calls the Join
func (i *InteractionHandlersImpl) JoinHandler(params interaction.InteractionJoinParams) middleware.Responder {
	defer trace.End(trace.Begin(""))

	handle := exec.HandleFromInterface(params.Config.Handle)
	if handle == nil {
		err := &models.Error{Message: "Failed to get the Handle"}
		return interaction.NewInteractionJoinInternalServerError().WithPayload(err)
	}

	handleprime, err := attach.Join(handle)
	if err != nil {
		log.Errorf("%s", err.Error())

		return interaction.NewInteractionJoinInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}
	res := &models.InteractionJoinResponse{
		Handle: exec.ReferenceFromHandle(handleprime),
	}
	return interaction.NewInteractionJoinOK().WithPayload(res)
}

// BindHandler calls the Bind
func (i *InteractionHandlersImpl) BindHandler(params interaction.InteractionBindParams) middleware.Responder {
	defer trace.End(trace.Begin(""))

	handle := exec.HandleFromInterface(params.Config.Handle)
	if handle == nil {
		err := &models.Error{Message: "Failed to get the Handle"}
		return interaction.NewInteractionBindInternalServerError().WithPayload(err)
	}

	handleprime, err := attach.Bind(handle)
	if err != nil {
		log.Errorf("%s", err.Error())

		return interaction.NewInteractionBindInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}

	res := &models.InteractionBindResponse{
		Handle: exec.ReferenceFromHandle(handleprime),
	}
	return interaction.NewInteractionBindOK().WithPayload(res)
}

// UnbindHandler calls the Unbind
func (i *InteractionHandlersImpl) UnbindHandler(params interaction.InteractionUnbindParams) middleware.Responder {
	defer trace.End(trace.Begin(""))

	handle := exec.HandleFromInterface(params.Config.Handle)
	if handle == nil {
		err := &models.Error{Message: "Failed to get the Handle"}
		return interaction.NewInteractionUnbindInternalServerError().WithPayload(err)
	}

	handleprime, err := attach.Unbind(handle)
	if err != nil {
		log.Errorf("%s", err.Error())

		return interaction.NewInteractionUnbindInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}

	res := &models.InteractionUnbindResponse{
		Handle: exec.ReferenceFromHandle(handleprime),
	}
	return interaction.NewInteractionUnbindOK().WithPayload(res)
}

// ContainerResizeHandler calls resize
func (i *InteractionHandlersImpl) ContainerResizeHandler(params interaction.ContainerResizeParams) middleware.Responder {
	// See whether there is an active session to the container
	session, err := i.attachServer.Get(context.Background(), params.ID, 0)
	if err != nil {
		// just note the warning and return, resize requires an active connection
		log.Warnf("No resize connection found (id: %s): %s", params.ID, err)

		return interaction.NewContainerResizeOK()
	}

	// Request a resize
	cWidth := uint32(params.Width)
	cHeight := uint32(params.Height)

	if err = session.Resize(cWidth, cHeight, 0, 0); err != nil {
		log.Errorf("%s", err.Error())

		return interaction.NewContainerResizeInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}

	return interaction.NewContainerResizeOK()
}

// ContainerSetStdinHandler returns the stdin
func (i *InteractionHandlersImpl) ContainerSetStdinHandler(params interaction.ContainerSetStdinParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	var ctxDeadline time.Time
	var timeout time.Duration

	// Calculate the timeout for the attach if the caller specified a deadline.  This deadline
	if params.Deadline != nil {
		ctxDeadline = time.Time(*params.Deadline)
		timeout = ctxDeadline.Sub(time.Now())
		log.Debugf("Attempting to get ssh session for container %s stdin with deadline %s", params.ID, ctxDeadline.Format(time.UnixDate))
		if timeout < 0 {
			e := &models.Error{Message: fmt.Sprintf("Deadline for stdin already passed for container %s", params.ID)}
			return interaction.NewContainerSetStdinInternalServerError().WithPayload(e)
		}
	} else {
		log.Debugf("Attempting to get ssh session for container %s stdin", params.ID)
		timeout = interactionTimeout
	}

	session, err := i.attachServer.Get(context.Background(), params.ID, timeout)
	if err != nil {
		log.Errorf("%s", err.Error())

		e := &models.Error{
			Message: fmt.Sprintf("No stdin connection found (id: %s): %s", params.ID, err.Error()),
		}
		return interaction.NewContainerSetStdinNotFound().WithPayload(e)
	}
	// Remove the connection from the map
	defer func() {
		// io.EOF is expected if the channel is already closed so ignore it
		if err := i.attachServer.Remove(params.ID); err != nil && err != io.EOF {
			log.Errorf("Removing the connection from the map failed with %s", err)
		}
	}()

	detachableIn := NewFlushingReaderWithInitBytes(params.RawStream, []byte(attachStdinInitString))
	_, err = io.Copy(session.Stdin(), detachableIn)
	if err != nil {
		log.Errorf("Copy@ContainerSetStdinHandler returned %s", err.Error())
		/*
			// FIXME(caglar10ur): need a way to differentiate detach from pipe
			// Close the stdin if we get an EOF in the middle of the stream
			if err == io.ErrUnexpectedEOF {
				if err = session.CloseStdin(); err != nil {
					log.Errorf("CloseStdin@ContainerSetStdinHandler failed with %s", err.Error())
				} else {
					log.Infof("CloseStdin@ContainerSetStdinHandler succeeded")
				}
			}
		*/

		// FIXME(caglar10ur): Do not return an error here - https://github.com/vmware/vic/issues/2594
		/*
			e := &models.Error{
				Message: fmt.Sprintf("Error copying stdin (id: %s): %s", params.ID, err.Error()),
			}
			return interaction.NewContainerSetStdinInternalServerError().WithPayload(e)
		*/
	}

	log.Debugf("Done copying stdin")

	return interaction.NewContainerSetStdinOK()
}

// ContainerCloseStdinHandler closes the stdin, it returns an error if there is no active connection between portlayer and the tether
func (i *InteractionHandlersImpl) ContainerCloseStdinHandler(params interaction.ContainerCloseStdinParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	session, err := i.attachServer.Get(context.Background(), params.ID, interactionTimeout)
	if err != nil {
		log.Errorf("%s", err.Error())

		e := &models.Error{
			Message: fmt.Sprintf("No stdin connection found (id: %s): %s", params.ID, err.Error()),
		}
		return interaction.NewContainerCloseStdinNotFound().WithPayload(e)
	}

	if err = session.CloseStdin(); err != nil {
		log.Errorf("%s", err.Error())

		return interaction.NewContainerCloseStdinInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}
	return interaction.NewContainerCloseStdinOK()
}

// ContainerGetStdoutHandler returns the stdout
func (i *InteractionHandlersImpl) ContainerGetStdoutHandler(params interaction.ContainerGetStdoutParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	var ctxDeadline time.Time
	var timeout time.Duration

	// Calculate the timeout for the attach if the caller specified a deadline
	if params.Deadline != nil {
		ctxDeadline = time.Time(*params.Deadline)
		timeout = ctxDeadline.Sub(time.Now())
		log.Debugf("Attempting to get ssh session for container %s stdout with deadline %s", params.ID, ctxDeadline.Format(time.UnixDate))
		if timeout < 0 {
			e := &models.Error{Message: fmt.Sprintf("Deadline for stdout already passed for container %s", params.ID)}
			return interaction.NewContainerGetStdoutInternalServerError().WithPayload(e)
		}
	} else {
		log.Debugf("Attempting to get ssh session for container %s stdout", params.ID)
		timeout = interactionTimeout
	}

	session, err := i.attachServer.Get(context.Background(), params.ID, timeout)
	if err != nil {
		log.Errorf("%s", err.Error())

		// FIXME (caglar10ur): Do not return an error here - https://github.com/vmware/vic/issues/2594
		/*
			e := &models.Error{
				Message: fmt.Sprintf("No stdout connection found (id: %s): %s", params.ID, err.Error()),
			}
			return interaction.NewContainerGetStdoutNotFound().WithPayload(e)
		*/
		return interaction.NewContainerGetStdoutNotFound()
	}

	return NewStreamOutputHandler("stdout").WithPayload(
		NewFlushingReader(
			session.Stdout(),
		),
		params.ID,
		nil,
	)
}

// ContainerGetStderrHandler returns the stderr
func (i *InteractionHandlersImpl) ContainerGetStderrHandler(params interaction.ContainerGetStderrParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	var ctxDeadline time.Time
	var timeout time.Duration

	// Calculate the timeout for the attach if the caller specified a deadline
	if params.Deadline != nil {
		ctxDeadline = time.Time(*params.Deadline)
		timeout = ctxDeadline.Sub(time.Now())
		log.Debugf("Attempting to get ssh session for container %s stderr with deadline %s", params.ID, ctxDeadline.Format(time.UnixDate))
		if timeout < 0 {
			e := &models.Error{Message: fmt.Sprintf("Deadline for stderr already passed for container %s", params.ID)}
			return interaction.NewContainerGetStderrInternalServerError().WithPayload(e)
		}
	} else {
		log.Debugf("Attempting to get ssh session for container %s stderr", params.ID)
		timeout = interactionTimeout
	}

	session, err := i.attachServer.Get(context.Background(), params.ID, timeout)
	if err != nil {
		log.Errorf("%s", err.Error())

		// FIXME (caglar10ur): Do not return an error here - https://github.com/vmware/vic/issues/2594
		/*
			e := &models.Error{
				Message: fmt.Sprintf("No stderr connection found (id: %s): %s", params.ID, err.Error()),
			}
			return interaction.NewContainerGetStderrNotFound().WithPayload(e)
		*/
		return interaction.NewContainerGetStderrNotFound()
	}

	return NewStreamOutputHandler("stderr").WithPayload(
		NewFlushingReader(
			session.Stderr(),
		),
		params.ID,
		nil,
	)
}
