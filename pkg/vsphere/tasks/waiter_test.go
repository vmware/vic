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

package tasks

import (
	"context"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
)

type MyTask struct {
	success bool
}

func (t *MyTask) Wait(ctx context.Context) error {
	_, err := t.WaitForResult(ctx, nil)
	return err
}
func (t *MyTask) WaitForResult(ctx context.Context, s progress.Sinker) (*types.TaskInfo, error) {
	if t.success {
		return nil, nil
	}
	return nil, errors.Errorf("Wait failed")
}

func createFailedTask(context.Context) (Task, error) {
	return nil, errors.Errorf("Create VM failed")
}

func createFailedResultWaiter(context.Context) (Task, error) {
	task := &MyTask{
		false,
	}
	return task, nil
}

func createResultWaiter(context.Context) (Task, error) {
	task := &MyTask{
		true,
	}
	return task, nil
}

func TestFailedInvokeResult(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.TODO()
	_, err := WaitForResult(ctx, func(ctx context.Context) (Task, error) {
		return createFailedTask(ctx)
	})
	if err == nil || !strings.Contains(err.Error(), "Create VM failed") {
		t.Errorf("Not expected error message")
	}
}

func TestFailedWaitResult(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.TODO()
	_, err := WaitForResult(ctx, func(ctx context.Context) (Task, error) {
		return createFailedResultWaiter(ctx)
	})
	log.Debugf("got error: %s", err.Error())
	if err == nil || !strings.Contains(err.Error(), "Wait failed") {
		t.Errorf("Not expected error message")
	}
}

func TestSuccessWaitResult(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.TODO()
	_, err := WaitForResult(ctx, func(ctx context.Context) (Task, error) {
		return createResultWaiter(ctx)
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func createFailed(context.Context) (Task, error) {
	return nil, errors.Errorf("Create VM failed")
}

func createFailedWaiter(context.Context) (Task, error) {
	task := &MyTask{
		false,
	}
	return task, nil
}

func createWaiter(context.Context) (Task, error) {
	task := &MyTask{
		true,
	}
	return task, nil
}

func TestFailedInvoke(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.TODO()
	err := Wait(ctx, func(ctx context.Context) (Task, error) {
		return createFailed(ctx)
	})
	if err == nil || !strings.Contains(err.Error(), "Create VM failed") {
		t.Errorf("Not expected error message")
	}
}

func TestFailedWait(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.TODO()
	err := Wait(ctx, func(ctx context.Context) (Task, error) {
		return createFailedWaiter(ctx)
	})
	log.Debugf("got error: %s", err.Error())
	if err == nil || !strings.Contains(err.Error(), "Wait failed") {
		t.Errorf("Not expected error message")
	}
}

func TestSuccessWait(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.TODO()
	err := Wait(ctx, func(ctx context.Context) (Task, error) {
		return createWaiter(ctx)
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}
