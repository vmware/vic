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
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/cmd/vic-machine/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/session"

	"github.com/vmware/govmomi/vim25/types"

	"golang.org/x/net/context"
)

type Validator struct {
	TargetPath            string
	DatacenterName        string
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

	issues []error
}

func NewValidator(ctx context.Context, input *data.Data) (*Validator, error) {
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
	v.DatacenterName = tURL.Path
	if v.DatacenterName != "" {
		sessionconfig.DatacenterPath = v.DatacenterName
		// path needs to be stripped before we can use it as a service url
		tURL.Path = ""
	}

	sessionconfig.Service = tURL.String()

	v.Session, err = session.NewSession(sessionconfig).Create(v.Context)
	if err != nil {
		log.Errorf("Failed to create session: %s", err)
		return nil, err
	}
	v.Session, err = v.Session.Connect(v.Context)
	if err != nil {
		return nil, err
	}

	v.Session.Populate(ctx)

	if v.Session.Datacenter == nil {
		detail := "Target should specify datacenter when there are multiple possibilities, e.g. https://addr/datacenter"
		log.Error(detail)
		// TODO: list available datacenters
		return nil, errors.New(detail)
	}

	return v, nil

}

func (v *Validator) NoteIssue(err error) {
	if err != nil {
		v.issues = append(v.issues, err)
	}
}

func (v *Validator) ListIssues() error {
	if len(v.issues) == 0 {
		return nil
	}

	for _, err := range v.issues {
		fmt.Println(err)
	}

	return errors.New("Validation of configuration failed")
}

// Validate runs through various validations, starting with basics such as naming, moving onto vSphere entities
// and then the compatibility between those entities. It assembles a set of issues that are found for reporting.
func (v *Validator) Validate(ctx context.Context, input *data.Data) (*metadata.VirtualContainerHostConfigSpec, error) {
	log.Infof("Validating supplied configuration")

	conf := &metadata.VirtualContainerHostConfigSpec{}

	v.Basics(ctx, input, conf)

	v.Target(ctx, input, conf)
	v.Compute(ctx, input, conf)
	v.Storage(ctx, input, conf)
	v.Network(ctx, input, conf)

	v.Certificate(ctx, input, conf)

	// Perform the higher level compatibility and consistency checks
	v.Compatibility(ctx, conf)

	// TODO: determine if this is where we should turn the noted issues into message
	return conf, v.ListIssues()

}

func (v *Validator) Basics(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	// TODO: ensure that displayname doesn't violate constraints (length, characters, etc)
	conf.SetName(input.DisplayName)

	conf.Name = input.DisplayName
}

func (v *Validator) Compute(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	// Compute
	pool, err := v.GetResourcePool(ctx, input.ComputeResourcePath)
	v.NoteIssue(err)
	moref := pool.Reference()
	conf.AddComputeResource(&moref)
}

func (v *Validator) Storage(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {

	// Image Store
	ds, err := v.datastoreHelper(ctx, input.ImageDatastoreName)
	v.NoteIssue(err)
	conf.AddImageStore(ds)

	// TODO: add volume locations

}

func (v *Validator) Network(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
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
		},
	})

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

	// Bridge net -
	//   vCenter: must exist and must be a DPG
	//   ESX: doesn't need to exist
	bridgeMoref, err := v.dpgHelper(ctx, input.BridgeNetworkName)
	if err != nil {
		if _, ok := err.(*find.NotFoundError); !ok || v.Session.Client.IsVC() {
			v.NoteIssue(err)
		}
	}

	gateway, ok := input.MappedNetworksGateway["bridge"]
	if !ok {
		// TODO: put the default somewhere less inline
		ip, mask, _ := net.ParseCIDR("172.16.0.1/24")
		gateway = &net.IPNet{
			IP:   ip,
			Mask: mask.Mask,
		}
	}
	conf.AddNetwork(&metadata.NetworkEndpoint{
		Common: metadata.Common{
			Name: "bridge",
		},
		Static: gateway,
		Network: metadata.ContainerNetwork{
			Common: metadata.Common{
				Name: "bridge",
				ID:   bridgeMoref,
			},
			Gateway: *gateway,
		},
	})

	// add mapped networks
	//   these should be a distributed port groups in vCenter
	for name, net := range input.MappedNetworks {
		moref, err := v.dpgHelper(ctx, net)
		v.NoteIssue(err)
		mappedNet := &metadata.ContainerNetwork{
			Common: metadata.Common{
				Name: name,
				ID:   moref,
			},
		}

		gateway, ok = input.MappedNetworksGateway[name]
		if ok {
			mappedNet.Gateway = *gateway
		}
		conf.AddContainerNetwork(mappedNet)
	}
}

func (v *Validator) Target(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	targetURL := *input.URL
	if targetURL.Scheme == "" {
		targetURL.Scheme = "https"
	}
	targetURL.Path = "sdk"
	targetURL.User = url.UserPassword(input.User, *input.Password)

	conf.Target = targetURL
	conf.Insecure = input.Insecure

	// TODO: more checks needed here if specifying service account for VCH
}

func (v *Validator) Certificate(ctx context.Context, input *data.Data, conf *metadata.VirtualContainerHostConfigSpec) {
	// check the cert can be loaded
	_, err := tls.X509KeyPair(input.CertPEM, input.KeyPEM)
	v.NoteIssue(err)

	conf.HostCertificate = metadata.RawCertificate{
		Key:  input.KeyPEM,
		Cert: input.CertPEM,
	}
}

