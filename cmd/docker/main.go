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
	"crypto/x509"
	"flag"
	"fmt"
	"net"
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

	vicbackends "github.com/vmware/vic/lib/apiservers/engine/backends"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/pprof"
	viclog "github.com/vmware/vic/pkg/log"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

type CliOptions struct {
	serverPort    uint
	portLayerAddr string
	proto         string
}

const (
	productName    = "vSphere Integrated Containers"
	clientHostName = "client.localhost"
)

var vchConfig config.VirtualContainerHostConfigSpec

func init() {
	log.SetFormatter(viclog.NewTextFormatter())
	trace.Logger.Level = log.DebugLevel
	pprof.StartPprof("docker personality", pprof.DockerPort)
}

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

	if err := vicbackends.Init(cli.portLayerAddr, productName, &vchConfig, vchConfig.InsecureRegistries); err != nil {
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

	_ = flag.String("serveraddr", "127.0.0.1", "Server address to listen") // ignored
	serverPort := flag.Uint("port", 9000, "Port to listen")
	portLayerAddr := flag.String("port-layer-addr", "127.0.0.1", "Port layer server address")
	portLayerPort := flag.Uint("port-layer-port", 9001, "Port Layer server port")

	debug := flag.Bool("debug", false, "Enable debuglevel logging")

	flag.Parse()

	// load the vch config
	src, err := extraconfig.GuestInfoSource()
	if err != nil {
		log.Fatalf("Unable to load configuration from guestinfo: %s", err)
	}
	extraconfig.Decode(src, &vchConfig)

	if *debug || vchConfig.Diagnostics.DebugLevel > 0 {
		log.SetLevel(log.DebugLevel)
	}

	cli := &CliOptions{
		serverPort:    *serverPort,
		portLayerAddr: fmt.Sprintf("%s:%d", *portLayerAddr, *portLayerPort),
		proto:         "tcp",
	}

	return cli, true
}

func loadCAPool() *x509.CertPool {
	// If we should verify the server, we need to load a trusted ca
	pool := x509.NewCertPool()

	pem := vchConfig.CertificateAuthorities
	if len(pem) == 0 {
		return nil
	}

	if !pool.AppendCertsFromPEM(vchConfig.CertificateAuthorities) {
		log.Fatalf("Unable to load CAs in config")
	}

	log.Debugf("Loaded %d CAs from config", len(pool.Subjects()))
	return pool
}

func startServerWithOptions(cli *CliOptions) *apiserver.Server {
	serverConfig := &apiserver.Config{
		Logging: true,
		Version: "1.22", //dockerversion.Version,
	}

	// FIXME: assignment copies lock value to tlsConfig: crypto/tls.Config contains sync.Once contains sync.Mutex
	tlsConfig := func(c *tls.Config) *tls.Config {
		return &tls.Config{
			Certificates:             c.Certificates,
			NameToCertificate:        c.NameToCertificate,
			GetCertificate:           c.GetCertificate,
			RootCAs:                  c.RootCAs,
			NextProtos:               c.NextProtos,
			ServerName:               c.ServerName,
			ClientAuth:               c.ClientAuth,
			ClientCAs:                c.ClientCAs,
			InsecureSkipVerify:       c.InsecureSkipVerify,
			CipherSuites:             c.CipherSuites,
			PreferServerCipherSuites: c.PreferServerCipherSuites,
			SessionTicketsDisabled:   c.SessionTicketsDisabled,
			SessionTicketKey:         c.SessionTicketKey,
			ClientSessionCache:       c.ClientSessionCache,
			MinVersion:               tls.VersionTLS12,
			MaxVersion:               c.MaxVersion,
			CurvePreferences:         c.CurvePreferences,
		}
	}(&tlsconfig.ServerDefault)

	if !vchConfig.HostCertificate.IsNil() {
		log.Info("TLS enabled")

		cert, err := vchConfig.HostCertificate.Certificate()
		if err != nil {
			// This is only viable because we've verified those certificates
			log.Fatalf("Could not load certificate from config and refusing to run without TLS with a host certificate specified: %s", err)
		}

		tlsConfig.Certificates = []tls.Certificate{*cert}
		serverConfig.TLSConfig = tlsConfig

		// Set options for TLS
		if len(vchConfig.CertificateAuthorities) > 0 {
			log.Info("Client verification enabled")
			// server requires and verifies client's certificate
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs = loadCAPool()
			tlsConfig.InsecureSkipVerify = false
		}
	}

	addr := "0.0.0.0"
	// exposing this on all interfaces
	if vchConfig.Diagnostics.DebugLevel <= 2 {

		// determine the address to listen on
		ips, err := net.LookupIP(clientHostName)
		if err != nil {
			// TODO: don't want to directly enter this into vchConfig.Sessions[].Started but no
			// structure currently to report back contents otherwise
			log.Fatalf("Unable to look up %s to serve docker API: %s", clientHostName, err)
		}

		if len(ips) == 0 {
			log.Fatalf("No IP found for %s during launch of docker API server", clientHostName)
		}

		if len(ips) > 1 {
			log.Fatalf("Multiple IPs found for %s during launch of docker API server: %v", clientHostName, ips)
		}

		addr = ips[0].String()
	}

	api := apiserver.New(serverConfig)
	fullserver := fmt.Sprintf("%s:%d", addr, cli.serverPort)
	l, err := listeners.Init(cli.proto, fullserver, "", serverConfig.TLSConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listener created for HTTP on %s//%s", addr, cli.proto)
	api.Accept(fullserver, l...)

	return api
}

func setAPIRoutes(api *apiserver.Server) {
	imageHandler := &vicbackends.Image{}
	containerHandler := vicbackends.NewContainerBackend()
	volumeHandler := &vicbackends.Volume{}
	networkHandler := &vicbackends.Network{}
	systemHandler := vicbackends.NewSystemBackend()

	api.InitRouter(false,
		image.NewRouter(imageHandler),
		container.NewRouter(containerHandler),
		volume.NewRouter(volumeHandler),
		network.NewRouter(networkHandler),
		system.NewRouter(systemHandler))
}
