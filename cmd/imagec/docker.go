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
	"archive/tar"
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

	ddigest "github.com/docker/distribution/digest"
	dlayer "github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/libtrust"

	"github.com/vmware/vic/pkg/trace"
)

// DigestSHA256EmptyTar is the canonical sha256 digest of empty tar file -
// (1024 NULL bytes)
const DigestSHA256EmptyTar = string(dlayer.DigestSHA256EmptyTar)

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
	Digest   string    `json:"digest,omitempty"`
	FSLayers []FSLayer `json:"fsLayers"`
	History  []History `json:"history"`
	// ignoring signatures
}

// LearnRegistryURL returns the registry URL after making sure that it responds to queries
func LearnRegistryURL(options ImageCOptions) (string, error) {
	defer trace.End(trace.Begin(options.registry))

	req := func(schema string) (string, error) {
		registry := fmt.Sprintf("%s://%s/v2/", schema, options.registry)

		url, err := url.Parse(registry)
		if err != nil {
			return "", err
		}
		log.Debugf("URL: %s", url)

		fetcher := NewURLFetcher(FetcherOptions{
			Timeout:            options.timeout,
			Username:           options.username,
			Password:           options.password,
			InsecureSkipVerify: options.insecureSkipVerify,
		})
		headers, err := fetcher.Head(url)
		if err != nil {
			return "", err
		}
		// v2 API requires this check
		if headers.Get("Docker-Distribution-API-Version") != "registry/2.0" {
			return "", fmt.Errorf("Missing Docker-Distribution-API-Version header")
		}
		return registry, nil
	}

	// first try https
	log.Debugf("Trying https scheme")
	registry, err := req("https")
	if err != nil && options.insecureAllowHTTP {
		// fallback to http if it's allowed
		log.Debugf("Falling back to http scheme")
		registry, err = req("http")
	}

	return registry, err
}

// LearnAuthURL returns the URL of the OAuth endpoint
func LearnAuthURL(options ImageCOptions) (*url.URL, error) {
	defer trace.End(trace.Begin(options.image + "/" + options.tag))

	url, err := url.Parse(options.registry)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, options.image, "manifests", options.tag)

	log.Debugf("URL: %s", url)

	fetcher := NewURLFetcher(FetcherOptions{
		Timeout:            options.timeout,
		Username:           options.username,
		Password:           options.password,
		InsecureSkipVerify: options.insecureSkipVerify,
	})
	// We expect docker registry to return a 401 to us - with a WWW-Authenticate header
	// We parse that header and learn the OAuth endpoint to fetch OAuth token.
	hdr, err := fetcher.Head(url)
	if err == nil && fetcher.IsStatusUnauthorized() {
		return fetcher.ExtractOAuthUrl(hdr.Get("www-authenticate"), url)
	}

	// Private registry returned the manifest directly as auth option is optional.
	// https://github.com/docker/distribution/blob/master/docs/configuration.md#auth
	if err == nil && options.registry != DefaultDockerURL && fetcher.IsStatusOK() {
		log.Debugf("%s does not support OAuth", url)
		return nil, nil
	}

	// Do we even have the image on that registry
	if err != nil && fetcher.IsStatusNotFound() {
		return nil, fmt.Errorf("%s:%s does not exist at %s", options.image, options.tag, options.registry)
	}

	return nil, fmt.Errorf("%s returned an unexpected response: %s", url, err)
}

// FetchToken fetches the OAuth token from OAuth endpoint
func FetchToken(url *url.URL) (*Token, error) {
	defer trace.End(trace.Begin(url.String()))

	log.Debugf("URL: %s", url)

	fetcher := NewURLFetcher(FetcherOptions{
		Timeout:            options.timeout,
		Username:           options.username,
		Password:           options.password,
		InsecureSkipVerify: options.insecureSkipVerify,
	})
	tokenFileName, err := fetcher.Fetch(url)
	if err != nil {
		err := fmt.Errorf("FetchToken (%s) failed: %s", url, err)
		log.Error(err)
		return nil, err
	}

	// Clenaup function
	defer func() {
		os.Remove(tokenFileName)
	}()

	// Read the file content into []byte for json.Unmarshal
	content, err := ioutil.ReadFile(tokenFileName)
	if err != nil {
		return nil, err
	}

	token := &Token{}

	err = json.Unmarshal(content, &token)
	if err != nil {
		return nil, err
	}

	if token.ExpiresIn == 0 {
		token.Expires = time.Now().Add(DefaultTokenExpirationDuration)
	} else {
		token.Expires = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	}

	return token, nil
}

