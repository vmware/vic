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

package certificate

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"crypto/tls"

	log "github.com/Sirupsen/logrus"
)

const (
	keyFile  = "./key.pem"
	certFile = "./cert.pem"
)

func TestCreateRawKeyPair(t *testing.T) {
	cert, key, err := CreateRawKeyPair()
	if err != nil {
		t.Errorf("CreateRawKeyPair failed with error %s", err)
	}

	certString := cert.String()
	keyString := key.String()

	log.Infof("cert: %s", certString)
	log.Infof("key: %s", keyString)

	if !strings.HasPrefix(certString, "-----BEGIN CERTIFICATE-----") {
		t.Errorf("Certificate lacks proper prefix; must not have been generated properly.")
	}

	if !strings.HasSuffix(certString, "-----END CERTIFICATE-----\n") {
		t.Errorf("Certificate lacks proper suffix; must not have been generated properly.")
	}

	if !strings.HasPrefix(keyString, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Errorf("Private key lacks proper prefix; must not have been generated properly.")
	}

	if !strings.HasSuffix(keyString, "-----END RSA PRIVATE KEY-----\n") {
		t.Errorf("Private key lacks proper suffix; must not have been generated properly.")
	}

	_, err = tls.X509KeyPair([]byte(certString), []byte(keyString))
	if err != nil {
		t.Errorf("Unable to load X509 key pair(%s,%s): %s", certString, keyString, err)
	}

}

func TestGenerate(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	if _, err := os.Stat(keyFile); err == nil {
		os.Remove(keyFile)
	}

	pair := NewKeyPair(true, keyFile, certFile)
	if err := pair.GetCertificate(); err != nil {
		t.Errorf("%s", err)
	}
	log.Infof("key: %s", pair.KeyPEM)
	log.Infof("Cert: %s", pair.CertPEM)

	if _, err := os.Stat(keyFile); err != nil {
		t.Errorf("key file is not generated")
	}
	if !strings.Contains(string(pair.KeyPEM), "RSA PRIVATE KEY") {
		t.Errorf("Key is not correctly generated")
	}
}

func TestGetCertificate(t *testing.T) {
	pair := NewKeyPair(true, keyFile, certFile)
	if err := pair.GetCertificate(); err != nil {
		t.Errorf("%s", err)
	}
	keyPEM := pair.KeyPEM

	pair = NewKeyPair(false, keyFile, certFile)
	if err := pair.GetCertificate(); err != nil {
		t.Errorf("%s", err)
	}

	if !bytes.Equal(pair.KeyPEM, keyPEM) {
		log.Errorf("Expected pem: %s", keyPEM)
		log.Errorf("Actual pem: %s", pair.KeyPEM)
		t.Errorf("key is not correctly read out")
	}
}
