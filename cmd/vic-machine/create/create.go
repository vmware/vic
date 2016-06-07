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

package create

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/flags"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
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

// Create has all input parameters for vic-machine create command
type Create struct {
	target              string
	user                string
	password            *string
	computeResourcePath string
	imageDatastoreName  string
	displayName         string

	containerDatastoreName string
	externalNetworkName    string
	managementNetworkName  string
	bridgeNetworkName      string

	numCPUs  int
	memoryMB int
	insecure bool

	applianceISO string
	bootstrapISO string

	cert string
	key  string

	force       bool
	tlsGenerate bool

	osType  string
	timeout time.Duration
	logfile string

	executor *management.Dispatcher
}

var (
	images = map[string][]string{
		ApplianceImageKey: []string{ApplianceImageName},
		LinuxImageKey:     []string{LinuxImageName},
	}
)

func NewCreate() *Create {
	return &Create{}
}

// Flags return all cli flags for create
func (c *Create) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "target",
			Value:       "",
			Usage:       "ESXi or vCenter FQDN or IPv4 address",
			Destination: &c.target,
		},
		cli.StringFlag{
			Name:        "user",
			Value:       "",
			Usage:       "ESX or vCenter user",
			Destination: &c.user,
		},
		cli.GenericFlag{
			Name:  "password, p",
			Value: flags.NewOptionalString(&c.password),
			Usage: "ESX or vCenter password",
		},
		cli.StringFlag{
			Name:        "cert",
			Value:       "",
			Usage:       "Virtual Container Host x509 certificate file",
			Destination: &c.cert,
		},
		cli.StringFlag{
			Name:        "key",
			Value:       "",
			Usage:       "Virtual Container Host private key file",
			Destination: &c.key,
		},
		cli.StringFlag{
			Name:        "compute-resource",
			Value:       "",
			Usage:       "Compute resource path, e.g. /ha-datacenter/host/myCluster/Resources/myRP",
			Destination: &c.computeResourcePath,
		},
		cli.StringFlag{
			Name:        "image-datastore",
			Value:       "",
			Usage:       "Image datastore name",
			Destination: &c.imageDatastoreName,
		},
		cli.StringFlag{
			Name:        "container-datastore",
			Value:       "",
			Usage:       "Container datastore name - defaults to image datastore",
			Destination: &c.containerDatastoreName,
		},
		cli.StringFlag{
			Name:        "name",
			Value:       "docker-appliance",
			Usage:       "The name of the Virtual Container Host",
			Destination: &c.displayName,
		},
		cli.StringFlag{
			Name:        "external-network",
			Value:       "",
			Usage:       "The external network (can see hub.docker.com)",
			Destination: &c.externalNetworkName,
		},
		cli.StringFlag{
			Name:        "management-network",
			Value:       "",
			Usage:       "The management network (can see target)",
			Destination: &c.managementNetworkName,
		},
		cli.StringFlag{
			Name:        "bridge-network",
			Value:       "",
			Usage:       "The bridge network",
			Destination: &c.bridgeNetworkName,
		},
		cli.StringFlag{
			Name:        "appliance-iso",
			Value:       "",
			Usage:       "The appliance iso",
			Destination: &c.applianceISO,
		},
		cli.StringFlag{
			Name:        "bootstrap-iso",
			Value:       "",
			Usage:       "The bootstrap iso",
			Destination: &c.bootstrapISO,
		},
		cli.BoolFlag{
			Name:        "force",
			Usage:       "Force the install, removing existing if present",
			Destination: &c.force,
		},
		cli.BoolTFlag{
			Name:        "generate-cert",
			Usage:       "Generate certificate for Virtual Container Host",
			Destination: &c.tlsGenerate,
		},
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for appliance initialization",
			Destination: &c.timeout,
		},
		cli.IntFlag{
			Name:        "appliance-memory",
			Value:       2048,
			Usage:       "Memory for the appliance VM, in MB",
			Destination: &c.memoryMB,
		},
		cli.IntFlag{
			Name:        "appliance-cpu",
			Value:       1,
			Usage:       "vCPUs for the appliance VM",
			Destination: &c.numCPUs,
		},
	}
}

