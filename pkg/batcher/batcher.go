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

package batcher

import (
	"context"
	"sync"
	"time"
)

// Batcher is a batching mechanism that can be configured with the following:
//  - assessor: callback to determine if a queued item is admissible to a batch, with control of when the rejection is returned if rejected
//  - processor: callback to actually process the batch
//  - latency: the maximum latency a queued item can tolerate (for more efficient grouping when possible)
//  - serialization: allows serialization of batches that share the same group ID (will wait for the current group to complete)
type Batcher interface {
	// QueueSync blocks until the member is processed. There is no guaranteed ordering within the batch, even if the caller
	// enforces ordering on the QueueX calls.
	QueueSync(ctx context.Context, groupID string, latency time.Duration, data interface{}) interface{}
	QueueAsync(ctx context.Context, groupID string, latency time.Duration, data interface{}, returnHandler func(interface{}))
	Start(ctx context.Context)
	Stop()
}

// Assessment is the outcome reported by the admission control assessor and is used internally by the batcher to control when errors are
// returned.
type Assessment int

const (
	// Accept means the member was accepted into the batch
	Accept = iota
	// RejectImmediate means the member was rejected from the batch and the payload should be returned immediately
	RejectImmediate
	// RejectWaitIssue means the member was rejected from the batch and the payload should be returned once the
	// batch is issued
	RejectWaitIssue
	// RejectWaitComplete means the member was rejected from the batch and the payload should be returned once the
	// batch completes
	RejectWaitComplete
)

// Assessor is an opaque function that determines whether a batch member will be accepted into the batch group.
// Returns nil if the member is accepted, otherwise the return is expected to be an error or similar and is
// propagated back to the caller of Queue.
// Candidate is the item to be batched
// Members are the existing items within a group
// If the assessment is accept then the interface return is passed back into the next call of assessor for that batch
// via the state field to allow for passing of working state
// If the assessment is a type of rejection then the interface returned is passed back to the caller of the Queue function
// with timing as indicated by the Assessment
type Assessor func(ctx context.Context, groupID string, candidate interface{}, members []interface{}, state interface{}) (Assessment, interface{})

// Contexts allows a caller to retrieve a specific context assocated with a member if available. This will be
// the context provided when calling Queue...
type Contexts func(member interface{}) context.Context

// Processor is an opaque function that actually performs the batched work. The ctxs function is provided to
// state is the return from the last call to the assessor function if any. This is provided so that if processing
// and assessment need to share logic then we can avoid repetition in the call to processor.
// ctxs is provided to allow this function to retrieve the context directly associated with the member data if
// any.
type Processor func(groupID string, members []interface{}, state interface{}, ctxs Contexts) interface{}

const (
	// Maximum number of non-blocking pending queue operations that have not yet been through admission control.
	pendingQueue = 100
	// Duration a group processor will wait for new operations before terminating
	gcDelay = 5 * time.Second
)

type batchMember struct {
	groupID string
	latency time.Duration
	ctx     context.Context
	ret     chan interface{}
	data    interface{}
}

type batchGroup struct {
	// id is the group ID
	id string
	// the channels through which to pass the final result
	ret []chan interface{}
	// the members of the batch and an associated context that may be nil
	data map[interface{}]context.Context
	// members is the key set of the data map, saved to avoid recreating it multiple times
	members []interface{}
	// working state passed between calls to the assessor function for the same group
	state interface{}
	// queue operations to be rejected just prior to the batch being issued
	rejectOnIssue map[chan interface{}]interface{}
	// queue operations to be rejected after the batch has completed
	rejectOnCompletion map[chan interface{}]interface{}

	// the batcher this group was created by
	batcher *batcher
	// queue for dispatch to this specific group instance
	queue chan *batchMember
}

type batcher struct {
	mu sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	queue chan *batchMember

	assessor  Assessor
	processor Processor

	batchSize    int
	maxLatency   time.Duration
	groupGcDelay time.Duration

	completions     chan *batchGroup
	serializeGroups bool
}

// NewBatcher provides a batching mechanism that can be configured with the following:
//  - assessor: optional callback to determine if a queued item is admissible to a batch
//  - processor: callback to actually process the batch
//  - latency: the maximum latency a queued item can tolerate (for more efficient grouping when possible)
//  - serialization: allows serialization of batches that share the same group ID (will wait for the current group to complete)
func NewBatcher(ctx context.Context, assessor Assessor, processor Processor, serializeGroups bool) Batcher {
	b := &batcher{
		queue:           make(chan *batchMember, pendingQueue),
		completions:     make(chan *batchGroup, pendingQueue), // no reasoning behind using pending queue except it's available
		assessor:        assessor,
		processor:       processor,
		serializeGroups: serializeGroups,
		groupGcDelay:    gcDelay,
	}

	return b
}

func (b *batcher) Start(ctx context.Context) {
	b.mu.Lock()

	if b.cancel == nil {
		b.ctx, b.cancel = context.WithCancel(ctx)
		go b.global(b.ctx)
	}

	b.mu.Unlock()
}

func (b *batcher) Stop() {
	b.mu.Lock()
	b.cancel()
	b.cancel = nil
	b.mu.Unlock()
}

func (b *batcher) QueueSync(ctx context.Context, groupID string, latency time.Duration, data interface{}) interface{} {
	member := &batchMember{
		groupID: groupID,
		latency: latency,
		ctx:     ctx,
		ret:     make(chan interface{}, 1),
		data:    data,
	}

	b.queue <- member

	return <-member.ret
}

