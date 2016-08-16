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
	h := newTrace(msg)

	t := Operation{
		tr:    *h,
		opNum: n,
	}

	// stash the value
	ctx = context.WithValue(ctx, OpTraceKey, t)

	Debugf(ctx, "[BEGIN] [%s]", h.frameName)
	return ctx
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
		Errorf(ctx, "[ END ] %s %s: error: %s", t.frameName, t.msg, err.Error())
		return err
	}

	Debugf(ctx, "[ END ] %s", t.msg)
	return nil
}

func header(ctx context.Context) string {
	t := FromContext(ctx)

	if Logger.Level >= logrus.DebugLevel {
		return fmt.Sprintf("op=%d (delta:%s)", t.opNum, time.Now().Sub(t.startTime))
	} else {
		return fmt.Sprintf("op=%d", t.opNum)
	}
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	Logger.Infof("%s: %s", header(ctx), fmt.Sprintf(format, args...))
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	Logger.Debugf("%s: %s", header(ctx), fmt.Sprintf(format, args...))
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	Logger.Errorf("%s: %s", header(ctx), fmt.Sprintf(format, args...))
}
