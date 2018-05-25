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
	"net/url"
	"net/http"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/errors"
)

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