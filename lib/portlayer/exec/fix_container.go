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
	"context"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

type ContainerFixer struct {
}

func InitErrorHandler() {
	tasks.RegisterErrorHandler(&ContainerFixer{})
}

func (h *ContainerFixer) Handle(ctx context.Context, err error) (bool, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("error: %s", err)))
	o := ctx.Value(tasks.VMContextObjectKey)
	if o == nil {
		log.Debugf("No vm object set, not vm operations.")
		return false, nil
	}

	vm, ok := o.(*vm.VirtualMachine)
	if !ok {
		log.Debugf("Not vm object, do not fix failure")
		return false, nil
	}
	if !vm.IsInvalidState(ctx) {
		log.Debugf("VM is not in invalid state, do not fix failure")
		return false, nil
	}
	log.Debugf("Try to fix failure %s", err)
	container := Containers.Container(vm.Reference().String())
	if container != nil {
		oldState := container.CurrentState()
		defer container.SetState(oldState)
		container.SetState(StateFixing)
	}
	if nerr := vm.FixVM(ctx); nerr != nil {
		log.Errorf("Failed to fix task failure: %s", nerr)
		return true, nerr
	}
	log.Debugf("Fixed")

	return true, nil
}
