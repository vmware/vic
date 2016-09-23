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

package trace

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestContextUnpack(t *testing.T) {
	Logger.Level = logrus.DebugLevel

	cnt := 100
	wg := &sync.WaitGroup{}
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func() {
			defer wg.Done()
			ctx := NewOperation(context.TODO(), "testmsg")

			// unpack an Operation via the context using it's Values fields
			c, err := FromContext(ctx)

			if !assert.NoError(t, err) || !assert.NotNil(t, c) {
				return
			}
			c.Infof("test info message %d", i)
		}()
	}
	wg.Wait()
}

// If we timeout a child, test a stack is printed of contexts
func TestNestedLogging(t *testing.T) {
	// create a buf to check the log
	buf := new(bytes.Buffer)
	Logger = logrus.StandardLogger()
	logrus.SetOutput(buf)

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
	lines := strings.Count(buf.String(), "\n")
	t.Log(buf.String())

	// Sample stack
	//
	// ERRO[0000] op=1: [ END ] [github.com/vmware/vic/pkg/trace.TestNestedLogging:60]: root error: context deadline exceeded
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging:60 root
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 0
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 1
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 2
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 3
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 4
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 5
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 6
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 7
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 8
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging.func1:70 level 9
	// ERRO[0000]      github.com/vmware/vic/pkg/trace.TestNestedLogging:80 Err

	// We arrive at 3 because we have the err line (line 0), then the root
	// (line 1), the then final "Err" line.
	if !assert.Equal(t, lines, levels+3) {
		return
	}
}

// Just checking behavior of the context package
func TestSanity(t *testing.T) {
	Logger.Level = logrus.InfoLevel
	levels := 10

	root, _ := context.WithDeadline(context.Background(), time.Time{})

	var ctxFunc func(parent context.Context, level int) context.Context

	ctxFunc = func(parent context.Context, level int) context.Context {
		if level == levels {
			return parent
		}

		child, _ := context.WithDeadline(parent, time.Now().Add(time.Hour))

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