func (v *Validator) Compatibility(ctx context.Context, conf *metadata.VirtualContainerHostConfigSpec) {
	// TODO: add checks such as datastore is acessible from target cluster
	_, err := v.Session.Datastore.AttachedClusterHosts(v.Context, v.Session.Cluster)
	v.NoteIssue(err)
}

func (v *Validator) networkHelper(ctx context.Context, path string) (string, error) {
	nets, err := v.Session.Finder.NetworkList(ctx, path)
	if err != nil {
		// TODO: error message about no such match and how to get a network list
		return "", errors.New("no such network " + path)
	}
	if len(nets) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return "", errors.New("ambiguous network " + path)
	}

	moref := nets[0].Reference()
	return moref.String(), nil
}

func (v *Validator) dpgMorefHelper(ctx context.Context, ref string) (string, error) {
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
	if v.Session.Client.IsVC() {
		_, dpg := net.(*object.DistributedVirtualPortgroup)
		if !dpg {
			return "", fmt.Errorf("%s is not a Distributed Port Group", ref)
		}
	}

	return ref, nil
}

func (v *Validator) dpgHelper(ctx context.Context, path string) (string, error) {
	nets, err := v.Session.Finder.NetworkList(ctx, path)
	if err != nil {
		// TODO: error message about no such match and how to get a network list
		return "", errors.New("no such network " + path)
	}
	if len(nets) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return "", errors.New("ambiguous network " + path)
	}

	// ensure that the type of the network is a Distributed Port Group if the target is a vCenter
	// if it's not then any network suffices
	if v.Session.Client.IsVC() {
		_, dpg := nets[0].(*object.DistributedVirtualPortgroup)
		if !dpg {
			return "", fmt.Errorf("%s is not a Distributed Port Group", path)
		}
	}

	moref := nets[0].Reference()
	return moref.String(), nil
}

func (v *Validator) datastoreHelper(ctx context.Context, path string) (*url.URL, error) {
	dsURL, err := url.Parse(path)
	if err != nil {
		// try treating it as a plain path
		pathElements := strings.Split(path, "/")
		if pathElements[0] == "" {
			// TODO: error about requiring datastore path and how to get a datastore list
			return nil, errors.New("requires datastore name")
		}

		dsURL.Scheme = "ds://"
		dsURL.Host = pathElements[0]
		dsURL.Path = strings.Join(pathElements[1:], "/")
	}

	// if a datastore name (e.g. "datastore1") is specifed with no decoration then this
	// is interpreted as the Path
	if dsURL.Host == "" && dsURL.Path != "" {
		dsURL.Host = dsURL.Path
		dsURL.Path = ""
	}

	stores, err := v.Session.Finder.DatastoreList(ctx, dsURL.Host)
	if err != nil {
		// TODO: error message about no such match and how to get a network list
		return nil, fmt.Errorf("no such datastore %#v", dsURL)
	}
	if len(stores) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return nil, errors.New("ambiguous datastore " + dsURL.Host)
	}

	// FIXME: commented out until components can consume moid
	// dsURL.Host = stores[0].Reference().Value

	return dsURL, nil
}

func (v *Validator) GetResourcePool(ctx context.Context, path string) (*object.ResourcePool, error) {
	pools, err := v.Session.Finder.ResourcePoolList(ctx, path)
	if err != nil {
		// TODO: error message about no such match and how to get a network list
		return nil, errors.New("no such compute resource " + path)
	}
	if len(pools) > 1 {
		// TODO: error about required disabmiguation and list entries in nets
		return nil, errors.New("ambiguous compute resource " + path)
	}

	return pools[0], nil
}

func (v *Validator) ParseComputeResourcePath(rp string) error {
	resources := strings.Split(rp, "/")
	if len(resources) < 2 || resources[1] == "" {
		err := errors.Errorf("Could not determine datacenter from specified --compute-resource path: %s", rp)
		return err
	}
	v.DatacenterName = resources[1]
	v.ClusterPath = strings.Split(rp, "/Resources")[0]

	if v.ClusterPath == "" {
		err := errors.Errorf("Could not determine cluster from specified --compute-resource path: %s", rp)
		return err
	}
	v.ResourcePoolPath = rp
	return nil
}

func (v *Validator) addDeprecatedFields(ctx context.Context, conf *metadata.VirtualContainerHostConfigSpec, input *data.Data) *management.InstallerData {

	dconfig := management.InstallerData{}

	dconfig.ApplianceSize.CPU.Limit = int64(input.NumCPUs)
	dconfig.ApplianceSize.Memory.Limit = int64(input.MemoryMB)

	dconfig.DatacenterName = v.Session.Datacenter.Name()
	dconfig.ClusterPath = v.Session.Cluster.InventoryPath

	// first element is the cluster or host, then resource pool path
	pElem := filepath.SplitList(input.ComputeResourcePath)
	path := ""
	if len(pElem) > 1 {
		pElem = pElem[1:]
		path = strings.Join(pElem, "/")
	}
	dconfig.ResourcePoolPath = fmt.Sprintf("%s/Resources/%s", dconfig.ClusterPath, path)

	return &dconfig
}
