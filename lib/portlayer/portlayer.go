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
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/network"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

const (
	// docker official ports are 2375 and 2376
	SerialOverLANPort = 8080
)

// XXX TODO(FA) use this in the _handlers the swagger server includes.
//
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

	f := find.NewFinder(sess.Vim25(), false)

	extraconfig.Decode(source, &exec.Config)
	log.Debugf("Decoded VCH config for execution: %#v", exec.Config)
	ccount := len(exec.Config.ComputeResources)
	if ccount != 1 {
		detail := fmt.Sprintf("expected singular compute resource element, found %d", ccount)
		log.Errorf(detail)
		return err
	}

	cr := exec.Config.ComputeResources[0]
	r, err := f.ObjectReference(ctx, cr)
	if err != nil {
		detail := fmt.Sprintf("could not get resource pool reference from %s: %s", cr.String(), err)
		log.Errorf(detail)
		return err
	}
	exec.Config.ResourcePool = r.(*object.ResourcePool)
	//FIXME: temporary injection of debug network for debug nic
	ne := exec.Config.Networks["client"]
	if ne == nil {
		detail := fmt.Sprintf("could not get client network reference for debug nic - this code can be removed once network mapping/dhcp client is present")
		log.Errorf(detail)
		return err
	}
	nr := new(types.ManagedObjectReference)
	nr.FromString(ne.Network.ID)
	r, err = f.ObjectReference(ctx, *nr)
	if err != nil {
		detail := fmt.Sprintf("could not get client network reference from %s: %s", nr.String(), err)
		log.Errorf(detail)
		return err
	}
	exec.Config.DebugNetwork = r.(object.NetworkReference)

	extraconfig.Decode(source, &network.Config)
	log.Debugf("Decoded VCH config for network: %#v", network.Config)
	for nn, n := range network.Config.ContainerNetworks {
		pgref := new(types.ManagedObjectReference)
		if !pgref.FromString(n.ID) {
			log.Errorf("Could not reacquire object reference from id for network %s: %s", nn, n.ID)
		}

		r, err = f.ObjectReference(ctx, *pgref)
		if err != nil {
			log.Warnf("could not get network reference for %s network", nn)
			continue
		}

		n.PortGroup = r.(object.NetworkReference)
	}

	// Grab the storage layer config blobs from extra config
	extraconfig.Decode(source, &storage.Config)
	log.Debugf("Decoded VCH config for storage: %#v", storage.Config)
	return nil
}
