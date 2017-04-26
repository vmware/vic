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

package create

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/certificate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/flags"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
)

const (
	// Max permitted length of Virtual Machine name
	MaxVirtualMachineNameLen = 80
	// Max permitted length of Virtual Switch name
	MaxDisplayNameLen = 31

	clientCert = "cert.pem"
	clientKey  = "key.pem"
	serverCert = "server-cert.pem"
	serverKey  = "server-key.pem"
	caCert     = "ca.pem"
	caKey      = "ca-key.pem"
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

	certPath   string
	scert      string
	skey       string
	ccert      string
	ckey       string
	cacert     string
	cakey      string
	clientCert *tls.Certificate

	envFile string

	cname   string
	org     cli.StringSlice
	keySize int

	noTLS           bool
	noTLSverify     bool
	advancedOptions bool

	clientCAs   cli.StringSlice
	registryCAs cli.StringSlice

	containerNetworks         cli.StringSlice
	containerNetworksGateway  cli.StringSlice
	containerNetworksIPRanges cli.StringSlice
	containerNetworksDNS      cli.StringSlice
	volumeStores              cli.StringSlice
	insecureRegistries        cli.StringSlice
	dns                       cli.StringSlice
	clientNetworkName         string
	clientNetworkGateway      string
	clientNetworkIP           string
	publicNetworkName         string
	publicNetworkGateway      string
	publicNetworkIP           string
	managementNetworkName     string
	managementNetworkGateway  string
	managementNetworkIP       string

	memoryReservLimits string
	cpuReservLimits    string

	BridgeIPRange string

	httpsProxy string
	httpProxy  string

	executor *management.Dispatcher
}

func NewCreate() *Create {
	create := &Create{}
	create.Data = data.NewData()

	return create
}

