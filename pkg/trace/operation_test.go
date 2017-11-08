// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

package trace

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestContextUnpack(t *testing.T) {
	Logger.Level = logrus.DebugLevel

	cnt := 100
	wg := &sync.WaitGroup{}
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func(i int) {
			defer wg.Done()
			ctx := NewOperation(context.TODO(), "testmsg")

			// unpack an Operation via the context using it's Values fields
			c := FromContext(ctx, "test")
			c.Infof("test info message %d", i)
		}(i) // fix race in test
	}
	wg.Wait()
}

// If we timeout a child, test a stack is printed of contexts
func TestNestedLogging(t *testing.T) {
	// create a buf to check the log
	buf := new(bytes.Buffer)
	Logger.Out = buf

	root := NewOperation(context.Background(), "root")

	var ctxFunc func(parent Operation, level int) Operation

	levels := 10
	ctxFunc = func(parent Operation, level int) Operation {
		if level == levels {
			return parent
		}

		child, _ := WithDeadline(&parent, time.Time{}, fmt.Sprintf("level %d", level))

		return ctxFunc(child, level+1)
	}

	child := ctxFunc(root, 0)

	// Assert the child has an error and prints a stack.  The parent doesn't
	// see this and should not have an error.  Only cancelation trickles up the
	// stack to the parent.
	if !assert.NoError(t, root.Err()) || !assert.Error(t, child.Err()) {
		return
	}

	// Assert we got a stack trace in the log
	log := buf.String()
	lines := strings.Count(log, "\n")
	t.Log(log)

	// Sample stack
	//
	//        ERRO[0000] op=21598.101: github.com/vmware/vic/pkg/trace.TestNestedLogging: level 9 error: context deadline exceeded
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 9
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 8
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 7
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 6
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 5
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 4
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 3
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 2
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 1
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:71 level 0
	//                        github.com/vmware/vic/pkg/trace.TestNestedLogging:61 root

	// We arrive at 2 because we have the err line (line 0), then the root
	// (line 11) of where we created the ctx.
	if assert.False(t, lines < levels) {
		t.Logf("exepected at least %d and got %d", levels, lines)
		return
	}
}

// Just checking behavior of the context package
func TestSanity(t *testing.T) {
	Logger.Level = logrus.InfoLevel
	levels := 10

	root, cancel := context.WithDeadline(context.Background(), time.Time{})
	defer cancel()

	var ctxFunc func(parent context.Context, level int) context.Context

	ctxFunc = func(parent context.Context, level int) context.Context {
		if level == levels {
			return parent
		}

		child, cancel := context.WithDeadline(parent, time.Now().Add(time.Hour))
		defer cancel()

		return ctxFunc(child, level+1)
	}

	child := ctxFunc(root, 0)

	if !assert.Error(t, child.Err()) {
		t.FailNow()
	}

	err := root.Err()
	if !assert.Error(t, err) {
		return
	}
}
