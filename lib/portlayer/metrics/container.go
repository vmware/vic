// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package metrics

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/pkg/pubsub"

	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
)

const (
	// number of samples per collection
	sampleSize = int32(2)
	// number of seconds between sample collection
	sampleInterval = int32(20)
)

// CPUUsage provides individual CPU metrics
type CPUUsage struct {
	// processor id (0,1,2)
	ID int
	// MhzUsage is the MhZ consumed by a specific processor
	MhzUsage int64
}

// CPUMetrics encapsulates available vm CPU metrics
type CPUMetrics struct {
	// CPUs are the individual CPU metrics
	CPUs []CPUUsage
	// Usage is the percentage of total vm CPU usage
	Usage float32
}

// MemoryMetrics encapsulates available vm memory metrics
type MemoryMetrics struct {
	// Consumed memory of vm in bytes
	Consumed int64
	// Active memory of vm in bytes
	Active int64
	// Provisioned memory of vm in bytes
	Provisioned int64
}

// NetworkUsage provides metrics for specific networks
type NetworkUsage struct {
	//TBD
}

// NetworkMetrics encapsulates available vm Network metrics
type NetworkMetrics struct {
	//TBD
}

// DiskUsage provides metrics for specific disks
type DiskUsage struct {
	//TBD
}

// NetworkMetrics encapsulates available vm Storage metrics
type StorageMetrics struct {
	//TBD
}

// VMMetrics encapsulates vm metrics available
type VMMetrics struct {
	CPU        CPUMetrics
	Memory     MemoryMetrics
	Network    NetworkMetrics
	Storage    StorageMetrics
	SampleTime time.Time
	// interval of collection in seconds
	Interval int32
}

// collectorVM is the VM metrics collector
type collectorVM struct {
	perfMgr *performance.Manager

	// subscribers to streaming
	mu   sync.RWMutex
	subs map[types.ManagedObjectReference]*subscriber
}

// subscriber is the receiver of the metrics
type subscriber struct {
	id        string
	ref       types.ManagedObjectReference
	publisher *pubsub.Publisher
}

func newVMCollector(session *session.Session) *collectorVM {
	s := make(map[types.ManagedObjectReference]*subscriber)
	cc := &collectorVM{
		subs:    s,
		perfMgr: performance.NewManager(session.Vim25()),
	}

	// kick off the subscription sampler
	go cc.sampler()

	return cc
}

// sampler will check for subscriptions and sample when needed
func (cc *collectorVM) sampler() {
	// TODO: replace this with service that only runs when there are
	// active subscribers
	d := time.Duration(int64(sampleInterval)) * time.Second
	for range time.Tick(d) {
		// if we have no subscribers skip
		if cc.SubscriberCount() == 0 {
			continue
		}
		// collect metrics for current subscribers
		cc.collect(cc.Subscribers(), false)
	}
}

// Sample returns a single metrics collection
func (cc *collectorVM) Sample(contain interface{}) (chan interface{}, error) {
	defer trace.End(trace.Begin(""))

	containerSub, err := containerSubscriber(contain)
	if err != nil {
		return nil, err
	}

	// create a publisher and subscribe
	containerSub.publisher = pubsub.NewPublisher(100*time.Millisecond, 0)
	ch := containerSub.publisher.Subscribe()

	// create a map with this single subscriber
	single := make(map[types.ManagedObjectReference]*subscriber)
	single[containerSub.ref] = containerSub

	// collect the metrics
	go cc.collect(single, true)

	return ch, nil
}

