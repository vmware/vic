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
	"net/url"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"

	"github.com/vmware/vic/lib/apiservers/portlayer/restapi"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	ploptions "github.com/vmware/vic/lib/apiservers/portlayer/restapi/options"
	"github.com/vmware/vic/lib/dns"
	"github.com/vmware/vic/lib/pprof"
	"github.com/vmware/vic/lib/vspc"
	viclog "github.com/vmware/vic/pkg/log"
	"github.com/vmware/vic/pkg/log/syslog"
	"github.com/vmware/vic/pkg/trace"
)

var (
	options = dns.ServerOptions{}
	parser  *flags.Parser
	server  *restapi.Server
)

func init() {
	trace.Logger = log.StandardLogger()

	pprof.StartPprof("portlayer server", pprof.PortlayerPort)

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

	logcfg := viclog.NewLoggingConfig()
	if ploptions.PortLayerOptions.Debug {
		logcfg.Level = log.DebugLevel
		syslog.Logger.Level = log.DebugLevel
	}

	if ploptions.PortLayerOptions.SyslogAddr != nil {
		u, err := url.Parse(*ploptions.PortLayerOptions.SyslogAddr)
		if err != nil {
			log.Fatalln(err)
		}

		logcfg.Syslog = &viclog.SyslogConfig{
			Network:  u.Scheme,
			RAddr:    u.Host,
			Priority: syslog.LOG_INFO | syslog.LOG_DAEMON,
		}
	}

	log.Infof("%+v", *logcfg)
	viclog.Init(logcfg)

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

	vspc := vspc.NewVspc()

	if vspc == nil {
		log.Fatalln("cannot initialize virtual serial port concentrator")
	}
	vspc.Start()

	go func() {
		<-sig

		vspc.Stop()
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
