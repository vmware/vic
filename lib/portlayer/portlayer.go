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
	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

var Config metadata.VirtualContainerHostConfigSpec

type API interface {
	storage.ImageStorer
}

func Init(ctx context.Context, sess *session.Session) error {
	source, err := extraconfig.GuestInfoSource()
	if err != nil {
		return err
	}

	extraconfig.DecodeWithPrefix(source, &Config, "guestinfo.vch")
	log.Debugf("decoded vch config: %#v", Config)
	f := find.NewFinder(sess.Vim25(), false)
	for nn, n := range Config.Networks {
		r, err := f.ObjectReference(ctx, n.PortGroupRef)
		if err != nil {
			log.Warnf("could not get network reference for %s network", nn)
			continue
		}

		n.PortGroup = r.(object.NetworkReference)
	}

	return nil
}
