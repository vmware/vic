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

package validate

import (
	"fmt"
	"net/url"
	"strings"

	"context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
)

func (v *Validator) storage(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// Image Store
	imageDSpath, ds, err := v.DatastoreHelper(ctx, input.ImageDatastorePath, "", "--image-store")

	if err != nil {
		v.NoteIssue(err)
		return
	}

	// provide a default path if only a DS name is provided
	if imageDSpath.Path == "" {
		log.Debug("No path element specified for image store - will use default")
	}

	if ds != nil {
		v.SetDatastore(ds, imageDSpath)
		conf.AddImageStore(imageDSpath)
	}

	if conf.VolumeLocations == nil {
		conf.VolumeLocations = make(map[string]*url.URL)
	}

	for label, targetURL := range input.VolumeLocations {
		var vsErr error
		switch targetURL.Scheme {
		case common.NfsScheme:
			vsErr := validateNFSTarget(targetURL)
			v.NoteIssue(vsErr)
		case common.DsScheme:
			// TODO: change v.DatastoreHelper to take url struct instead of string and modify tests.
			targetURL, _, vsErr = v.DatastoreHelper(ctx, targetURL.String(), label, "--volume-store")
			v.NoteIssue(vsErr)
		default:
			// We should not reach here, if we do we will attempt to treat this as a vsphere datastore
			targetURL, _, vsErr = v.DatastoreHelper(ctx, targetURL.String(), label, "--volume-store")
			v.NoteIssue(vsErr)
		}

		// skip adding volume stores that we know will fail
		if vsErr != nil {
			continue
		}

		conf.VolumeLocations[label] = targetURL
	}
}

func validateNFSTarget(nfsURL *url.URL) error {
	if nfsURL.Host == "" {
		return fmt.Errorf("volume store target (%s) is missing the host field. format: <nfs://<user>:<password>@<host>/<share point path>:label", nfsURL.String())
	}

	if nfsURL.Path == "" {
		return fmt.Errorf("volume store target (%s) is missing the path field. format: <nfs://<host>/<share point path>?<mount options as query vars>:label", nfsURL.String())
	}

	return nil
}

func (v *Validator) DatastoreHelper(ctx context.Context, path string, label string, flag string) (*url.URL, *object.Datastore, error) {
	defer trace.End(trace.Begin(path))
	dsURL, dsErr := url.Parse(path)
	if dsErr != nil {
		return nil, nil, errors.Errorf("error parsing datastore path: %s", dsErr)
	}

	// url scheme does not contain ://, so remove it to make url work
	if dsURL.Scheme != "" && dsURL.Scheme != "ds" {
		return nil, nil, errors.Errorf("bad scheme %q provided for datastore", dsURL.Scheme)
	}

	dsURL.Scheme = common.DsScheme

	// if a datastore name (e.g. "datastore1") is specified with no decoration then this
	// is interpreted as the Path
	if dsURL.Host == "" && dsURL.Path != "" {
		pathElements := strings.SplitN(path, "/", 2)
		dsURL.Host = pathElements[0]
		if len(pathElements) > 1 {
			dsURL.Path = pathElements[1]
		} else {
			dsURL.Path = ""
		}
	}

	if dsURL.Host == "" {
		// see if we can find a default datastore
		store, err := v.Session.Finder.DatastoreOrDefault(ctx, "*")
		if err != nil {
			v.suggestDatastore("*", label, flag)
			return nil, nil, errors.New("datastore empty")
		}

		dsURL.Host = store.Name()
		log.Infof("Using default datastore: %s", dsURL.Host)
	}

	stores, err := v.Session.Finder.DatastoreList(ctx, dsURL.Host)
	if err != nil {
		log.Debugf("no such datastore %#v", dsURL)
		v.suggestDatastore(path, label, flag)
		// TODO: error message about no such match and how to get a datastore list
		// we return err directly here so we can check the type
		return nil, nil, err
	}
	if len(stores) > 1 {
		// TODO: error about required disabmiguation and list entries in stores
		v.suggestDatastore(path, label, flag)
		return nil, nil, errors.New("ambiguous datastore " + dsURL.Host)
	}

	// temporary until session is extracted
	// FIXME: commented out until components can consume moid
	// dsURL.Host = stores[0].Reference().Value

	// make sure the vsphere ds format fits the right format
	if _, err := datastore.ToURL(fmt.Sprintf("[%s] %s", dsURL.Host, dsURL.Path)); err != nil {
		return nil, nil, err
	}

	return dsURL, stores[0], nil
}

func (v *Validator) SetDatastore(ds *object.Datastore, path *url.URL) {
	v.Session.Datastore = ds
	v.Session.DatastorePath = path.Host
}

// suggestDatastore suggests all datastores present on target in datastore:label format if applicable
func (v *Validator) suggestDatastore(path string, label string, flag string) {
	defer trace.End(trace.Begin(""))

	var val string
	if label != "" {
		val = fmt.Sprintf("%s:%s", path, label)
	} else {
		val = path
	}
	log.Infof("Suggesting valid values for %s based on %q", flag, val)

	dss, err := v.Session.Finder.DatastoreList(v.Context, "*")
	if err != nil {
		log.Errorf("Unable to list datastores: %s", err)
		return
	}

	if len(dss) == 0 {
		log.Info("No datastores found")
		return
	}

	matches := make([]string, len(dss))
	for i, d := range dss {
		if label != "" {
			matches[i] = fmt.Sprintf("%s:%s", d.Name(), label)
		} else {
			matches[i] = d.Name()
		}
	}

	if matches != nil {
		log.Infof("Suggested values for %s:", flag)
		for _, d := range matches {
			log.Infof("  %q", d)
		}
	}
}