// Flags return all cli flags for create
func (c *Create) Flags() []cli.Flag {
	create := []cli.Flag{
		// credentials
		cli.StringFlag{
			Name:        "ops-user",
			Value:       "",
			Usage:       "The user with which the VCH operates after creation. Defaults to the credential supplied with target",
			Destination: &c.OpsUser,
			Hidden:      true,
		},
		cli.GenericFlag{
			Name:   "ops-password",
			Value:  flags.NewOptionalString(&c.OpsPassword),
			Usage:  "Password or token for the operations user. Defaults to the credential supplied with target",
			Hidden: true,
		},

		// images
		cli.StringFlag{
			Name:        "image-store, i",
			Value:       "",
			Usage:       "Image datastore path in format \"datastore/path\"",
			Destination: &c.ImageDatastorePath,
		},
		cli.StringFlag{
			Name:        "base-image-size",
			Value:       "8GB",
			Usage:       "Specify the size of the base image from which all other images are created e.g. 8GB/8000MB",
			Destination: &c.ScratchSize,
			Hidden:      true,
		},

		// container disk
		cli.StringFlag{
			Name:        "container-store, cs",
			Value:       "",
			Usage:       "Container datastore path - not supported yet, defaults to image datastore",
			Destination: &c.ContainerDatastoreName,
			Hidden:      true,
		},

		// volume
		cli.StringSliceFlag{
			Name:  "volume-store, vs",
			Value: &c.volumeStores,
			Usage: "Specify a list of location and label for volume store, e.g. \"datastore/path:label\" or \"datastore:label\".",
		},

		// bridge
		cli.StringFlag{
			Name:        "bridge-network, b",
			Value:       "",
			Usage:       "The bridge network port group name (private port group for containers). Defaults to the Virtual Container Host name",
			Destination: &c.BridgeNetworkName,
		},
		cli.StringFlag{
			Name:        "bridge-network-range, bnr",
			Value:       "172.16.0.0/12",
			Usage:       "The IP range from which bridge networks are allocated",
			Destination: &c.BridgeIPRange,
			Hidden:      true,
		},

		// client
		cli.StringFlag{
			Name:        "client-network, cln",
			Value:       "",
			Usage:       "The client network port group name (restricts DOCKER_API access to this network). Defaults to DCHP - see advanced help (-x)",
			Destination: &c.clientNetworkName,
		},
		cli.StringFlag{
			Name:        "client-network-gateway",
			Value:       "",
			Usage:       "Gateway for the VCH on the client network, including one or more routing destinations in a comma separated list, e.g. 10.1.0.0/16,10.2.0.0/16:10.0.0.1",
			Destination: &c.clientNetworkGateway,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "client-network-ip",
			Value:       "",
			Usage:       "IP address with a network mask for the VCH on the client network, e.g. 10.0.0.2/24",
			Destination: &c.clientNetworkIP,
			Hidden:      true,
		},

		// public
		cli.StringFlag{
			Name:        "public-network, pn",
			Value:       "VM Network",
			Usage:       "The public network port group name (port forwarding and default route). Defaults to 'VM Network' and DHCP -- see advanced help (-x)",
			Destination: &c.publicNetworkName,
		},
		cli.StringFlag{
			Name:        "public-network-gateway",
			Value:       "",
			Usage:       "Gateway for the VCH on the public network, e.g. 10.0.0.1",
			Destination: &c.publicNetworkGateway,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "public-network-ip",
			Value:       "",
			Usage:       "IP address with a network mask for the VCH on the public network, e.g. 10.0.1.2/24",
			Destination: &c.publicNetworkIP,
			Hidden:      true,
		},

		// management
		cli.StringFlag{
			Name:        "management-network, mn",
			Value:       "",
			Usage:       "The management network port group name (provides route to target hosting vSphere). Defaults to DCHP - see advanced help (-x)",
			Destination: &c.managementNetworkName,
		},
		cli.StringFlag{
			Name:        "management-network-gateway",
			Value:       "",
			Usage:       "Gateway for the VCH on the management network, including one or more routing destinations in a comma separated list, e.g. 10.1.0.0/16,10.2.0.0/16:10.0.0.1",
			Destination: &c.managementNetworkGateway,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "management-network-ip",
			Value:       "",
			Usage:       "IP address with a network mask for the VCH on the management network, e.g. 10.0.2.2/24",
			Destination: &c.managementNetworkIP,
			Hidden:      true,
		},

		// general DNS
		cli.StringSliceFlag{
			Name:   "dns-server",
			Value:  &c.dns,
			Usage:  "DNS server for the client, public, and management networks. Defaults to 8.8.8.8 and 8.8.4.4 when VCH uses static IP",
			Hidden: true,
		},

		// container networks - mapped from vSphere
		cli.StringSliceFlag{
			Name:  "container-network, cn",
			Value: &c.containerNetworks,
			Usage: "vSphere network list that containers can use directly with labels, e.g. vsphere-net:backend. Defaults to DCHP - see advanced help (-x).",
		},
		cli.StringSliceFlag{
			Name:   "container-network-gateway, cng",
			Value:  &c.containerNetworksGateway,
			Usage:  "Gateway for the container network's subnet in CONTAINER-NETWORK:SUBNET format, e.g. vsphere-net:172.16.0.1/16.",
			Hidden: true,
		},
		cli.StringSliceFlag{
			Name:   "container-network-ip-range, cnr",
			Value:  &c.containerNetworksIPRanges,
			Usage:  "IP range for the container network in CONTAINER-NETWORK:IP-RANGE format, e.g. vsphere-net:172.16.0.0/24, vsphere-net:172.16.0.10-20.",
			Hidden: true,
		},
		cli.StringSliceFlag{
			Name:   "container-network-dns, cnd",
			Value:  &c.containerNetworksDNS,
			Usage:  "DNS servers for the container network in CONTAINER-NETWORK:DNS format, e.g. vsphere-net:8.8.8.8. Ignored if no static IP assigned.",
			Hidden: true,
		},

		// memory
		cli.IntFlag{
			Name:        "memory, mem",
			Value:       0,
			Usage:       "VCH resource pool memory limit in MB (unlimited=0)",
			Destination: &c.VCHMemoryLimitsMB,
		},
		cli.IntFlag{
			Name:        "memory-reservation, memr",
			Value:       0,
			Usage:       "VCH resource pool memory reservation in MB",
			Destination: &c.VCHMemoryReservationsMB,
			Hidden:      true,
		},
		cli.GenericFlag{
			Name:   "memory-shares, mems",
			Value:  flags.NewSharesFlag(&c.VCHMemoryShares),
			Usage:  "VCH resource pool memory shares in level or share number, e.g. high, normal, low, or 163840",
			Hidden: true,
		},
		cli.IntFlag{
			Name:        "endpoint-memory",
			Value:       2048,
			Usage:       "Memory for the VCH endpoint VM, in MB. Does not impact resources allocated per container.",
			Hidden:      true,
			Destination: &c.MemoryMB,
		},

		// cpu
		cli.IntFlag{
			Name:        "cpu",
			Value:       0,
			Usage:       "VCH resource pool vCPUs limit in MHz (unlimited=0)",
			Destination: &c.VCHCPULimitsMHz,
		},
		cli.IntFlag{
			Name:        "cpu-reservation, cpur",
			Value:       0,
			Usage:       "VCH resource pool reservation in MHz",
			Destination: &c.VCHCPUReservationsMHz,
			Hidden:      true,
		},
		cli.GenericFlag{
			Name:   "cpu-shares, cpus",
			Value:  flags.NewSharesFlag(&c.VCHCPUShares),
			Usage:  "VCH VCH resource pool vCPUs shares, in level or share number, e.g. high, normal, low, or 4000",
			Hidden: true,
		},
		cli.IntFlag{
			Name:        "endpoint-cpu",
			Value:       1,
			Usage:       "vCPUs for the VCH endpoint VM. Does not impact resources allocated per container.",
			Hidden:      true,
			Destination: &c.NumCPUs,
		},

		// TLS
		cli.StringFlag{
			Name:        "tls-cname",
			Value:       "",
			Usage:       "Common Name to use in generated CA certificate when requiring client certificate authentication",
			Destination: &c.cname,
		},
		cli.StringSliceFlag{
			Name:   "organization",
			Usage:  "A list of identifiers to record in the generated certificates. Defaults to VCH name and IP/FQDN if provided.",
			Value:  &c.org,
			Hidden: true,
		},
		cli.BoolFlag{
			Name:        "no-tlsverify, kv",
			Usage:       "Disable authentication via client certificates - for more tls options see advanced help (-x)",
			Destination: &c.noTLSverify,
		},
		cli.BoolFlag{
			Name:        "no-tls, k",
			Usage:       "Disable TLS support completely",
			Destination: &c.noTLS,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "key",
			Value:       "",
			Usage:       "Virtual Container Host private key file (server certificate)",
			Destination: &c.skey,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "cert",
			Value:       "",
			Usage:       "Virtual Container Host x509 certificate file (server certificate)",
			Destination: &c.scert,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "cert-path",
			Value:       "",
			Usage:       "The path to check for existing certificates and in which to save generated certificates. Defaults to './<vch name>/'",
			Destination: &c.certPath,
			Hidden:      true,
		},
		cli.StringSliceFlag{
			Name:   "tls-ca, ca",
			Usage:  "Specify a list of certificate authority files to use for client verification",
			Value:  &c.clientCAs,
			Hidden: true,
		},
		cli.IntFlag{
			Name:        "certificate-key-size, ksz",
			Usage:       "Size of key to use when generating certificates",
			Value:       2048,
			Destination: &c.keySize,
			Hidden:      true,
		},

		// registries
		cli.StringSliceFlag{
			Name:   "registry-ca, rc",
			Usage:  "Specify a list of additional certificate authority files to use to verify secure registry servers",
			Value:  &c.registryCAs,
			Hidden: true,
		},
		cli.StringSliceFlag{
			Name:  "insecure-registry, dir",
			Value: &c.insecureRegistries,
			Usage: "Specify a list of permitted insecure registry server URLs",
		},

		// proxies
		cli.StringFlag{
			Name:        "https-proxy, sproxy",
			Value:       "",
			Usage:       "An HTTPS proxy for use when fetching images, in the form https://fqdn_or_ip:port",
			Destination: &c.httpsProxy,
			Hidden:      true,
		},

		cli.StringFlag{
			Name:        "http-proxy, hproxy",
			Value:       "",
			Usage:       "An HTTP proxy for use when fetching images, in the form http://fqdn_or_ip:port",
			Destination: &c.httpProxy,
			Hidden:      true,
		},
	}

	util := []cli.Flag{
		// miscellaneous
		cli.BoolFlag{
			Name:        "use-rp",
			Usage:       "Use resource pool for vch parent in VC instead of a vApp",
			Destination: &c.UseRP,
			Hidden:      true,
		},

		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Force the install, removing existing if present",
			Destination: &c.Force,
		},
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for create",
			Destination: &c.Timeout,
		},
		cli.BoolFlag{
			Name:        "asymmetric-routes, ar",
			Usage:       "Set up the Virtual Container Host for asymmetric routing",
			Destination: &c.AsymmetricRouting,
			Hidden:      true,
		},
	}

	help := []cli.Flag{
		// help options
		cli.BoolFlag{
			Name:        "extended-help, x",
			Usage:       "Show all options - this must be specified instead of --help",
			Destination: &c.advancedOptions,
		},
	}

	target := c.TargetFlags()
	compute := c.ComputeFlags()
	iso := c.ImageFlags(true)
	debug := c.DebugFlags()

	// flag arrays are declared, now combined
	var flags []cli.Flag
	for _, f := range [][]cli.Flag{target, compute, create, iso, util, debug, help} {
		flags = append(flags, f...)
	}

	return flags
}

