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
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/migration"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	log "github.com/Sirupsen/logrus"
)

// NotYetExistError is returned when a call that requires a VM exist is made
type NotYetExistError struct {
	ID string
}

func (e NotYetExistError) Error() string {
	return fmt.Sprintf("%s is not completely created", e.ID)
}

// containerBase holds fields common between Handle and Container. The fields and
// methods in containerBase should not require locking as they're primary use is:
// a. for read-only reference when used in Container
// b. single use/no-concurrent modification when used in Handle
type containerBase struct {
	ExecConfig *executor.ExecutorConfig

	// Migrated is used during in memory migration to assign whether an execConfig is viable for a commit phase
	Migrated bool
	// MigrationError means the errors happens during data migration, some operation might fail for we cannot extract the whole container configuration
	MigrationError error
	DataVersion    int

	// original - can be pointers so long as refreshes
	// use different instances of the structures
	Config  *types.VirtualMachineConfigInfo
	Runtime *types.VirtualMachineRuntimeInfo

	// doesn't change so can be copied here
	vm *vm.VirtualMachine
}

func newBase(vm *vm.VirtualMachine, c *types.VirtualMachineConfigInfo, r *types.VirtualMachineRuntimeInfo) *containerBase {
	base := &containerBase{
		ExecConfig: &executor.ExecutorConfig{},
		Config:     c,
		Runtime:    r,
		vm:         vm,
	}

	// construct a working copy of the exec config
	if c != nil && c.ExtraConfig != nil {
		var migratedConf map[string]string
		containerExecKeyValues := vmomi.OptionValueMap(c.ExtraConfig)
		// #nosec: Errors unhandled.
		base.DataVersion, _ = migration.ContainerDataVersion(containerExecKeyValues)
		migratedConf, base.Migrated, base.MigrationError = migration.MigrateContainerConfig(containerExecKeyValues)
		extraconfig.Decode(extraconfig.MapSource(migratedConf), base.ExecConfig)
	}

	return base
}

// VMReference will provide the vSphere vm managed object reference
func (c *containerBase) VMReference() types.ManagedObjectReference {
	var moref types.ManagedObjectReference
	if c.vm != nil {
		moref = c.vm.Reference()
	}
	return moref
}

// unlocked refresh of container state
func (c *containerBase) refresh(ctx context.Context) error {
	base, err := c.updates(ctx)
	if err != nil {
		log.Errorf("Update: unable to update container %s", c.ExecConfig.ID)
		return err
	}

	// copy over the new state
	*c = *base
	return nil
}

// updates acquires updates from the infrastructure without holding a lock
func (c *containerBase) updates(ctx context.Context) (*containerBase, error) {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	var o mo.VirtualMachine

	// make sure we have vm
	if c.vm == nil {
		return nil, NotYetExistError{c.ExecConfig.ID}
	}

	if c.Config != nil {
		log.Debugf("Update: refreshing from change version %s", c.Config.ChangeVersion)
	}

	if err := c.vm.Properties(ctx, c.vm.Reference(), []string{"config", "runtime"}, &o); err != nil {
		return nil, err
	}

	base := &containerBase{
		vm:         c.vm,
		Config:     o.Config,
		Runtime:    &o.Runtime,
		ExecConfig: &executor.ExecutorConfig{},
	}

	// Get the ExtraConfig
	var migratedConf map[string]string
	containerExecKeyValues := vmomi.OptionValueMap(o.Config.ExtraConfig)
	if containerExecKeyValues["guestinfo.vice./common/id"] == "" {
		return nil, fmt.Errorf("Update: change version %s failed assertion extraconfig id != nil", o.Config.ChangeVersion)
	}

	log.Debugf("Update: change version %s, extraconfig id: %+v", o.Config.ChangeVersion, containerExecKeyValues["guestinfo.vice./common/id"])
	// #nosec: Errors unhandled.
	base.DataVersion, _ = migration.ContainerDataVersion(containerExecKeyValues)
	migratedConf, base.Migrated, base.MigrationError = migration.MigrateContainerConfig(containerExecKeyValues)
	extraconfig.Decode(extraconfig.MapSource(migratedConf), base.ExecConfig)

	return base, nil
}

