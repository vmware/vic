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
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/groupcache/lru"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/session"
)

const (
	serialOverLANPort  = 8080
	ManagementHostName = "management.localhost"
)

// ContainerCreateConfig defines the parameters for Create call
type ContainerCreateConfig struct {
	Metadata executor.ExecutorConfig

	ParentImageID  string
	ImageStoreName string
	VCHName        string
}

var handles *lru.Cache
var handlesLock sync.Mutex

const handleLen = 16
const lruSize = 1000

func init() {
	handles = lru.New(lruSize)
}

type Handle struct {
	Spec       *spec.VirtualMachineConfigSpec
	ExecConfig executor.ExecutorConfig
	Container  *Container
	State      *State

	key       string
	committed bool
}

func newHandleKey() string {
	b := make([]byte, handleLen)
	rand.Read(b)
	return hex.EncodeToString(b)

}

func newHandle(con *Container) *Handle {
	h := &Handle{
		key:        newHandleKey(),
		committed:  false,
		Container:  con,
		ExecConfig: *con.ExecConfig,
	}

	handlesLock.Lock()
	defer handlesLock.Unlock()

	handles.Add(h.key, h)

	return h
}

func GetHandle(key string) *Handle {
	handlesLock.Lock()
	defer handlesLock.Unlock()

	if h, ok := handles.Get(key); ok {
		return h.(*Handle)
	}

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

func (h *Handle) Commit(ctx context.Context, sess *session.Session) error {
	if h.committed {
		return nil // already committed
	}

	// make sure there is a spec
	h.SetSpec(nil)
	cfg := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(cfg), h.ExecConfig)
	s := h.Spec.Spec()
	s.ExtraConfig = append(s.ExtraConfig, vmomi.OptionValueFromMap(cfg)...)

	if err := h.Container.Commit(ctx, sess, h); err != nil {
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
	// add create time to config
	h.ExecConfig.Common.Created = time.Now().UTC().Unix()
	// configure with debug
	h.ExecConfig.Diagnostics.DebugLevel = VCHConfig.DebugLevel
	// Convert the management hostname to IP
	ips, err := net.LookupIP(ManagementHostName)
	if err != nil {
		log.Errorf("Unable to look up %s during create of %s: %s", ManagementHostName, config.Metadata.ID, err)
		return err
	}

	if len(ips) == 0 {
		log.Errorf("No IP found for %s during create of %s", ManagementHostName, config.Metadata.ID)
		return fmt.Errorf("No IP found on %s", ManagementHostName)
	}

	if len(ips) > 1 {
		log.Errorf("Multiple IPs found for %s during create of %s: %v", ManagementHostName, config.Metadata.ID, ips)
		return fmt.Errorf("Multiple IPs found on %s: %#v", ManagementHostName, ips)
	}

	URI := fmt.Sprintf("tcp://%s:%d", ips[0], serialOverLANPort)

	//FIXME: remove debug network
	backing, err := VCHConfig.DebugNetwork.EthernetCardBackingInfo(ctx)
	if err != nil {
		detail := fmt.Sprintf("unable to generate backing info for debug network - this code can be removed once network mapping/dhcp client are available: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}
	specconfig := &spec.VirtualMachineConfigSpecConfig{
		// FIXME: hardcoded values
		NumCPUs:  2,
		MemoryMB: 2048,

		ConnectorURI: URI,

		ID:   config.Metadata.ID,
		Name: config.Metadata.Name,

		ParentImageID: config.ParentImageID,
		BootMediaPath: VCHConfig.BootstrapImagePath,
		VMPathName:    fmt.Sprintf("[%s]", sess.Datastore.Name()),
		DebugNetwork:  backing,

		ImageStoreName: config.ImageStoreName,

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
