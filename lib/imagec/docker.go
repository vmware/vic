// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

package imagec

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	ddigest "github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	dlayer "github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/reference"
	"github.com/docker/libtrust"

	"github.com/docker/docker/pkg/ioutils"
	urlfetcher "github.com/vmware/vic/pkg/fetcher"
	registryutils "github.com/vmware/vic/pkg/registry"
	"github.com/vmware/vic/pkg/trace"
)

const (
	// DigestSHA256EmptyTar is the canonical sha256 digest of empty tar file -
	// (1024 NULL bytes)
	DigestSHA256EmptyTar = string(dlayer.DigestSHA256EmptyTar)

	MaxMountAttempts = 4
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
	Digest   string    `json:"digest,omitempty"`
	FSLayers []FSLayer `json:"fsLayers"`
	History  []History `json:"history"`
	// ignoring signatures
}

// LearnRegistryURL returns the registry URL after making sure that it responds to queries
func LearnRegistryURL(options *Options) (string, error) {
	defer trace.End(trace.Begin(options.Registry))

	log.Debugf("Trying https scheme for %#v", options)

	registry, err := registryutils.Reachable(options.Registry, "https", options.Username, options.Password, options.RegistryCAs, options.Timeout, options.InsecureSkipVerify)

	if err != nil && options.InsecureAllowHTTP {
		// try https without verification
		log.Debugf("Trying https without verification, last error: %+v", err)
		registry, err = registryutils.Reachable(options.Registry, "https", options.Username, options.Password, options.RegistryCAs, options.Timeout, true)
		if err == nil {
			// Success, set InsecureSkipVerify to true
			options.InsecureSkipVerify = true
		} else {
			// try http
			log.Debugf("Falling back to http")
			registry, err = registryutils.Reachable(options.Registry, "http", options.Username, options.Password, options.RegistryCAs, options.Timeout, options.InsecureSkipVerify)
		}
	}

	return registry, err
}

// LearnAuthURL returns the URL of the OAuth endpoint
func LearnAuthURL(options Options) (*url.URL, error) {
	defer trace.End(trace.Begin(options.Reference.String()))

	url, err := url.Parse(options.Registry)
	if err != nil {
		return nil, err
	}

	tagOrDigest := tagOrDigest(options.Reference, options.Tag)
	url.Path = path.Join(url.Path, options.Image, "manifests", tagOrDigest)
	log.Debugf("URL: %s", url)

	fetcher := urlfetcher.NewURLFetcher(urlfetcher.Options{
		Timeout:            options.Timeout,
		Username:           options.Username,
		Password:           options.Password,
		InsecureSkipVerify: options.InsecureSkipVerify,
		RootCAs:            options.RegistryCAs,
	})

	// We expect docker registry to return a 401 to us - with a WWW-Authenticate header
	// We parse that header and learn the OAuth endpoint to fetch OAuth token.
	hdr, err := fetcher.Head(url)
	if err == nil && fetcher.IsStatusUnauthorized() {
		return fetcher.ExtractOAuthURL(hdr.Get("www-authenticate"), url)
	}

	// Private registry returned the manifest directly as auth option is optional.
	// https://github.com/docker/distribution/blob/master/docs/configuration.md#auth
	if err == nil && options.Registry != DefaultDockerURL && fetcher.IsStatusOK() {
		log.Debugf("%s does not support OAuth", url)
		return nil, nil
	}

	// Do we even have the image on that registry
	if err != nil && fetcher.IsStatusNotFound() {
		err = fmt.Errorf("image not found")
		return nil, urlfetcher.ImageNotFoundError{Err: err}
	}

	return nil, fmt.Errorf("%s returned an unexpected response: %s", url, err)
}

// FetchToken fetches the OAuth token from OAuth endpoint
func FetchToken(ctx context.Context, options Options, url *url.URL, progressOutput progress.Output) (*urlfetcher.Token, error) {
	defer trace.End(trace.Begin(url.String()))

	log.Debugf("URL: %s", url)

	fetcher := urlfetcher.NewURLFetcher(urlfetcher.Options{
		Timeout:            options.Timeout,
		Username:           options.Username,
		Password:           options.Password,
		InsecureSkipVerify: options.InsecureSkipVerify,
		RootCAs:            options.RegistryCAs,
	})

	token, err := fetcher.FetchAuthToken(url)
	if err != nil {
		err := fmt.Errorf("FetchToken (%s) failed: %s", url, err)
		log.Error(err)
		return nil, err
	}

	return token, nil
}

