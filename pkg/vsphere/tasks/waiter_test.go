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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)

	m.Run()
}

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
	ctx := context.TODO()
	_, err := WaitForResult(ctx, func(ctx context.Context) (Task, error) {
		return createFailedTask(ctx)
	})
	if err == nil || !strings.Contains(err.Error(), "Create VM failed") {
		t.Errorf("Not expected error message")
	}
}

func TestFailedWaitResult(t *testing.T) {
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
	ctx := context.TODO()
	err := Wait(ctx, func(ctx context.Context) (Task, error) {
		return createFailed(ctx)
	})
	if err == nil || !strings.Contains(err.Error(), "Create VM failed") {
		t.Errorf("Not expected error message")
	}
}

func TestFailedWait(t *testing.T) {
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
	ctx := context.TODO()
	err := Wait(ctx, func(ctx context.Context) (Task, error) {
		return createWaiter(ctx)
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

var taskInProgressFault = task.Error{
	LocalizedMethodFault: &types.LocalizedMethodFault{
		Fault: &types.TaskInProgress{},
	},
}

type taskInProgressTask struct {
	cur, max int
	err      error
	info     *types.TaskInfo
}

func (t *taskInProgressTask) Wait(ctx context.Context) error {
	t.cur++
	if t.cur == t.max {
		return t.err
	}

	return taskInProgressFault
}

func (t *taskInProgressTask) WaitForResult(ctx context.Context, s progress.Sinker) (*types.TaskInfo, error) {
	return t.info, t.Wait(ctx)
}

func mustRunInTime(t *testing.T, d time.Duration, f func()) {
	done := make(chan bool)

	go func() {
		f()
		close(done)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	select {
	case <-done: // ran within alloted time
	case <-ctx.Done():
		t.Fatalf("test did not run in alloted time %s", d)
	}
}

func TestRetry(t *testing.T) {
	mustRunInTime(t, 2*time.Second, func() {
		ctx := context.Background()
		i := 0
		ti, err := WaitForResult(ctx, func(_ context.Context) (Task, error) {
			i++
			return nil, assert.AnError
		})

		assert.Nil(t, ti)
		assert.Equal(t, i, 1)
		assert.Error(t, err)
		assert.Equal(t, err, assert.AnError)

		// error != TaskInProgress during task creation
		i = 0
		e := &task.Error{
			LocalizedMethodFault: &types.LocalizedMethodFault{
				Fault:            &types.RuntimeFault{}, // random fault != TaskInProgress
				LocalizedMessage: "random fault",
			},
		}
		ti, err = WaitForResult(ctx, func(_ context.Context) (Task, error) {
			i++
			return nil, e
		})

		assert.Nil(t, ti)
		assert.Equal(t, i, 1)
		assert.Error(t, err)
		assert.Equal(t, err, e)

		// context cancelled after two retries
		i = 0
		ctx, cancel := context.WithCancel(ctx)
		ti, err = WaitForResult(ctx, func(_ context.Context) (Task, error) {
			i++
			if i == 2 {
				cancel()
			}
			return nil, taskInProgressFault
		})

		assert.Nil(t, ti)
		assert.Equal(t, i, 2)
		assert.Error(t, err)
		assert.Equal(t, err, ctx.Err())

		// TaskInProgress from task creation for 2 iterations and
		// then nil error
		tsk := &taskInProgressTask{
			max: 1,
			info: &types.TaskInfo{
				Task: types.ManagedObjectReference{
					Type:  "task",
					Value: "foo",
				},
			},
		}
		i = 0
		ti, err = WaitForResult(context.Background(), func(_ context.Context) (Task, error) {
			i++
			if i == 2 {
				return tsk, nil
			}
			return nil, taskInProgressFault
		})

		assert.Equal(t, tsk.info, ti)
		assert.Equal(t, i, 2)
		assert.NoError(t, err)

		// return TaskInPregress from task.WaitForResult for 2 iterations
		// and then return assert.AnError
		tsk = &taskInProgressTask{
			max: 2,
			err: assert.AnError,
			info: &types.TaskInfo{
				Task: types.ManagedObjectReference{
					Type:  "task",
					Value: "foo",
				},
			},
		}
		ti, err = WaitForResult(context.Background(), func(_ context.Context) (Task, error) {
			return tsk, nil
		})

		assert.Equal(t, tsk.info, ti)
		assert.Equal(t, tsk.max, tsk.cur)
		assert.Error(t, err)
		assert.Equal(t, err, tsk.err)

		// return TaskInPregress from task.WaitForResult for 2 iterations
		// and then return nil error
		tsk.cur = 0
		tsk.err = nil
		ti, err = WaitForResult(context.Background(), func(_ context.Context) (Task, error) {
			return tsk, nil
		})

		assert.Equal(t, tsk.info, ti)
		assert.Equal(t, tsk.info, ti)
		assert.Equal(t, tsk.cur, tsk.max)
		assert.NoError(t, err)
	})
}