func (c *Create) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := c.HasCredentials(); err != nil {
		return err
	}

	// prevent usage of special characters for certain user provided values
	if err := common.CheckUnsupportedChars(c.DisplayName); err != nil {
		return cli.NewExitError(fmt.Sprintf("--name contains unsupported characters: %s Allowed characters are alphanumeric, space and symbols - _ ( )", err), 1)
	}

	if len(c.DisplayName) > MaxDisplayNameLen {
		return cli.NewExitError(fmt.Sprintf("Display name %s exceeds the permitted 31 characters limit. Please use a shorter -name parameter", c.DisplayName), 1)
	}

	if c.BridgeNetworkName == "" {
		c.BridgeNetworkName = c.DisplayName
	}

	if err := c.processOpsCredentials(); err != nil {
		return err
	}

	if err := c.processContainerNetworks(); err != nil {
		return err
	}

	if err := c.processBridgeNetwork(); err != nil {
		return err
	}

	if err := c.processNetwork(&c.Data.ClientNetwork, "client", c.clientNetworkName,
		c.clientNetworkIP, c.clientNetworkGateway); err != nil {
		return err
	}

	if err := c.processNetwork(&c.Data.PublicNetwork, "public", c.publicNetworkName,
		c.publicNetworkIP, c.publicNetworkGateway); err != nil {
		return err
	}

	if err := c.processNetwork(&c.Data.ManagementNetwork, "management", c.managementNetworkName,
		c.managementNetworkIP, c.managementNetworkGateway); err != nil {
		return err
	}

	if err := c.processDNSServers(); err != nil {
		return err
	}

	// must come after client network processing as it checks for static IP on that interface
	if err := c.processCertificates(); err != nil {
		return err
	}

	if err := common.CheckUnsupportedCharsDatastore(c.ImageDatastorePath); err != nil {
		return cli.NewExitError(fmt.Sprintf("--image-store contains unsupported characters: %s Allowed characters are alphanumeric, space and symbols - _ ( ) / :", err), 1)
	}

	if err := c.processVolumeStores(); err != nil {
		return errors.Errorf("Error occurred while processing volume stores: %s", err)
	}

	if err := c.processRegistries(); err != nil {
		return err
	}

	if err := c.processProxies(); err != nil {
		return err
	}

	return nil
}

