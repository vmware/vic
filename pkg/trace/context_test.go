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
	"sync"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestContextDeadline(t *testing.T) {
	Logger.Level = logrus.DebugLevel

	cnt := 100
	wg := &sync.WaitGroup{}
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func() {
			ctx := NewOperation(context.TODO(), "testmsg")

			// unpack an Operation via the context using it's Values fields
			c := FromContext(ctx)

			assert.NotNil(t, c)
			c.Infof("foo %d", i)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestDeadlineLogging(t *testing.T) {
	Logger.Level = logrus.InfoLevel

	ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond)
	ctx = NewOperation(ctx, "testmsg")
	err := ctx.Err()
	if !assert.Error(t, err) {
		return
	}
}
