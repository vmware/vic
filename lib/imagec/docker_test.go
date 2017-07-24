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

package imagec

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	dmetadata "github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/reference"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/portlayer/storage"
	urlfetcher "github.com/vmware/vic/pkg/fetcher"
	"github.com/vmware/vic/pkg/uid"
)

const (
	UbuntuTaggedRef      = "library/ubuntu:latest"
	UbuntuDigest         = "ubuntu@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2"
	UbuntuDigestSHA      = "sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2"
	UbuntuDigestManifest = `{
   "schemaVersion": 1,
   "name": "library/ubuntu",
   "tag": "14.04",
   "architecture": "amd64",
   "fsLayers": [
      {
         "blobSum": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
      },
      {
         "blobSum": "sha256:28a2f68d1120598986362662445c47dce7ec13c2662479e7aab9f0ecad4a7416"
      },
      {
         "blobSum": "sha256:fd2731e4c50ce221d785d4ce26a8430bca9a95bfe4162fafc997a1cc65682cce"
      },
      {
         "blobSum": "sha256:5a132a7e7af11f304041e93efb9cb2a0a7839bccaec5a03cfbdc9a3f5d0eb481"
      }
   ],
   "history": [
      {
         "v1Compatibility": "{\"id\":\"56063ad57855f2632f641a622efa81a0feda1731bfadda497b1288d11feef4da\",\"parent\":\"4e1f7c524148bf80fcc4ce9991e88d708048d38440e3e3a14d56e72c17ddccc7\",\"created\":\"2016-03-03T21:38:53.80360049Z\",\"container\":\"b6361ab0a2a82f71c5bd3becbb9c854331f8259b9c3fe466bf6e7e073c735a2c\",\"container_config\":{\"Hostname\":\"c24ffee5b2b8\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) CMD [\\\"/bin/bash\\\"]\"],\"Image\":\"4e1f7c524148bf80fcc4ce9991e88d708048d38440e3e3a14d56e72c17ddccc7\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"docker_version\":\"1.9.1\",\"config\":{\"Hostname\":\"c24ffee5b2b8\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[],\"Cmd\":[\"/bin/bash\"],\"Image\":\"4e1f7c524148bf80fcc4ce9991e88d708048d38440e3e3a14d56e72c17ddccc7\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"architecture\":\"amd64\",\"os\":\"linux\"}"
      },
      {
         "v1Compatibility": "{\"id\":\"4e1f7c524148bf80fcc4ce9991e88d708048d38440e3e3a14d56e72c17ddccc7\",\"parent\":\"38112156678df7d8001ae944f118d283009565540dc0cd88fb39fccc88c3c7f2\",\"created\":\"2016-03-03T21:38:53.085760873Z\",\"container\":\"ccc6ec8b31df981344b4107bd890394be35564adb8d13df34d1cb1849c7f0c3e\",\"container_config\":{\"Hostname\":\"c24ffee5b2b8\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[],\"Cmd\":[\"/bin/sh\",\"-c\",\"sed -i 's/^#\\\\s*\\\\(deb.*universe\\\\)$/\\\\1/g' /etc/apt/sources.list\"],\"Image\":\"38112156678df7d8001ae944f118d283009565540dc0cd88fb39fccc88c3c7f2\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"docker_version\":\"1.9.1\",\"config\":{\"Hostname\":\"c24ffee5b2b8\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[],\"Cmd\":null,\"Image\":\"38112156678df7d8001ae944f118d283009565540dc0cd88fb39fccc88c3c7f2\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":1895}"
      },
      {
         "v1Compatibility": "{\"id\":\"38112156678df7d8001ae944f118d283009565540dc0cd88fb39fccc88c3c7f2\",\"parent\":\"454970bd163ba95435b50e963edd63b2b2fff4c1845e5d3cd03d5ba8afb8a08d\",\"created\":\"2016-03-03T21:38:51.45368726Z\",\"container\":\"3c8556d1a209f22cfbc87f3cbd9bcb6674c5f9645a14aa488756d129f6987f40\",\"container_config\":{\"Hostname\":\"c24ffee5b2b8\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[],\"Cmd\":[\"/bin/sh\",\"-c\",\"echo '#!/bin/sh' \\u003e /usr/sbin/policy-rc.d \\t\\u0026\\u0026 echo 'exit 101' \\u003e\\u003e /usr/sbin/policy-rc.d \\t\\u0026\\u0026 chmod +x /usr/sbin/policy-rc.d \\t\\t\\u0026\\u0026 dpkg-divert --local --rename --add /sbin/initctl \\t\\u0026\\u0026 cp -a /usr/sbin/policy-rc.d /sbin/initctl \\t\\u0026\\u0026 sed -i 's/^exit.*/exit 0/' /sbin/initctl \\t\\t\\u0026\\u0026 echo 'force-unsafe-io' \\u003e /etc/dpkg/dpkg.cfg.d/docker-apt-speedup \\t\\t\\u0026\\u0026 echo 'DPkg::Post-Invoke { \\\"rm -f /var/cache/apt/archives/*.deb /var/cache/apt/archives/partial/*.deb /var/cache/apt/*.bin || true\\\"; };' \\u003e /etc/apt/apt.conf.d/docker-clean \\t\\u0026\\u0026 echo 'APT::Update::Post-Invoke { \\\"rm -f /var/cache/apt/archives/*.deb /var/cache/apt/archives/partial/*.deb /var/cache/apt/*.bin || true\\\"; };' \\u003e\\u003e /etc/apt/apt.conf.d/docker-clean \\t\\u0026\\u0026 echo 'Dir::Cache::pkgcache \\\"\\\"; Dir::Cache::srcpkgcache \\\"\\\";' \\u003e\\u003e /etc/apt/apt.conf.d/docker-clean \\t\\t\\u0026\\u0026 echo 'Acquire::Languages \\\"none\\\";' \\u003e /etc/apt/apt.conf.d/docker-no-languages \\t\\t\\u0026\\u0026 echo 'Acquire::GzipIndexes \\\"true\\\"; Acquire::CompressionTypes::Order:: \\\"gz\\\";' \\u003e /etc/apt/apt.conf.d/docker-gzip-indexes\"],\"Image\":\"454970bd163ba95435b50e963edd63b2b2fff4c1845e5d3cd03d5ba8afb8a08d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"docker_version\":\"1.9.1\",\"config\":{\"Hostname\":\"c24ffee5b2b8\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[],\"Cmd\":null,\"Image\":\"454970bd163ba95435b50e963edd63b2b2fff4c1845e5d3cd03d5ba8afb8a08d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":194533}"
      },
      {
         "v1Compatibility": "{\"id\":\"454970bd163ba95435b50e963edd63b2b2fff4c1845e5d3cd03d5ba8afb8a08d\",\"created\":\"2016-03-03T21:38:46.169812943Z\",\"container\":\"c24ffee5b2b808674d335e2c42c9c47aa9aff1b368eb5920018cde7dd26f2046\",\"container_config\":{\"Hostname\":\"c24ffee5b2b8\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ADD file:b9504126dc55908988977286e89c43c7ea73a506d43fae82c29ef132e21b6ece in /\"],\"Image\":\"\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":null},\"docker_version\":\"1.9.1\",\"config\":{\"Hostname\":\"c24ffee5b2b8\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":null,\"Image\":\"\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":null},\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":187763841}"
      }
   ],
   "signatures": [
      {
         "header": {
            "jwk": {
               "crv": "P-256",
               "kid": "G2JA:NPRD:EWRM:EFEJ:4PHQ:TRZR:6W6O:D5LC:UJ36:RHOE:ZN7D:N55I",
               "kty": "EC",
               "x": "NrETepARqTLeOBcTdBCE8K8jbQJgiTH1p7XJ78zBxjk",
               "y": "ay0SmJatkJs-JdnW80807CcNPWHElsh6MW_JTh7NdbU"
            },
            "alg": "ES256"
         },
         "signature": "kHbWdD1NQw2RIAQ8uYnKmolU3Z_WeUW8DfRtJRprVDzK7AZaF-ChI4V9Lh74HnjSNwoNZ_QRhUQDl_Nezb0Hgw",
         "protected": "eyJmb3JtYXRMZW5ndGgiOjY3NzksImZvcm1hdFRhaWwiOiJDbjAiLCJ0aW1lIjoiMjAxNy0wNS0xN1QxODowMjozMFoifQ"
      }
   ]
}
`
	WrongDigest        = "sha256:12345dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2"
	MockImage          = "test/ubuntu"
	MockUploadLocation = "/v2/test/ubuntu/blobs/uploads/" +
		"00f4bfb1-682e-4a2b-86c5-8ace83e70cba?_state=rh9dXg5Vn5dZ_LQa5u9HN7RHdo3DyyuJ3GQowFn2jvt7Ik5hbWUiOiJjaGVuZy10ZXN0L2J1c3lib3h2NSIsIlVVSUQiOiIwMGY0YmZiMS02ODJlLTRhMmItODZjNS04YWNlODNlNzBjYmEiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTctMDctMTdUMTk6MDg6MjMuNjE0NDQyNzg4WiJ9"
	WrongUploadLocation = "/v2/wrong/ubuntu/blobs/uploads/" +
		"00f4bfb1-682e-4a2b-86c5-8ace83e70cba?_state=rh9dXg5Vn5dZ_LQa5u9HN7RHdo3DyyuJ3GQowFn2jvt7Ik5hbWUiOiJjaGVuZy10ZXN0L2J1c3lib3h2NSIsIlVVSUQiOiIwMGY0YmZiMS02ODJlLTRhMmItODZjNS04YWNlODNlNzBjYmEiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTctMDctMTdUMTk6MDg6MjMuNjE0NDQyNzg4WiJ9"
	RepoMounted    = "test/ubuntu1"
	RepoNotMounted = "test/busybox"
	RepoRandom     = "randomRepo"
)