func (c *containerBase) ReloadConfig(ctx context.Context) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	return c.startGuestProgram(ctx, "reload", "")
}

// WaitForExec waits exec'ed task to set started field or timeout
func (c *containerBase) WaitForExec(ctx context.Context, id string) error {
	defer trace.End(trace.Begin(id))

	return c.waitForExec(ctx, id)
}

// WaitForSession waits non-exec'ed task to set started field or timeout
func (c *containerBase) WaitForSession(ctx context.Context, id string) error {
	defer trace.End(trace.Begin(id))

	return c.waitForSession(ctx, id)
}

func (c *containerBase) startGuestProgram(ctx context.Context, name string, args string) error {
	// make sure we have vm
	if c.vm == nil {
		return NotYetExistError{c.ExecConfig.ID}
	}

	defer trace.End(trace.Begin(c.ExecConfig.ID + ":" + name))
	o := guest.NewOperationsManager(c.vm.Client.Client, c.vm.Reference())
	m, err := o.ProcessManager(ctx)
	if err != nil {
		return err
	}

	spec := types.GuestProgramSpec{
		ProgramPath: name,
		Arguments:   args,
	}

	auth := types.NamePasswordAuthentication{
		Username: c.ExecConfig.ID,
	}

	_, err = m.StartProgram(ctx, &auth, &spec)

	return err
}

func (c *containerBase) start(ctx context.Context) error {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	// make sure we have vm
	if c.vm == nil {
		return NotYetExistError{c.ExecConfig.ID}
	}

	// Power on
	_, err := c.vm.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return c.vm.PowerOn(ctx)
	})

	return err
}

func (c *containerBase) stop(ctx context.Context, waitTime *int32) error {
	// make sure we have vm
	if c.vm == nil {
		return NotYetExistError{c.ExecConfig.ID}
	}

	// get existing state and set to stopping
	// if there's a failure we'll revert to existing

	err := c.shutdown(ctx, waitTime)
	if err == nil {
		return nil
	}

	log.Warnf("stopping %s via hard power off due to: %s", c.ExecConfig.ID, err)

	return c.poweroff(ctx)
}

func (c *containerBase) kill(ctx context.Context) error {
	// make sure we have vm
	if c.vm == nil {
		return NotYetExistError{c.ExecConfig.ID}
	}

	wait := 10 * time.Second // default
	timeout, cancel := context.WithTimeout(ctx, wait)
	defer cancel()

	sig := string(ssh.SIGKILL)
	log.Infof("sending kill -%s %s", sig, c.ExecConfig.ID)

	err := c.startGuestProgram(timeout, "kill", sig)
	if err == nil && timeout.Err() != nil {
		log.Warnf("timeout (%s) waiting for %s to power off via SIG%s", wait, c.ExecConfig.ID, sig)
	}
	if err != nil {
		log.Warnf("killing %s attempt resulted in: %s", c.ExecConfig.ID, err)
	}

	// Even if startGuestProgram failed above, it may actually have executed.  If the container came up and then
	// we kill it before VC gets a chance to detect the toolbox, vSphere can execute the kill but report an
	// error 3016 indicating the guest toolbox wasn't found.  If we then try to poweroff, it may throw vSphere
	// into an invalid transition and will need to recover.  If we try to grab properties at this time, the
	// power state may be incorrect.  We work around this by waiting on the power state, regardless of error
	// from startGuestProgram. https://github.com/vmware/vic/issues/5803
	log.Infof("waiting %s for %s to power off", wait, c.ExecConfig.ID)
	err = c.vm.WaitForPowerState(timeout, types.VirtualMachinePowerStatePoweredOff)
	if err == nil {
		return nil // VM has powered off
	}

	log.Warnf("killing %s via hard power off", c.ExecConfig.ID)

	// stop wait time is not applied for the hard kill
	return c.poweroff(ctx)
}

