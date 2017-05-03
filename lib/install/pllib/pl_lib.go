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

package pllib

import (
	"context"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/logging"
	"github.com/vmware/vic/lib/portlayer/task"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var (
	stopTimeout = int32(60) //60s
)

const (
	Running = "RUNNING"
	Stopped = "STOPPED"
	Created = "CREATED"
)

// PLClient is the client to call portlayer handlers as library, for in vic-machine portlayer is not started as a service
type Client struct {
	s *session.Session
}

func NewClient(s *session.Session) *Client {
	// this portlayer is used to manage VCH
	exec.Config.ManagingVCH = true

	// instantiate the container cache
	exec.NewContainerCache()

	pl := &Client{
		s: s,
	}
	return pl
}

func (pl *Client) SetParentResources(vApp *object.VirtualApp, pool *object.ResourcePool) {
	exec.Config.VirtualApp = vApp
	if vApp != nil {
		exec.Config.ResourcePool = vApp.ResourcePool
	} else {
		exec.Config.ResourcePool = pool
	}
}

// CreateVchHandle returns portlayer create container handle, this handle should be used to reconfigure and commit
func (pl *Client) CreateVchHandle(ctx context.Context, conf *executor.ExecutorConfig, cpu, memory int64) (interface{}, error) {
	defer trace.End(trace.Begin(""))

	// Create the executor.ExecutorCreateConfig
	c := &exec.ContainerCreateConfig{
		Metadata:       conf,
		ParentImageID:  "",
		ImageStoreName: "",
		Resources: exec.Resources{
			NumCPUs:  cpu,
			MemoryMB: memory,
		},
	}

	return exec.Create(ctx, pl.s, c)
}

// SetVCHMoref replace VCH configuration ID from creating id to VM mob ref
// This method will commit changes and update container cache, to avoid id inconsistence issue in portlayer
func (pl *Client) SetVCHMoref(ctx context.Context, id string) (string, error) {
	defer trace.End(trace.Begin(id))
	handle := pl.NewHandle(ctx, id)
	if handle == nil {
		return "", errors.Errorf("%s is not found", id)
	}

	c := exec.Containers.Container(id)
	mobID := handle.VMReference().String()
	handle.ExecConfig.ID = mobID

	err := handle.Commit(ctx, pl.s, &stopTimeout)
	if err != nil {
		return "", err
	}
	exec.Containers.Remove(id)
	c.ExecConfig.ID = mobID
	exec.Containers.Put(c)
	return mobID, nil
}

// NewHandle creates new handle
func (pl *Client) NewHandle(ctx context.Context, id string) *exec.Handle {
	c := exec.Containers.Container(id)
	if c == nil {
		return nil
	}
	return c.NewHandle(ctx)
}

func (pl *Client) AddTask(ctx context.Context, h interface{}, t *executor.SessionConfig) (interface{}, error) {
	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}
	defer trace.End(trace.Begin(fmt.Sprintf("%s: %s", handle.String(), t.Cmd.Path)))
	op := trace.NewOperation(ctx, "task.Join(%s, %s)", handle.String(), t.ID)

	handleprime, err := task.Join(&op, handle, t)
	if err != nil {
		return handleprime, err
	}
	op = trace.NewOperation(ctx, "task.Bind(%s, %s)", handle.String(), t.ID)
	return task.Bind(&op, handleprime, t.ID)
}