var (
	RepoNotMountedEncoded  string
	RepoMountedEncoded     string
	UbuntuDigestSHAEncoded string
)

func init() {
	RepoNotMountedEncoded = url.Values{
		"from": {RepoNotMounted},
	}.Encode()

	RepoMountedEncoded = url.Values{
		"from": {RepoMounted},
	}.Encode()

	UbuntuDigestSHAEncoded = url.Values{
		"": {UbuntuDigestSHA},
	}.Encode()
}

func TestGetManifestDigest(t *testing.T) {
	// Get the manifest content when the image is not pulled by digest
	ref, err := reference.ParseNamed(UbuntuTaggedRef)
	if err != nil {
		t.Errorf(err.Error())
	}
	digest, err := getManifestDigest([]byte(UbuntuDigestManifest), ref)
	assert.NoError(t, err)
	assert.Equal(t, digest, UbuntuDigestSHA)

	// Get and verify the manifest content with the correct digest
	ref, err = reference.ParseNamed(UbuntuDigest)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, ok := ref.(reference.Canonical)
	assert.True(t, ok)
	digest, err = getManifestDigest([]byte(UbuntuDigestManifest), ref)
	assert.NoError(t, err)
	assert.Equal(t, digest, UbuntuDigestSHA)

	// Attempt to get and verify an incorrect manifest content with the digest
	digest, err = getManifestDigest([]byte(DefaultManifest), ref)
	assert.NotNil(t, err)
}

