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
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	maxBackoffFactor = int64(16)
)

//FIXME: remove this type and refactor to use object.Task from govmomi
//       this will require a lot of code being touched in a lot of places.
type Task interface {
	Reference() types.ManagedObjectReference
	Wait(ctx context.Context) error
	WaitForResult(ctx context.Context, s progress.Sinker) (*types.TaskInfo, error)
}

// Wait wraps govmomi operations and wait the operation to complete
// Sample usage:
//    info, err := Wait(ctx, func(ctx), (*object.Reference, *TaskInfo, error) {
//       return vm, vm.Reconfigure(ctx, config)
//    })
func Wait(ctx context.Context, f func(context.Context) (Task, error)) error {
	_, err := WaitForResult(ctx, f)
	return err
}

// WaitForResult wraps govmomi operations and wait the operation to complete.
// Return the operation result
// Sample usage:
//    info, err := WaitForResult(ctx, func(ctx) (*TaskInfo, error) {
//       return vm, vm.Reconfigure(ctx, config)
//    })
func WaitForResult(ctx context.Context, f func(context.Context) (Task, error)) (*types.TaskInfo, error) {
	var info *types.TaskInfo
	var backoffFactor int64 = 1

	for {
		t, err := f(ctx)
		log.Debugf("New task: %s", taskMoid(t))

		if err == nil {
			info, err = t.WaitForResult(ctx, nil)
			if err == nil {
				return info, err
			}
		}

		if !isTaskInProgress(t, err) {
			return info, err
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

func isTaskInProgress(t Task, err error) bool {
	if soap.IsSoapFault(err) {
		switch f := soap.ToSoapFault(err).VimFault().(type) {
		case types.TaskInProgress:
			return true
		default:
			logSoapFault(t, f)
		}
	}

	if soap.IsVimFault(err) {
		switch f := soap.ToVimFault(err).(type) {
		case *types.TaskInProgress:
			return true
		default:
			logFault(t, f)
		}
	}

	switch err := err.(type) {
	case task.Error:
		if _, ok := err.Fault().(*types.TaskInProgress); ok {
			return true
		}
		logFault(t, err.Fault())
	default:
		if f, ok := err.(types.HasFault); ok {
			logFault(t, f.Fault())
		} else {
			logError(t, err)
		}
	}
	return false
}

// Helper Functions
func logFault(t Task, fault types.BaseMethodFault) {
	log.Errorf("%s: unexpected fault on task retry : %#v", taskMoid(t), fault)
}

func logSoapFault(t Task, fault types.AnyType) {
	log.Errorf("%s: unexpected soap fault on task retry : %#v", taskMoid(t), fault)
}

func logError(t Task, err error) {
	log.Errorf("%s: unexpected error on task retry : %#v", taskMoid(t), err)
}

func taskMoid(t Task) string {
	if t == nil {
		return "Unknown task"
	}

	return t.Reference().Value
}
