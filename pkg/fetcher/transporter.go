// Copyright 2017 VMware, Inc. All Rights Reserved.
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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/progress"
	"golang.org/x/net/context/ctxhttp"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
)

const maxTransportAttempts = 5

// Transporter interface
type Transporter interface {
	Put(ctx context.Context, url *url.URL, body io.Reader, reqHdrs *http.Header, po progress.Output, ids ...string) (http.Header, error)
	Post(ctx context.Context, url *url.URL, body io.Reader, reqHdrs *http.Header, po progress.Output, ids ...string) (http.Header, error)
	Delete(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, error)
	Head(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, error)
	Get(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, io.ReadCloser, error)
	GetBytes(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, []byte, error)
	GetHeaderOnly(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, error)

	IsStatusUnauthorized() bool
	IsStatusOK() bool
	IsStatusNotFound() bool
	IsStatusAccepted() bool
	IsStatusCreated() bool
	IsStatusNoContent() bool
	IsStatusBadGateway() bool
	IsStatusServiceUnavailable() bool
	IsStatusGatewayTimeout() bool
	Status() int

	ExtractOAuthURL(hdr string, repository *url.URL) (*url.URL, error)
}

// URLTransporter struct
type URLTransporter struct {
	client *http.Client

	OAuthEndpoint *url.URL

	StatusCode int

	options Options
}

// NewURLTransporter creates a new URLTransporter
func NewURLTransporter(options Options) *URLTransporter {
	/* #nosec */
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.InsecureSkipVerify,
			RootCAs:            options.RootCAs,
		},
	}
	client := &http.Client{Transport: tr}

	return &URLTransporter{
		client:  client,
		options: options,
	}
}

// Put implements a PUT request
func (u *URLTransporter) Put(ctx context.Context, url *url.URL, body io.Reader, reqHdrs *http.Header, po progress.Output, ids ...string) (http.Header, error) {
	hdr, _, err := u.requestWithRetry(ctx, url, body, reqHdrs, http.MethodPut, po)
	return hdr, err
}

// Post implements a POST request
func (u *URLTransporter) Post(ctx context.Context, url *url.URL, body io.Reader, reqHdrs *http.Header, po progress.Output, ids ...string) (http.Header, error) {
	hdr, _, err := u.requestWithRetry(ctx, url, body, reqHdrs, http.MethodPost, po)
	return hdr, err
}

// Delete implements a DELETE request
func (u *URLTransporter) Delete(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, error) {
	hdr, _, err := u.requestWithRetry(ctx, url, nil, reqHdrs, http.MethodDelete, po)
	return hdr, err
}

// Head implements a HEAD request
func (u *URLTransporter) Head(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, error) {
	hdr, _, err := u.requestWithRetry(ctx, url, nil, reqHdrs, http.MethodHead, po)
	return hdr, err
}

// Get implements a GET request; the caller is responsible for closing the stream of the response body after use
func (u *URLTransporter) Get(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, io.ReadCloser, error) {
	return u.requestWithRetry(ctx, url, nil, reqHdrs, http.MethodGet, po)
}

// GetBytes returns the response body as []byte
func (u *URLTransporter) GetBytes(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, []byte, error) {
	hdr, rdr, err := u.requestWithRetry(ctx, url, nil, reqHdrs, http.MethodGet, po)

	defer rdr.Close()

	out := bytes.NewBuffer(nil)

	// Stream into it
	_, err = io.Copy(out, rdr)
	if err != nil {
		return nil, nil, err
	}
	return hdr, out.Bytes(), err
}

// GetHeaderOnly only returns the header of the response of a GET request
func (u *URLTransporter) GetHeaderOnly(ctx context.Context, url *url.URL, reqHdrs *http.Header, po progress.Output) (http.Header, error) {
	hdr, rdr, err := u.requestWithRetry(ctx, url, nil, reqHdrs, http.MethodGet, po)
	if rdr != nil {
		rdr.Close()
	}
	return hdr, err
}

