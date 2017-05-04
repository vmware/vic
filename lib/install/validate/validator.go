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
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	units "github.com/docker/go-units"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	vmomisession "github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/session"
)

type Validator struct {
	TargetPath        string
	DatacenterPath    string
	ClusterPath       string
	ResourcePoolPath  string
	ImageStorePath    string
	PublicNetworkPath string
	BridgeNetworkPath string
	BridgeNetworkName string

	Session *session.Session
	Context context.Context

	isVC   bool
	issues []error

	DisableFirewallCheck bool
	DisableDRSCheck      bool
	allowEmptyDC         bool
}

func CreateFromVCHConfig(ctx context.Context, vch *config.VirtualContainerHostConfigSpec, sess *session.Session) (*Validator, error) {
	defer trace.End(trace.Begin(""))

	v := &Validator{}
	v.Session = sess
	v.Context = ctx
	v.isVC = v.Session.IsVC()

	return v, nil
}

func NewValidator(ctx context.Context, input *data.Data) (*Validator, error) {
	defer trace.End(trace.Begin(""))
	var err error

	v := &Validator{}
	v.Context = ctx
	tURL := input.URL

	// normalize the path - strip trailing /
	input.URL.Path = strings.TrimSuffix(input.URL.Path, "/")

	// default to https scheme
	if tURL.Scheme == "" {
		tURL.Scheme = "https"
	}

	// if they specified only an IP address the parser for some reason considers that a path
	if tURL.Host == "" {
		tURL.Host = tURL.Path
		tURL.Path = ""
	}

	if tURL.Scheme == "https" && input.Thumbprint == "" {
		var cert object.HostCertificateInfo
		if err = cert.FromURL(tURL, new(tls.Config)); err != nil {
			return nil, err
		}

		if cert.Err != nil {
			if !input.Force {
				// TODO: prompt user / check ./known_hosts
				log.Errorf("Failed to verify certificate for target=%s (thumbprint=%s)",
					tURL.Host, cert.ThumbprintSHA1)
				return nil, cert.Err
			}
		}

		input.Thumbprint = cert.ThumbprintSHA1
		log.Debugf("Accepting host %q thumbprint %s", tURL.Host, input.Thumbprint)
	}

	sessionconfig := &session.Config{
		Thumbprint: input.Thumbprint,
		Insecure:   input.Force,
	}

	// if a datacenter was specified, set it
	v.DatacenterPath = tURL.Path
	if v.DatacenterPath != "" {
		v.DatacenterPath = strings.TrimPrefix(v.DatacenterPath, "/")
		sessionconfig.DatacenterPath = v.DatacenterPath
		// path needs to be stripped before we can use it as a service url
		tURL.Path = ""
	}

	sessionconfig.Service = tURL.String()

	v.Session = session.NewSession(sessionconfig)
	v.Session.UserAgent = version.UserAgent("vic-machine")
	v.Session, err = v.Session.Connect(v.Context)
	if err != nil {
		return nil, err
	}

	// cached here to allow a modicum of testing while session is still in use.
	v.isVC = v.Session.IsVC()
	finder := find.NewFinder(v.Session.Client.Client, false)
	v.Session.Finder = finder

	// Intentionally ignore any error returned by Populate
	_, err = v.Session.Populate(ctx)
	if err != nil {
		log.Debugf("new validator Session.Populate: %s", err)
	}

	if strings.Contains(sessionconfig.DatacenterPath, "/") {
		detail := "--target should only specify datacenter in the path (e.g. https://addr/datacenter) - specify cluster, resource pool, or folder with --compute-resource"
		log.Error(detail)
		v.suggestDatacenter()
		return nil, errors.New(detail)
	}

	return v, nil
}

var schemeMatch = regexp.MustCompile(`^\w+://`)

