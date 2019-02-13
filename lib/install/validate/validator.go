// Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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
	"strings"
	"time"

	"github.com/docker/go-units"

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
	"github.com/vmware/vic/lib/config/dynamic"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/opsuser"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/registry"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/optmanager"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const defaultSyslogPort = 514
const registryValidationTime = 10 * time.Second

type Validator struct {
	session *session.Session

	issues []error
}

func (v *Validator) Session() *session.Session {
	return v.session
}

func CreateFromSession(ctx context.Context, sess *session.Session) (*Validator, error) {
	defer trace.End(trace.Begin("", ctx))

	v := &Validator{}
	v.session = sess

	return v, nil
}

func NewSessionConfig(ctx context.Context, input *data.Data) (*session.Config, error) {
	op := trace.FromContext(ctx, "NewSessionConfig")

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
		if err := cert.FromURL(tURL, new(tls.Config)); err != nil {
			return nil, err
		}

		if cert.Err != nil {
			if !input.Force {
				// TODO: prompt user / check ./known_hosts
				op.Errorf("Failed to verify certificate for target=%s (thumbprint=%s)",
					tURL.Host, cert.ThumbprintSHA1)
				return nil, cert.Err
			}
		}

		input.Thumbprint = cert.ThumbprintSHA1
		op.Debugf("Accepting host %q thumbprint %s", tURL.Host, input.Thumbprint)
	}

	sessionconfig := &session.Config{
		Thumbprint: input.Thumbprint,
		Insecure:   input.Force,
	}

	// if a datacenter was specified, set it
	datacenterPath := tURL.Path
	if datacenterPath != "" {
		datacenterPath = strings.TrimPrefix(datacenterPath, "/")
		sessionconfig.DatacenterPath = datacenterPath
		// path needs to be stripped before we can use it as a service url
		tURL.Path = ""
	}

	sessionconfig.ClusterPath = input.ComputeResourcePath
	sessionconfig.Service = tURL.String()

	sessionconfig.CloneTicket = input.CloneTicket

	sessionconfig.UserAgent = version.UserAgent("vic-machine")

	return sessionconfig, nil
}

func NewSession(ctx context.Context, input *data.Data) (*session.Session, error) {
	op := trace.FromContext(ctx, "NewSession")

	sessionconfig, err := NewSessionConfig(op, input)
	if err != nil {
		return nil, err
	}

	s := session.NewSession(sessionconfig)
	s, err = s.Connect(ctx)
	if err != nil {
		return nil, err
	}

	finder := find.NewFinder(s.Client.Client, false)
	s.Finder = finder

	// Intentionally ignore any error returned by Populate
	_, err = s.Populate(ctx)
	if err != nil {
		op.Debugf("new validator Session.Populate: %s", err)
	}

	return s, nil
}

func NewValidator(ctx context.Context, input *data.Data) (*Validator, error) {
	op := trace.FromContext(ctx, "NewValidator")
	defer trace.End(trace.Begin("", op))

	s, err := NewSession(op, input)
	if err != nil {
		return nil, err
	}

	v, err := CreateFromSession(op, s)
	if err != nil {
		return nil, err
	}

	if strings.Contains(s.Config.DatacenterPath, "/") {
		detail := "--target should only specify datacenter in the path (e.g. https://addr/datacenter) - specify cluster, resource pool, or folder with --compute-resource"
		op.Error(detail)
		v.suggestDatacenter(op)
		return nil, errors.New(detail)
	}

	return v, nil
}

func (v *Validator) datacenter(op trace.Operation, allowEmptyDC bool) error {
	if allowEmptyDC && v.session.DatacenterPath == "" {
		return nil
	}
	if v.session.Datacenter != nil {
		v.session.DatacenterPath = v.session.Datacenter.InventoryPath
		return nil
	}
	var detail string
	if v.session.DatacenterPath != "" {
		detail = fmt.Sprintf("Datacenter %q in --target is not found", path.Base(v.session.DatacenterPath))
	} else {
		// this means multiple datacenter exists, but user did not specify it in --target
		detail = "Datacenter must be specified in --target (e.g. https://addr/datacenter)"
	}
	op.Error(detail)
	v.suggestDatacenter(op)
	return errors.New(detail)
}

