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

package validate

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"path"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/govc/host/esxcli"
	"github.com/vmware/govmomi/license"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

type Validator struct {
	TargetPath            string
	DatacenterPath        string
	ClusterPath           string
	ResourcePoolPath      string
	ImageStorePath        string
	ExternalNetworkPath   string
	BridgeNetworkPath     string
	BridgeNetworkName     string
	ManagementNetworkPath string
	ManagementNetworkName string

	Session *session.Session
	Context context.Context

	isVC   bool
	issues []error

	DisableFirewallCheck bool
	DisableDRSCheck      bool
}

func NewValidator(ctx context.Context, input *data.Data) (*Validator, error) {
	defer trace.End(trace.Begin(""))
	var err error

	v := &Validator{}
	v.Context = ctx
	tURL := input.URL

	// default to https scheme
	if tURL.Scheme == "" {
		tURL.Scheme = "https"
	}

	// if they specified only an IP address the parser for some reason considers that a path
	if tURL.Host == "" {
		tURL.Host = tURL.Path
		tURL.Path = ""
	}

	sessionconfig := &session.Config{
		Insecure: input.Insecure,
	}

	// only allow the datacenter to be specified in the taget url, if any
	pElems := filepath.SplitList(tURL.Path)
	if len(pElems) > 1 {
		detail := "Target should specify only the datacenter in the path component (e.g. https://addr/datacenter) - resource pools or folders are separate arguments"
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// if a datacenter was specified, set it
	v.DatacenterPath = tURL.Path
	if v.DatacenterPath != "" {
		sessionconfig.DatacenterPath = v.DatacenterPath
		// path needs to be stripped before we can use it as a service url
		tURL.Path = ""
	}

	sessionconfig.Service = tURL.String()

	v.Session = session.NewSession(sessionconfig)
	v.Session, err = v.Session.Connect(v.Context)
	if err != nil {
		return nil, err
	}

	// cached here to allow a modicum of testing while session is still in use.
	v.isVC = v.Session.IsVC()
	finder := find.NewFinder(v.Session.Client.Client, false)
	v.Session.Finder = finder

	v.Session.Populate(ctx)

	if v.Session.Datacenter == nil {
		detail := "Target should specify datacenter when there are multiple possibilities, e.g. https://addr/datacenter"
		log.Error(detail)
		// TODO: list available datacenters
		return nil, errors.New(detail)
	}

	v.DatacenterPath = v.Session.Datacenter.InventoryPath

	return v, nil
}

func (v *Validator) NoteIssue(err error) {
	if err != nil {
		v.issues = append(v.issues, err)
	}
}

func (v *Validator) ListIssues() error {
	defer trace.End(trace.Begin(""))

	if len(v.issues) == 0 {
		return nil
	}

	log.Error("--------------------")
	for _, err := range v.issues {
		log.Error(err)
	}

	return errors.New("validation of configuration failed")
}

// Validate runs through various validations, starting with basics such as naming, moving onto vSphere entities
// and then the compatibility between those entities. It assembles a set of issues that are found for reporting.
func (v *Validator) Validate(ctx context.Context, input *data.Data) (*metadata.VirtualContainerHostConfigSpec, error) {
	defer trace.End(trace.Begin(""))
	log.Infof("Validating supplied configuration")

	conf := &metadata.VirtualContainerHostConfigSpec{}

	v.basics(ctx, input, conf)

	v.target(ctx, input, conf)
	v.compute(ctx, input, conf)
	v.storage(ctx, input, conf)
	v.network(ctx, input, conf)
	v.firewall(ctx)
	v.license(ctx)
	v.drs(ctx)

	v.certificate(ctx, input, conf)

	// Perform the higher level compatibility and consistency checks
	v.compatibility(ctx, conf)

	// TODO: determine if this is where we should turn the noted issues into message
	return conf, v.ListIssues()

}

func (v *Validator) basics(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// TODO: ensure that displayname doesn't violate constraints (length, characters, etc)
	conf.SetName(input.DisplayName)
	conf.SetDebug(input.Debug.Debug)

	conf.Name = input.DisplayName

	conf.BootstrapImagePath = fmt.Sprintf("[%s] %s/%s",
		input.ImageDatastoreName,
		input.DisplayName,
		path.Base(input.BootstrapISO),
	)
}

func (v *Validator) compute(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// Compute

	// compute resource looks like <toplevel>/<sub/path>
	// this should map to /datacenter-name/host/<toplevel>/Resources/<sub/path>
	// we need to validate that <toplevel> exists and then that the combined path exists.

	pool, err := v.ResourcePoolHelper(ctx, input.ComputeResourcePath)
	v.NoteIssue(err)
	if pool == nil {
		return
	}

	// stash the pool for later use
	v.ResourcePoolPath = pool.InventoryPath

	// some hoops for while we're still using session package
	v.Session.Pool = pool
	v.Session.PoolPath = pool.InventoryPath
	v.Session.ClusterPath = v.inventoryPathToCluster(pool.InventoryPath)

	clusters, err := v.Session.Finder.ComputeResourceList(v.Context, v.Session.ClusterPath)
	if err != nil {
		log.Errorf("Unable to acquire reference to cluster %q: %s", path.Base(v.Session.ClusterPath), err)
		v.NoteIssue(err)
		return
	}

	if len(clusters) != 1 {
		err := fmt.Errorf("Unable to acquire unambiguous reference to cluster %q", path.Base(v.Session.ClusterPath))
		log.Error(err)
		v.NoteIssue(err)
	}

	v.Session.Cluster = clusters[0]

	// TODO: for vApp creation assert that the name doesn't exist
	// TODO: for RP creation assert whatever we decide about the pool - most likely that it's empty
}

func (v *Validator) storage(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// Image Store
	imageDSpath, ds, err := v.DatastoreHelper(ctx, input.ImageDatastoreName)
	v.NoteIssue(err)
	if ds == nil {
		log.Errorf("Datastore for images is nil")
		return
	}
	v.SetDatastore(ds, imageDSpath)
	conf.AddImageStore(imageDSpath)

	if conf.VolumeLocations == nil {
		conf.VolumeLocations = make(map[string]*url.URL)
	}

	// TODO: add volume locations
	for label, volDSpath := range input.VolumeLocations {

		dsURL, _, err := v.DatastoreHelper(ctx, volDSpath)
		v.NoteIssue(err)
		if err != nil {
			log.Errorf("Validator could not look up volume datastore due to error: %s", err)
			return
		}
		conf.VolumeLocations[label] = dsURL

	}
}

func (v *Validator) network(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// External net
	extMoref, err := v.networkHelper(ctx, input.ExternalNetworkName)
	v.NoteIssue(err)
	conf.AddNetwork(&metadata.NetworkEndpoint{
		Common: metadata.Common{
			Name: "external",
		},
		Network: metadata.ContainerNetwork{
			Common: metadata.Common{
				Name: "external",
				ID:   extMoref,
			},
			Default: true, // external network is default for appliance
		},
	})
	// Bridge network should be different with all other networks
	if input.BridgeNetworkName == input.ExternalNetworkName {
		v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - external also uses %q", input.BridgeNetworkName))
	}

	// Client net
	if input.ClientNetworkName == "" {
		input.ClientNetworkName = input.ExternalNetworkName
	}
	clientMoref, err := v.networkHelper(ctx, input.ClientNetworkName)
	v.NoteIssue(err)
	conf.AddNetwork(&metadata.NetworkEndpoint{
		Common: metadata.Common{
			Name: "client",
		},
		Network: metadata.ContainerNetwork{
			Common: metadata.Common{
				Name: "client",
				ID:   clientMoref,
			},
		},
	})
	if input.BridgeNetworkName == input.ClientNetworkName {
		v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - client also uses %q", input.BridgeNetworkName))
	}

	// Management net
	if input.ManagementNetworkName == "" {
		input.ManagementNetworkName = input.ClientNetworkName
	}
	managementMoref, err := v.networkHelper(ctx, input.ManagementNetworkName)
	v.NoteIssue(err)
	conf.AddNetwork(&metadata.NetworkEndpoint{
		Network: metadata.ContainerNetwork{
			Common: metadata.Common{
				Name: "management",
				ID:   managementMoref,
			},
		},
	})
	if input.BridgeNetworkName == input.ManagementNetworkName {
		v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - management also uses %q", input.BridgeNetworkName))
	}

	// Bridge net -
	//   vCenter: must exist and must be a DPG
	//   ESX: doesn't need to exist
	//
	// for now we're hardcoded to "bridge" for the container host name
	conf.BridgeNetwork = "bridge"
	endpointMoref, err := v.dpgHelper(ctx, input.BridgeNetworkName)
	netMoref := endpointMoref
	if err != nil {
		if _, ok := err.(*find.NotFoundError); !ok || v.IsVC() {
			v.NoteIssue(fmt.Errorf("An existing distributed port group must be specified for bridge network on vCenter: %s", err))
		}

		// this allows the dispatcher to create the network with corresponding name
		// if BridgeNetworkName doesn't already exist then we set the ContainerNetwork
		// ID to the name, but leaving the NetworkEndpoint moref as ""
		netMoref = input.BridgeNetworkName
	}

	bridgeNet := &metadata.NetworkEndpoint{
		Common: metadata.Common{
			Name: "bridge",
			ID:   endpointMoref,
		},
		Static: &net.IPNet{IP: net.IPv4zero}, // static but managed externally
		Network: metadata.ContainerNetwork{
			Common: metadata.Common{
				Name: "bridge",
				ID:   netMoref,
			},
		},
	}
	// we  need to have the bridge network identified as an available container
	// network
	conf.AddContainerNetwork(&bridgeNet.Network)
	// we also need to have the appliance attached to the bridge network to allow
	// port forwarding
	conf.AddNetwork(bridgeNet)

	// add mapped networks
	//   these should be a distributed port groups in vCenter
	for name, net := range input.MappedNetworks {
		// "bridge" is reserved
		if name == "bridge" {
			v.NoteIssue(fmt.Errorf("Cannot use reserved name \"bridge\" for container network"))
			continue
		}

		gw := input.MappedNetworksGateways[name]
		pools := input.MappedNetworksIPRanges[name]
		dns := input.MappedNetworksDNS[name]
		if len(pools) != 0 && nilIPNet(gw) {
			v.NoteIssue(fmt.Errorf("IP range specified without gateway for container network %q", name))
			continue
		}

		err = nil
		// verify ip ranges are within subnet,
		// and don't overlap with each other
		for i, r := range pools {
			if !gw.Contains(r.FirstIP) || !gw.Contains(r.LastIP) {
				err = fmt.Errorf("IP range %q is not in subnet %q", r, gw)
				break
			}

			for _, r2 := range pools[i+1:] {
				if r2.Overlaps(r) {
					err = fmt.Errorf("Overlapping ip ranges: %q %q", r2, r)
					break
				}
			}

			if err != nil {
				break
			}
		}

		if err != nil {
			v.NoteIssue(err)
			continue
		}

		moref, err := v.dpgHelper(ctx, net)
		v.NoteIssue(err)
		mappedNet := &metadata.ContainerNetwork{
			Common: metadata.Common{
				Name: name,
				ID:   moref,
			},
			Gateway:     gw,
			Nameservers: dns,
			Pools:       pools,
		}
		if input.BridgeNetworkName == net {
			v.NoteIssue(errors.Errorf("the bridge network must not be shared with another network role - %q also mapped as container network %q", input.BridgeNetworkName, name))
		}

		conf.AddContainerNetwork(mappedNet)
	}
}

// generateBridgeName returns a name that can be used to create a switch/pg pair on ESX
func (v *Validator) generateBridgeName(ctx, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) string {
	defer trace.End(trace.Begin(""))

	return input.DisplayName
}

func (v *Validator) checkSessionSet() []string {
	var errs []string

	if v.Session.Datastore == nil {
		errs = append(errs, "datastore not set")
	}
	if v.Session.Cluster == nil {
		errs = append(errs, "cluster not set")
	}

	return errs
}

func (v *Validator) sessionValid(errMsg string) bool {
	if c := v.checkSessionSet(); len(c) > 0 {
		log.Error(errMsg)
		for _, e := range c {
			log.Errorf("  %s", e)
		}
		v.NoteIssue(errors.New(errMsg))
		return false
	}
	return true
}

func (v *Validator) firewall(ctx context.Context) {
	if v.DisableFirewallCheck {
		return
	}
	defer trace.End(trace.Begin(""))

	var hosts []*object.HostSystem
	var err error

	rule := types.HostFirewallRule{
		Port:      8080, // serialOverLANPort
		PortType:  types.HostFirewallRulePortTypeDst,
		Protocol:  string(types.HostFirewallRuleProtocolTcp),
		Direction: types.HostFirewallRuleDirectionOutbound,
	}

	errMsg := "Firewall check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	if hosts, err = v.Session.Datastore.AttachedClusterHosts(ctx, v.Session.Cluster); err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		v.NoteIssue(err)
		return
	}

	var misconfiguredEnabled []string
	var misconfiguredDisabled []string
	var correct []string

	for _, host := range hosts {
		fs, err := host.ConfigManager().FirewallSystem(ctx)
		if err != nil {
			v.NoteIssue(err)
			break
		}

		disabled := false
		esxfw, err := esxcli.GetFirewallInfo(host)
		if err != nil {
			v.NoteIssue(err)
			break
		}
		if !esxfw.Enabled {
			disabled = true
			log.Infof("Firewall status: DISABLED on %q", host.InventoryPath)
		} else {
			log.Infof("Firewall status: ENABLED on %q", host.InventoryPath)
		}

		info, err := fs.Info(ctx)
		if err != nil {
			v.NoteIssue(err)
			break
		}

		rs := object.HostFirewallRulesetList(info.Ruleset)
		_, err = rs.EnabledByRule(rule, true)
		if err != nil {
			if !disabled {
				misconfiguredEnabled = append(misconfiguredEnabled, host.InventoryPath)
			} else {
				misconfiguredDisabled = append(misconfiguredDisabled, host.InventoryPath)
			}
		} else {
			correct = append(correct, host.InventoryPath)
		}
	}

	if len(correct) > 0 {
		log.Info("Firewall configuration OK on hosts:")
		for _, h := range correct {
			log.Infof("  %q", h)
		}
	}
	if len(misconfiguredEnabled) > 0 {
		log.Error("Firewall configuration incorrect on hosts:")
		for _, h := range misconfiguredEnabled {
			log.Errorf("  %q", h)
		}
		// TODO: when we can intelligently place containerVMs on hosts with proper config, install
		// can proceed if there is at least one host properly configured. For now this prevents install.
		msg := "Firewall must permit 8080/tcp outbound to use VIC"
		log.Error(msg)
		v.NoteIssue(errors.New(msg))
		return
	}
	if len(misconfiguredDisabled) > 0 {
		log.Warning("Firewall configuration will be incorrect if firewall is reenabled on hosts:")
		for _, h := range misconfiguredDisabled {
			log.Warningf("  %q", h)
		}
		log.Warning("Firewall must permit 8080/tcp outbound if firewall is reenabled")
	}
}

