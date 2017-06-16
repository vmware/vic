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

package dynamic

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/vsphere/tags"
)

func NewAdmiralSource(target *url.URL, insecure bool) (Source, error) {
	return &admiralSource{
		c: tags.NewClient(target, insecure),
	}, nil
}

type admiralSource struct {
	mu       sync.Mutex
	c        *tags.RestClient
	loggedIn bool
	lastCfg  *config.VirtualContainerHostConfigSpec
}

// Get returns the dynamic config portion from an Admiral instance. For now,
// this is empty pending details from the Admiral team.
func (a *admiralSource) Get(context.Context) (*config.VirtualContainerHostConfigSpec, error) {
	if err := a.login(); err != nil {
		return nil, ErrSourceUnavailable
	}

	tags, _ := a.c.ListTagsForCategory("OVA")
	fmt.Println("tags: %v", tags)

	objs, err := a.c.ListAttachedObjects("OVA")
	if err != nil {
		return nil, ErrSourceUnavailable
	}

	for _, o := range objs {
		fmt.Println("o.Id: %v, o.Type: %v", o.Id, o.Type)
	}

	return nil, nil
}

func (a *admiralSource) login() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.loggedIn {
		return nil
	}

	if err := a.c.Login(); err != nil {
		return err
	}

	a.loggedIn = true
	return nil
}
