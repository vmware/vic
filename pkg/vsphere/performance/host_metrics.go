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

package performance

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
)

const (
	// cpuUsage measures the actively used CPU of the host, as a percentage of the total available CPU.
	cpuUsage = "cpu.usage.average"

	// memActive measures the sum of all active metrics for all powered-on VMs plus vSphere services on the host.
	memActive = "mem.active.average"

	// memConsumed measures the amount of machine memory used on the host, including vSphere services, VMkernel,
	// the service console and the total consumed memory metrics for all running VMs.
	memConsumed = "mem.consumed.average"

	// memTotalCapacity measures the total amount of memory reservation used by and available for powered-on VMs
	// and vSphere services on the host.
	memTotalCapacity = "mem.totalCapacity.average"

	// memOverhead measures the total of all overhead metrics for powered-on VMs, plus the overhead of running
	// vSphere services on the host.
	memOverhead = "mem.overhead.average"
)

// HostMemory stores an ESXi host's memory metrics.
type HostMemory struct {
	ActiveKB   int64
	ConsumedKB int64
	OverheadKB int64
	TotalKB    int64
}

// HostCPU stores an ESXi host's CPU metrics.
type HostCPU struct {
	UsagePercent float64
}

// HostMetricsInfo stores an ESXi host's memory and CPU metrics.
type HostMetricsInfo struct {
	Memory HostMemory
	CPU    HostCPU
}

// HostMetrics returns CPU and memory metrics for all ESXi hosts in the input session's cluster.
func HostMetrics(op trace.Operation, session *session.Session) (map[*object.HostSystem]*HostMetricsInfo, error) {
	if session == nil {
		return nil, fmt.Errorf("session not set")
	}

	// Gather hosts from the session cluster and then obtain their morefs.
	hosts, err := gatherHosts(op.Context, session)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain host morefs from session: %s", err)
	}
	morefToHost := make(map[types.ManagedObjectReference]*object.HostSystem)
	morefs := make([]types.ManagedObjectReference, len(hosts))
	for i, host := range hosts {
		moref := host.Reference()
		morefToHost[moref] = host
		morefs[i] = moref
	}

	// Query CPU and memory metrics for the morefs.
	spec := types.PerfQuerySpec{
		Format:     string(types.PerfFormatNormal),
		IntervalId: sampleInterval,
	}

	counters := []string{cpuUsage, memActive, memConsumed, memTotalCapacity, memOverhead}
	perfMgr := performance.NewManager(session.Vim25())
	sample, err := perfMgr.SampleByName(op.Context, spec, counters, morefs)
	if err != nil {
		errStr := "unable to get metric sample: %s"
		op.Errorf(errStr, err)
		return nil, fmt.Errorf(errStr, err)
	}

	results, err := perfMgr.ToMetricSeries(op.Context, sample)
	if err != nil {
		errStr := "unable to convert metric sample to metric series: %s"
		op.Errorf(errStr, err)
		return nil, fmt.Errorf(errStr, err)
	}

	metrics := assembleMetrics(op, morefToHost, results)
	return metrics, nil
}

// assembleMetrics processes the metric samples received from govmomi and returns a finalized metrics map
// keyed by the hosts.
func assembleMetrics(op trace.Operation, morefToHost map[types.ManagedObjectReference]*object.HostSystem,
	results []performance.EntityMetric) map[*object.HostSystem]*HostMetricsInfo {
	metrics := make(map[*object.HostSystem]*HostMetricsInfo)

	for _, host := range morefToHost {
		metrics[host] = &HostMetricsInfo{}
	}

	for i := range results {
		res := results[i]
		host, exists := morefToHost[res.Entity]
		if !exists {
			op.Warnf("moref %s does not exist in requested morefs, skipping", res.Entity.String())
			continue
		}

		// Process each value and assign it directly to the corresponding metric field
		// since there is only one sample.
		for _, v := range res.Value {

			// We don't need to collect non-aggregate (non-empty Instance) metrics.
			if v.Instance != "" {
				continue
			}

			if len(v.Value) == 0 {
				op.Warnf("metric %s for moref %s has no value, skipping", v.Name, res.Entity.String())
				continue
			}

			switch v.Name {
			case cpuUsage:
				// Convert percent units from 1/100th of a percent (100 = 1%) to a human-readable percentage.
				metrics[host].CPU.UsagePercent = float64(v.Value[0]) / 100.0
			case memActive:
				metrics[host].Memory.ActiveKB = v.Value[0]
			case memConsumed:
				metrics[host].Memory.ConsumedKB = v.Value[0]
			case memOverhead:
				metrics[host].Memory.OverheadKB = v.Value[0]
			case memTotalCapacity:
				// Total capacity is in MB, convert to KB so as to have all memory values in KB.
				metrics[host].Memory.TotalKB = v.Value[0] * 1024
			}
		}
	}

	return metrics
}

// gatherHosts gathers ESXi host(s) from the input session's cluster.
func gatherHosts(ctx context.Context, session *session.Session) ([]*object.HostSystem, error) {
	if session.Cluster == nil {
		return nil, fmt.Errorf("session cluster not set")
	}

	hosts, err := session.Cluster.Hosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain hosts from session cluster: %s", err)
	}
	if hosts == nil {
		return nil, fmt.Errorf("no hosts found from session cluster")
	}

	return hosts, nil
}