// FetchImageBlob fetches the image blob
func FetchImageBlob(options ImageCOptions, image *ImageWithMeta) (string, error) {
	defer trace.End(trace.Begin(options.image + "/" + image.layer.BlobSum))

	id := image.ID
	layer := image.layer.BlobSum
	meta := image.meta
	diffID := ""

	url, err := url.Parse(options.registry)
	if err != nil {
		return diffID, err
	}
	url.Path = path.Join(url.Path, options.image, "blobs", layer)

	log.Debugf("URL: %s\n ", url)

	progress.Update(po, image.String(), "Pulling fs layer")

	fetcher := NewURLFetcher(FetcherOptions{
		Timeout:            options.timeout,
		Username:           options.username,
		Password:           options.password,
		Token:              options.token,
		InsecureSkipVerify: options.insecureSkipVerify,
	})
	imageFileName, err := fetcher.Fetch(url, image.String())
	if err != nil {
		return diffID, err
	}

	// Cleanup function for the error case
	defer func() {
		if err != nil {
			os.Remove(imageFileName)
		}
	}()

	// Open the file so that we can use it as a io.Reader for sha256 calculation
	imageFile, err := os.Open(string(imageFileName))
	if err != nil {
		return diffID, err
	}
	defer imageFile.Close()

	// blobSum is the sha of the compressed layer
	blobSum := sha256.New()

	// diffIDSum is the sha of the uncompressed layer
	diffIDSum := sha256.New()

	// blobTr is an io.TeeReader that writes bytes to blobSum that it reads from imageFile
	// see https://golang.org/pkg/io/#TeeReader
	blobTr := io.TeeReader(imageFile, blobSum)

	progress.Update(po, image.String(), "Verifying Checksum")
	decompressedTar, err := archive.DecompressStream(blobTr)
	if err != nil {
		return diffID, err
	}

	// Copy bytes from decompressed layer into diffIDSum to calculate diffID
	_, cerr := io.Copy(diffIDSum, decompressedTar)
	if cerr != nil {
		return diffID, cerr
	}

	bs := fmt.Sprintf("sha256:%x", blobSum.Sum(nil))
	if bs != layer {
		return diffID, fmt.Errorf("Failed to validate layer checksum. Expected %s got %s", layer, bs)
	}

	diffID = fmt.Sprintf("sha256:%x", diffIDSum.Sum(nil))

	// this isn't an empty layer, so we need to calculate the size
	if diffID != string(DigestSHA256EmptyTar) {
		var layerSize int64

		// seek to the beginning of the file
		imageFile.Seek(0, 0)

		// recreate the decompressed tar Reader
		decompressedTar, err := archive.DecompressStream(imageFile)
		if err != nil {
			return "", err
		}

		// get a tar reader for access to the files in the archive
		tr := tar.NewReader(decompressedTar)

		// iterate through tar headers to get file sizes
		for {
			tarHeader, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", err
			}
			layerSize += tarHeader.Size
		}

		image.size = layerSize
	}

	log.Infof("diffID for layer %s: %s", id, diffID)

	// Ensure the parent directory exists
	destination := path.Join(DestinationDirectory(), id)
	err = os.MkdirAll(destination, 0755) /* #nosec */
	if err != nil {
		return diffID, err
	}

	// Move(rename) the temporary file to its final destination
	err = os.Rename(string(imageFileName), path.Join(destination, id+".tar"))
	if err != nil {
		return diffID, err
	}

	// Dump the history next to it
	err = ioutil.WriteFile(path.Join(destination, id+".json"), []byte(meta), 0644)
	if err != nil {
		return diffID, err
	}

	progress.Update(po, image.String(), "Download complete")

	return diffID, nil
}

// FetchImageManifest fetches the image manifest file
func FetchImageManifest(options ImageCOptions) (*Manifest, error) {
	defer trace.End(trace.Begin(options.image + "/" + options.tag))

	url, err := url.Parse(options.registry)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, options.image, "manifests", options.tag)

	log.Debugf("URL: %s", url)

	fetcher := NewURLFetcher(FetcherOptions{
		Timeout:            10 * time.Second,
		Username:           options.username,
		Password:           options.password,
		Token:              options.token,
		InsecureSkipVerify: options.insecureSkipVerify,
	})
	manifestFileName, err := fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}

	// Cleanup function for the error case
	defer func() {
		if err != nil {
			os.Remove(manifestFileName)
		}
	}()

	// Read the entire file into []byte for json.Unmarshal
	content, err := ioutil.ReadFile(manifestFileName)
	if err != nil {
		return nil, err
	}

	manifest := &Manifest{}

	err = json.Unmarshal(content, manifest)
	if err != nil {
		return nil, err
	}

	if manifest.Name != options.image {
		return nil, fmt.Errorf("name doesn't match what was requested, expected: %s, downloaded: %s", options.image, manifest.Name)
	}

	if manifest.Tag != options.tag {
		return nil, fmt.Errorf("tag doesn't match what was requested, expected: %s, downloaded: %s", options.tag, manifest.Tag)
	}

	digest, err := getManifestDigest(content)
	if err != nil {
		return nil, err
	}

	manifest.Digest = digest

	// Ensure the parent directory exists
	destination := DestinationDirectory()
	err = os.MkdirAll(destination, 0755) /* #nosec */
	if err != nil {
		return nil, err
	}

	// Move(rename) the temporary file to its final destination
	err = os.Rename(string(manifestFileName), path.Join(destination, "manifest.json"))
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func getManifestDigest(content []byte) (string, error) {
	jsonSig, err := libtrust.ParsePrettySignature(content, "signatures")
	if err != nil {
		return "", err
	}

	// Resolve the payload in the manifest.
	bytes, err := jsonSig.Payload()
	if err != nil {
		return "", err
	}

	log.Debugf("Canonical Bytes: %d", len(bytes))
	digest := ddigest.FromBytes(bytes)
	// Correct Manifest Digest
	log.Debugf("Manifest Digest: %v", digest)
	return string(digest), nil
}
