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
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"

	"github.com/vmware/vic/lib/apiservers/portlayer/restapi"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/dns"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/portlayer/vspc"
	"github.com/vmware/vic/lib/pprof"
	viclog "github.com/vmware/vic/pkg/log"
)

var (
	options = dns.ServerOptions{}
	parser  *flags.Parser
	server  *restapi.Server
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

func main() {
	if _, err := parser.Parse(); err != nil {
		if err := err.(*flags.Error); err != nil && err.Type == flags.ErrHelp {
			os.Exit(0)
		}

		os.Exit(1)
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

	defer server.Shutdown()

	go func() {
		vchIP, err := lookupVCHIP()
		if err != nil {
			log.Fatalf("cannot retrieve vch-endpoint ip: %v", err)
		}
		log.Infof("vSPC started...")
		doneCh := make(chan bool)
		vspc := vspc.NewVspc(vchIP.String(), constants.SerialOverLANPort, "127.0.0.1", constants.AttachServerPort, doneCh)
		for {
			_, err := vspc.Accept()
			if err != nil {
				log.Errorf("vSPC cannot accept connection: %v", err)
				doneCh <- true
				log.Errorf("vSPC exiting...")
				return
			}
		}
	}()
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

func lookupVCHIP() (net.IP, error) {
	// FIXME: THERE MUST BE ANOTHER WAY
	// following is from Create@exec.go
	ips, err := net.LookupIP(constants.ManagementHostName)
	if err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("No IP found on %s", constants.ManagementHostName)
	}

	if len(ips) > 1 {
		return nil, fmt.Errorf("Multiple IPs found on %s: %#v", constants.ManagementHostName, ips)
	}
	return ips[0], nil
}
