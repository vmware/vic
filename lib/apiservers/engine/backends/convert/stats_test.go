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

package convert

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/ioutils"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/portlayer/metrics"
)

const (
	vcpuMhz        = 3300
	vcpuCount      = 1
	vchMhzTotal    = 3300
	memConsumed    = 1024 * 1024 * 500
	memProvisioned = 1024 * 1024 * 1024
)

func TestContainerConverter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	r, o := io.Pipe()
	defer o.Close()
	out := io.Writer(o)
	// Outstream modification (from Docker's code) so the stream is streamed with the
	// necessary headers that the CLI expects.  This is Docker's scheme.
	wf := ioutils.NewWriteFlusher(out)
	defer wf.Close()
	wf.Flush()
	out = io.Writer(wf)

	config := ContainerStatsConfig{
		VchMhz:      int64(vchMhzTotal),
		Ctx:         ctx,
		Cancel:      cancel,
		ContainerID: "1234",
		Out:         out,
		Stream:      true,
		Memory:      2048,
	}

	cStats := NewContainerStats(config)
	assert.NotNil(t, cStats)

	// this writer goes is provided to the PL
	writer := cStats.Listen()
	assert.NotNil(t, writer)

	w2 := cStats.Listen()
	assert.Nil(t, w2)

	initCPU := 1000
	vmBefore := vmMetrics(vcpuCount, initCPU)
	time.Sleep(1 * time.Millisecond)
	vmm := vmMetrics(vcpuCount, initCPU)

	// first metric sent, should return nil
	js, err := cStats.ToContainerStats(vmm)
	assert.NoError(t, err)
	assert.Nil(t, js)

	// send the same stat should return nil
	js, err = cStats.ToContainerStats(vmm)
	assert.Nil(t, err)
	assert.Nil(t, js)

	// send stat before the previous
	js, err = cStats.ToContainerStats(vmBefore)
	assert.NotNil(t, err)
	assert.Nil(t, js)

	secondCPU := 250
	// create a new metric
	vmmm := vmMetrics(vcpuCount, secondCPU)
	// sample will be 20 seconds apart..
	vmmm.SampleTime.Add(time.Second * 20)
	js, err = cStats.ToContainerStats(vmmm)
	assert.NoError(t, err)
	assert.NotZero(t, js.Read, js.PreRead)
	assert.Equal(t, uint64(vchMhzTotal*2), js.CPUStats.SystemUsage)
	assert.Equal(t, uint64(secondCPU+initCPU), js.CPUStats.CPUUsage.TotalUsage)
	assert.Equal(t, uint64(initCPU), js.PreCPUStats.CPUUsage.TotalUsage)
	assert.Equal(t, uint64(vchMhzTotal), js.PreCPUStats.SystemUsage)

	// this reading should show 250mhz of 3300mhz used -- 7.58%
	cpuPercent := fmt.Sprintf("%2.2f", calculateCPUPercentUnix(js.PreCPUStats.CPUUsage.TotalUsage, js.PreCPUStats.SystemUsage, js))
	assert.Equal(t, "7.58", cpuPercent)

	// reset listener, so reader/writer operates
	cStats.currentMetrics = nil
	cStats.dockerStats = &types.StatsJSON{}

	// simulate portLayer
	plEnc := json.NewEncoder(writer)
	err = plEnc.Encode(vmm)
	assert.NoError(t, err)
	err = plEnc.Encode(vmmm)
	assert.NoError(t, err)

	// simulate docker client
	docClient := json.NewDecoder(r)
	dstat := &types.StatsJSON{}
	err = docClient.Decode(dstat)
	assert.NoError(t, err)

	// ensure stop closes reader / writer
	cStats.Stop()
	_, err = cStats.reader.Read([]byte{0, 0, 0})
	assert.Error(t, err)

	config.Stream = false

	cStats = NewContainerStats(config)
	assert.NotNil(t, cStats)

	writer = cStats.Listen()

	// simulate portLayer
	plEnc = json.NewEncoder(writer)
	err = plEnc.Encode(vmm)
	assert.NoError(t, err)
	err = plEnc.Encode(vmmm)
	assert.NoError(t, err)

	// simulate docker client
	dstat = &types.StatsJSON{}
	err = docClient.Decode(dstat)
	assert.NoError(t, err)

}

func vmMetrics(count int, vcpuMhz int) *metrics.VMMetrics {
	vmm := &metrics.VMMetrics{}
	vmm.SampleTime = time.Now()
	vmm.CPU = cpuUsageMetrics(count, vcpuMhz)
	vmm.Memory = metrics.MemoryMetrics{
		Consumed:    int64(memConsumed),
		Provisioned: int64(memProvisioned),
	}
	return vmm
}

// cpuUsageMetrics will return a populated CPUMetrics struct
func cpuUsageMetrics(count int, cpuMhz int) metrics.CPUMetrics {
	vmCPUs := make([]metrics.CPUUsage, count, count)
	total := count * cpuMhz
	for i := range vmCPUs {
		vmCPUs[i] = metrics.CPUUsage{
			ID:       i,
			MhzUsage: int64(cpuMhz),
		}
	}

	return metrics.CPUMetrics{
		CPUs:  vmCPUs,
		Usage: calcVCPUUsage(total),
	}
}

// calcUsage is a helper function that will take the total provdied usage
// and convert to percentage of total vCPU usage
func calcVCPUUsage(total int) float32 {
	return float32(total) / (vcpuMhz * vcpuCount)
}

// calculateCPUPercentUnix is a copy from docker to test the percentage calculations
func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}