func TestLearnAuthURLForPush(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("www-authenticate",
				"Bearer realm=\"https://auth.docker.io/token\",service=\"registry.docker.io\",scope=\"repository:test/ubuntu:pull,push\"")
			http.Error(w, "You shall not pass", http.StatusUnauthorized)
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = MockImage
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	url, err := LearnAuthURLForPush(ic.Options, ic.progressOutput)
	if err != nil {
		t.Errorf(err.Error())
	}

	if url.String() != "https://auth.docker.io/token?scope=repository%3Atest%2Fubuntu%3Apull%2Cpush&service=registry.docker.io" {
		t.Errorf("Returned url %s is different than expected", url)
	}
}

func TestLearnAuthURLForBlobMount(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("www-authenticate",
				"Bearer realm=\"https://auth.docker.io/token\",service=\"registry.docker.io\",scope=\"repository:test/ubuntu:pull,push repository:test/ubuntu1:pull\"")
			http.Error(w, "You shall not pass", http.StatusUnauthorized)
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = MockImage
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	url, err := LearnAuthURLForBlobMount(ic.Options, UbuntuDigestSHA, RepoMounted, ic.progressOutput)
	if err != nil {
		t.Errorf(err.Error())
	}

	if url.String() != "https://auth.docker.io/token?scope=repository%3Atest%2Fubuntu%3Apull%2Cpush+repository%3Atest%2Fubuntu1%3Apull&service=registry.docker.io" {
		t.Errorf("Returned url %s is different than expected", url)
	}
}