// Create a basic request without retry
func (u *URLTransporter) request(ctx context.Context, url *url.URL, body io.Reader, reqHdrs *http.Header, operation string, po progress.Output) (http.Header, io.ReadCloser, error) {
	req, err := http.NewRequest(operation, url.String(), body)
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

	defer func() {
		if operation == http.MethodHead || operation == http.MethodDelete ||
			operation == http.MethodPost || operation == http.MethodPut {
			res.Body.Close()
		}
	}()

	log.Debugf("URLTransporter.request() - statuscode: %d, body: %#v, header: %#v", res.StatusCode, res.Body, res.Header)

	u.StatusCode = res.StatusCode

	if u.options.Token == nil && u.IsStatusUnauthorized() {
		// this is the case when fetching the auth token
		hdr := res.Header.Get("www-authenticate")
		if hdr == "" {
			return nil, nil, DoNotRetry{Err: fmt.Errorf("www-authenticate header is missing")}
		}
		return res.Header, nil, nil
	}

	if u.IsStatusUnauthorized() {
		hdr := res.Header.Get("www-authenticate")
		return nil, nil, DoNotRetry{Err: fmt.Errorf("unauthorized: %s", hdr)}
	}

	if u.IsStatusBadGateway() || u.IsStatusGatewayTimeout() || u.IsStatusServiceUnavailable() {
		return nil, nil, RetryOnErr{Err: fmt.Errorf("Network failure: statuscode: %d", res.StatusCode)}
	}

	if u.IsStatusOK() || u.IsStatusCreated() || u.IsStatusNoContent() || u.IsStatusAccepted() || u.IsStatusNotFound() {
		return res.Header, res.Body, nil
	}

	return nil, nil, DoNotRetry{Err: fmt.Errorf("Unexpected http code: %d, URL: %s", u.StatusCode, url)}
}

// Create a request with retry logic
func (u *URLTransporter) requestWithRetry(ctx context.Context, url *url.URL, body io.Reader, reqHdrs *http.Header, operation string, po progress.Output, ids ...string) (http.Header, io.ReadCloser, error) {
	defer trace.End(trace.Begin(operation + " " + url.Path))

	// extract ID from ids. Existence of an ID enables progress reporting
	ID := ""
	if len(ids) > 0 {
		ID = ids[0]
	}

	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), u.options.Timeout)
	defer cancel()

	var retries int

	for {
		hdr, bdr, err := u.request(ctx, url, body, reqHdrs, operation, po)
		if err == nil {
			return hdr, bdr, nil
		}

		// If an error was returned because the context was cancelled, we shouldn't retry.
		select {
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("cancelled during transporting")
		default:
		}

		retries++
		// give up if we reached maxTransportAttempts
		if retries == maxTransportAttempts {
			log.Debugf("Hit max transport attempts. Failed: %v", err)
			return nil, nil, err
		}

		switch err := err.(type) {
		case DoNotRetry:
			log.Debugf("Error: %s", err.Error())
			return nil, nil, err
		}

		// retry transporting again
		log.Debugf("Transporting failed, retrying: %v", err)

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
				return nil, nil, fmt.Errorf("cancelled during retry delay")
			}
		}
	}

}

// IsStatusUnauthorized returns true if status code is StatusUnauthorized
func (u *URLTransporter) IsStatusUnauthorized() bool {
	return u.StatusCode == http.StatusUnauthorized
}

// IsStatusOK returns true if status code is StatusOK
func (u *URLTransporter) IsStatusOK() bool {
	return u.StatusCode == http.StatusOK
}

// IsStatusNotFound returns true if status code is StatusNotFound
func (u *URLTransporter) IsStatusNotFound() bool {
	return u.StatusCode == http.StatusNotFound
}

// IsStatusAccepted returns true if status code is StatusAccepted
func (u *URLTransporter) IsStatusAccepted() bool {
	return u.StatusCode == http.StatusAccepted
}

