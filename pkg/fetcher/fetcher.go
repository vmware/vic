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

package fetcher

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
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

	"golang.org/x/net/context/ctxhttp"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
)

const (
	maxDownloadAttempts = 5

	// DefaultTokenExpirationDuration specifies the default token expiration
	DefaultTokenExpirationDuration = 60 * time.Second
)

// Fetcher interface
type Fetcher interface {
	Fetch(ctx context.Context, url *url.URL, reqHdrs *http.Header, toFile bool, po progress.Output, id ...string) (string, error)
	FetchAuthToken(url *url.URL) (*Token, error)

	Head(url *url.URL) (http.Header, error)

	ExtractOAuthURL(hdr string, repository *url.URL) (*url.URL, error)

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

// Options struct
type Options struct {
	Timeout time.Duration

	Username string
	Password string

	InsecureSkipVerify bool

	Token *Token

	// RootCAs will not be modified by fetcher.
	RootCAs *x509.CertPool
}

// URLFetcher struct
type URLFetcher struct {
	client *http.Client

	OAuthEndpoint *url.URL

	StatusCode int

	options Options
}

// NewURLFetcher creates a new URLFetcher
func NewURLFetcher(options Options) Fetcher {
	/* #nosec */
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.InsecureSkipVerify,
			RootCAs:            options.RootCAs,
		},
	}
	client := &http.Client{Transport: tr}

	return &URLFetcher{
		client:  client,
		options: options,
	}
}