func (b *batcher) QueueAsync(ctx context.Context, groupID string, latency time.Duration, data interface{}, returnHandler func(interface{})) {
	member := &batchMember{
		groupID: groupID,
		latency: latency,
		ctx:     ctx,
		ret:     make(chan interface{}, 1),
		data:    data,
	}

	b.queue <- member

	// this will run the error processing in the background
	go func() {
		val := <-member.ret
		if returnHandler != nil {
			returnHandler(val)
		}
	}()
}

// global is responsible for routing incoming queued operations to the group processors, and for lifecycle management of
// those group processors
func (b *batcher) global(ctx context.Context) {
	groups := make(map[string]*batchGroup)

	for {
		// block and wait for first request
		select {
		case req, ok := <-b.queue:
			if !ok {
				return // channel closed, quit
			}

			if req.ret == nil {
				// precaution vs misuse - all batch members should supply a channel even if they discard the result
				continue
			}

			// create group and start processor if not already present
			group, exits := groups[req.groupID]
			if !exits {
				group = &batchGroup{
					batcher: b,
					id:      req.groupID,
					queue:   make(chan *batchMember, 100),
				}
				groups[req.groupID] = group

				go group.start(ctx)
			}

			group.queue <- req

		case group, ok := <-b.completions:
			// after we've processed the exit notification either the group processor and queue are completely gone, or
			// it's been fully restored and will issue another notification on subsequent exit. This is necessary so the
			// only lifecycle operation request dequeuing case needs to handle is a full create/start.
			if !ok {
				return // channel closed, quit
			}

			if len(group.queue) == 0 {
				delete(groups, group.id)
			} else {
				// restart group processor to process work in queue
				go group.start(ctx)
			}

		case <-ctx.Done():
			// when parent context is cancelled, quit
			return
		}
	}
}

// add assesses whether the member could be added to the current batch and either adds it or handles the rejection
// as requested by the assessor function for the batcher.
// Returns true if accepted into the batch, false otherwise
func (g *batchGroup) add(req *batchMember) bool {
	if g.batcher.assessor != nil {
		assessment, result := g.batcher.assessor(req.ctx, g.id, req.data, g.members, g.state)
		if assessment != Accept {
			switch assessment {
			case RejectImmediate:
				req.ret <- result
			case RejectWaitIssue:
				g.rejectOnIssue[req.ret] = result
			case RejectWaitComplete:
				g.rejectOnCompletion[req.ret] = result
			}

			return false
		}

		g.state = result
	}

	g.data[req.data] = req.ctx
	g.ret = append(g.ret, req.ret)
	g.members = append(g.members, req.data)

	return true
}

// start is the per-group processor loop. It will continue processing batches until the GC delay expires, and will then
// exit.
func (g *batchGroup) start(ctx context.Context) {
	// TODO: add max latency bounds - for now we add no latency which is compatible with any expressed tolerance.

	for {
		// we reinitialize these structures each loop rather than attempt to clean the previous batch with this groupID
		g.data = make(map[interface{}]context.Context)
		g.ret = make([]chan interface{}, 0, 1)
		g.members = make([]interface{}, 0, 1)
		g.rejectOnIssue = make(map[chan interface{}]interface{})
		g.rejectOnCompletion = make(map[chan interface{}]interface{})

		// ran into hangs when trying to use stop/drain/reset as per timer.Reset doc so just new timer every time
		gc := time.NewTimer(g.batcher.groupGcDelay)

		// block and wait for first request
		select {
		case req, ok := <-g.queue:
			if !ok {
				return // channel closed, quit
			}
			if req.ret == nil {
				// precaution vs misuse - all batch members should supply a channel even if they discard the result
				continue
			}

			if !g.add(req) {
				// this wasn't added to a group so we still need to block
				continue
			}

		case <-gc.C:
			// no traffic during the GC delay, so exit and let the router know. If any traffic arrives before the
			// group is fully GC'd the router will restart this routine.
			g.batcher.completions <- g
			return

		case <-ctx.Done():
			// when parent context is cancelled, quit
			return
		}

		// fetch additional requests
		for len(g.queue) > 0 {
			req := <-g.queue
			if req.ret == nil {
				continue
			}

			g.add(req)
		}

		// ensure that we've done this rejection step just prior to dispatching the group. If not done here then
		// go routine scheduling may mean this happens after the group routine has exited and I don't think that's
		// optimal.
		for k, v := range g.rejectOnIssue {
			k <- v
		}

		// process requests
		if g.batcher.serializeGroups {
			g.dispatch()
		} else {
			// create a copy of the group so that state relevant to dispatch remains available after this routine
			// reinitializes them
			g2 := *g
			go g2.dispatch()
		}
	}
}

// dispatch is responsible for calling the Batcher Processor function, rejection of any members with an assessment of
// RejectOnCompletion, and fan out of the result of the Processor.
func (g *batchGroup) dispatch() {
	retVal := g.batcher.processor(g.id, g.members, g.state, func(key interface{}) context.Context {
		return g.data[key]
	})

	for k, v := range g.rejectOnCompletion {
		k <- v
	}

	// signal batched operations and throw back result
	for i := range g.ret {
		g.ret[i] <- retVal
		close(g.ret[i])
	}
}
