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

	"github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
)

const OpTraceKey = "traceKey"
const OpNumKey = "numKey"

// monotonic counter which inrements on Start()
var opCount uint64

type Operation struct {
	context.Context

	t     tr
	opNum uint64
}

// Add tracing info to the context.
func NewOperation(ctx context.Context, msg string) context.Context {
	// inc the counter
	n := atomic.AddUint64(&opCount, 1)

	// start the trace
	h := newTrace(msg)

	// We need to be able to identify this operation across API (and process)
	// boundaries.  So add the trace as a value to the embedded context.  We
	// stash the values individually in the context because we can't assign
	// the operation itself as a value to the embedded context (it's circular)
	ctx = context.WithValue(ctx, OpTraceKey, *h)
	ctx = context.WithValue(ctx, OpNumKey, n)

	o := Operation{
		Context: ctx,
		t:       *h,
		opNum:   n,
	}

	o.Debugf(o.t.beginHdr())
	return o
}

// Creates a header string to be printed.
func (o *Operation) header() string {
	if Logger.Level >= logrus.DebugLevel {
		return fmt.Sprintf("op=%d (delta:%s)", o.opNum, o.t.delta())
	} else {
		return fmt.Sprintf("op=%d", o.opNum)
	}
}

// Err returns a non-nil error value after Done is closed.  Err returns
// Canceled if the context was canceled or DeadlineExceeded if the
// context's deadline passed.  No other values for Err are defined.
// After Done is closed, successive calls to Err return the same value.
func (o Operation) Err() error {

	if err := o.Context.Err(); err != nil {
		o.Errorf("%s: %s error: %s", o.t.endHdr(), o.t.msg, err.Error())
		return err
	}

	return nil
}

func (o Operation) Infof(format string, args ...interface{}) {
	Logger.Infof("%s: %s", o.header(), fmt.Sprintf(format, args...))
}

func (o Operation) Debugf(format string, args ...interface{}) {
	Logger.Debugf("%s: %s", o.header(), fmt.Sprintf(format, args...))
}

func (o Operation) Errorf(format string, args ...interface{}) {
	Logger.Errorf("%s: %s", o.header(), fmt.Sprintf(format, args...))
}

// FromContext unpacks the values in the ctx to create an Operation
func FromContext(ctx context.Context) *Operation {

	o := &Operation{
		Context: ctx,
	}

	traceContext := ctx.Value(OpTraceKey)
	switch traceContext.(type) {
	case tr:
		o.t = traceContext.(tr)
	default:
		return nil
	}

	opNum := ctx.Value(OpNumKey)
	switch opNum.(type) {
	case uint64:
		o.opNum = opNum.(uint64)
	default:
		return nil
	}

	return o
}
