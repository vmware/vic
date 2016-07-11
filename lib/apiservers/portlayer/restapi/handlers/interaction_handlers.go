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
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/interaction"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/options"
	"github.com/vmware/vic/lib/portlayer/attach"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/vsphere/session"
)

// ExecHandlersImpl is the receiver for all of the exec handler methods
type InteractionHandlersImpl struct {
	attachServer *attach.Server
}

var (
	interactionSession = &session.Session{}
)

const interactionTimeout time.Duration = 600 * time.Second

func (i *InteractionHandlersImpl) Configure(api *operations.PortLayerAPI, _ *HandlerContext) {
	var err error

	api.InteractionContainerResizeHandler = interaction.ContainerResizeHandlerFunc(i.ContainerResizeHandler)
	api.InteractionContainerSetStdinHandler = interaction.ContainerSetStdinHandlerFunc(i.ContainerSetStdinHandler)
	api.InteractionContainerGetStdoutHandler = interaction.ContainerGetStdoutHandlerFunc(i.ContainerGetStdoutHandler)
	api.InteractionContainerGetStderrHandler = interaction.ContainerGetStderrHandlerFunc(i.ContainerGetStderrHandler)

	sessionconfig := &session.Config{
		Service:        options.PortLayerOptions.SDK,
		Insecure:       options.PortLayerOptions.Insecure,
		Keepalive:      options.PortLayerOptions.Keepalive,
		DatacenterPath: options.PortLayerOptions.DatacenterPath,
		ClusterPath:    options.PortLayerOptions.ClusterPath,
		PoolPath:       options.PortLayerOptions.PoolPath,
		DatastorePath:  options.PortLayerOptions.DatastorePath,
	}

	ctx := context.Background()
	interactionSession, err = session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		log.Fatalf("InteractionHandler ERROR: %s", err)
	}

	i.attachServer = attach.NewAttachServer(exec.ManagementHostName, 0)

	if err := i.attachServer.Start(); err != nil {
		log.Fatalf("Attach server unable to start: %s", err)
	}
}

func (i *InteractionHandlersImpl) ContainerResizeHandler(params interaction.ContainerResizeParams) middleware.Responder {
	// Get the ssh session to the container
	connContainer, err := i.attachServer.Get(context.Background(), params.ID, interactionTimeout)
	if err != nil {
		retErr := &models.Error{Message: fmt.Sprintf("No such container: %s", params.ID)}
		return interaction.NewContainerResizeNotFound().WithPayload(retErr)
	}

	// Request a resize
	cWidth := uint32(params.Width)
	cHeight := uint32(params.Height)

	err = connContainer.Resize(cWidth, cHeight, 0, 0)
	if err != nil {
		return interaction.NewContainerResizeInternalServerError()
	}

	return interaction.NewContainerResizeOK()
}

func (i *InteractionHandlersImpl) ContainerSetStdinHandler(params interaction.ContainerSetStdinParams) middleware.Responder {
	log.Printf("Attempting to get ssh session for container %s stdin", params.ID)
	sshConn, err := i.attachServer.Get(context.Background(), params.ID, interactionTimeout)
	if err != nil {
		err = fmt.Errorf("No stdin found (id:%s): %s", params.ID, err.Error())
		log.Errorf("%s", err.Error())

		return interaction.NewContainerSetStdinNotFound()
	}

	detachableIn := NewFlushingReader(params.RawStream)
	_, err = io.Copy(sshConn.Stdin(), detachableIn)
	if err != nil {
		err = fmt.Errorf("Error copying stdin (id:%s): %s", params.ID, err.Error())
		log.Errorf("%s", err.Error())
	}

	log.Printf("Done copying stdin")

	return interaction.NewContainerSetStdinOK()
}

func (i *InteractionHandlersImpl) ContainerGetStdoutHandler(params interaction.ContainerGetStdoutParams) middleware.Responder {
	log.Printf("Attempting to get ssh session for container %s stdout", params.ID)
	sshConn, err := i.attachServer.Get(context.Background(), params.ID, interactionTimeout)
	if err != nil {

		err = fmt.Errorf("No stdout found for %s: %s", params.ID, err.Error())
		log.Errorf("%s", err.Error())

		return interaction.NewContainerGetStdoutNotFound()
	}

	detachableOut := NewFlushingReader(sshConn.Stdout())
	return NewContainerOutputHandler("stdout").WithPayload(detachableOut, params.ID)
}

func (i *InteractionHandlersImpl) ContainerGetStderrHandler(params interaction.ContainerGetStderrParams) middleware.Responder {
	log.Printf("Attempting to get ssh session for container %s stderr", params.ID)
	sshConn, err := i.attachServer.Get(context.Background(), params.ID, interactionTimeout)
	if err != nil {

		err = fmt.Errorf("No stderr found for %s: %s", params.ID, err.Error())
		log.Errorf("%s", err.Error())

		return interaction.NewContainerGetStderrNotFound()
	}

	detachableErr := NewFlushingReader(sshConn.Stderr())
	return NewContainerOutputHandler("stderr").WithPayload(detachableErr, params.ID)
}

// Custom reader to allow us to detach cleanly during an io.Copy

type GenericFlusher interface {
	Flush()
}

type FlushingReader struct {
	io.Reader
	io.WriterTo

	flusher GenericFlusher
}

func NewFlushingReader(rdr io.Reader) *FlushingReader {
	return &FlushingReader{Reader: rdr, flusher: nil}
}

func (d *FlushingReader) AddFlusher(flusher GenericFlusher) {
	d.flusher = flusher
}

// Derived from go's io.Copy.  We use a smaller buffer so as to not hold up
// writing out data.  Go's version allocates 32k, and the Read will wait till
// buffer is filled (unless EOF is encountered).  Also, we force a flush if
// a flusher is added.  We've seen cases where the last bit of data for a
// screen doesn't reach the docker engine api server.  The flush solves that
// issue.
func (d *FlushingReader) WriteTo(w io.Writer) (written int64, err error) {
	buf := make([]byte, 64)

	for {
		nr, er := d.Read(buf)
		if nr > 0 {
			nw, ew := w.Write(buf[0:nr])
			if d.flusher != nil {
				d.flusher.Flush()
			}
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
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
			err = er
			break
		}
	}
	return written, err
}

// Custom return handlers for stdout/stderr

type ContainerOutputHandler struct {
	outputStream *FlushingReader
	containerID  string
	outputName   string
}

// NewContainerSetStdinInternalServerError creates ContainerSetStdinInternalServerError with default headers values
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
func (c *ContainerOutputHandler) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/octet-streaming")
	rw.Header().Set("Transfer-Encoding", "chunked")
	if f, ok := rw.(http.Flusher); ok {
		c.outputStream.AddFlusher(f)
	}
	_, err := io.Copy(rw, c.outputStream)

	if err != nil {
		log.Printf("Error copying %s stream for container %s: %s", c.outputName, c.containerID, err)
	} else {
		log.Printf("Finished copying %s stream for container %s", c.outputName, c.containerID)
	}
}
