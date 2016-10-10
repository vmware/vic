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

package disk

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

const (
	// The duration waitForPath will tolerate before timing out.
	// TODO FIXME see GH issues 2340 and 2385
	// TODO We need to add a vSphere cancellation step to cancel calls that are taking too long
	// TODO Remove these TODOs after 2385 is completed
	pathTimeout = 60 * time.Second
)

func waitForPath(op trace.Operation, path string) error {
	defer trace.End(trace.Begin(path))
	timeout := time.Duration(pathTimeout)

	op, _ = trace.WithTimeout(&op, timeout, path)
	done := make(chan struct{})

	go func() {
		t := time.NewTicker(200 * time.Microsecond)
		defer t.Stop()
		for range t.C {
			if _, err := os.Stat(path); err == nil {
				close(done)
				break
			}

			// We've timed out.
			if op.Err() != nil {
				break
			}
		}
	}()

	log.Debugf("Waiting for attached disk to appear in /dev/disk/by-path, or timeout")
	select {
	case <-done:
		log.Infof("Attached disk present at %s", path)
	case <-op.Done():
		if op.Err() != nil {
			return errors.Errorf("timeout waiting for layer to present as %s", path)
		}
	}

	return nil
}

// ensures that a paravirtual scsi controller is present and determines the
// base path of disks attached to it returns a handle to the controller and a
// format string, with a single decimal for the disk unit number which will
// result in the /dev/disk/by-path path
func verifyParavirtualScsiController(op trace.Operation, vm *object.VirtualMachine) (*types.ParaVirtualSCSIController, string, error) {
	devices, err := vm.Device(op)
	if err != nil {
		log.Errorf("vmware driver failed to retrieve device list for VM %s: %s", vm, errors.ErrorStack(err))
		return nil, "", errors.Trace(err)
	}

	controller, ok := devices.PickController((*types.ParaVirtualSCSIController)(nil)).(*types.ParaVirtualSCSIController)
	if controller == nil || !ok {
		err = errors.Errorf("vmware driver failed to find a paravirtual SCSI controller - ensure setup ran correctly")
		log.Error(err.Error())
		return nil, "", errors.Trace(err)
	}

	// build the base path
	// first we determine which label we're looking for (requires VMW hardware version >=10)
	targetLabel := fmt.Sprintf("SCSI%d", controller.BusNumber)
	log.Debugf("Looking for scsi controller with label %s", targetLabel)

	pciBase := "/sys/bus/pci/devices"
	pciBus, err := os.Open(pciBase)
	if err != nil {
		log.Errorf("Failed to open %s for reading: %s", pciBase, errors.ErrorStack(err))
		return controller, "", errors.Trace(err)
	}
	defer pciBus.Close()

	pciDevices, err := pciBus.Readdirnames(0)
	if err != nil {
		log.Errorf("Failed to read contents of %s: %s", pciBase, errors.ErrorStack(err))
		return controller, "", errors.Trace(err)
	}

	var buf = make([]byte, len(targetLabel))
	var controllerName string

	for _, n := range pciDevices {
		nlabel := fmt.Sprintf("%s/%s/label", pciBase, n)
		flabel, err := os.Open(nlabel)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Errorf("Unable to read label from %s: %s", nlabel, errors.ErrorStack(err))
			}
			continue
		}
		defer flabel.Close()

		_, err = flabel.Read(buf)
		if err != nil {
			log.Errorf("Unable to read label from %s: %s", nlabel, errors.ErrorStack(err))
			continue
		}

		if targetLabel == string(buf) {
			// we've found our controller
			controllerName = n
			log.Debugf("Found pvscsi controller directory: %s", controllerName)

			break
		}
	}

	if controllerName == "" {
		err := errors.Errorf("Failed to locate pvscsi controller directory")
		log.Errorf(err.Error())
		return controller, "", errors.Trace(err)
	}

	formatString := fmt.Sprintf("/dev/disk/by-path/pci-%s-scsi-0:0:%%d:0", controllerName)
	log.Debugf("Disk location format: %s", formatString)
	return controller, formatString, nil
}

// Find the disk by name attached to the given vm.
func findDisk(op trace.Operation, vm *object.VirtualMachine, name string) (*types.VirtualDisk, error) {
	defer trace.End(trace.Begin(vm.String()))

	log.Debugf("Looking for attached disk matching filename %s", name)

	devices, err := vm.Device(op)
	if err != nil {
		return nil, fmt.Errorf("Failed to refresh devices for vm: %s", errors.ErrorStack(err))
	}

	candidates := devices.Select(func(device types.BaseVirtualDevice) bool {
		db := device.GetVirtualDevice().Backing
		if db == nil {
			return false
		}

		backing, ok := device.GetVirtualDevice().Backing.(*types.VirtualDiskFlatVer2BackingInfo)
		if !ok {
			return false
		}

		log.Debugf("backing file name %s", backing.VirtualDeviceFileBackingInfo.FileName)
		match := strings.HasSuffix(backing.VirtualDeviceFileBackingInfo.FileName, name)
		if match {
			log.Debugf("Found candidate disk for %s at %s", name, backing.VirtualDeviceFileBackingInfo.FileName)
		}

		return match
	})

	if len(candidates) == 0 {
		log.Warnf("No disks match name: %s", name)
		return nil, os.ErrNotExist
	}

	if len(candidates) > 1 {
		return nil, errors.Errorf("Too many disks match name: %s", name)
	}

	return candidates[0].(*types.VirtualDisk), nil
}