func (c *Create) processOpsCredentials() error {
	if c.OpsUser == "" && c.OpsPassword != nil {
		return errors.New("Password for operations user specified without user having been specified")
	}

	if c.OpsUser == "" {
		log.Warn("Using administrative user for VCH operation - use --ops-user to improve security (see -x for advanced help)")
		c.OpsUser = c.Target.User
		if c.Target.Password == nil {
			return errors.New("Unable to use nil password from administrative user for operations user")
		}

		c.OpsPassword = c.Target.Password
		return nil
	}

	if c.OpsPassword != nil {
		return nil
	}

	log.Printf("vSphere password for %s: ", c.OpsUser)
	b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		message := fmt.Sprintf("Failed to read password from stdin: %s", err)
		cli.NewExitError(message, 1)
	}
	sb := string(b)
	c.OpsPassword = &sb

	return nil
}

func (c *Create) processCertificates() error {
	// set up the locations for the certificates and env file
	if c.certPath == "" {
		c.certPath = c.DisplayName
	}
	c.envFile = fmt.Sprintf("%s/%s.env", c.certPath, c.DisplayName)

	// check for insecure case
	if c.noTLS {
		log.Warn("Configuring without TLS - all communications will be insecure")
		return nil
	}

	if c.scert != "" && c.skey == "" {
		return cli.NewExitError("key and cert should be specified at the same time", 1)
	}
	if c.scert == "" && c.skey != "" {
		return cli.NewExitError("key and cert should be specified at the same time", 1)
	}

	// if we've not got a specific CommonName but do have a static IP then go with that.
	if c.cname == "" {
		if c.clientNetworkIP != "" {
			c.cname = c.clientNetworkIP
			log.Infof("Using client-network-ip as cname where needed - use --tls-cname to override: %s", c.cname)
		} else if c.publicNetworkIP != "" && (c.publicNetworkName == c.clientNetworkName || c.clientNetworkName == "") {
			c.cname = c.publicNetworkIP
			log.Infof("Using public-network-ip as cname where needed - use --tls-cname to override: %s", c.cname)
		} else if c.managementNetworkIP != "" && (c.managementNetworkName == c.clientNetworkName || (c.clientNetworkName == "" && c.managementNetworkName == c.publicNetworkName)) {
			c.cname = c.managementNetworkIP
			log.Infof("Using management-network-ip as cname where needed - use --tls-cname to override: %s", c.cname)
		}

		if c.cname != "" {
			// Strip network mask from IP address if set
			if cnameIP, _, _ := net.ParseCIDR(c.cname); cnameIP != nil {
				c.cname = cnameIP.String()
			}
		}
	}

	// load what certificates we can
	cas, keypair, err := c.loadCertificates()
	if err != nil {
		log.Errorf("Unable to load certificates: %s", err)
		if !c.Force {
			return err
		}

		log.Warnf("Ignoring error loading certificates due to --force")
		cas = nil
		keypair = nil
		err = nil
	}

	// we need to generate some part of the certificate configuration
	gcas, gkeypair, err := c.generateCertificates(keypair == nil, !c.noTLSverify && len(cas) == 0)
	if err != nil {
		log.Error("Create cannot continue: unable to generate certificates")
		return err
	}

	if keypair != nil {
		c.KeyPEM = keypair.KeyPEM
		c.CertPEM = keypair.CertPEM
	} else if gkeypair != nil {
		c.KeyPEM = gkeypair.KeyPEM
		c.CertPEM = gkeypair.CertPEM
	}

	if len(cas) == 0 {
		cas = gcas
	}

	if len(c.KeyPEM) == 0 {
		return errors.New("Failed to load or generate server certificates")
	}

	if len(cas) == 0 && !c.noTLSverify {
		return errors.New("Failed to load or generate certificate authority")
	}

	// do we have key, cert, and --no-tlsverify
	if c.noTLSverify || len(cas) == 0 {
		log.Warnf("Configuring without TLS verify - certificate-based authentication disabled")
		return nil
	}

	c.ClientCAs = cas
	return nil
}

func (c *Create) processBridgeNetwork() error {
	// bridge network params
	var err error

	_, c.Data.BridgeIPRange, err = net.ParseCIDR(c.BridgeIPRange)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error parsing bridge network ip range: %s. Range must be in CIDR format, e.g., 172.16.0.0/12", err), 1)
	}
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

		delete(gws, vnet)
		delete(pools, vnet)
		delete(dns, vnet)
	}

	var hasError bool
	fmtMsg := "The following container network %s is set, but CONTAINER-NETWORK cannot be found. Please check the --container-network and %s settings"
	if len(gws) > 0 {
		log.Error(fmt.Sprintf(fmtMsg, "gateway", "--container-network-gateway"))
		for key, value := range gws {
			mask, _ := value.Mask.Size()
			log.Errorf("\t%s:%s/%d, %q should be vSphere network name", key, value.IP, mask, key)
		}
		hasError = true
	}
	if len(pools) > 0 {
		log.Error(fmt.Sprintf(fmtMsg, "ip range", "--container-network-ip-range"))
		for key, value := range pools {
			log.Errorf("\t%s:%s, %q should be vSphere network name", key, value, key)
		}
		hasError = true
	}
	if len(dns) > 0 {
		log.Errorf(fmt.Sprintf(fmtMsg, "dns", "--container-network-dns"))
		for key, value := range dns {
			log.Errorf("\t%s:%s, %q should be vSphere network name", key, value, key)
		}
		hasError = true
	}
	if hasError {
		return cli.NewExitError("Inconsistent container network configuration.", 1)
	}
	return nil
}

