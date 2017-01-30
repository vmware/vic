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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"

	"github.com/vmware/vic/lib/apiservers/portlayer/restapi"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/dns"
	"github.com/vmware/vic/lib/pprof"
	viclog "github.com/vmware/vic/pkg/log"
)

var (
	options    = dns.ServerOptions{}
	parser     *flags.Parser
	server     *restapi.Server
	argsParsed bool
)

func init() {
	log.SetFormatter(viclog.NewTextFormatter())

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewPortLayerAPI(swaggerSpec)
	server = restapi.NewServer(api)

	parser = flags.NewParser(server, flags.Default)
	parser.ShortDescription = `Port Layer API`
	parser.LongDescription = `Port Layer API`

	server.ConfigureFlags()

	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

// parseArgs parses command line options using the go-flags package,
// retaining options that are unknown in a format that can be
// understood by the standard flags package
func parseArgs(unkFlags map[string]interface{}) ([]string, error) {
	argsParsed = true
	var unkArgs []string
	parser.UnknownOptionHandler = func(option string, arg flags.SplitArgument, args []string) ([]string, error) {
		if _, ok := unkFlags[option]; ok {
			unkArgs = append(unkArgs, "-"+option)
			val, exists := arg.Value()
			if exists {
				unkArgs = append(unkArgs, val)
			}

			return args, nil
		}

		return nil, fmt.Errorf("unknown option %s", option)
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

	return unkArgs, nil
}

func main() {

	defer server.Shutdown()

	if !argsParsed {
		parseArgs(nil)
	}

	pprof.StartPprof("portlayer server", pprof.PortlayerPort)

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
		restapi.StopAPIServers()
	}()

	go func() {
		dnsserver.Wait()
	}()
	// END

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
