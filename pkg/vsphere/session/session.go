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

package session

import (
	"crypto/tls"
	"time"

	"golang.org/x/net/context"

	"github.com/juju/errors"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25/soap"
)

// Config contains the configuration used to create a Session.
type Config struct {
	// SDK URL or proxy
	Service string
	// Allow insecure connection to Service
	Insecure bool
	// Keep alive duration
	Keepalive time.Duration

	ClusterPath    string
	DatacenterPath string
	DatastorePath  string
	HostPath       string
	NetworkPath    string
	PoolPath       string

	CertFile string
	KeyFile  string
}

// HasCertificate checks for presence of a certificate and keyfile
func (c *Config) HasCertificate() bool {
	return c.CertFile != "" && c.KeyFile != ""
}

// Session caches vSphere objects obtained by querying the SDK.
type Session struct {
	*govmomi.Client

	*Config

	Cluster    *object.ComputeResource
	Datacenter *object.Datacenter
	Datastore  *object.Datastore
	Host       *object.HostSystem
	Network    object.NetworkReference
	Pool       *object.ResourcePool
}

// NewSession creates a new Session struct
func NewSession(config *Config) *Session {
	return &Session{Config: config}
}

// Create accepts a Config and returns a Session with the cached vSphere resources.
func (s *Session) Create(ctx context.Context) (*Session, error) {
	soapURL, err := soap.ParseURL(s.Service)
	if err != nil {
		return nil, errors.Errorf("SDK URL (%s) could not be parsed: %s", s.Service, err)
	}

	// we can't set a keep alive if we log in directly with client creation
	user := soapURL.User
	soapURL.User = nil

	// 1st connect without any userinfo to get the API type
	s.Client, err = govmomi.NewClient(ctx, soapURL, s.Insecure)
	if err != nil {
		return nil, errors.Errorf("Failed to connect to %s: %s", soapURL.String(), err)
	}

	if s.HasCertificate() {
		if !s.Client.IsVC() {
			return nil, errors.Errorf("Certificate based authentication not yet supported with ESXi")
		}

		// load the certificates
		cert, err2 := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err2 != nil {
			return nil, errors.Errorf("Unable to load X509 key pair(%s,%s): %s", s.CertFile, s.KeyFile, err2)
		}

		// create the new client
		s.Client, err = govmomi.NewClientWithCertificate(ctx, soapURL, s.Insecure, cert)
		if err != nil {
			return nil, errors.Errorf("Failed to connect to %s: %s", soapURL.String(), err)
		}
	}

	if s.Keepalive != 0 {
		// now that we've verified everything, enable keepalive
		s.RoundTripper = session.KeepAlive(s.Client.RoundTripper, s.Keepalive)
	}

	// and now that the keepalive is registered we can log in to trigger it
	if !s.HasCertificate() {
		err = s.Client.Login(ctx, user)
	} else {
		err = s.LoginExtensionByCertificate(ctx, user.Username(), "")
	}
	if err != nil {
		return nil, errors.Errorf("Failed to log in to %s: %s", soapURL.String(), err)
	}

	// Populate s
	finder := find.NewFinder(s.Client.Client, true)

	s.Datacenter, err = finder.DatacenterOrDefault(ctx, s.DatacenterPath)
	if err != nil {
		return nil, err
	}
	finder.SetDatacenter(s.Datacenter)

	s.Cluster, err = finder.ComputeResourceOrDefault(ctx, s.ClusterPath)
	if err != nil {
		return nil, err
	}

	s.Datastore, err = finder.DatastoreOrDefault(ctx, s.DatastorePath)
	if err != nil {
		return nil, err
	}

	s.Host, err = finder.HostSystemOrDefault(ctx, s.HostPath)
	if err != nil {
		if _, ok := err.(*find.DefaultMultipleFoundError); !ok || !s.IsVC() {
			return nil, err
		}
	}

	s.Network, err = finder.NetworkOrDefault(ctx, s.NetworkPath)
	if err != nil {
		return nil, err
	}

	s.Pool, err = finder.ResourcePoolOrDefault(ctx, s.PoolPath)
	if err != nil {
		return nil, err
	}

	return s, nil
}
