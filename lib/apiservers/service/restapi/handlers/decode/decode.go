// Copyright 2018 VMware, Inc. All Rights Reserved.
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

// Package decode converts model objects used by the API into objects used by
// the management package. Functions are grouped by area.
package decode

import (
	"net/http"
	"net/url"

	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/client"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/errors"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
)

func ProcessVCH(op trace.Operation, d *data.Data, vch *models.VCH, finder client.Finder) error {
	var err error

	if vch != nil {
		err = checkVersion(vch.Version)
		if err != nil {
			return err
		}

		err := processVCHName(op, d, vch)
		if err != nil {
			return err
		}

		debug := int(vch.Debug)
		d.Debug.Debug = &debug

		err = ProcessCompute(op, d, vch, finder)
		if err != nil {
			return err
		}

		err = ProcessNetworks(op, d, vch, finder)
		if err != nil {
			return err
		}

		err = ProcessStorage(op, d, vch, finder)
		if err != nil {
			return err
		}

		err = ProcessCertificates(op, d, vch)
		if err != nil {
			return err
		}

		err = ProcessEndpoint(op, d, vch)
		if err != nil {
			return err
		}

		err = ProcessRegistry(op, d, vch)
		if err != nil {
			return err
		}

		err = ProcessSyslogs(op, d, vch)
		if err != nil {
			return err
		}

		err = ProcessContainer(op, d, vch)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO [AngieCris]: processing of all other fields (ones don't belong to other categories) should be moved to some other file

func checkVersion(ver models.Version) error {
	if ver != "" && version.String() != string(ver) {
		return errors.NewError(http.StatusBadRequest, "invalid version: %s", string(ver))
	}

	return nil
}

// TODO [AngieCris] (#6710): move validation to swagger
func processVCHName(op trace.Operation, d *data.Data, vch *models.VCH) error {
	d.DisplayName = vch.Name
	err := CheckUnsupportedChars(d.DisplayName)

	if err != nil {
		return errors.NewError(http.StatusBadRequest, "invalid display name: %s", err)
	}

	if len(d.DisplayName) > MaxDisplayNameLen {
		return errors.NewError(http.StatusBadRequest, "invalid display name: length exceeds %d characters", MaxDisplayNameLen)
	}

	return nil
}

func ProcessSyslogs(op trace.Operation, d *data.Data, vch *models.VCH) error {
	if vch.SyslogAddr != "" {
		addr := vch.SyslogAddr.String()
		if len(addr) == 0 {
			return nil
		}

		u, err := url.Parse(addr)
		if err != nil {
			return errors.WrapError(http.StatusBadRequest, err)
		}

		d.SyslogConfig.Addr = u
	}

	return nil
}

func ProcessContainer(op trace.Operation, d *data.Data, vch *models.VCH) error {
	if vch.Container != nil && vch.Container.NameConvention != "" {
		d.ContainerNameConvention = vch.Container.NameConvention
	}

	return nil
}