func (v *Validator) listDatacenters(op trace.Operation) ([]string, error) {
	dcs, err := v.session.Finder.DatacenterList(op, "*")
	if err != nil {
		return nil, fmt.Errorf("unable to list datacenters: %s", err)
	}

	if len(dcs) == 0 {
		return nil, nil
	}

	matches := make([]string, len(dcs))
	for i, d := range dcs {
		matches[i] = d.Name()
	}
	return matches, nil
}

// suggestDatacenter suggests all datacenters on the target
func (v *Validator) suggestDatacenter(op trace.Operation) {
	defer trace.End(trace.Begin("", op))

	op.Info("Suggesting valid values for datacenter in --target")

	dcs, err := v.listDatacenters(op)
	if err != nil {
		op.Error(err)
		return
	}

	if len(dcs) == 0 {
		op.Info("No datacenters found")
		return
	}

	op.Info("Suggested datacenters:")
	for _, d := range dcs {
		op.Infof("  %q", d)
	}
	return
}

func (v *Validator) NoteIssue(err error) {
	if err != nil {
		v.issues = append(v.issues, err)
	}
}

func (v *Validator) ListIssues(ctx context.Context) error {
	op := trace.FromContext(ctx, "ListIssues")
	defer trace.End(trace.Begin("", op))

	if len(v.issues) == 0 {
		return nil
	}

	op.Error("--------------------")
	for _, err := range v.issues {
		op.Error(err)
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
func (v *Validator) Validate(ctx context.Context, input *data.Data, allowEmptyDC bool) (*config.VirtualContainerHostConfigSpec, error) {
	op := trace.FromContext(ctx, "Validate")
	defer trace.End(trace.Begin("", op))
	op.Info("Validating supplied configuration")

	conf := &config.VirtualContainerHostConfigSpec{}

	if err := v.datacenter(op, allowEmptyDC); err != nil {
		return conf, err
	}

	v.basics(op, input, conf)
	v.target(op, input, conf)
	v.credentials(op, input, conf)
	v.compute(op, input, conf)
	v.storage(op, input, conf)
	v.network(op, input, conf)
	// FIXME ATC DEBT setting this value needs to be moved to Dispatcher
	// https://github.com/vmware/vic/issues/6803
	ok := v.checkPersistNetworkBacking(op, true)
	if !ok {
		err := v.configureVCenter(op)
		if err != nil {
			op.Errorf("%s", err)
			op.Errorf("vCenter settings update FAILED")
		}
	}
	v.checkFirewall(op, conf)
	v.checkPersistNetworkBacking(op, false)
	v.CheckLicense(op)
	v.checkDRS(op, input)
	v.checkVMGroup(op, input, conf) // Depends on a side-effect of the checkDRS method.

	v.certificate(op, input, conf)
	v.certificateAuthorities(op, input, conf)
	v.registries(op, input, conf)

	// Perform the higher level compatibility and consistency checks
	v.compatibility(op, conf)

	v.syslog(op, conf, input)

	// TODO: determine if this is where we should turn the noted issues into message
	return conf, v.ListIssues(op)

}

func (v *Validator) ValidateTarget(ctx context.Context, input *data.Data, allowEmptyDC bool) (*config.VirtualContainerHostConfigSpec, error) {
	op := trace.FromContext(ctx, "ValidateTarget")
	defer trace.End(trace.Begin("", op))

	conf := &config.VirtualContainerHostConfigSpec{}

	op.Info("Validating target")
	if err := v.datacenter(op, allowEmptyDC); err != nil {
		return conf, err
	}
	v.target(op, input, conf)
	return conf, v.ListIssues(op)
}

func (v *Validator) ValidateStorageQuota(ctx context.Context, quotaGB int, conf *config.VirtualContainerHostConfigSpec, vch *vm.VirtualMachine) (int64, error) {
	op := trace.FromContext(ctx, "ValidateStorageQuota")
	defer trace.End(trace.Begin("", op))
	op.Info("Validating storage quota")

	vmStorageUsage := int64(0)
	imageStorageUsage := int64(0)
	var err error
	if vch != nil {
		vmStorageUsage, err = v.getVMStorageUsage(op, vch)
		if err != nil {
			return 0, err
		}
		imageStorageUsage, err = v.getImageStorageUsage(op, conf, vch)
		if err != nil {
			return 0, err
		}
	}
	quota := int64(quotaGB) * units.GiB
	if quota > v.getDatastoreFreeSpace(op)+vmStorageUsage+imageStorageUsage {
		return 0, errors.New("Storage quota exceeds datastore free space")
	}
	return quota, nil
}

func (v *Validator) getVMStorageUsage(op trace.Operation, vch *vm.VirtualMachine) (int64, error) {
	vchFolder, err := vch.Folder(op)
	if err != nil {
		return 0, err
	}
	children, err := vchFolder.Children(op)
	if err != nil {
		return 0, err
	}

	total := (int64)(0)
	for _, child := range children {
		vmObj, ok := child.(*object.VirtualMachine)
		if vmObj.Reference().Value == vch.Reference().Value {
			continue
		}

		if ok {
			vm2 := vm.NewVirtualMachineFromVM(op, v.session, vmObj)

			folder, err := vm2.DatastoreFolderName(op)
			if err != nil {
				op.Warn("Failed to get container vm information")
				continue
			}
			// get a ds helper for this ds url
			dsHelper, err := datastore.NewHelper(trace.NewOperation(op, "datastore helper creation"), v.session,
				v.session.Datastore, folder)
			if err != nil {
				op.Warn("Failed to get datastore information for container vm")
				continue
			}

			size, err := dsHelper.GetFilesSize(op, "", true, "*.vmdk", "*.vswp", "*.log")
			if err != nil {
				op.Warn("Failed to browse container vm files")
				continue
			}
			total += size
		}
	}

	return total, nil
}

func (v *Validator) getImageStorageUsage(op trace.Operation, conf *config.VirtualContainerHostConfigSpec, vch *vm.VirtualMachine) (int64, error) {
	uuid, err := vch.UUID(op)
	if err != nil {
		return 0, err
	}

	imageURL := conf.ImageStores[0]
	// get a ds helper for this ds url
	dsHelper, err := datastore.NewHelper(trace.NewOperation(op, "datastore helper creation"), v.session,
		v.session.Datastore, fmt.Sprintf("%s/%s/%s", imageURL.Path, constants.StorageParentDir, uuid))
	if err != nil {
		return 0, err
	}

	size, err := dsHelper.GetFilesSize(op, "", true, "*.vmdk")
	if err != nil {
		return 0, err
	}

	return size, nil
}

func (v *Validator) basics(op trace.Operation, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin("", op))

	// TODO: ensure that displayname doesn't violate constraints (length, characters, etc)
	conf.SetName(input.DisplayName)

	if input.Debug.Debug != nil {
		conf.SetDebug(*input.Debug.Debug)
	}

	conf.Name = input.DisplayName
	conf.Version = version.GetBuild()

	scratchSize, err := units.FromHumanSize(input.ScratchSize)
	if err != nil { // TODO set minimum size of scratch disk
		v.NoteIssue(errors.Errorf("Invalid default image size %s provided; error from parser: %s", input.ScratchSize, err.Error()))
	} else {
		conf.ScratchSize = scratchSize / units.KB
		op.Debugf("Setting scratch image size to %d KB in VCHConfig", conf.ScratchSize)
	}

	if input.ContainerNameConvention != "" {
		// ensure token is present
		if !strings.Contains(input.ContainerNameConvention, string(config.IDToken)) && !strings.Contains(input.ContainerNameConvention, string(config.NameToken)) {
			v.NoteIssue(errors.Errorf("Container name convention must include %s or %s token", config.IDToken, config.NameToken))
		}

		// coarse check - only enforce that there's enough capcity for a shortID
		// name lengths are many and vary significantly so much harder to provide sanity checks for - they will truncate when convention is applied.
		if len(input.ContainerNameConvention)-len(config.IDToken) >= constants.MaxVMNameLength-constants.ShortIDLen {
			v.NoteIssue(errors.Errorf("Container name convetion exceeds maximum length (%d, discounting tokens)", constants.MaxVMNameLength-constants.ShortIDLen))
		}
	}

	conf.ContainerNameConvention = input.ContainerNameConvention
}

func (v *Validator) checkSessionSet() []string {
	var errs []string

	if v.session.Datastore == nil {
		errs = append(errs, "datastore not set")
	}
	if v.session.Cluster == nil {
		errs = append(errs, "cluster not set")
	}

	return errs
}

func (v *Validator) sessionValid(op trace.Operation, errMsg string) bool {
	if c := v.checkSessionSet(); len(c) > 0 {
		op.Error(errMsg)
		for _, e := range c {
			op.Errorf("  %s", e)
		}
		v.NoteIssue(errors.New(errMsg))
		return false
	}
	return true
}

func (v *Validator) target(op trace.Operation, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin("", op))

	// check if host is managed by VC
	v.managedbyVC(op)
}

