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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/groupcache/lru"
	"golang.org/x/net/context"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
)

// Resources describes the resource allocation for the containerVM
type Resources struct {
	NumCPUs  int64
	MemoryMB int64
}

// ContainerCreateConfig defines the parameters for Create call
type ContainerCreateConfig struct {
	Metadata *executor.ExecutorConfig

	ParentImageID  string
	ImageStoreName string
	Resources      Resources
}

var handles *lru.Cache
var handlesLock sync.Mutex

const (
	handleLen = 16
	lruSize   = 1000
)

func init() {
	handles = lru.New(lruSize)
}

type Handle struct {
	// copy from container cache
	containerBase

	// desired spec
	Spec *spec.VirtualMachineConfigSpec
	// desired state
	targetState State

	// allow for passing outside of the process
	key string
}

func newHandleKey() string {
	b := make([]byte, handleLen)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err) // This shouldn't happen
	}
	return hex.EncodeToString(b)
}

// Added solely to support testing - need a better way to do this
func TestHandle(id string) *Handle {
	defer trace.End(trace.Begin("Handle.Create"))

	h := newHandle(&Container{})
	h.ExecConfig.ID = id

	return h
}

// newHandle creates a handle for an existing container
// con must not be nil
func newHandle(con *Container) *Handle {
	h := &Handle{
		key:           newHandleKey(),
		targetState:   StateUnknown,
		containerBase: *newBase(con.vm, con.Config, con.Runtime),
		// currently every operation has a spec, because even the power operations
		// make changes to extraconfig for timestamps and session status
		Spec: &spec.VirtualMachineConfigSpec{
			VirtualMachineConfigSpec: &types.VirtualMachineConfigSpec{},
		},
	}

	handlesLock.Lock()
	defer handlesLock.Unlock()

	handles.Add(h.key, h)

	return h
}

func (h *Handle) TargetState() State {
	return h.targetState
}

func (h *Handle) SetTargetState(s State) {
	h.targetState = s
}

// GetHandle finds and returns the handle that is referred by key
func GetHandle(key string) *Handle {
	handlesLock.Lock()
	defer handlesLock.Unlock()

	if h, ok := handles.Get(key); ok {
		return h.(*Handle)
	}

	return nil
}

// HandleFromInterface returns the Handle
func HandleFromInterface(key interface{}) *Handle {
	defer trace.End(trace.Begin(""))

	if h, ok := key.(string); ok {
		return GetHandle(h)
	}

	log.Errorf("Type assertion failed for %#+v", key)
	return nil
}

// ReferenceFromHandle returns the reference of the given handle
func ReferenceFromHandle(handle interface{}) interface{} {
	defer trace.End(trace.Begin(""))

	if h, ok := handle.(*Handle); ok {
		return h.String()
	}

	log.Errorf("Type assertion failed for %#+v", handle)
	return nil
}

func removeHandle(key string) {
	handlesLock.Lock()
	defer handlesLock.Unlock()

	handles.Remove(key)
}

func (h *Handle) String() string {
	return h.key
}

func (h *Handle) Commit(ctx context.Context, sess *session.Session, waitTime *int32) error {
	if err := Commit(ctx, sess, h, waitTime); err != nil {
		return err
	}

	removeHandle(h.key)
	return nil
}

func (h *Handle) Close() {
	removeHandle(h.key)
}

// Create returns a new handle that can be Committed to create a new container.
// At this time the config is *not* deep copied so should not be changed once passed
//
// TODO: either deep copy the configuration, or provide an alternative means of passing the data that
// avoids the need for the caller to unpack/repack the parameters
func Create(ctx context.Context, sess *session.Session, config *ContainerCreateConfig) (*Handle, error) {
	defer trace.End(trace.Begin("Handle.Create"))

	h := &Handle{
		key:         newHandleKey(),
		targetState: StateCreated,
		containerBase: containerBase{
			ExecConfig: config.Metadata,
		},
	}

	// configure with debug
	h.ExecConfig.Diagnostics.DebugLevel = Config.DebugLevel

	// Convert the management hostname to IP
	ips, err := net.LookupIP(constants.ManagementHostName)
	if err != nil {
		log.Errorf("Unable to look up %s during create of %s: %s", constants.ManagementHostName, config.Metadata.ID, err)
		return nil, err
	}

	if len(ips) == 0 {
		log.Errorf("No IP found for %s during create of %s", constants.ManagementHostName, config.Metadata.ID)
		return nil, fmt.Errorf("No IP found on %s", constants.ManagementHostName)
	}

	if len(ips) > 1 {
		log.Errorf("Multiple IPs found for %s during create of %s: %v", constants.ManagementHostName, config.Metadata.ID, ips)
		return nil, fmt.Errorf("Multiple IPs found on %s: %#v", constants.ManagementHostName, ips)
	}

	uuid, err := instanceUUID(config.Metadata.ID)
	if err != nil {
		detail := fmt.Sprintf("unable to get instance UUID: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	specconfig := &spec.VirtualMachineConfigSpecConfig{
		NumCPUs:  int32(config.Resources.NumCPUs),
		MemoryMB: config.Resources.MemoryMB,

		ID:       config.Metadata.ID,
		Name:     config.Metadata.Name,
		BiosUUID: uuid,

		ParentImageID: config.ParentImageID,
		BootMediaPath: Config.BootstrapImagePath,
		VMPathName:    fmt.Sprintf("[%s]", sess.Datastore.Name()),

		ImageStoreName: config.ImageStoreName,
		ImageStorePath: &Config.ImageStores[0],

		Metadata: config.Metadata,
	}
	log.Debugf("Config: %#v", specconfig)

	// Create a linux guest
	linux, err := guest.NewLinuxGuest(ctx, sess, specconfig)
	if err != nil {
		log.Errorf("Failed during linux specific spec generation during create of %s: %s", config.Metadata.ID, err)
		return nil, err
	}

	h.Spec = linux.Spec()

	handlesLock.Lock()
	defer handlesLock.Unlock()

	handles.Add(h.key, h)

	return h, nil
}
