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

package configuration

import (
	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/session"

	"golang.org/x/net/context"
)

// TODO: this object will merge with Virtual Container Host
type Configuration struct {
	DisplayName string

	ApplianceISO string
	BootstrapISO string

	CertPEM string
	KeyPEM  string

	TargetPath            string
	DatacenterName        string
	ClusterPath           string
	ImageStorePath        string
	ImageDatastoreName    string
	ResourcePoolPath      string
	ExtenalNetworkPath    string
	ExtenalNetworkName    string
	ManagementNetworkPath string
	ManagementNetworkName string
	BridgeNetworkName     string
	BridgeNetworkPath     string
	NumCPUs               int32
	MemoryMB              int64
	ImageFiles            []string

	Pool   *compute.ResourcePool
	Hosts  []*object.HostSystem
	HostIP string

	Session  *session.Session
	Context  context.Context
	insecure bool
}

type DiagnosticLog struct {
	Key     string
	Name    string
	Start   int32
	Host    *object.HostSystem
	Collect bool
}

func NewConfig() *Configuration {
	return &Configuration{
		insecure: true,
	}
}

func (conf *Configuration) ValidateConfiguration() error {
	log.Infof("Validating supplied configuration")

	var err error
	conf.Context = context.TODO()

	sessionconfig := &session.Config{
		Service:        conf.TargetPath,
		Insecure:       conf.insecure,
		DatacenterPath: conf.DatacenterName,
		ClusterPath:    conf.ClusterPath,
		DatastorePath:  conf.ImageStorePath,
		PoolPath:       conf.ResourcePoolPath,
	}

	conf.Session, err = session.NewSession(sessionconfig).Create(conf.Context)
	if err != nil {
		log.Errorf("Failed to create session: %s", err)
		return err
	}
	_, err = conf.Session.Connect(conf.Context)
	if err != nil {
		return err
	}

	if _, err = conf.Session.Populate(conf.Context); err != nil {
		log.Errorf("Failed to get resources: %s", err)
		return err
	}

	if _, err = conf.Session.Finder.NetworkOrDefault(conf.Context, conf.ExtenalNetworkPath); err != nil {
		log.Errorf("Unable to get network: %s", err)
		return err
	}

	// find the host(s) attached to given storage
	if conf.Hosts, err = conf.Session.Datastore.AttachedClusterHosts(conf.Context, conf.Session.Cluster); err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		return err
	}

	//TODO: Add more configuration validation
	return nil
}
