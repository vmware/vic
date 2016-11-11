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

package vsphere

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
)

const mapFile = "parentMap"

// Parent relationships This file will go away when First Class Disk
// support is added to vsphere.  Currently, we can't get a disk spec for a
// disk after creating the disk (and the spec) from vsphere.  Basically, if
// we have a vmdk in the datastore, we have no way of getting the delta
// disk's (if it even is a delta disk) spec to find it's immediate parent.
// This map is used to persist the parent relationship for the disk which
// we maintain outside of the vsphere API.  So, for now, persist this data
// in the datastore and look it up when we need it.  We write a map file
// every time a new disk is created, and swap it with the original map
// file.  Then, at startup, we read the file, and rebuild this map in
// memory.  At runtime, we consult the map to find which disk is the parent
// of a given disk.

// Implements the cache used to lookup an image's parent
type parentM struct {
	// map of image ID to parent ID
	db map[string]string

	// roots where the map is stored
	ds *datastore.Helper

	parentMFile string

	l sync.Mutex
}

// Starts here.  Tries to create a new parentM or load an existing one.
func restoreParentMap(op trace.Operation, ds *datastore.Helper, storeName string) (*parentM, error) {
	p := &parentM{
		ds: ds,
	}
	p.parentMFile = path.Join(storeName, mapFile)

	// Download the map file
	if err := p.download(op); err != nil {
		log.Infof("err = %#v", err)
		return nil, err
	}

	return p, nil
}

// Add sets the parent for image i to parent
func (p *parentM) Add(i string, parent string) {
	p.l.Lock()
	defer p.l.Unlock()

	p.db[i] = parent
}

// Get gets a given image's parent
func (p *parentM) Get(i string) string {
	p.l.Lock()
	defer p.l.Unlock()

	return p.db[i]
}

// Save persists the parent map to the datastore
func (p *parentM) Save(op trace.Operation) error {
	p.l.Lock()
	defer p.l.Unlock()

	buf, err := json.Marshal(p.db)
	if err != nil {
		return err
	}

	// upload to an ephemeral file
	tmpURI := p.parentMFile + ".tmp"

	r := bytes.NewReader(buf)
	if err = p.ds.Upload(op, r, tmpURI); err != nil {
		log.Errorf("Error uploading %s: %s", tmpURI, err)
		return err
	}

	log.Infof("Saving parent map (%s)", p.parentMFile)
	if err := p.ds.Mv(op, tmpURI, p.parentMFile); err != nil {
		log.Errorf("Error moving %s: %s", tmpURI, err)
		return err
	}

	return nil
}

func (p *parentM) download(op trace.Operation) error {
	p.l.Lock()
	defer p.l.Unlock()

	p.db = make(map[string]string)

	rc, err := p.ds.Download(op, p.parentMFile)
	if err != nil {
		// We need to check for 404 vs something else here.
		return nil
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf, &p.db); err != nil {
		return err
	}

	return nil
}