func parseGatewaySpec(gw string) (cidrs []net.IPNet, gwIP net.IPNet, err error) {
	ss := strings.Split(gw, ":")
	if len(ss) > 2 {
		err = fmt.Errorf("gateway %s specified incorrectly", gw)
		return
	}

	gwStr := ss[0]
	cidrsStr := ""
	if len(ss) > 1 {
		gwStr = ss[1]
		cidrsStr = ss[0]
	}

	if gwIP.IP = net.ParseIP(gwStr); gwIP.IP == nil {
		err = fmt.Errorf("Provided gateway IP address is not valid: %s", gwStr)
	}

	if err != nil {
		return
	}

	if cidrsStr != "" {
		for _, c := range strings.Split(cidrsStr, ",") {
			var ipnet *net.IPNet
			_, ipnet, err = net.ParseCIDR(c)
			if err != nil {
				err = fmt.Errorf("invalid CIDR in gateway specification: %s", err)
				return
			}
			cidrs = append(cidrs, *ipnet)
		}
	}

	return
}

// processNetwork parses network args if present
func (c *Create) processNetwork(network *data.NetworkConfig, netName, pgName, staticIP, gateway string) error {
	var err error

	network.Name = pgName

	if gateway != "" && staticIP == "" {
		return fmt.Errorf("Gateway provided without static IP for %s network", netName)
	}

	defer func(net *data.NetworkConfig) {
		if err == nil {
			log.Debugf("%s network: IP %s gateway %s dest: %s", netName, net.IP, net.Gateway.IP, net.Destinations)
		}
	}(network)

	var ipNet *net.IPNet
	if staticIP != "" {
		var ipAddr net.IP
		ipAddr, ipNet, err = net.ParseCIDR(staticIP)
		if err != nil {
			return fmt.Errorf("Failed to parse the provided %s network IP address %s: %s", netName, staticIP, err)
		}

		network.IP.IP = ipAddr
		network.IP.Mask = ipNet.Mask
	}

	if gateway != "" {
		network.Destinations, network.Gateway, err = parseGatewaySpec(gateway)
		if err != nil {
			return fmt.Errorf("Invalid %s network gateway: %s", netName, err)
		}

		if !network.IP.Contains(network.Gateway.IP) {
			return fmt.Errorf("%s gateway with IP %s is not reachable from %s", netName, network.Gateway.IP, ipNet.String())
		}

		// TODO(vburenin): this seems ugly, and it actually is. The reason is that a gateway required to specify
		// a network mask for it, which is just not how network configuration should be done. Network mask has to
		// be provided separately or with the IP address. It is hard to change all dependencies to keep mask
		// with IP address, so it will be stored with network gateway as it was previously.
		network.Gateway.Mask = network.IP.Mask
	}

	return nil
}

// processDNSServers parses DNS servers used for client, public, mgmt networks
func (c *Create) processDNSServers() error {
	for _, d := range c.dns {
		s := net.ParseIP(d)
		if s == nil {
			return errors.New("Invalid DNS server specified")
		}
		c.Data.DNS = append(c.Data.DNS, s)
	}

	if len(c.Data.DNS) > 3 {
		log.Warn("Maximum of 3 DNS servers allowed. Additional servers specified will be ignored.")
	}

	log.Debugf("VCH DNS servers: %s", c.Data.DNS)
	return nil
}

func (c *Create) processVolumeStores() error {
	defer trace.End(trace.Begin(""))
	c.VolumeLocations = make(map[string]string)
	for _, arg := range c.volumeStores {
		if err := common.CheckUnsupportedCharsDatastore(arg); err != nil {
			return fmt.Errorf("--volume-store contains unsupported characters: %s Allowed characters are alphanumeric, space and symbols - _ ( ) / :", err)
		}
		splitMeta := strings.SplitN(arg, ":", 2)
		if len(splitMeta) != 2 {
			return errors.New("Volume store input must be in format datastore/path:label")
		}
		c.VolumeLocations[splitMeta[1]] = splitMeta[0]
	}

	return nil
}

func (c *Create) processRegistries() error {
	// load additional certificate authorities for use with registries
	if len(c.registryCAs) > 0 {
		registryCAs, err := c.loadRegistryCAs()
		if err != nil {
			return errors.Errorf("Unable to load CA certificates for registry logins: %s", err)
		}

		c.RegistryCAs = registryCAs
	}

	// load a list of insecure registries
	for _, registry := range c.insecureRegistries {
		regurl, err := validate.ParseURL(registry)

		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s is an invalid format for registry url", registry), 1)
		}
		c.InsecureRegistries = append(c.InsecureRegistries, *regurl)
	}

	return nil
}

