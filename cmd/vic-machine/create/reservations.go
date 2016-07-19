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
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
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

	var errs []string

	elements := strings.Split(m, ":")
	if len(elements) != 2 {
		return cli.NewExitError("--cpu-limits value must be reservations:limits format, e.g. e.g. 800MHz:20GHz, 800MHz:, :20GHz", 1)
	}
	if elements[0] == "" && elements[1] == "" {
		return cli.NewExitError(fmt.Sprintf("nothing is set with --cpu-limits %q", m), 1)
	}
	log.Debugf("cpu reservation: %q, limit: %q", elements[0], elements[1])

	var err error
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
	if len(m) < 4 {
		return 0, errors.New("CPU size should end with MHz or GHz")
	}

	suffix := m[len(m)-3 : len(m)]
	sn := m[:len(m)-3]

	var errs []string
	cpu, err := strconv.Atoi(sn)
	if err != nil {
		err = errors.Errorf("CPU size %q is not int number, %s", sn, err)
		errs = append(errs, err.Error())
	}
	var times int
	switch suffix = strings.ToLower(suffix); suffix {
	case "mhz":
		times = 1
	case "ghz":
		times = 1000
	default:
		err = errors.Errorf("Invalid size unit %q, only MHz or GHz is accepted", suffix)
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return 0, errors.New(strings.Join(errs, "\n"))
	}
	return cpu * times, nil
}

func (c *Create) handleMemeoryReservation(m string) error {
	defer trace.End(trace.Begin(m))
	if m == "" {
		return nil
	}

	var errs []string

	elements := strings.Split(m, ":")
	if len(elements) != 2 {
		return cli.NewExitError("--memory-limits value must be reservations:limits format, e.g. 800MB:8GB, 800MB:, :8GB", 1)
	}
	if elements[0] == "" && elements[1] == "" {
		return cli.NewExitError(fmt.Sprintf("nothing is set with --memory-limits %q", m), 1)
	}
	log.Debugf("memory reservation: %q, limit: %q", elements[0], elements[1])

	var err error
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
	if len(m) < 3 {
		return 0, errors.New("Memory size should end with MB or GB")
	}

	suffix := m[len(m)-2 : len(m)]
	sn := m[:len(m)-2]

	var errs []string
	memory, err := strconv.Atoi(sn)
	if err != nil {
		err = errors.Errorf("Memeory size %q is not int number, %s", sn, err)
		errs = append(errs, err.Error())
	}
	var times int
	switch suffix = strings.ToLower(suffix); suffix {
	case "mb":
		times = 1
	case "gb":
		times = 1024
	default:
		err = errors.Errorf("Invalid size unit %q, only MB or GB is accepted", suffix)
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return 0, errors.New(strings.Join(errs, "\n"))
	}
	return memory * times, nil
}