func (v *Validator) managedbyVC(op trace.Operation) {
	defer trace.End(trace.Begin("", op))

	if v.isVC() {
		return
	}
	host, err := v.session.Finder.DefaultHostSystem(op)
	if err != nil {
		v.NoteIssue(fmt.Errorf("Failed to get host system: %s", err))
		return
	}

	var mh mo.HostSystem

	if err = host.Properties(op, host.Reference(), []string{"summary.managementServerIp"}, &mh); err != nil {
		v.NoteIssue(fmt.Errorf("Failed to get host properties: %s", err))
		return
	}

	if ip := mh.Summary.ManagementServerIp; ip != "" {
		v.NoteIssue(fmt.Errorf("Target is managed by vCenter server %q, please change --target to vCenter server address or select a standalone ESXi", ip))
	}
	return
}

func (v *Validator) credentials(op trace.Operation, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	// empty string for password is horrific, but a legitimate scenario especially in isolated labs
	if input.OpsCredentials.OpsUser == nil || input.OpsCredentials.OpsPassword == nil {
		v.NoteIssue(errors.New("User/password for the operations user has not been set"))
		return
	}

	// check target with ops credentials
	u := input.Target.URL

	conf.Username = *input.OpsCredentials.OpsUser
	conf.Token = *input.OpsCredentials.OpsPassword
	conf.TargetThumbprint = input.Thumbprint

	// If Grant Perms has  been explicitly requested (either true or false)
	// set it to the new value. Otherwise leave the value in conf as it is
	if input.OpsCredentials.GrantPerms != nil {
		if *input.OpsCredentials.GrantPerms {
			// Set Grant Permissions level
			conf.SetGrantPerms()
		} else {
			conf.ClearGrantPerms()
		}
	}

	// If Grant Perms is set trying adding ReadOnly role to the Datacenter
	// for the ops-user. This is necessary since the Login operation below
	// fails if the ops-user has no permissions.
	//
	// FIXME DEBT.
	// https://github.com/vmware/vic/issues/6870
	// Notice that this operation should not be performed from the Validator.
	// Eventually, this must be moved to the Dispatcher as the Validator
	// should not modify VC configuration.
	if conf.ShouldGrantPerms() {
		err := opsuser.GrantDCReadOnlyPerms(op, v.session, conf)
		if err != nil {
			v.NoteIssue(fmt.Errorf("Failed to validate operations credentials: %s", err))
			return
		}
	}

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

	vimClient, err := vim25.NewClient(op, soapClient)
	if err != nil {
		v.NoteIssue(fmt.Errorf("Failed to create client for validation of operations credentials: %s", err))
		return
	}

	client := &govmomi.Client{
		Client:         vimClient,
		SessionManager: vmomisession.NewManager(vimClient),
	}

	err = client.Login(op, url.UserPassword(conf.Username, conf.Token))
	if err != nil {
		v.NoteIssue(fmt.Errorf("Failed to validate operations credentials: %s", err))
		return
	}
	client.Logout(op)

	// confirm the RBAC configuration of the provided user
	// TODO: this can be dropped once we move to configuration the RBAC during creation
}

