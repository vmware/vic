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
	"sync"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/stringid"

	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/pkg/i18n"
)

var (
	options = ImageCOptions{}

	// https://raw.githubusercontent.com/docker/docker/master/distribution/pull_v2.go
	po = streamformatter.NewJSONStreamFormatter().NewProgressOutput(os.Stdout, false)
)

// ImageCOptions wraps the cli arguments
type ImageCOptions struct {
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
}

// ImageWithMeta wraps the models.Image with some additional metadata
type ImageWithMeta struct {
	*models.Image

	layer   FSLayer
	history History
}

const (
	// DefaultDockerURL holds the URL of Docker registry
	DefaultDockerURL = "https://registry-1.docker.io/v2/"
	// DefaultDockerImage holds the default image name
	DefaultDockerImage = "library/photon"
	// DefaultDockerDigest holds the default digest name
	DefaultDockerDigest = "latest"

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

	flag.StringVar(&options.registry, "registry", DefaultDockerURL, i18n.T("Address of the registry"))
	flag.StringVar(&options.image, "image", DefaultDockerImage, i18n.T("Name of the image"))
	flag.StringVar(&options.digest, "digest", DefaultDockerDigest, i18n.T("Tag name or image digest"))

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

	flag.Parse()
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

func main() {
	// Open the log file
	f, err := os.OpenFile(options.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open the logfile %s: %s", options.logfile, err)
	}
	defer f.Close()

	// Set the log level
	if options.debug {
		log.SetLevel(log.DebugLevel)
	}

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})

	// SetOutput to log file and/or stdout
	log.SetOutput(f)
	if options.stdout {
		log.SetOutput(io.MultiWriter(os.Stdout, f))
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

	progress.Message(po, options.digest, "Pulling from "+options.image)

	// List of ImageWithMeta to hold Image structs
	images := make([]ImageWithMeta, len(manifest.FSLayers))

	v1 := V1Compatibility{}
	// iterate from parent to children
	for i := len(manifest.History) - 1; i >= 0; i-- {
		history := manifest.History[i]
		layer := manifest.FSLayers[i]

		// unmarshall V1Compatibility to get the image ID
		if err := json.Unmarshal([]byte(history.V1Compatibility), &v1); err != nil {
			log.Fatalf("Failed to unmarshall image history: %s", err)
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

	var existingImages map[string]*models.Image

	if !options.standalone {
		// Create the image store
		err = CreateImageStore(hostname)
		if err != nil {
			log.Fatalf("Failed to create image store: %s", err)
		}

		// Get the list of existing images
		existingImages, err = ListImages(hostname, images)
		if err != nil {
			log.Fatalf("Failed to obtain list of images: %s", err)
		}
		for i := range existingImages {
			log.Debugf("Existing image: %#v", existingImages[i])
		}
	}

	// iterate from parent to children
	// so that we can delete from the slice
	// while iterating over it
	for i := len(images) - 1; i >= 0; i-- {
		ID := images[i].ID
		if _, ok := existingImages[ID]; ok {
			log.Debugf("%s already exists", ID)
			// delete existing image from images
			images = append(images[:i], images[i+1:]...)

			progress.Update(po, stringid.TruncateID(ID), "Already exists")
		}
	}

	var wg sync.WaitGroup

	wg.Add(len(images))

	// iterate from parent to children
	// so that portlayer can extract each layer
	// on top of previous one
	results := make(chan error, len(images))
	for i := len(images) - 1; i >= 0; i-- {
		go func(image ImageWithMeta) {
			defer wg.Done()

			err := FetchImageBlob(options, &image)
			if err != nil {
				results <- fmt.Errorf("%s/%s returned %s", options.image, image.layer.BlobSum, err)
			} else {
				results <- nil
			}
		}(images[i])
	}
	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			log.Fatalf("Failed to fetch image blob: %s", err)
		}
	}

	if !options.standalone {

		// iterate from parent to children
		// so that portlayer can extract each layer
		// on top of previous one
		destination := DestinationDirectory()
		for i := len(images) - 1; i >= 0; i-- {
			image := images[i]

			id := image.Image.ID
			f, err := os.Open(path.Join(destination, id, id+".tar"))
			if err != nil {
				log.Fatalf("Failed to open file: %s", err)
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				log.Fatalf("Failed to stat file: %s", err)
			}

			in := progress.NewProgressReader(
				ioutils.NewCancelReadCloser(
					context.Background(), f),
				po,
				fi.Size(),
				stringid.TruncateID(id),
				"Extracting",
			)
			defer in.Close()

			// Write the image
			// FIXME: send metadata when portlayer supports it
			err = WriteImage(&image, in)
			if err != nil {
				log.Fatalf("Failed to write to image store: %s", err)
			}
			progress.Update(po, stringid.TruncateID(id), "Pull complete")
		}
	}
	// FIXME: Dump the digest
	//progress.Message(po, "", "Digest: 0xDEAD:BEEF")

	if len(images) > 0 {
		progress.Message(po, "", "Status: Downloaded newer image for "+options.image+":"+options.digest)
	} else {
		progress.Message(po, "", "Status: Image is up to date for "+options.image+":"+options.digest)

	}
}
