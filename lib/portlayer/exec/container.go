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
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
	"golang.org/x/net/context"
)

var containers map[ID]*Container
var containersLock sync.Mutex

func init() {
	containers = make(map[ID]*Container)
}

type State int

const (
	StateRunning = iota + 1
	StateStopped

	propertyCollectorTimeout = 3 * time.Minute
)

type Container struct {
	sync.Mutex

	ID ID

	ExecConfig *metadata.ExecutorConfig
	State      State
	// friendly description of state
	Status string

	VMUnsharedDisk int64

	vm *vm.VirtualMachine
}

func NewContainer(id ID) *Handle {
	con := &Container{
		ID:         id,
		ExecConfig: &metadata.ExecutorConfig{},
		State:      StateStopped,
	}

	containersLock.Lock()
	containers[id] = con
	containersLock.Unlock()

	return con.newHandle()
}

func GetContainer(id ID) *Handle {
	containersLock.Lock()
	defer containersLock.Unlock()

	if con, ok := containers[id]; ok {
		return con.newHandle()
	}

	return nil
}

func (c *Container) newHandle() *Handle {
	return newHandle(c)
}

func (c *Container) cacheExecConfig(ec *metadata.ExecutorConfig) {
	c.Lock()
	defer c.Unlock()

	c.ExecConfig = ec
}

func (c *Container) Commit(ctx context.Context, sess *session.Session, h *Handle) error {
	defer trace.End(trace.Begin("Committing handle"))

	c.Lock()
	defer c.Unlock()

	if h.Spec != nil {
		s := h.Spec.Spec()
		if c.vm != nil {
			// reconfigure
			_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
				return c.vm.Reconfigure(ctx, *s)
			})

			if err != nil {
				return err
			}
		} else {
			// Find the Virtual Machine folder that we use
			folders, err := sess.Datacenter.Folders(ctx)
			if err != nil {
				return err
			}
			parent := folders.VmFolder

			// FIXME: Replace this simple logic with DRS placement
			// Pick a random host
			hosts, err := sess.Datastore.AttachedClusterHosts(ctx, sess.Cluster)
			if err != nil {
				return err
			}
			host := hosts[rand.Intn(len(hosts))]
			// Create the vm
			res, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
				return parent.CreateVM(ctx, *s, Config.ResourcePool, host)
			})
			if err != nil {
				return err
			}

			c.vm = vm.NewVirtualMachine(ctx, sess, res.Result.(types.ManagedObjectReference))
		}
	}

	c.ExecConfig = &h.ExecConfig

	if h.State != nil {
		if c.vm == nil {
			return fmt.Errorf("no VM to do state change")
		}

		switch *h.State {
		case StateRunning:
			// start the container
			if err := h.Container.start(ctx); err != nil {
				return err
			}

		case StateStopped:
			// stop the container
			if err := h.Container.stop(ctx); err != nil {
				return err
			}
		}

		c.State = *h.State
	}

	return nil
}

// Start starts a container vm with the given params
func (c *Container) start(ctx context.Context) error {
	defer trace.End(trace.Begin("Container.start"))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}

	// Power on
	_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return c.vm.PowerOn(ctx)
	})
	if err != nil {
		return err
	}

	// guestinfo key that we want to wait for
	key := fmt.Sprintf("guestinfo..sessions|%s.started", c.ID)
	var detail string

	// Wait some before giving up...
	ctx, cancel := context.WithTimeout(ctx, propertyCollectorTimeout)
	defer cancel()

	detail, err = c.vm.WaitForKeyInExtraConfig(ctx, key)
	if err != nil {
		return fmt.Errorf("unable to wait for process launch status: %s", err.Error())
	}

	if detail != "true" {
		return errors.New(detail)
	}

	return nil
}

func (c *Container) stop(ctx context.Context) error {
	defer trace.End(trace.Begin("Container.stop"))

	if c.vm == nil {
		return fmt.Errorf("vm not set")
	}

	//TODO: make the shutdown much cleaner, right now we just pull the plug on the vm.(may need corresponding work in tether.)

	// Power off
	_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return c.vm.PowerOff(ctx)
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Container) Remove(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()

	if c.vm == nil {
		return fmt.Errorf("VM has already been removed")
	}

	//removes the vm from vsphere
	_, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return c.vm.Destroy(ctx)
	})
	if err != nil {
		return err
	}

	//removes container from map
	containersLock.Lock()
	delete(containers, c.ID)
	containersLock.Unlock()

	return nil
}

// return a list of container attributes
func List(ctx context.Context, sess *session.Session, all *bool) ([]Container, error) {

	// for now we'll go to the infrastructure
	// future iteration will utilize cache & event stream
	moVMs, err := infraContainers(ctx, sess)
	if err != nil {
		return nil, err
	}
	// convert to container
	containers := convertInfraContainers(moVMs, all)
	return containers, nil
}

// get the containerVMs from infrastructure
func infraContainers(ctx context.Context, sess *session.Session) ([]mo.VirtualMachine, error) {
	var rp mo.ResourcePool
	var vms []mo.VirtualMachine

	// popluate the vm property of the vch resource pool
	if err := Config.ResourcePool.Properties(ctx, Config.ResourcePool.Reference(), []string{"vm"}, &rp); err != nil {
		log.Errorf("List failed to get %s resource pool child vms: %s", Config.ResourcePool.Name(), err)
		return nil, err
	}

	// We now have morefs for the pools VMs, so
	//define attributes do we need from the VMs
	attrib := []string{"config", "runtime.powerState", "summary"}

	// populate the vm properties
	sess.Retrieve(ctx, rp.Vm, attrib, &vms)

	return vms, nil
}

// convert the infra containers to a container object
func convertInfraContainers(vms []mo.VirtualMachine, all *bool) []Container {
	var containerVMs []Container

	for i := range vms {
		// poweredOn or all states
		if !*all && vms[i].Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff {
			// don't want it
			log.Debugf("Skipping poweredOff VM %s", vms[i].Config.Name)
			continue
		}

		container := &Container{ExecConfig: &metadata.ExecutorConfig{}}
		source := extraconfig.OptionValueSource(vms[i].Config.ExtraConfig)
		extraconfig.Decode(source, container.ExecConfig)

		// check extraConfig to see if we have a containerVM -- assumes
		// that ID will always be populated for each containerVM
		if container.ExecConfig == nil || container.ExecConfig.ID == "" {
			log.Debugf("Skipping non-container vm %s", vms[i].Config.Name)
			continue
		}

		// set state & friendly status
		if vms[i].Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
			container.State = StateRunning
			container.Status = "Running"
		} else {
			container.State = StateStopped
			container.Status = "Stopped"
		}
		container.VMUnsharedDisk = vms[i].Summary.Storage.Unshared

		containerVMs = append(containerVMs, *container)

	}

	return containerVMs
}
