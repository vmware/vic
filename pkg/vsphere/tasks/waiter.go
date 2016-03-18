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
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"

	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/types"
)

type Waiter interface {
	Wait(ctx context.Context) error
}

type ResultWaiter interface {
	WaitForResult(ctx context.Context, s progress.Sinker) (*types.TaskInfo, error)
}

// Wait wraps govmomi operations and wait the operation to complete
// Sample usage:
//    info, err := Wait(ctx, func(ctx) (*TaskInfo, error) {
//       return vm.Reconfigure(ctx, config)
//    })
func Wait(ctx context.Context, f func(context.Context) (Waiter, error)) error {
	task, err := f(ctx)
	if err != nil {
		cerr := errors.Errorf("Failed to invoke operation: %s", errors.ErrorStack(err))
		log.Errorf(cerr.Error())
		return cerr
	}

	err = task.Wait(ctx)
	if err != nil {
		cerr := errors.Errorf("Operation failed: %s", errors.ErrorStack(err))
		log.Errorf(cerr.Error())
		return cerr
	}
	return nil
}

// WaitForResult wraps govmomi operations and wait the operation to complete.
// Return the operation result
// Sample usage:
//    info, err := WaitForResult(ctx, func(ctx) (*TaskInfo, error) {
//       return vm.Reconfigure(ctx, config)
//    })
func WaitForResult(ctx context.Context, f func(context.Context) (ResultWaiter, error)) (*types.TaskInfo, error) {
	task, err := f(ctx)
	if err != nil {
		cerr := errors.Errorf("Failed to invoke operation: %s", errors.ErrorStack(err))
		log.Errorf(cerr.Error())
		return nil, cerr
	}

	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		cerr := errors.Errorf("Operation failed: %s", errors.ErrorStack(err))
		if info != nil && info.Error != nil {
			cerr = errors.Errorf("%s - (%s)", cerr, info.Error)
		}

		log.Errorf(cerr.Error())
		return nil, cerr
	}
	return info, nil
}
