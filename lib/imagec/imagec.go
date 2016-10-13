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

package imagec

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	docker "github.com/docker/docker/image"
	dockerLayer "github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/reference"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/metadata"
	urlfetcher "github.com/vmware/vic/pkg/fetcher"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/sys"
)

// ImageC is responsible for pulling docker images from a repository
type ImageC struct {
	Options

	// https://raw.githubusercontent.com/docker/docker/master/distribution/pull_v2.go
	sf             *streamformatter.StreamFormatter
	progressOutput progress.Output
}

// NewImageC returns a new instance of ImageC
func NewImageC(options Options, strfmtr *streamformatter.StreamFormatter) *ImageC {
	return &ImageC{
		Options:        options,
		sf:             strfmtr,
		progressOutput: strfmtr.NewProgressOutput(options.Outstream, false),
	}
}

// Options contain all options for a single instance of imagec
type Options struct {
	Reference string

	Registry string
	Image    string
	Tag      string

	Destination string

	Host      string
	Storename string

	Logfile string

	Username string
	Password string

	Token *urlfetcher.Token

	Timeout time.Duration

	Stdout bool
	Debug  bool

	Outstream io.Writer

	InsecureSkipVerify bool
	InsecureAllowHTTP  bool

	Standalone bool

	Profiling string
	Tracing   bool

	ImageManifest *Manifest
}

// ImageWithMeta wraps the models.Image with some additional metadata
type ImageWithMeta struct {
	*models.Image

	diffID string
	layer  FSLayer
	meta   string
	size   int64
}

func (i *ImageWithMeta) String() string {
	return stringid.TruncateID(i.layer.BlobSum)
}

var (
	ldm *LayerDownloader
)

const (
	// DefaultDockerURL holds the URL of Docker registry
	DefaultDockerURL = "registry-1.docker.io"

	// DefaultDestination specifies the default directory to use
	DefaultDestination = "images"

	// DefaultPortLayerHost specifies the default port layer server
	DefaultPortLayerHost = "localhost:2377"

	// DefaultLogfile specifies the default log file name
	DefaultLogfile = "imagec.log"

	// DefaultHTTPTimeout specifies the default HTTP timeout
	DefaultHTTPTimeout = 3600 * time.Second

	// attribute update actions
	Add = iota + 1
	Remove
)

func init() {
	ldm = NewLayerDownloader()
}

// ParseReference parses the -reference parameter and populate options struct
func (ic *ImageC) ParseReference() error {
	// Validate and parse reference name
	ref, err := reference.ParseNamed(ic.Reference)
	if err != nil {
		log.Warn("Error while parsing reference %s: %#v", ic.Reference, err)
		return err
	}

	ic.Tag = reference.DefaultTag
	if !reference.IsNameOnly(ref) {
		if tagged, ok := ref.(reference.NamedTagged); ok {
			ic.Tag = tagged.Tag()
		}
	}

	ic.Registry = DefaultDockerURL
	if ref.Hostname() != reference.DefaultHostname {
		ic.Registry = ref.Hostname()
	}

	ic.Image = ref.RemoteName()

	return nil
}

// DestinationDirectory returns the path of the output directory
func DestinationDirectory(options Options) string {
	u, _ := url.Parse(options.Registry)

	// Use a hierarchy like following so that we can support multiple schemes, registries and versions
	/*
		https/
		├── 192.168.218.5:5000
		│   └── v2
		│       └── busybox
		│           └── latest
		...
		│               ├── fef924a0204a00b3ec67318e2ed337b189c99ea19e2bf10ed30a13b87c5e17ab
		│               │   ├── fef924a0204a00b3ec67318e2ed337b189c99ea19e2bf10ed30a13b87c5e17ab.json
		│               │   └── fef924a0204a00b3ec67318e2ed337b189c99ea19e2bf10ed30a13b87c5e17ab.tar
		│               └── manifest.json
		└── registry-1.docker.io
		    └── v2
		        └── library
		            └── golang
		                └── latest
		                    ...
		                    ├── f61ebe2817bb4e6a7f0a4cf249a5316223f7ecc886feac24b9887a490feaed57
		                    │   ├── f61ebe2817bb4e6a7f0a4cf249a5316223f7ecc886feac24b9887a490feaed57.json
		                    │   └── f61ebe2817bb4e6a7f0a4cf249a5316223f7ecc886feac24b9887a490feaed57.tar
		                    └── manifest.json

	*/
	return path.Join(
		options.Destination,
		u.Scheme,
		u.Host,
		u.Path,
		options.Image,
		options.Tag,
	)
}

