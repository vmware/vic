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
	"fmt"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/registry"
)

type merger struct {
}

func NewMerger() Merger {
	return &merger{}
}

// Merge merges two config objects together. For now only
// whitelist registries are merged.
func (m *merger) Merge(orig, other *config.VirtualContainerHostConfigSpec) (*config.VirtualContainerHostConfigSpec, error) {
	otherWl, err := ParseRegistries(other.RegistryWhitelist)
	if err != nil {
		return nil, err
	}

	origWl, err := ParseRegistries(orig.RegistryWhitelist)
	if err != nil {
		return nil, err
	}

	var wl registry.Set
	if len(otherWl) > 0 {
		if wl, err = origWl.Merge(otherWl, &whitelistMerger{}); err != nil {
			return nil, err
		}
	}

	// whitelist should not grow
	if len(wl) > len(origWl) {
		return nil, fmt.Errorf("whitelist cannot grow")
	}

	res := *orig
	res.RegistryWhitelist = wl.Strings()
	return &res, nil
}

func ParseRegistries(regs []string) (registry.Set, error) {
	var s registry.Set
	for _, r := range regs {
		e := registry.ParseEntry(r)
		if e != nil {
			s = append(s, e)
			continue
		}

		return nil, fmt.Errorf("could not parse entry %s", r)
	}

	return s, nil
}

type whitelistMerger struct{}

// Merge merges two registry entries. The merge fails if merging orig and other would
// broaden orig's scope. The result of the merge is other if that is more restrictive.
// if orig equals other, the result is orig.
func (w *whitelistMerger) Merge(orig, other registry.Entry) (registry.Entry, error) {
	if orig.Equal(other) {
		return orig, nil
	}

	if other.Contains(orig) {
		return nil, fmt.Errorf("merge of %s and %s would broaden %s", orig, other, orig)
	}

	// more restrictive result is OK
	if orig.Contains(other) {
		return other, nil
	}

	// no merge
	return nil, nil
}
