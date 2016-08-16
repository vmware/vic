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
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	units "github.com/docker/go-units"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
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

func CreateFromVCHConfig(ctx context.Context, vch *config.VirtualContainerHostConfigSpec, sess *session.Session) (*Validator, error) {
	defer trace.End(trace.Begin(""))

	v := &Validator{}
	v.Session = sess
	v.Context = ctx

	return v, nil
}

func NewValidator(ctx context.Context, input *data.Data) (*Validator, error) {
	v, err := CreateNoDCCheck(ctx, input)
	if err != nil {
		return nil, err
	}

	return v, v.datacenter()
}

func CreateNoDCCheck(ctx context.Context, input *data.Data) (*Validator, error) {
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

	// only allow the datacenter to be specified in the taget url, if any
	pElems := strings.Split(v.DatacenterPath, "/")
	if len(pElems) > 2 {
		detail := "--target should only specify datacenter in the path (e.g. https://addr/datacenter) - specify cluster, resource pool, or folder with --compute-resource"
		log.Error(detail)
		v.suggestDatacenter()
		return nil, errors.New(detail)
	}

	return v, nil
}

func (v *Validator) datacenter() error {
	if v.Session.Datacenter == nil {
		detail := "Datacenter must be specified in --target (e.g. https://addr/datacenter)"
		log.Error(detail)
		v.suggestDatacenter()
		return errors.New(detail)
	}
	v.DatacenterPath = v.Session.Datacenter.InventoryPath
	return nil
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

	v.basics(ctx, input, conf)

	v.target(ctx, input, conf)
	v.compute(ctx, input, conf)
	v.storage(ctx, input, conf)
	v.network(ctx, input, conf)
	v.CheckFirewall(ctx)
	v.CheckLicense(ctx)
	v.CheckDrs(ctx)

	v.certificate(ctx, input, conf)

	// Perform the higher level compatibility and consistency checks
	v.compatibility(ctx, conf)

	// TODO: determine if this is where we should turn the noted issues into message
	return conf, v.ListIssues()

}

func (v *Validator) ValidateTarget(ctx context.Context, input *data.Data) (*config.VirtualContainerHostConfigSpec, error) {
	defer trace.End(trace.Begin(""))
	conf := &config.VirtualContainerHostConfigSpec{}

	log.Infof("Validating target")
	v.target(ctx, input, conf)
	return conf, v.ListIssues()
}

func (v *Validator) basics(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// TODO: ensure that displayname doesn't violate constraints (length, characters, etc)
	conf.SetName(input.DisplayName)
	conf.SetDebug(input.Debug.Debug)
	conf.Name = input.DisplayName

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

	targetURL := input.Target.URLWithoutPassword()
	if !v.IsVC() {
		var err error
		targetURL, err = url.Parse(v.Session.Service)
		if err != nil {
			v.NoteIssue(fmt.Errorf("Error processing target after transformation to SOAP endpoint: %q: %s", v.Session.Service, err))
			return
		}

		// ESXi requires user/password to be encoded in the Target URL
		// However, this gets lost when the URL is Marshaled
		conf.UserPassword = targetURL.User.String()
	}

	// check if host is managed by VC
	v.managedbyVC(ctx)

	conf.Target = *targetURL
	conf.Insecure = input.Insecure

	// TODO: more checks needed here if specifying service account for VCH
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

func (v *Validator) compatibility(ctx context.Context, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// TODO: add checks such as datastore is acessible from target cluster
	errMsg := "Compatibility check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	_, err := v.Session.Datastore.AttachedClusterHosts(v.Context, v.Session.Cluster)
	v.NoteIssue(err)
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
