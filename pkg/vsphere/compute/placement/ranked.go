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

package placement

import (
	"fmt"
	"sort"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/performance"
	"github.com/vmware/vic/pkg/vsphere/session"
)

type rankedHost struct {
	HostReference string
	*performance.HostMetricsInfo
	score float64
}

type rankedHosts []rankedHost

func (r rankedHosts) Len() int           { return len(r) }
func (r rankedHosts) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r rankedHosts) Less(i, j int) bool { return r[i].score > r[j].score }

// WeightConfiguration holds user-provided weights for different host metrics. These weights are
// used to determine a host ranking.
type WeightConfiguration struct {
	memUnconsumedWeight float64
	memInactiveWeight   float64
}

// RankedHostPolicy uses data from a MetricsProvider to decide on which host to power-on a VM.
type RankedHostPolicy struct {
	source performance.MetricsProvider
	config WeightConfiguration
}

// NewRankedHostPolicy returns a RandomHostPolicy instance using the supplied MetricsProvider with
// the default weighting configuration.
func NewRankedHostPolicy(s performance.MetricsProvider) *RankedHostPolicy {
	return NewRankedHostPolicyWithConfig(s, WeightConfiguration{
		memInactiveWeight:   MemDefaultInactiveWeight,
		memUnconsumedWeight: MemDefaultUnconsumedWeight,
	})
}

// NewRankedHostPolicyWithConfig returns a RandomHostPolicy instance using the supplied MetricsProvider and
// WeightConfiguration.
func NewRankedHostPolicyWithConfig(s performance.MetricsProvider, wc WeightConfiguration) *RankedHostPolicy {
	return &RankedHostPolicy{
		source: s,
		config: wc,
	}
}

// CheckHost returns true if the host has adequate capacity to power on the VM, false otherwise.
func (r *RankedHostPolicy) CheckHost(op trace.Operation, host *object.HostSystem) bool {
	// TODO(jzt): return false until we have host checking logic decided
	return false
}

// RecommendHost recommends an ideal host on which to place a newly created VM.
func (r *RankedHostPolicy) RecommendHost(op trace.Operation, sess *session.Session, hosts []*object.HostSystem) ([]*object.HostSystem, error) {
	var (
		err error
		hm  map[string]*performance.HostMetricsInfo
	)

	if len(hosts) == 0 {
		op.Debugf("no hosts specified - gathering metrics on cluster")
		hm, err = r.source.GetMetricsForComputeResource(op, sess.Cluster)
	} else {
		hm, err = r.source.GetMetricsForHosts(op, hosts)
	}
	if err != nil {
		return nil, err
	}

	if len(hm) == 0 {
		return nil, fmt.Errorf("no candidate hosts to rank")
	}

	ranked := r.rankHosts(op, hm)
	result := make([]*object.HostSystem, 0, len(ranked))
	for _, h := range ranked {
		ref := types.ManagedObjectReference{}
		if ok := ref.FromString(h.HostReference); !ok {
			return nil, fmt.Errorf("could not restore serialized managed object reference: %s", h.HostReference)
		}

		result = append(result, object.NewHostSystem(sess.Vim25(), ref))
	}

	return result, nil
}

func (r *RankedHostPolicy) rankHosts(op trace.Operation, hm map[string]*performance.HostMetricsInfo) []rankedHost {
	ranking := []rankedHost{}
	for h, m := range hm {
		rh := rankedHost{
			HostReference:   h,
			HostMetricsInfo: m,
			score:           r.rankMemory(m) * (1 - m.CPU.UsagePercent),
		}
		ranking = append(ranking, rh)
	}
	sort.Sort(rankedHosts(ranking))
	return ranking
}

func (r *RankedHostPolicy) rankMemory(hm *performance.HostMetricsInfo) float64 {
	free := float64(hm.Memory.TotalKB-hm.Memory.ConsumedKB) / 1024.0
	return free * r.config.memUnconsumedWeight
}