// Fetch fetches from a url and stores its content in a temporary file.
//	hdrs is optional.
func (u *URLFetcher) Fetch(ctx context.Context, url *url.URL, reqHdrs *http.Header, toFile bool, po progress.Output, ids ...string) (string, error) {
	defer trace.End(trace.Begin(url.String()))

	// extract ID from ids. Existence of an ID enables progress reporting
	ID := ""
	if len(ids) > 0 {
		ID = ids[0]
	}

	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), u.options.Timeout)
	defer cancel()

	var data string
	var err error
	var retries int
	for {
		if toFile {
			data, err = u.fetchToFile(ctx, url, reqHdrs, ID, po)
		} else {
			data, err = u.fetchToString(ctx, url, reqHdrs, ID)
		}
		if err == nil {
			return data, nil
		}

		// If an error was returned because the context was cancelled, we shouldn't retry.
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("download cancelled during download")
		default:
		}

		retries++
		// give up if we reached maxDownloadAttempts
		if retries == maxDownloadAttempts {
			log.Debugf("Hit max download attempts. Download failed: %v", err)
			return "", err
		}

		switch err := err.(type) {
		case DoNotRetry, TagNotFoundError, ImageNotFoundError:
			log.Debugf("Error: %s", err.Error())
			return "", err
		}

		// retry downloading again
		log.Debugf("Download failed, retrying: %v", err)

		delay := retries * 5
		ticker := time.NewTicker(time.Second)

	selectLoop:
		for {
			// Do not report progress back if ID is empty
			if ID != "" && po != nil {
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

func (u *URLFetcher) FetchAuthToken(url *url.URL) (*Token, error) {
	defer trace.End(trace.Begin(url.String()))

	data, err := u.Fetch(context.Background(), url, nil, false, nil)
	if err != nil {
		log.Errorf("Download failed: %v", err)
		return nil, err
	}

	token := &Token{}

	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		log.Errorf("Incorrect token format: %v", err)
		return nil, err
	}

	if token.ExpiresIn == 0 {
		token.Expires = time.Now().Add(DefaultTokenExpirationDuration)
	} else {
		token.Expires = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	}

	return token, nil
}

func (u *URLFetcher) fetch(ctx context.Context, url *url.URL, reqHdrs *http.Header, ID string) (io.ReadCloser, http.Header, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	u.setBasicAuth(req)

	u.setAuthToken(req)

	u.setUserAgent(req)

	// Add optional request headers
	if reqHdrs != nil {
		for k, values := range *reqHdrs {
			for _, v := range values {
				req.Header.Add(k, v)
			}
		}
	}

	res, err := ctxhttp.Do(ctx, u.client, req)
	if err != nil {
		return nil, nil, err
	}

	u.StatusCode = res.StatusCode

	if u.options.Token == nil && u.IsStatusUnauthorized() {
		hdr := res.Header.Get("www-authenticate")
		if hdr == "" {
			return nil, nil, fmt.Errorf("www-authenticate header is missing")
		}
		u.OAuthEndpoint, err = u.ExtractOAuthURL(hdr, url)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, DoNotRetry{Err: fmt.Errorf("Authentication required")}
	}

	if u.IsStatusNotFound() {
		err = fmt.Errorf("Not found: %d, URL: %s", u.StatusCode, url)
		return nil, nil, TagNotFoundError{Err: err}
	}

	if u.IsStatusUnauthorized() {
		hdr := res.Header.Get("www-authenticate")

		// check if image is non-existent (#757)
		if strings.Contains(hdr, "error=\"insufficient_scope\"") {
			err = fmt.Errorf("image not found")
			return nil, nil, ImageNotFoundError{Err: err}
		} else if strings.Contains(hdr, "error=\"invalid_token\"") {
			return nil, nil, fmt.Errorf("not authorized")
		} else {
			return nil, nil, fmt.Errorf("Unexpected http code: %d, URL: %s", u.StatusCode, url)
		}
	}

	// FIXME: handle StatusTemporaryRedirect and StatusFound
	if !u.IsStatusOK() {
		return nil, nil, fmt.Errorf("Unexpected http code: %d, URL: %s", u.StatusCode, url)
	}

	log.Debugf("URLFetcher.fetch() - %#v, %#v", res.Body, res.Header)
	return res.Body, res.Header, nil
}

// fetch fetches the given URL using ctxhttp. It also streams back the progress bar only when ID is not an empty string.
func (u *URLFetcher) fetchToFile(ctx context.Context, url *url.URL, reqHdrs *http.Header, ID string, po progress.Output) (string, error) {
	rdr, hdrs, err := u.fetch(ctx, url, reqHdrs, ID)
	if err != nil {
		return "", err
	}
	defer rdr.Close()

	// stream progress as json and body into a file - only if we have an ID and a Content-Length header
	if contLen := hdrs.Get("Content-Length"); ID != "" && contLen != "" {
		cl, cerr := strconv.ParseInt(contLen, 10, 64)
		if cerr != nil {
			return "", cerr
		}

		if po != nil {
			rdr = progress.NewProgressReader(
				ioutils.NewCancelReadCloser(ctx, rdr), po, cl, ID, "Downloading",
			)
			defer rdr.Close()
		} else {
			rdr = ioutils.NewCancelReadCloser(ctx, rdr)
		}
	}

	// Create a temporary file and stream the res.Body into it
	out, err := ioutil.TempFile(os.TempDir(), ID)
	if err != nil {
		return "", DoNotRetry{Err: err}
	}
	defer out.Close()

	// Stream into it
	_, err = io.Copy(out, rdr)
	if err != nil {
		log.Errorf("Fetch (%s) to file failed to stream to file: %s", url.String(), err)

		// cleanup
		defer os.Remove(out.Name())
		return "", DoNotRetry{Err: err}
	}

	// Return the temporary file name
	return out.Name(), nil
}

// fetch fetches the given URL using ctxhttp. It also streams back the progress bar only when ID is not an empty string.
func (u *URLFetcher) fetchToString(ctx context.Context, url *url.URL, reqHdrs *http.Header, ID string) (string, error) {
	rdr, _, err := u.fetch(ctx, url, reqHdrs, ID)
	if err != nil {
		log.Errorf("Fetch (%s) to string error: %s", url.String(), err)
		return "", err
	}
	defer rdr.Close()

	out := bytes.NewBuffer(nil)

	// Stream into it
	_, err = io.Copy(out, rdr)
	if err != nil {
		// cleanup
		return "", DoNotRetry{Err: err}
	}

	// Return the string
	return string(out.Bytes()), nil
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

func (u *URLFetcher) setUserAgent(req *http.Request) {
	log.Debugf("Setting user-agent to vic/%s", version.Version)
	req.Header.Set("User-Agent", "vic/"+version.Version)
}

func (u *URLFetcher) setBasicAuth(req *http.Request) {
	if u.options.Username != "" && u.options.Password != "" {
		log.Debugf("Setting BasicAuth: %s", u.options.Username)
		req.SetBasicAuth(u.options.Username, u.options.Password)
	}
}

func (u *URLFetcher) setAuthToken(req *http.Request) {
	if u.options.Token != nil {
		req.Header.Set("Authorization", "Bearer "+u.options.Token.Token)
	}
}

// ExtractOAuthURL extracts the OAuth url from the www-authenticate header
func (u *URLFetcher) ExtractOAuthURL(hdr string, repository *url.URL) (*url.URL, error) {
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