// FetchImageBlob fetches the image blob
func FetchImageBlob(ctx context.Context, options Options, image *ImageWithMeta, progressOutput progress.Output) (string, error) {
	defer trace.End(trace.Begin(options.Image + "/" + image.Layer.BlobSum))

	id := image.ID
	layer := image.Layer.BlobSum
	meta := image.Meta
	diffID := ""

	url, err := url.Parse(options.Registry)
	if err != nil {
		return diffID, err
	}
	url.Path = path.Join(url.Path, options.Image, "blobs", layer)

	log.Debugf("URL: %s\n ", url)

	fetcher := urlfetcher.NewURLFetcher(urlfetcher.Options{
		Timeout:            options.Timeout,
		Username:           options.Username,
		Password:           options.Password,
		Token:              options.Token,
		InsecureSkipVerify: options.InsecureSkipVerify,
		RootCAs:            options.RegistryCAs,
	})

	// ctx
	ctx, cancel := context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	imageFileName, err := fetcher.Fetch(ctx, url, nil, true, progressOutput, image.String())
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

	progress.Update(progressOutput, image.String(), "Verifying Checksum")
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

		image.Size = layerSize
	}

	log.Infof("diffID for layer %s: %s", id, diffID)

	// Ensure the parent directory exists
	destination := path.Join(DestinationDirectory(options), id)
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

	progress.Update(progressOutput, image.String(), "Download complete")

	log.Infof("Layer %s blobsum, diff id = (%s, %s)", id, layer, diffID)

	return diffID, nil
}

// tagOrDigest returns an image's digest if it's pulled by digest, or its tag
// otherwise.
func tagOrDigest(r reference.Named, tag string) string {
	if digested, ok := r.(reference.Canonical); ok {
		return digested.Digest().String()
	}

	return tag
}

// FetchImageManifest fetches the image manifest file
func FetchImageManifest(ctx context.Context, options Options, schemaVersion int, progressOutput progress.Output) (interface{}, string, error) {
	defer trace.End(trace.Begin(options.Reference.String()))

	if schemaVersion != 1 && schemaVersion != 2 {
		return nil, "", fmt.Errorf("Unknown schema version %d requested!", schemaVersion)
	}

	url, err := url.Parse(options.Registry)
	if err != nil {
		return nil, "", err
	}

	tagOrDigest := tagOrDigest(options.Reference, options.Tag)
	url.Path = path.Join(url.Path, options.Image, "manifests", tagOrDigest)
	log.Debugf("URL: %s", url)

	fetcher := urlfetcher.NewURLFetcher(urlfetcher.Options{
		Timeout:            options.Timeout,
		Username:           options.Username,
		Password:           options.Password,
		Token:              options.Token,
		InsecureSkipVerify: options.InsecureSkipVerify,
		RootCAs:            options.RegistryCAs,
	})

	reqHeaders := make(http.Header)
	if schemaVersion == 2 {
		reqHeaders.Add("Accept", schema2.MediaTypeManifest)
		reqHeaders.Add("Accept", schema1.MediaTypeManifest)
	}

	manifestFileName, err := fetcher.Fetch(ctx, url, &reqHeaders, true, progressOutput)
	if err != nil {
		return nil, "", err
	}

	// Cleanup function for the error case
	defer func() {
		if err != nil {
			os.Remove(manifestFileName)
		}
	}()

	switch schemaVersion {
	case 1: //schema 1, signed manifest
		return decodeManifestSchema1(manifestFileName, options, url.Hostname())
	case 2: //schema 2
		return decodeManifestSchema2(manifestFileName, options)
	}

	//We shouldn't really get here
	return nil, "", fmt.Errorf("Unknown schema version %d requested!", schemaVersion)
}

