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
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime"
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

	i.attachServer = attach.NewAttachServer(constants.ManagementHostName, 0)

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

	return NewContainerOutputHandler("stdout").WithPayload(
		NewFlushingReader(
			session.Stdout(),
		),
		params.ID,
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

	return NewContainerOutputHandler("stderr").WithPayload(
		NewFlushingReader(
			session.Stderr(),
		),
		params.ID,
	)
}

// GenericFlusher is a custom reader to allow us to detach cleanly during an io.Copy
type GenericFlusher interface {
	Flush()
}

type FlushingReader struct {
	io.Reader
	io.WriterTo

	flusher   GenericFlusher
	initBytes []byte
}

func NewFlushingReader(rdr io.Reader) *FlushingReader {
	return &FlushingReader{Reader: rdr, flusher: nil, initBytes: nil}
}

func NewFlushingReaderWithInitBytes(rdr io.Reader, initBytes []byte) *FlushingReader {
	return &FlushingReader{Reader: rdr, flusher: nil, initBytes: initBytes}
}

func (d *FlushingReader) AddFlusher(flusher GenericFlusher) {
	d.flusher = flusher
}

// readDetectInit() is used by WriteTo() which is used by io.Copy.  It attempts
// to detect a init byte buffer.  If it finds that init byte sequence, it is
// ignored.  This reader does not care about the init sequeunce.  The init sequence
// maybe used by the higher level interaction, which in this case is the Swagger
// establishing initial connection for stdin.
//
// Panics if the buf is smaller than the initBytes
func (d *FlushingReader) readDetectInit(buf []byte) (int, error) {
	initLen := len(d.initBytes)

	// fast path - len(nil) return 0
	if initLen == 0 {
		return d.Read(buf)
	}

	// make sure we have enough room
	if len(buf) < initLen {
		panic("Read buffer is smaller than the initialization byte sequence")
	}

	total := 0
	upto := 0
	for total < initLen {
		nr, err := d.Read(buf[total:])
		if nr > 0 {
			total += nr
			// we are only interested with the first initLen bytes
			upto = total
			if upto > initLen {
				upto = initLen
			}
			if bytes.Compare(d.initBytes[0:upto], buf[0:upto]) != 0 {
				// First bytes aren't part of init bytes so client must not be
				// the docker personality so break and ignore looking for the
				// init bytes.
				log.Debugf("Did not find primer bytes, stopping watch")
				return total, err
			}
		}
		if err != nil && total < initLen {
			log.Debugf("Primer bytes read %d bytes, err %s, stopping watch", nr, err)
			return 0, err
		}
	}

	// would have returned in the compare clause if not matching init bytes
	copy(buf[0:], buf[initLen:])
	log.Debugf("Found primer bytes, port layer client might be personality server")

	// no risk of returning <0
	return total - initLen, nil
}

// Derived from go's io.Copy.  We use a smaller buffer so as to not hold up
// writing out data.  Go's version allocates 32k, and the Read will wait till
// buffer is filled (unless EOF is encountered).  Also, we force a flush if
// a flusher is added.  We've seen cases where the last bit of data for a
// screen doesn't reach the docker engine api server.  The flush solves that
// issue.
func (d *FlushingReader) WriteTo(w io.Writer) (written int64, err error) {
	buf := make([]byte, ioCopyBufferSize)

	nr, er := d.readDetectInit(buf)
	for {
		log.Debugf("[%p] nr: %d", d, nr)
		if nr > 0 {
			log.Debugf("[%p] buf: %s", d, string(buf[:nr]))
			nw, ew := w.Write(buf[0:nr])
			if d.flusher != nil {
				d.flusher.Flush()
			}
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				log.Debugf("[%p] ew: %s", d, ew)
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			log.Debugf("[%p] er: %s", d, er)
			err = er
			break
		}
		nr, er = d.Read(buf)
	}
	log.Debugf("[%p] written: %d err: %s", d, written, err)
	return written, err
}

// ContainerOutputHandler is custom return handlers for stdout/stderr
type ContainerOutputHandler struct {
	outputStream *FlushingReader
	containerID  string
	outputName   string
}

// NewContainerOutputHandler creates ContainerOutputHandler with default headers values
func NewContainerOutputHandler(name string) *ContainerOutputHandler {
	return &ContainerOutputHandler{outputName: name}
}

// WithPayload adds the payload to the container set stdin internal server error response
func (c *ContainerOutputHandler) WithPayload(payload *FlushingReader, id string) *ContainerOutputHandler {
	c.outputStream = payload
	c.containerID = id
	return c
}

// WriteResponse to the client
func (c *ContainerOutputHandler) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusOK)
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
		c.outputStream.AddFlusher(f)
	}

	_, err := io.Copy(rw, c.outputStream)
	if err != nil {
		log.Debugf("Error copying %s stream for container %s: %s", c.outputName, c.containerID, err)
	} else {
		log.Debugf("Finished copying %s stream for container %s", c.outputName, c.containerID)
	}
}