func (v *Validator) certificate(op trace.Operation, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin("", op))

	if len(input.CertPEM) == 0 && len(input.KeyPEM) == 0 {
		// if there's no data supplied then we're configuring without TLS
		op.Debug("Configuring without TLS due to empty key and cert buffers")
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

func (v *Validator) certificateAuthorities(op trace.Operation, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin("", op))

	if len(input.ClientCAs) == 0 {
		// if there's no data supplied then we're configuring without client verification
		op.Debug("Configuring without client verification due to empty certificate authorities")
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

func (v *Validator) registries(op trace.Operation, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin("", op))

	// Check if CAs can be loaded
	pool := x509.NewCertPool()
	if len(input.RegistryCAs) > 0 {
		if !pool.AppendCertsFromPEM(input.RegistryCAs) {
			v.NoteIssue(errors.New("Unable to load certificate authority data for registry"))
			return
		}
	}

	conf.RegistryCertificateAuthorities = input.RegistryCAs

	// test reachability
	insecureRegistries, whitelistRegistries, err := v.reachableRegistries(op, input, pool)
	if err != nil {
		v.NoteIssue(err)
		return
	}

	// copy the list of insecure registries
	conf.InsecureRegistries = insecureRegistries

	// copy the list of whitelist registries
	conf.RegistryWhitelist = whitelistRegistries

	// create vic-machine info message
	msg := v.friendlyRegistryList("Insecure registries", conf.InsecureRegistries)
	if msg != "" {
		op.Info(msg)
	}
	msg = v.friendlyRegistryList("Whitelist registries", conf.RegistryWhitelist)
	if msg != "" {
		op.Info(msg)
	}

	if len(input.RegistryCAs) == 0 {
		return
	}
}

func (v *Validator) friendlyRegistryList(registryType string, registryList []string) string {
	if len(registryList) == 0 {
		return ""
	}

	return registryType + " = " + strings.Join(registryList, ", ")
}

// Validate registries are reachable.  Secure registries that are not specified as insecure are validated with the
// CA certs passed into vic-machine.
func (v *Validator) reachableRegistries(op trace.Operation, input *data.Data, pool *x509.CertPool) (insecureRegistries []string, whitelistRegistries []string, err error) {
	secureRegistriesSet, err := dynamic.ParseRegistries(input.WhitelistRegistries)
	if err != nil {
		return nil, nil, err
	}

	insecureRegistriesSet, err := dynamic.ParseRegistries(input.InsecureRegistries)
	if err != nil {
		return nil, nil, err
	}

	// Test insecure registries' reachability
	for _, r := range insecureRegistriesSet {
		r, ok := r.(registry.URLEntry)
		if !ok {
			err = fmt.Errorf("invalid insecure registry entry: %s", r)
			v.NoteIssue(err)
			return nil, nil, err
		}

		// Remove intersection between insecure registries and whitelist registries from whitelist set so
		// we can ensure we test the exclusion set with certs
		for idx, s := range secureRegistriesSet {
			if s.IsURL() && r.Match(s.String()) {
				// remove the insecure registry from list of registries to get validated against certs
				secureRegistriesSet = append(secureRegistriesSet[:idx], secureRegistriesSet[idx+1:]...)
				break
			}
		}

		// Make sure address is not a wildcard domain or CIDR.  If it is, do not validate.
		if strings.HasPrefix(r.URL().Host, "*") {
			op.Debugf("Skipping registry validation for %s", r)
			continue
		}

		schemes := []string{""}
		if r.URL().Scheme == "" {
			schemes = []string{"https", "http"}
		}

		rs := r.String()
		for _, s := range schemes {
			if _, err = registry.Reachable(op, rs, s, "", "", nil, registryValidationTime, true); err == nil {
				break
			}
		}

		if err != nil {
			op.Warnf("Unable to confirm insecure registry %s is a valid registry at this time.", r)
		} else {
			op.Debugf("Insecure registry %s confirmed.", r)
		}
	}

	// Test secure registries' reachability
	for _, w := range secureRegistriesSet {
		// Make sure address is not a wildcard domain or CIDR.  If it is, do not validate.
		if w.IsCIDR() {
			op.Debugf("Skipping registry validation for %s", w)
			continue
		}

		w, ok := w.(registry.URLEntry)
		if !ok {
			op.Debugf("Skipping registry validation for %s", w)
			continue
		}

		if strings.HasPrefix(w.URL().Host, "*") {
			op.Debugf("Skipping registry validation for %s", w)
			continue
		}

		scheme := w.URL().Scheme
		if scheme == "" {
			scheme = "https"
		}

		if _, err = registry.Reachable(op, w.String(), scheme, "", "", pool, registryValidationTime, false); err != nil {
			op.Warnf("Unable to confirm secure registry %s is a valid registry at this time.", w)
		} else {
			op.Debugf("Secure registry %s confirmed.", w)
		}
	}

	// Return output
	insecureRegistries = input.InsecureRegistries
	// If vic-machine had whitelist registry specified
	if len(input.WhitelistRegistries) > 0 {
		// ignoring error since default merge policy is union, so should never return
		// an error
		// #nosec: Errors unhandled.
		m, _ := secureRegistriesSet.Merge(insecureRegistriesSet, nil)
		whitelistRegistries = m.Strings()
	}

	err = nil
	return
}

func (v *Validator) compatibility(op trace.Operation, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin("", op))

	// TODO: add checks such as datastore is acessible from target cluster
	errMsg := "Compatibility check SKIPPED"
	if !v.sessionValid(op, errMsg) {
		return
	}

	// check session's datastore(s) exist
	_, err := v.session.Datastore.AttachedClusterHosts(op, v.session.Cluster)
	v.NoteIssue(err)

	v.checkDatastoresAreWriteable(op, conf)
}