func (pl *Client) AddNetworks(ctx context.Context, h interface{}, networks map[string]*executor.NetworkEndpoint) (interface{}, error) {
	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}
	defer trace.End(trace.Begin(handle.String()))

	var devices object.VirtualDeviceList
	// network name:alias, to avoid create multiple devices for same network
	slots := make(map[int32]bool)
	nets := make(map[string]*executor.NetworkEndpoint)

	for name, endpoint := range networks {
		if pnic, ok := nets[endpoint.Network.Common.ID]; ok {
			// there's already a NIC on this network
			endpoint.Common.ID = pnic.Common.ID
			log.Infof("Network role %q is sharing NIC with %q", name, pnic.Network.Common.Name)
			continue
		}

		moref := new(types.ManagedObjectReference)
		if ok := moref.FromString(endpoint.Network.ID); !ok {
			return handle, errors.Errorf("serialized managed object reference in unexpected format: %q", endpoint.Network.ID)
		}
		obj, err := pl.s.Finder.ObjectReference(ctx, *moref)
		if err != nil {
			return handle, errors.Errorf("unable to reacquire reference for network %q from serialized form: %q", endpoint.Network.Name, endpoint.Network.ID)
		}
		network, ok := obj.(object.NetworkReference)
		if !ok {
			return handle, errors.Errorf("reacquired reference for network %q, from serialized form %q, was not a network: %T", endpoint.Network.Name, endpoint.Network.ID, obj)
		}

		backing, err := network.EthernetCardBackingInfo(ctx)
		if err != nil {
			err = errors.Errorf("Failed to get network backing info for %q: %s", network, err)
			return handle, err
		}

		nic, err := devices.CreateEthernetCard("vmxnet3", backing)
		if err != nil {
			err = errors.Errorf("Failed to create Ethernet Card spec for %s", err)
			return handle, err
		}

		slot := handle.Spec.AssignSlotNumber(nic, slots)
		if slot == spec.NilSlot {
			err = errors.Errorf("Failed to assign stable PCI slot for %q network card", name)
		}

		endpoint.Common.ID = strconv.Itoa(int(slot))
		slots[slot] = true
		log.Debugf("Setting %q to slot %d", name, slot)

		devices = append(devices, nic)

		nets[endpoint.Network.Common.ID] = endpoint
	}

	deviceChange, err := devices.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		return handle, err
	}

	handle.Spec.DeviceChange = append(handle.Spec.DeviceChange, deviceChange...)
	return handle, err
}

func (pl *Client) AddLogging(h interface{}) (interface{}, error) {
	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}
	defer trace.End(trace.Begin(handle.String()))
	return logging.Join(handle)
}

func (pl *Client) Commit(ctx context.Context, h interface{}) error {
	handle, ok := h.(*exec.Handle)
	if !ok {
		return fmt.Errorf("Type assertion failed for %#+v", handle)
	}
	return handle.Commit(ctx, pl.s, &stopTimeout)
}

func (pl *Client) VCHFolderName(ctx context.Context, id string) (string, error) {
	c := exec.Containers.Container(id)
	if c == nil {
		return "", errors.Errorf("%s is not found", id)
	}
	return c.VMFolder(ctx)
}

func (pl *Client) UpdateApplianceISOFiles(h interface{}, newFile string) (interface{}, error) {
	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}
	defer trace.End(trace.Begin(fmt.Sprintf("%s: %s", handle.String(), newFile)))

	// get the virtual device list
	devices := object.VirtualDeviceList(handle.Config.Hardware.Device)

	// find the single cdrom
	cd, err := devices.FindCdrom("")
	if err != nil {
		log.Errorf("Failed to get CD rom device from appliance: %s", err)
		//		return handle, err
		ide, err := devices.FindIDEController("")
		if err != nil {
			log.Errorf("Failed to find IDE controller for appliance: %s", err)
			return nil, err
		}
		cdrom, err := devices.CreateCdrom(ide)
		if err != nil {
			log.Errorf("Failed to create Cdrom device for appliance: %s", err)
			return nil, err
		}
		cdrom = devices.InsertIso(cdrom, newFile)

		config := &types.VirtualDeviceConfigSpec{
			Device:    cdrom,
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
		}
		handle.Spec.DeviceChange = append(handle.Spec.DeviceChange, config)
		return handle, nil
	}

	oldApplianceISO := cd.Backing.(*types.VirtualCdromIsoBackingInfo).FileName
	if oldApplianceISO == newFile {
		log.Debugf("Target file name %q is same to old one, no need to change.")
		return handle, nil
	}
	cd = devices.InsertIso(cd, newFile)

	config := &types.VirtualDeviceConfigSpec{
		Device:    cd,
		Operation: types.VirtualDeviceConfigSpecOperationEdit,
	}
	handle.Spec.DeviceChange = append(handle.Spec.DeviceChange, config)
	return handle, nil
}

func (pl *Client) UpdateExtraConfig(h interface{}, data map[string]string) (interface{}, error) {
	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}
	defer trace.End(trace.Begin(handle.String()))

	s := handle.Spec.Spec()
	s.ExtraConfig = append(s.ExtraConfig, vmomi.OptionValueFromMap(data)...)
	return handle, nil
}

func (pl *Client) ChangeState(ctx context.Context, h interface{}, state string) (interface{}, error) {
	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}

	var plState exec.State
	switch state {
	case "RUNNING":
		plState = exec.StateRunning
	case "STOPPED":
		plState = exec.StateStopped
	case "CREATED":
		plState = exec.StateCreated
	default:
		return handle, errors.New("unknown state")
	}

	handle.SetTargetState(plState)
	handle.TargetState()
	return handle, nil
}
