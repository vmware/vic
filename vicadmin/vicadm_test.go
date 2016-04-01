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
)

// use an http client which we modify in init()
// to be permissive with certificates so we can
// use a self-signed cert hardcoded into these tests
var insecureClient *http.Client

func init() {
	sdk := os.Getenv("VIC_ESX_TEST_URL")
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
}

func testLogTar(t *testing.T, plainHTTP bool) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}

	logFileDir = "."

	s := &server{
		addr: "127.0.0.1:0",
	}

	err := s.listen(!plainHTTP)

	if err != nil {
		log.Fatal(err)
	}

	port := s.listenPort()

	go s.serve()

	var res *http.Response
	if !plainHTTP {
		res, err = insecureClient.Get(fmt.Sprintf("https://localhost:%d/container-logs.tar.gz", port))
	} else {
		res, err = http.Get(fmt.Sprintf("http://localhost:%d/container-logs.tar.gz", port))
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

	f, err := ioutil.TempFile("", "vicadm")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString("# not much here yet\n")

	logFileDir = "."
	name := filepath.Base(f.Name())

	s := &server{
		addr: "127.0.0.1:0",
	}

	err = s.listen(true)
	if err != nil {
		log.Fatal(err)
	}

	port := s.listenPort()

	go s.serve()

	out := ioutil.Discard
	if testing.Verbose() {
		out = os.Stdout
	}

	paths := []string{
		"/logs/tail/" + name,
		"/logs/tail",
		"/logs/" + name,
		"/",
	}

	i := 0

	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("localhost:%d", port),
	}

	for _, path := range paths {
		u.Path = path
		log.Printf("GET %s:\n", u)
		res, err := insecureClient.Get(u.String())
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			for j := 0; j < 512; j++ {
				i++
				f.WriteString(fmt.Sprintf("this is line %d\n", i))
			}
		}()

		size := int64(256)
		n, _ := io.CopyN(out, res.Body, size)
		out.Write([]byte("...\n"))
		res.Body.Close()

		if n != size {
			t.Errorf("expected %d, got %d", size, n)
		}
	}
}
