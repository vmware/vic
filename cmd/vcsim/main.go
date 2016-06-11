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
	"log"
	"runtime"

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

	isESX := flag.Bool("esx", false, "Simulate standalone ESX")

	flag.Parse()

	f := flag.Lookup("httptest.serve")
	if f.Value.String() == "" {
		f.Value.Set("localhost:8989")
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

	defer model.Remove()

	err := model.Create()
	if err != nil {
		log.Fatal(err)
	}

	model.Service.NewServer()
}