func (c *Create) processParams() error {
	if c.target == "" {
		return cli.NewExitError("--target argument must be specified", 1)
	}

	if c.user == "" {
		return cli.NewExitError("--user User to login target must be specified", 1)
	}

	if c.computeResourcePath == "" {
		return cli.NewExitError("--compute-resource Compute resource path must be specified", 1)
	}

	if c.imageDatastoreName == "" {
		return cli.NewExitError("--image-datastore Image datastore name must be specified", 1)
	}

	if c.cert != "" && c.key == "" {
		return cli.NewExitError("key and cert should be specified at the same time", 1)
	}
	if c.cert == "" && c.key != "" {
		return cli.NewExitError("key and cert should be specified at the same time", 1)
	}

	if c.externalNetworkName == "" {
		c.externalNetworkName = "VM Network"
	}

	if c.bridgeNetworkName == "" {
		c.bridgeNetworkName = c.displayName
	}

	if len(c.displayName) > MaxDisplayNameLen {
		return cli.NewExitError(fmt.Sprintf("Display name %s exceeds the permitted 31 characters limit. Please use a shorter -name parameter", c.displayName), 1)
	}

	//prompt for passwd if not specified
	if c.password == nil {
		log.Print("Please enter ESX or vCenter password: ")
		b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			message := fmt.Sprintf("Failed to read password from stdin: %s", err)
			cli.NewExitError(message, 1)
		}
		sb := string(b)
		c.password = &sb
	}

	// FIXME: add parameters for these configurations
	c.osType = "linux"
	c.logfile = "install.log"

	c.insecure = true
	return nil
}

func (c *Create) loadCertificate() (*Keypair, error) {
	var keypair *Keypair
	if c.cert != "" && c.key != "" {
		log.Infof("Loading certificate/key pair - private key in %s", c.key)
		keypair = NewKeyPair(false, c.key, c.cert)
	} else if c.tlsGenerate {
		c.key = fmt.Sprintf("./%s-key.pem", c.displayName)
		c.cert = fmt.Sprintf("./%s-cert.pem", c.displayName)
		log.Infof("Generating certificate/key pair - private key in %s", c.key)
		keypair = NewKeyPair(true, c.key, c.cert)
	}
	if keypair == nil {
		log.Warnf("Configuring without TLS - to enable use -generate-cert or -key/-cert parameters")
		return nil, nil
	}
	if err := keypair.GetCertificate(); err != nil {
		log.Errorf("Failed to read/generate certificate: %s", err)
		return nil, err
	}
	return keypair, nil
}

func (c *Create) checkImagesFiles() ([]string, error) {
	// detect images files
	osImgs, ok := images[c.osType]
	if !ok {
		return nil, fmt.Errorf("Specified OS \"%s\" is not known to this installer", c.osType)
	}

	var imgs []string
	var result []string
	if c.applianceISO != "" {
		imgs = append(imgs, c.applianceISO)
	} else {
		imgs = append(imgs, images[ApplianceImageKey]...)
	}
	if c.bootstrapISO != "" {
		imgs = append(imgs, c.bootstrapISO)
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
			return nil, err
		}
		result = append(result, img)
	}
	return result, nil
}

func (c *Create) Run(cli *cli.Context) error {
	var err error
	if err = c.processParams(); err != nil {
		return err
	}

	var images []string
	if images, err = c.checkImagesFiles(); err != nil {
		return err
	}

	// Open log file
	f, err := os.OpenFile(c.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		err = errors.Errorf("Error opening logfile %s: %v", c.logfile, err)
		return err
	}
	defer f.Close()

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	// SetOutput to io.MultiWriter so that we can log to stdout and a file
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	log.Infof("### Installing VCH ####")

	var keypair *Keypair
	if keypair, err = c.loadCertificate(); err != nil {
		err = errors.Errorf("Loading certificate failed with %s. Exiting...", err)
		return err
	}

	validator := NewValidator()
	vchConfig, err := validator.Validate(c)
	if err != nil {
		err = errors.Errorf("%s. Exiting...", err)
		return err
	}

	if keypair != nil {
		vchConfig.KeyPEM = keypair.KeyPEM
		vchConfig.CertPEM = keypair.CertPEM
	}
	vchConfig.ImageFiles = images

	var cancel context.CancelFunc
	validator.Context, cancel = context.WithTimeout(validator.Context, c.timeout)
	defer cancel()
	executor := management.NewDispatcher(validator.Context, validator.Session, vchConfig, c.force)
	if err = executor.Dispatch(vchConfig); err != nil {
		executor.CollectDiagnosticLogs()
		return err
	}

	log.Infof("Initialization of appliance successful")

	log.Infof("")
	log.Infof("SSH to appliance (default=root:password)")
	log.Infof("ssh root@%s", executor.HostIP)
	log.Infof("")
	log.Infof("Log server:")
	log.Infof("%s://%s:2378", executor.VICAdminProto, executor.HostIP)
	log.Infof("")
	if c.key != "" {
		log.Infof("Connect to docker:")
		log.Infof("docker -H %s:%s --tls --tlscert='%s' --tlskey='%s' info", executor.HostIP, executor.DockerPort, c.cert, c.key)
	} else {
		log.Infof("DOCKER_HOST=%s:%s", executor.HostIP, executor.DockerPort)
	}

	log.Infof("Installer completed successfully...")
	return nil
}