func (c *Create) processProxies() error {
	var err error
	if c.httpProxy != "" {
		c.HTTPProxy, err = url.Parse(c.httpProxy)
		if err != nil || c.HTTPProxy.Host == "" || c.HTTPProxy.Scheme != "http" {
			return cli.NewExitError(fmt.Sprintf("Could not parse HTTP proxy - expected format http://fqnd_or_ip:port: %s", c.httpProxy), 1)
		}
	}

	if c.httpsProxy != "" {
		c.HTTPSProxy, err = url.Parse(c.httpsProxy)
		if err != nil || c.HTTPSProxy.Host == "" || c.HTTPSProxy.Scheme != "https" {
			return cli.NewExitError(fmt.Sprintf("Could not parse HTTPS proxy - expected format https://fqnd_or_ip:port: %s", c.httpsProxy), 1)
		}
	}

	return nil
}

// loadCertificates returns the client CA pool and the keypair for server certificates on success
func (c *Create) loadCertificates() ([]byte, *certificate.KeyPair, error) {
	defer trace.End(trace.Begin(""))

	// reads each of the files specified, assuming that they are PEM encoded certs,
	// and constructs a byte array suitable for passing to CertPool.AppendCertsFromPEM
	var certs []byte
	for _, f := range c.clientCAs {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			err = errors.Errorf("Failed to load authority from file %s: %s", f, err)
			return nil, nil, err
		}

		certs = append(certs, b...)
		log.Infof("Loaded CA from %s", f)
	}

	var keypair *certificate.KeyPair
	// default names
	skey := filepath.Join(c.certPath, serverKey)
	scert := filepath.Join(c.certPath, serverCert)
	ca := filepath.Join(c.certPath, caCert)
	ckey := filepath.Join(c.certPath, clientKey)
	ccert := filepath.Join(c.certPath, clientCert)

	// if specific files are supplied, use those
	explicit := false
	if c.scert != "" && c.skey != "" {
		explicit = true
		skey = c.skey
		scert = c.scert
	}

	// load the server certificate
	keypair = certificate.NewKeyPair(scert, skey, nil, nil)
	if err := keypair.LoadCertificate(); err != nil {
		if explicit || !os.IsNotExist(err) {
			// if these files were explicit paths, or anything other than not found, fail
			log.Errorf("Failed to load certificate: %s", err)
			return certs, nil, err
		}

		log.Debugf("Unable to locate existing server certificate in cert path")
		return nil, nil, nil
	}

	// check that any supplied cname matches the server cert CN
	cert, err := keypair.Certificate()
	if err != nil {
		log.Errorf("Failed to parse certificate: %s", err)
		return certs, nil, err
	}

	if cert.Leaf == nil {
		log.Warnf("Failed to load x509 leaf: Unable to confirm server certificate cname matches provided cname %q. Continuing...", c.cname)
	} else {
		// We just do a direct equality check here - trying to be clever is liable to lead to hard
		// to diagnose errors
		if cert.Leaf.Subject.CommonName != c.cname {
			log.Errorf("Provided cname does not match that in existing server certificate: %s", cert.Leaf.Subject.CommonName)
			if c.Debug.Debug > 2 {
				log.Debugf("Certificate does not match provided cname: %#+v", cert.Leaf)
			}
			return certs, nil, fmt.Errorf("cname option doesn't match existing server certificate in certificate path %s", c.certPath)
		}
	}

	log.Infof("Loaded server certificate %s", scert)
	c.skey = skey
	c.scert = scert

	// only try for CA certificate if no-tlsverify has NOT been specified and we haven't already loaded an authority cert
	if !c.noTLSverify && len(certs) == 0 {
		b, err := ioutil.ReadFile(ca)
		if err != nil {
			if os.IsNotExist(err) {
				log.Debugf("Unable to locate existing CA in cert path")
				return certs, keypair, nil
			}

			// if the CA exists but cannot be loaded then it's an error
			log.Errorf("Failed to load authority from certificate path %s: %s", c.certPath, err)
			return certs, keypair, errors.New("failed to load certificate authority")
		}

		c.cacert = ca

		log.Infof("Loaded CA with default name from certificate path %s", c.certPath)
		certs = b

		// load client certs - we ensure the client certs validate with the provided CA or ignore any we find
		cpair := certificate.NewKeyPair(ccert, ckey, nil, nil)
		if err := cpair.LoadCertificate(); err != nil {
			log.Warnf("Unable to load client certificate - validation of API endpoint will be best effort only: %s", err)
		}

		clientCert, err := cpair.Certificate()
		if err != nil || clientCert.Leaf == nil {
			log.Debugf("Unable to parse client certificate: %s", err)
			return certs, keypair, nil
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(certs) {
			log.Debugf("Unable to create CA pool to check client certificate")
			return certs, keypair, nil
		}

		opts := x509.VerifyOptions{
			Roots:     pool,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		}

		_, err = clientCert.Leaf.Verify(opts)
		if err != nil {
			log.Warnf("Client certificate in certificate path does not validate with provided CA - continuing without client certificate")
			return certs, keypair, nil
		}

		c.ckey = ckey
		c.ccert = ccert
		c.clientCert = clientCert

		log.Infof("Loaded client certificate with default name from certificate path %s", c.certPath)
	}

	return certs, keypair, nil
}

