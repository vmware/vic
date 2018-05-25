// Copyright 2018 VMware, Inc. All Rights Reserved.
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

package decode

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/client"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/errors"
	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/trace"
)

// TODO [AngieCris]: duplicate of common/utils.go
const (
	// scheme string for nfs volume store targets
	NfsScheme = "nfs"

	// scheme string for ds volume store targets
	DsScheme = "ds"

	// scheme string for volume store targets without a scheme
	EmptyScheme = ""

	dsInputFormat = "<datastore url w/ path>"

	nfsInputFormat = "nfs://<host>/<url-path>?<mount option as query parameters>"
)

func ProcessStorage(op trace.Operation, d *data.Data, vch *models.VCH, finder client.Finder) error {
	if vch.Storage != nil {
		// image stores
		// TODO [AngieCris]: is imagestore required?
		if vch.Storage.ImageStores != nil && len(vch.Storage.ImageStores) > 0 {
			path, err := processImageStorePath(op, vch.Storage.ImageStores) // TODO (#6712): many vs. one mismatch
			if err != nil {
				return errors.WrapError(http.StatusBadRequest, err)
			}
			d.ImageDatastorePath = path
		} else {
			return errors.NewError(http.StatusBadRequest, "image store is required")
		}

		if vch.Storage.VolumeStores != nil {
			volumeLocations, err := processVolumeStores(op, vch.Storage.VolumeStores)
			if err != nil {
				return errors.WrapError(http.StatusBadRequest, err)
			}
			d.VolumeLocations = volumeLocations
		}

		d.ScratchSize = constants.DefaultBaseImageScratchSize

		if vch.Storage.BaseImageSize != nil {
			d.ScratchSize = FromValueBytesMetric(vch.Storage.BaseImageSize)
		}
	}
	return nil
}

func processImageStorePath(op trace.Operation, imageStores []string) (string, error) {
	path := imageStores[0]

	err := CheckUnsupportedCharsDatastore(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func processVolumeStores(op trace.Operation, volumeStores []*models.VCHStorageVolumeStoresItems0) (map[string]*url.URL, error) {
	volumeLocations := make(map[string]*url.URL)

	for _, v := range volumeStores {
		label, url, err := processVolumeStore(op, v)
		if err != nil {
			return nil, err
		}
		volumeLocations[label] = url
	}

	return volumeLocations, nil
}

func processVolumeStore(op trace.Operation, volumeStore *models.VCHStorageVolumeStoresItems0) (string, *url.URL, error) {
	var url *url.URL

	label := volumeStore.Label
	rawTarget := volumeStore.Datastore

	err := CheckUnsupportedChars(label)
	if err != nil {
		return "", nil, err
	}

	// we strip the scheme off before parsing to url
	// this is to avoid url.Parse break when there is empty space in input path. Ex: ds://datastore [1]/foo/bar
	stripTarget := strings.Replace(rawTarget, DsScheme+"://", "", -1)
	url, err = url.Parse(stripTarget)
	if err != nil {
		return "", nil, fmt.Errorf("volume store path %s cannot be parsed into url: %s", rawTarget, err)
	}

	switch url.Scheme {
	case NfsScheme:
		// no further parsing needed for nfs target
	case EmptyScheme, DsScheme:
		if err := CheckUnsupportedCharsDatastore(rawTarget); err != nil {
			return "", nil, err
		}
		url.Scheme = DsScheme // add the scheme back on
		if len(url.RawQuery) > 0 {
			return "", nil, fmt.Errorf("volume store input must be in format %s or %s", dsInputFormat, nfsInputFormat)
		}
	default:
		return "", nil, fmt.Errorf("please specify a datastore or nfs target")
	}

	return label, url, nil
}
