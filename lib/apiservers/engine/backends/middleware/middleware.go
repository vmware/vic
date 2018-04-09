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

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	vicdns "github.com/vmware/vic/lib/dns"

	context "golang.org/x/net/context"
)

// HostCheckMiddleware provides middleware for Host: field checking
type HostCheckMiddleware struct {
	ValidDomains vicdns.FQDNs
}

// WrapHandler satisfies the Docker middleware interface for HostCheckMiddleware
func (h HostCheckMiddleware) WrapHandler(f func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {

		hostname := strings.Split(r.Host, ":")[0] // trim port if it's there. r.Host should never contain a scheme so this should be fine
		if h.ValidDomains[hostname] {
			return f(ctx, w, r, vars)
		}

		return fmt.Errorf("request from %s with HTTP Host header \"%s\" rejected because %s is an invalid hostname for this endpoint", r.RemoteAddr, r.Host, r.Host)
	}
}
