// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package task

import (
	"fmt"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
)

// Wait waits the task to start
func Wait(op *trace.Operation, h interface{}, id string) error {
	defer trace.End(trace.Begin(id, op))

	handle, ok := h.(*exec.Handle)
	if !ok {
		return fmt.Errorf("type assertion failed for %#+v", handle)
	}

	if handle.Runtime != nil && handle.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOn {
		err := fmt.Errorf("unable to wait for task when container %s is not running", handle.ExecConfig.ID)
		op.Errorf("%s", err)
		return TaskPowerStateError{Err: err}
	}

	_, okS := handle.ExecConfig.Sessions[id]
	_, okE := handle.ExecConfig.Execs[id]

	if !okS && !okE {
		return fmt.Errorf("unknown task ID: %s", id)
	}

	timeout, cancel := trace.WithTimeout(op, constants.PropertyCollectorTimeout, "Wait")
	defer cancel()

	c := exec.Containers.Container(handle.ExecConfig.ID)
	if c == nil {
		return fmt.Errorf("unknown container ID: %s", handle.ExecConfig.ID)
	}

	if okS {
		return c.WaitForSession(timeout, id)
	}

	// This should get a handle from the callers and wait on the handle since they only have a snap shot of the world. We need to determine a change from what the caller knows... This should really only be during an Exec, we do not want to affect the current WaitForSession functionality. At least for now.

	// we should not wait for an exec that is failed or stopped. it will not change again.
	currentState, err := handle.ExecState(*op, id)
	if err != nil {
		op.Errorf("Unable to find state of exec: exec(%s), container(%s)", id, handle.ExecConfig.ID)
		return err
	}

	if currentState == exec.ExecFailedState || currentState == exec.ExecStoppedState || currentState == exec.ExecUnknownState {
		op.Debugf("Exec already halted: exec(%s), container(%s)", id, handle.ExecConfig.ID)
		return nil
	}

	return handle.WaitForExec(timeout, id)
}

type TaskPowerStateError struct {
	Err error
}

func (t TaskPowerStateError) Error() string {
	return t.Err.Error()
}
