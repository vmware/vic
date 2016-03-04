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
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/progress"

	"github.com/vmware/vic/pkg/trace"
)

// FSLayer is a container struct for BlobSums defined in an image manifest
type FSLayer struct {
	// BlobSum is the tarsum of the referenced filesystem image layer
	BlobSum string `json:"blobSum"`
}

// History is a container struct for V1Compatibility defined in an image manifest
type History struct {
	V1Compatibility string `json:"v1Compatibility"`
}

// Manifest represents the Docker Manifest file
type Manifest struct {
	Name     string    `json:"name"`
	Tag      string    `json:"tag"`
	FSLayers []FSLayer `json:"fsLayers"`
	History  []History `json:"history"`
	// ignoring signatures
}

// V1Compatibility represents some parts of V1Compatibility
type V1Compatibility struct {
	ID        string    `json:"id"`
	Parent    string    `json:"parent,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	Created   time.Time `json:"created"`
	Container string    `json:"container,omitempty"`
}

// LearnAuthURL returns the URL of the OAuth endpoint
func LearnAuthURL(options ImageCOptions) (*url.URL, error) {
	defer trace.End(trace.Begin(options.image + "/" + options.digest))

	url, err := url.Parse(options.registry)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, options.image, "manifests", options.digest)

	log.Debugf("URL: %s", url)

	fetcher := NewFetcher(FetcherOptions{
		Timeout:            options.timeout,
		Username:           options.username,
		Password:           options.password,
		InsecureSkipVerify: options.insecure,
	})
	// We expect docker registry to return a 401 to us - with a WWW-Authenticate header
	// We parse that header and learn the OAuth endpoint to fetch OAuth token.
	_, err = fetcher.Fetch(url)
	if err != nil && fetcher.IsStatusUnauthorized() {
		return fetcher.AuthURL(), nil
	}

	// Private registry returned the manifest directly as auth option is optional.
	// https://github.com/docker/distribution/blob/master/docs/configuration.md#auth
	if err == nil && options.registry != DefaultDockerURL && fetcher.IsStatusOK() {
		log.Debugf("%s does not support OAuth", url)
		return nil, nil
	}

	// Do we even have the image on that registry
	if err != nil && fetcher.IsStatusNotFound() {
		return nil, fmt.Errorf("%s:%s does not exists at %s", options.image, options.digest, options.registry)
	}

	return nil, fmt.Errorf("%s returned an unexpected response: %s", url, err)
}

// FetchToken fetches the OAuth token from OAuth endpoint
func FetchToken(url *url.URL) (*Token, error) {
	defer trace.End(trace.Begin(url.String()))

	log.Debugf("URL: %s", url)

	fetcher := NewFetcher(FetcherOptions{
		Timeout:            options.timeout,
		Username:           options.username,
		Password:           options.password,
		InsecureSkipVerify: options.insecure,
	})
	body, err := fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}

	token := &Token{}

	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	if token.Expires.IsZero() {
		token.Expires = time.Now().Add(DefaultTokenExpirationDuration)
	}

	return token, err
}

// FetchImageBlob fetches the image blob
func FetchImageBlob(options ImageCOptions, image *ImageWithMeta) error {
	defer trace.End(trace.Begin(options.image + "/" + image.layer.BlobSum))

	id := image.ID
	layer := image.layer.BlobSum
	history := image.history.V1Compatibility

	url, err := url.Parse(options.registry)
	if err != nil {
		return err
	}
	url.Path = path.Join(url.Path, options.image, "blobs", layer)

	log.Debugf("URL: %s\n ", url)

	progress.Update(po, id[:12], "Pulling fs layer")

	fetcher := NewFetcher(FetcherOptions{
		Timeout:            options.timeout,
		Username:           options.username,
		Password:           options.password,
		Token:              options.token,
		InsecureSkipVerify: options.insecure,
	})
	blob, err := fetcher.FetchWithProgress(url, id[:12])
	if err != nil {
		return err
	}

	in := bytes.NewReader(blob)

	destination := path.Join(DestinationDirectory(), id)
	err = os.MkdirAll(destination, 0755)
	if err != nil {
		return err
	}

	progress.Update(po, id[:12], "Verifying Checksum")

	h := sha256.New()
	t := io.TeeReader(in, h)

	layerFile, err := os.Create(path.Join(destination, id+".tar"))
	if err != nil {
		return err
	}
	defer layerFile.Close()

	_, err = io.Copy(layerFile, t)
	if err != nil {
		return err
	}

	if err := layerFile.Sync(); err != nil {
		return err
	}

	sum := fmt.Sprintf("sha256:%x", h.Sum(nil))
	if sum != layer {
		return fmt.Errorf("Failed to validate layer checksum. Expected %s got %s", layer, sum)
	}
	ioutil.WriteFile(path.Join(destination, id+".json"), []byte(history), 0644)

	progress.Update(po, id[:12], "Download complete")

	return nil
}

// FetchImageManifest fetches the image manifest file
func FetchImageManifest(options ImageCOptions) (*Manifest, error) {
	defer trace.End(trace.Begin(options.image + "/" + options.digest))

	url, err := url.Parse(options.registry)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, options.image, "manifests", options.digest)

	log.Debugf("URL: %s", url)

	fetcher := NewFetcher(FetcherOptions{
		Timeout:            10 * time.Second,
		Username:           options.username,
		Password:           options.password,
		Token:              options.token,
		InsecureSkipVerify: options.insecure,
	})

	blob, err := fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}

	manifest := &Manifest{}

	err = json.Unmarshal(blob, manifest)
	if err != nil {
		return nil, err
	}

	if manifest.Name != options.image {
		return nil, fmt.Errorf("name doesn't match what was requested, expected: %s, downloaded: %s", options.image, manifest.Name)
	}

	if manifest.Tag != options.digest {
		return nil, fmt.Errorf("tag doesn't match what was requested, expected: %s, downloaded: %s", options.digest, manifest.Tag)
	}

	destination := DestinationDirectory()
	err = os.MkdirAll(destination, 0755)
	if err != nil {
		return nil, err
	}
	ioutil.WriteFile(path.Join(destination, "manifest.json"), blob, 0644)

	return manifest, nil
}
