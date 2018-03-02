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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/performance"

	units "github.com/docker/go-units"
)

// MockMetricsProvider mocks the MetricsProvider interface.
type MockMetricsProvider struct{}

// GetMetricsForComputeResource not yet implemented.
func (m MockMetricsProvider) GetMetricsForComputeResource(op trace.Operation, cr *object.ComputeResource) (map[string]*performance.HostMetricsInfo, error) {
	return nil, nil
}

var (
	low = &performance.HostMetricsInfo{
		Memory: performance.HostMemory{
			TotalKB:    1 * units.GiB,
			ConsumedKB: 900 * units.MiB,
		},
		CPU: performance.HostCPU{
			UsagePercent: 0.25,
		},
	}
	// slightly higher CPU usage than medium
	lowMedium = &performance.HostMetricsInfo{
		Memory: performance.HostMemory{
			TotalKB:    16 * units.GiB,
			ConsumedKB: 9 * units.GiB,
		},
		CPU: performance.HostCPU{
			UsagePercent: 0.50,
		},
	}
	medium = &performance.HostMetricsInfo{
		Memory: performance.HostMemory{
			TotalKB:    16 * units.GiB,
			ConsumedKB: 9 * units.GiB,
		},
		CPU: performance.HostCPU{
			UsagePercent: 0.25,
		},
	}
	high = &performance.HostMetricsInfo{
		Memory: performance.HostMemory{
			TotalKB:    32 * units.GiB,
			ConsumedKB: 24 * units.GiB,
		},
		CPU: performance.HostCPU{
			UsagePercent: 0.3,
		},
	}

	lh = &object.HostSystem{
		Common: object.NewCommon(nil, types.ManagedObjectReference{
			Type:  "low_type",
			Value: "low_value",
		}),
	}
	lmh = &object.HostSystem{
		Common: object.NewCommon(nil, types.ManagedObjectReference{
			Type:  "lowmedium_type",
			Value: "lowmedium_value",
		}),
	}
	mh = &object.HostSystem{
		Common: object.NewCommon(nil, types.ManagedObjectReference{
			Type:  "medium_type",
			Value: "medium_value",
		}),
	}
	hh = &object.HostSystem{
		Common: object.NewCommon(nil, types.ManagedObjectReference{
			Type:  "high_type",
			Value: "high_value",
		}),
	}
)

func (m MockMetricsProvider) GetMetricsForHosts(op trace.Operation, hosts []*object.HostSystem) (map[string]*performance.HostMetricsInfo, error) {
	fakeHostMetrics := make(map[string]*performance.HostMetricsInfo)
	fakeHostMetrics[lh.Reference().String()] = low
	fakeHostMetrics[lmh.Reference().String()] = lowMedium
	fakeHostMetrics[mh.Reference().String()] = medium
	fakeHostMetrics[hh.Reference().String()] = high

	return fakeHostMetrics, nil
}

func TestRankHosts(t *testing.T) {
	op := trace.NewOperation(context.Background(), "TestRankHosts")

	model, server, _ := vpxModelSetup(op, t)
	defer func() {
		model.Remove()
		server.Close()
	}()

	m := MockMetricsProvider{}
	hm, err := m.GetMetricsForHosts(op, []*object.HostSystem{})
	assert.NoError(t, err)

	rhp := NewRankedHostPolicy(m)
	result := rhp.rankHosts(op, hm)

	for _, r := range result {
		op.Infof("%s: %f", r.HostReference, r.score)
	}

	expected := lh.Reference().String()
	actual := result[0].HostReference
	assert.NotEqual(t, expected, actual)

	expected = hh.Reference().String()
	actual = result[0].HostReference
	assert.Equal(t, expected, actual)
}
