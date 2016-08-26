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

const (
	maxDownloadAttempts = 5
)

// Fetcher interface
type Fetcher interface {
	Fetch(url *url.URL, id ...string) (string, error)

	Head(url *url.URL) (http.Header, error)

	ExtractOAuthUrl(hdr string, repository *url.URL) (*url.URL, error)

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
	Expires   time.Time
	ExpiresIn int       `json:"expires_in"`
	IssueAt   time.Time `json:"issued_at"`
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

// NewURLFetcher creates a new URLFetcher
func NewURLFetcher(options FetcherOptions) Fetcher {
	/* #nosec */
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

// Fetch fetches from a url and stores its content in a temporary file.
func (u *URLFetcher) Fetch(url *url.URL, ids ...string) (string, error) {
	defer trace.End(trace.Begin(url.String()))

	// extract ID from ids. Existence of an ID enables progress reporting
	ID := ""
	if len(ids) > 0 {
		ID = ids[0]
	}

	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), u.options.Timeout)
	defer cancel()

	var name string
	var err error
	var retries int
	for {
		name, err = u.fetch(ctx, url, ID)
		if err == nil {
			return name, nil
		}

		// If an error was returned because the context was cancelled, we shouldn't retry.
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("download cancelled during download")
		default:
		}

		retries++
		// give up if we reached maxDownloadAttempts or got a DNR
		if _, isDNR := err.(DoNotRetry); isDNR || retries == maxDownloadAttempts {
			log.Debugf("Download failed: %v", err)
			return "", err
		}

		// retry downloading again
		log.Debugf("Download failed, retrying: %v", err)

		delay := retries * 5
		ticker := time.NewTicker(time.Second)

	selectLoop:
		for {
			// Do not report progress back if ID is empty
			if ID != "" {
				progress.Updatef(po, ID, "Retrying in %d second%s", delay, (map[bool]string{true: "s"})[delay != 1])
			}

			select {
			case <-ticker.C:
				delay--
				if delay == 0 {
					ticker.Stop()
					break selectLoop
				}
			case <-ctx.Done():
				ticker.Stop()
				return "", fmt.Errorf("download cancelled during retry delay")
			}
		}
	}
}

// fetch fetches the given URL using ctxhttp. It also streams back the progress bar only when ID is not an empty string.
func (u *URLFetcher) fetch(ctx context.Context, url *url.URL, ID string) (string, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", err
	}

	u.setBasicAuth(req)

	u.setAuthToken(req)

	res, err := ctxhttp.Do(ctx, u.client, req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	u.StatusCode = res.StatusCode

	if u.options.Token == nil && u.IsStatusUnauthorized() {
		hdr := res.Header.Get("www-authenticate")
		if hdr == "" {
			return "", fmt.Errorf("www-authenticate header is missing")
		}
		u.OAuthEndpoint, err = u.ExtractOAuthUrl(hdr, url)
		if err != nil {
			return "", err
		}
		return "", DoNotRetry{Err: fmt.Errorf("Authentication required")}
	}

	if u.IsStatusNotFound() {
		return "", fmt.Errorf("Not found: %d, URL: %s", u.StatusCode, url)
	}

	if u.IsStatusUnauthorized() {
		hdr := res.Header.Get("www-authenticate")

		// check if image is non-existent (#757)
		if strings.Contains(hdr, "error=\"insufficient_scope\"") {
			return "", DoNotRetry{Err: fmt.Errorf("image not found")}
		} else if strings.Contains(hdr, "error=\"invalid_token\"") {
			return "", fmt.Errorf("not authorized")
		} else {
			return "", fmt.Errorf("Unexpected http code: %d, URL: %s", u.StatusCode, url)
		}
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
		return "", DoNotRetry{Err: err}
	}
	defer out.Close()

	// Stream into it
	_, err = io.Copy(out, in)
	if err != nil {
		// cleanup
		defer os.Remove(out.Name())
		return "", DoNotRetry{Err: err}
	}

	// Return the temporary file name
	return out.Name(), nil
}

// Head sends a HEAD request to url
func (u *URLFetcher) Head(url *url.URL) (http.Header, error) {
	defer trace.End(trace.Begin(url.String()))

	ctx, cancel := context.WithTimeout(context.Background(), u.options.Timeout)
	defer cancel()

	return u.head(ctx, url)
}

func (u *URLFetcher) head(ctx context.Context, url *url.URL) (http.Header, error) {
	res, err := ctxhttp.Head(ctx, u.client, url.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	u.StatusCode = res.StatusCode

	if u.IsStatusUnauthorized() || u.IsStatusOK() {
		return res.Header, nil
	}
	return nil, fmt.Errorf("Unexpected http code: %d, URL: %s", u.StatusCode, url)
}

// AuthURL returns the Oauth endpoint URL
func (u *URLFetcher) AuthURL() *url.URL {
	return u.OAuthEndpoint
}

// IsStatusUnauthorized returns true if status code is StatusUnauthorized
func (u *URLFetcher) IsStatusUnauthorized() bool {
	return u.StatusCode == http.StatusUnauthorized
}

// IsStatusOK returns true if status code is StatusOK
func (u *URLFetcher) IsStatusOK() bool {
	return u.StatusCode == http.StatusOK
}

// IsStatusNotFound returns true if status code is StatusNotFound
func (u *URLFetcher) IsStatusNotFound() bool {
	return u.StatusCode == http.StatusNotFound
}

func (u *URLFetcher) setBasicAuth(req *http.Request) {
	if u.options.Username != "" && u.options.Password != "" {
		log.Debugf("Setting BasicAuth: %s", u.options.Username)
		req.SetBasicAuth(u.options.Username, u.options.Password)
	}
}

func (u *URLFetcher) setAuthToken(req *http.Request) {
	if u.options.Token != nil {
		log.Debugf("Setting AuthToken: %s", u.options.Token.Token)
		req.Header.Set("Authorization", "Bearer "+u.options.Token.Token)
	}
}

func (u *URLFetcher) ExtractOAuthUrl(hdr string, repository *url.URL) (*url.URL, error) {
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
