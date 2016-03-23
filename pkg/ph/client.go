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

package ph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/vmware/vic/pkg/errors"

	log "github.com/Sirupsen/logrus"
)

var (
	ProductServer = "ph"
	StagingServer = "ph-stg"

	CollectorID        = "vic.1_0"
	DefaultContentType = "application/json"

	addrFmt = "https://vcsa.vmware.com/%s/api/hyper/send?_v=1.0&_c=%s&_i=%s"
)

type Client struct {
	// phone home address.
	phURL     *url.URL
	Transport *http.Transport
	doneChan  chan bool
}

// NewClient is to create new Phone Home client
// instance ID should be a VC object UUID string without the dashes, for example, the vApp UUID
func NewClient(instanceID string) *Client {
	address := fmt.Sprintf(addrFmt, StagingServer, CollectorID, instanceID)
	log.Debugf("Phone Home address: %s", address)

	u, err := url.Parse(address)
	if err != nil {
		log.Errorf("Failed to parse Phone Home URL for %s", errors.ErrorStack(err))
		return nil
	}

	return &Client{
		phURL: u,
	}
}

func (phc *Client) SetPHAddress(addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		log.Errorf("Invalid Phone Home address %s, error ", addr, errors.ErrorStack(err))
		return errors.Trace(err)
	}

	phc.phURL = u
	return nil
}

func (phc *Client) httpClient() *http.Client {
	if phc.Transport == nil {
		// using DefaultTransport
		return &http.Client{}
	}
	return &http.Client{Transport: phc.Transport}
}

func (phc *Client) encodeData(data interface{}) (*bytes.Buffer, error) {
	params := bytes.NewBuffer(nil)
	if data != nil {
		if err := json.NewEncoder(params).Encode(data); err != nil {
			return nil, errors.Trace(err)
		}
	}
	return params, nil
}

func (phc *Client) buildHeaders(req *http.Request) {
	req.Header.Set("Content-Type", DefaultContentType)
	req.Header.Set("Accept", DefaultContentType)
}

// POST method posts metrics data to VMware phone home
func (phc *Client) POST(data interface{}) error {
	log.Debugf("data: %s", data)
	in, err := phc.encodeData(data)
	if err != nil {
		return errors.Trace(err)
	}
	log.Debugf("encoding json: %s", in)

	// create request
	req, err := http.NewRequest("POST", phc.phURL.String(), in)
	if err != nil {
		return errors.Trace(err)
	}

	// build request headers
	phc.buildHeaders(req)

	resp, err := phc.httpClient().Do(req)
	statusCode := -1
	if resp != nil {
		statusCode = resp.StatusCode

		if resp.Body != nil {
			defer resp.Body.Close()
		}
	}

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return errors.New("cannot connect to PhoneHome server")
		}

		return errors.Errorf("an error occurred trying to connect phone home: %s", errors.ErrorStack(err))
	}

	if statusCode != http.StatusCreated {
		log.Debugf("Response status code is %s", statusCode)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Trace(err)
		}
		return errors.Errorf("error response: %s", bytes.TrimSpace(body))
	}

	return nil
}