func (c *containerBase) shutdown(ctx context.Context, waitTime *int32) error {
	// make sure we have vm
	if c.vm == nil {
		return NotYetExistError{c.ExecConfig.ID}
	}

	wait := 10 * time.Second // default
	if waitTime != nil && *waitTime > 0 {
		wait = time.Duration(*waitTime) * time.Second
	}

	cs := c.ExecConfig.Sessions[c.ExecConfig.ID]
	stop := []string{cs.StopSignal, string(ssh.SIGKILL)}
	if stop[0] == "" {
		stop[0] = string(ssh.SIGTERM)
	}

	for _, sig := range stop {
		msg := fmt.Sprintf("sending kill -%s %s", sig, c.ExecConfig.ID)
		log.Info(msg)

		timeout, cancel := context.WithTimeout(ctx, wait)
		defer cancel()

		err := c.startGuestProgram(timeout, "kill", sig)
		if err != nil {
			// Just warn and proceed to waiting for power state per issue https://github.com/vmware/vic/issues/5803
			// Description above in function kill()
			log.Warnf("%s: %s", msg, err)
		}

		log.Infof("waiting %s for %s to power off", wait, c.ExecConfig.ID)
		err = c.vm.WaitForPowerState(timeout, types.VirtualMachinePowerStatePoweredOff)
		if err == nil {
			return nil // VM has powered off
		}

		if timeout.Err() == nil {
			return err // error other than timeout
		}

		log.Warnf("timeout (%s) waiting for %s to power off via SIG%s", wait, c.ExecConfig.ID, sig)
	}

	return fmt.Errorf("failed to shutdown %s via kill signals %s", c.ExecConfig.ID, stop)
}

func (c *containerBase) poweroff(ctx context.Context) error {
	// make sure we have vm
	if c.vm == nil {
		return NotYetExistError{c.ExecConfig.ID}
	}

	_, err := c.vm.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return c.vm.PowerOff(ctx)
	})

	if err != nil {

		// It is possible the VM has finally shutdown in between, ignore the error in that case
		if terr, ok := err.(task.Error); ok {
			switch terr := terr.Fault().(type) {
			case *types.InvalidPowerState:
				if terr.ExistingState == types.VirtualMachinePowerStatePoweredOff {
					log.Warnf("power off %s task skipped (state was already %s)", c.ExecConfig.ID, terr.ExistingState)
					return nil
				}
				log.Warnf("invalid power state during power off: %s", terr.ExistingState)

			case *types.GenericVmConfigFault:

				// Check if the poweroff task was canceled due to a concurrent guest shutdown
				if len(terr.FaultMessage) > 0 {
					k := terr.FaultMessage[0].Key
					if k == vmNotSuspendedKey || k == vmPoweringOffKey {
						log.Infof("power off %s task skipped due to guest shutdown", c.ExecConfig.ID)
						return nil
					}
				}
				log.Warnf("generic vm config fault during power off: %#v", terr)

			default:
				log.Warnf("hard power off failed due to: %#v", terr)
			}
		}

		return err
	}

	return nil
}

func (c *containerBase) waitForPowerState(ctx context.Context, max time.Duration, state types.VirtualMachinePowerState) (bool, error) {
	defer trace.End(trace.Begin(c.ExecConfig.ID))

	timeout, cancel := context.WithTimeout(ctx, max)
	defer cancel()

	err := c.vm.WaitForPowerState(timeout, state)
	if err != nil {
		return timeout.Err() != nil, err
	}

	return false, nil
}

func (c *containerBase) waitForSession(ctx context.Context, id string) error {
	defer trace.End(trace.Begin(id))

	// guestinfo key that we want to wait for
	key := extraconfig.CalculateKeys(c.ExecConfig, fmt.Sprintf("Sessions.%s.Started", id), "")[0]
	return c.waitFor(ctx, key)
}

func (c *containerBase) waitForExec(ctx context.Context, id string) error {
	defer trace.End(trace.Begin(id))

	// guestinfo key that we want to wait for
	key := extraconfig.CalculateKeys(c.ExecConfig, fmt.Sprintf("Execs.%s.Started", id), "")[0]
	return c.waitFor(ctx, key)
}

func (c *containerBase) waitFor(ctx context.Context, key string) error {
	detail, err := c.vm.WaitForKeyInExtraConfig(ctx, key)
	if err != nil {
		return fmt.Errorf("unable to wait for process launch status: %s", err)
	}

	if detail != "true" {
		return fmt.Errorf("%s", detail)
	}

	return nil
}
