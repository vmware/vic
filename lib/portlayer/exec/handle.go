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
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/groupcache/lru"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/session"
)

// ContainerCreateConfig defines the parameters for Create call
type ContainerCreateConfig struct {
	Metadata executor.ExecutorConfig

	ParentImageID  string
	ImageStoreName string
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
	Spec *spec.VirtualMachineConfigSpec

	// desired
	ExecConfig executor.ExecutorConfig
	State      *State

	Container *Container

	key       string
	committed bool
}

func newHandleKey() string {
	b := make([]byte, handleLen)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err) // This shouldn't happen
	}
	return hex.EncodeToString(b)
}

func newHandle(con *Container) *Handle {
	h := &Handle{
		key:        newHandleKey(),
		committed:  false,
		Container:  con,
		ExecConfig: *con.ExecConfig,
		State:      new(State),
	}

	*h.State = con.State

	handlesLock.Lock()
	defer handlesLock.Unlock()

	handles.Add(h.key, h)

	return h
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

func (h *Handle) IsCommitted() bool {
	return h.committed
}

func (h *Handle) SetSpec(s *spec.VirtualMachineConfigSpec) error {
	if h.Spec != nil {
		if s != nil {
			return fmt.Errorf("spec is already set")
		}

		return nil
	}

	if s == nil {
		// initialization
		s = &spec.VirtualMachineConfigSpec{
			VirtualMachineConfigSpec: &types.VirtualMachineConfigSpec{},
		}
	}

	h.Spec = s
	return nil
}

func (h *Handle) String() string {
	return h.key
}

func (h *Handle) Commit(ctx context.Context, sess *session.Session, waitTime *int32) error {
	if h.committed {
		return nil // already committed
	}

	// make sure there is a spec
	h.SetSpec(nil)
	cfg := make(map[string]string)

	// Set timestamps based on target state
	switch *h.State {
	case StateRunning:
		se := h.ExecConfig.Sessions[h.ExecConfig.ID]
		se.StartTime = time.Now().UTC().Unix()
		h.ExecConfig.Sessions[h.ExecConfig.ID] = se
	case StateStopped:
		se := h.ExecConfig.Sessions[h.ExecConfig.ID]
		se.StopTime = time.Now().UTC().Unix()
		h.ExecConfig.Sessions[h.ExecConfig.ID] = se
	}

	extraconfig.Encode(extraconfig.MapSink(cfg), h.ExecConfig)
	s := h.Spec.Spec()
	s.ExtraConfig = append(s.ExtraConfig, vmomi.OptionValueFromMap(cfg)...)

	if err := h.Container.Commit(ctx, sess, h, waitTime); err != nil {
		return err
	}

	h.committed = true
	removeHandle(h.key)
	return nil
}

func (h *Handle) Close() {
	removeHandle(h.key)
}

func (h *Handle) SetState(s State) {
	h.State = new(State)
	*h.State = s
}

func (h *Handle) Create(ctx context.Context, sess *session.Session, config *ContainerCreateConfig) error {
	defer trace.End(trace.Begin("Handle.Create"))

	if h.Spec != nil {
		log.Debugf("spec has already been set on handle %p during create of %s", h, config.Metadata.ID)
		return fmt.Errorf("spec already set")
	}

	// update the handle with Metadata
	h.ExecConfig = config.Metadata
	// configure with debug
	h.ExecConfig.Diagnostics.DebugLevel = Config.DebugLevel
	// Convert the management hostname to IP
	ips, err := net.LookupIP(constants.ManagementHostName)
	if err != nil {
		log.Errorf("Unable to look up %s during create of %s: %s", constants.ManagementHostName, config.Metadata.ID, err)
		return err
	}

	if len(ips) == 0 {
		log.Errorf("No IP found for %s during create of %s", constants.ManagementHostName, config.Metadata.ID)
		return fmt.Errorf("No IP found on %s", constants.ManagementHostName)
	}

	if len(ips) > 1 {
		log.Errorf("Multiple IPs found for %s during create of %s: %v", constants.ManagementHostName, config.Metadata.ID, ips)
		return fmt.Errorf("Multiple IPs found on %s: %#v", constants.ManagementHostName, ips)
	}

	uuid, err := instanceUUID(config.Metadata.ID)
	if err != nil {
		detail := fmt.Sprintf("unable to get instance UUID: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}
	specconfig := &spec.VirtualMachineConfigSpecConfig{
		// FIXME: hardcoded values
		NumCPUs:  2,
		MemoryMB: 2048,

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
		return err
	}

	h.SetSpec(linux.Spec())
	return nil
}

func (h *Handle) Update(ctx context.Context, sess *session.Session) (*executor.ExecutorConfig, error) {
	defer trace.End(trace.Begin("Handle.Create"))

	ec, err := h.Container.Update(ctx, sess)
	if err != nil {
		return nil, err
	}

	h.ExecConfig = *ec
	return ec, nil
}