// looks up a datastore and adds it to the set
func (v *Validator) getDatastore(op trace.Operation, u *url.URL, datastoreSet map[types.ManagedObjectReference]*object.Datastore) map[types.ManagedObjectReference]*object.Datastore {
	if datastoreSet == nil {
		datastoreSet = make(map[types.ManagedObjectReference]*object.Datastore)
	}

	datastores, err := v.session.Finder.DatastoreList(op, u.Host)
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
func (v *Validator) getAllDatastores(op trace.Operation, conf *config.VirtualContainerHostConfigSpec) map[types.ManagedObjectReference]*object.Datastore {
	// note that ImageStores, ContainerStores, and VolumeLocations
	// have just-different-enough types/structures that this cannot be made more succinct
	var datastoreSet map[types.ManagedObjectReference]*object.Datastore
	for _, u := range conf.ImageStores {
		datastoreSet = v.getDatastore(op, &u, datastoreSet)
	}
	for _, u := range conf.ContainerStores {
		datastoreSet = v.getDatastore(op, &u, datastoreSet)
	}
	for _, u := range conf.VolumeLocations {
		//skip non datastore volume stores
		if u.Scheme != common.DsScheme {
			continue
		}
		datastoreSet = v.getDatastore(op, u, datastoreSet)
	}
	return datastoreSet
}

func (v *Validator) checkDatastoresAreWriteable(op trace.Operation, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin("", op))

	// gather compute host references
	clusterDatastores, err := v.session.Cluster.Datastores(op)
	v.NoteIssue(err)

	// check that the cluster can see all of the datastores in question
	requestedDatastores := v.getAllDatastores(op, conf)
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

	clusterHosts, err := v.session.Cluster.Hosts(op)
	justOneHost := len(clusterHosts) == 1
	v.NoteIssue(err)

	for _, store := range validStores {
		aHosts, err := store.AttachedHosts(op)
		v.NoteIssue(err)
		clusterHosts = intersect(clusterHosts, aHosts)
	}

	if len(clusterHosts) == 0 {
		v.NoteIssue(errors.New("No single host can access all of the requested datastores. Installation cannot continue."))
	}

	if len(clusterHosts) == 1 && v.session.IsVC() && !justOneHost {
		// if we have a cluster with >1 host to begin with, on VC, and only one host can talk to all the datastores, warn
		// on ESX and clusters with only one host to begin with, this warning would be redundant/irrelevant
		op.Warn("Only one host can access all of the image/volume datastores. This may be a point of contention/performance degradation and HA/DRS may not work as intended.")
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

func (v *Validator) isVC() bool {
	return v.session.IsVC()
}

func (v *Validator) AddDeprecatedFields(ctx context.Context, conf *config.VirtualContainerHostConfigSpec, input *data.Data) *data.InstallerData {
	op := trace.FromContext(ctx, "AddDeprecatedFields")
	defer trace.End(trace.Begin("", op))

	dconfig := data.InstallerData{}

	cpuLimit := int64(input.NumCPUs)
	memLimit := int64(input.MemoryMB)
	dconfig.ApplianceSize.CPU.Limit = &cpuLimit
	dconfig.ApplianceSize.Memory.Limit = &memLimit

	if v.session.Datacenter != nil {
		dconfig.Datacenter = v.session.Datacenter.Reference()
		dconfig.DatacenterName = v.session.Datacenter.Name()
	} else {
		op.Debug("session datacenter is nil")
	}

	if v.session.Cluster != nil {
		dconfig.Cluster = v.session.Cluster.Reference()
		dconfig.ClusterPath = v.session.Cluster.InventoryPath
	} else {
		op.Debug("session cluster is nil")
	}

	dconfig.ResourcePoolPath = v.session.PoolPath

	op.Debugf("Datacenter: %q, Cluster: %q, Resource Pool: %q", dconfig.DatacenterName, dconfig.Cluster, dconfig.ResourcePoolPath)

	if input.VCHCPUReservationsMHz != nil {
		cpuReserve := int64(*input.VCHCPUReservationsMHz)
		dconfig.VCHSize.CPU.Reservation = &cpuReserve
	}
	if input.VCHCPULimitsMHz != nil {
		cpuLimit := int64(*input.VCHCPULimitsMHz)
		dconfig.VCHSize.CPU.Limit = &cpuLimit
	}
	dconfig.VCHSize.CPU.Shares = input.VCHCPUShares

	if input.VCHMemoryReservationsMB != nil {
		memReserve := int64(*input.VCHMemoryReservationsMB)
		dconfig.VCHSize.Memory.Reservation = &memReserve
	}
	if input.VCHMemoryLimitsMB != nil {
		memLimit := int64(*input.VCHMemoryLimitsMB)
		dconfig.VCHSize.Memory.Limit = &memLimit
	}
	dconfig.VCHSize.Memory.Shares = input.VCHMemoryShares

	return &dconfig
}

func (v *Validator) syslog(op trace.Operation, conf *config.VirtualContainerHostConfigSpec, input *data.Data) {
	defer trace.End(trace.Begin("", op))

	if input.SyslogConfig.Addr == nil {
		return
	}

	u := input.SyslogConfig.Addr
	network := u.Scheme
	if len(network) == 0 {
		v.NoteIssue(errors.New("syslog address does not have network specified"))
		return
	}

	switch network {
	case "udp", "tcp":
	default:
		v.NoteIssue(fmt.Errorf("syslog address transport should be udp or tcp"))
		return
	}

	host := u.Host
	if len(host) == 0 {
		v.NoteIssue(errors.New("syslog address host not specified"))
		return
	}

	if u.Port() == "" {
		host += fmt.Sprintf(":%d", defaultSyslogPort)
	}

	conf.Diagnostics.SysLogConfig = &executor.SysLogConfig{
		Network: network,
		RAddr:   host,
	}
}

// FIXME ATC DEBT setting this value needs to be moved to Dispatcher
// https://github.com/vmware/vic/issues/6803
// set PersistNetworkBacking key to "true"
func (v *Validator) configureVCenter(ctx context.Context) error {
	op := trace.FromContext(ctx, "Set vCenter serial port backing")
	defer trace.End(trace.Begin("", op))

	errMsg := "Set vCenter settings SKIPPED"
	if !v.sessionValid(op, errMsg) {
		return nil
	}
	if !v.isVC() {
		op.Debug(errMsg)
		return nil
	}

	err := optmanager.UpdateOptionValue(ctx, v.session, persistNetworkBackingKey, "true")
	if err != nil {
		msg := fmt.Sprintf("Failed to set required value \"true\" for %s: %s", persistNetworkBackingKey, err)
		return errors.New(msg)
	}
	op.Infof("Set vCenter settings OK")
	return nil
}