func (v *Validator) license(ctx context.Context) {
	var err error

	errMsg := "License check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	if v.IsVC() {
		if err = v.checkAssignedLicenses(ctx); err != nil {
			v.NoteIssue(err)
			return
		}
	} else {
		if err = v.checkLicense(ctx); err != nil {
			v.NoteIssue(err)
			return
		}
	}
}

func (v *Validator) assignedLicenseHasFeature(la []types.LicenseAssignmentManagerLicenseAssignment, feature string) bool {
	for _, a := range la {
		if license.HasFeature(a.AssignedLicense, feature) {
			return true
		}
	}
	return false
}

func (v *Validator) checkAssignedLicenses(ctx context.Context) error {
	var hosts []*object.HostSystem
	var invalidLic []string
	var validLic []string
	var err error
	client := v.Session.Client.Client

	if hosts, err = v.Session.Datastore.AttachedClusterHosts(ctx, v.Session.Cluster); err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		return err
	}

	lm := license.NewManager(client)

	am, err := lm.AssignmentManager(ctx)
	if err != nil {
		return err
	}

	features := []string{"serialuri", "dvs"}

	for _, host := range hosts {
		valid := true
		la, err := am.QueryAssigned(ctx, host.Reference().Value)
		if err != nil {
			return err
		}

		for _, feature := range features {
			if !v.assignedLicenseHasFeature(la, feature) {
				valid = false
				msg := fmt.Sprintf("%q - license missing feature %q", host.InventoryPath, feature)
				invalidLic = append(invalidLic, msg)
			}
		}

		if valid == true {
			validLic = append(validLic, host.InventoryPath)
		}
	}

	if len(validLic) > 0 {
		log.Infof("License check OK on hosts:")
		for _, h := range validLic {
			log.Infof("  %q", h)
		}
	}
	if len(invalidLic) > 0 {
		log.Errorf("License check FAILED on hosts:")
		for _, h := range invalidLic {
			log.Errorf("  %q", h)
		}
		msg := "License does not meet minimum requirements to use VIC"
		return errors.New(msg)
	}
	return nil
}

