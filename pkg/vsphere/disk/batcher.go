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

package disk

import (
	"context"
	"time"

	"github.com/vmware/vic/pkg/trace"
)

const (
	// BatchLatency is the duration to wait for other work before acting on lazy operations
	BatchLatency = 10 * time.Second
	// MaxBatchSize is the number of items to batch and is set sufficient to completely cycle attached disks
	MaxBatchSize = MaxAttachedDisks * 2
)

type batchMember struct {
	op   trace.Operation
	err  chan error
	data interface{}
}

// lazyDeviceChange adds a lazy deferral mechanism for device change operations (specifically disk at this time).
// This is due to the fact that reconfigure operations are unintentionally serializing parallel operations and
// causing performance impacts (concurrent volume create as a primary example which impacts concurrent container start if
// there are anonymous volumes - worse if there are multiple volumes per container)
func lazyDeviceChange(ctx context.Context, batch chan batchMember, operation func(operation trace.Operation, data []interface{}) error) {
	op := trace.FromContext(ctx, "Lazy batching of disk operations")

	for {
		data := make([]interface{}, 0, MaxBatchSize)  // batching queue for arguments
		errors := make([]chan error, 0, MaxBatchSize) // batching queue for error returns

		// block and wait for first request
		select {
		case req, ok := <-batch:
			if !ok {
				return // channel closed, quit
			}
			if req.err == nil {
				continue
			}
			data = append(data, req.data)
			errors = append(errors, req.err)
			req.op.Debugf("Dispatching queued operation")
		case <-op.Done(): // when parent context is cancelled, quit
			return
		}

		// fetch batched requests
		// TODO: I want to add some optional latency in here so that the attach/detach pair from
		// pull make use of it, but for now it's purely opportunistic for non-serial operations.
		for len(batch) > 0 {
			req := <-batch
			if req.err != nil {
				data = append(data, req.data)
				errors = append(errors, req.err)
				req.op.Debugf("Dispatching queued operation")
			}
		}

		// process requests
		err := operation(op, data)

		// signal batched operations and throw back result
		for _, member := range errors {
			member <- err
			close(member)
		}
	}
}
