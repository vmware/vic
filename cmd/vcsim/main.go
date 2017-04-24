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
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vmware/vic/pkg/vsphere/simulator"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
)

func main() {
	model := simulator.VPX()

	flag.IntVar(&model.Datacenter, "dc", model.Datacenter, "Number of datacenters")
	flag.IntVar(&model.Cluster, "cluster", model.Cluster, "Number of clusters")
	flag.IntVar(&model.ClusterHost, "host", model.ClusterHost, "Number of hosts per cluster")
	flag.IntVar(&model.Host, "standalone-host", model.Host, "Number of standalone hosts")
	flag.IntVar(&model.Datastore, "ds", model.Datastore, "Number of local datastores")
	flag.IntVar(&model.Machine, "vm", model.Machine, "Number of virtual machines per resource pool")
	flag.IntVar(&model.Pool, "pool", model.Pool, "Number of resource pools per compute resource")
	flag.IntVar(&model.App, "app", model.App, "Number of virtual apps per compute resource")
	flag.IntVar(&model.Pod, "pod", model.Pod, "Number of storage pods per datacenter")
	flag.IntVar(&model.Portgroup, "pg", model.Portgroup, "Number of port groups")
	flag.IntVar(&model.Folder, "folder", model.Folder, "Number of folders")

	isESX := flag.Bool("esx", false, "Simulate standalone ESX")
	isTLS := flag.Bool("tls", true, "Enable TLS")
	cert := flag.String("tlscert", "", "Path to TLS certificate file")
	key := flag.String("tlskey", "", "Path to TLS key file")
	flag.BoolVar(&simulator.Trace, "trace", simulator.Trace, "Trace SOAP to stderr")

	flag.Parse()

	f := flag.Lookup("httptest.serve")
	if f.Value.String() == "" {
		_ = f.Value.Set("127.0.0.1:8989")
	}

	if *isESX {
		opts := model
		model = simulator.ESX()
		// Preserve options that also apply to ESX
		model.Datastore = opts.Datastore
		model.Machine = opts.Machine
	}

	tag := " (govmomi simulator)"
	model.ServiceContent.About.Name += tag
	model.ServiceContent.About.OsType = runtime.GOOS + "-" + runtime.GOARCH

	esx.HostSystem.Summary.Hardware.Vendor += tag

	err := model.Create()
	if err != nil {
		log.Fatal(err)
	}

	if *isTLS {
		model.Service.TLS = new(tls.Config)
		if *cert != "" {
			c, err := tls.LoadX509KeyPair(*cert, *key)
			if err != nil {
				log.Fatal(err)
			}

			model.Service.TLS.Certificates = []tls.Certificate{c}
		}
	}

	s := model.Service.NewServer()

	fmt.Printf("GOVC_URL=%s\n", s.URL)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	model.Remove()
}
