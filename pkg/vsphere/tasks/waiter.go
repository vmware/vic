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
	"fmt"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/types"
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
//    info, err := Wait(ctx, func(ctx) (*TaskInfo, error) {
//       return vm.Reconfigure(ctx, config)
//    })
func Wait(ctx context.Context, f func(context.Context) (Task, error)) error {
	_, err := WaitForResult(ctx, f)
	return err
}

// WaitForResult wraps govmomi operations and wait the operation to complete.
// Return the operation result
// Sample usage:
//    info, err := WaitForResult(ctx, func(ctx) (*TaskInfo, error) {
//       return vm.Reconfigure(ctx, config)
//    })
func WaitForResult(ctx context.Context, f func(context.Context) (Task, error)) (*types.TaskInfo, error) {
	var err error
	var taskInfo *types.TaskInfo
	backoffFactor := int64(1)

	t, err := f(ctx)
	if err != nil {
		log.Errorf("Failed to invoke task operation : %s", err)
		return nil, err
	}

	for {
		taskInfo, werr := t.WaitForResult(ctx, nil)

		if werr == nil {
			return taskInfo, nil
		}

		if terr, ok := werr.(*task.Error); ok {
			if _, ok := terr.Fault().(*types.TaskInProgress); ok {
				sleepValue := time.Duration(backoffFactor * (rand.Int63n(100) + int64(50)))
				select {
				case <-time.After(sleepValue * time.Millisecond):
					if backoffFactor*2 > maxBackoffFactor {
						backoffFactor = maxBackoffFactor
					} else {
						backoffFactor *= 2
					}
				case <-ctx.Done():
					err = fmt.Errorf("%s while retrying task %#v", ctx.Err(), taskInfo)
					log.Error(err)
					return nil, err
				}
				log.Warnf("Retrying Task due to TaskInProgressFault: %s", taskInfo.Task.Reference())
				continue
			}
		}
		err = werr
		break
	}
	log.Errorf("Task failed with error : %s", err)
	return taskInfo, err

}
