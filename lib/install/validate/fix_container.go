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

package validate

import (
	"context"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// init register fixVMInvalidStateHandler to tasks error handler. So if any container vm operations fail for invalid state issue, this handler will be executed to fix the error.
func init() {
	tasks.RegisterErrorHandler(fixVMInvalidStateHandler)
}

func fixVMInvalidStateHandler(ctx context.Context, err error) (bool, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("error: %s", err)))
	o := ctx.Value(tasks.VMObjectKey)
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
	if nerr := vm.FixVM(ctx); nerr != nil {
		log.Errorf("Failed to fix task failure: %s", nerr)
		return true, nerr
	}
	log.Debugf("Fixed")

	return true, nil
}