func TestCheckLayerExistence(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, UbuntuDigestSHA) {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = MockImage
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            ic.Options.Timeout,
		Username:           ic.Options.Username,
		Password:           ic.Options.Password,
		Token:              ic.Options.Token,
		InsecureSkipVerify: ic.Options.InsecureSkipVerify,
		RootCAs:            ic.Options.RegistryCAs,
	})

	layerID := "MockLayer"

	registryURL, err := url.Parse(ic.Options.Registry)
	if err != nil {
		t.Errorf(err.Error())
	}

	// this layer should exist
	pushDigest := UbuntuDigestSHA
	exist, err := CheckLayerExistence(ctx, transporter, options.Image, pushDigest, registryURL, ic.progressOutput)
	if err != nil {
		t.Errorf("failed to check for presence of layer %s (%s) in %s: %s", layerID, pushDigest, options.Image, err)
	}
	assert.Equal(t, true, exist, "Layer should exist!")

	// this layer should not exist since the digest is wrong
	exist, err = CheckLayerExistence(ctx, transporter, options.Image, WrongDigest, registryURL, ic.progressOutput)
	if err != nil {
		t.Errorf("failed to check for presence of layer %s (%s) in %s: %s", layerID, pushDigest, options.Image, err)
	}
	assert.Equal(t, false, exist, "Layer should not exist!")
}

func TestObtainUploadURL(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Location", MockUploadLocation)
			w.WriteHeader(http.StatusAccepted)
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = MockImage
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            ic.Options.Timeout,
		Username:           ic.Options.Username,
		Password:           ic.Options.Password,
		Token:              ic.Options.Token,
		InsecureSkipVerify: ic.Options.InsecureSkipVerify,
		RootCAs:            ic.Options.RegistryCAs,
	})

	registryURL, err := url.Parse(ic.Options.Registry)
	if err != nil {
		t.Errorf(err.Error())
	}

	uploadURL, err := ObtainUploadURL(ctx, transporter, registryURL, options.Image, ic.progressOutput)
	if err != nil {
		t.Errorf("failed to obtain url for uploading layer: %s", err)
	}
	assert.Equal(t, MockUploadLocation, uploadURL, "UploadURL is wrong!")
}

func TestCancelUpload(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.URL.String())
			if strings.Contains(r.URL.String(), MockUploadLocation) {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = MockImage
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            ic.Options.Timeout,
		Username:           ic.Options.Username,
		Password:           ic.Options.Password,
		Token:              ic.Options.Token,
		InsecureSkipVerify: ic.Options.InsecureSkipVerify,
		RootCAs:            ic.Options.RegistryCAs,
	})

	url := fmt.Sprintf("%s%s", s.URL, MockUploadLocation)
	fmt.Println(url)

	err = CancelUpload(ctx, transporter, url, ic.progressOutput)
	assert.NoError(t, err, "CancelUpload is expected to succeed!")

	url = fmt.Sprintf("%s%s", s.URL, WrongUploadLocation)
	fmt.Println(url)

	err = CancelUpload(ctx, transporter, url, ic.progressOutput)
	assert.Error(t, err, "CancelUpload is expected to fail due to wrong upload location!")
}

func TestCompletedUpload(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("The request url is: ", r.URL.String())
			if strings.Contains(r.URL.String(), MockUploadLocation) && strings.Contains(r.URL.String(), UbuntuDigestSHAEncoded) {
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = MockImage
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            ic.Options.Timeout,
		Username:           ic.Options.Username,
		Password:           ic.Options.Password,
		Token:              ic.Options.Token,
		InsecureSkipVerify: ic.Options.InsecureSkipVerify,
		RootCAs:            ic.Options.RegistryCAs,
	})

	url := fmt.Sprintf("%s%s", s.URL, MockUploadLocation)
	fmt.Println(url)

	err = CompletedUpload(ctx, transporter, UbuntuDigestSHA, url, ic.progressOutput)
	assert.NoError(t, err, "CompletedUpload is expected to succeed!")

	err = CompletedUpload(ctx, transporter, WrongDigest, url, ic.progressOutput)
	assert.Error(t, err, "CompletedUpload is expected to fail due to wrong digest!")
}

