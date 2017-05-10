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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"time"

	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/groupcache/lru"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
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

	// should this change trigger a reload in the target container
	reload bool

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

func (h *Handle) Reload() {
	h.reload = true
}

// Rename updates the container name in ExecConfig as well as the vSphere display name
func (h *Handle) Rename(newName string) *Handle {
	defer trace.End(trace.Begin(newName))

	h.ExecConfig.Name = newName

	s := &spec.VirtualMachineConfigSpecConfig{
		ID:   h.ExecConfig.ID,
		Name: newName,
	}
	h.Spec.Spec().Name = util.DisplayName(s)

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

func (h *Handle) String() string {
	return h.key
}

func (h *Handle) Commit(ctx context.Context, sess *session.Session, waitTime *int32) error {
	cfg := make(map[string]string)

	// Set timestamps based on target state
	switch h.TargetState() {
	case StateRunning:
		for _, sc := range h.ExecConfig.Sessions {
			sc.StartTime = time.Now().UTC().Unix()
			sc.Started = ""
			sc.ExitStatus = 0
		}
	case StateStopped:
		for _, sc := range h.ExecConfig.Sessions {
			sc.StopTime = time.Now().UTC().Unix()
			sc.Started = ""
		}
	}

	// this is the first step to use portlayer manage VCH, to avoid upgrade VCH configuration, using different prefix to encoding atm.
	// TODO: after VCH bootstrap is updated to tether, instead of vic-init, this can be removed

	s := h.Spec.Spec()
	// if runtime is nil, should be fresh container create
	if h.Runtime == nil || h.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff || h.TargetState() == StateStopped {
		extraconfig.EncodeWithPrefix(extraconfig.MapSink(cfg), h.ExecConfig, h.ConfigPrefix)
		s.ExtraConfig = append(s.ExtraConfig, vmomi.OptionValueFromMap(cfg)...)
	} else {
		extraconfig.EncodeWithPrefix(extraconfig.ScopeFilterSink(extraconfig.NonPersistent|extraconfig.Hidden, extraconfig.MapSink(cfg)), h.ExecConfig, h.ConfigPrefix)
		s.ExtraConfig = append(s.ExtraConfig, vmomi.OptionValueFromMap(cfg)...)
	}

	if err := Commit(ctx, sess, h, waitTime); err != nil {
		return err
	}

	h.Close()
	return nil
}

// CommitWithoutSpec sets the handle's spec to nil so that Commit operation only does a state change and won't touch the extraconfig
func (h *Handle) CommitWithoutSpec(ctx context.Context, sess *session.Session, waitTime *int32) error {
	h.Spec = nil

	if err := Commit(ctx, sess, h, waitTime); err != nil {
		return err
	}

	h.Close()
	return nil
}
func (h *Handle) Close() {
	handlesLock.Lock()
	defer handlesLock.Unlock()

	handles.Remove(h.key)
}

// Create returns a new handle that can be Committed to create a new container.
// At this time the config is *not* deep copied so should not be changed once passed
//
// TODO: either deep copy the configuration, or provide an alternative means of passing the data that
// avoids the need for the caller to unpack/repack the parameters
func Create(ctx context.Context, vmomiSession *session.Session, conf *ContainerCreateConfig) (*Handle, error) {
	defer trace.End(trace.Begin("Handle.Create"))

	prefix := extraconfig.DefaultPrefix
	if Config.ManagingVCH {
		prefix = config.VCHPrefix
	}
	h := &Handle{
		key:         newHandleKey(),
		targetState: StateCreated,
		containerBase: containerBase{
			ExecConfig:   conf.Metadata,
			ConfigPrefix: prefix,
		},
	}

	// configure with debug
	h.ExecConfig.Diagnostics.DebugLevel = Config.DebugLevel

	specconfig := &spec.VirtualMachineConfigSpecConfig{
		NumCPUs:  int32(conf.Resources.NumCPUs),
		MemoryMB: conf.Resources.MemoryMB,

		ID:   conf.Metadata.ID,
		Name: conf.Metadata.Name,

		ParentImageID: conf.ParentImageID,
		BootMediaPath: Config.BootstrapImagePath,
		VMPathName:    fmt.Sprintf("[%s]", vmomiSession.Datastore.Name()),

		ImageStoreName: conf.ImageStoreName,

		Metadata:     conf.Metadata,
		ConfigPrefix: h.ConfigPrefix,
	}
	if Config.ImageStores != nil {
		specconfig.ImageStorePath = &Config.ImageStores[0]
	}

	// if not vsan, set the datastore folder name to containerID
	if !vmomiSession.IsVSAN(ctx) {
		specconfig.VMPathName = fmt.Sprintf("[%s] %s/%s.vmx", vmomiSession.Datastore.Name(), specconfig.ID, specconfig.ID)
		if Config.ManagingVCH {
			specconfig.VMPathName = fmt.Sprintf("[%s] %s/%s.vmx", vmomiSession.Datastore.Name(), conf.Metadata.Name, conf.Metadata.Name)
		}
	}

	specconfig.VMFullName = conf.Metadata.Name
	specconfig.AlternateGuestName = constants.DefaultAltVCHGuestName()
	if !Config.ManagingVCH {
		specconfig.VMFullName = util.DisplayName(specconfig)
		specconfig.AlternateGuestName = constants.DefaultAltContainerGuestName()
	}

	// log only core portions
	s := specconfig
	log.Debugf("id: %s, name: %s, cpu: %d, mem: %d, parent: %s, os: %s, path: %s", s.ID, s.Name, s.NumCPUs, s.MemoryMB, s.ParentImageID, s.BootMediaPath, s.VMPathName)
	m := s.Metadata
	log.Debugf("annotations: %#v, reponame: %s", m.Annotations, m.RepoName)
	for name, sess := range m.Sessions {
		log.Debugf("session: %s, path: %s, dir: %s, runblock: %t, tty: %t, restart: %t, stdin: %t, stopsig: %s",
			name, sess.Cmd.Path, sess.Cmd.Dir, sess.RunBlock, sess.Tty, sess.Restart, sess.OpenStdin, sess.StopSignal)
	}

	// If the debug level is high, dump everything
	// we still do the logging above for consistency so searching the logs for common strings works.
	// TODO: move this into a debug level aware structure renderer
	if Config.DebugLevel > 2 {
		log.Debugf("Config: %#v", specconfig)
		log.Debugf("Executor spec: %#v", *specconfig.Metadata)
		for _, sess := range m.Sessions {
			log.Debugf("Session spec: %#v", *sess)
		}
	}

	// Create a linux guest
	linux, err := guest.NewLinuxGuest(ctx, vmomiSession, specconfig)
	if err != nil {
		log.Errorf("Failed during linux specific spec generation during create of %s: %s", conf.Metadata.ID, err)
		return nil, err
	}

	h.Spec = linux.Spec()

	handlesLock.Lock()
	defer handlesLock.Unlock()

	handles.Add(h.key, h)

	return h, nil
}
