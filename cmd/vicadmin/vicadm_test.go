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

	vchconfig "github.com/vmware/vic/lib/config"
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

	hostCertFile := "fixtures/vicadmin_test_cert.pem"
	hostKeyFile := "fixtures/vicadmin_test_pkey.pem"

	cert, cerr := ioutil.ReadFile(hostCertFile)
	key, kerr := ioutil.ReadFile(hostKeyFile)
	if kerr != nil || cerr != nil {
		panic("unable to load test certificate")
	}
	vchConfig.HostCertificate = &vchconfig.RawCertificate{
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
	}

	err := s.listen()
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
	}

	err := s.listen()
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
	}

	err := s.listen()
	assert.NoError(t, err)

	port := s.listenPort()

	go s.serve()
	defer s.stop()

	var res *http.Response
	res, err = insecureClient.Get(fmt.Sprintf("https://root:thisisinsecure@localhost:%d/container-logs.tar.gz", port))
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

	err = s.listen()
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
	}

	u := url.URL{
		// User:   url.UserPassword("root", "thisisinsecure"),
		Scheme: "https",
		Host:   fmt.Sprintf("localhost:%d", port),
	}

	str := "The quick brown fox jumps over the lazy dog.\n"

	// Pre-populate the log file
	for i := 0; i < tailLines; i++ {
		f.WriteString(str)
	}

	log.Printf("Testing TestLogTail\n")
	for _, path := range paths {
		u.Path = path
		log.Printf("GET %s:\n", u.String())
		res, err := insecureClient.Get(u.String())
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		// Each line written to the log file has enough bytes to ensure
		// that 8 lines make up at least 256 bytes. This is in case this
		// goroutine finishes writing all lines before tail kicks in.
		go func() {
			for j := 1; j < 512; j++ {
				f.WriteString(str)
			}
			f.Sync()
		}()

		size := int64(256)
		n, err := io.CopyN(out, res.Body, size)
		assert.NoError(t, err)
		out.Write([]byte("...\n"))

		assert.Equal(t, size, n)
	}
}

type seekTest struct {
	input  string
	output int
}

func testSeek(t *testing.T, st seekTest, td string) {
	f, err := ioutil.TempFile(td, "FindSeekPos")
	defer f.Close()
	if err != nil {
		log.Printf("Unable to create temporary file: %s", err)
		t.Fatal(err)
	}
	n, err := f.WriteString(st.input)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(st.input) {
		t.Fatal(fmt.Errorf("Incorrect byte count on write: %d/%d", n, len(st.input)))
	}

	if ret := findSeekPos(f); ret != int64(st.output) {
		t.Fatal(fmt.Errorf("Incorrect seek position: %d/%d", ret, st.output))
	}
	log.Printf("Successfully seeked to position %d", st.output)
	os.Remove(f.Name())
}
func TestFindSeekPos(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}

	seekTests := []seekTest{
		{"abcd\nabcd\n", 0},
		{"abcd\nabcd\nabcd\nabcd\nabcd\nabcd\nabcd\nabcd\nabcd\nabcd\n", 10},
	}

	str := "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	str = fmt.Sprintf("%s%s", str, str)
	str2 := fmt.Sprintf("%s\n%s\n%s\n", str, str, str)
	str2 = fmt.Sprintf("%s%s", str2, str2)
	// Verify we don't have overlapping reads at beginning of file
	// 6 lines, 1206 characters. Should come back with seek position 0
	seekTests = append(seekTests, seekTest{str2, 0})

	for seg := 0; seg < 3; seg++ {
		str = str + str
	}
	str2 = str + "\n"
	fmt.Printf("str length is %d\n", len(str))
	// str is 1,601 chars long now
	for line := 0; line < 2; line++ {
		str2 = str2 + str2
	}

	// str2 is 4 lines long.  Should seek to beginning of file
	seekTests = append(seekTests, seekTest{str2, 0})

	// str2 is now 12 lines long.  Should seek to position 6,404
	str2 = fmt.Sprintf("%s%s%s", str2, str2, str2)
	seekTests = append(seekTests, seekTest{str2, 6404})

	td := os.TempDir()
	for i, st := range seekTests {
		log.Printf("Test case #%d: ", i)
		testSeek(t, st, td)
	}
}
