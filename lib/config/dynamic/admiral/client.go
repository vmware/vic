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
	"net/url"
	"sync"

	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/lib/config/dynamic"
	admclient "github.com/vmware/vic/lib/config/dynamic/admiral/client"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tags"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

type config struct {
	Registries []string
}

type client interface {
	Config(context.Context) (*config, error)
}

func newClient() client {
	return &defaultClient{}
}

const (
	ovaTagCategory     = "VsphereIntegratedContainers"
	ovaTagName         = "ProductVM"
	ovaTokenKey        = "some/key" // as yet undetermined
	admiralEndpointKey = "some/key" // undetermined
)

var projectsFilter = "customProperties.enableContentTrust eq 'true'"

type defaultClient struct {
	mu       sync.Mutex
	c        *tags.RestClient
	adm      *admclient.ContainerManagement
	loggedIn bool
	sess     *session.Session
	vm       *vm.VirtualMachine
	token    string
}

func (c *defaultClient) Config(ctx context.Context) (*config, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// try a couple of times
	for i := 0; i < 2; i++ {
		if err := c.init(ctx); err != nil {
			return nil, err
		}

		// TODO: hash this out when info about admiral is in guestinfo of ova vm
		//
		// params := projects.NewGetProjectsParamsWithContext(ctx).WithDollarFilter(&projectsFilter)
		// var res *projects.GetProjectsOK
		// res, err = c.adm.Projects.GetProjects(params)
		// if err != nil {
		// 	c.reset()
		// 	continue
		// }

		// // do something with res
		// for _, r := range res.Payload.DocumentLinks {
		// }
	}

	// for now just return an error
	return nil, dynamic.ErrSourceUnavailable
}

func (c *defaultClient) reset() {
	c.vm = nil
}

func (c *defaultClient) init(ctx context.Context) error {
	var err error
	if c.sess == nil {
		if c.sess, err = session.NewSession(&session.Config{}).Create(ctx); err != nil {
			return fmt.Errorf("could not create session: %s", err)
		}

		u, err := url.Parse(c.sess.Service)
		if err != nil {
			return err
		}

		c.c = tags.NewClient(u, c.sess.Insecure)
	}

	if err = c.populate(ctx); err != nil {
		return err
	}

	return nil
}

func (c *defaultClient) login() error {
	if c.loggedIn {
		return nil
	}

	if err := c.c.Login(); err != nil {
		return err
	}

	c.loggedIn = true
	return nil
}

func (c *defaultClient) findOVATag() (string, error) {
	cats, err := c.c.GetCategoriesByName(ovaTagCategory)
	if err != nil {
		return "", err
	}

	// just use the first one
	if len(cats) == 0 {
		return "", errors.New("could not find tag")
	}

	cat := cats[0]
	tags, err := c.c.GetTagByNameForCategory(ovaTagName, cat.ID)
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", errors.New("could not find tag")
	}

	return tags[0].ID, nil
}

func (c *defaultClient) populate(ctx context.Context) error {
	if c.vm != nil {
		return nil
	}

	if err := c.login(); err != nil {
		return fmt.Errorf("error logging in: %s", err)
	}

	t, err := c.findOVATag()
	if err != nil {
		return fmt.Errorf("could not find ova tag: %s", err)
	}

	objs, err := c.c.ListAttachedObjects(t)
	if err != nil || len(objs) == 0 {
		return fmt.Errorf("could not find ova vm")
	}

	for _, o := range objs {
		if o.Type != nil && *o.Type != "VirtualMachine" {
			continue
		}

		v := vm.NewVirtualMachine(ctx, c.sess, types.ManagedObjectReference{Type: *o.Type, Value: *o.ID})
		values, err := keys(ctx, v, []string{ovaTokenKey, admiralEndpointKey})
		if err != nil {
			continue
		}

		c.token = values[ovaTokenKey]
		u, err := url.Parse(values[admiralEndpointKey])
		if err != nil {
			return err
		}

		c.adm = admclient.NewHTTPClientWithConfig(
			nil,
			&admclient.TransportConfig{
				Host:     u.Host,
				BasePath: u.Path,
				Schemes:  []string{u.Scheme},
			})

		c.vm = v
		break
	}

	if c.vm == nil {
		return fmt.Errorf("could not find ova vm")
	}

	return nil
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
			if k == ov.GetOptionValue().Key {
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
