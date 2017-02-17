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

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
)

// MigrateConfig migrate old VCH configuration to new version. Currently check required fields only
func (v *Validator) ValidateMigratedConfig(ctx context.Context, conf *config.VirtualContainerHostConfigSpec) (*config.VirtualContainerHostConfigSpec, error) {
	defer trace.End(trace.Begin(conf.Name))

	v.assertBasics(conf)
	v.assertTarget(conf)
	v.assertDatastore(conf)
	v.assertNetwork(conf)

	return conf, v.ListIssues()
}

func (v *Validator) assertNetwork(conf *config.VirtualContainerHostConfigSpec) {
	// minimum network configuration check
}

// assertDatastore check required datastore configuration only
func (v *Validator) assertDatastore(conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))
	if len(conf.ImageStores) == 0 {
		v.NoteIssue(errors.New("Image store is not set"))
	}
}

func (v *Validator) assertTarget(conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	if conf.Target == "" {
		v.NoteIssue(errors.New("target is not set"))
	}

	if conf.Username == "" {
		v.NoteIssue(errors.New("target username is not set"))
	}

	if conf.Token == "" {
		v.NoteIssue(errors.New("target token is not set"))
	}
}

func (v *Validator) assertBasics(conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))
	v.assertVersion(conf)
}

func (v *Validator) assertVersion(conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))
	if conf.Version == nil {
		v.NoteIssue(errors.Errorf("Unknown version of VCH %q", conf.Name))
		return
	}
	installerBuild := version.GetBuild()
	if installerBuild.Equal(conf.Version) {
		v.NoteIssue(errors.Errorf("%q has same version as installer. No upgrade is available.", conf.Name))
		return
	}
	older, err := installerBuild.IsOlder(conf.Version)
	if err != nil {
		v.NoteIssue(errors.Errorf("Failed to compare VCH version %q with installer version %q: %s", conf.Version.ShortVersion(), installerBuild.ShortVersion(), err))
		return
	}
	if older {
		v.NoteIssue(errors.Errorf("VCH version %q is newer than installer version %q", conf.Version.ShortVersion(), installerBuild.ShortVersion()))
		return
	}
}
