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
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/api/types"

	"github.com/vmware/vic/lib/portlayer/metrics"
)

// ContainerStats encapsulates the conversion of VMMetrics to
// docker specific metrics
type ContainerStats struct {
	config *ContainerStatsConfig

	totalVCHMhz uint64
	dblVCHMhz   uint64
	preTotalMhz uint64

	preDockerStat *types.StatsJSON
	curDockerStat *types.StatsJSON
	currentMetric *metrics.VMMetrics

	mu        sync.Mutex
	reader    *io.PipeReader
	writer    *io.PipeWriter
	listening bool
}

type ContainerStatsConfig struct {
	Ctx            context.Context
	Cancel         context.CancelFunc
	Out            io.Writer
	ContainerID    string
	ContainerState *types.ContainerState
	Memory         int64
	Stream         bool
	VchMhz         int64
}

type InvalidOrderError struct {
	current  time.Time
	previous time.Time
}

func (iso InvalidOrderError) Error() string {
	return fmt.Sprintf("The current sample time (%s) is before the previous time (%s)", iso.current, iso.previous)
}

// NewContainerStats will return a new instance of ContainerStats
func NewContainerStats(config *ContainerStatsConfig) *ContainerStats {
	return &ContainerStats{
		config:        config,
		curDockerStat: &types.StatsJSON{},
		totalVCHMhz:   uint64(config.VchMhz),
		dblVCHMhz:     uint64(config.VchMhz * 2),
	}
}

// IsListening returns the listening flag
func (cs *ContainerStats) IsListening() bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.listening
}

// Stop will clean up the pipe and flip listening flag
func (cs *ContainerStats) Stop() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.listening {
		cs.reader.Close()
		cs.writer.Close()
		cs.listening = false
	}
}

// newPipe will initialize the pipe for encoding / decoding and
// set the listening flag
func (cs *ContainerStats) newPipe() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// create a new reader / writer
	cs.reader, cs.writer = io.Pipe()
	cs.listening = true
}

// Listen will listen for new metrics from the portLayer, convert to docker format
// and encode to the configured Writer.  The returned PipeWriter is the source of
// the vic metrics that will be transformed to docker stats
func (cs *ContainerStats) Listen() *io.PipeWriter {
	// Are we already listening?
	if cs.IsListening() {
		return nil
	}

	// create pipe for encode/decode
	cs.newPipe()

	dec := json.NewDecoder(cs.reader)
	doc := json.NewEncoder(cs.config.Out)

	// channel to transfer metric from decoder to encoder
	metric := make(chan metrics.VMMetrics)

	// if we aren't streaming and the container is not running, then create an empty
	// docker stat to return
	if !cs.config.Stream && !cs.config.ContainerState.Running {
		cs.preDockerStat = &types.StatsJSON{}
	}

	// go routine will decode metrics received from the portLayer and
	// send them to the encoding routine
	go func() {
		for {
			select {
			case <-cs.config.Ctx.Done():
				close(metric)
				cs.Stop()
				return
			default:
				for dec.More() {
					var vmm metrics.VMMetrics
					err := dec.Decode(&vmm)
					if err != nil {
						log.Errorf("container metric decoding error for container(%s): %s", cs.config.ContainerID, err)
						cs.config.Cancel()
					}
					// send the decoded metric for transform and encoding
					metric <- vmm
				}
			}
		}

	}()

	// go routine will convert incoming metrics to docker specific stats and encode for the docker client.
	go func() {
		// docker needs updates quicker than vSphere can produce metrics, so we'll send a minimum of 1 metric/sec
		ticker := time.NewTicker(time.Millisecond * 500)
		for range ticker.C {
			select {
			case <-cs.config.Ctx.Done():
				cs.Stop()
				ticker.Stop()
				return
			case nm := <-metric:
				// convert the Stat to docker struct
				stat, err := cs.ToContainerStats(&nm)
				if err != nil {
					log.Errorf("container metric conversion error for container(%s): %s", cs.config.ContainerID, err)
					cs.config.Cancel()
				}
				if stat != nil {
					cs.preDockerStat = stat
				}
			default:
				if cs.IsListening() && cs.preDockerStat != nil {
					err := doc.Encode(cs.preDockerStat)
					if err != nil {
						log.Warnf("container metric encoding error for container(%s): %s", cs.config.ContainerID, err)
						cs.config.Cancel()
					}
					// if we aren't streaming then cancel
					if !cs.config.Stream {
						cs.config.Cancel()
					}
				}
			}
		}
	}()

	return cs.writer
}