// loadRegistryCAs loads additional CA certs for docker registry usage
func (c *Create) loadRegistryCAs() ([]byte, error) {
	defer trace.End(trace.Begin(""))

	var registryCerts []byte
	for _, f := range c.registryCAs {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			err = errors.Errorf("Failed to load authority from file %s: %s", f, err)
			return nil, err
		}

		registryCerts = append(registryCerts, b...)
		log.Infof("Loaded registry CA from %s", f)
	}

	return registryCerts, nil
}

func (c *Create) generateCertificates(server bool, client bool) ([]byte, *certificate.KeyPair, error) {
	defer trace.End(trace.Begin(""))

	if !server && !client {
		log.Debugf("Not generating server or client certs, nothing for generateCertificates to do")
		return nil, nil, nil
	}

	var certs []byte
	// generate the certs and keys with names conforming the default the docker client expects
	files, err := ioutil.ReadDir(c.certPath)
	if len(files) > 0 {
		return nil, nil, fmt.Errorf("Specified directory to store certificates is not empty. Specify a new path in which to store generated certificates using --cert-path or remove the contents of \"%s\" and run vic-machine again.", c.certPath)
	}

	err = os.MkdirAll(c.certPath, 0700)
	if err != nil {
		log.Errorf("Unable to make directory to hold certificates (set via --cert-path)")
		return nil, nil, err
	}

	c.skey = filepath.Join(c.certPath, serverKey)
	c.scert = filepath.Join(c.certPath, serverCert)

	c.ckey = filepath.Join(c.certPath, clientKey)
	c.ccert = filepath.Join(c.certPath, clientCert)

	cakey := filepath.Join(c.certPath, caKey)
	c.cacert = filepath.Join(c.certPath, caCert)

	if server && !client {
		log.Infof("Generating self-signed certificate/key pair - private key in %s", c.skey)
		keypair := certificate.NewKeyPair(c.scert, c.skey, nil, nil)
		err := keypair.CreateSelfSigned(c.cname, nil, c.keySize)
		if err != nil {
			log.Errorf("Failed to generate self-signed certificate: %s", err)
			return nil, nil, err
		}
		if err = keypair.SaveCertificate(); err != nil {
			log.Errorf("Failed to save server certificates: %s", err)
			return nil, nil, err
		}

		return certs, keypair, nil
	}

	// client auth path
	if c.cname == "" {
		log.Error("Common Name must be provided when generating certificates for client authentication:")
		log.Info("  --tls-cname=<FQDN or static IP> # for the appliance VM")
		log.Info("  --tls-cname=<*.yourdomain.com>  # if DNS has entries in that form for DHCP addresses (less secure)")
		log.Info("  --no-tlsverify                  # disables client authentication (anyone can connect to the VCH)")
		log.Info("  --no-tls                        # disables TLS entirely")
		log.Info("")

		return certs, nil, errors.New("provide Common Name for server certificate")
	}

	// for now re-use the display name as the organisation if unspecified
	if len(c.org) == 0 {
		c.org = []string{c.DisplayName}
	}
	if len(c.org) == 1 && !strings.HasPrefix(c.cname, "*") {
		// Add in the cname if it's not a wildcard
		c.org = append(c.org, c.cname)
	}

	// Certificate authority
	log.Infof("Generating CA certificate/key pair - private key in %s", cakey)
	cakp := certificate.NewKeyPair(c.cacert, cakey, nil, nil)
	err = cakp.CreateRootCA(c.cname, c.org, c.keySize)
	if err != nil {
		log.Errorf("Failed to generate CA: %s", err)
		return nil, nil, err
	}
	if err = cakp.SaveCertificate(); err != nil {
		log.Errorf("Failed to save CA certificates: %s", err)
		return nil, nil, err
	}

	// Server certificates
	var skp *certificate.KeyPair
	if server {
		log.Infof("Generating server certificate/key pair - private key in %s", c.skey)
		skp = certificate.NewKeyPair(c.scert, c.skey, nil, nil)
		err = skp.CreateServerCertificate(c.cname, c.org, c.keySize, cakp)
		if err != nil {
			log.Errorf("Failed to generate server certificates: %s", err)
			return nil, nil, err
		}
		if err = skp.SaveCertificate(); err != nil {
			log.Errorf("Failed to save server certificates: %s", err)
			return nil, nil, err
		}
	}

	// Client certificates
	if client {
		log.Infof("Generating client certificate/key pair - private key in %s", c.ckey)
		ckp := certificate.NewKeyPair(c.ccert, c.ckey, nil, nil)
		err = ckp.CreateClientCertificate(c.cname, c.org, c.keySize, cakp)
		if err != nil {
			log.Errorf("Failed to generate client certificates: %s", err)
			return nil, nil, err
		}
		if err = ckp.SaveCertificate(); err != nil {
			log.Errorf("Failed to save client certificates: %s", err)
			return nil, nil, err
		}

		c.clientCert, err = ckp.Certificate()
		if err != nil {
			log.Warnf("Failed to stash client certificate for later application level validation: %s", err)
		}

		// If openssl is present, try to generate a browser friendly pfx file (a bundle of the public certificate AND the private key)
		// The pfx file can be imported directly into keychains for client certificate authentication
		certPath := filepath.Clean(c.certPath)
		args := strings.Split(fmt.Sprintf("pkcs12 -export -out %[1]s/cert.pfx -inkey %[1]s/key.pem -in %[1]s/cert.pem -certfile %[1]s/ca.pem -password pass:", certPath), " ")
		// #nosec: Subprocess launching with variable
		pfx := exec.Command("openssl", args...)
		out, err := pfx.CombinedOutput()
		if err != nil {
			log.Debug(out)
			log.Warnf("Failed to generate browser friendly PFX client certificate: %s", err)
		} else {
			log.Infof("Generated browser friendly PFX client certificate - certificate in %s/cert.pfx", certPath)
		}
	}

	return cakp.CertPEM, skp, nil
}

