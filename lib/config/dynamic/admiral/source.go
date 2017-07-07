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
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime"
	rtclient "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/vim25/types"

	"strings"

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

	clusterFilter = "(address eq '%s' and customProperties.__containerHostType eq 'VCH')"
)

var (
	trueStr        = "true"
	projectsFilter = "customProperties.__enableContentTrust eq 'true'"
)

// NewSource creates a new Admiral dynamic config source. sess
// is a valid vsphere session object. vchID is a unique identifier
// that will be used to lookup the VCH in the Admiral instance; currently
// this is a URI to the docker endpoint in the VCH.
func NewSource(sess *session.Session, vchID string) dynamic.Source {
	return &source{
		d:     &productDiscovery{},
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
		return nil, err
	}

	var projs []string
	projs, err = a.projects(ctx)
	if err != nil {
		return nil, err
	}

	var wl []string
	wl, err = a.whitelist(ctx, projs)
	if err != nil {
		return nil, err
	}

	return &vchcfg.VirtualContainerHostConfigSpec{
		Registry: vchcfg.Registry{RegistryWhitelist: wl},
	}, nil
}

func (a *source) discover(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.c != nil {
		return nil
	}

	token, u, err := a.d.Discover(ctx, a.sess)
	if err != nil {
		return err
	}

	// copied from http.DefaultTransport
	tp := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	cl := &http.Client{
		Transport: tp,
	}

	if u.Scheme == "https" {
		tp.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	rt := rtclient.NewWithClient(u.Host, u.Path, []string{u.Scheme}, cl)
	rt.DefaultAuthentication = &admiralAuth{token: token}
	a.c = client.New(rt, strfmt.Default)

	return nil
}

type admiralAuth struct {
	token string
}

func (a *admiralAuth) AuthenticateRequest(req runtime.ClientRequest, _ strfmt.Registry) error {
	return req.SetHeaderParam("x-xenon-auth-token", a.token)
}

func (a *source) projects(ctx context.Context) ([]string, error) {
	ids := []string{a.vchID}
	if u, err := url.Parse(a.vchID); err == nil {
		if u.Scheme == "" {
			ids = append(ids, "https://"+a.vchID)
		} else {
			ids = append(ids, strings.TrimPrefix(a.vchID, u.Scheme))
		}
	}

	var err error
	var comps *resources_compute.GetResourcesComputeOK
	for _, vchID := range ids {
		filter := fmt.Sprintf(clusterFilter, vchID)
		log.Debugf("getting compute resources with filter %s", filter)
		comps, err = a.c.ResourcesCompute.GetResourcesCompute(resources_compute.NewGetResourcesComputeParamsWithContext(ctx).WithDollarFilter(&filter))
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if comps.Payload.DocumentCount == 0 {
		return nil, errors.Errorf("no admiral instances have host %s registered", a.vchID)
	}

	comp := &models.ComVmwarePhotonControllerModelResourcesComputeServiceComputeState{}
	if err := mapstructure.Decode(comps.Payload.Documents[comps.Payload.DocumentLinks[0]], comp); err != nil {
		return nil, err
	}

	return comp.TenantLinks, nil
}

func (a *source) whitelist(ctx context.Context, hostProjs []string) ([]string, error) {
	// find at least one project with enable content trust
	// that also contains the vch
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
		// no project with enable content trust and vch
		return nil, nil
	}

	regs, err := a.c.ConfigRegistries.GetConfigRegistries(config_registries.NewGetConfigRegistriesParamsWithContext(ctx).WithExpand(&trueStr))
	if err != nil {
		return nil, nil
	}

	var wl []string
	for _, r := range regs.Payload.Documents {
		m := &models.ComVmwareAdmiralServiceCommonRegistryServiceRegistryState{}
		if err := mapstructure.Decode(r, m); err != nil {
			log.Warnf("skipping registry: %s", err)
			continue
		}

		wl = append(wl, m.Address)
	}

	return wl, nil
}

type discovery interface {
	Discover(ctx context.Context, sess *session.Session) (token string, u *url.URL, err error)
}

type productDiscovery struct {
}

func (o *productDiscovery) Discover(ctx context.Context, sess *session.Session) (token string, u *url.URL, err error) {
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
		err = errors.Errorf("could not find ova tag: %s", err)
		return
	}

	objs, err := t.ListAttachedObjects(ctx, tag)
	if err != nil || len(objs) == 0 {
		err = errors.Errorf("could not find ova vm: %s", err)
		return
	}

	for _, o := range objs {
		if o.Type != nil && *o.Type != "VirtualMachine" {
			// not a virtual machine
			continue
		}

		log.Debugf("%v", o)
		v := vm.NewVirtualMachine(ctx, sess, types.ManagedObjectReference{Type: *o.Type, Value: *o.ID})
		var values map[string]string
		values, err = keys(ctx, v, []string{admiralEndpointKey, admiralTokenKey})
		if err != nil {
			log.Debugf("keys not found in %q: %s", v, err)
			err = nil // keys not found
			continue
		}

		token = values[admiralTokenKey]
		if token == "" {
			// not a useable product installation
			continue
		}

		u, err = url.Parse(values[admiralEndpointKey])
		if err != nil {
			log.Warnf("ignoring bad admiral endpoint %s: %s", values[admiralEndpointKey], err)
			err = nil // ignore bad endpoint
			continue
		}
	}

	if u == nil || token == "" {
		err = errors.Errorf("could not find admiral")
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
			return nil, errors.Errorf("key not found: %s", k)
		}
	}

	return res, nil
}
