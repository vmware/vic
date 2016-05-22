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

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
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
	StateRunning = iota
	StateStopped = iota

	propertyCollectorTimeout = 3 * time.Minute
)

type Container struct {
	sync.Mutex

	ID ID

	vm *vm.VirtualMachine
}

func NewContainer(id ID) *Handle {
	con := &Container{
		ID: id,
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

func (c *Container) Commit(ctx context.Context, sess *session.Session, s *types.VirtualMachineConfigSpec) error {
	c.Lock()
	defer c.Unlock()

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
			return parent.CreateVM(ctx, *s, sess.Pool, host)
		})
		if err != nil {
			return err
		}

		c.vm = vm.NewVirtualMachine(ctx, sess, res.Result.(types.ManagedObjectReference))
	}

	return nil
}

// Start starts a container vm with the given params
func (c *Container) Start(ctx context.Context) error {
	defer trace.End(trace.Begin("Container.Start"))

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

	var detail string
	waitFunc := func(pc []types.PropertyChange) bool {
		// guestinfo key that we want to wait for
		key := fmt.Sprintf("guestinfo..sessions|%s.started", c.ID)

		for _, c := range pc {
			if c.Op != types.PropertyChangeOpAssign {
				continue
			}

			values := c.Val.(types.ArrayOfOptionValue).OptionValue
			for _, value := range values {
				// check the status of the key and return true if it's been set to non-nil
				if key == value.GetOptionValue().Key {
					detail = value.GetOptionValue().Value.(string)
					return detail != "" && detail != "<nil>"
				}
			}
		}
		return false
	}

	// Wait some before giving up...
	ctx, cancel := context.WithTimeout(ctx, propertyCollectorTimeout)
	defer cancel()

	err = c.vm.WaitForExtraConfig(ctx, waitFunc)
	if err != nil {
		return fmt.Errorf("unable to wait for process launch status: %s", err.Error())
	}

	if detail != "true" {
		return errors.New(detail)
	}

	return nil
}
