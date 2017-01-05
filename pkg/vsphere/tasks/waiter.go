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
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
)

const (
	maxBackoffFactor = int64(16)

	VMObjectKey = "VMObject"
)

var (
	m             sync.RWMutex
	errorHandlers []ErrorHandler
)

type ErrorHandler func(ctx context.Context, err error) (bool, error)

func taskInProgressHandler(ctx context.Context, err error) (bool, error) {
	if isTaskInProgress(err) {
		// TaskInProgress error, no need to fail here
		log.Debugf("TaskInProgress error, continue")
		return true, nil
	}
	return false, nil
}

func init() {
	RegisterErrorHandler(taskInProgressHandler)
}

func RegisterErrorHandler(handler ErrorHandler) {
	defer trace.End(trace.Begin(""))
	m.Lock()
	defer m.Unlock()
	errorHandlers = append(errorHandlers, handler)
}

//FIXME: remove this type and refactor to use object.Task from govmomi
//       this will require a lot of code being touched in a lot of places.
type Task interface {
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

		handled, herr := HandleError(ctx, err)
		if herr != nil {
			log.Debugf("Handler failed: %s", herr)
			return info, err
		}
		if !handled {
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

func HandleError(ctx context.Context, err error) (bool, error) {
	m.RLock()
	defer m.RUnlock()

	handled := false
	for i := range errorHandlers {
		expected, herr := errorHandlers[i](ctx, err)
		if !expected {
			continue
		}
		handled = true
		if herr != nil {
			return handled, herr
		}
		break
	}
	return handled, nil
}

func isTaskInProgress(err error) bool {
	if soap.IsSoapFault(err) {
		switch f := soap.ToSoapFault(err).VimFault().(type) {
		case types.TaskInProgress:
			return true
		default:
			logSoapFault(f)
		}
	}

	if soap.IsVimFault(err) {
		switch f := soap.ToVimFault(err).(type) {
		case *types.TaskInProgress:
			return true
		default:
			logFault(f)
		}
	}

	switch err := err.(type) {
	case task.Error:
		if _, ok := err.Fault().(*types.TaskInProgress); ok {
			return true
		}
		logFault(err.Fault())
	default:
		if f, ok := err.(types.HasFault); ok {
			logFault(f.Fault())
		} else {
			logError(err)
		}
	}
	return false
}

// Helper Functions
func logFault(fault types.BaseMethodFault) {
	log.Debugf("unexpected fault on task retry : %#v", fault)
}

func logSoapFault(fault types.AnyType) {
	log.Debugf("unexpected soap fault on task retry : %#v", fault)
}

func logError(err error) {
	log.Debugf("unexpected error on task retry : %#v", err)
}
