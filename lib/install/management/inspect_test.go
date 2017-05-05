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
	"os"
	"testing"

	"github.com/vmware/vic/pkg/certificate"

	"github.com/stretchr/testify/assert"
)

func TestVerifyClientCert(t *testing.T) {
	cacert, cakey, err := certificate.CreateRootCA("foo.com", []string{"FooOrg"}, 2048)
	assert.NoError(t, err)

	cert, key, err := certificate.CreateClientCertificate("foo.com", []string{"FooOrg"}, 2048, cacert.Bytes(), cakey.Bytes())
	assert.NoError(t, err)

	kp := certificate.NewKeyPair(ClientCert, ClientKey, cert.Bytes(), key.Bytes())
	err = kp.SaveCertificate()
	assert.NoError(t, err)
	defer func() {
		os.Remove(ClientCert)
		os.Remove(ClientKey)
	}()

	// Validate client certificate keypair created with the right CA
	_, err = VerifyClientCert(cacert.Bytes(), kp)
	assert.NoError(t, err)

	cacert, cakey, err = certificate.CreateRootCA("bar.com", []string{"BarOrg"}, 2048)
	assert.NoError(t, err)

	// Attempt to validate client certificate keypair created with a different CA
	_, err = VerifyClientCert(cacert.Bytes(), kp)
	assert.NotNil(t, err)
}
