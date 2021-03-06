// Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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

package portlayer

import (
	"context"
	"fmt"
	"path"

	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/portlayer/attach"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/logging"
	"github.com/vmware/vic/lib/portlayer/metrics"
	"github.com/vmware/vic/lib/portlayer/network"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/store"
	"github.com/vmware/vic/pkg/retry"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// Init initializes portlayer components at startup
func Init(ctx context.Context, sess *session.Session) error {
	defer trace.End(trace.Begin(""))

	op := trace.ChildFromContext(ctx, "port layer initialization")

	source, err := extraconfig.GuestInfoSource()
	if err != nil {
		return err
	}

	sink, err := extraconfig.GuestInfoSink()
	if err != nil {
		return err
	}

	// create or restore a portlayer k/v store in the VCH's directory.
	vch, err := guest.GetSelf(op, sess)
	if err != nil {
		return err
	}

	vchvm := vm.NewVirtualMachineFromVM(op, sess, vch)
	vmPath, err := vchvm.VMPathName(op)
	if err != nil {
		return err
	}

	// vmPath is set to the vmx.  Grab the directory from that.
	vmFolder, err := datastore.ToURL(path.Dir(vmPath))
	if err != nil {
		return err
	}

	vmParentPool, err := vchvm.ResourcePool(op)
	if err != nil {
		return err
	}

	// store the reference to the vch inventory folder before portlayer init
	vchFolder, err := vchvm.Folder(op)
	if err != nil {
		return err
	}
	sess.VCHFolder = vchFolder

	if err = storage.Init(op, sess, vmParentPool, source, sink); err != nil {
		return err
	}

	if err = store.Init(op, sess, vmFolder); err != nil {
		return err
	}

	if err := exec.Init(op, sess, source, sink, vchvm.Reference()); err != nil {
		return err
	}

	if err = network.Init(op, sess, source, sink); err != nil {
		return err
	}

	if err = logging.Init(op); err != nil {
		return err
	}

	if err = metrics.Init(op, sess); err != nil {
		return err
	}

	// Unbind containerVM serial ports configured with the old VCH IP.
	// Useful when the appliance restarts and the VCH has a different IP.
	TakeCareOfSerialPorts(op, sess)

	return nil
}

func Finalize(ctx context.Context) error {
	op := trace.NewOperation(context.Background(), "Shutdown")
	defer trace.End(trace.Begin("", op))

	storage.Finalize(op)
	store.Finalize(ctx)
	exec.Finalize(op)
	network.Finalize(ctx)
	logging.Finalize(ctx)
	metrics.Finalize(op)

	return nil
}

// TakeCareOfSerialPorts disconnects serial ports backed by network on the VCH's old IP and connects serial ports backed by file.
// This is useful when the appliance or the portlayer restarts and the VCH has a new IP or container vms gets migrated
// Any errors are logged and portlayer init proceeds as usual.
func TakeCareOfSerialPorts(op trace.Operation, sess *session.Session) {
	op = trace.FromOperation(op, "SerialPorts")
	defer trace.End(trace.Begin("", op))
	// Get all running containers from the portlayer cache
	// Including starting containers here as well
	// TODO: for starting containers, if using the runblocking mechanism present as of this date, we should cause the
	// unbind change to blocking status to propagate into the container and release the process for start
	containers := exec.Containers.Containers([]exec.State{exec.StateRunning, exec.StateStarting})

	for i := range containers {
		var containerID string

		if containers[i].ExecConfig != nil {
			containerID = containers[i].ExecConfig.ID
		}
		op.Infof("unbinding serial port for running container %s", containerID)

		operation := func() error {
			// Obtain a container handle
			handle, err := containers[i].NewHandle(op)
			if err != nil {
				err := fmt.Errorf("unable to obtain a handle for container %s: %s", containerID, err)
				op.Errorf("%s", err)

				return err
			}

			// Unbind the network backed VirtualSerialPort
			unbindHandle, err := attach.Unbind(handle, containerID)
			if err != nil {
				err := fmt.Errorf("unable to unbind serial port for container %s: %s", containerID, err)
				op.Errorf("%s", err)

				return err
			}

			execHandle, ok := unbindHandle.(*exec.Handle)
			if !ok {
				err := fmt.Errorf("handle type assertion failed for container %s", containerID)
				op.Errorf("%s", err)

				return err
			}

			// Bind the file backed VirtualSerialPort
			bindHandle, err := logging.Bind(execHandle)
			if err != nil {
				err := fmt.Errorf("unable to unbind serial port for container %s: %s", containerID, err)
				op.Errorf("%s", err)

				return err
			}

			execHandle, ok = bindHandle.(*exec.Handle)
			if !ok {
				err := fmt.Errorf("handle type assertion failed for container %s", containerID)
				op.Errorf("%s", err)

				return err
			}

			// Commit the handle
			if err := execHandle.Commit(op, sess, nil); err != nil {
				op.Errorf("unable to commit handle for container %s: %s", containerID, err)
				return err
			}
			return nil
		}

		if err := retry.Do(op, operation, exec.IsConcurrentAccessError); err != nil {
			op.Errorf("Multiple attempts failed for committing the handle with %s", err)
		}
	}
}
