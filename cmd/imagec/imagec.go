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
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"runtime/trace"
	"strings"
	"sync"
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
	"github.com/vmware/vic/pkg/i18n"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/sys"

	"github.com/pkg/profile"
)

var (
	options = ImageCOptions{}

	// https://raw.githubusercontent.com/docker/docker/master/distribution/pull_v2.go
	sf = streamformatter.NewJSONStreamFormatter()
	po = sf.NewProgressOutput(os.Stdout, false)
)

// ImageCOptions wraps the cli arguments
type ImageCOptions struct {
	reference string

	registry string
	image    string
	tag      string

	destination string

	host string

	logfile string

	username string
	password string

	token *Token

	timeout time.Duration

	stdout bool
	debug  bool

	insecureSkipVerify bool
	insecureAllowHTTP  bool

	standalone bool
	resolv     bool

	profiling string
	tracing   bool
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

	// DefaultTokenExpirationDuration specifies the default token expiration
	DefaultTokenExpirationDuration = 60 * time.Second

	// attribute update actions
	Add = iota + 1
	Remove
)

func init() {
	// TODO: get language from host OS
	lang := i18n.DefaultLang
	languageFile := fmt.Sprintf("messages/%s", lang.String())
	data, err := Asset(languageFile)
	if err != nil {
		panic(fmt.Sprintf("Invalid language asset: %v", languageFile))
	}
	i18n.LoadLanguageBytes(lang, data)

	flag.StringVar(&options.reference, "reference", "", i18n.T("Name of the reference"))

	flag.StringVar(&options.destination, "destination", DefaultDestination, i18n.T("Destination directory"))

	flag.StringVar(&options.host, "host", DefaultPortLayerHost, i18n.T("Host that runs portlayer API (FQDN:port format)"))

	flag.StringVar(&options.logfile, "logfile", DefaultLogfile, i18n.T("Path of the imagec log file"))

	flag.StringVar(&options.username, "username", "", i18n.T("Username"))
	flag.StringVar(&options.password, "password", "", i18n.T("Password"))

	flag.DurationVar(&options.timeout, "timeout", DefaultHTTPTimeout, i18n.T("HTTP timeout"))

	flag.BoolVar(&options.stdout, "stdout", false, i18n.T("Enable writing to stdout"))
	flag.BoolVar(&options.debug, "debug", false, i18n.T("Show debug logging"))
	flag.BoolVar(&options.insecureSkipVerify, "insecure-skip-verify", false, i18n.T("Don't verify certificates when fetching images"))
	flag.BoolVar(&options.insecureAllowHTTP, "insecure-allow-http", false, i18n.T("Uses unencrypted connections when fetching images"))
	flag.BoolVar(&options.standalone, "standalone", false, i18n.T("Disable port-layer integration"))

	flag.BoolVar(&options.resolv, "resolv", false, i18n.T("Return the name of the vmdk from given reference"))

	flag.StringVar(&options.profiling, "profile.mode", "", i18n.T("Enable profiling mode, one of [cpu, mem, block]"))
	flag.BoolVar(&options.tracing, "tracing", false, i18n.T("Enable runtime tracing"))

	flag.Parse()
}

// ParseReference parses the -reference parameter and populate options struct
func ParseReference() error {
	// Validate and parse reference name
	ref, err := reference.ParseNamed(options.reference)
	if err != nil {
		log.Warn("Error while parsing reference %s: %#v", options.reference, err)
		return err
	}

	options.tag = reference.DefaultTag
	if !reference.IsNameOnly(ref) {
		if tagged, ok := ref.(reference.NamedTagged); ok {
			options.tag = tagged.Tag()
		}
	}

	options.registry = DefaultDockerURL
	if ref.Hostname() != reference.DefaultHostname {
		options.registry = ref.Hostname()
	}

	options.image = ref.RemoteName()

	return nil
}

// DestinationDirectory returns the path of the output directory
func DestinationDirectory() string {
	u, _ := url.Parse(options.registry)

	// Use a hierarchy like following so that we can support multiple schemes, registries and versions
	/*
		https/
		├── 192.168.218.5:5000
		│   └── v2
		│       └── busybox
		│           └── latest
		...
		│               ├── fef924a0204a00b3ec67318e2ed337b189c99ea19e2bf10ed30a13b87c5e17ab
		│               │   ├── fef924a0204a00b3ec67318e2ed337b189c99ea19e2bf10ed30a13b87c5e17ab.json
		│               │   └── fef924a0204a00b3ec67318e2ed337b189c99ea19e2bf10ed30a13b87c5e17ab.tar
		│               └── manifest.json
		└── registry-1.docker.io
		    └── v2
		        └── library
		            └── golang
		                └── latest
		                    ...
		                    ├── f61ebe2817bb4e6a7f0a4cf249a5316223f7ecc886feac24b9887a490feaed57
		                    │   ├── f61ebe2817bb4e6a7f0a4cf249a5316223f7ecc886feac24b9887a490feaed57.json
		                    │   └── f61ebe2817bb4e6a7f0a4cf249a5316223f7ecc886feac24b9887a490feaed57.tar
		                    └── manifest.json

	*/
	return path.Join(
		options.destination,
		u.Scheme,
		u.Host,
		u.Path,
		options.image,
		options.tag,
	)
}

