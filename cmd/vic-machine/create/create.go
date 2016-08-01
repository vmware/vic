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
	"encoding"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/certificate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/flags"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"

	"bytes"

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

var EntireOptionHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

// Create has all input parameters for vic-machine create command
type Create struct {
	*data.Data

	cert string
	key  string

	noTLS           bool
	advancedOptions bool

	osType  string
	logfile string

	containerNetworks         cli.StringSlice
	containerNetworksGateway  cli.StringSlice
	containerNetworksIPRanges cli.StringSlice
	containerNetworksDNS      cli.StringSlice
	volumeStores              cli.StringSlice

	memoryReservLimits string
	cpuReservLimits    string

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
			Name:        "image-datastore, i",
			Value:       "",
			Usage:       "REQUIRED. Image datastore name",
			Destination: &c.ImageDatastoreName,
		},
		cli.StringFlag{
			Name:        "container-datastore, cs",
			Value:       "",
			Usage:       "Container datastore name - not supported yet, default to image datastore",
			Destination: &c.ContainerDatastoreName,
		},
		cli.StringSliceFlag{
			Name:  "volume-store, vs",
			Value: &c.volumeStores,
			Usage: "Specify location and label for volume store; path optional: \"datastore/path:label\" or \"datastore:label\"",
		},
		cli.StringFlag{
			Name:        "bridge-network, b",
			Value:       "",
			Usage:       "The bridge network (private port group for containers)",
			Destination: &c.BridgeNetworkName,
		},
		cli.StringFlag{
			Name:        "external-network, en",
			Value:       "",
			Usage:       "The external network (can see hub.docker.com)",
			Destination: &c.ExternalNetworkName,
		},
		cli.StringFlag{
			Name:        "management-network, mn",
			Value:       "",
			Usage:       "The management network (provides route to target hosting vSphere)",
			Destination: &c.ManagementNetworkName,
		},
		cli.StringFlag{
			Name:        "client-network, cln",
			Value:       "",
			Usage:       "The client network (restricts DOCKER_API access to this network)",
			Destination: &c.ClientNetworkName,
		},
		cli.StringSliceFlag{
			Name:  "container-network, cn",
			Value: &c.containerNetworks,
			Usage: "Networks that containers can use",
		},
		cli.StringSliceFlag{
			Name:  "container-network-gateway, cng",
			Value: &c.containerNetworksGateway,
			Usage: "Gateway for the container network's subnet in CONTAINER-NETWORK:SUBNET format, e.g. a_network:172.16.0.0/16",
		},
		cli.StringSliceFlag{
			Name:  "container-network-ip-range, cnr",
			Value: &c.containerNetworksIPRanges,
			Usage: "IP range for the container network in CONTAINER-NETWORK:IP-RANGE format, e.g. a_network:172.16.0.0/24, a_network:172.16.0.10-20",
		},
		cli.StringSliceFlag{
			Name:  "container-network-dns, cnd",
			Value: &c.containerNetworksDNS,
			Usage: "DNS servers for the container network in CONTAINER-NETWORK:DNS format, e.g. a_network:8.8.8.8",
		},
		cli.IntFlag{
			Name:        "pool-memory-reservation, pmr",
			Value:       0,
			Usage:       "VCH Memory reservation in MB",
			Destination: &c.VCHMemoryReservationsMB,
		},
		cli.IntFlag{
			Name:        "pool-memory-limit, pml",
			Value:       0,
			Usage:       "VCH Memory limit in MB",
			Destination: &c.VCHMemoryLimitsMB,
		},
		cli.GenericFlag{
			Name:  "pool-memory-shares, pms",
			Value: flags.NewSharesFlag(&c.VCHMemoryShares),
			Usage: "VCH Memory shares in level or share number, e.g. high, normal, low, or 163840",
		},
		cli.IntFlag{
			Name:        "pool-cpu-reservation, pcr",
			Value:       0,
			Usage:       "VCH vCPUs reservation in MHz",
			Destination: &c.VCHCPUReservationsMHz,
		},
		cli.IntFlag{
			Name:        "pool-cpu-limit, pcl",
			Value:       0,
			Usage:       "VCH vCPUs limit in MHz",
			Destination: &c.VCHCPULimitsMHz,
		},
		cli.GenericFlag{
			Name:  "pool-cpu-shares, pcs",
			Value: flags.NewSharesFlag(&c.VCHCPUShares),
			Usage: "VCH vCPUs shares, in level or share number, e.g. high, normal, low, or 4000",
		},
		cli.StringFlag{
			Name:        "bridge-network-range, bnr",
			Value:       "172.16.0.0/12",
			Usage:       "The ip range from which bridge networks are allocated",
			Destination: &c.BridgeIPRange,
		},
		cli.StringFlag{
			Name:        "appliance-iso, ai",
			Value:       "",
			Usage:       "The appliance iso",
			Destination: &c.ApplianceISO,
		},
		cli.StringFlag{
			Name:        "bootstrap-iso, bi",
			Value:       "",
			Usage:       "The bootstrap iso",
			Destination: &c.BootstrapISO,
		},
		cli.StringFlag{
			Name:        "key",
			Value:       "",
			Usage:       "Virtual Container Host private key file",
			Destination: &c.key,
		},
		cli.StringFlag{
			Name:        "cert",
			Value:       "",
			Usage:       "Virtual Container Host x509 certificate file",
			Destination: &c.cert,
		},
		cli.BoolFlag{
			Name:        "no-tls, k",
			Usage:       "Disable TLS support",
			Destination: &c.noTLS,
		},
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Force the install, removing existing if present",
			Destination: &c.Force,
		},
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for appliance initialization",
			Destination: &c.Timeout,
		},
		cli.BoolFlag{
			Name:        "advanced-options, x",
			Usage:       "Show all options",
			Destination: &c.advancedOptions,
		},
		cli.IntFlag{
			Name:        "appliance-memory",
			Value:       2048,
			Usage:       "Memory for the appliance VM, in MB",
			Hidden:      true,
			Destination: &c.MemoryMB,
		},
		cli.IntFlag{
			Name:        "appliance-cpu",
			Value:       1,
			Usage:       "vCPUs for the appliance VM",
			Hidden:      true,
			Destination: &c.NumCPUs,
		},
		cli.BoolFlag{
			Name:        "use-rp",
			Usage:       "Use resource pool for vch parent in VC",
			Destination: &c.UseRP,
			Hidden:      true,
		},
	}
	preFlags := append(c.TargetFlags(), c.ComputeFlags()...)
	flags = append(preFlags, flags...)
	flags = append(flags, c.DebugFlags()...)
	return flags
}