func TestUploadLayer(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("The request url is: ", r.URL.String())
			if strings.Contains(r.URL.String(), MockUploadLocation) && strings.Contains(r.URL.String(), UbuntuDigestSHAEncoded) {
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = MockImage
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            ic.Options.Timeout,
		Username:           ic.Options.Username,
		Password:           ic.Options.Password,
		Token:              ic.Options.Token,
		InsecureSkipVerify: ic.Options.InsecureSkipVerify,
		RootCAs:            ic.Options.RegistryCAs,
	})

	url := fmt.Sprintf("%s%s", s.URL, MockUploadLocation)
	fmt.Println(url)

	err = UploadLayer(ctx, transporter, UbuntuDigestSHA, url, bytes.NewReader([]byte("")), ic.progressOutput)
	assert.NoError(t, err, "UploadLayer is expected to succeed!")

	err = UploadLayer(ctx, transporter, WrongDigest, url, bytes.NewReader([]byte("")), ic.progressOutput)
	assert.Error(t, err, "UploadLayer is expected to fail due to wrong digest!")

	url = fmt.Sprintf("%s%s", s.URL, WrongUploadLocation)
	fmt.Println(url)

	err = UploadLayer(ctx, transporter, UbuntuDigestSHA, url, bytes.NewReader([]byte("")), ic.progressOutput)
	assert.Error(t, err, "UploadLayer is expected to fail due to wrong upload location!")
}

func TestMountBlobToRepo(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("The request url is: ", r.URL.String())
			if strings.Contains(r.URL.String(), RepoMountedEncoded) {
				w.WriteHeader(http.StatusCreated)
			} else if strings.Contains(r.URL.String(), RepoNotMountedEncoded) {
				w.WriteHeader(http.StatusAccepted)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = MockImage
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	transporter := urlfetcher.NewURLTransporter(urlfetcher.Options{
		Timeout:            ic.Options.Timeout,
		Username:           ic.Options.Username,
		Password:           ic.Options.Password,
		Token:              ic.Options.Token,
		InsecureSkipVerify: ic.Options.InsecureSkipVerify,
		RootCAs:            ic.Options.RegistryCAs,
	})

	registryURL, err := url.Parse(ic.Registry)
	if err != nil {
		t.Errorf(err.Error())
	}

	mounted, _, err := MountBlobToRepo(ctx, transporter, registryURL, UbuntuDigestSHA, ic.Image, RepoMounted, ic.progressOutput)
	assert.NoError(t, err, "MountBlobToRepo is expected to succeed!")
	assert.Equal(t, true, mounted, "The layer should have been mounted!")

	mounted, _, err = MountBlobToRepo(ctx, transporter, registryURL, UbuntuDigestSHA, ic.Image, RepoNotMounted, ic.progressOutput)
	assert.NoError(t, err, "MountBlobToRepo is expected to succeed!")
	assert.Equal(t, false, mounted, "The layer should not have been mounted!")

	mounted, _, err = MountBlobToRepo(ctx, transporter, registryURL, UbuntuDigestSHA, ic.Image, RepoRandom, ic.progressOutput)
	assert.Error(t, err, "MountBlobToRepo is expected to fail!")
	assert.Equal(t, false, mounted, "The layer should not have been mounted!")
}

func TestObtainSourceRepoList(t *testing.T) {
	layerCache = &LCache{
		layers: make(map[string]*ImageWithMeta),
	}

	// Add some fake data to the layer cache
	layer1 := &ImageWithMeta{
		Image: &models.Image{
			ID:     uid.New().String(),
			Parent: storage.Scratch.ID,
		},
		V2Meta: []dmetadata.V2Metadata{{
			SourceRepository: "docker.io/library/busybox",
		}},
	}
	layerCache.Add(layer1)

	ref, err := reference.ParseNamed("busybox:latest")
	if err != nil {
		t.Error(err.Error())
	}
	imageID := uid.New()

	err = cache.RepositoryCache().AddReference(ref, imageID.String(), false, layer1.ID, false)
	if err != nil {
		t.Error(err.Error())
	}

	newSourceRepo := "docker.io/test/busybox"
	err = UpdateV2MetaData(ref, newSourceRepo)
	if err != nil {
		t.Error(err.Error())
	}

	targetRepo, err := reference.ParseNamed("busybox:latest")
	if err != nil {
		t.Error(err.Error())
	}
	rl, err := ObtainSourceRepoList(layer1.ID, targetRepo)
	assert.NoError(t, err, "Failed to obtain source repository list: %s", err)
	assert.Equal(t, 1, len(rl), "SourceRepoList should contain one element!")

	newSourceRepo = "harbor-registry.com/test/busybox"
	err = UpdateV2MetaData(ref, newSourceRepo)
	if err != nil {
		t.Error(err.Error())
	}

	rl, err = ObtainSourceRepoList(layer1.ID, targetRepo)
	assert.NoError(t, err, "Failed to obtain source repository list: %s", err)
	assert.Equal(t, 1, len(rl), "SourceRepoList should contain one element!")
}