func (v *Validator) checkLicense(ctx context.Context) error {
	var invalidLic []string
	client := v.Session.Client.Client

	lm := license.NewManager(client)
	licenses, err := lm.List(ctx)
	if err != nil {
		return err
	}
	v.checkEvalLicense(licenses)

	features := []string{"serialuri"}

	for _, feature := range features {
		if len(licenses.WithFeature(feature)) == 0 {
			msg := fmt.Sprintf("Host license missing feature %q", feature)
			invalidLic = append(invalidLic, msg)
		}
	}

	if len(invalidLic) > 0 {
		log.Errorf("License check FAILED:")
		for _, h := range invalidLic {
			log.Errorf("  %q", h)
		}
		msg := "License does not meet minimum requirements to use VIC"
		return errors.New(msg)
	}
	log.Infof("License check OK")
	return nil
}

func (v *Validator) checkEvalLicense(licenses []types.LicenseManagerLicenseInfo) {
	for _, l := range licenses {
		if l.EditionKey == "eval" {
			log.Warning("Evaluation license detected. VIC may not function if evaluation expires or insufficient license is later assigned.")
		}
	}
}

// isStandaloneHost checks if host is ESX or vCenter with single host
func (v *Validator) isStandaloneHost() bool {
	cl := v.Session.Cluster.Reference()

	if cl.Type != "ClusterComputeResource" {
		return true
	}
	return false
}

