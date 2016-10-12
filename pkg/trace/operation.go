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
	"os"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
)

const OpTraceKey = "traceKey"

var opIDPrefix = os.Getpid()

// monotonic counter which inrements on Start()
var opCount uint64

type Operation struct {
	context.Context
	operation
}

type operation struct {
	t  []Message
	id string
}

func newOperation(ctx context.Context, id string, skip int, msg string) Operation {
	op := operation{

		// Can be used to trace based on this number which is unique per chain
		// of operations
		id: id,

		// Start the trace.
		t: []Message{*newTrace(msg, skip)},
	}

	// We need to be able to identify this operation across API (and process)
	// boundaries.  So add the trace as a value to the embedded context.  We
	// stash the values individually in the context because we can't assign
	// the operation itself as a value to the embedded context (it's circular)
	ctx = context.WithValue(ctx, OpTraceKey, op)

	o := Operation{
		Context:   ctx,
		operation: op,
	}

	o.Debugf(o.t[0].beginHdr())
	return o
}

// Creates a header string to be printed.
func (o *Operation) header() string {
	if Logger.Level >= logrus.DebugLevel {
		return fmt.Sprintf("op=%s (delta:%s)", o.id, o.t[0].delta())
	}
	return fmt.Sprintf("op=%s", o.id)
}

// Err returns a non-nil error value after Done is closed.  Err returns
// Canceled if the context was canceled or DeadlineExceeded if the
// context's deadline passed.  No other values for Err are defined.
// After Done is closed, successive calls to Err return the same value.
func (o Operation) Err() error {

	// Walk up the contexts from which this context was created and get their errors
	if err := o.Context.Err(); err != nil {
		// Print the error
		o.Errorf("%s: %s error: %s", o.t[0].endHdr(), o.t[0].msg, err)

		// Walk the stack and end with where this was called from
		for _, t := range append(o.t, *newTrace("Err", 2)) {
			Logger.Errorf("\t%s:%d %s", t.funcName, t.lineNo, t.msg)
		}

		return err
	}

	return nil
}

func (o *Operation) Infof(format string, args ...interface{}) {
	Logger.Infof("%s: %s", o.header(), fmt.Sprintf(format, args...))
}

func (o *Operation) Debugf(format string, args ...interface{}) {
	Logger.Debugf("%s: %s", o.header(), fmt.Sprintf(format, args...))
}

func (o *Operation) Errorf(format string, args ...interface{}) {
	Logger.Errorf("%s: %s", o.header(), fmt.Sprintf(format, args...))
}

func (o *Operation) newChild(ctx context.Context, msg string) Operation {
	child := newOperation(ctx, o.id, 4, msg)
	t := child.t[0]
	child.t = append(o.t, t)
	return child
}

func opID(opNum uint64) string {
	return fmt.Sprintf("%d.%d", opIDPrefix, opNum)
}

// Add tracing info to the context.
func NewOperation(ctx context.Context, msg string) Operation {
	return newOperation(ctx, opID(atomic.AddUint64(&opCount, 1)), 3, msg)
}

// WithTimeout
func WithTimeout(parent *Operation, timeout time.Duration, msg string) (Operation, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(parent.Context, timeout)
	op := parent.newChild(ctx, msg)

	return op, cancelFunc
}

// WithDeadline
func WithDeadline(parent *Operation, expiration time.Time, msg string) (Operation, context.CancelFunc) {
	ctx, cancelFunc := context.WithDeadline(parent.Context, expiration)
	op := parent.newChild(ctx, msg)

	return op, cancelFunc
}

// FromContext unpacks the values in the ctx to create an Operation
func FromContext(ctx context.Context) (Operation, error) {

	o := Operation{
		Context: ctx,
	}

	op := ctx.Value(OpTraceKey)
	switch val := op.(type) {
	case operation:
		o.operation = val
	default:
		return Operation{}, fmt.Errorf("not an Operation")
	}

	return o, nil
}
