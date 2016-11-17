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

package portlayer

import (
	"path"

	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/portlayer/attach"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/network"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/store"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// API defines the interface the REST server used by the portlayer expects the
// implementation side to export
type API interface {
	storage.ImageStorer
	storage.VolumeStorer
}

func Init(ctx context.Context, sess *session.Session) error {
	source, err := extraconfig.GuestInfoSource()
	if err != nil {
		return err
	}

	sink, err := extraconfig.GuestInfoSink()
	if err != nil {
		return err
	}

	// Grab the storage layer config blobs from extra config
	extraconfig.Decode(source, &storage.Config)
	log.Debugf("Decoded VCH config for storage: %#v", storage.Config)

	// create or restore a portlayer k/v store in the VCH's directory.
	vch, err := guest.GetSelf(ctx, sess)
	if err != nil {
		return err
	}

	vchvm := vm.NewVirtualMachineFromVM(ctx, sess, vch)
	vmPath, err := vchvm.VMPathName(ctx)
	if err != nil {
		return err
	}

	// vmPath is set to the vmx.  Grab the directory from that.
	vmFolder, err := datastore.ToURL(path.Dir(vmPath))
	if err != nil {
		return err
	}

	if err = store.Init(ctx, sess, vmFolder); err != nil {
		return err
	}

	if err := exec.Init(ctx, sess, source, sink); err != nil {
		return err
	}

	if err = network.Init(ctx, sess, source, sink); err != nil {
		return err
	}

	// Unbind containerVM serial ports configured with the old VCH IP.
	// Useful when the appliance restarts and the VCH has a different IP.
	unbindSerialPorts(sess)

	return nil
}

// unbindSerialPorts disconnects serial ports backed by network on the VCH's old IP
// for running containers. This is useful when the appliance or the portlayer restarts
// and the VCH has a new IP. Any errors are logged and portlayer init proceeds as usual.
func unbindSerialPorts(sess *session.Session) {
	// Get all running containers from the cache
	runningState := new(exec.State)
	*runningState = exec.StateRunning
	containers := exec.Containers.Containers(runningState)

	for i := range containers {
		if containers[i].ExecConfig != nil {
			log.Infof("unbinding serial port for running container: %s", containers[i].ExecConfig.ID)
		}

		// Obtain a container handle
		handle := containers[i].NewHandle(context.Background())
		if handle == nil {
			log.Error("unable to obtain a handle for container")
			continue
		}

		// Unbind the VirtualSerialPort
		newHandle, err := attach.Unbind(handle)
		if err != nil {
			log.Errorf("unable to unbind serial port for container: %s", err)
			continue
		}

		// Commit the handle
		if execHandle, ok := newHandle.(*exec.Handle); !ok {
			log.Error("handle type assertion failed for container")
		} else {
			if err := execHandle.Commit(context.Background(), sess, nil); err != nil {
				log.Errorf("unable to commit handle for container: %s", err)
			}
		}
	}
}
