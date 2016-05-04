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

package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"

	"github.com/vmware/vic/pkg/trace"
)

// Fetcher interface
type Fetcher interface {
	Fetch(url *url.URL) (string, error)
	FetchWithProgress(url *url.URL, ID string) (string, error)

	IsStatusUnauthorized() bool
	IsStatusOK() bool
	IsStatusNotFound() bool

	AuthURL() *url.URL
}

// Token represents https://docs.docker.com/registry/spec/auth/token/
type Token struct {
	// An opaque Bearer token that clients should supply to subsequent requests in the Authorization header.
	Token string `json:"token"`
	// (Optional) The duration in seconds since the token was issued that it will remain valid. When omitted, this defaults to 60 seconds.
	Expires time.Time `json:"expires_in"`
}

// FetcherOptions struct
type FetcherOptions struct {
	Timeout time.Duration

	Username string
	Password string

	InsecureSkipVerify bool

	Token *Token
}

// URLFetcher struct
type URLFetcher struct {
	client *http.Client

	OAuthEndpoint *url.URL

	StatusCode int

	options FetcherOptions
}

// NewFetcher creates a new Fetcher instance
func NewFetcher(options FetcherOptions) Fetcher {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.InsecureSkipVerify,
		},
	}
	client := &http.Client{Transport: tr}

	return &URLFetcher{
		client:  client,
		options: options,
	}
}

// Fetch fetches a web page from url and stores in a temporary file.
func (u *URLFetcher) Fetch(url *url.URL) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.options.Timeout)
	defer cancel()

	return u.fetch(ctx, url, "")
}

// FetchWithProgress fetches a web page from url and stores in a temporary file while showing a progress bar.
func (u *URLFetcher) FetchWithProgress(url *url.URL, ID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.options.Timeout)
	defer cancel()

	return u.fetch(ctx, url, ID)
}

func (u *URLFetcher) fetch(ctx context.Context, url *url.URL, ID string) (string, error) {
	defer trace.End(trace.Begin(url.String()))

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", err
	}

	u.SetBasicAuth(req)

	u.SetAuthToken(req)

	res, err := ctxhttp.Do(ctx, u.client, req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	u.StatusCode = res.StatusCode

	if u.IsStatusUnauthorized() {
		hdr := res.Header.Get("www-authenticate")
		if hdr == "" {
			return "", fmt.Errorf("www-authenticate header is missing")
		}
		u.OAuthEndpoint, err = u.ExtractQueryParams(hdr, url)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("Authentication required")
	}

	// FIXME: handle StatusTemporaryRedirect and StatusFound
	if !u.IsStatusOK() {
		return "", fmt.Errorf("Unexpected http code: %d, URL: %s", u.StatusCode, url)
	}

	in := res.Body
	// stream progress as json and body into a file - only if we have an ID and a Content-Length header
	if hdr := res.Header.Get("Content-Length"); ID != "" && hdr != "" {
		cl, cerr := strconv.ParseInt(hdr, 10, 64)
		if cerr != nil {
			return "", cerr
		}

		in = progress.NewProgressReader(
			ioutils.NewCancelReadCloser(ctx, res.Body), po, cl, ID, "Downloading",
		)
		defer in.Close()
	}

	// Create a temporary file and stream the res.Body into it
	out, err := ioutil.TempFile(os.TempDir(), ID)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Stream into it
	_, err = io.Copy(out, in)
	if err != nil {
		return "", err
	}

	// Return the temporary file name
	return out.Name(), nil
}

func (u *URLFetcher) AuthURL() *url.URL {
	return u.OAuthEndpoint
}

func (u *URLFetcher) IsStatusUnauthorized() bool {
	return u.StatusCode == http.StatusUnauthorized
}

func (u *URLFetcher) IsStatusOK() bool {
	return u.StatusCode == http.StatusOK
}

func (u *URLFetcher) IsStatusNotFound() bool {
	return u.StatusCode == http.StatusNotFound
}

func (u *URLFetcher) SetBasicAuth(req *http.Request) {
	if u.options.Username != "" && u.options.Password != "" {
		log.Debugf("Setting BasicAuth: %s", u.options.Username)
		req.SetBasicAuth(u.options.Username, u.options.Password)
	}
}

func (u *URLFetcher) SetAuthToken(req *http.Request) {
	if u.options.Token != nil {
		log.Debugf("Setting AuthToken: %s", u.options.Token.Token)
		req.Header.Set("Authorization", "Bearer "+u.options.Token.Token)
	}
}

func (u *URLFetcher) ExtractQueryParams(hdr string, repository *url.URL) (*url.URL, error) {
	tokens := strings.Split(hdr, " ")
	if len(tokens) != 2 || strings.ToLower(tokens[0]) != "bearer" {
		return nil, fmt.Errorf("www-authenticate header is corrupted")
	}

	tokens = strings.Split(tokens[1], ",")

	var realm, service, scope string
	for _, token := range tokens {
		if strings.HasPrefix(token, "realm") {
			realm = strings.Trim(token[len("realm="):], "\"")
		}
		if strings.HasPrefix(token, "service") {
			service = strings.Trim(token[len("service="):], "\"")
		}
		if strings.HasPrefix(token, "scope") {
			scope = strings.Trim(token[len("scope="):], "\"")
		}
	}

	if realm == "" {
		return nil, fmt.Errorf("missing realm in bearer auth challenge")
	}
	if service == "" {
		return nil, fmt.Errorf("missing service in bearer auth challenge")
	}
	// The scope can be empty if we're not getting a token for a specific repo
	if scope == "" && repository != nil {
		return nil, fmt.Errorf("missing scope in bearer auth challenge")
	}

	auth, err := url.Parse(realm)
	if err != nil {
		return nil, err
	}

	q := auth.Query()
	q.Add("service", service)
	if scope != "" {
		q.Add("scope", scope)
	}
	auth.RawQuery = q.Encode()

	return auth, nil
}