// IsStatusCreated returns true if status code is StatusCreated
func (u *URLTransporter) IsStatusCreated() bool {
	return u.StatusCode == http.StatusCreated
}

// IsStatusNoContent returns true if status code is StatusNoContent
func (u *URLTransporter) IsStatusNoContent() bool {
	return u.StatusCode == http.StatusNoContent
}

// IsStatusBadGateway returns true if status code is StatusBadGateway
func (u *URLTransporter) IsStatusBadGateway() bool {
	return u.StatusCode == http.StatusBadGateway
}

// IsStatusServiceUnavailable returns true if status code is StatusServiceUnavailable
func (u *URLTransporter) IsStatusServiceUnavailable() bool {
	return u.StatusCode == http.StatusServiceUnavailable
}

// IsStatusGatewayTimeout returns true if status code is StatusGatewayTimeout
func (u *URLTransporter) IsStatusGatewayTimeout() bool {
	return u.StatusCode == http.StatusGatewayTimeout
}

func (u *URLTransporter) setUserAgent(req *http.Request) {
	log.Debugf("Setting user-agent to vic/%s", version.Version)
	req.Header.Set("User-Agent", "vic/"+version.Version)
}

func (u *URLTransporter) setBasicAuth(req *http.Request) {
	if u.options.Username != "" && u.options.Password != "" {
		log.Debugf("Setting BasicAuth: %s", u.options.Username)
		req.SetBasicAuth(u.options.Username, u.options.Password)
	}
}

func (u *URLTransporter) setAuthToken(req *http.Request) {
	if u.options.Token != nil {
		req.Header.Set("Authorization", "Bearer "+u.options.Token.Token)
	}
}

// ExtractOAuthURL extracts the OAuth url from the www-authenticate header
func (u *URLTransporter) ExtractOAuthURL(hdr string, repository *url.URL) (*url.URL, error) {

	log.Infof("the hdr in ExtractAuthURL is: %s", hdr)
	tokens := strings.Split(hdr, " ")
	if strings.ToLower(tokens[0]) != "bearer" {
		err := fmt.Errorf("www-authenticate header is corrupted")
		return nil, DoNotRetry{Err: err}
	}
	// Example for tokens[1]:
	// bearer realm=\"https://kang.eng.vmware.com/service/token\",
	// service=\"harbor-registry\",
	// scope=\"repository:test/busybox:pull,push repository:test/ubuntu:pull\"
	if len(tokens) == 3 {
		tokens[1] += " " + tokens[2]
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
			if len(tokens) == 4 {
				scope += ","
				scope += strings.Trim(tokens[len(tokens)-1], "\"")
			}
		}
	}

	if realm == "" {
		err := fmt.Errorf("missing realm in bearer auth challenge")
		return nil, DoNotRetry{Err: err}
	}
	if service == "" {
		err := fmt.Errorf("missing service in bearer auth challenge")
		return nil, DoNotRetry{Err: err}
	}
	// The scope can be empty if we're not getting a token for a specific repo
	if scope == "" && repository != nil {
		err := fmt.Errorf("missing scope in bearer auth challenge")
		return nil, DoNotRetry{Err: err}
	}
	log.Debugf("The service is: %s", service)
	log.Debugf("The realm is: %s", realm)
	log.Debugf("The scope is: %s", scope)
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

// URLDeepCopy deep-copys a URL
func URLDeepCopy(src *url.URL) *url.URL {
	dest := &url.URL{
		Scheme:     src.Scheme,
		Opaque:     src.Opaque,
		User:       src.User,
		Host:       src.Host,
		Path:       src.Path,
		RawPath:    src.RawPath,
		ForceQuery: src.ForceQuery,
		RawQuery:   src.RawQuery,
		Fragment:   src.Fragment,
	}

	return dest
}

// Status returns the status code
func (u *URLTransporter) Status() int {
	return u.StatusCode
}