// Starting from Go 1.8 the URL parser does not
// work properly with URLs with no Scheme,
// this function adds "https" as Scheme if necessary
func ParseURL(s string) (*url.URL, error) {
	var err error
	var u *url.URL

	if s != "" {
		// Default the scheme to https
		if !schemeMatch.MatchString(s) {
			s = "https://" + s
		}

		u, err = url.Parse(s)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (v *Validator) AllowEmptyDC() {
	v.allowEmptyDC = true
}

func (v *Validator) datacenter() error {
	if v.allowEmptyDC && v.DatacenterPath == "" {
		return nil
	}
	if v.Session.Datacenter != nil {
		v.DatacenterPath = v.Session.Datacenter.InventoryPath
		return nil
	}
	var detail string
	if v.DatacenterPath != "" {
		detail = fmt.Sprintf("Datacenter %q in --target is not found", path.Base(v.DatacenterPath))
	} else {
		// this means multiple datacenter exists, but user did not specify it in --target
		detail = "Datacenter must be specified in --target (e.g. https://addr/datacenter)"
	}
	log.Error(detail)
	v.suggestDatacenter()
	return errors.New(detail)
}

// suggestDatacenter suggests all datacenters on the target
func (v *Validator) suggestDatacenter() {
	defer trace.End(trace.Begin(""))

	log.Info("Suggesting valid values for datacenter in --target")

	dcs, err := v.Session.Finder.DatacenterList(v.Context, "*")
	if err != nil {
		log.Errorf("Unable to list datacenters: %s", err)
		return
	}

	if len(dcs) == 0 {
		log.Info("No datacenters found")
		return
	}

	matches := make([]string, len(dcs))
	for i, d := range dcs {
		matches[i] = d.Name()
	}

	if matches != nil {
		log.Info("Suggested datacenters:")
		for _, d := range matches {
			log.Infof("  %q", d)
		}
		return
	}
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

func (v *Validator) GetIssues() []error {
	return v.issues
}

func (v *Validator) ClearIssues() {
	v.issues = []error{}
}

// Validate runs through various validations, starting with basics such as naming, moving onto vSphere entities
// and then the compatibility between those entities. It assembles a set of issues that are found for reporting.
func (v *Validator) Validate(ctx context.Context, input *data.Data) (*config.VirtualContainerHostConfigSpec, error) {
	defer trace.End(trace.Begin(""))
	log.Infof("Validating supplied configuration")

	conf := &config.VirtualContainerHostConfigSpec{}

	if err := v.datacenter(); err != nil {
		return conf, err
	}

	v.basics(ctx, input, conf)

	v.target(ctx, input, conf)
	v.credentials(ctx, input, conf)
	v.compute(ctx, input, conf)
	v.storage(ctx, input, conf)
	v.network(ctx, input, conf)
	v.CheckFirewall(ctx, conf)
	v.CheckLicense(ctx)
	v.CheckDrs(ctx)

	v.certificate(ctx, input, conf)
	v.certificateAuthorities(ctx, input, conf)
	v.registries(ctx, input, conf)

	// Perform the higher level compatibility and consistency checks
	v.compatibility(ctx, conf)

	// TODO: determine if this is where we should turn the noted issues into message
	return conf, v.ListIssues()

}

func (v *Validator) ValidateTarget(ctx context.Context, input *data.Data) (*config.VirtualContainerHostConfigSpec, error) {
	defer trace.End(trace.Begin(""))
	conf := &config.VirtualContainerHostConfigSpec{}

	log.Infof("Validating target")
	if err := v.datacenter(); err != nil {
		return conf, err
	}
	v.target(ctx, input, conf)
	return conf, v.ListIssues()
}

func (v *Validator) basics(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// TODO: ensure that displayname doesn't violate constraints (length, characters, etc)
	conf.SetName(input.DisplayName)
	conf.SetDebug(input.Debug.Debug)
	conf.Name = input.DisplayName
	conf.Version = version.GetBuild()

	scratchSize, err := units.FromHumanSize(input.ScratchSize)
	if err != nil { // TODO set minimum size of scratch disk
		v.NoteIssue(errors.Errorf("Invalid default image size %s provided; error from parser: %s", input.ScratchSize, err.Error()))
	} else {
		conf.ScratchSize = scratchSize / units.KB
		log.Debugf("Setting scratch image size to %d KB in VCHConfig", conf.ScratchSize)
	}

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

func (v *Validator) target(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// check if host is managed by VC
	v.managedbyVC(ctx)
}

func (v *Validator) managedbyVC(ctx context.Context) {
	defer trace.End(trace.Begin(""))

	if v.IsVC() {
		return
	}
	host, err := v.Session.Finder.DefaultHostSystem(ctx)
	if err != nil {
		v.NoteIssue(fmt.Errorf("Failed to get host system: %s", err))
		return
	}

	var mh mo.HostSystem

	if err = host.Properties(ctx, host.Reference(), []string{"summary.managementServerIp"}, &mh); err != nil {
		v.NoteIssue(fmt.Errorf("Failed to get host properties: %s", err))
		return
	}

	if ip := mh.Summary.ManagementServerIp; ip != "" {
		v.NoteIssue(fmt.Errorf("Target is managed by vCenter server %q, please change --target to vCenter server address or select a standalone ESXi", ip))
	}
	return
}

func (v *Validator) credentials(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	// empty string for password is horrific, but a legitimate scenario especially in isolated labs
	if input.OpsPassword == nil {
		v.NoteIssue(errors.New("Password for operations user has not been set"))
		return
	}

	// check target with ops credentials
	u := input.Target.URL

	conf.Username = input.OpsUser
	conf.Token = *input.OpsPassword
	conf.TargetThumbprint = input.Thumbprint

	// Discard anything other than these URL fields for the target
	stripped := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}
	conf.Target = stripped.String()

	// validate that the provided operations credentials are valid
	stripped.Path = "/sdk"

	var soapClient *soap.Client

	if input.Thumbprint != "" {
		// if any thumprint is specified, then object if there's a mismatch
		soapClient = soap.NewClient(stripped, false)
		soapClient.SetThumbprint(stripped.Host, conf.TargetThumbprint)
	} else {
		soapClient = soap.NewClient(stripped, input.Force)
	}
	soapClient.UserAgent = "vice-validator"

	vimClient, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		v.NoteIssue(fmt.Errorf("Failed to create client for validation of operations credentials: %s", err))
		return
	}

	client := &govmomi.Client{
		Client:         vimClient,
		SessionManager: vmomisession.NewManager(vimClient),
	}

	err = client.Login(ctx, url.UserPassword(conf.Username, conf.Token))
	if err != nil {
		v.NoteIssue(fmt.Errorf("Failed to validate operations credentials: %s", err))
		return
	}
	client.Logout(ctx)

	// confirm the RBAC configuration of the provided user
	// TODO: this can be dropped once we move to configuration the RBAC during creation
}

func (v *Validator) certificate(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	if len(input.CertPEM) == 0 && len(input.KeyPEM) == 0 {
		// if there's no data supplied then we're configuring without TLS
		log.Debug("Configuring without TLS due to empty key and cert buffers")
		return
	}

	// check the cert can be loaded
	_, err := tls.X509KeyPair(input.CertPEM, input.KeyPEM)
	v.NoteIssue(err)

	conf.HostCertificate = &config.RawCertificate{
		Key:  input.KeyPEM,
		Cert: input.CertPEM,
	}
}

func (v *Validator) certificateAuthorities(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	if len(input.ClientCAs) == 0 {
		// if there's no data supplied then we're configuring without client verification
		log.Debug("Configuring without client verification due to empty certificate authorities")
		return
	}

	// ensure TLS is configurable
	if len(input.CertPEM) == 0 {
		v.NoteIssue(errors.New("Certificate authority specified, but no TLS certificate provided"))
		return
	}

	// check a CA can be loaded
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(input.ClientCAs) {
		v.NoteIssue(errors.New("Unable to load certificate authority data"))
		return
	}

	conf.CertificateAuthorities = input.ClientCAs
}

func (v *Validator) registries(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// copy the list of insecure registries
	conf.InsecureRegistries = input.InsecureRegistries

	if len(input.RegistryCAs) == 0 {
		return
	}

	// Check if CAs can be loaded
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(input.RegistryCAs) {
		v.NoteIssue(errors.New("Unable to load certificate authority data for registry"))
		return
	}

	conf.RegistryCertificateAuthorities = input.RegistryCAs
}

func (v *Validator) compatibility(ctx context.Context, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// TODO: add checks such as datastore is acessible from target cluster
	errMsg := "Compatibility check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	// check session's datastore(s) exist
	_, err := v.Session.Datastore.AttachedClusterHosts(v.Context, v.Session.Cluster)
	v.NoteIssue(err)

	v.checkDatastoresAreWriteable(ctx, conf)
}

// looks up a datastore and adds it to the set
func (v *Validator) getDatastore(ctx context.Context, u *url.URL, datastoreSet map[types.ManagedObjectReference]*object.Datastore) map[types.ManagedObjectReference]*object.Datastore {
	if datastoreSet == nil {
		datastoreSet = make(map[types.ManagedObjectReference]*object.Datastore)
	}

	datastores, err := v.Session.Finder.DatastoreList(ctx, u.Host)
	v.NoteIssue(err)

	if len(datastores) != 1 {
		v.NoteIssue(errors.Errorf("Looking up datastore %s returned %d results.", u.String(), len(datastores)))
	}
	for _, d := range datastores {
		datastoreSet[d.Reference()] = d
	}
	return datastoreSet
}

// populates the v.datastoreSet "set" with datastore references generated from config
func (v *Validator) getAllDatastores(ctx context.Context, conf *config.VirtualContainerHostConfigSpec) map[types.ManagedObjectReference]*object.Datastore {
	// note that ImageStores, ContainerStores, and VolumeLocations
	// have just-different-enough types/structures that this cannot be made more succinct
	var datastoreSet map[types.ManagedObjectReference]*object.Datastore
	for _, u := range conf.ImageStores {
		datastoreSet = v.getDatastore(ctx, &u, datastoreSet)
	}
	for _, u := range conf.ContainerStores {
		datastoreSet = v.getDatastore(ctx, &u, datastoreSet)
	}
	for _, u := range conf.VolumeLocations {
		//skip non datastore volume stores
		if u.Scheme != common.DsScheme {
			continue
		}
		datastoreSet = v.getDatastore(ctx, u, datastoreSet)
	}
	return datastoreSet
}

func (v *Validator) checkDatastoresAreWriteable(ctx context.Context, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// gather compute host references
	clusterDatastores, err := v.Session.Cluster.Datastores(ctx)
	v.NoteIssue(err)

	// check that the cluster can see all of the datastores in question
	requestedDatastores := v.getAllDatastores(ctx, conf)
	validStores := make(map[types.ManagedObjectReference]*object.Datastore)
	// remove any found datastores from requested datastores
	for _, cds := range clusterDatastores {
		if requestedDatastores[cds.Reference()] != nil {
			delete(requestedDatastores, cds.Reference())
			validStores[cds.Reference()] = cds
		}
	}
	// if requestedDatastores is not empty, some requested datastores are not writable
	for _, store := range requestedDatastores {
		v.NoteIssue(errors.Errorf("Datastore %q is not accessible by the compute resource", store.Name()))
	}

	clusterHosts, err := v.Session.Cluster.Hosts(ctx)
	justOneHost := len(clusterHosts) == 1
	v.NoteIssue(err)

	for _, store := range validStores {
		aHosts, err := store.AttachedHosts(ctx)
		v.NoteIssue(err)
		clusterHosts = intersect(clusterHosts, aHosts)
	}

	if len(clusterHosts) == 0 {
		v.NoteIssue(errors.New("No single host can access all of the requested datastores. Installation cannot continue."))
	}

	if len(clusterHosts) == 1 && v.Session.IsVC() && !justOneHost {
		// if we have a cluster with >1 host to begin with, on VC, and only one host can talk to all the datastores, warn
		// on ESX and clusters with only one host to begin with, this warning would be redundant/irrelevant
		log.Warnf("Only one host can access all of the image/container/volume datastores. This may be a point of contention/performance degradation and HA/DRS may not work as intended.")
	}
}

// finds the intersection between two sets of HostSystem objects
func intersect(one []*object.HostSystem, two []*object.HostSystem) []*object.HostSystem {
	var result []*object.HostSystem
	for _, o := range one {
		for _, t := range two {
			if o.Reference() == t.Reference() {
				result = append(result, o)
			}
		}
	}
	return result
}

func (v *Validator) IsVC() bool {
	return v.isVC
}

func (v *Validator) AddDeprecatedFields(ctx context.Context, conf *config.VirtualContainerHostConfigSpec, input *data.Data) *data.InstallerData {
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

	dconfig.VCHSize.CPU.Reservation = int64(input.VCHCPUReservationsMHz)
	dconfig.VCHSize.CPU.Limit = int64(input.VCHCPULimitsMHz)
	dconfig.VCHSize.CPU.Shares = input.VCHCPUShares

	dconfig.VCHSize.Memory.Reservation = int64(input.VCHMemoryReservationsMB)
	dconfig.VCHSize.Memory.Limit = int64(input.VCHMemoryLimitsMB)
	dconfig.VCHSize.Memory.Shares = input.VCHMemoryShares

	return &dconfig
}
