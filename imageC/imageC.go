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
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
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

	stdout bool
	debug  bool
}

const (
	DefaultDockerURL    = "https://registry-1.docker.io/v2/"
	DefaultDockerImage  = "library/photon"
	DefaultDockerDigest = "latest"

	DefaultDestination = "."

	DefaultHost = "localhost:80"

	DefaultLogfile = "imageC.log"

	DefaultHTTPTimeout             = 60 * time.Second
	DefaultTokenExpirationDuration = 60 * time.Second
)

func init() {
	flag.StringVar(&options.registry, "registry", DefaultDockerURL, "Address of the registry")
	flag.StringVar(&options.image, "image", DefaultDockerImage, "Name of the image")
	flag.StringVar(&options.digest, "digest", DefaultDockerDigest, "Tag name or image digest")

	flag.StringVar(&options.destination, "destination", DefaultDestination, "Destination directory")

	flag.StringVar(&options.host, "host", DefaultHost, "Host that runs portlayer API (FQDN:port format)")

	flag.StringVar(&options.logfile, "logfile", DefaultLogfile, "Path of the installer log file")

	flag.StringVar(&options.username, "username", "", "Username")
	flag.StringVar(&options.password, "password", "", "Password")

	flag.DurationVar(&options.timeout, "timeout", DefaultHTTPTimeout, "HTTP timeout")

	flag.BoolVar(&options.stdout, "stdout", false, "Enable writing to stdout")
	flag.BoolVar(&options.debug, "debug", false, "Enable debugging")

	flag.Parse()
}

func main() {
	// Open log file
	f, err := os.OpenFile(options.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening logfile %s: %v", options.logfile, err)
	}
	defer f.Close()

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

	// FIXME: Ping the portlayer to make sure that it's up and running

	url, err := LearnAuthURL(options)
	if err != nil {
		log.Fatalf("%s", err)
	}

	token, err := FetchToken(url)
	if err != nil {
		log.Fatalf("%s", err)
	}
	options.token = token

	manifest, err := FetchImageManifest(options)
	if err != nil {
		log.Fatalf("%s", err)
	}

	var wg sync.WaitGroup

	layers := manifest.FSLayers
	histories := manifest.History

	progress.Message(po, options.digest, "Pulling from "+options.image)

	wg.Add(len(layers))
	results := make(chan error, len(layers))
	for i := 0; i < len(layers); i++ {
		go func(layer string, history string) {
			defer wg.Done()

			err := FetchImageBlob(options, layer, history)
			if err != nil {
				results <- fmt.Errorf("%s/%s returned %s", options.image, layer, err)
			} else {
				results <- nil
			}
		}(layers[i].BlobSum, histories[i].V1Compatibility)
	}

	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	// FIXME: Dump the digest
	//progress.Message(po, "", "Digest: 0xDEAD:BEEF")

	progress.Message(po, "", "Status: Downloaded newer image for "+options.image+":"+options.digest)
}