// LayersToDownload creates a slice of ImageWithMeta for the layers that need to be downloaded
func (ic *ImageC) LayersToDownload() ([]*ImageWithMeta, error) {
	images := make([]*ImageWithMeta, len(ic.ImageManifest.FSLayers))

	manifest := ic.ImageManifest
	v1 := docker.V1Image{}

	// iterate from parent to children
	for i := len(ic.ImageManifest.History) - 1; i >= 0; i-- {
		history := manifest.History[i]
		layer := manifest.FSLayers[i]

		// unmarshall V1Compatibility to get the image ID
		if err := json.Unmarshal([]byte(history.V1Compatibility), &v1); err != nil {
			return nil, fmt.Errorf("Failed to unmarshall image history: %s", err)
		}

		// if parent is empty set it to scratch
		parent := "scratch"
		if v1.Parent != "" {
			parent = v1.Parent
		}

		// add image to ImageWithMeta list
		images[i] = &ImageWithMeta{
			Image: &models.Image{
				ID:     v1.ID,
				Parent: &parent,
				Store:  ic.Storename,
			},
			meta:   history.V1Compatibility,
			layer:  layer,
			diffID: "",
		}
	}

	return images, nil
}

// Will updated metadata based on the current image being requested
//
// Currently only tags are updated, but once imagec supports pulling by digest
// then digest will need to be updated as well
func updateImageMetadata(ic *ImageC, imageLayer *ImageWithMeta, manifest *Manifest) error {
	// if standalone no need to continue
	if ic.Standalone {
		return nil
	}
	// updated Image slice
	var updatedImages []*ImageWithMeta

	// assuming calls to the portlayer are cheap, lets get all existing images
	existingImages, err := ListImages(ic.Host, imageLayer.Store, nil)
	if err != nil {
		return fmt.Errorf("updateImageMetadata failed to obtain list of images: %s", err)
	}

	// iterate overall the images and remove the tag that was just downloaded
	for id := range existingImages {
		// utlizing index to avoid copy
		existingImage := existingImages[id]

		imageMeta := &metadata.ImageConfig{}
		// tag / digest info only resides in the metadata so must unmarshall meta to determine
		if err := json.Unmarshal([]byte(existingImage.Metadata[metadata.MetaDataKey]), imageMeta); err != nil {
			return fmt.Errorf("updateImageMetadata failed to get existing metadata: Layer(%s) %s", id, err)
		}

		// if this isn't an image we can skip..
		if imageMeta.ImageID == "" {
			continue
		}

		// is this the same repo (busybox, etc)
		if imageMeta.Name == ic.Image {

			// default action to remove and update if necessary
			action := Remove
			if id == imageLayer.ID {
				// add to the array if this layer is the requested layer
				action = Add
			}

			// TODO: imagec doesn't support pulling by digest - add digest support once
			// issue 1186 is implemented
			if newTags, ok := arrayUpdate(manifest.Tag, imageMeta.Tags, action); ok {
				imageMeta.Tags = newTags

				// We need the parent id this is the layer ID of the parent disk
				// Parent is returned from the PL as a full URL, so we only need the base
				layerParent := path.Base(*existingImage.Parent)

				// marshall back to dumb byte array
				meta, err := json.Marshal(imageMeta)
				if err != nil {
					return fmt.Errorf("updateImageMetadata unable to marshal modified metadata: %s", err.Error())
				}
				updatedImage := &ImageWithMeta{Image: &models.Image{
					ID:     id,
					Parent: &layerParent,
					Store:  imageLayer.Store,
				},
					meta: string(meta),
					// create a dummy sum as this is required by the portlayer,
					// but since we are only updating metadata it will be ignored
					layer: FSLayer{BlobSum: "this-is-required"},
				}
				updatedImages = append(updatedImages, updatedImage)
			}

		}
	}

	log.Debugf("%d images require metadata updates", len(updatedImages))
	for id := range updatedImages {
		// utlizing index to avoid copy
		update := updatedImages[id]
		log.Debugf("updating ImageMetadata: layerID(%s)", update.ID)
		//Send the Metadata to the PL to write to disk
		err = WriteImage(ic.Host, update, nil)
		if err != nil {
			return fmt.Errorf("updateImageMetadata failed to writeMeta to the portLayer: %s", err.Error())
		}

	}
	return nil
}

