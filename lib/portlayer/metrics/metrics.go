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
	"sync"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var (
	Supervisor *super

	initializer struct {
		err  error
		once sync.Once
	}
)

// super manages the lifecycle and access to the
// available metrics collectors
type super struct {
	vms *collectorVM
}

type Collector interface {
	// Subscribe to a stream of this collectors metrics
	Subscribe(interface{}) (chan interface{}, error)
	// Unsubscribe to a stream of this collectors metrics
	Unsubscribe(interface{}, chan interface{})
	// SubscriberCount returns the number of subscribers for the collector
	SubscriberCount() int
	// Sample will return metrics without a subscription
	Sample(interface{}) (chan interface{}, error)
}

func Init(ctx context.Context, session *session.Session) error {
	defer trace.End(trace.Begin(""))
	initializer.once.Do(func() {
		var err error
		defer func() {
			if err != nil {
				initializer.err = err
			}
		}()
		Supervisor = newSupervisor(session)

	})
	return initializer.err

}

func newSupervisor(session *session.Session) *super {
	defer trace.End(trace.Begin(""))
	// create the vm metric collector
	v := newVMCollector(session)
	return &super{
		vms: v,
	}
}

// VMCollector will return the vm metrics collector
func (s *super) VMCollector() Collector {
	return s.vms
}