// decodeManifestSchema1() reads a manifest schema 1 and creates an imageC
// defined Manifest structure and returns the digest of the manifest as a string.
// For historical reason, we did not use the Docker's defined schema1.Manifest
// instead of our own and probably should do so in the future.
func decodeManifestSchema1(filename string, options Options, registry string) (interface{}, string, error) {
	// Read the entire file into []byte for json.Unmarshal
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}

	manifest := &Manifest{}
	err = json.Unmarshal(content, manifest)
	if err != nil {
		return nil, "", err
	}

	digest, err := getManifestDigest(content, options.Reference)
	if err != nil {
		return nil, "", err
	}

	manifest.Digest = digest

	// Verify schema 1 manifest's fields per docker/docker/distribution/pull_v2.go
	numFSLayers := len(manifest.FSLayers)
	if numFSLayers == 0 {
		return nil, "", fmt.Errorf("no FSLayers in manifest")
	}
	if numFSLayers != len(manifest.History) {
		return nil, "", fmt.Errorf("length of history not equal to number of layers")
	}

	return manifest, digest, nil
}

// verifyManifestDigest checks the manifest digest against the received payload.
func verifyManifestDigest(digested reference.Canonical, bytes []byte) error {
	verifier, err := ddigest.NewDigestVerifier(digested.Digest())
	if err != nil {
		return err
	}
	if _, err = verifier.Write(bytes); err != nil {
		return err
	}
	if !verifier.Verified() {
		return fmt.Errorf("image manifest verification failed for digest %s", digested.Digest())
	}

	return nil
}

// decodeManifestSchema2() reads a manifest schema 2 and creates a Docker
// defined Manifest structure and returns the digest of the manifest as a string.
func decodeManifestSchema2(filename string, options Options) (interface{}, string, error) {
	// Read the entire file into []byte for json.Unmarshal
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}

	manifest := &schema2.DeserializedManifest{}

	err = json.Unmarshal(content, manifest)
	if err != nil {
		return nil, "", err
	}

	_, canonical, err := manifest.Payload()
	if err != nil {
		return nil, "", err
	}

	digest := ddigest.FromBytes(canonical)

	return manifest, string(digest), nil
}

func getManifestDigest(content []byte, ref reference.Named) (string, error) {
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

	// Verify the manifest digest if the image is pulled by digest. If the image
	// is not pulled by digest, we proceed without this check because we don't
	// have a digest to verify the received content with.
	// https://docs.docker.com/registry/spec/api/#content-digests
	if digested, ok := ref.(reference.Canonical); ok {
		if err := verifyManifestDigest(digested, bytes); err != nil {
			return "", err
		}
	}

	digest := ddigest.FromBytes(bytes)
	// Correct Manifest Digest
	log.Debugf("Manifest Digest: %v", digest)
	return string(digest), nil
}
func PushImageBlob(ctx context.Context, options Options, as *ArchiveStream, layerReader io.ReadCloser, po progress.Output) (err error) {
	defer trace.End(trace.Begin(options.Image))

	layerID := ShortID(as.layerID)

	// The workflow: (https://docs.docker.com/registry/spec/api/#pushing-an-image)
	// 1. check if layer already exists
	// 2. try cross-repo-blob-mount
	// 3. if layer doesn't exist and cross-repo-blob-mount fails, obtain an upload url from the registry
	// 4. upload the layer to the registry
	// 5. if upload fails, cancel the upload process
	// TODO: skip foreign layers before checking layer existence

	log.Debugf("The registry in use is: %s", options.Registry)
	registryURL, err := url.Parse(options.Registry)
	if err != nil {
		return err
	}

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            options.Timeout,
		Username:           options.Username,
		Password:           options.Password,
		Token:              options.Token,
		InsecureSkipVerify: options.InsecureSkipVerify,
		RootCAs:            options.RegistryCAs,
	})

	pushDigest := as.digest

	log.Debugf("Checking for presence of layer %s (%s) in %s", layerID, pushDigest, options.Image)
	exist, err := CheckLayerExistence(ctx, transporter, options.Image, pushDigest, registryURL, po)
	if err != nil {
		return fmt.Errorf("failed to check for presence of layer %s (%s) in %s: %s", layerID, pushDigest, options.Image, err)
	}
	if exist {
		progress.Update(po, layerID, "Layer already exists")
		return nil
	}

	repoList, err := ObtainSourceRepoList(as.layerID, options.Reference)
	if err != nil {
		log.Errorf("Failed to fetch repo list: %s", err)
	} else {
		if repoList != nil && len(repoList) > 0 {
			mounted, err := CrossRepoBlobMount(ctx, layerID, pushDigest, options, repoList, po)
			if err != nil {
				return fmt.Errorf("failed during CrossRepoBlobMount: %s", err)
			}
			if mounted {
				return nil
			}
		}
	}

	// obtain upload url to start upload process
	// which would require obtaining a list of the repositories that the current user has access to.
	// See https://docs.docker.com/registry/spec/api/#pushing-an-image
	uploadURL, err := ObtainUploadURL(ctx, transporter, registryURL, options.Image, po)
	if err != nil {
		return err
	}
	log.Infof("The upload url is: %s", uploadURL)

	var reader io.ReadCloser

	reader = progress.NewProgressReader(ioutils.NewCancelReadCloser(ctx, layerReader), po, as.size, layerID, "Pushing")

	if err = UploadLayer(ctx, transporter, pushDigest, uploadURL, reader, po, layerID); err != nil {
		if err2 := CancelUpload(ctx, transporter, uploadURL, po); err2 != nil {
			log.Errorf("Failed during CancelUpload: %s", err2)
		}
		return err
	}
	progress.Update(po, layerID, "Pushed")

	return nil
}

