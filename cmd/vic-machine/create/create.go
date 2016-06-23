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
	"github.com/vmware/vic/cmd/vic-machine/data"
	"github.com/vmware/vic/cmd/vic-machine/validate"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"

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
	*data.Data

	applianceISO string
	bootstrapISO string

	cert string
	key  string

	tlsGenerate bool

	osType  string
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
	create := &Create{}
	create.Data = data.NewData()

	// FIXME: make this a flag
	create.logfile = "create.log"

	return create
}

// Flags return all cli flags for create
func (c *Create) Flags() []cli.Flag {
	flags := []cli.Flag{
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
			Usage:       "Compute resource path, e.g. myCluster/Resources/myRP",
			Destination: &c.ComputeResourcePath,
		},
		cli.StringFlag{
			Name:        "image-datastore",
			Value:       "",
			Usage:       "Image datastore name",
			Destination: &c.ImageDatastoreName,
		},
		cli.StringFlag{
			Name:        "container-datastore",
			Value:       "",
			Usage:       "Container datastore name - defaults to image datastore",
			Destination: &c.ContainerDatastoreName,
		},
		cli.StringFlag{
			Name:        "name",
			Value:       "docker-appliance",
			Usage:       "The name of the Virtual Container Host",
			Destination: &c.DisplayName,
		},
		cli.StringFlag{
			Name:        "external-network",
			Value:       "",
			Usage:       "The external network (can see hub.docker.com)",
			Destination: &c.ExternalNetworkName,
		},
		cli.StringFlag{
			Name:        "management-network",
			Value:       "",
			Usage:       "The management network (provides route to target hosting vSphere)",
			Destination: &c.ManagementNetworkName,
		},
		cli.StringFlag{
			Name:        "client-network",
			Value:       "",
			Usage:       "The client network (restricts DOCKER_API access to this network)",
			Destination: &c.ClientNetworkName,
		},
		cli.StringFlag{
			Name:        "bridge-network",
			Value:       "",
			Usage:       "The bridge network (private port group for containers)",
			Destination: &c.BridgeNetworkName,
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
			Destination: &c.Force,
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
			Destination: &c.Timeout,
		},
		cli.IntFlag{
			Name:        "appliance-memory",
			Value:       2048,
			Usage:       "Memory for the appliance VM, in MB",
			Destination: &c.MemoryMB,
		},
		cli.IntFlag{
			Name:        "appliance-cpu",
			Value:       1,
			Usage:       "vCPUs for the appliance VM",
			Destination: &c.NumCPUs,
		},
	}
	flags = append(flags, c.TargetFlags()...)
	flags = append(flags, c.DebugFlags()...)
	return flags
}

func (c *Create) processParams() error {
	if err := c.HasCredentials(); err != nil {
		return err
	}

	if c.ImageDatastoreName == "" {
		return cli.NewExitError("--image-datastore Image datastore name must be specified", 1)
	}

	if c.cert != "" && c.key == "" {
		return cli.NewExitError("key cert should be specified at the same time", 1)
	}
	if c.cert == "" && c.key != "" {
		return cli.NewExitError("key cert should be specified at the same time", 1)
	}

	if c.ExternalNetworkName == "" {
		c.ExternalNetworkName = "VM Network"
	}

	if c.BridgeNetworkName == "" {
		c.BridgeNetworkName = c.DisplayName
	}

	if len(c.DisplayName) > MaxDisplayNameLen {
		return cli.NewExitError(fmt.Sprintf("Display name %s exceeds the permitted 31 characters limit. Please use a shorter -name parameter", c.DisplayName), 1)
	}

	// FIXME: add parameters for these configurations
	c.osType = "linux"

	c.Insecure = true
	return nil
}

func (c *Create) loadCertificate() (*Keypair, error) {
	var keypair *Keypair
	if c.cert != "" && c.key != "" {
		log.Infof("Loading certificate/key pair - private key in %s", c.key)
		keypair = NewKeyPair(false, c.key, c.cert)
	} else if c.tlsGenerate && c.DisplayName != "" {
		c.key = fmt.Sprintf("./%s-key.pem", c.DisplayName)
		c.cert = fmt.Sprintf("./%s-cert.pem", c.DisplayName)
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

	if c.Debug.Debug {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if err = c.processParams(); err != nil {
		return err
	}

	var images []string
	if images, err = c.checkImagesFiles(); err != nil {
		return err
	}

	log.Infof("### Installing VCH ####")

	var keypair *Keypair
	if keypair, err = c.loadCertificate(); err != nil {
		log.Error("Creation cannot continue: unable to load certificate")
		return err
	}

	if keypair != nil {
		c.KeyPEM = keypair.KeyPEM
		c.CertPEM = keypair.CertPEM
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	validator, err := validate.NewValidator(ctx, c.Data)
	if err != nil {
		log.Error("Creation cannot continue: failed to create validator")
		return err
	}

	vchConfig, err := validator.Validate(ctx, c.Data)
	if err != nil {
		log.Error("Creation cannot continue: configuration validaton failed")
		return err
	}

	vConfig := validator.AddDeprecatedFields(ctx, vchConfig, c.Data)
	vConfig.ImageFiles = images

	executor := management.NewDispatcher(ctx, validator.Session, vchConfig, c.Force)
	if err = executor.Dispatch(vchConfig, vConfig); err != nil {

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
	tls := ""
	if c.key != "" {
		// if we're generating then there's no CA currently
		if len(vchConfig.CertificateAuthorities) > 0 {
			tls = fmt.Sprintf(" --tls --tlscert='%s' --tlskey='%s'", c.cert, c.key)
		} else {
			tls = " --tls"
		}
	}
	log.Infof("DOCKER_HOST=%s:%s", executor.HostIP, executor.DockerPort)
	log.Infof("DOCKER_OPTS=\"-H %s:%s%s\"", executor.HostIP, executor.DockerPort, tls)
	log.Infof("")
	log.Infof("Connect to docker:")
	log.Infof("docker -H %s:%s%s info", executor.HostIP, executor.DockerPort, tls)

	log.Infof("Installer completed successfully")
	return nil
}