// drs checks that DRS is enabled
func (v *Validator) drs(ctx context.Context) {
	if v.DisableDRSCheck {
		return
	}
	defer trace.End(trace.Begin(""))

	errMsg := "DRS check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	cl := v.Session.Cluster
	ref := cl.Reference()

	if v.isStandaloneHost() {
		log.Info("DRS check SKIPPED - target is standalone host")
		return
	}

	var ccr mo.ClusterComputeResource

	err := cl.Properties(ctx, ref, []string{"configurationEx"}, &ccr)
	if err != nil {
		msg := fmt.Sprintf("Failed to validate DRS config: %s", err)
		v.NoteIssue(errors.New(msg))
		return
	}

	z := ccr.ConfigurationEx.(*types.ClusterConfigInfoEx).DrsConfig

	if !(*z.Enabled) {
		log.Error("DRS check FAILED")
		log.Errorf("  DRS must be enabled on cluster %q", v.Session.Pool.InventoryPath)
		v.NoteIssue(errors.New("DRS must be enabled to use VIC"))
		return
	}
	log.Info("DRS check OK on:")
	log.Infof("  %q", v.Session.Pool.InventoryPath)
}

func (v *Validator) target(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	targetURL := input.Target.URLWithoutPassword()
	if !v.IsVC() {
		var err error
		targetURL, err = url.Parse(v.Session.Service)
		if err != nil {
			v.NoteIssue(fmt.Errorf("Error processing target after transformation to SOAP endpoint: %q: %s", v.Session.Service, err))
			return
		}
	}

	conf.Target = *targetURL
	conf.Insecure = input.Insecure

	// TODO: more checks needed here if specifying service account for VCH
}

