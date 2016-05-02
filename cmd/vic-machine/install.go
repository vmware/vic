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
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/install/configuration"
	"github.com/vmware/vic/install/management"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	// Max permitted length of Virtual Machine name
	MaxVirtualMachineNameLen = 80
	// Max permitted length of Virtual Switch name
	MaxDisplayNameLen  = 31
	ApplianceImageKey  = "core"
	LinuxImageKey      = "linux"
	ApplianceImageName = "appliance.iso"
	LinuxImageName     = "bootstrap.iso"
)

// Data is pulled from govmomi tests where we needed a global config as we could pass it around
type Data struct {
	target              string
	user                string
	passwd              string
	computeResourcePath string
	imageDatastoreName  string
	displayName         string

	containerDatastoreName string
	externalNetworkName    string
	managementNetworkName  string
	bridgeNetworkName      string

	numCPUs  int32
	memoryMB int64

	applianceISO string
	bootstrapISO string

	cert string
	key  string

	force       bool
	tlsGenerate bool

	osType  string
	timeout time.Duration
	logfile string

	conf     *configuration.Configuration
	executor *management.Dispatcher
}

var (
	BuildID string
	data    = &Data{}
	images  = map[string][]string{
		ApplianceImageKey: []string{ApplianceImageName},
		LinuxImageKey:     []string{LinuxImageName},
	}
)

func init() {
	flag.StringVar(&data.target, "target", "", "ESXi or vCenter FQDN or IPv4 address")
	flag.StringVar(&data.user, "user", "", "ESX or vCenter user")
	flag.StringVar(&data.passwd, "passwd", "", "ESX or vCenter password")
	flag.StringVar(&data.cert, "cert", "", "Virtual Container Host x509 certificate file")
	flag.StringVar(&data.key, "key", "", "Virtual Container Host private key file")
	flag.StringVar(&data.computeResourcePath, "compute-resource", "", "Compute resource path, e.g. /ha-datacenter/host/myCluster/Resources/myRP")
	flag.StringVar(&data.imageDatastoreName, "image-store", "", " Image datastore name")
	flag.StringVar(&data.containerDatastoreName, "container-store", "", " Container datastore name - defaults to image datastore")
	flag.StringVar(&data.displayName, "name", "docker-appliance", "The name of the Virtual Container Host")
	flag.StringVar(&data.externalNetworkName, "external-network", "", "The external network (can see hub.docker.com)")
	flag.StringVar(&data.managementNetworkName, "management-network", "", "The management network (can see target)")
	flag.StringVar(&data.bridgeNetworkName, "bridge-network", "", "The bridge network")
	flag.StringVar(&data.applianceISO, "appliance-iso", "", "The appliance iso")
	flag.StringVar(&data.bootstrapISO, "bootstrap-iso", "", "The bootstrap iso")
	flag.BoolVar(&data.force, "force", false, "Force the install, removing existing if present")
	flag.BoolVar(&data.tlsGenerate, "generate-cert", false, "Generate certificate for Virtual Container Host")

	flag.Parse()
}

func processParams() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s BUILD ID: %s\n", os.Args[0], BuildID)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if data.target == "" {
		log.Errorf("-target argument must be specified")
		flag.Usage()
	}

	if data.user == "" {
		log.Errorf("-user User to login target must be specified")
		flag.Usage()
	}

	if data.computeResourcePath == "" {
		log.Errorf("-compute-resource Compute resource path must be specified")
		flag.Usage()
	}

	if data.imageDatastoreName == "" {
		log.Errorf("-image-store Image datastore name must be specified")
		flag.Usage()
	}

	if data.cert != "" && data.key == "" {
		log.Errorf("key cert should be specified at the same time")
	}
	if data.cert == "" && data.key != "" {
		log.Errorf("key cert should be specified at the same time")
	}

	if data.externalNetworkName == "" {
		data.externalNetworkName = "VM Network"
	}

	if data.bridgeNetworkName == "" {
		data.bridgeNetworkName = data.displayName
	}

	if len(data.displayName) > MaxDisplayNameLen {
		log.Fatalf("Display name %s exceeds the permitted 31 characters limit. Please use a shorter -name parameter", data.displayName)
	}

	// FIXME: add parameters for these configurations
	data.osType = "linux"
	data.timeout = 3 * time.Minute
	data.logfile = "install.log"

	data.conf = configuration.NewConfig()
	// FIXME: add parameters for these configurations
	data.conf.NumCPUs = 1
	data.conf.MemoryMB = 2048
	data.conf.DisplayName = data.displayName
	resources := strings.Split(data.computeResourcePath, "/")
	if len(resources) < 2 || resources[1] == "" {
		log.Fatalf("Could not determine datacenter from specified -compute path, %s", data.computeResourcePath)
	}
	data.conf.DatacenterName = resources[1]
	data.conf.ClusterPath = strings.Split(data.computeResourcePath, "/Resources")[0]

	if data.conf.ClusterPath == "" {
		log.Fatalf("Could not determine cluster from specified -compute path, %s", data.computeResourcePath)
	}

	data.conf.ResourcePoolPath = data.computeResourcePath
	data.conf.ImageDatastoreName = data.imageDatastoreName
	data.conf.ImageStorePath = fmt.Sprintf("/%s/datastore/%s", data.conf.DatacenterName, data.imageDatastoreName)
	data.conf.ExtenalNetworkPath = fmt.Sprintf("/%s/network/%s", data.conf.DatacenterName, data.externalNetworkName)
	data.conf.ExtenalNetworkName = data.externalNetworkName
	data.conf.BridgeNetworkPath = fmt.Sprintf("/%s/network/%s", data.conf.DatacenterName, data.bridgeNetworkName)
	data.conf.BridgeNetworkName = data.bridgeNetworkName
	if data.managementNetworkName != "" {
		data.conf.ManagementNetworkPath = fmt.Sprintf("/%s/network/%s", data.conf.DatacenterName, data.managementNetworkName)
		data.conf.ManagementNetworkName = data.managementNetworkName
	}
	//prompt for passwd if not specified
	if data.passwd == "" {
		log.Print("Please enter ESX or vCenter password: ")
		b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("Failed to read password from stdin: %s", err)
		}
		data.passwd = string(b)
	}
	data.conf.TargetPath = fmt.Sprintf("%s:%s@%s/sdk", data.user, data.passwd, data.target)
}

