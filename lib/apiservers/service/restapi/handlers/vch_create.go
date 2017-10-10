// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package handlers

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"net/url"
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-units"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/cmd/vic-machine/create"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/util"
	"github.com/vmware/vic/lib/apiservers/service/restapi/operations"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/lib/install/vchlog"
	"github.com/vmware/vic/pkg/ip"
	viclog "github.com/vmware/vic/pkg/log"
	"github.com/vmware/vic/pkg/version"
)

const (
	LogFile = "vic-machine.log" // name of local log file
)

// VCHCreate is the handler for creating a VCH
type VCHCreate struct {
}

// VCHCreate is the handler for creating a VCH within a Datacenter
type VCHDatacenterCreate struct {
}

func (h *VCHCreate) Handle(params operations.PostTargetTargetVchParams, principal interface{}) middleware.Responder {
	// Set up VCH create logger
	localLogFile := setUpLogger()
	// Close the two logging streams when done
	defer vchlog.Close()
	defer localLogFile.Close()

	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		nil,
		nil)
	if err != nil {
		return operations.NewPostTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	validator, err := validateTarget(params.HTTPRequest.Context(), d)
	if err != nil {
		return operations.NewPostTargetTargetVchDefault(400).WithPayload(&models.Error{Message: err.Error()})
	}

	c, err := buildCreate(params.HTTPRequest.Context(), d, validator.Session.Finder, params.Vch)
	if err != nil {
		return operations.NewPostTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	task, err := handleCreate(params.HTTPRequest.Context(), c, validator)
	if err != nil {
		return operations.NewPostTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPostTargetTargetVchCreated().WithPayload(operations.PostTargetTargetVchCreatedBody{Task: task})
}

func (h *VCHDatacenterCreate) Handle(params operations.PostTargetTargetDatacenterDatacenterVchParams, principal interface{}) middleware.Responder {
	// Set up VCH create logger
	localLogFile := setUpLogger()
	// Close the two logging streams when done
	defer vchlog.Close()
	defer localLogFile.Close()

	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		&params.Datacenter,
		nil)
	if err != nil {
		return operations.NewPostTargetTargetDatacenterDatacenterVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	validator, err := validateTarget(params.HTTPRequest.Context(), d)
	if err != nil {
		return operations.NewPostTargetTargetDatacenterDatacenterVchDefault(400).WithPayload(&models.Error{Message: err.Error()})
	}

	c, err := buildCreate(params.HTTPRequest.Context(), d, validator.Session.Finder, params.Vch)
	if err != nil {
		return operations.NewPostTargetTargetDatacenterDatacenterVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	task, err := handleCreate(params.HTTPRequest.Context(), c, validator)
	if err != nil {
		return operations.NewPostTargetTargetDatacenterDatacenterVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPostTargetTargetDatacenterDatacenterVchCreated().WithPayload(operations.PostTargetTargetDatacenterDatacenterVchCreatedBody{Task: task})
}

func setUpLogger() *os.File {
	vchlog.Init()
	logs := []io.Writer{}

	// Write to local log file
	// #nosec: Expect file permissions to be 0600 or less
	localLogFile, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		logs = append(logs, localLogFile)
	}

	// Also write logs to pipe streaming to VCH datastore
	logs = append(logs, vchlog.GetPipe())
	// Set log level to debug
	log.SetLevel(log.DebugLevel)
	// Initiliaze logger with default TextFormatter
	log.SetFormatter(viclog.NewTextFormatter())
	// SetOutput to io.MultiWriter so that we can log to stdout and a file
	log.SetOutput(io.MultiWriter(logs...))

	// Fire the logger
	go vchlog.Run()

	return localLogFile
}

func buildCreate(ctx context.Context, d *data.Data, finder *find.Finder, vch *models.VCH) (*create.Create, error) {
	c := &create.Create{Data: d}

	// TODO: deduplicate with create.processParams

	if vch != nil {
		if vch.Version != "" && version.String() != string(vch.Version) {
			return nil, util.NewError(400, fmt.Sprintf("Invalid version: %s", vch.Version))
		}

		c.DisplayName = vch.Name

		// TODO: move validation to swagger
		if err := common.CheckUnsupportedChars(c.DisplayName); err != nil {
			return nil, util.NewError(400, fmt.Sprintf("Invalid display name: %s", err))
		}
		if len(c.DisplayName) > create.MaxDisplayNameLen {
			return nil, util.NewError(400, fmt.Sprintf("Invalid display name: length exceeds %d characters", create.MaxDisplayNameLen))
		}

		debug := int(vch.Debug)
		c.Debug.Debug = &debug

		if vch.Compute != nil {
			if vch.Compute.CPU != nil {
				c.ResourceLimits.VCHCPULimitsMHz = mhzFromValueHertz(vch.Compute.CPU.Limit)
				c.ResourceLimits.VCHCPUReservationsMHz = mhzFromValueHertz(vch.Compute.CPU.Reservation)
				c.ResourceLimits.VCHCPUShares = fromShares(vch.Compute.CPU.Shares)
			}

			if vch.Compute.Memory != nil {
				c.ResourceLimits.VCHMemoryLimitsMB = mbFromValueBytes(vch.Compute.Memory.Limit)
				c.ResourceLimits.VCHMemoryReservationsMB = mbFromValueBytes(vch.Compute.Memory.Reservation)
				c.ResourceLimits.VCHMemoryShares = fromShares(vch.Compute.Memory.Shares)
			}

			resourcePath, err := fromManagedObject(ctx, finder, "ResourcePool", vch.Compute.Resource) // TODO: Do we need to handle clusters differently?
			if err != nil {
				return nil, util.NewError(400, fmt.Sprintf("Error finding resource pool: %s", err))
			}
			if resourcePath == "" {
				return nil, util.NewError(400, "Resource pool must be specified (by name or id)")
			}
			c.ComputeResourcePath = resourcePath
		}

		if vch.Network != nil {
			if vch.Network.Bridge != nil {
				path, err := fromManagedObject(ctx, finder, "Network", vch.Network.Bridge.PortGroup)
				if err != nil {
					return nil, util.NewError(400, fmt.Sprintf("Error finding bridge network: %s", err))
				}
				if path == "" {
					return nil, util.NewError(400, "Bridge network portgroup must be specified (by name or id)")
				}
				c.BridgeNetworkName = path
				c.BridgeIPRange = fromCIDR(&vch.Network.Bridge.IPRange.CIDR)

				if err := c.ProcessBridgeNetwork(); err != nil {
					return nil, util.WrapError(400, err)
				}
			}

			if vch.Network.Client != nil {
				path, err := fromManagedObject(ctx, finder, "Network", vch.Network.Client.PortGroup)
				if err != nil {
					return nil, util.NewError(400, fmt.Sprintf("Error finding client network portgroup: %s", err))
				}
				if path == "" {
					return nil, util.NewError(400, "Client network portgroup must be specified (by name or id)")
				}
				c.ClientNetworkName = path
				c.ClientNetworkGateway = fromGateway(vch.Network.Client.Gateway)
				c.ClientNetworkIP = fromNetworkAddress(vch.Network.Client.Static)

				if err := c.ProcessNetwork(&c.Data.ClientNetwork, "client", c.ClientNetworkName, c.ClientNetworkIP, c.ClientNetworkGateway); err != nil {
					return nil, util.WrapError(400, err)
				}
			}

			if vch.Network.Management != nil {
				path, err := fromManagedObject(ctx, finder, "Network", vch.Network.Management.PortGroup)
				if err != nil {
					return nil, util.NewError(400, fmt.Sprintf("Error finding management network portgroup: %s", err))
				}
				if path == "" {
					return nil, util.NewError(400, "Management network portgroup must be specified (by name or id)")
				}
				c.ManagementNetworkName = path
				c.ManagementNetworkGateway = fromGateway(vch.Network.Management.Gateway)
				c.ManagementNetworkIP = fromNetworkAddress(vch.Network.Management.Static)

				if err := c.ProcessNetwork(&c.Data.ManagementNetwork, "management", c.ManagementNetworkName, c.ManagementNetworkIP, c.ManagementNetworkGateway); err != nil {
					return nil, util.WrapError(400, err)
				}
			}

			if vch.Network.Public != nil {
				path, err := fromManagedObject(ctx, finder, "Network", vch.Network.Public.PortGroup)
				if err != nil {
					return nil, util.NewError(400, fmt.Sprintf("Error finding public network portgroup: %s", err))
				}
				if path == "" {
					return nil, util.NewError(400, "Public network portgroup must be specified (by name or id)")
				}
				c.PublicNetworkName = path
				c.PublicNetworkGateway = fromGateway(vch.Network.Public.Gateway)
				c.PublicNetworkIP = fromNetworkAddress(vch.Network.Public.Static)

				if err := c.ProcessNetwork(&c.Data.PublicNetwork, "public", c.PublicNetworkName, c.PublicNetworkIP, c.PublicNetworkGateway); err != nil {
					return nil, util.WrapError(400, err)
				}
			}

			if vch.Network.Container != nil {
				containerNetworks := common.ContainerNetworks{
					MappedNetworks:         make(map[string]string),
					MappedNetworksGateways: make(map[string]net.IPNet),
					MappedNetworksIPRanges: make(map[string][]ip.Range),
					MappedNetworksDNS:      make(map[string][]net.IP),
				}

				for _, cnetwork := range vch.Network.Container {
					alias := cnetwork.Alias

					path, err := fromManagedObject(ctx, finder, "Network", cnetwork.PortGroup)
					if err != nil {
						return nil, util.NewError(400, fmt.Sprintf("Error finding portgroup for container network %s: %s", alias, err))
					}
					if path == "" {
						return nil, util.NewError(400, fmt.Sprintf("Container network %s portgroup must be specified (by name or id)", alias))
					}
					containerNetworks.MappedNetworks[alias] = path

					address := net.ParseIP(string(cnetwork.Gateway.Address))
					_, mask, err := net.ParseCIDR(string(cnetwork.Gateway.RoutingDestinations[0].CIDR))
					if err != nil {
						return nil, util.NewError(400, fmt.Sprintf("Error parsing network mask for container network %s: %s", alias, err))
					}
					containerNetworks.MappedNetworksGateways[alias] = net.IPNet{
						IP:   address,
						Mask: mask.Mask,
					}

					ipranges := make([]ip.Range, 0, len(cnetwork.IPRanges))
					for _, ipRange := range cnetwork.IPRanges {
						r := ip.ParseRange(string(ipRange.CIDR))

						ipranges = append(ipranges, *r)
					}
					containerNetworks.MappedNetworksIPRanges[alias] = ipranges

					nameservers := make([]net.IP, 0, len(cnetwork.Nameservers))
					for _, nameserver := range cnetwork.Nameservers {
						n := net.ParseIP(string(nameserver))
						nameservers = append(nameservers, n)
					}
					containerNetworks.MappedNetworksDNS[alias] = nameservers
				}

				c.ContainerNetworks = containerNetworks
			}
		}

		if vch.Storage != nil {
			if vch.Storage.ImageStores != nil && len(vch.Storage.ImageStores) > 0 {
				c.ImageDatastorePath = vch.Storage.ImageStores[0] // TODO: many vs. one mismatch
			}

			if err := common.CheckUnsupportedCharsDatastore(c.ImageDatastorePath); err != nil {
				return nil, util.WrapError(400, err)
			}

			if vch.Storage.VolumeStores != nil {
				volumes := make([]string, 0, len(vch.Storage.VolumeStores))
				for _, v := range vch.Storage.VolumeStores {
					volumes = append(volumes, fmt.Sprintf("%s:%s", v.Datastore, v.Label))
				}

				vs := common.VolumeStores{VolumeStores: cli.StringSlice(volumes)}
				volumeLocations, err := vs.ProcessVolumeStores()
				if err != nil {
					return nil, util.NewError(400, fmt.Sprintf("Error processing volume stores: %s", err))
				}
				c.VolumeLocations = volumeLocations
			}

			c.ScratchSize = "8GB"
			if vch.Storage.BaseImageSize != nil {
				c.ScratchSize = fromValueBytes(vch.Storage.BaseImageSize)
			}
		}

		if vch.Auth != nil {
			c.Certs.NoTLS = vch.Auth.NoTLS

			if vch.Auth.Client != nil {
				c.Certs.NoTLSverify = vch.Auth.Client.NoTLSVerify
				c.Certs.ClientCAs = fromPemCertificates(vch.Auth.Client.CertificateAuthorities)
				c.ClientCAs = c.Certs.ClientCAs
			}

			if vch.Auth.Server != nil {

				if vch.Auth.Server.Generate != nil {
					c.Certs.Cname = vch.Auth.Server.Generate.Cname
					c.Certs.Org = vch.Auth.Server.Generate.Organization
					c.Certs.KeySize = fromValueBits(vch.Auth.Server.Generate.Size)

					if err := c.Certs.ProcessCertificates(c.DisplayName, c.Force, 0); err != nil {
						return nil, util.NewError(400, fmt.Sprintf("Error generating certificates: %s", err))
					}
				} else {
					c.Certs.CertPEM = []byte(vch.Auth.Server.Certificate.Pem)
					c.Certs.KeyPEM = []byte(vch.Auth.Server.PrivateKey.Pem)
				}

				c.CertPEM = c.Certs.CertPEM
				c.KeyPEM = c.Certs.KeyPEM
			}
		}

		c.MemoryMB = 2048
		if vch.Endpoint != nil {
			c.UseRP = vch.Endpoint.UseResourcePool
			if vch.Endpoint.Memory != nil {
				c.MemoryMB = *mbFromValueBytes(vch.Endpoint.Memory)
			}
			if vch.Endpoint.CPU != nil {
				c.NumCPUs = int(vch.Endpoint.CPU.Sockets)
			}

			if vch.Endpoint.OperationsCredentials != nil {
				opsPassword := string(vch.Endpoint.OperationsCredentials.Password)
				c.OpsCredentials = common.OpsCredentials{
					OpsUser:     &vch.Endpoint.OperationsCredentials.User,
					OpsPassword: &opsPassword,
				}
			}
			if err := c.OpsCredentials.ProcessOpsCredentials(true, c.Target.User, c.Target.Password); err != nil {
				return nil, util.WrapError(400, err)
			}
		}

		if vch.Registry != nil {
			c.InsecureRegistries = vch.Registry.Insecure
			c.WhitelistRegistries = vch.Registry.Whitelist

			//params.Vch.Registry.Blacklist

			c.RegistryCAs = fromPemCertificates(vch.Registry.CertificateAuthorities)

			if vch.Registry.ImageFetchProxy != nil {
				c.Proxies = fromImageFetchProxy(vch.Registry.ImageFetchProxy)
				_, _, err := c.Proxies.ProcessProxies()
				if err != nil {
					return nil, util.NewError(400, fmt.Sprintf("Error processing proxies: %s", err))
				}
			}
		}
	}

	return c, nil
}

func handleCreate(ctx context.Context, c *create.Create, validator *validate.Validator) (*strfmt.URI, error) {
	vchConfig, err := validator.Validate(validator.Context, c.Data)
	vConfig := validator.AddDeprecatedFields(validator.Context, vchConfig, c.Data)

	// TODO: make this configurable
	images := common.Images{}
	vConfig.ImageFiles, err = images.CheckImagesFiles(true)
	vConfig.ApplianceISO = path.Base(images.ApplianceISO)
	vConfig.BootstrapISO = path.Base(images.BootstrapISO)

	executor := management.NewDispatcher(validator.Context, validator.Session, nil, false)
	err = executor.CreateVCH(vchConfig, vConfig)
	if err != nil {
		return nil, util.NewError(500, fmt.Sprintf("Failed to create VCH: %s", err))
	}

	return nil, nil
}

func fromIPRanges(m *[]models.IPRange) *[]string {
	s := make([]string, 0, len(*m))
	for _, d := range *m {
		s = append(s, string(d.CIDR))
	}

	return &s
}

func fromNetworkAddress(m *models.NetworkAddress) string {
	if m == nil {
		return ""
	}

	if m.IP != "" {
		return string(m.IP)
	}

	return string(m.Hostname)
}

func fromManagedObject(ctx context.Context, finder *find.Finder, t string, m *models.ManagedObject) (string, error) {
	if m.ID != "" {
		managedObjectReference := types.ManagedObjectReference{Type: t, Value: m.ID}
		element, err := finder.Element(ctx, managedObjectReference)

		if err != nil {
			return "", err
		}

		return element.Path, nil
	}

	return m.Name, nil
}

func fromCIDR(m *models.CIDR) string {
	if m == nil {
		return ""
	}

	return string(*m)
}

func fromGateway(m *models.Gateway) string {
	if m == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", // TODO: what if RoutingDestinations is empty?
		strings.Join(*fromIPRanges(&m.RoutingDestinations), ","),
		m.Address,
	)
}

func fromValueBytes(m *models.ValueBytes) string {
	v := float64(m.Value.Value)

	var bytes float64
	switch m.Value.Units {
	case models.ValueBytesUnitsB:
		bytes = v
	case models.ValueBytesUnitsKiB:
		bytes = v * float64(units.KiB)
	case models.ValueBytesUnitsMiB:
		bytes = v * float64(units.MiB)
	case models.ValueBytesUnitsGiB:
		bytes = v * float64(units.GiB)
	case models.ValueBytesUnitsTiB:
		bytes = v * float64(units.TiB)
	case models.ValueBytesUnitsPiB:
		bytes = v * float64(units.PiB)
	}

	return fmt.Sprintf("%d B", int64(bytes))
}

func mbFromValueBytes(m *models.ValueBytes) *int {
	if m == nil {
		return nil
	}

	v := float64(m.Value.Value)

	var mbs float64
	switch m.Value.Units {
	case models.ValueBytesUnitsB:
		mbs = v / float64(units.MiB)
	case models.ValueBytesUnitsKiB:
		mbs = v / (float64(units.MiB) / float64(units.KiB))
	case models.ValueBytesUnitsMiB:
		mbs = v
	case models.ValueBytesUnitsGiB:
		mbs = v * (float64(units.GiB) / float64(units.MiB))
	case models.ValueBytesUnitsTiB:
		mbs = v * (float64(units.TiB) / float64(units.MiB))
	case models.ValueBytesUnitsPiB:
		mbs = v * (float64(units.PiB) / float64(units.MiB))
	}

	i := int(math.Ceil(mbs))

	return &i
}

func mhzFromValueHertz(m *models.ValueHertz) *int {
	if m == nil {
		return nil
	}

	v := float64(m.Value.Value)

	var mhzs float64
	switch m.Units {
	case models.ValueHertzUnitsHz:
		mhzs = v / float64(units.MB)
	case models.ValueHertzUnitsKHz:
		mhzs = v / (float64(units.MB) / float64(units.KB))
	case models.ValueHertzUnitsMHz:
		mhzs = v
	case models.ValueHertzUnitsGHz:
		mhzs = v * (float64(units.GB) / float64(units.MB))
	}

	i := int(math.Ceil(mhzs))

	return &i
}

func fromShares(m *models.Shares) *types.SharesInfo {
	if m == nil {
		return nil
	}

	var level types.SharesLevel
	switch types.SharesLevel(m.Level) {
	case types.SharesLevelLow:
		level = types.SharesLevelLow
	case types.SharesLevelNormal:
		level = types.SharesLevelNormal
	case types.SharesLevelHigh:
		level = types.SharesLevelHigh
	default:
		level = types.SharesLevelCustom
	}

	return &types.SharesInfo{
		Level:  level,
		Shares: int32(m.Number),
	}
}

func fromValueBits(m *models.ValueBits) int {
	return int(m.Value.Value)
}

func fromPemCertificates(m []*models.X509Data) []byte {
	var b []byte

	for _, ca := range m {
		c := []byte(ca.Pem)
		b = append(b, c...)
	}

	return b
}

func fromImageFetchProxy(p *models.VCHRegistryImageFetchProxy) common.Proxies {
	http := string(p.HTTP)
	https := string(p.HTTPS)

	return common.Proxies{
		HTTPProxy:  &http,
		HTTPSProxy: &https,
	}
}