// Will conditionally update array based on presence of value
func arrayUpdate(val string, list []string, action int) ([]string, bool) {

	for index, item := range list {
		if val == item {
			if action == Remove {
				// remove item requested
				list = append(list[:index], list[index+1:]...)
				return list, true
			}
			// item already in array
			return list, false
		}
	}
	if action == Add {
		list = append(list, val)
		return list, true
	}
	return list, false
}

// WriteImageBlob writes the image blob to the storage layer
func (ic *ImageC) WriteImageBlob(image *ImageWithMeta, progressOutput progress.Output, cleanup bool) error {
	defer trace.End(trace.Begin(image.Image.ID))
	if ic.Standalone {
		log.Debugf("Running in standalone, skipping WriteImageBlob")
		return nil
	}

	destination := DestinationDirectory(ic.Options)

	id := image.Image.ID
	log.Infof("Path: %s", path.Join(destination, id, id+".targ"))
	f, err := os.Open(path.Join(destination, id, id+".tar"))
	if err != nil {
		return fmt.Errorf("Failed to open file: %s", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Failed to stat file: %s", err)
	}

	in := progress.NewProgressReader(
		ioutils.NewCancelReadCloser(context.Background(), f),
		progressOutput,
		fi.Size(),
		image.String(),
		"Extracting",
	)
	defer in.Close()

	// Write the image
	err = WriteImage(ic.Host, image, in)
	if err != nil {
		return fmt.Errorf("Failed to write to image store: %s", err)
	}

	progress.Update(progressOutput, image.String(), "Pull complete")

	if cleanup {
		if err := os.RemoveAll(destination); err != nil {
			return fmt.Errorf("Failed to remove download directory: %s", err)
		}
	}
	return nil
}

// CreateImageConfig constructs the image metadata from layers that compose the image
func (ic *ImageC) CreateImageConfig(images []*ImageWithMeta) (metadata.ImageConfig, error) {

	if len(images) == 0 {
		return metadata.ImageConfig{}, nil
	}

	manifest := ic.ImageManifest
	imageLayer := images[0] // the layer that represents the actual image
	image := docker.V1Image{}
	rootFS := docker.NewRootFS()
	history := make([]docker.History, 0, len(images))
	diffIDs := make(map[string]string)
	var size int64

	// step through layers to get command history and diffID from oldest to newest
	for i := len(images) - 1; i >= 0; i-- {
		layer := images[i]
		if err := json.Unmarshal([]byte(layer.meta), &image); err != nil {
			return metadata.ImageConfig{}, fmt.Errorf("Failed to unmarshall layer history: %s", err)
		}
		h := docker.History{
			Created:   image.Created,
			Author:    image.Author,
			CreatedBy: strings.Join(image.ContainerConfig.Cmd, " "),
			Comment:   image.Comment,
		}
		history = append(history, h)
		rootFS.DiffIDs = append(rootFS.DiffIDs, dockerLayer.DiffID(layer.diffID))
		diffIDs[layer.diffID] = layer.ID
		size += layer.size
	}

	// result is constructed without unused fields
	result := docker.Image{
		V1Image: docker.V1Image{
			Comment:         image.Comment,
			Created:         image.Created,
			Container:       image.Container,
			ContainerConfig: image.ContainerConfig,
			DockerVersion:   image.DockerVersion,
			Author:          image.Author,
			Config:          image.Config,
			Architecture:    image.Architecture,
			OS:              image.OS,
		},
		RootFS:  rootFS,
		History: history,
	}

	bytes, err := result.MarshalJSON()
	if err != nil {
		return metadata.ImageConfig{}, fmt.Errorf("Failed to marshall image metadata: %s", err)
	}

	// calculate image ID
	sum := fmt.Sprintf("%x", sha256.Sum256(bytes))
	log.Infof("Image ID: sha256:%s", sum)

	// prepare metadata
	result.V1Image.Parent = image.Parent
	result.Size = size
	result.V1Image.ID = imageLayer.ID
	metaData := metadata.ImageConfig{
		V1Image: result.V1Image,
		ImageID: sum,
		// TODO: this will change when issue 1186 is
		// implemented -- only populate the digests when pulled by digest
		Digests:   []string{manifest.Digest},
		Tags:      []string{ic.Tag},
		Name:      manifest.Name,
		DiffIDs:   diffIDs,
		History:   history,
		Reference: ic.Reference,
	}

	blob, err := json.Marshal(metaData)
	if err != nil {
		return metadata.ImageConfig{}, fmt.Errorf("Failed to marshal image metadata: %s", err)
	}

	// store metadata
	imageLayer.meta = string(blob)

	return metaData, nil
}

// PullImage pulls an image from docker hub
func (ic *ImageC) PullImage() error {

	// ctx
	ctx, cancel := context.WithTimeout(ctx, ic.Options.Timeout)
	defer cancel()

	// Parse the -reference parameter
	if err := ic.ParseReference(); err != nil {
		log.Errorf(err.Error())
		return err
	}

	// Host is either the host's UUID (if run on vsphere) or the hostname of
	// the system (if run standalone)
	host, err := sys.UUID()
	if host != "" {
		log.Infof("Using UUID (%s) for imagestore name", host)
	} else if ic.Standalone {
		host, err = os.Hostname()
		log.Infof("Using host (%s) for imagestore name", host)
	}

	ic.Storename = host

	if err != nil {
		log.Errorf("Failed to return the UUID or host name: %s", err)
		return err
	}

	if !ic.Standalone {
		log.Debugf("Running with portlayer")

		// Ping the server to ensure it's at least running
		ok, err := PingPortLayer(ic.Host)
		if err != nil || !ok {
			log.Errorf("Failed to ping portlayer: %s", err)
			return err
		}
	} else {
		log.Debugf("Running standalone")
	}

	// Calculate (and overwrite) the registry URL and make sure that it responds to requests
	ic.Registry, err = LearnRegistryURL(ic.Options)
	if err != nil {
		log.Errorf("Error while pulling image: %s", err)
		return err
	}

	// Get the URL of the OAuth endpoint
	url, err := LearnAuthURL(ic.Options)
	if err != nil {
		log.Infof(err.Error())
		switch err := err.(type) {
		case urlfetcher.ImageNotFoundError:
			return fmt.Errorf("Error: image %s not found", ic.Reference)
		default:
			return fmt.Errorf("Failed to obtain OAuth endpoint: %s", err)
		}
	}

	// Get the OAuth token - if only we have a URL
	if url != nil {
		token, err := FetchToken(ctx, ic.Options, url, ic.progressOutput)
		if err != nil {
			log.Errorf("Failed to fetch OAuth token: %s", err)
			return err
		}
		ic.Token = token
	}

	progress.Message(ic.progressOutput, "", "Pulling from "+ic.Image)

	// Get the manifest
	manifest, err := FetchImageManifest(ctx, ic.Options, ic.progressOutput)
	if err != nil {
		log.Infof(err.Error())
		switch err := err.(type) {
		case urlfetcher.ImageNotFoundError:
			return fmt.Errorf("Error: image %s not found", ic.Image)
		case urlfetcher.TagNotFoundError:
			return fmt.Errorf("Tag %s not found in repository %s", ic.Tag, ic.Image)
		default:
			return fmt.Errorf("Error while pulling image manifest: %s", err)
		}
	}

	ic.ImageManifest = manifest

	err = ldm.DownloadLayers(ctx, ic)
	if err != nil {
		return err
	}

	return nil
}
