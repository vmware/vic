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

package management

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

var (
	errSecretKeyNotFound = fmt.Errorf("unable to find guestinfo secret")
	errNilDatastore      = fmt.Errorf("session's datastore is not set")
)

// extractSecretFromFile reads and extracts the GuestInfoSecretKey value from the input.
func extractSecretFromFile(rc io.ReadCloser) (string, error) {

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()

		// The line is of the format: key = "value"
		if strings.HasPrefix(line, extraconfig.GuestInfoSecretKey) {

			tokens := strings.SplitN(line, "=", 2)
			if len(tokens) < 2 {
				return "", fmt.Errorf("parse error: unexpected token count in line")
			}

			// Ensure that the key fully matches the secret key
			if strings.Trim(tokens[0], ` `) != extraconfig.GuestInfoSecretKey {
				continue
			}

			// Trim double quotes and spaces
			return strings.Trim(tokens[1], `" `), nil
		}
	}

	return "", errSecretKeyNotFound
}

// GuestInfoSecret downloads the VCH's .vmx file and returns the GuestInfoSecretKey
// value. This function expects the datastore in the dispatcher's session to be set.
func (d *Dispatcher) GuestInfoSecret(vchName string) (*extraconfig.SecretKey, error) {
	defer trace.End(trace.Begin(""))

	if d.session.Datastore == nil {
		return nil, errNilDatastore
	}

	helper, err := datastore.NewHelper(d.ctx, d.session, d.session.Datastore, d.vmPathName)
	if err != nil {
		return nil, err
	}

	// Download the VCH's .vmx file
	path := fmt.Sprintf("%s.vmx", vchName)
	rc, err := helper.Download(d.ctx, path)
	if err != nil {
		return nil, err
	}

	secret, err := extractSecretFromFile(rc)
	if err != nil {
		return nil, err
	}

	secretKey := &extraconfig.SecretKey{}
	if err = secretKey.FromString(secret); err != nil {
		return nil, err
	}

	return secretKey, nil
}
