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
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/performance"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test"

	units "github.com/docker/go-units"
)

// MockMetricsProvider mocks the MetricsProvider interface.
type MockMetricsProvider struct{}

// GetMetricsForComputeResource not yet implemented.
func (m *MockMetricsProvider) GetMetricsForComputeResource(op trace.Operation, cr *object.ComputeResource) (map[*object.HostSystem]*performance.HostMetricsInfo, error) {
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
		Common: object.Common{
			InventoryPath: "low_rank",
		},
	}
	lmh = &object.HostSystem{
		Common: object.Common{
			InventoryPath: "lowmedium_rank",
		},
	}
	mh = &object.HostSystem{
		Common: object.Common{
			InventoryPath: "medium_rank",
		},
	}
	hh = &object.HostSystem{
		Common: object.Common{
			InventoryPath: "high_rank",
		},
	}
)

func (m *MockMetricsProvider) GetMetricsForHosts(op trace.Operation, hosts []*object.HostSystem) (map[*object.HostSystem]*performance.HostMetricsInfo, error) {
	fakeHostMetrics := make(map[*object.HostSystem]*performance.HostMetricsInfo)
	fakeHostMetrics[lh] = low
	fakeHostMetrics[lmh] = lowMedium
	fakeHostMetrics[mh] = medium
	fakeHostMetrics[hh] = high

	return fakeHostMetrics, nil
}

// vpxModelSetup creates a VPX model, starts its server and populates the session. The caller must
// clean up the model and the server once it is done using them.
func vpxModelSetup(ctx context.Context, t *testing.T) (*simulator.Model, *simulator.Server, *session.Session) {
	model := simulator.VPX()
	if err := model.Create(); err != nil {
		t.Fatal(err)
	}

	server := model.Service.NewServer()
	sess, err := test.SessionWithVPX(ctx, server.URL.String())
	if err != nil {
		t.Fatal(err)
	}

	return model, server, sess
}

func TestRecommendHost(t *testing.T) {
	t.Skip("Not yet implemented")
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
	result := rankHosts(op, hm)

	for _, r := range result {
		op.Infof("%s: %f", r.InventoryPath, r.score)
	}

	assert.NotEqual(t, lh, result[0], "Expected %s, got %s", lh.InventoryPath, result[0].InventoryPath)
	assert.Equal(t, hh.InventoryPath, result[0].InventoryPath, "Expected %s, got %s", hh.InventoryPath, result[0].InventoryPath)
}