func (c *Create) processVolumeStores() error {
	defer trace.End(trace.Begin(""))
	c.VolumeLocations = make(map[string]string)
	for _, arg := range c.volumeStores {
		splitMeta := strings.SplitN(arg, ":", 2)
		if len(splitMeta) != 2 {
			return errors.New("Volume store input must be in format datastore-path:label")
		}
		c.VolumeLocations[splitMeta[1]] = splitMeta[0]
	}

	return nil

}

func (c *Create) processParams() error {
	defer trace.End(trace.Begin(""))

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

	if err := c.processContainerNetworks(); err != nil {
		return err
	}

	if err := c.processVolumeStores(); err != nil {
		return errors.Errorf("Error occurred while processing volume stores: %s", err)
	}

	//	if err := c.processReservations(); err != nil {
	//		return err
	//	}
	// FIXME: add parameters for these configurations
	c.osType = "linux"

	c.Insecure = true
	return nil
}

func (c *Create) processContainerNetworks() error {
	gws, err := parseContainerNetworkGateways([]string(c.containerNetworksGateway))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	pools, err := parseContainerNetworkIPRanges([]string(c.containerNetworksIPRanges))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	dns, err := parseContainerNetworkDNS([]string(c.containerNetworksDNS))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// parse container networks
	for _, cn := range c.containerNetworks {
		vnet, v, err := splitVnetParam(cn)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		vicnet := vnet
		if v != "" {
			vicnet = v
		}

		c.MappedNetworks[vicnet] = vnet
		c.MappedNetworksGateways[vicnet] = gws[vnet]
		c.MappedNetworksIPRanges[vicnet] = pools[vnet]
		c.MappedNetworksDNS[vicnet] = dns[vnet]

	}

	return nil
}

func (c *Create) loadCertificate() (*certificate.Keypair, error) {
	defer trace.End(trace.Begin(""))

	var keypair *certificate.Keypair
	if c.cert != "" && c.key != "" {
		log.Infof("Loading certificate/key pair - private key in %s", c.key)
		keypair = certificate.NewKeyPair(false, c.key, c.cert)
	} else if !c.noTLS && c.DisplayName != "" {
		c.key = fmt.Sprintf("./%s-key.pem", c.DisplayName)
		c.cert = fmt.Sprintf("./%s-cert.pem", c.DisplayName)
		log.Infof("Generating certificate/key pair - private key in %s", c.key)
		keypair = certificate.NewKeyPair(true, c.key, c.cert)
	}
	if keypair == nil {
		log.Warnf("Configuring without TLS - to enable drop --no-tls or use --key/--cert parameters")
		return nil, nil
	}
	if err := keypair.GetCertificate(); err != nil {
		log.Errorf("Failed to read/generate certificate: %s", err)
		return nil, err
	}
	return keypair, nil
}

