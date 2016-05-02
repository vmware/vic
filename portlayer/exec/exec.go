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
	"fmt"
	"math/rand"
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/guest"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/spec"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
	"golang.org/x/net/context"
)

const (
	serialOverLANPort  = 2377
	managementHostName = "management.localhost"
)

// Executor type for exec operations
type Executor struct {
	// Pointer to the session struct
	session *session.Session
}

// ExecutorCreateConfig defines the parameters for Create call
type ExecutorCreateConfig struct {
	Metadata metadata.ExecutorConfig

	ParentImageID  string
	ImageStoreName string
	VCHName        string
}

// ExecutorStartConfig defines the parameters for Start call
type ExecutorStartConfig struct {
	ID string
}

// NewExecutor returns a new Executor
func NewExecutor(ctx context.Context, session *session.Session) (*Executor, error) {
	e := &Executor{
		session: session,
	}

	return e, nil
}

// Create creates a container vm with the given params
func (e *Executor) Create(ctx context.Context, config *ExecutorCreateConfig) error {
	defer trace.End(trace.Begin("Create"))

	// Convert the management hostname to IP
	ips, err := net.LookupIP(managementHostName)
	if err != nil {
		return err
	}
	if len(ips) == 0 {
		return fmt.Errorf("No IP found on %s", managementHostName)
	}
	if len(ips) > 1 {
		return fmt.Errorf("Multiple IPs found on %s: %#v", managementHostName, ips)
	}
	URI := fmt.Sprintf("tcp://%s:%d", ips[0], serialOverLANPort)

	specconfig := &spec.VirtualMachineConfigSpecConfig{
		// FIXME: hardcoded values
		NumCPUs:  2,
		MemoryMB: 2048,

		ConnectorURI: URI,

		ID:   config.Metadata.ID,
		Name: config.Metadata.Name,

		ParentImageID: config.ParentImageID,

		// FIXME: hardcoded value
		BootMediaPath: e.session.Datastore.Path(fmt.Sprintf("%s/bootstrap.iso", config.VCHName)),
		VMPathName:    fmt.Sprintf("[%s]", e.session.Datastore.Name()),
		NetworkName:   strings.Split(e.session.Network.Reference().Value, "-")[1],

		ImageStoreName: config.ImageStoreName,

		Metadata: config.Metadata,
	}
	log.Debugf("Config: %#v", specconfig)

	// Create a linux guest
	linux, err := guest.NewLinuxGuest(ctx, e.session, specconfig)
	if err != nil {
		return err
	}

	// Find the Virtual Machine folder that we use
	folders, err := e.session.Datacenter.Folders(ctx)
	if err != nil {
		return err
	}
	parent := folders.VmFolder

	// FIXME: Replace this simple logic with DRS placement
	// Pick a random host
	hosts, err := e.session.Datastore.AttachedClusterHosts(ctx, e.session.Cluster)
	if err != nil {
		return err
	}
	host := hosts[rand.Intn(len(hosts))]

	// Create the vm
	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return parent.CreateVM(ctx, *linux.Spec(), e.session.Pool, host)
	})
	if err != nil {
		return err
	}

	return nil
}

// Start starts a container vm with the given params
func (e *Executor) Start(ctx context.Context, config *ExecutorStartConfig) error {
	defer trace.End(trace.Begin("Start"))

	// Locate the Virtual Machine
	foundvm, err := e.session.Finder.VirtualMachine(ctx, config.ID)
	if err != nil {
		return err
	}

	// Wrap the result with our version of VirtualMachine
	vm := vm.NewVirtualMachine(ctx, e.session, foundvm.Reference())

	// Power on
	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return vm.PowerOn(ctx)
	})
	if err != nil {
		return err
	}

	return nil
}
