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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/vmware/vic/pkg/errors"
)

type Keypair struct {
	tlsGenerate bool
	keyFile     string
	certFile    string

	CertPEM []byte
	KeyPEM  []byte
}

func NewKeyPair(tlsGenerate bool, keyFile, certFile string) *Keypair {
	return &Keypair{
		tlsGenerate: tlsGenerate,
		keyFile:     keyFile,
		certFile:    certFile,
	}
}

// CreateRawKeyPair generates a default certificate / key and returns them as bytes buffers
// If you wish to save them to files as a side effect, use GetCertificate() instead
func CreateRawKeyPair() (cert bytes.Buffer, key bytes.Buffer, err error) {
	org := "VMware"
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return cert, key, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		err = errors.Errorf("Failed to generate random number: %s", err)
		return cert, key, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{org},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		err = errors.Errorf("Failed to generate x509 certificate: %s", err)
		return cert, key, err
	}

	err = pem.Encode(&cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		err = errors.Errorf("Failed to encode x509 certificate: %s", err)
		return cert, key, err
	}

	err = pem.Encode(&key, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	if err != nil {
		err = errors.Errorf("Failed to encode tls key pairs: %s", err)
		return cert, key, err
	}

	return cert, key, nil
}

func (k *Keypair) generate() error {
	cert, key, err := CreateRawKeyPair()
	if err != nil {
		return err
	}

	certFile, err := os.Create(k.certFile)
	if err != nil {
		err = errors.Errorf("Failed to create key/cert file %s: %s", k.certFile, err)
		return err
	}
	defer certFile.Close()

	_, err = certFile.Write(cert.Bytes())
	if err != nil {
		err = errors.Errorf("Failed to write certificate: %s", err)
		return err
	}

	keyFile, err := os.Create(k.keyFile)
	if err != nil {
		err = errors.Errorf("Failed to create key/cert file %s: %s", k.keyFile, err)
		return err
	}
	defer keyFile.Close()
	_, err = keyFile.Write(key.Bytes())
	if err != nil {
		err = errors.Errorf("Failed to write certificate: %s", err)
		return err
	}
	k.KeyPEM = key.Bytes()
	k.CertPEM = cert.Bytes()
	return nil
}

func (k *Keypair) GetCertificate() error {
	if k.tlsGenerate {
		return k.generate()
	}

	b, err := ioutil.ReadFile(k.certFile)
	if err != nil {
		err = errors.Errorf("Failed to read certificate file %s: %s", k.certFile, err)
		return err
	}

	k.CertPEM = b

	if b, err = ioutil.ReadFile(k.keyFile); err != nil {
		err = errors.Errorf("Failed to read key file %s: %s", k.keyFile, err)
		return err
	}

	k.KeyPEM = b
	return nil
}
