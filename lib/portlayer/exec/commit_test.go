// Copyright 2018 VMware, Inc. All Rights Reserved.
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

package exec

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

func TestBatchBlockOnFuncSerialize(t *testing.T) {
	ctx := context.Background()

	var totalRequest int
	var batchCount int
	var reqInterval time.Duration
	var waitTime time.Duration

	countBatch := func(op trace.Operation, waitTime time.Duration) error {
		time.Sleep(waitTime)
		batchCount++
		return nil
	}

	// serial test: only one request
	batchCount = 0
	testFetchOne(ctx, t, func(op trace.Operation) error {
		return countBatch(op, 0)
	})
	assert.Equal(t, 1, batchCount)

	// serial test: batch size 1, 5 requests
	// batch size 1 means no batching
	// total # of batches = total # of requests
	batchCount = 0
	totalRequest = 5
	testMultipleBatch(ctx, t, 1, totalRequest, 0, func(op trace.Operation) error {
		return countBatch(op, 0)
	}, nil)
	assert.Equal(t, totalRequest, batchCount)

	// serial test: batch size 10, 20 requests
	// requests come in slower than operation processing time, making all requests serialized
	// total # of batches = total # of requests
	batchCount = 0
	totalRequest = 20
	reqInterval = 15 * time.Millisecond
	waitTime = 10 * time.Millisecond
	testMultipleBatch(ctx, t, 10, totalRequest, reqInterval, func(op trace.Operation) error {
		return countBatch(op, waitTime)
	}, nil)
	assert.Equal(t, totalRequest, batchCount)
}

func TestBatchBlockOnFuncConcurrent(t *testing.T) {
	ctx := context.Background()

	var batchCount int
	var totalRequest int

	operation := func(op trace.Operation) error {
		time.Sleep(10 * time.Millisecond)
		batchCount++
		return nil
	}

	// batch size 10, 5 requests
	// requests come in tight loop. All 5 should be processed in one single batch
	// total time spent should be within: 1 * (operation wait time)
	batchCount = 0
	totalRequest = 5
	testMultipleBatch(ctx, t, 10, totalRequest, 0, operation, nil)
	assert.True(t, batchCount > 1 && batchCount < totalRequest)

	// batch size 100, 50 requests
	// the 50 requests should be processed in one batch
	// total time spent should be within: 1 * (operation wait time)
	batchCount = 0
	totalRequest = 50
	testMultipleBatch(ctx, t, 100, totalRequest, 0, operation, nil)
	assert.True(t, batchCount > 1 && batchCount < totalRequest)

	// stress test: batch size 100, 200 requests
	// the 200 requests should be processed within 2 batches
	// total time spent <= 2 * (operation wait time)
	batchCount = 0
	totalRequest = 200
	testMultipleBatch(ctx, t, 100, totalRequest, 0, operation, nil)
	assert.True(t, batchCount > 1 && batchCount < totalRequest)

	// stress test: batch size 100, 500 requests
	// the 500 requests should be processed within 5 batches
	// total time spent <= 5 * (operation wait time)
	batchCount = 0
	totalRequest = 500
	testMultipleBatch(ctx, t, 100, totalRequest, 0, operation, nil)
	assert.True(t, batchCount > 1 && batchCount < totalRequest)
}

func TestBatchBlockOnFuncResultPropagate(t *testing.T) {
	ctx := context.Background()

	err := errors.New("test")
	operation := func(op trace.Operation) error {
		time.Sleep(10 * time.Millisecond)
		return err
	}

	testMultipleBatch(ctx, t, 10, 5, 0, operation, err)
}

func testFetchOne(ctx context.Context, t *testing.T, operation func(op trace.Operation) error) {
	batch := make(chan chan error, 5) // batch size 5

	// fire background reader
	go batchBlockOnFunc(ctx, batch, operation)

	// send only 1 request
	sendRequest(t, batch, nil, nil)
	close(batch)
}

func testMultipleBatch(ctx context.Context, t *testing.T, batchSize int, totalRequest int, interval time.Duration, operation func(op trace.Operation) error, expected error) {
	batch := make(chan chan error, batchSize)

	// fire background request reader
	go batchBlockOnFunc(ctx, batch, operation)

	// send requests concurrently with a time interval between requests
	done := sendMultiRequests(t, totalRequest, batch, interval, expected)

	// wait until all requests are processed, close batch channel and quit background receiver
	quitBatchUntilDone(t, done, batch)
}

func sendMultiRequests(t *testing.T, totalRequest int, batch chan chan error, interval time.Duration, expected error) []chan bool {
	done := make([]chan bool, totalRequest)

	for i := 0; i < totalRequest; i++ {
		done[i] = make(chan bool, 1)
		go sendRequest(t, batch, done[i], expected)
		time.Sleep(interval)
	}

	return done
}

func sendRequest(t *testing.T, batch chan chan error, done chan bool, expected error) {
	req := make(chan error)
	batch <- req
	err := <-req
	assert.Equal(t, expected, err)
	if done != nil {
		done <- true
	}
}

func quitBatchUntilDone(t *testing.T, done []chan bool, batch chan chan error) {
	for _, c := range done {
		select {
		case _ = <-c:
			close(c)
			continue
		case <-time.After(30 * time.Second):
			t.Fail()
		}
	}
	close(batch)
}
