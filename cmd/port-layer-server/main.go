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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"

	"github.com/vmware/vic/lib/apiservers/portlayer/restapi"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/dns"
	"github.com/vmware/vic/lib/pprof"
)

var (
	options = dns.ServerOptions{}
)

func init() {
	pprof.StartPprof("portlayer server", pprof.PortlayerPort)
}

func main() {

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewPortLayerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = `Port Layer API`
	parser.LongDescription = `Port Layer API`

	server.ConfigureFlags()

	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	// BEGIN
	// Set the Interface name to instruct listeners to bind on this interface
	options.Interface = "bridge"

	// Start the DNS Server
	dnsserver := dns.NewServer(options)
	if dnsserver != nil {
		dnsserver.Start()
	}

	// handle the signals and gracefully shutdown the server
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig

		dnsserver.Stop()

		server.Shutdown()
	}()

	go func() {
		dnsserver.Wait()
	}()
	// END

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