func (v *Validator) certificate(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	if len(input.CertPEM) == 0 && len(input.KeyPEM) == 0 {
		// if there's no data supplied then we're configuring without TLS
		log.Debug("Configuring without TLS due to empty key and cert buffers")
		return
	}

	// check the cert can be loaded
	_, err := tls.X509KeyPair(input.CertPEM, input.KeyPEM)
	v.NoteIssue(err)

	conf.HostCertificate = &metadata.RawCertificate{
		Key:  input.KeyPEM,
		Cert: input.CertPEM,
	}
}

func (v *Validator) compatibility(ctx context.Context, conf *metadata.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// TODO: add checks such as datastore is acessible from target cluster
	errMsg := "Compatibility check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	_, err := v.Session.Datastore.AttachedClusterHosts(v.Context, v.Session.Cluster)
	v.NoteIssue(err)
}

func (v *Validator) networkHelper(ctx context.Context, path string) (string, error) {
	defer trace.End(trace.Begin(path))

	nets, err := v.Session.Finder.NetworkList(ctx, path)
	if err != nil {
		log.Debugf("no such network %q", path)
		// TODO: error message about no such match and how to get a network list
		// we return err directly here so we can check the type
		return "", err
	}
	if len(nets) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return "", errors.New("ambiguous network " + path)
	}

	moref := nets[0].Reference()
	return moref.String(), nil
}

