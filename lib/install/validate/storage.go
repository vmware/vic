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

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"golang.org/x/net/context"
)

func (v *Validator) storage(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// Image Store
	imageDSpath, ds, err := v.DatastoreHelper(ctx, input.ImageDatastoreName, "", "--image-datastore")
	v.NoteIssue(err)
	if ds != nil {
		v.SetDatastore(ds, imageDSpath)
		conf.AddImageStore(imageDSpath)
	}

	// Volume Store(s)
	if conf.VolumeLocations == nil {
		conf.VolumeLocations = make(map[string]*url.URL)
	}

	// TODO: add volume locations
	for label, volDSpath := range input.VolumeLocations {
		dsURL, _, err := v.DatastoreHelper(ctx, volDSpath, label, "--volume-store")
		v.NoteIssue(err)
		if dsURL != nil {
			conf.VolumeLocations[label] = dsURL
		}
	}
}

func (v *Validator) DatastoreHelper(ctx context.Context, path string, label string, flag string) (*url.URL, *object.Datastore, error) {
	defer trace.End(trace.Begin(path))

	dsURL, err := url.Parse(path)
	if err != nil {
		return nil, nil, errors.Errorf("error parsing datastore path: %s", err)
	}

	// url scheme does not contain ://, so remove it to make url work
	if dsURL.Scheme != "" && dsURL.Scheme != "ds" {
		return nil, nil, errors.Errorf("bad scheme %q provided for datastore", dsURL.Scheme)
	}

	dsURL.Scheme = "ds"

	// if a datastore name (e.g. "datastore1") is specifed with no decoration then this
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
		v.suggestDatastore("*", label, flag)
		return nil, nil, errors.New("datastore hostname empty")
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
		return
	}
}