func ObtainSourceRepoList(layerID string, targetRepo reference.Named) ([]string, error) {
	defer trace.End(trace.Begin(layerID))

	layer, err := LayerCache().Get(layerID)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain source repo list: %s", err)
	}
	if layer.V2Meta == nil || len(layer.V2Meta) == 0 {
		log.Debugf("layer.V2Meta does not exist or is empty")
		return nil, nil
	}

	var repoList []string
	var repo string

	for _, meta := range layer.V2Meta {
		repo = meta.SourceRepository

		// Do not consider repos not in the same registry
		sourceRepo, err := reference.ParseNamed(meta.SourceRepository)
		if err != nil || targetRepo.Hostname() != sourceRepo.Hostname() {
			continue
		}
		// Do not consider the target repository
		if repo == targetRepo.FullName() {
			continue
		}
		repoList = append(repoList, sourceRepo.RemoteName())
	}

	if repoList != nil && len(repoList) > MaxMountAttempts {
		repoList = append(repoList[:MaxMountAttempts-1])
	}

	if repoList != nil {
		log.Debugf("RepoList: %+v", repoList)
	}

	return repoList, nil
}

func CrossRepoBlobMount(ctx context.Context, layerID, digest string, options Options, repoList []string, po progress.Output) (bool, error) {
	defer trace.End(trace.Begin(options.Image))

	var (
		mounted  bool
		authURL  *url.URL
		newToken *urlfetcher.Token
	)

	registry, err := url.Parse(options.Registry)
	if err != nil {
		return false, err
	}

	for _, repo := range repoList {
		log.Debugf("Attempting to mount layer %s (%s) from %s", layerID, digest, repo)

		authURL, err = LearnAuthURLForBlobMount(options, digest, repo, po)
		if err != nil {
			return false, err
		}

		// Get the OAuth token - if only we have a URL
		if authURL != nil {
			newToken, err = FetchToken(ctx, options, authURL, po)
			if err != nil {
				log.Errorf("Failed to fetch OAuth token: %s", err)
				return false, err
			}
		}

		newTransporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
			Timeout:            options.Timeout,
			Username:           options.Username,
			Password:           options.Password,
			Token:              newToken,
			InsecureSkipVerify: options.InsecureSkipVerify,
			RootCAs:            options.RegistryCAs,
		})

		// if mount fails, the registry will fall back to the standard upload behavior
		// and return a 202 Accepted with the upload URL in the Location header
		mounted, _, err = MountBlobToRepo(ctx, newTransporter, registry, digest, options.Image, repo, po)
		if err != nil {
			log.Errorf("Mount layer to repo %s failed: %s", repo, err)
		}
		if mounted {
			progress.Updatef(po, layerID, "Mounted from %s", repo)
			log.Debugf("Layer %s mounted from %s", layerID, repo)
			break
		}
	}

	if mounted {
		return true, nil
	}

	log.Debugf("Layer cannot be mounted")
	return false, nil
}