func (v *Validator) dpgMorefHelper(ctx context.Context, ref string) (string, error) {
	defer trace.End(trace.Begin(ref))

	moref := new(types.ManagedObjectReference)
	ok := moref.FromString(ref)
	if !ok {
		// TODO: error message about no such match and how to get a network list
		return "", errors.New("could not restore serialized managed object reference: " + ref)
	}

	net, err := v.Session.Finder.ObjectReference(ctx, *moref)
	if err != nil {
		// TODO: error message about no such match and how to get a network list
		return "", errors.New("unable to locate network from moref: " + ref)
	}

	// ensure that the type of the network is a Distributed Port Group if the target is a vCenter
	// if it's not then any network suffices
	if v.IsVC() {
		_, dpg := net.(*object.DistributedVirtualPortgroup)
		if !dpg {
			return "", fmt.Errorf("%q is not a Distributed Port Group", ref)
		}
	}

	return ref, nil
}

func (v *Validator) dpgHelper(ctx context.Context, path string) (string, error) {
	defer trace.End(trace.Begin(path))

	nets, err := v.Session.Finder.NetworkList(ctx, path)
	if err != nil {
		log.Debugf("no such network %q", path)
		// TODO: error message about no such match and how to get a network list
		// we return err directly here so we can check the type
		return "", err
	}
	if len(nets) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return "", errors.New("ambiguous network " + path)
	}

	// ensure that the type of the network is a Distributed Port Group if the target is a vCenter
	// if it's not then any network suffices
	if v.IsVC() {
		_, dpg := nets[0].(*object.DistributedVirtualPortgroup)
		if !dpg {
			return "", fmt.Errorf("%q is not a Distributed Port Group", path)
		}
	}

	moref := nets[0].Reference()
	return moref.String(), nil
}

func (v *Validator) DatastoreHelper(ctx context.Context, path string) (*url.URL, *object.Datastore, error) {
	defer trace.End(trace.Begin(path))

	dsURL, err := url.Parse(path)
	if err != nil {
		return nil, nil, errors.Errorf("parsing error occured while parsing datastore path: %s", err)
	}

	// url scheme does not contain ://, so remove it to make url work
	if dsURL.Scheme != "" && dsURL.Scheme != "ds" {
		return nil, nil, errors.Errorf("bad scheme %q provided for datastore", dsURL.Scheme)
	}

	dsURL.Scheme = "ds"

	// if a datastore name (e.g. "datastore1") is specifed with no decoration then this
	// is interpreted as the Path
	if dsURL.Host == "" && dsURL.Path != "" {
		pathElements := strings.SplitN(path, "/", 2)
		dsURL.Host = pathElements[0]
		if len(pathElements) > 1 {
			dsURL.Path = pathElements[1]
		} else {
			dsURL.Path = ""
		}
	}
	if dsURL.Host == "" {
		return nil, nil, errors.New("datastore hostname came back empty")
	}

	stores, err := v.Session.Finder.DatastoreList(ctx, dsURL.Host)
	if err != nil {
		log.Debugf("no such datastore %#v", dsURL)
		// TODO: error message about no such match and how to get a datastore list
		// we return err directly here so we can check the type
		return nil, nil, err
	}
	if len(stores) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return nil, nil, errors.New("ambiguous datastore " + dsURL.Host)
	}

	// temporary until session is extracted
	// FIXME: commented out until components can consume moid
	// dsURL.Host = stores[0].Reference().Value

	return dsURL, stores[0], nil
}

func (v *Validator) SetDatastore(ds *object.Datastore, path *url.URL) {
	v.Session.Datastore = ds
	v.Session.DatastorePath = path.Host
}

