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

package admiral

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime"
	"github.com/vmware/govmomi/vim25/types"

	vchcfg "github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/dynamic"
	"github.com/vmware/vic/lib/config/dynamic/admiral/client"
	"github.com/vmware/vic/lib/config/dynamic/admiral/client/config_registries"
	"github.com/vmware/vic/lib/config/dynamic/admiral/client/projects"
	"github.com/vmware/vic/lib/config/dynamic/admiral/client/resources_compute"
	"github.com/vmware/vic/lib/config/dynamic/admiral/models"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tags"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	VicProductCategory = "VsphereIntegratedContainers"
	ProductVMTag       = "ProductVM"
	admiralTokenKey    = "guestinfo.vicova.admiral.token"
	admiralEndpointKey = "guestinfo.vicova.admiral.endpoint"

	clusterFilter = "(address eq %s and customProperties.__containerHostType eq 'VCH'"
)

var trueStr = "true"
var projectsFilter = "customProperties.__enableContentTrust eq 'true'"
var errHostNotFound = errors.New("host not found")

func NewSource(sess *session.Session, vchID string) dynamic.Source {
	return &source{
		d:     &ovaDiscovery{},
		sess:  sess,
		vchID: vchID,
	}
}

type source struct {
	mu    sync.Mutex
	d     discovery
	sess  *session.Session
	vchID string
	c     *client.Admiral
}

// Get returns the dynamic config portion from an Admiral instance. For now,
// this is empty pending details from the Admiral team.
func (a *source) Get(ctx context.Context) (*vchcfg.VirtualContainerHostConfigSpec, error) {
	var err error
	if err = a.discover(ctx); err != nil {
		log.Debugf(err.Error())
		return nil, transformErr(err)
	}

	var projs []string
	projs, err = a.projects(ctx)
	if err != nil {
		return nil, transformErr(err)
	}

	var wl []string
	wl, err = a.whitelist(ctx, projs)
	if err != nil {
		return nil, transformErr(err)
	}

	return &vchcfg.VirtualContainerHostConfigSpec{
		Registry: vchcfg.Registry{RegistryWhitelist: wl},
	}, nil
}

func transformErr(err interface{}) error {
	switch err := err.(type) {
	case runtime.APIError:
		switch err.Code {
		case http.StatusForbidden, http.StatusUnauthorized:
			return dynamic.ErrAccessDenied
		}
	}

	return dynamic.ErrSourceUnavailable
}

func (a *source) discover(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.c != nil {
		return nil
	}
	var u *url.URL
	_, u, err := a.d.Discover(ctx, a.sess)
	if err != nil {
		return err
	}

	a.c = client.NewHTTPClientWithConfig(
		nil,
		&client.TransportConfig{
			Host:     u.Host,
			BasePath: u.Path,
			Schemes:  []string{u.Scheme},
		})

	return nil
}

func (a *source) projects(ctx context.Context) ([]string, error) {
	filter := fmt.Sprintf(clusterFilter, a.vchID)
	comps, err := a.c.ResourcesCompute.GetResourcesCompute(resources_compute.NewGetResourcesComputeParamsWithContext(ctx).WithExpand(&trueStr).WithDollarFilter(&filter))
	if err != nil {
		return nil, err
	}

	if comps.Payload.DocumentCount != 1 {
		return nil, errHostNotFound
	}

	comp := &models.ComVmwarePhotonControllerModelResourcesComputeServiceComputeState{}
	comp.UnmarshalBinary([]byte(comps.Payload.Documents[comps.Payload.DocumentLinks[0]]))
	return comp.TenantLinks, nil
}

func (a *source) whitelist(ctx context.Context, hostProjs []string) ([]string, error) {
	projs, err := a.c.Projects.GetProjects(projects.NewGetProjectsParamsWithContext(ctx).WithDollarFilter(&projectsFilter))
	if err != nil {
		return nil, err
	}

	trust := false
	for _, t := range hostProjs {
		for _, p := range projs.Payload.DocumentLinks {
			if t == p {
				trust = true
				break
			}
		}

		if trust {
			break
		}
	}

	if !trust {
		return nil, nil
	}

	regs, err := a.c.ConfigRegistries.GetConfigRegistries(config_registries.NewGetConfigRegistriesParamsWithContext(ctx).WithExpand(&trueStr))
	if err != nil {
		return nil, nil
	}

	wl := make([]string, regs.Payload.DocumentCount)
	i := 0
	for _, r := range regs.Payload.Documents {
		m := &models.ComVmwareAdmiralServiceCommonRegistryServiceRegistryState{}
		m.UnmarshalBinary([]byte(r))
		wl[i] = m.Address
		i++
	}

	return wl, nil
}

type discovery interface {
	Discover(ctx context.Context, sess *session.Session) (token string, u *url.URL, err error)
}

type ovaDiscovery struct {
}

func (o *ovaDiscovery) Discover(ctx context.Context, sess *session.Session) (token string, u *url.URL, err error) {
	service, err := url.Parse(sess.Service)
	if err != nil {
		return
	}

	service.User = sess.User
	t := tags.NewClient(service, sess.Insecure, sess.Thumbprint)
	if err = t.Login(ctx); err != nil {
		return
	}

	var tag string
	tag, err = findOVATag(ctx, t)
	if err != nil {
		err = fmt.Errorf("could not find ova tag: %s", err)
		return
	}

	objs, err := t.ListAttachedObjects(ctx, tag)
	if err != nil || len(objs) == 0 {
		err = fmt.Errorf("could not find ova vm: %s", err)
		return
	}

	for _, o := range objs {
		if o.Type != nil && *o.Type != "VirtualMachine" {
			// not a virtual machine
			continue
		}

		v := vm.NewVirtualMachine(ctx, sess, types.ManagedObjectReference{Type: *o.Type, Value: *o.ID})
		var values map[string]string
		values, err = keys(ctx, v, []string{admiralEndpointKey})
		if err != nil {
			log.Debugf("keys not found in %q: %s", v, err)
			err = nil // keys not found
			continue
		}

		token = values[admiralTokenKey]
		u, err = url.Parse(values[admiralEndpointKey])
		if err != nil {
			log.Warnf("ignoring bad admiral endpoint %s: %s", values[admiralEndpointKey], err)
			err = nil // ignore bad endpoint
			continue
		}

		return
	}

	if u == nil {
		err = fmt.Errorf("could not find admiral")
		log.Debugf(err.Error())
	}

	return
}

func findOVATag(ctx context.Context, t *tags.RestClient) (string, error) {
	cats, err := t.GetCategoriesByName(ctx, VicProductCategory)
	if err != nil {
		return "", err
	}

	// just use the first one
	if len(cats) == 0 {
		return "", errors.New("could not find tag")
	}

	cat := cats[0]
	tags, err := t.GetTagByNameForCategory(ctx, ProductVMTag, cat.ID)
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", errors.New("could not find tag")
	}

	return tags[0].ID, nil
}

func keys(ctx context.Context, v *vm.VirtualMachine, keys []string) (map[string]string, error) {
	ovs, err := v.FetchExtraConfigBaseOptions(ctx)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for _, k := range keys {
		found := false
		for _, ov := range ovs {
			log.Debugf("key: %s", ov.GetOptionValue().Key)
			if k == ov.GetOptionValue().Key {
				log.Debugf("found %s", ov.GetOptionValue().Key)
				res[k] = ov.GetOptionValue().Value.(string)
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("key not found: %s", k)
		}
	}

	return res, nil
}
