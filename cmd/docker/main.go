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
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	apiserver "github.com/docker/docker/api/server"
	"github.com/docker/docker/api/server/router/container"
	"github.com/docker/docker/api/server/router/image"
	"github.com/docker/docker/api/server/router/network"
	"github.com/docker/docker/api/server/router/system"
	"github.com/docker/docker/api/server/router/volume"
	"github.com/docker/docker/docker/listeners"
	"github.com/docker/docker/pkg/signal"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/vmware/vic/lib/apiservers/engine/backends"
)

type CliOptions struct {
	enableTLS     bool
	verifyTLS     bool
	cafile        string
	certfile      string
	keyfile       string
	serverAddr    string
	serverPort    uint
	fullserver    string
	portLayerAddr string
	proto         string
}

const productName = "vSphere Integrated Containers"

func Usage() {
	fmt.Fprintf(os.Stderr, "\nvSphere Integrated Container Daemon Usage:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	// Get flags
	cli, ok := handleFlags()

	if !ok {
		os.Exit(1)
	}

	if err := vicbackends.Init(cli.portLayerAddr); err != nil {
		log.Fatalf("failed to initialize backend: %s", err)
	}

	// Start API server wit options from command line args
	api := startServerWithOptions(cli)

	setAPIRoutes(api)

	serveAPIWait := make(chan error)
	go api.Wait(serveAPIWait)

	signal.Trap(func() {
		api.Close()
	})

	<-serveAPIWait
}

func handleFlags() (*CliOptions, bool) {
	flag.Usage = Usage

	enableTLS := flag.Bool("TLS", false, "Use TLS; implied by --tlsverify")
	verifyTLS := flag.Bool("tlsverify", false, "Use TLS and verify the remote")
	cafile := flag.String("tls-ca-certificate", "", "Trust certs signed only by this CA")
	certfile := flag.String("tls-certificate", "", "Path to TLS certificate file")
	keyfile := flag.String("tls-key", "", "Path to TLS Key file")
	serverAddr := flag.String("serveraddr", "127.0.0.1", "Server address to listen")
	serverPort := flag.Uint("port", 9000, "Port to listen")
	portLayerAddr := flag.String("port-layer-addr", "127.0.0.1", "Port layer server address")
	portLayerPort := flag.Uint("port-layer-port", 9001, "Port Layer server port")

	flag.Parse()

	if *enableTLS && (len(*certfile) == 0 || len(*keyfile) == 0) {
		fmt.Fprintf(os.Stderr, "TLS requested, but tls-certificate and tls-key were all not specified\n")
		return nil, false
	}

	if *verifyTLS {
		*enableTLS = true

		if len(*certfile) == 0 || len(*keyfile) == 0 || len(*cafile) == 0 {
			fmt.Fprintf(os.Stderr, "tlsverfiy requested, but tls-ca-certificate, tls-certificate, tls-key were all not specified\n")
			return nil, false
		}
	}

	cli := &CliOptions{
		enableTLS:     *enableTLS,
		verifyTLS:     *verifyTLS,
		cafile:        *cafile,
		certfile:      *certfile,
		keyfile:       *keyfile,
		serverAddr:    *serverAddr,
		serverPort:    *serverPort,
		fullserver:    fmt.Sprintf("%s:%d", *serverAddr, *serverPort),
		portLayerAddr: fmt.Sprintf("%s:%d", *portLayerAddr, *portLayerPort),
		proto:         "tcp",
	}

	return cli, true
}

func startServerWithOptions(cli *CliOptions) *apiserver.Server {
	serverConfig := &apiserver.Config{
		Logging: true,
		Version: "1.22", //dockerversion.Version,
	}

	// Set options for TLS
	if cli.enableTLS {
		tlsOptions := tlsconfig.Options{
			CAFile:   cli.cafile,
			CertFile: cli.certfile,
			KeyFile:  cli.keyfile,
		}

		if cli.verifyTLS {
			// server requires and verifies client's certificate
			tlsOptions.ClientAuth = tls.RequireAndVerifyClientCert
		}

		tlsConfig, err := tlsconfig.Server(tlsOptions)

		if err != nil {
			log.Fatal(err)
		}
		serverConfig.TLSConfig = tlsConfig
	}

	api := apiserver.New(serverConfig)

	l, err := listeners.Init(cli.proto, cli.fullserver, "", serverConfig.TLSConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listener created for HTTP on TCP", cli.fullserver)
	api.Accept(cli.fullserver, l...)

	return api
}

func setAPIRoutes(api *apiserver.Server) {
	imageHandler := &vicbackends.Image{ProductName: productName}
	containerHandler := &vicbackends.Container{ProductName: productName, HackMap: make(map[string]vicbackends.V1Compatibility)}
	volumeHandler := &vicbackends.Volume{ProductName: productName}
	networkHandler := &vicbackends.Network{ProductName: productName}
	systemHandler := &vicbackends.System{ProductName: productName}

	api.InitRouter(false,
		image.NewRouter(imageHandler),
		container.NewRouter(containerHandler),
		volume.NewRouter(volumeHandler),
		network.NewRouter(networkHandler),
		system.NewRouter(systemHandler))
}
