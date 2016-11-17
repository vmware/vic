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

package create

import (
	"flag"
	"fmt"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	create = NewCreate()
)

func TestGenKey(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	os.Args = []string{"cmd", "create"}
	flag.Parse()
	create.noTLS = false
	create.certPath = "install-test"
	create.cname = "common name"
	create.keySize = 1024

	ca, kp, err := create.generateCertificates(true, true)
	defer os.RemoveAll(fmt.Sprintf("./%s", create.certPath))

	assert.NoError(t, err, "Expected to cleanly generate certificates")
	assert.NotEmpty(t, ca, "Expected CA to contain data")
	assert.NotNil(t, kp, "Expected keypair to contain data")
	assert.NotEmpty(t, kp.CertPEM, "Expected certificate to contain data")
	assert.NotEmpty(t, kp.CertPEM, "Expected key to contain data")

	ca, kp, err = create.loadCertificates()
	assert.NoError(t, err, "Expected to cleanly load certificates")
	assert.NotEmpty(t, ca, "Expected CA to contain data")
	assert.NotNil(t, kp, "Expected keypair to contain data")
	assert.NotEmpty(t, kp.CertPEM, "Expected certificate to contain data")
	assert.NotEmpty(t, kp.CertPEM, "Expected key to contain data")
}
