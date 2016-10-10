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

package management

import (
	"context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

func (d *Dispatcher) DebugVCH(vch *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec, password, authorizedKey string) error {
	defer trace.End(trace.Begin(conf.Name))

	op, err := trace.FromContext(d.ctx)
	if err != nil {
		op = trace.NewOperation(d.ctx, "enable appliance debug")
	}

	err = d.enableSSH(op, vch, password, authorizedKey)
	if err != nil {
		op.Errorf("Unable to enable ssh on the VCH appliance VM: %s", err)
		return err
	}

	d.sshEnabled = true

	return nil
}

func (d *Dispatcher) enableSSH(ctx context.Context, vch *vm.VirtualMachine, password, authorizedKey string) error {
	op, err := trace.FromContext(ctx)
	if err != nil {
		op = trace.NewOperation(ctx, "enable ssh in appliance")
	}

	state, err := vch.PowerState(op)
	if err != nil {
		log.Errorf("Failed to get appliance power state, service might not be available at this moment.")
	}
	if state != types.VirtualMachinePowerStatePoweredOn {
		err = errors.Errorf("VCH appliance is not powered on, state %s", state)
		op.Errorf("%s", err)
		return err
	}

	running, err := vch.IsToolsRunning(op)
	if err != nil || !running {
		err = errors.New("Tools is not running in the appliance, unable to continue")
		op.Errorf("%s", err)
		return err
	}

	manager := guest.NewOperationsManager(d.session.Client.Client, vch.Reference())
	processManager, err := manager.ProcessManager(op)
	if err != nil {
		err = errors.Errorf("Unable to manage processes in appliance VM: %s", err)
		op.Errorf("%s", err)
		return err
	}

	auth := types.NamePasswordAuthentication{}

	spec := types.GuestProgramSpec{
		ProgramPath:      "enable-ssh",
		Arguments:        string(authorizedKey),
		WorkingDirectory: "/",
		EnvVariables:     []string{},
	}

	_, err = processManager.StartProgram(op, &auth, &spec)
	if err != nil {
		err = errors.Errorf("Unable to enable SSH in appliance VM: %s", err)
		op.Errorf("%s", err)
		return err
	}

	if password == "" {
		return nil
	}

	// set the password as well
	spec = types.GuestProgramSpec{
		ProgramPath:      "passwd",
		Arguments:        password,
		WorkingDirectory: "/",
		EnvVariables:     []string{},
	}

	_, err = processManager.StartProgram(op, &auth, &spec)
	if err != nil {
		err = errors.Errorf("Unable to enable in appliance VM: %s", err)
		op.Errorf("%s", err)
		return err
	}

	return nil
}
