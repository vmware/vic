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

package main

import (
	"io"

	log "github.com/Sirupsen/logrus"

	"github.com/go-swagger/go-swagger/httpkit"
	httptransport "github.com/go-swagger/go-swagger/httpkit/client"

	apiclient "github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/misc"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/pkg/trace"
)

const (
	MetaDataKey = "metaData"
)

// PingPortLayer calls the _ping endpoint of the portlayer
func PingPortLayer() (bool, error) {
	defer trace.End(trace.Begin(options.host))

	transport := httptransport.New(options.host, "/", []string{"http"})
	client := apiclient.New(transport, nil)

	ok, err := client.Misc.Ping(misc.NewPingParams())
	if err != nil {
		return false, err
	}
	return ok.Payload == "OK", nil
}

// CreateImageStore creates an image store
func CreateImageStore(storename string) error {
	defer trace.End(trace.Begin(storename))

	transport := httptransport.New(options.host, "/", []string{"http"})
	client := apiclient.New(transport, nil)

	log.Debugf("Creating a store")

	body := &models.ImageStore{Name: storename}

	_, err := client.Storage.CreateImageStore(
		storage.NewCreateImageStoreParams().WithBody(body),
	)
	if _, ok := err.(*storage.CreateImageStoreConflict); ok {
		log.Debugf("Store already exists")
		return nil
	}
	if err != nil {
		log.Debugf("Creating a store failed: %s", err)

		return err
	}
	log.Debugf("Created a store %#v", body)

	return nil
}

// ListImages lists the images from given image store
func ListImages(storename string, images []*ImageWithMeta) (map[string]*models.Image, error) {
	defer trace.End(trace.Begin(storename))

	transport := httptransport.New(options.host, "/", []string{"http"})
	client := apiclient.New(transport, nil)

	ids := make([]string, len(images))

	for i := range images {
		ids = append(ids, images[i].ID)
	}

	imageList, err := client.Storage.ListImages(
		storage.NewListImagesParams().WithStoreName(storename).WithIds(ids),
	)
	if err != nil {
		return nil, err
	}

	existingImages := make(map[string]*models.Image)
	for i := range imageList.Payload {
		v := imageList.Payload[i]
		existingImages[v.ID] = v
	}
	return existingImages, nil
}

// WriteImage writes the image to given image store
func WriteImage(image *ImageWithMeta, data io.ReadCloser) error {
	defer trace.End(trace.Begin(image.ID))

	transport := httptransport.New(options.host, "/", []string{"http"})
	client := apiclient.New(transport, nil)

	transport.Consumers["application/json"] = httpkit.JSONConsumer()
	transport.Producers["application/json"] = httpkit.JSONProducer()
	transport.Consumers["application/octet-stream"] = httpkit.ByteStreamConsumer()
	transport.Producers["application/octet-stream"] = httpkit.ByteStreamProducer()

	key := new(string)
	blob := new(string)

	*key = MetaDataKey
	*blob = image.meta

	r, err := client.Storage.WriteImage(
		storage.NewWriteImageParams().
			WithImageID(image.ID).
			WithParentID(*image.Parent).
			WithStoreName(image.Store).
			WithMetadatakey(key).
			WithMetadataval(blob).
			WithImageFile(data).
			WithSum(image.layer.BlobSum),
	)
	if err != nil {
		log.Debugf("Creating an image failed: %s", err)
		return err
	}
	log.Printf("Created an image %#v", r.Payload)

	return nil

}
