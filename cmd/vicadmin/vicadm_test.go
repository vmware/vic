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
	"archive/tar"
	"compress/gzip"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/vsphere/test/env"
)

// use an http client which we modify in init()
// to be permissive with certificates so we can
// use a self-signed cert hardcoded into these tests
var insecureClient *http.Client

func init() {
	sdk := env.URL(nil)
	if sdk != "" {
		flag.Set("sdk", sdk)
		flag.Set("vm-path", "docker-appliance")
		flag.Set("cluster", os.Getenv("GOVC_CLUSTER"))
	}

	// fake up a docker-host for pprof collection
	u := url.URL{Scheme: "http", Host: "127.0.0.1:6060"}

	go func() {
		log.Println(http.ListenAndServe(u.Host, nil))
	}()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	insecureClient = &http.Client{Transport: transport}
	flag.Set("docker-host", u.Host)

	config.hostCertFile = "fixtures/vicadmin_test_cert.pem"
	config.hostKeyFile = "fixtures/vicadmin_test_pkey.pem"

	cert, cerr := ioutil.ReadFile(config.hostCertFile)
	key, kerr := ioutil.ReadFile(config.hostKeyFile)
	if kerr != nil || cerr != nil {
		panic("unable to load test certificate")
	}
	vchConfig.HostCertificate = &metadata.RawCertificate{
		Cert: cert,
		Key:  key,
	}
}

type credentials struct {
	// the expected user
	username string

	// the expected password
	password string
}

// Checks credentials for the given user/password combo
func (c *credentials) Validate(u string, p string) bool {
	return u == c.username && p == c.password
}

func TestLoginFailure(t *testing.T) {
	// Authentication not yet implemented
	t.SkipNow()
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}

	s := &server{
		addr: "127.0.0.1:0",
		auth: &credentials{"root", "thisisinsecure"},
	}

	err := s.listen(true)
	assert.NoError(t, err)

	port := s.listenPort()

	go s.serve()
	defer s.stop()

	var res *http.Response
	res, err = insecureClient.Get(fmt.Sprintf("https://root:notthepassword@localhost:%d/", port))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestNoAuth(t *testing.T) {
	// Authentication not yet supported
	t.SkipNow()
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}

	s := &server{
		addr: "127.0.0.1:0",
		auth: nil,
	}

	err := s.listen(true)
	assert.NoError(t, err)

	port := s.listenPort()

	go s.serve()
	defer s.stop()

	var res *http.Response
	res, err = insecureClient.Get(fmt.Sprintf("https://localhost:%d/", port))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func testLogTar(t *testing.T, plainHTTP bool) {
	// Authentication not yet supported
	if plainHTTP {
		t.SkipNow()
	}

	if runtime.GOOS != "linux" {
		t.SkipNow()
	}

	logFileDir = "."

	s := &server{
		addr: "127.0.0.1:0",
		auth: &credentials{"root", "thisisinsecure"},
	}

	err := s.listen(!plainHTTP)
	assert.NoError(t, err)

	port := s.listenPort()

	go s.serve()
	defer s.stop()

	var res *http.Response
	if !plainHTTP {
		res, err = insecureClient.Get(fmt.Sprintf("https://root:thisisinsecure@localhost:%d/container-logs.tar.gz", port))
	} else {
		res, err = http.Get(fmt.Sprintf("http://root:thisisinsecure@localhost:%d/container-logs.tar.gz", port))
	}
	if err != nil {
		t.Fatal(err)
	}

	z, err := gzip.NewReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	tz := tar.NewReader(z)

	for {
		h, err := tz.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatal(err)
		}

		name, err := url.QueryUnescape(h.Name)
		if err != nil {
			t.Fatal(err)
		}

		if testing.Verbose() {
			fmt.Printf("\n%s...\n", name)
			io.CopyN(os.Stdout, tz, 150)
			fmt.Printf("...\n")
		}
	}
}

func TestLogTar(t *testing.T) {
	testLogTar(t, false)
	testLogTar(t, true)
}

func TestLogTail(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}

	f, err := os.OpenFile("./vicadmin.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString("# not much here yet\n")

	logFileDir = "."
	name := filepath.Base(f.Name())

	s := &server{
		addr: "127.0.0.1:0",
		// auth: &credentials{"root", "thisisinsecure"},
	}

	err = s.listen(true)
	assert.NoError(t, err)

	port := s.listenPort()

	go s.serve()
	defer s.stop()

	out := ioutil.Discard
	if testing.Verbose() {
		out = os.Stdout
	}

	paths := []string{
		"/logs/tail/" + name,
		"/logs/" + name,
	}

	i := 0

	u := url.URL{
		// User:   url.UserPassword("root", "thisisinsecure"),
		Scheme: "https",
		Host:   fmt.Sprintf("localhost:%d", port),
	}

	f.WriteString("this is line 0\n")
	log.Printf("Testing TestLogTail\n")
	for _, path := range paths {
		u.Path = path
		log.Printf("GET %s:\n", u.String())
		res, err := insecureClient.Get(u.String())
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			for j := 1; j < 512; j++ {
				i++
				f.WriteString(fmt.Sprintf("this is line %d\n", i))
			}
		}()

		size := int64(256)
		n, _ := io.CopyN(out, res.Body, size)
		out.Write([]byte("...\n"))
		res.Body.Close()

		assert.Equal(t, size, n)
	}
}
