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
	"bytes"
	"context"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/docker/docker/opts"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/util"
	"github.com/vmware/vic/lib/apiservers/service/restapi/operations"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// VCHGet is the handler for inspecting a VCH
type VCHGet struct {
}

// VCHGet is the handler for inspecting a VCH within a Datacenter
type VCHDatacenterGet struct {
}

func (h *VCHGet) Handle(params operations.GetTargetTargetVchVchIDParams, principal interface{}) middleware.Responder {
	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		nil,
		nil)
	if err != nil {
		return operations.NewGetTargetTargetVchVchIDDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	d.ID = params.VchID

	vch, err := getVCH(params.HTTPRequest.Context(), d)

	if err != nil {
		return operations.NewGetTargetTargetVchVchIDDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewGetTargetTargetVchVchIDOK().WithPayload(vch)
}

func (h *VCHDatacenterGet) Handle(params operations.GetTargetTargetDatacenterDatacenterVchVchIDParams, principal interface{}) middleware.Responder {
	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		&params.Datacenter,
		nil)
	if err != nil {
		return operations.NewGetTargetTargetDatacenterDatacenterVchVchIDDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	d.ID = params.VchID

	vch, err := getVCH(params.HTTPRequest.Context(), d)

	if err != nil {
		return operations.NewGetTargetTargetDatacenterDatacenterVchVchIDDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewGetTargetTargetDatacenterDatacenterVchVchIDOK().WithPayload(vch)
}

func getVCH(ctx context.Context, d *data.Data) (*models.VCH, error) {
	validator, err := validateTarget(ctx, d)
	if err != nil {
		return nil, util.WrapError(400, err)
	}

	executor := management.NewDispatcher(validator.Context, validator.Session, nil, false)
	vch, err := executor.NewVCHFromID(d.ID)
	if err != nil {
		return nil, util.NewError(500, fmt.Sprintf("Failed to inspect VCH: %s", err))
	}

	err = validate.SetDataFromVM(validator.Context, validator.Session.Finder, vch, d)
	if err != nil {
		return nil, util.NewError(500, fmt.Sprintf("Failed to load VCH data: %s", err))
	}

	model, err := vchToModel(vch, d, executor)
	if err != nil {
		return nil, util.WrapError(500, err)
	}

	return model, nil
}

func vchToModel(vch *vm.VirtualMachine, d *data.Data, executor *management.Dispatcher) (*models.VCH, error) {
	vchConfig, err := executor.GetNoSecretVCHConfig(vch)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve VCH information: %s", err)
	}

	model := &models.VCH{}
	model.Version = models.Version(vchConfig.Version.ShortVersion())
	model.Name = vchConfig.Name

	// compute
	model.Compute = &models.VCHCompute{
		CPU: &models.VCHComputeCPU{
			Limit:       asMHz(d.ResourceLimits.VCHCPULimitsMHz),
			Reservation: asMHz(d.ResourceLimits.VCHCPUReservationsMHz),
			Shares:      asShares(d.ResourceLimits.VCHCPUShares),
		},
		Memory: &models.VCHComputeMemory{
			Limit:       asMiB(d.ResourceLimits.VCHMemoryLimitsMB),
			Reservation: asMiB(d.ResourceLimits.VCHMemoryReservationsMB),
			Shares:      asShares(d.ResourceLimits.VCHMemoryShares),
		},
		Resource: &models.ManagedObject{
			ID: mobidToID(vchConfig.Container.ComputeResources[0].String()),
		},
	}

	// network
	model.Network = &models.VCHNetwork{
		Bridge: &models.VCHNetworkBridge{
			PortGroup: &models.ManagedObject{
				ID: mobidToID(vchConfig.ExecutorConfig.Networks[vchConfig.Network.BridgeNetwork].Network.Common.ID),
			},
			IPRange: asIPRange(vchConfig.Network.BridgeIPRange),
		},
		Client:     asNetwork(vchConfig.ExecutorConfig.Networks["client"]),
		Management: asNetwork(vchConfig.ExecutorConfig.Networks["management"]),
		Public:     asNetwork(vchConfig.ExecutorConfig.Networks["public"]),
	}

	containerNetworks := make([]*models.ContainerNetwork, 0, len(vchConfig.Network.ContainerNetworks))
	for key, value := range vchConfig.Network.ContainerNetworks {
		if key != "bridge" {
			containerNetworks = append(containerNetworks, &models.ContainerNetwork{
				Alias: value.Name,
				PortGroup: &models.ManagedObject{
					ID: mobidToID(value.Common.ID),
				},
				Nameservers: *asIPAddresses(&value.Nameservers),
				Gateway: &models.Gateway{
					Address:             asIPAddress(value.Gateway.IP),
					RoutingDestinations: []models.IPRange{asIPRange(&value.Gateway)},
				},
				IPRanges: *asIPRanges(&value.Destinations),
			})
		}
	}
	model.Network.Container = containerNetworks

	// storage
	scratchSize := int(vchConfig.Storage.ScratchSize)
	model.Storage = &models.VCHStorage{
		BaseImageSize: asKiB(&scratchSize),
	}

	volumeLocations := make([]*models.VCHStorageVolumeStoresItems0, 0, len(vchConfig.Storage.VolumeLocations))
	for label, path := range vchConfig.Storage.VolumeLocations {
		volume := models.VCHStorageVolumeStoresItems0{Datastore: path.String(), Label: label}
		volumeLocations = append(volumeLocations, &volume)
	}
	model.Storage.VolumeStores = volumeLocations

	imageStores := make([]string, 0, len(vchConfig.Storage.ImageStores))
	for _, value := range vchConfig.Storage.ImageStores {
		imageStores = append(imageStores, value.String())
	}
	model.Storage.ImageStores = imageStores

	// auth
	model.Auth = &models.VCHAuth{
		Client: &models.VCHAuthClient{},
	}

	if vchConfig.Certificate.HostCertificate != nil {
		model.Auth.Server = &models.VCHAuthServer{
			Certificate: asPemCertificate(vchConfig.Certificate.HostCertificate.Cert),
		}
	}

	model.Auth.Client.CertificateAuthorities = asPemCertificates(vchConfig.Certificate.CertificateAuthorities)

	// endpoint
	model.Endpoint = &models.VCHEndpoint{
		UseResourcePool: d.UseRP,
		Memory:          asMiB(&d.MemoryMB),
		CPU: &models.VCHEndpointCPU{
			Sockets: int64(d.NumCPUs),
		},
		OperationsCredentials: &models.VCHEndpointOperationsCredentials{
			User: vchConfig.Connection.Username,
			// Password intentionally excluded from GET responses for security reasons!
		},
	}

	// registry
	model.Registry = &models.VCHRegistry{
		Insecure:               vchConfig.Registry.InsecureRegistries,
		Whitelist:              vchConfig.Registry.RegistryWhitelist,
		Blacklist:              vchConfig.Registry.RegistryBlacklist,
		CertificateAuthorities: asPemCertificates(vchConfig.Certificate.RegistryCertificateAuthorities),
		ImageFetchProxy:        asImageFetchProxy(vchConfig.ExecutorConfig.Sessions[config.VicAdminService], config.VICAdminHTTPProxy, config.VICAdminHTTPSProxy),
	}

	// runtime
	model.Runtime = &models.VCHRuntime{}

	installerVer := version.GetBuild()
	upgradeStatus := upgradeStatusMessage(context.Background(), vch, installerVer, vchConfig.Version)
	model.Runtime.UpgradeStatus = upgradeStatus

	powerState, err := vch.PowerState(context.Background())
	if err != nil {
		powerState = "error"
	}
	model.Runtime.PowerState = string(powerState)

	if public := vchConfig.ExecutorConfig.Networks["public"]; public != nil {
		if publicIP := public.Assigned.IP; publicIP != nil {
			var dockerPort string
			if !vchConfig.HostCertificate.IsNil() {
				dockerPort = fmt.Sprintf("%d", opts.DefaultTLSHTTPPort)
			} else {
				dockerPort = fmt.Sprintf("%d", opts.DefaultHTTPPort)
			}

			model.Runtime.DockerHost = fmt.Sprintf("%s:%s", publicIP, dockerPort)
			model.Runtime.AdminPortal = fmt.Sprintf("https://%s:2378", publicIP)
		}
	}

	return model, nil
}

func asBytes(value *int, units string) *models.ValueBytes {
	if value == nil {
		return nil
	}

	if *value == 0 {
		return nil
	}

	return &models.ValueBytes{
		Value: models.Value{
			Value: int64(*value),
			Units: units,
		},
	}
}

func asKiB(value *int) *models.ValueBytes {
	return asBytes(value, models.ValueBytesUnitsKiB)
}

func asMiB(value *int) *models.ValueBytes {
	return asBytes(value, models.ValueBytesUnitsMiB)
}

func asMHz(value *int) *models.ValueHertz {
	if value == nil {
		return nil
	}

	if *value == 0 {
		return nil
	}

	return &models.ValueHertz{
		Value: models.Value{
			Value: int64(*value),
			Units: models.ValueHertzUnitsMHz,
		},
	}
}

func asShares(shares *types.SharesInfo) *models.Shares {
	if shares == nil {
		return nil
	}

	return &models.Shares{
		Level:  string(shares.Level),
		Number: int64(shares.Shares),
	}
}

func asIPAddress(address net.IP) models.IPAddress {
	return models.IPAddress(address.String())
}

func asIPAddresses(addresses *[]net.IP) *[]models.IPAddress {
	m := make([]models.IPAddress, 0, len(*addresses))
	for _, value := range *addresses {
		m = append(m, asIPAddress(value))
	}

	return &m
}

func asIPRange(network *net.IPNet) models.IPRange {
	if network == nil {
		return models.IPRange{}
	}

	return models.IPRange{CIDR: models.CIDR(network.String())}
}

func asIPRanges(networks *[]net.IPNet) *[]models.IPRange {
	m := make([]models.IPRange, 0, len(*networks))
	for _, value := range *networks {
		m = append(m, asIPRange(&value))
	}

	return &m
}

func asNetwork(network *executor.NetworkEndpoint) *models.Network {
	if network == nil {
		return nil
	}

	m := &models.Network{
		PortGroup: &models.ManagedObject{
			ID: mobidToID(network.Network.Common.ID),
		},
		Nameservers: *asIPAddresses(&network.Network.Nameservers),
	}

	if network.Network.Gateway.IP != nil {
		m.Gateway = &models.Gateway{
			Address:             asIPAddress(network.Network.Gateway.IP),
			RoutingDestinations: *asIPRanges(&network.Network.Destinations),
		}
	}

	return m
}

func mobidToID(mobid string) string {
	moref := new(types.ManagedObjectReference)
	ok := moref.FromString(mobid)
	if !ok {
		return "" // TODO: Handle? (We probably don't want to let this fail the request, but may want to convey that something unexpected happened.)
	}

	return moref.Value
}

func asPemCertificates(certificates []byte) []*models.X509Data {
	var buf bytes.Buffer

	m := make([]*models.X509Data, 0)
	for c := &certificates; len(*c) > 0; {
		b, rest := pem.Decode(*c)

		err := pem.Encode(&buf, b)
		if err != nil {
			continue // TODO: Handle? (We probably don't want to let this fail the request, but may want to convey that something unexpected happened.)
		}

		m = append(m, &models.X509Data{
			Pem: models.PEM(buf.String()),
		})

		c = &rest
	}

	return m
}

func asPemCertificate(certificates []byte) *models.X509Data {
	m := asPemCertificates(certificates)

	if len(m) > 1 {
		// TODO: Error? (We probably don't want to let this fail the request, but may want to convey that something unexpected happened.)
	}

	return m[0]
}

func asImageFetchProxy(sessionConfig *executor.SessionConfig, http, https string) *models.VCHRegistryImageFetchProxy {
	var httpProxy, httpsProxy strfmt.URI
	for _, env := range sessionConfig.Cmd.Env {
		if strings.HasPrefix(env, http+"=") {
			httpProxy = strfmt.URI(strings.SplitN(env, "=", 2)[1])
		}
		if strings.HasPrefix(env, https+"=") {
			httpsProxy = strfmt.URI(strings.SplitN(env, "=", 2)[1])
		}
	}

	if httpProxy == "" && httpsProxy == "" {
		return nil
	}

	return &models.VCHRegistryImageFetchProxy{HTTP: httpProxy, HTTPS: httpsProxy}
}