// ImagesToDownload creates a slice of ImageWithMeta for the images that needs to be downloaded
func ImagesToDownload(manifest *Manifest, storeName string) ([]*ImageWithMeta, *ImageWithMeta, error) {
	images := make([]*ImageWithMeta, len(manifest.FSLayers))

	v1 := docker.V1Image{}
	// iterate from parent to children
	for i := len(manifest.History) - 1; i >= 0; i-- {
		history := manifest.History[i]
		layer := manifest.FSLayers[i]

		// unmarshall V1Compatibility to get the image ID
		if err := json.Unmarshal([]byte(history.V1Compatibility), &v1); err != nil {
			return nil, nil, fmt.Errorf("Failed to unmarshall image history: %s", err)
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
				Store:  storeName,
			},
			meta:   history.V1Compatibility,
			layer:  layer,
			diffID: "",
		}
		log.Debugf("ImagesToDownload Manifest image: %#v", images[i])
	}

	// return early if -standalone set
	if options.standalone {
		return images, nil, nil
	}

	// Get the list of known images from the storage layer
	existingImages, err := ListImages(storeName, images)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to obtain list of images: %s", err)
	}

	for i := range existingImages {
		log.Debugf("Existing image: %#v", existingImages[i])
	}

	// grab the imageLayer for use in later evaluation of metadata update
	imageLayer := images[0]

	// iterate from parent to children
	// so that we can delete from the slice
	// while iterating over it
	for i := len(images) - 1; i >= 0; i-- {
		ID := images[i].ID
		// Check whether storage layer knows this image ID
		if _, ok := existingImages[ID]; ok {
			log.Debugf("%s already exists", ID)
			// update the progress before deleting it from the slice
			progress.Update(po, images[i].String(), "Already exists")

			// delete existing image from images
			images = append(images[:i], images[i+1:]...)
		}
	}

	return images, imageLayer, nil
}

