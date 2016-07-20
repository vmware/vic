// Copyright 2016 VMware, Inc. All Rights Reserved.
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

package create

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

var (
	MemoryPattern = "(?i)^([0-9]+(g|m)b)?:([0-9]+(g|m)b)?"
	CPUPattern    = "(?i)^([0-9]+(g|m)hz)?:([0-9]+(g|m)hz)?"
)

func (c *Create) processReservations() error {
	defer trace.End(trace.Begin(""))
	if err := c.handleMemeoryReservation(c.memoryReservLimits); err != nil {
		return err
	}

	if err := c.handleCPUReservation(c.cpuReservLimits); err != nil {
		return err
	}
	return nil
}

func (c *Create) handleCPUReservation(m string) error {
	defer trace.End(trace.Begin(m))
	if m == "" {
		return nil
	}
	var match bool
	var err error
	if match, err = regexp.MatchString(CPUPattern, m); err != nil {
		err = errors.Errorf("Failed to compile pattern %q: %s", MemoryPattern, err)
		log.Error(err)
		return cli.NewExitError(err.Error(), 1)
	}

	if !match {
		err = errors.Errorf("--cpu-limits %q must be reservations:limits format, e.g. 800MHz:20GHz, 800MHz:, :20GHz", m)
		return cli.NewExitError(err.Error(), 1)
	}

	var errs []string

	elements := strings.Split(m, ":")
	if elements[0] == "" && elements[1] == "" {
		return cli.NewExitError(fmt.Sprintf("nothing is set with --cpu-limits %q", m), 1)
	}
	log.Debugf("cpu reservation: %q, limit: %q", elements[0], elements[1])

	if c.Data.VCHCPUReservationsMHz, err = c.handleCPUSize(elements[0]); err != nil {
		errs = append(errs, err.Error())
	}
	if c.Data.VCHCPULimitsMHz, err = c.handleCPUSize(elements[1]); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return cli.NewExitError(fmt.Sprintf("--cpu-limits value format error:\n %s", strings.Join(errs, "\n")), 1)
	}
	return nil
}

func (c *Create) handleCPUSize(m string) (int, error) {
	defer trace.End(trace.Begin(m))
	if m == "" {
		return 0, nil
	}

	suffix := m[len(m)-3 : len(m)]
	sn := m[:len(m)-3]

	cpu, err := strconv.Atoi(sn)
	if err != nil {
		err = errors.Errorf("CPU size %q is not int number, %s", sn, err)
		return 0, err
	}
	var times int
	switch suffix = strings.ToLower(suffix); suffix {
	case "mhz":
		times = 1
	case "ghz":
		times = 1000
	default:
	}
	return cpu * times, nil
}

func (c *Create) handleMemeoryReservation(m string) error {
	defer trace.End(trace.Begin(m))
	if m == "" {
		return nil
	}

	var match bool
	var err error
	if match, err = regexp.MatchString(MemoryPattern, m); err != nil {
		return err
	}
	if !match {
		err = errors.Errorf("--memory-limits %q must be reservations:limits format, e.g. 800MHz:20GHz, 800MHz:, :20GHz", m)
		return cli.NewExitError(err.Error(), 1)
	}

	var errs []string

	elements := strings.Split(m, ":")
	if elements[0] == "" && elements[1] == "" {
		return cli.NewExitError(fmt.Sprintf("nothing is set with --memory-limits %q", m), 1)
	}
	log.Debugf("memory reservation: %q, limit: %q", elements[0], elements[1])

	if c.Data.VCHMemoryReservationsMB, err = c.handleMemeorySize(elements[0]); err != nil {
		errs = append(errs, err.Error())
	}
	if c.Data.VCHMemoryLimitsMB, err = c.handleMemeorySize(elements[1]); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return cli.NewExitError(fmt.Sprintf("--memory-limits value format error:\n %s", strings.Join(errs, "\n")), 1)
	}
	return nil
}

func (c *Create) handleMemeorySize(m string) (int, error) {
	defer trace.End(trace.Begin(m))
	if m == "" {
		return 0, nil
	}

	suffix := m[len(m)-2 : len(m)]
	sn := m[:len(m)-2]

	memory, err := strconv.Atoi(sn)
	if err != nil {
		err = errors.Errorf("Memeory size %q is not int number, %s", sn, err)
		return 0, err
	}
	var times int
	switch suffix = strings.ToLower(suffix); suffix {
	case "mb":
		times = 1
	case "gb":
		times = 1024
	default:
	}
	return memory * times, nil
}