// collect will query vSphere for VM metrics and return to the subscribers
func (cc *collectorVM) collect(currentSubs map[types.ManagedObjectReference]*subscriber, single bool) {
	defer trace.End(trace.Begin(""))

	ctx := context.Background()
	var mos []types.ManagedObjectReference

	// metrics we are currently interested in monitoring
	names := []string{"cpu.usagemhz.average", "mem.active.average"}

	// create the spec
	spec := types.PerfQuerySpec{
		Format:     string(types.PerfFormatNormal),
		MaxSample:  sampleSize,
		IntervalId: sampleInterval,
	}

	// create a collection of object references
	for mo := range currentSubs {
		mos = append(mos, mo)
	}

	// if this is a single request close publisher when complete
	if single {
		defer currentSubs[mos[0]].publisher.Close()
	}

	// get the sample..
	sample, err := cc.perfMgr.SampleByName(ctx, spec, names, mos)
	if err != nil {
		log.Errorf("unable to get metric sample: %s", err)
		return
	}
	// convert to metrics
	result, err := cc.perfMgr.ToMetricSeries(ctx, sample)
	if err != nil {
		log.Errorf("unable to convert metric sample to metric series: %s", err)
		return
	}

	for i := range result {
		met := result[i]
		sub := currentSubs[met.Entity]
		for s := range met.SampleInfo {

			// assume we can publish, but if we hit an issue with a specific value
			// we will not publish the metric
			publish := true

			metric := VMMetrics{
				CPU: CPUMetrics{
					CPUs: []CPUUsage{},
				},
				Memory:     MemoryMetrics{},
				SampleTime: met.SampleInfo[s].Timestamp,
				Interval:   sampleInterval,
			}

			// the series will have values for each sample
			for _, v := range met.Value {
				switch v.Name {
				case "cpu.usagemhz.average":
					// vSphere returns individual cpu metrics and
					// the aggregate for all cpus.  We want to skip
					// the aggregate
					if v.Instance == "" {
						break
					}
					// we aren't on the aggregate so convert to int
					iid, err := strconv.Atoi(v.Instance)
					if err != nil {
						// I don't expect this to ever happen, but if it does log and don't publish
						log.Errorf("metrics failed to convert container(%s) CPU id to an int - value(%#v): %s", sub.id, v, err)
						publish = false
						break
					}
					// specific vCPU metric
					cpu := CPUUsage{
						ID:       iid,
						MhzUsage: v.Value[s],
					}
					metric.CPU.CPUs = append(metric.CPU.CPUs, cpu)
				case "mem.active.average":
					metric.Memory.Active = v.Value[s]
				}

			}

			if publish {
				sub.publisher.Publish(metric)
			}
		}
	}
}

// Subscribe subscribes to a publisher that publishes metrics
func (cc *collectorVM) Subscribe(contain interface{}) (chan interface{}, error) {
	defer trace.End(trace.Begin(""))
	containerSub, err := containerSubscriber(contain)
	if err != nil {
		return nil, err
	}
	cc.mu.Lock()
	defer cc.mu.Unlock()

	sub, exists := cc.subs[containerSub.ref]
	if !exists {
		// create a new publisher with 100ms timeout and no buffer
		containerSub.publisher = pubsub.NewPublisher(100*time.Millisecond, 0)
		cc.subs[containerSub.ref] = containerSub
		sub = containerSub
	}
	// subscribe to this publisher -- a publisher can have multiple subscribers
	ch := sub.publisher.Subscribe()
	return ch, nil
}

// Unsubscribe unsubscribes from the container metrics publisher
func (cc *collectorVM) Unsubscribe(contain interface{}, ch chan interface{}) {
	defer trace.End(trace.Begin(""))
	containerSub, err := containerSubscriber(contain)
	if err != nil {
		return
	}
	cc.mu.Lock()
	defer cc.mu.Unlock()
	sub, exists := cc.subs[containerSub.ref]
	if exists {
		sub.publisher.Evict(ch)
		if sub.publisher.Len() == 0 {
			delete(cc.subs, sub.ref)
		}
	}
}

// Subscribers returns the current subscribers
func (cc *collectorVM) Subscribers() map[types.ManagedObjectReference]*subscriber {
	current := make(map[types.ManagedObjectReference]*subscriber)
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for ref := range cc.subs {
		current[ref] = cc.subs[ref]
	}
	return current
}

// SubscribeCount will return the current number of subscribers
func (cc *collectorVM) SubscriberCount() int {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return len(cc.subs)
}

// containerSubscriber is a helper func to convert the interface to a subscriber
func containerSubscriber(contain interface{}) (*subscriber, error) {
	container, ok := contain.(*exec.Container)
	if !ok {
		return nil, fmt.Errorf("invalid type provided for container stats request - got: %#v", contain)
	}
	info := container.Info()
	moRef := info.VMReference()

	// ensure we have a valid moRef..we won't worry about inspecting the details
	if moRef.String() == "" {
		return nil, fmt.Errorf("no vm associated with provided container(%s) stats request", container.ExecConfig.ID)
	}

	sub := &subscriber{
		id:  info.ExecConfig.ID,
		ref: moRef,
	}

	return sub, nil
}
