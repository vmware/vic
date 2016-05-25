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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"runtime/trace"
	"sync"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/reference"

	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/pkg/i18n"

	"github.com/pkg/profile"
)

var (
	options = ImageCOptions{}

	// https://raw.githubusercontent.com/docker/docker/master/distribution/pull_v2.go
	po = streamformatter.NewJSONStreamFormatter().NewProgressOutput(os.Stdout, false)
)

// ImageCOptions wraps the cli arguments
type ImageCOptions struct {
	reference string

	registry string
	image    string
	digest   string

	destination string

	host string

	logfile string

	username string
	password string

	token *Token

	timeout time.Duration

	stdout     bool
	debug      bool
	insecure   bool
	standalone bool
	resolv     bool

	profiling string
	tracing   bool
}

// ImageWithMeta wraps the models.Image with some additional metadata
type ImageWithMeta struct {
	*models.Image

	diffID  string
	layer   FSLayer
	history History
}

func (i *ImageWithMeta) String() string {
	return stringid.TruncateID(i.layer.BlobSum)
}

const (
	// DefaultDockerURL holds the URL of Docker registry
	DefaultDockerURL = "https://registry-1.docker.io/v2/"

	// DefaultDestination specifies the default directory to use
	DefaultDestination = "images"

	// DefaultPortLayerHost specifies the default port layer server
	DefaultPortLayerHost = "localhost:8080"

	// DefaultLogfile specifies the default log file name
	DefaultLogfile = "imagec.log"

	// DefaultHTTPTimeout specifies the default HTTP timeout
	DefaultHTTPTimeout = 3600 * time.Second

	// DefaultTokenExpirationDuration specifies the default token expiration
	DefaultTokenExpirationDuration = 60 * time.Second
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
	flag.BoolVar(&options.insecure, "insecure", false, i18n.T("Skip certificate verification checks"))
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
		return err
	}

	options.digest = reference.DefaultTag
	if !reference.IsNameOnly(ref) {
		if tagged, ok := ref.(reference.NamedTagged); ok {
			options.digest = tagged.Tag()
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

	// Use a hierachy like following so that we can support multiple schemes, registries and versions
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
		options.digest,
	)
}

// ImagesToDownload creates a slice of ImageWithMeta for the images that needs to be downloaded
func ImagesToDownload(manifest *Manifest, hostname string) ([]ImageWithMeta, error) {
	images := make([]ImageWithMeta, len(manifest.FSLayers))

	v1 := V1Compatibility{}
	// iterate from parent to children
	for i := len(manifest.History) - 1; i >= 0; i-- {
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
		images[i] = ImageWithMeta{
			Image: &models.Image{
				ID:     v1.ID,
				Parent: &parent,
				Store:  hostname,
			},
			history: history,
			layer:   layer,
		}
		log.Debugf("Manifest image: %#v", images[i])
	}

	// return early if -standalone set
	if options.standalone {
		return images, nil
	}

	// Create the image store just in case
	err := CreateImageStore(hostname)
	if err != nil {
		return nil, fmt.Errorf("Failed to create image store: %s", err)
	}

	// Get the list of known images from the storage layer
	existingImages, err := ListImages(hostname, images)
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain list of images: %s", err)
	}
	for i := range existingImages {
		log.Debugf("Existing image: %#v", existingImages[i])
	}

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

	return images, nil
}

// DownloadImageBlobs downloads the image blobs concurrently
func DownloadImageBlobs(images []ImageWithMeta) error {
	var wg sync.WaitGroup

	wg.Add(len(images))

	// iterate from parent to children
	// so that portlayer can extract each layer
	// on top of previous one
	results := make(chan error, len(images))
	for i := len(images) - 1; i >= 0; i-- {
		go func(image ImageWithMeta) {
			defer wg.Done()

			diffID, err := FetchImageBlob(options, &image)
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
func WriteImageBlobs(images []ImageWithMeta) error {
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
			ioutils.NewCancelReadCloser(
				context.Background(), f),
			po,
			fi.Size(),
			image.String(),
			"Extracting",
		)
		defer in.Close()

		// Write the image
		// FIXME: send metadata when portlayer supports it
		err = WriteImage(&image, in)
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

func main() {
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

	if err = ParseReference(); err != nil {
		log.Fatalf("Failed to parse -reference: %s", err)
	}

	// Hostname is our storename
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to return the host name: %s", err)
	}

	if !options.standalone {
		log.Debugf("Running with portlayer")

		// Ping the server to ensure it's at least running
		ok, err2 := PingPortLayer()
		if err2 != nil || !ok {
			log.Fatalf("Failed to ping portlayer: %s", err2)
		}
	} else {
		log.Debugf("Running standalone")
	}

	// Get the URL of the OAuth endpoint
	url, err := LearnAuthURL(options)
	if err != nil {
		log.Fatalf("Failed to obtain OAuth endpoint: %s", err)
	}

	// Get the OAuth token - if only we have a URL
	if url != nil {
		token, err2 := FetchToken(url)
		if err != nil {
			log.Fatalf("Failed to fetch OAuth token: %s", err2)
		}
		options.token = token
	}

	// Get the manifest
	manifest, err := FetchImageManifest(options)
	if err != nil {
		log.Fatalf("Failed to fetch image manifest: %s", err)
	}

	if !options.resolv {
		progress.Message(po, options.digest, "Pulling from "+options.image)
	}

	// Create the ImageWithMeta slice to hold Image structs
	images, err := ImagesToDownload(manifest, hostname)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if options.resolv {
		if len(images) > 0 {
			fmt.Printf("%s", images[0].history.V1Compatibility)
			os.Exit(0)
		}
		os.Exit(1)
	}

	// Fetch the blobs from registry
	if err := DownloadImageBlobs(images); err != nil {
		log.Fatalf(err.Error())
	}

	// Write blobs to the storage layer
	if err := WriteImageBlobs(images); err != nil {
		log.Fatalf(err.Error())
	}

	// FIXME: Dump the digest
	//progress.Message(po, "", "Digest: 0xDEAD:BEEF")
	if len(images) > 0 {
		progress.Message(po, "", "Status: Downloaded newer image for "+options.image+":"+options.digest)
	} else {
		progress.Message(po, "", "Status: Image is up to date for "+options.image+":"+options.digest)
	}
}