func (v *Validator) ResourcePoolHelper(ctx context.Context, path string) (*object.ResourcePool, error) {
	defer trace.End(trace.Begin(path))

	// if compute-resource is unspecified is there a default
	if path == "" || path == "/" {
		if v.Session.Pool != nil {
			log.Debugf("Using default resource pool for compute resource: %q", v.Session.Pool.InventoryPath)
			return v.Session.Pool, nil
		}

		// if no path specified and no default available the show all
		v.suggestComputeResource("*")
		return nil, errors.New("No unambiguous default compute resource available: --compute-resource must be specified")
	}

	ipath := v.computePathToInventoryPath(path)
	log.Debugf("Converted original path %q to %q", path, ipath)

	// first try the path directly without any processing
	pools, err := v.Session.Finder.ResourcePoolList(ctx, path)
	if err != nil {
		log.Debugf("Failed to look up compute resource as absolute path %q: %s", path, err)
		if _, ok := err.(*find.NotFoundError); !ok {
			// we return err directly here so we can check the type
			return nil, err
		}

		// if it starts with datacenter then we know it's absolute and invalid
		if strings.HasPrefix(path, "/"+v.Session.DatacenterPath) {
			v.suggestComputeResource(path)
			return nil, err
		}
	}

	if len(pools) == 0 {
		// assume it's a cluster specifier - that's the formal case, e.g. /cluster/resource/pool
		// not /cluster/Resources/resource/pool
		// everything from now on will use this assumption

		pools, err = v.Session.Finder.ResourcePoolList(ctx, ipath)
		if err != nil {
			log.Debugf("failed to look up compute resource as cluster path %q: %s", ipath, err)
			if _, ok := err.(*find.NotFoundError); !ok {
				// we return err directly here so we can check the type
				return nil, err
			}
		}
	}

	if len(pools) == 1 {
		log.Debugf("Selected compute resource %q", pools[0].InventoryPath)
		return pools[0], nil
	}

	// both cases we want to suggest options
	v.suggestComputeResource(ipath)

	if len(pools) == 0 {
		log.Debugf("no such compute resource %q", path)
		// we return err directly here so we can check the type
		return nil, err
	}

	// TODO: error about required disabmiguation and list entries in nets
	return nil, errors.New("ambiguous compute resource " + path)
}

func (v *Validator) suggestComputeResource(path string) {
	defer trace.End(trace.Begin(path))

	log.Infof("Suggesting valid values for --compute-resource based on %q", path)

	// allow us to work on inventory paths
	path = v.computePathToInventoryPath(path)

	var matches []string
	for matches = nil; matches == nil; matches = v.findValidPool(path) {
		// back up the path until we find a pool
		newpath := filepath.Dir(path)
		if newpath == "." {
			// filepath.Dir returns . which has no meaning for us
			newpath = "/"
		}
		if newpath == path {
			break
		}
		path = newpath
	}

	if matches == nil {
		// Backing all the way up didn't help
		log.Info("Failed to find resource pool in the provided path, showing all top level resource pools.")
		matches = v.findValidPool("*")
	}

	if matches != nil {
		// we've collected recommendations - displayname
		log.Info("Suggested values for --compute-resource:")
		for _, p := range matches {
			log.Infof("  %q", v.inventoryPathToComputePath(p))
		}
		return
	}

	log.Info("No resource pools found")
}

func (v *Validator) findValidPool(path string) []string {
	defer trace.End(trace.Begin(path))

	// list pools in path
	matches := v.listResourcePools(path)
	if matches != nil {
		sort.Strings(matches)
		return matches
	}

	// no pools in path, but if path is cluster, list pools in cluster
	clusters, err := v.Session.Finder.ComputeResourceList(v.Context, path)
	if len(clusters) == 0 {
		// not a cluster
		log.Debugf("Path %q does not identify a cluster (or clusters) or the list could not be obtained: %s", path, err)
		return nil
	}

	if len(clusters) > 1 {
		log.Debugf("Suggesting clusters as there are multiple matches")
		matches = make([]string, len(clusters))
		for i, c := range clusters {
			matches[i] = c.InventoryPath
		}
		sort.Strings(matches)
		return matches
	}

	log.Debugf("Suggesting pools for cluster %q", clusters[0].Name())
	matches = v.listResourcePools(fmt.Sprintf("%s/Resources/*", clusters[0].InventoryPath))
	if matches == nil {
		// no child pools so recommend cluster directly
		return []string{clusters[0].InventoryPath}
	}

	return matches
}

func (v *Validator) listResourcePools(path string) []string {
	defer trace.End(trace.Begin(path))

	pools, err := v.Session.Finder.ResourcePoolList(v.Context, path+"/*")
	if err != nil {
		log.Debugf("Unable to list pools for %q: %s", path, err)
		return nil
	}

	if len(pools) == 0 {
		return nil
	}

	matches := make([]string, len(pools))
	for i, p := range pools {
		matches[i] = p.InventoryPath
	}

	return matches
}

