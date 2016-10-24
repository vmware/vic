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

// Package tasks wraps the operation of VC. It will invoke the operation and wait
// until it's finished, and then return the execution result or error message.
package tasks

import (
	"context"
	"errors"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	maxBackoffFactor = int64(16)
)

//FIXME: remove this type and refactor to use object.Task from govmomi
//       this will require a lot of code being touched in a lot of places.
type Task interface {
	Wait(ctx context.Context) error
	WaitForResult(ctx context.Context, s progress.Sinker) (*types.TaskInfo, error)
}

// Wait wraps govmomi operations and wait the operation to complete
// Sample usage:
//    info, err := Wait(ctx, func(ctx), sess *session.Session, *object.Reference, (*object.Reference, *TaskInfo, error) {
//       return vm, vm.Reconfigure(ctx, sess, target, config)
//    })
func Wait(ctx context.Context, sess *session.Session, target object.Reference, f func(context.Context) (Task, error)) error {
	_, err := WaitForResult(ctx, sess, target, f)
	return err
}

// WaitForResult wraps govmomi operations and wait the operation to complete.
// Return the operation result
// Sample usage:
//    info, err := WaitForResult(ctx, sess *session.Session, *object.Reference, func(ctx) (*TaskInfo, error) {
//       return vm, vm.Reconfigure(ctx, config)
//    })
func WaitForResult(ctx context.Context, sess *session.Session, target object.Reference, f func(context.Context) (Task, error)) (*types.TaskInfo, error) {
	var err error
	var info *types.TaskInfo
	var backoffFactor int64 = 1

	for {
		var t Task
		if t, err = f(ctx); err == nil {
			info, err = t.WaitForResult(ctx, nil)
			if err == nil {
				return info, err
			}
		}

		log.Errorf("task failed: %s", err)
		if !needsRetry(target, err) {
			return info, err
		}

		if needsFix(target, err) {
			if nerr := fixTask(ctx, sess, target); nerr != nil {
				log.Errorf("Failed to fix task failure: %s", nerr)
				return info, err
			}
			log.Debugf("Fixed error: %s", err)
		}

		sleepValue := time.Duration(backoffFactor * (rand.Int63n(100) + int64(50)))
		select {
		case <-time.After(sleepValue * time.Millisecond):
			backoffFactor *= 2
			if backoffFactor > maxBackoffFactor {
				backoffFactor = maxBackoffFactor
			}
		case <-ctx.Done():
			return info, ctx.Err()
		}

		log.Warnf("retrying task")
	}
}

func fixTask(ctx context.Context, sess *session.Session, target object.Reference) error {
	switch vmm := target.(type) {
	case *object.VirtualMachine:
		vm := vm.NewVirtualMachine(ctx, sess, target.Reference())
		err := vm.FixInvalidState(ctx)
		// make sure original object is updated, to make sure task retry works
		vmm.Common = vm.Common
		return err
	case *vm.VirtualMachine:
		return vmm.FixInvalidState(ctx)
	default:
		return errors.New("Unable to fix non-vm object")
	}
}

// check if task is in progress, or vm is in invalid state.
func needsRetry(target object.Reference, err error) bool {
	if isTaskInProgress(err) {
		return true
	}
	if needsFix(target, err) {
		return true
	}
	return false
}

func needsFix(target object.Reference, err error) bool {
	f, ok := err.(types.HasFault)
	if !ok {
		return false
	}
	switch f.Fault().(type) {
	case *types.InvalidState:
	default:
		log.Debugf("Do not fix non invalid state error")
		return false
	}
	if target == nil {
		log.Debugf("Do not fix nil object")
		return false
	}
	switch target.(type) {
	case *object.VirtualMachine, *vm.VirtualMachine:
		return true
	default:
		log.Debugf("Unable to fix non-vm object: %#v", target)
		return false
	}
}

func isTaskInProgress(err error) bool {
	if soap.IsSoapFault(err) {
		if _, ok := soap.ToSoapFault(err).VimFault().(types.TaskInProgress); ok {
			return true
		}
	}

	if soap.IsVimFault(err) {
		switch soap.ToVimFault(err).(type) {
		case *types.TaskInProgress:
			return true
		}
	}

	switch err := err.(type) {
	case task.Error:
		if _, ok := err.Fault().(*types.TaskInProgress); ok {
			return true
		}
	}

	return false
}