func validateConfiguration() error {
	if err := loadCertificate(); err != nil {
		return err
	}

	return data.conf.ValidateConfiguration()
}

func loadCertificate() error {
	var keypair *Keypair
	if data.cert != "" && data.key != "" {
		log.Infof("Loading certificate/key pair - private key in %s", data.key)
		keypair = NewKeyPair(false, data.key, data.cert)
	} else if data.tlsGenerate {
		data.key = fmt.Sprintf("./%s-key.pem", data.displayName)
		data.cert = fmt.Sprintf("./%s-cert.pem", data.displayName)
		log.Infof("Generating certificate/key pair - private key in %s", data.key)
		keypair = NewKeyPair(true, data.key, data.cert)
	}
	if keypair == nil {
		log.Warnf("Configuring without TLS - to enable use -gen or -key/-cert parameters")
		return nil
	}
	if err := keypair.GetCertificate(); err != nil {
		log.Errorf("Failed to read/generate certificate: %s", err)
		return err
	}
	data.conf.KeyPEM = keypair.KeyPEM
	data.conf.CertPEM = keypair.CertPEM
	return nil
}

func checkImagesFiles() error {
	// detect images files
	osImgs, ok := images[data.osType]
	if !ok {
		return fmt.Errorf("Specified OS \"%s\" is not known to this installer", data.osType)
	}

	var imgs []string
	if data.applianceISO != "" {
		imgs = append(imgs, data.applianceISO)
	} else {
		imgs = append(imgs, images[ApplianceImageKey]...)
	}
	if data.bootstrapISO != "" {
		imgs = append(imgs, data.bootstrapISO)
	} else {
		imgs = append(imgs, osImgs...)
	}

	for _, img := range imgs {
		_, err := os.Stat(img)
		if os.IsNotExist(err) {
			var dir string
			dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
			_, err = os.Stat(filepath.Join(dir, img))
			if err == nil {
				img = filepath.Join(dir, img)
			}
		}

		if os.IsNotExist(err) {
			log.Warnf("\t\tUnable to locate %s in the current or installer directory.", img)
			return err
		}
		data.conf.ImageFiles = append(data.conf.ImageFiles, img)
	}
	return nil
}

func main() {
	var err error
	processParams()

	if err = checkImagesFiles(); err != nil {
		log.Fatalf("%s", err)
	}

	// Open log file
	f, err := os.OpenFile(data.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening logfile %s: %v", data.logfile, err)
	}
	defer f.Close()

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	// SetOutput to io.MultiWriter so that we can log to stdout and a file
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	log.Infof("### Installing VCH ####")

	if err = validateConfiguration(); err != nil {
		log.Fatalf("Validating supplied configuration failed. Exiting...")
	}

	executor := management.NewDispatcher(data.conf, data.force)
	if err = executor.Dispatch(data.conf, data.timeout); err != nil {
		executor.CollectDiagnosticLogs()
		log.Fatal(err)
	}

	log.Infof("Initialization of appliance successful")

	log.Infof("")
	log.Infof("SSH to appliance (default=root:password)")
	log.Infof("ssh root@%s", executor.HostIP)
	log.Infof("")
	log.Infof("Log server:")
	log.Infof("https://%s:2378", executor.HostIP)
	log.Infof("")
	if data.key != "" {
		log.Infof("Connect to docker:")
		log.Infof("docker -H %s:%s --tls --tlscert='%s' --tlskey='%s' info", executor.HostIP, executor.DockerPort, data.cert, data.key)
	} else {
		log.Infof("DOCKER_HOST=%s:%s", executor.HostIP, executor.DockerPort)
	}

	log.Infof("Installer completed successfully...")
}