// LearnAuthURLForPush returns the URL of the OAuth endpoint
func LearnAuthURLForPush(options Options, po progress.Output) (*url.URL, error) {
	defer trace.End(trace.Begin(options.Reference.String()))

	url, err := url.Parse(options.Registry)
	if err != nil {
		return nil, err
	}

	url.Path = path.Join(url.Path, options.Image, "blobs", "uploads")
	url.Path += "/"
	log.Debugf("LearnAuthURLForPush URL: %s", url)

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            options.Timeout,
		Username:           options.Username,
		Password:           options.Password,
		InsecureSkipVerify: options.InsecureSkipVerify,
		RootCAs:            options.RegistryCAs,
	})

	// We expect docker registry to return a 401 to us - with a WWW-Authenticate header
	// We parse that header and learn the OAuth endpoint to fetch OAuth token.
	hdr, err := transporter.Post(ctx, url, bytes.NewReader([]byte("")), nil, po)

	if err != nil {
		return nil, err
	}

	if err == nil && transporter.IsStatusUnauthorized() {
		return transporter.ExtractOAuthURL(hdr.Get("www-authenticate"), url)
	}

	return nil, fmt.Errorf("Unexpected http code: %d, URL: %s", transporter.Status(), url)
}

func LearnAuthURLForBlobMount(options Options, digest, repo string, po progress.Output) (*url.URL, error) {
	defer trace.End(trace.Begin(options.Reference.String()))

	composedURL, _ := url.Parse(options.Registry)
	composedURL.Path = path.Join(composedURL.Path, options.Image, "blobs/uploads")
	composedURL.Path += "/"

	q := composedURL.Query()
	q.Add("mount", digest)
	q.Add("from", repo)
	composedURL.RawQuery = q.Encode()

	log.Debugf("The url for MountBlobToRepo is: %s", composedURL)

	// POST /v2/<name>/blobs/uploads/?mount=<digest>&from=<repository name>
	// Content-Length: 0
	reqHdrs := &http.Header{
		"Content-Length": {"0"},
	}

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            options.Timeout,
		Username:           options.Username,
		Password:           options.Password,
		InsecureSkipVerify: options.InsecureSkipVerify,
		RootCAs:            options.RegistryCAs,
	})

	hdr, err := transporter.Post(ctx, composedURL, bytes.NewReader([]byte("")), reqHdrs, po)

	if err != nil {
		return nil, err
	}

	if transporter.IsStatusUnauthorized() {
		authURL, err := transporter.ExtractOAuthURL(hdr.Get("www-authenticate"), composedURL)
		if err != nil {
			return nil, err
		}

		log.Debugf("The url for authenticating CrossRepoBlobMount is: %+v", authURL)
		return authURL, nil
	}

	return nil, fmt.Errorf("Unexpected http code: %d, URL: %s", transporter.Status(), composedURL)
}

// PutImageManifest simply pushes the manifest up to the registry.
func PutImageManifest(ctx context.Context, pusher Pusher, options Options, progressOutput progress.Output) error {
	defer trace.End(trace.Begin(""))

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            options.Timeout,
		Username:           options.Username,
		Password:           options.Password,
		InsecureSkipVerify: options.InsecureSkipVerify,
		RootCAs:            options.RegistryCAs,
		Token:              options.Token,
	})

	// Create manifest push URL
	url, err := url.Parse(options.Registry)
	if err != nil {
		return err
	}

	// check configJSON existence
	exist, err := CheckLayerExistence(ctx, transporter, options.Image, pusher.PushManifest.Config.Digest.String(), url, progressOutput)
	if err != nil {
		return fmt.Errorf("failed to check configJSON existence: %s", err)
	}

	if !exist {
		// obtain uploadURL for configJSON
		uploadURL, err := ObtainUploadURL(ctx, transporter, url, options.Image, progressOutput)
		if err != nil {
			return fmt.Errorf("failed to obtain uploadURL for configJSON: %s", err)
		}

		// upload configJSON
		configReader := bytes.NewReader(pusher.configJSON)
		err = UploadLayer(ctx, transporter, pusher.PushManifest.Config.Digest.String(), uploadURL, configReader, progressOutput)
		if err != nil {
			return fmt.Errorf("failed to upload configJSON: %s", err)
		}
	}

	// upload manifest
	tagOrDigest := tagOrDigest(options.Reference, options.Tag)
	url.Path = path.Join(url.Path, options.Image, "manifests", tagOrDigest)
	log.Debugf("URL for PutIamgeManifest: %s", url)

	// Add content type headers
	reqHeaders := make(http.Header)
	var dataReader io.Reader

	reqHeaders.Add("Content-Type", schema2.MediaTypeManifest)
	dataReader, dmanifest, err := getManifestSchema2Reader(pusher.PushManifest)
	if err != nil {
		log.Errorf("Failed to read manifest schema 2: %s", err.Error())
	}

	_, err = transporter.Put(ctx, url, dataReader, &reqHeaders, progressOutput)
	if err != nil {
		return fmt.Errorf("failed to upload image manifest: %s", err)
	}

	if transporter.IsStatusCreated() {
		log.Infof("The manifest is uploaded successfully")

		_, canonicalManifest, err := dmanifest.Payload()
		if err != nil {
			return err
		}

		manifestDigest := ddigest.FromBytes(canonicalManifest)
		msg := fmt.Sprintf("Digest: %s size: %d", manifestDigest, len(canonicalManifest))
		progress.Message(progressOutput, tagOrDigest, msg)
		return nil
	}

	return fmt.Errorf("unexpected http code during PushManifest: %d, URL: %s", transporter.StatusCode, url)
}