func (c *Create) checkImagesFiles() ([]string, error) {
	defer trace.End(trace.Begin(""))

	// detect images files
	osImgs, ok := images[c.osType]
	if !ok {
		return nil, fmt.Errorf("Specified OS \"%s\" is not known to this installer", c.osType)
	}

	var imgs []string
	var result []string
	if c.ApplianceISO == "" {
		c.ApplianceISO = images[ApplianceImageKey][0]
	}
	imgs = append(imgs, c.ApplianceISO)

	if c.BootstrapISO == "" {
		c.BootstrapISO = osImgs[0]
	}
	imgs = append(imgs, c.BootstrapISO)

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

func (c *Create) Run(cliContext *cli.Context) error {
	var err error

	if c.advancedOptions {
		cli.HelpPrinter(cliContext.App.Writer, EntireOptionHelpTemplate, cliContext.Command)
		return nil
	}

	if c.Debug.Debug > 0 {
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

	if len(cliContext.Args()) > 0 {
		log.Error("Create cannot continue: invalid CLI arguments")
		log.Errorf("Unknown argument: %s", cliContext.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	log.Infof("### Installing VCH ####")

	var keypair *certificate.Keypair
	if keypair, err = c.loadCertificate(); err != nil {
		log.Error("Create cannot continue: unable to load certificate")
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
		log.Error("Create cannot continue: failed to create validator")
		return err
	}

	vchConfig, err := validator.Validate(ctx, c.Data)
	if err != nil {
		log.Error("Create cannot continue: configuration validation failed")
		return err
	}

	if keypair != nil {
		vchConfig.UserKeyPEM = string(keypair.KeyPEM)
		vchConfig.UserCertPEM = string(keypair.CertPEM)
	}

	vConfig := validator.AddDeprecatedFields(ctx, vchConfig, c.Data)
	vConfig.ImageFiles = images

	{ // create certificates for VCH extension
		var certbuffer, keybuffer bytes.Buffer
		if certbuffer, keybuffer, err = certificate.CreateRawKeyPair(); err != nil {
			return errors.Errorf("Failed to create certificate for VIC vSphere extension: %s", err)
		}
		vchConfig.ExtensionCert = certbuffer.String()
		vchConfig.ExtensionKey = keybuffer.String()
	}

	executor := management.NewDispatcher(ctx, validator.Session, vchConfig, c.Force)
	if err = executor.CreateVCH(vchConfig, vConfig); err != nil {

		executor.CollectDiagnosticLogs()
		return err
	}

	log.Infof("Initialization of appliance successful")

	executor.ShowVCH(vchConfig, c.key, c.cert)
	log.Infof("Installer completed successfully")
	return nil
}

type ipNetUnmarshaler struct {
	ipnet *net.IPNet
	ip    net.IP
}

func (m *ipNetUnmarshaler) UnmarshalText(text []byte) error {
	s := string(text)
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}

	m.ipnet = ipnet
	m.ip = ip
	return nil
}
func parseContainerNetworkGateways(cgs []string) (map[string]net.IPNet, error) {
	gws := make(map[string]net.IPNet)
	for _, cg := range cgs {
		m := &ipNetUnmarshaler{}
		vnet, err := parseVnetParam(cg, m)
		if err != nil {
			return nil, err
		}

		if _, ok := gws[vnet]; ok {
			return nil, fmt.Errorf("Duplicate gateway specified for container network %s", vnet)
		}

		gws[vnet] = net.IPNet{IP: m.ip, Mask: m.ipnet.Mask}
	}

	return gws, nil
}

func parseContainerNetworkIPRanges(cps []string) (map[string][]ip.Range, error) {
	pools := make(map[string][]ip.Range)
	for _, cp := range cps {
		ipr := &ip.Range{}
		vnet, err := parseVnetParam(cp, ipr)
		if err != nil {
			return nil, err
		}

		pools[vnet] = append(pools[vnet], *ipr)
	}

	return pools, nil
}

func parseContainerNetworkDNS(cds []string) (map[string][]net.IP, error) {
	dns := make(map[string][]net.IP)
	for _, cd := range cds {
		var ip net.IP
		vnet, err := parseVnetParam(cd, &ip)
		if err != nil {
			return nil, err
		}

		if ip == nil {
			return nil, fmt.Errorf("DNS IP not specified for container network %s", vnet)
		}

		dns[vnet] = append(dns[vnet], ip)
	}

	return dns, nil
}

func splitVnetParam(p string) (vnet string, value string, err error) {
	mapped := strings.Split(p, ":")
	if len(mapped) == 0 || len(mapped) > 2 {
		err = fmt.Errorf("Invalid value for parameter %s", p)
		return
	}

	vnet = mapped[0]
	if vnet == "" {
		err = fmt.Errorf("Container network not specified in parameter %s", p)
		return
	}

	if len(mapped) > 1 {
		value = mapped[1]
	}

	return
}

func parseVnetParam(p string, m encoding.TextUnmarshaler) (vnet string, err error) {
	vnet, v, err := splitVnetParam(p)
	if err != nil {
		return "", fmt.Errorf("Error parsing container network parameter %s: %s", p, err)
	}

	if err = m.UnmarshalText([]byte(v)); err != nil {
		return "", fmt.Errorf("Error parsing container network parameter %s: %s", p, err)
	}

	return vnet, nil
}