func (v *Validator) computePathToInventoryPath(path string) string {
	defer trace.End(trace.Begin(path))

	// if it opens with the datacenter prefix the assume it's an absolute
	if strings.HasPrefix(path, v.DatacenterPath) {
		log.Debugf("Path is treated as absolute given datacenter prefix %q", v.DatacenterPath)
		return path
	}

	parts := []string{
		v.DatacenterPath, // has leading /
		"host",
		"*", // easy for ESX
		"Resources",
	}

	// normalize the path - strip leading /
	path = strings.TrimPrefix(path, "/")

	// if it's vCenter the first element is the cluster or host, then resource pool path
	if v.IsVC() {
		pElem := strings.SplitN(path, "/", 2)
		if pElem[0] != "" {
			parts[2] = pElem[0]
		}
		if len(pElem) > 1 {
			parts = append(parts, pElem[1])
		}
	} else if path != "" {
		// for ESX, first element is a pool
		parts = append(parts, path)
	}

	return strings.Join(parts, "/")
}

func (v *Validator) inventoryPathToComputePath(path string) string {
	defer trace.End(trace.Begin(path))

	// sanity check datacenter
	if !strings.HasPrefix(path, v.DatacenterPath) {
		log.Debugf("Expected path to be within target datacenter %q: %q", v.DatacenterPath, path)
		v.NoteIssue(errors.New("inventory path was not in datacenter scope"))
		return ""
	}

	// inventory path is always /dc/host/computeResource/Resources/path/to/pool
	// NOTE: all of the indexes are +1 because the leading / means we have an empty string for [0]
	pElems := strings.Split(path, "/")
	if len(pElems) < 4 {
		log.Debugf("Expected path to be fully qualified, e.g. /dcName/host/clusterName/Resources/poolName: %s", path)
		v.NoteIssue(errors.New("inventory path format was not recognised"))
		return ""
	}

	if len(pElems) == 4 || len(pElems) == 5 {
		// cluster only or cluster/Resources
		return pElems[3]
	}

	// messy but avoid reallocation - overwrite Resources with cluster name
	pElems[4] = pElems[3]

	// /dc/host/cluster/Resources/path/to/pool
	return strings.Join(pElems[4:], "/")
}

// inventoryPathToCluster is a convenience method that will return the cluster
// path prefix or "" in the case of unexpected path structure
func (v *Validator) inventoryPathToCluster(path string) string {
	defer trace.End(trace.Begin(path))

	// inventory path is always /dc/host/computeResource/Resources/path/to/pool
	pElems := strings.Split(path, "/")
	if len(pElems) < 3 {
		log.Debugf("Expected path to be fully qualified, e.g. /dcName/host/clusterName/Resources/poolName: %s", path)
		v.NoteIssue(errors.New("inventory path format was not recognised"))
		return ""
	}

	// /dc/host/cluster/Resources/path/to/pool
	return strings.Join(pElems[:4], "/")
}

func (v *Validator) IsVC() bool {
	return v.isVC
}

func (v *Validator) GetResourcePool(input *data.Data) (*object.ResourcePool, error) {
	defer trace.End(trace.Begin(""))

	return v.ResourcePoolHelper(v.Context, input.ComputeResourcePath)
}

func (v *Validator) AddDeprecatedFields(ctx context.Context, conf *metadata.VirtualContainerHostConfigSpec, input *data.Data) *data.InstallerData {
	defer trace.End(trace.Begin(""))

	dconfig := data.InstallerData{}

	dconfig.ApplianceSize.CPU.Limit = int64(input.NumCPUs)
	dconfig.ApplianceSize.Memory.Limit = int64(input.MemoryMB)

	dconfig.Datacenter = v.Session.Datacenter.Reference()
	dconfig.DatacenterName = v.Session.Datacenter.Name()

	dconfig.Cluster = v.Session.Cluster.Reference()
	dconfig.ClusterPath = v.Session.Cluster.InventoryPath

	dconfig.ResourcePoolPath = v.ResourcePoolPath
	dconfig.UseRP = input.UseRP

	log.Debugf("Datacenter: %q, Cluster: %q, Resource Pool: %q", dconfig.DatacenterName, dconfig.ClusterPath, dconfig.ResourcePoolPath)

	return &dconfig
}

func nilIPNet(n net.IPNet) bool {
	return n.IP == nil && n.Mask == nil
}