// ToContainerStats will convert the vic VMMetrics to a docker stats struct -- a complete docker stats
// struct requires two samples.  Func will return nil until a complete stat is available
func (cs *ContainerStats) ToContainerStats(current *metrics.VMMetrics) (*types.StatsJSON, error) {
	// if we have a current metric then validate and transform
	if cs.currentMetric != nil {
		// do we have the same metric as before?
		if cs.currentMetric.SampleTime.Equal(current.SampleTime) {
			// we've already got this as current, so skip and wait for the
			// next sample
			return nil, nil
		}
		// we have new current stats so need to move the previous CPU
		err := cs.previousCPU(current)
		if err != nil {
			return nil, err
		}
	}
	cs.currentMetric = current

	// create the current CPU stats
	cs.currentCPU()

	// create memory stats
	cs.memory()

	// set sample time
	cs.curDockerStat.Read = cs.currentMetric.SampleTime

	// PreRead will be zero if we don't have two samples
	if cs.curDockerStat.PreRead.IsZero() {
		return nil, nil
	}
	return cs.curDockerStat, nil
}

func (cs *ContainerStats) memory() {
	// given MB (i.e. 2048) convert to GB
	cs.curDockerStat.MemoryStats.Limit = uint64(cs.config.Memory * 1024 * 1024)
	// given KB (i.e. 384.5) convert to Bytes
	cs.curDockerStat.MemoryStats.Usage = uint64(cs.currentMetric.Memory.Active * 1024)
}

// previousCPU will move the current stats to the previous CPU location
func (cs *ContainerStats) previousCPU(current *metrics.VMMetrics) error {
	// validate that the sampling is in the correct order
	if current.SampleTime.Before(cs.curDockerStat.Read) {
		err := InvalidOrderError{
			current:  current.SampleTime,
			previous: cs.curDockerStat.Read,
		}
		return err
	}

	// move the stats
	cs.curDockerStat.PreCPUStats = cs.curDockerStat.CPUStats

	// set the previousTotal -- this will be added to the current CPU
	cs.preTotalMhz = cs.curDockerStat.PreCPUStats.CPUUsage.TotalUsage

	cs.curDockerStat.PreRead = cs.curDockerStat.Read
	// previous systemUsage will always be the VCH total
	// see note in func currentCPU() for detail
	cs.curDockerStat.PreCPUStats.SystemUsage = cs.totalVCHMhz

	return nil
}

// currentCPU will convert the VM CPU metrics to docker CPU stats
func (cs *ContainerStats) currentCPU() {
	cpuCount := len(cs.currentMetric.CPU.CPUs)
	dockerCPU := types.CPUStats{
		CPUUsage: types.CPUUsage{
			PercpuUsage: make([]uint64, cpuCount, cpuCount),
		},
	}

	// collect the current CPU Metrics
	for ci, current := range cs.currentMetric.CPU.CPUs {
		dockerCPU.CPUUsage.PercpuUsage[ci] = uint64(current.MhzUsage)
		dockerCPU.CPUUsage.TotalUsage += uint64(current.MhzUsage)
	}

	// vSphere will report negative usage for a starting VM, lets
	// set to zero
	if dockerCPU.CPUUsage.TotalUsage < 0 {
		dockerCPU.CPUUsage.TotalUsage = 0
	}

	// The first stat available for a VM will be missing detail
	if cpuCount > 0 {
		// TotalUsage is the sum of the individual vCPUs Mhz
		// consumption this reading.  We must divide that by the
		// number of vCPUs to get the average across both, since
		// the cpuUsage calc (explained below) will multiply by
		// the number of CPUs to get the cpuUsage percent
		dockerCPU.CPUUsage.TotalUsage /= uint64(cpuCount)
	}

	// Set the current systemUsage to double the VCH as the
	// previous systemUsage is the VCH total.  The docker
	// client formula creates a SystemDelta which is the following:
	// systemDelta = currentSystemUsage - previousSystemUsage
	// We always need systemDelta to equal the total amount of
	// VCH Mhz thus the doubling here.
	dockerCPU.SystemUsage = cs.dblVCHMhz

	// Much like systemUsage (above) totalCPUUsage and previous
	// totalCPUUsage will be used to create a CPUUsage delta as such:
	// CPUDelta = currentTotalCPUUsage - previousTotalCPUUsage
	// This amount will then be divided by the systemDelta
	// (explained above) as part of the CPU % Usage calculation
	// cpuUsage = (CPUDelta / SystemDelta) * cpuCount * 100
	// This will require the addition of the previous total usage
	dockerCPU.CPUUsage.TotalUsage += cs.preTotalMhz
	cs.curDockerStat.CPUStats = dockerCPU
}