// Will updated metadata based on the current image being requested
//
// Currently on tags are updated, but once imagec supports pulling by digest
// then digest will need to be updated as well
func updateImageMetadata(imageLayer *ImageWithMeta, manifest *Manifest) error {
	// if standalone no need to continue
	if options.standalone {
		return nil
	}
	// updated Image slice
	var updatedImages []*ImageWithMeta

	// assuming calls to the portlayer are cheap, lets get all existing images
	existingImages, err := ListImages(imageLayer.Store, nil)
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
		if imageMeta.Name == options.image {

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
		err = WriteImage(update, nil)
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

// DownloadImageBlobs downloads the image blobs concurrently
func DownloadImageBlobs(images []*ImageWithMeta) error {
	var wg sync.WaitGroup

	wg.Add(len(images))

	// iterate from parent to children
	// so that portlayer can extract each layer
	// on top of previous one
	results := make(chan error, len(images))
	for i := len(images) - 1; i >= 0; i-- {
		go func(image *ImageWithMeta) {
			defer wg.Done()

			diffID, err := FetchImageBlob(options, image)
			if err != nil {
				results <- fmt.Errorf("%s/%s returned %s", options.image, image.layer.BlobSum, err)
			} else {
				image.diffID = diffID
				results <- nil
			}
		}(images[i])
	}
	wg.Wait()
	close(results)

	// iterate over results chan to see whether we have a failed download
	for err := range results {
		if err != nil {
			return fmt.Errorf("Failed to fetch image blob: %s", err)
		}
	}

	return nil
}

// WriteImageBlobs writes the image blob to the storage layer
func WriteImageBlobs(images []*ImageWithMeta) error {
	if options.standalone {
		return nil
	}

	// iterate from parent to children
	// so that portlayer can extract each layer
	// on top of previous one
	destination := DestinationDirectory()
	for i := len(images) - 1; i >= 0; i-- {
		image := images[i]

		id := image.Image.ID
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
			po,
			fi.Size(),
			image.String(),
			"Extracting",
		)
		defer in.Close()

		// Write the image
		err = WriteImage(image, in)
		if err != nil {
			return fmt.Errorf("Failed to write to image store: %s", err)
		}
		progress.Update(po, image.String(), "Pull complete")
	}
	if err := os.RemoveAll(destination); err != nil {
		return fmt.Errorf("Failed to remove download directory: %s", err)
	}
	return nil
}

// CreateImageConfig constructs the image metadata from layers that compose the image
func CreateImageConfig(images []*ImageWithMeta, manifest *Manifest, reference string) error {

	if len(images) == 0 {
		return nil
	}

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
			return fmt.Errorf("Failed to unmarshall layer history: %s", err)
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
		return fmt.Errorf("Failed to marshall image metadata: %s", err)
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
		Tags:      []string{options.tag},
		Name:      manifest.Name,
		DiffIDs:   diffIDs,
		History:   history,
		Reference: reference,
	}

	blob, err := json.Marshal(metaData)
	if err != nil {
		return fmt.Errorf("Failed to marshal image metadata: %s", err)
	}

	// store metadata
	imageLayer.meta = string(blob)

	return nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, string(sf.FormatError(fmt.Errorf("%s : %s", r, debug.Stack()))))
		}
	}()

	if version.Show() {
		fmt.Fprintf(os.Stdout, "%s\n", version.String())
		return
	}

	// Enable profiling if mode is set
	switch options.profiling {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.Quiet).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.Quiet).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile, profile.ProfilePath("."), profile.Quiet).Stop()
	default:
		// do nothing
	}

	// Register our custom Error hook
	log.AddHook(NewErrorHook(os.Stderr))

	// Enable runtime tracing if tracing is true
	if options.tracing {
		tracing, err := os.Create(time.Now().Format("2006-01-02T150405.pprof"))
		if err != nil {
			log.Fatalf("Failed to create tracing logfile: %s", err)
		}
		defer tracing.Close()

		if err := trace.Start(tracing); err != nil {
			log.Fatalf("Failed to start tracing: %s", err)
		}
		defer trace.Stop()
	}

	// Open the log file
	f, err := os.OpenFile(options.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open the logfile %s: %s", options.logfile, err)
	}
	defer f.Close()

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})

	// Set the log level
	if options.debug {
		log.SetLevel(log.DebugLevel)
	}

	// SetOutput to log file and/or stdout
	log.SetOutput(f)
	if options.stdout {
		log.SetOutput(io.MultiWriter(os.Stdout, f))
	}

	// Parse the -reference parameter
	if err = ParseReference(); err != nil {
		log.Fatalf(err.Error())
	}

	// Host is either the host's UUID (if run on vsphere) or the hostname of
	// the system (if run standalone)
	host, err := sys.UUID()
	if host != "" {
		log.Infof("Using UUID (%s) for imagestore name", host)
	} else if options.standalone {
		host, err = os.Hostname()
		log.Infof("Using host (%s) for imagestore name", host)
	}

	if err != nil {
		log.Fatalf("Failed to return the UUID or host name: %s", err)
	}

	if !options.standalone {
		log.Debugf("Running with portlayer")

		// Ping the server to ensure it's at least running
		ok, err := PingPortLayer()
		if err != nil || !ok {
			log.Fatalf("Failed to ping portlayer: %s", err)
		}
	} else {
		log.Debugf("Running standalone")
	}

	// Calculate (and overwrite) the registry URL and make sure that it responds to requests
	options.registry, err = LearnRegistryURL(options)
	if err != nil {
		log.Fatalf("Error while pulling image: %s", err)
	}

	// Get the URL of the OAuth endpoint
	url, err := LearnAuthURL(options)
	if err != nil {
		log.Fatalf("Failed to obtain OAuth endpoint: %s", err)
	}

	// Get the OAuth token - if only we have a URL
	if url != nil {
		token, err := FetchToken(url)
		if err != nil {
			log.Fatalf("Failed to fetch OAuth token: %s", err)
		}
		options.token = token
	}

	// HACK: Required to learn the name of the vmdk from given reference
	// Used by docker personality until metadata support lands
	if !options.resolv {
		progress.Message(po, "", "Pulling from "+options.image)
	}

	// Get the manifest
	manifest, err := FetchImageManifest(options)
	if err != nil {
		switch err := err.(type) {
		case ImageNotFoundError:
			log.Fatalf("Error: image %s not found", options.image)
		case TagNotFoundError:
			log.Fatalf("Tag %s not found in repository %s", options.tag, options.image)
		default:
			log.Fatalf("Error while pulling image manifest: %s", err)
		}
	}

	// Create the ImageWithMeta slice to hold Image structs
	images, imageLayer, err := ImagesToDownload(manifest, host)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// HACK: Required to learn the name of the vmdk from given reference
	// Used by docker personality until metadata support lands
	if options.resolv {
		if len(images) > 0 {
			fmt.Printf("%s", images[0].meta)
			os.Exit(0)
		}
		os.Exit(1)
	}

	// Fetch the blobs from registry
	if err := DownloadImageBlobs(images); err != nil {
		log.Fatalf(err.Error())
	}

	if err := CreateImageConfig(images, manifest, options.reference); err != nil {
		log.Fatalf(err.Error())
	}

	// Write blobs to the storage layer
	if err := WriteImageBlobs(images); err != nil {
		log.Fatalf(err.Error())
	}

	if err := updateImageMetadata(imageLayer, manifest); err != nil {
		log.Fatalf(err.Error())
	}

	progress.Message(po, "", "Digest: "+manifest.Digest)
	if len(images) > 0 {
		progress.Message(po, "", "Status: Downloaded newer image for "+options.image+":"+options.tag)
	} else {
		progress.Message(po, "", "Status: Image is up to date for "+options.image+":"+options.tag)
	}
}