// Upload the layer (monolithic upload)
func UploadLayer(ctx context.Context, u *urlfetcher.URLTransporter, digest, uploadURL string, layer io.Reader, po progress.Output, ids ...string) error {
	defer trace.End(trace.Begin(uploadURL))

	// PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	// The uuid is from the `location` header
	// in the response of the first step if successful
	composedURL, err := url.Parse(uploadURL)
	if err != nil {
		return fmt.Errorf("failed to parse uploadURL: %s", err)
	}
	q := composedURL.Query()
	q.Add("digest", digest)
	composedURL.RawQuery = q.Encode()

	log.Debugf("The url for UploadLayer is: %s", composedURL)

	reqHdrs := &http.Header{
		"Content-Type": {"application/octet-stream"},
	}

	id := ""
	if len(ids) > 0 {
		id = ids[0]
	}
	_, err = u.Put(ctx, composedURL, layer, reqHdrs, po, id)
	if err != nil {
		return fmt.Errorf("failed to upload layer: %s", err)
	}

	if u.IsStatusCreated() {
		log.Infof("The uploadLayer finished successfully")
		return nil
	}

	return fmt.Errorf("unexpected http code during UploadLayer: %d, URL: %s", u.StatusCode, composedURL)
}

// Cancel the upload process
func CancelUpload(ctx context.Context, u *urlfetcher.URLTransporter, uploadURL string, po progress.Output) error {
	defer trace.End(trace.Begin(uploadURL))

	// DELETE /v2/<name>/blobs/uploads/<uuid>
	composedURL, err := url.Parse(uploadURL)
	if err != nil {
		return fmt.Errorf("failed to parse uploadURL: %s", err)
	}

	log.Debugf("The url for CancelUpload is: %s\n ", composedURL)

	_, err = u.Delete(ctx, composedURL, nil, po)
	if err != nil {
		return fmt.Errorf("failed to cancel upload: %s", err)
	}

	if u.IsStatusNoContent() {
		log.Infof("The upload process is cancelled successfully")
		return nil
	}

	return fmt.Errorf("unexpected http code during CancelUpload: %d, URL: %s", u.StatusCode, composedURL)
}

// Notify the registry that the upload process is completed
// Currently this is not used since we only use monolithic upload
// However, if the image layer is too large, chunk upload has to be implemented and this method should be called to complete the process
func CompletedUpload(ctx context.Context, u *urlfetcher.URLTransporter, digest, uploadURL string, po progress.Output) error {
	defer trace.End(trace.Begin(uploadURL))

	// PUT /v2/<name>/blob/uploads/<uuid>?digest=<digest>
	composedURL, err := url.Parse(uploadURL)
	q := composedURL.Query()
	q.Add("digest", digest)
	composedURL.RawQuery = q.Encode()

	log.Debugf("The url for CompletedUpload is: %s", composedURL)

	reqHdrs := &http.Header{
		"Content-Length": {"0"},
		"Content-Type":   {"application/octet-stream"},
	}
	_, err = u.Put(ctx, composedURL, bytes.NewReader([]byte("")), reqHdrs, po)
	if err != nil {
		return fmt.Errorf("failed to complete upload: %s", err)
	}

	if u.IsStatusCreated() {
		log.Infof("The upload process completed successfully")
		return nil
	}

	return fmt.Errorf("unexpected http code during CompletedUpload: %d, URL: %s", u.StatusCode, composedURL)
}

