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
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
)

const OpTraceKey = "traceKey"

// monotonic counter which inrements on Start()
var opCount uint64

type Operation struct {
	tr
	opNum uint64
}

// Add tracing info to the context.
func Start(ctx context.Context, msg string) context.Context {
	// inc the counter
	n := atomic.AddUint64(&opCount, 1)

	// start the trace
	h := Begin(msg)

	t := Operation{
		tr:    *h,
		opNum: n,
	}

	// stash the value
	return context.WithValue(ctx, OpTraceKey, t)
}

func FromContext(ctx context.Context) *Operation {
	traceContext := ctx.Value(OpTraceKey)

	switch traceContext.(type) {
	case Operation:
		t := traceContext.(Operation)
		return &t
	}

	return nil
}

func Done(ctx context.Context) error {
	t := FromContext(ctx)
	if err := ctx.Err(); err != nil {
		Errorf(ctx, "%s %s: error: %s", t.frameName, t.msg, err.Error())
		return err
	}

	End(&t.tr)
	return nil
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	t := FromContext(ctx)

	logrus.Infof("op=%d (delta:%s): %s", t.opNum, time.Now().Sub(t.startTime), fmt.Sprintf(format, args...))
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	t := FromContext(ctx)

	logrus.Debugf("op=%d (delta:%s): %s", t.opNum, time.Now().Sub(t.startTime), fmt.Sprintf(format, args...))
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	t := FromContext(ctx)

	logrus.Errorf("op=%d (delta:%s): %s", t.opNum, time.Now().Sub(t.startTime), fmt.Sprintf(format, args...))
}
