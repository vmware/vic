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

package storage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"sync"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
)

const parentMFile = "parentMap"

// Implements the cache used to lookup an image's parent
type parentM struct {
	// location in the datastore in datastore URI format
	mFilePath string

	// map of image ID to parent ID
	db map[string]string

	sess *session.Session

	l sync.Mutex
}

// Starts here.  Tries to create a new parentM or load an existing one.
func restoreParentMap(ctx context.Context, s *session.Session) (*parentM, error) {
	p := &parentM{
		mFilePath: path.Join(datastoreParentPath, parentMFile),
		sess:      s,
	}

	// Download the map file
	if err := p.download(ctx); err != nil {
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
func (p *parentM) Save(ctx context.Context) error {
	p.l.Lock()
	defer p.l.Unlock()

	buf, err := json.Marshal(p.db)
	if err != nil {
		return err
	}

	// upload to an ephemeral file
	tmp := p.mFilePath + ".tmp"
	tmpURI := p.sess.Datastore.Path(tmp)

	parentMURI := p.sess.Datastore.Path(p.mFilePath)

	r := bytes.NewReader(buf)
	if err = p.sess.Datastore.Upload(ctx, r, tmp, &soap.DefaultUpload); err != nil {
		log.Errorf("Error uploading %s: %s", tmp, err)
		return err
	}

	fm := object.NewFileManager(p.sess.Vim25())
	err = tasks.Wait(ctx, func(context.Context) (tasks.Waiter, error) {
		log.Infof("Saving parent map (%s)", p.mFilePath)
		return fm.MoveDatastoreFile(ctx, tmpURI, nil, parentMURI, nil, true)
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *parentM) download(ctx context.Context) error {
	p.l.Lock()
	defer p.l.Unlock()

	p.db = make(map[string]string)

	rc, _, err := p.sess.Datastore.Download(ctx, p.mFilePath, &soap.DefaultDownload)
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