// Check if a layer exists
func CheckLayerExistence(ctx context.Context, u *urlfetcher.URLTransporter, image, digest string, registry *url.URL, po progress.Output) (bool, error) {
	defer trace.End(trace.Begin(digest))

	// HEAD /v2/<name>/blobs/<digest>
	composedURL := urlfetcher.URLDeepCopy(registry)
	composedURL.Path = path.Join(registry.Path, image, "blobs", digest)

	log.Debugf("The url for checking layer existence is: %s", composedURL)

	_, err := u.Head(ctx, composedURL, nil, po)
	if err != nil {
		return false, fmt.Errorf("failed to check layer existence: %s", err)
	}

	if u.IsStatusOK() {
		log.Debugf("The layer already exists")
		return true, nil
	}

	if u.IsStatusNotFound() {
		log.Infof("The layer does not exist")
		return false, nil
	}
	return false, fmt.Errorf("unexpected http code during CheckLayerExistence: %d, URL: %s", u.StatusCode, composedURL)
}

// obtain the upload url
func ObtainUploadURL(ctx context.Context, u *urlfetcher.URLTransporter, registry *url.URL, image string, po progress.Output) (string, error) {
	defer trace.End(trace.Begin(image))

	// POST /v2/<name>/blobs/uploads
	composedURL := urlfetcher.URLDeepCopy(registry)
	composedURL.Path = path.Join(registry.Path, image, "blobs/uploads/")
	composedURL.Path += "/"

	log.Debugf("The url for ObtainUploadURL is: %s", composedURL)

	hdr, err := u.Post(ctx, composedURL, nil, nil, po)

	if err != nil {
		return "", err
	}

	// even if the image does not exist (push a new image), we should still be able to get a location for upload
	if u.IsStatusAccepted() {
		log.Debugf("The location is: %s", hdr.Get("Location"))
		return hdr.Get("Location"), nil
	}

	return "", fmt.Errorf("unexpected http code during ObtainUploadURL: %d, URL: %s", u.StatusCode, composedURL)
}

func MountBlobToRepo(ctx context.Context, u *urlfetcher.URLTransporter, registry *url.URL, digest, image, repo string, po progress.Output) (bool, string, error) {
	defer trace.End(trace.Begin("image: " + image + ", repo: " + repo))

	// POST /v2/<name>/blobs/uploads/?mount=<digest>&from=<repository name>
	// Content-Length: 0
	composedURL := urlfetcher.URLDeepCopy(registry)
	composedURL.Path = path.Join(registry.Path, image, "blobs/uploads")
	composedURL.Path += "/"

	q := composedURL.Query()
	q.Add("mount", digest)
	q.Add("from", repo)
	composedURL.RawQuery = q.Encode()

	log.Debugf("The url for MountBlobToRepo is: %s\n ", composedURL)

	reqHdrs := &http.Header{
		"Content-Length": {"0"},
	}

	hdr, err := u.Post(ctx, composedURL, bytes.NewReader([]byte("")), reqHdrs, po)
	if err != nil {
		return false, "", fmt.Errorf("failed to mount blob to repo: %s", err)
	}

	if u.IsStatusCreated() {
		log.Infof("The blob is already mounted to the repo!")
		return true, "", nil
	}

	if u.IsStatusAccepted() {
		log.Infof("The blob is not mounted to repo '%s' yet", repo)
		log.Infof("The location is: %s", hdr.Get("Location"))
		return false, hdr.Get("Location:"), nil
	}

	return false, "", fmt.Errorf("unexpected http code during ObtainUploadURL: %d, URL: %s", u.StatusCode, composedURL)
}

func ShortID(id string) string {
	return stringid.TruncateID(id)
}

func getManifestSchema2Reader(manifest schema2.Manifest) (io.Reader, *schema2.DeserializedManifest, error) {
	log.Debugf("Constructing manifest schema2 reader")

	dm, err := schema2.FromStruct(manifest)
	if err != nil {
		msg := fmt.Sprintf("unable to construct DeserializedManifest: %s", err)
		log.Error(msg)
		return nil, nil, fmt.Errorf(msg)
	}

	data, err := dm.MarshalJSON()
	//data, err := json.Marshal(manifest)
	if err != nil {
		msg := fmt.Sprintf("unable to marshal DeserializedManifest: %s", err)
		log.Error(msg)
		return nil, nil, fmt.Errorf(msg)
	}

	return bytes.NewReader(data), dm, nil
}
