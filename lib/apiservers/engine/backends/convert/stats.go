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
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/api/types"

	"github.com/vmware/vic/lib/portlayer/metrics"
)

// ContainerStats encapsulates the conversion of VMMetrics to
// docker specific metrics
type ContainerStats struct {
	config         ContainerStatsConfig
	dockerStats    *types.StatsJSON
	currentMetrics *metrics.VMMetrics
	totalVCHMhz    uint64
	dblVCHMhz      uint64
	preTotalMhz    uint64

	// reader/writer for stream
	reader *io.PipeReader
	writer *io.PipeWriter
}

type ContainerStatsConfig struct {
	Ctx         context.Context
	Cancel      context.CancelFunc
	Out         io.Writer
	ContainerID string
	Memory      int64
	Stream      bool
	VchMhz      int64
}

type InvalidOrderError struct {
	current  time.Time
	previous time.Time
}

func (iso InvalidOrderError) Error() string {
	return fmt.Sprintf("The current sample time (%s) is before the previous time (%s)", iso.current, iso.previous)
}

// NewContainerStats will return a new instance of ContainerStats
func NewContainerStats(config ContainerStatsConfig) *ContainerStats {
	return &ContainerStats{
		config:      config,
		dockerStats: &types.StatsJSON{},
		totalVCHMhz: uint64(config.VchMhz),
		dblVCHMhz:   uint64(config.VchMhz * 2),
	}
}

// Stop will clean up remaining conversion resources
func (cs *ContainerStats) Stop() {
	if cs.reader != nil && cs.writer != nil {
		cs.reader.Close()
		cs.writer.Close()
	}
}

// Listen will listen for new metrics from the portLayer, convert to docker format
// and encode to the configured Writer.  The returned PipeWriter is the source of
// the vic metrics that will be transformed to docker stats
func (cs *ContainerStats) Listen() *io.PipeWriter {
	// TODO: could split decode / encode into separate funcs -- would provide for easier
	// unit testing

	// we already are listening
	if cs.reader != nil {
		return nil
	}

	cs.reader, cs.writer = io.Pipe()

	dec := json.NewDecoder(cs.reader)
	doc := json.NewEncoder(cs.config.Out)

	// channel to transfer metric from decoder to encoder
	// closed w/in the decoder
	metric := make(chan metrics.VMMetrics)

	// signal to decoder / encoder that we are done
	finished := make(chan struct{})

	var vmm metrics.VMMetrics
	var previousStat *types.StatsJSON

	go func() {
		for {
			select {
			case <-cs.config.Ctx.Done():
				close(finished)
			case <-finished:
				close(metric)
				cs.Stop()
				return
			default:
				for dec.More() {
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

	go func() {
		for {
			select {
			case <-cs.config.Ctx.Done():
				return
			case <-finished:
				return
			case nm := <-metric:
				// convert the Stat to docker struct
				stats, err := cs.ToContainerStats(&nm)
				if err != nil {
					log.Errorf("container metric conversion error for container(%s): %s", cs.config.ContainerID, err)
					cs.config.Cancel()
				}
				// Do we have a complete stat that can be sent to the client?
				if stats != nil {
					err = doc.Encode(stats)
					if err != nil {
						log.Warnf("container metric encoding error for container(%s): %s", cs.config.ContainerID, err)
						cs.config.Cancel()
					}
					// if we aren't streaming then cancel
					if !cs.config.Stream {
						cs.config.Cancel()
					}
					// set to previous stat so we can reuse
					previousStat = stats
				}
			default:
				// the docker client expects updates quicker than vSphere can produce them, so
				// we need to send the previous stats to avoid intermittent empty output
				time.Sleep(time.Second * 1)
				if previousStat != nil && cs.reader != nil {
					err := doc.Encode(previousStat)
					if err != nil {
						log.Warnf("container previous metric encoding error for container(%s): %s", cs.config.ContainerID, err)
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
	if cs.currentMetrics != nil {
		// do we have the same metric as before?
		if cs.currentMetrics.SampleTime.Equal(current.SampleTime) {
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
	cs.currentMetrics = current

	// create the current CPU stats
	cs.currentCPU()

	// create memory stats
	cs.memory()

	// set sample time
	cs.dockerStats.Read = cs.currentMetrics.SampleTime

	// PreRead will be zero if we don't have two samples
	if cs.dockerStats.PreRead.IsZero() {
		return nil, nil
	}
	return cs.dockerStats, nil
}

func (cs *ContainerStats) memory() {
	// given MB (i.e. 2048) convert to GB
	cs.dockerStats.MemoryStats.Limit = uint64(cs.config.Memory * 1024 * 1024)
	// given KB (i.e. 384.5) convert to Bytes
	cs.dockerStats.MemoryStats.Usage = uint64(cs.currentMetrics.Memory.Active * 1024)
}

// previousCPU will move the current stats to the previous CPU location
func (cs *ContainerStats) previousCPU(current *metrics.VMMetrics) error {
	// validate that the sampling is in the correct order
	if current.SampleTime.Before(cs.dockerStats.Read) {
		err := InvalidOrderError{
			current:  current.SampleTime,
			previous: cs.dockerStats.Read,
		}
		return err
	}

	// move the stats
	cs.dockerStats.PreCPUStats = cs.dockerStats.CPUStats

	// set the previousTotal -- this will be added to the current CPU
	cs.preTotalMhz = cs.dockerStats.PreCPUStats.CPUUsage.TotalUsage

	cs.dockerStats.PreRead = cs.dockerStats.Read
	// previous systemUsage will always be the VCH total
	// see note in func currentCPU() for detail
	cs.dockerStats.PreCPUStats.SystemUsage = cs.totalVCHMhz

	return nil
}

// currentCPU will convert the VM CPU metrics to docker CPU stats
func (cs *ContainerStats) currentCPU() {
	cpuCount := len(cs.currentMetrics.CPU.CPUs)
	dockerCPU := types.CPUStats{
		CPUUsage: types.CPUUsage{
			PercpuUsage: make([]uint64, cpuCount, cpuCount),
		},
	}

	// collect the current CPU Metrics
	for ci, current := range cs.currentMetrics.CPU.CPUs {
		dockerCPU.CPUUsage.PercpuUsage[ci] = uint64(current.MhzUsage)
		dockerCPU.CPUUsage.TotalUsage += uint64(current.MhzUsage)
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
	cs.dockerStats.CPUStats = dockerCPU
}
