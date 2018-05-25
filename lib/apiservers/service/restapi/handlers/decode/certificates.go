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

package decode

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/errors"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/trace"
)

// TODO [AngieCris]: needs common.Networks and common.CertFactory. Not good enough (decode is only responsible for translating from model to data)
// it needs common.Networks and certFactory to do certificate generating in certs.ProcessCertificate
// to avoid duplicate code, here it uses cert processing code from cmd/common
// need to find a better pattern
func ProcessCertificates(op trace.Operation, d *data.Data, vch *models.VCH) error {
	if vch.Auth != nil {
		// TODO [AngieCris]: use package in cmd/common to avoid duplicate logic and code
		certs := common.CertFactory{}
		certs.NoTLS = vch.Auth.NoTLS

		if vch.Auth.Client != nil {
			certs.NoTLSverify = vch.Auth.Client.NoTLSVerify
			certs.ClientCAs = FromPemCertificates(vch.Auth.Client.CertificateAuthorities)
			d.ClientCAs = certs.ClientCAs
		}

		if vch.Auth.Server != nil {
			if vch.Auth.Server.Generate != nil {
				certs.Cname = vch.Auth.Server.Generate.Cname
				certs.Org = vch.Auth.Server.Generate.Organization
				certs.KeySize = FromValueBits(vch.Auth.Server.Generate.Size)
				certs.NoSaveToDisk = true
				certs.Networks = buildNetwork(d) // TODO [AngieCris]: figure out a plan so we don't have to rely on common.Network to generate cert

				// TODO [AngieCris]: VCH API does not set force flag (?)
				if err := certs.ProcessCertificates(op, d.DisplayName, d.Force, 0); err != nil {
					return errors.NewError(http.StatusBadRequest, "error generating certificates: %s", err)
				}
			} else {
				certs.CertPEM = []byte(vch.Auth.Server.Certificate.Pem)
				certs.KeyPEM = []byte(vch.Auth.Server.PrivateKey.Pem)
			}

			d.KeyPEM = certs.KeyPEM
			d.CertPEM = certs.CertPEM
			d.ClientCAs = certs.ClientCAs
		}
	}

	return nil
}

func ProcessRegistry(op trace.Operation, d *data.Data, vch *models.VCH) error {
	if vch.Registry != nil {
		d.InsecureRegistries = vch.Registry.Insecure
		d.WhitelistRegistries = vch.Registry.Whitelist
		d.RegistryCAs = FromPemCertificates(vch.Registry.CertificateAuthorities)

		if vch.Registry.ImageFetchProxy != nil {
			httpProxy := string(vch.Registry.ImageFetchProxy.HTTP)
			httpsProxy := string(vch.Registry.ImageFetchProxy.HTTPS)
			hpurl, err := processProxy(&httpProxy, "http")
			if err != nil {
				return errors.WrapError(http.StatusBadRequest, err)
			}
			d.HTTPProxy = hpurl

			hpsurl, err := processProxy(&httpsProxy, "https")
			if err != nil {
				return errors.WrapError(http.StatusBadRequest, err)
			}
			d.HTTPSProxy = hpsurl
		}
	}

	return nil
}

func processProxy(proxy *string, scheme string) (*url.URL, error) {
	p, err := url.Parse(*proxy)
	if err != nil || p.Host == "" || p.Scheme != scheme {
		return nil, fmt.Errorf("could not parse %s proxy - expected format: %s://fqdn_or_ip:port: %s", scheme, scheme, err)
	}

	return p, nil
}

func buildNetwork(d *data.Data) common.Networks {
	return common.Networks{
		ClientNetworkName:        d.ClientNetwork.Name,
		ClientNetworkIP:          d.ClientNetwork.IP.String(),
		ClientNetworkGateway:     d.ClientNetwork.Gateway.String(),
		PublicNetworkName:        d.PublicNetwork.Name,
		PublicNetworkIP:          d.PublicNetwork.IP.String(),
		PublicNetworkGateway:     d.PublicNetwork.Gateway.String(),
		ManagementNetworkName:    d.ManagementNetwork.Name,
		ManagementNetworkIP:      d.ManagementNetwork.IP.String(),
		ManagementNetworkGateway: d.ManagementNetwork.Gateway.String(),
	}
}

func FromPemCertificates(m []*models.X509Data) []byte {
	var b []byte

	for _, ca := range m {
		c := []byte(ca.Pem)
		b = append(b, c...)
	}

	return b
}