func logArguments(cliContext *cli.Context) {
	for _, f := range cliContext.FlagNames() {
		if !cliContext.IsSet(f) {
			continue
		}

		// avoid logging sensitive data
		if f == "user" || f == "password" || f == "ops-user" || f == "ops-password" {
			log.Debugf("--%s=<censored>", f)
			continue
		}

		if f == "target" {
			url, err := url.Parse(cliContext.String(f))
			if err != nil {
				log.Debugf("Unable to re-parse target url for logging")
				continue
			}
			url.User = nil
			log.Debugf("--target=%s", url.String())
			continue
		}

		i := cliContext.Int(f)
		if i != 0 {
			log.Debugf("--%s=%d", f, i)
			continue
		}
		d := cliContext.Duration(f)
		if d != 0 {
			log.Debugf("--%s=%s", f, d.String())
			continue
		}
		x := cliContext.Float64(f)
		if x != 0 {
			log.Debugf("--%s=%f", f, x)
			continue
		}
		s := cliContext.String(f)
		if s != "" {
			log.Debugf("--%s=%s", f, s)
			continue
		}
		b := cliContext.Bool(f)
		bT := cliContext.BoolT(f)
		if b && !bT {
			log.Debugf("--%s=%t", f, true)
			continue
		}

		// put the slices at the end as they cause panics
		match := func() (result bool) {
			result = false
			defer func() { recover() }()
			ss := cliContext.StringSlice(f)
			if ss != nil {
				log.Debugf("--%s=%#v", f, ss)
				return true
			}
			return
		}()
		if match {
			continue
		}

		match = func() (result bool) {
			result = false
			defer func() { recover() }()
			is := cliContext.IntSlice(f)
			if is != nil {
				log.Debugf("--%s=%#v", f, is)
				return true
			}
			return
		}()
		if match {
			continue
		}

		// generic last because it matches everything
		g := cliContext.Generic(f)
		if g != nil {
			log.Debugf("--%s=%#v", f, g)
			continue
		}
	}
}

func (c *Create) Run(clic *cli.Context) (err error) {

	if c.advancedOptions {
		cli.HelpPrinter(clic.App.Writer, EntireOptionHelpTemplate, clic.Command)
		return nil
	}

	log.Infof("### Installing VCH ####")

	if c.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	// urfave/cli will print out exit in error handling, so no more information in main method can be printed out.
	defer func() {
		err = common.LogErrorIfAny(clic, err)
	}()

	if err = c.processParams(); err != nil {
		return err
	}

	logArguments(clic)

	var images map[string]string
	if images, err = c.CheckImagesFiles(c.Force); err != nil {
		return err
	}

	if len(clic.Args()) > 0 {
		log.Errorf("Unknown argument: %s", clic.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	// all these operations will be executed without timeout
	ctx := context.Background()
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

	vConfig := validator.AddDeprecatedFields(ctx, vchConfig, c.Data)
	vConfig.ImageFiles = images
	vConfig.ApplianceISO = path.Base(c.ApplianceISO)
	vConfig.BootstrapISO = path.Base(c.BootstrapISO)

	vConfig.HTTPProxy = c.HTTPProxy
	vConfig.HTTPSProxy = c.HTTPSProxy

	vConfig.Timeout = c.Data.Timeout

	// separate initial validation from dispatch of creation task
	log.Info("")

	executor := management.NewDispatcher(ctx, validator.Session, vchConfig, c.Force)
	if err = executor.CreateVCH(vchConfig, vConfig); err != nil {
		executor.CollectDiagnosticLogs()
		log.Error(err)
		return err
	}

	// timeoout start to work from here, to make sure user does not wait forever
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	defer func() {
		if ctx.Err() == context.DeadlineExceeded {
			//context deadline exceeded, replace returned error message
			err = errors.Errorf("Creating VCH exceeded time limit of %s. Please increase the timeout using --timeout to accommodate for a busy vSphere target", c.Timeout)
		}
	}()

	if err = executor.CheckServiceReady(ctx, vchConfig, c.clientCert); err != nil {
		executor.CollectDiagnosticLogs()
		cmd, _ := executor.GetDockerAPICommand(vchConfig, c.ckey, c.ccert, c.cacert, c.certPath)
		log.Info("\tAPI may be slow to start - try to connect to API after a few minutes:")
		if cmd != "" {
			log.Infof("\t\tRun command: %s", cmd)
		} else {
			log.Infof("\t\tRun %s inspect to find API connection command and run the command if ip address is ready", clic.App.Name)
		}
		log.Info("\t\tIf command succeeds, VCH is started. If command fails, VCH failed to install - see documentation for troubleshooting.")
		return err
	}

	log.Infof("Initialization of appliance successful")

	executor.ShowVCH(vchConfig, c.ckey, c.ccert, c.cacert, c.envFile, c.certPath)
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
