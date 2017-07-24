// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest/schema2"
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
	OAuthToken = "Top_Secret_Token"
	Image      = "library/photon"
	Tag        = "latest"
	Reference  = Image + ":" + Tag

	BusyboxImage = "library/busybox"

	// fake content
	LayerContent = "Cannot_Contain_Myself"
	// fake ID
	LayerID = "f9767cae14f372c98900f15bb07cb40b8e1a6d1507912489e1342db499313d32"
	// fake history
	LayerHistory = "{\"id\":\"f9767cae14f372c98900f15bb07cb40b8e1a6d1507912489e1342db499313d32\"" + "," +
		"\"parent\":\"09a5baea69e9c781d64df5366c36492d53d507048035abd68632264dc23a1edb\"}"
	// fake store
	Storename = "PetStore"

	// sha256 sum of LayerContent
	DigestSHA256LayerContent = "sha256:1f1f9635040c465c7f7b32a396e56e26c1396dbadfbed744b0aab9337a24ad5a"

	//DigestSHA256EmptyData is the canonical sha256 digest of empty data
	DigestSHA256EmptyData = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	DefaultManifest = `
{
	"schemaVersion": 1,
	"name": "library/photon",
	"tag": "latest",
	"architecture": "amd64",
	"fsLayers": [
			{
				"blobSum": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
			}
	],
	"history": [
			{
				"v1Compatibility": "{\"architecture\":\"amd64\",\"config\":{\"Hostname\":\"156e10b83429\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":[\"sh\"],\"Image\":\"56ed16bd6310cca65920c653a9bb22de6b235990dcaa1742ff839867aed730e5\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"container\":\"5f8098ec29947b5bea80483cd3275008911ce87438fed628e34ec0c522665510\",\"container_config\":{\"Hostname\":\"156e10b83429\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) CMD [\\\"sh\\\"]\"],\"Image\":\"56ed16bd6310cca65920c653a9bb22de6b235990dcaa1742ff839867aed730e5\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"created\":\"2016-03-18T18:22:48.810791943Z\",\"docker_version\":\"1.9.1\",\"id\":\"437595becdebaaaf3a4fc3db02c59a980f955dee825c153308c670610bb694e1\",\"os\":\"linux\",\"parent\":\"920777304d1d5e337bc59877253e946f224df5aae64c72538672eb74637b3c9e\"}"
			}
	],
	"signatures": [
			{
				"header": {
						"jwk": {
							"crv": "P-256",
							"kid": "LUQI:WBTB:JRDU:TTD2:FUVY:EMCB:64HP:MZF6:SGFS:XAB6:JPUK:6PK4",
							"kty": "EC",
							"x": "zjkAuFGpCuWOBl-iMMzZqgl_1cid-04S04-k-A1qEeU",
							"y": "9HWcOMfVFUMXJGeNajIAlPicL4UOsCJSpqRcIxpUl0Q"
						},
						"alg": "ES256"
				},
				"signature": "dTXnnt3IkTScpZhyyqRlmZFcQV1QzD7lWDqnjlD4Cj-KsMuGd1pl5QpFL2Cadw-8KeTBlSleSecxjHU4t3yhCQ",
				"protected": "eyJmb3JtYXRMZW5ndGgiOjE5MjIsImZvcm1hdFRhaWwiOiJDbjAiLCJ0aW1lIjoiMjAxNi0wNi0wOVQxNzoyMzo1NFoifQ"
			}
	]
}
	`
)

func TestParseReference(t *testing.T) {

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	ref, err := reference.ParseNamed("index.docker.io/library/busybox")
	ic.Options.Reference = ref
	ic.ParseReference()
	assert.Equal(t, ic.Tag, reference.DefaultTag)
	assert.Equal(t, ic.Image, BusyboxImage)
	assert.Equal(t, ic.Registry, DefaultDockerURL)

	ref, err = reference.ParseNamed("vmware/photon")
	ic.Options.Reference = ref
	ic.ParseReference()
	assert.Equal(t, ic.Tag, reference.DefaultTag)
	assert.Equal(t, ic.Image, "vmware/photon")
	assert.Equal(t, ic.Registry, DefaultDockerURL)

	ref, err = reference.ParseNamed("busybox")
	if err != nil {
		t.Errorf(err.Error())
	}
	ic.Options.Reference = ref
	ic.ParseReference()
	assert.Equal(t, ic.Tag, reference.DefaultTag)
	assert.Equal(t, ic.Image, BusyboxImage)
	assert.Equal(t, ic.Registry, DefaultDockerURL)

	ref, err = reference.ParseNamed("library/busybox")
	if err != nil {
		t.Errorf(err.Error())
	}
	ic.Options.Reference = ref
	ic.ParseReference()
	assert.Equal(t, ic.Tag, reference.DefaultTag)
	assert.Equal(t, ic.Image, BusyboxImage)
	assert.Equal(t, ic.Registry, DefaultDockerURL)

	ref, err = reference.ParseNamed("library/busybox:latest")
	if err != nil {
		t.Errorf(err.Error())
	}
	ic.Options.Reference = ref
	ic.ParseReference()
	assert.Equal(t, ic.Tag, reference.DefaultTag)
	assert.Equal(t, ic.Image, BusyboxImage)
	assert.Equal(t, ic.Registry, DefaultDockerURL)

	ic = NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)
	ref, err = reference.ParseNamed("busybox")
	if err != nil {
		t.Errorf(err.Error())
	}
	digest, err := digest.ParseDigest("sha256:c79345819a6882c31b41bc771d9a94fc52872fa651b36771fbe0c8461d7ee558")
	if err != nil {
		t.Errorf(err.Error())
	}
	ref, err = reference.WithDigest(reference.TrimNamed(ref), digest)
	if err != nil {
		t.Errorf(err.Error())
	}
	ic.Options.Reference = ref
	ic.ParseReference()
	assert.Equal(t, ic.Tag, "")
	assert.Equal(t, ic.Image, BusyboxImage)
	assert.Equal(t, ic.Registry, DefaultDockerURL)
}

func TestLearnRegistryURL(t *testing.T) {

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
			http.Error(w, "You shall not pass", http.StatusUnauthorized)
		}))
	defer s.Close()

	ic.Options.Registry = s.URL[7:]
	ic.Options.Image = Image
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout

	// should fail
	_, err := LearnRegistryURL(&ic.Options)
	if err == nil {
		t.Errorf(err.Error())
	}

	// should pass
	ic.Options.InsecureAllowHTTP = true
	_, err = LearnRegistryURL(&ic.Options)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestLearnAuthURL(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("www-authenticate",
				"Bearer realm=\"https://auth.docker.io/token\",service=\"registry.docker.io\",scope=\"repository:library/photon:pull\"")
			http.Error(w, "You shall not pass", http.StatusUnauthorized)
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = Image
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	url, err := LearnAuthURL(ic.Options)
	if err != nil {
		t.Errorf(err.Error())
	}

	if url.String() != "https://auth.docker.io/token?scope=repository%3Alibrary%2Fphoton%3Apull&service=registry.docker.io" {
		t.Errorf("Returned url %s is different than expected", url)
	}
}

func TestFetchToken(t *testing.T) {

	options := Options{
		Outstream: os.Stdout,
		Timeout:   DefaultHTTPTimeout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			body, err := json.Marshal(&urlfetcher.Token{Token: OAuthToken})
			if err != nil {
				t.Errorf(err.Error())
			}
			w.Write(body)

		}))
	defer s.Close()

	url, err := url.Parse(s.URL)
	if err != nil {
		t.Errorf(err.Error())
	}
	url.Path = path.Join(url.Path, "token?scope=repository%3Alibrary%2Fphoton%3Apull&service=registry.docker.io")

	ctx := context.TODO()
	token, err := FetchToken(ctx, ic.Options, url, ic.progressOutput)
	if err != nil {
		t.Errorf(err.Error())
	}

	if token.Token != OAuthToken {
		t.Errorf("Returned token %s is different than expected", token.Token)
	}
}

func TestFetchImageManifest(t *testing.T) {
	var err error

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(DefaultManifest))

		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = Image
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Token = &urlfetcher.Token{Token: OAuthToken}
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	// create a temporary directory
	dir, err := ioutil.TempDir("", "imagec")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer os.RemoveAll(dir)

	ic.Options.Destination = dir

	ctx := context.TODO()
	manifest, _, err := FetchImageManifest(ctx, ic.Options, 1, ic.progressOutput)
	if err != nil {
		t.Errorf(err.Error())
	}
	if schema1, ok := manifest.(*Manifest); ok {
		if schema1.FSLayers[0].BlobSum != DigestSHA256EmptyData {
			t.Errorf("Returned manifest %#v is different than expected", manifest)
		}
	}
}

func TestFetchImageBlob(t *testing.T) {

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	// create a tar archive from our dummy data
	r := strings.NewReader(LayerContent)
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// create and populate tar header
	header := new(tar.Header)
	header.Size = r.Size()
	header.Name = LayerID
	header.Mode = 0644

	// write the header
	if err := tw.WriteHeader(header); err != nil {
		t.Errorf(err.Error())
	}

	// write the file into the tar archive
	if _, err := io.Copy(tw, r); err != nil {
		t.Error(err.Error())
	}

	attempt := 3
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempt--
			if attempt == 0 {
				w.Header().Set("Content-Type", "application/x-gzip")
				// return our tar archive
				w.Write(buf.Bytes())
			} else {
				// return an error to force caller to retry
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = Image
	ic.Options.Tag = Tag
	ic.Options.Timeout = DefaultHTTPTimeout
	ic.Options.Token = &urlfetcher.Token{Token: OAuthToken}

	// create a temporary directory
	dir, err := ioutil.TempDir("", "imagec")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer os.RemoveAll(dir)

	ic.Options.Destination = dir

	parent := "scratch"
	image := ImageWithMeta{
		Image: &models.Image{
			ID:     LayerID,
			Parent: parent,
			Store:  Storename,
		},
		Meta:  LayerHistory,
		Layer: FSLayer{BlobSum: DigestSHA256LayerContent},
	}
	diffID, err := FetchImageBlob(ctx, ic.Options, &image, ic.progressOutput)
	if err != nil {
		t.Errorf(err.Error())
	}
	if diffID == "" {
		t.Errorf("Expected a diffID, got nil.")
	}

	tarFile, err := ioutil.ReadFile(path.Join(DestinationDirectory(ic.Options), LayerID, LayerID+".tar"))
	if err != nil {
		t.Errorf(err.Error())
	}
	br := bytes.NewReader(tarFile)
	tr := tar.NewReader(br)
	out := new(bytes.Buffer)

	// extract the tar file
	for {
		_, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf(err.Error())
		}
		if _, err := io.Copy(out, tr); err != nil {
			t.Errorf(err.Error())
		}
	}

	// compare contents of tar file to dummy data
	if out.String() != LayerContent {
		t.Errorf(err.Error())
	}

	hist, err := ioutil.ReadFile(path.Join(DestinationDirectory(ic.Options), LayerID, LayerID+".json"))
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(hist) != LayerHistory {
		t.Errorf(err.Error())
	}
}

func TestPingPortLayer(t *testing.T) {

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("\"OK\""))
		}))
	defer s.Close()

	ic.Options.Host = s.URL[7:]

	b, err := PingPortLayer(ic.Host)
	if err != nil || b != true {
		t.Errorf(err.Error())
	}
}

func TestListImages(t *testing.T) {

	options := Options{
		Outstream: os.Stdout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := `
[
  {
    "ID": "7bd023c8937ded982c1b98da453b1a5afec86f390ffad8fa0f4fba244a6155f1",
    "Metadata": {
      "v1Compatibility": "{\"id\":\"7bd023c8937ded982c1b98da453b1a5afec86f390ffad8fa0f4fba244a6155f1\",\"parent\":\"b873f334fa5259acb24cf0e2cd2639d3a9fb3eb9bafbca06ed4f702c289b31c0\",\"created\":\"2016-05-27T14:15:02.359284074Z\",\"container\":\"b8bd6a8e8874a87f626871ce370f4775bdf598865637082da2949ee0f4786432\",\"container_config\":{\"Hostname\":\"914cf42a3e15\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) CMD [\\\"/bin/bash\\\"]\"],\"Image\":\"b873f334fa5259acb24cf0e2cd2639d3a9fb3eb9bafbca06ed4f702c289b31c0\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"docker_version\":\"1.9.1\",\"config\":{\"Hostname\":\"914cf42a3e15\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[],\"Cmd\":[\"/bin/bash\"],\"Image\":\"b873f334fa5259acb24cf0e2cd2639d3a9fb3eb9bafbca06ed4f702c289b31c0\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"architecture\":\"amd64\",\"os\":\"linux\"}"
    },
    "SelfLink": "http://Photon/storage/7bd023c8937ded982c1b98da453b1a5afec86f390ffad8fa0f4fba244a6155f1",
    "Store": "http://Photon/storage/PetStore"
  }
]
`
			w.Write([]byte(resp))
		}))
	defer s.Close()

	ic.Options.Host = s.URL[7:]

	m, err := ListImages(ic.Host, ic.Storename, nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	if m["7bd023c8937ded982c1b98da453b1a5afec86f390ffad8fa0f4fba244a6155f1"].ID != "7bd023c8937ded982c1b98da453b1a5afec86f390ffad8fa0f4fba244a6155f1" {
		t.Errorf("Returned list %#v is different than expected", m)
	}
}

func TestFetchScenarios(t *testing.T) {
	var err error
	ctx := context.TODO()

	options := Options{
		Outstream: os.Stdout,
		Timeout:   DefaultHTTPTimeout,
	}

	ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)

	wwwAuthenticate := false
	insufficientScope := false
	invalidToken := false

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				if wwwAuthenticate {
					w.Header().Set("www-authenticate", "Bearer realm=\"https://auth.docker.io/token\",service=\"registry.docker.io\",scope=\"repository:library/busybox:pull\"")
				}
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				if !insufficientScope && !invalidToken {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(DefaultManifest))
				} else if insufficientScope {
					w.Header().Set("www-authenticate", "error=\"insufficient_scope\"")
					w.WriteHeader(http.StatusUnauthorized)
				} else if invalidToken {
					w.Header().Set("www-authenticate", "error=\"invalid_token\"")
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
				}
			}
		}))
	defer s.Close()

	ic.Options.Registry = s.URL
	ic.Options.Image = Image
	ic.Options.Tag = Tag
	ic.Options.InsecureAllowHTTP = true
	ic.Options.Reference, err = reference.ParseNamed(Reference)
	if err != nil {
		t.Errorf(err.Error())
	}

	// create a temporary directory
	dir, err := ioutil.TempDir("", "imagec")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer os.RemoveAll(dir)

	ic.Options.Destination = dir

	ic.Options.Token = nil
	// try without token, response will miss www-authenticate so we will retry (and fail eventually)
	_, _, err = FetchImageManifest(ctx, ic.Options, 1, ic.progressOutput)
	if err == nil {
		t.Errorf("Condition didn't fail (testing failure)")
	}

	// set the www-authenticate
	wwwAuthenticate = true
	// try without token, response will carry a valid www-authenticate so we won't retry
	_, _, err = FetchImageManifest(ctx, ic.Options, 1, ic.progressOutput)
	if err != nil {
		// we should get a DNR error
		if _, isDNR := err.(urlfetcher.DoNotRetry); !isDNR {
			t.Errorf(err.Error())
		}
	}
	wwwAuthenticate = false

	// set a valid token
	ic.Options.Token = &urlfetcher.Token{Token: OAuthToken}

	// enable invalid token test
	invalidToken = true
	// valid token but faulty hub, we should retry but eventually fail
	_, _, err = FetchImageManifest(ctx, ic.Options, 1, ic.progressOutput)
	if err == nil {
		t.Errorf("Condition didn't fail (testing failure)")
	}
	invalidToken = false

	// enable insufficient_scope test
	insufficientScope = true
	// valid token but image is missing we shouldn't retry
	_, _, err = FetchImageManifest(ctx, ic.Options, 1, ic.progressOutput)
	if err != nil {
		// we should get a ImageNotFoundError
		if _, imageErr := err.(urlfetcher.ImageNotFoundError); !imageErr {
			t.Errorf(err.Error())
		}
	}
	insufficientScope = false

	// valid token, existing image, correct header. We should succeed
	_, _, err = FetchImageManifest(ctx, ic.Options, 1, ic.progressOutput)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPrepareManifestAndLayers(t *testing.T) {
	type PushMock struct {
		image string
	}
	// schema1Manifests := map[string]schema1.Manifest{
	// 	"ubuntu": {
	// 		Name: "library/busybox",
	// 		Tag:  "latest",
	// 		// Digest: "sha256:74155e76e3396bc23d5b8cedf606bc04f7f79660ee5bce14014b84651fc4f672",
	// 		FSLayers: []schema1.FSLayer{
	// 			{
	// 				BlobSum: "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4",
	// 			},
	// 			{
	// 				BlobSum: "sha256:27144aa8f1b9e066514d7f765909367584e552915d0d4bc2f5b7438ba7d1033a",
	// 			},
	// 		},
	// 		History: []schema1.History{
	// 			{
	// 				V1Compatibility: "{\"architecture\":\"amd64\",\"config\":{\"Hostname\":\"c673fc810c50\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"sh\"],\"ArgsEscaped\":true,\"Image\":\"sha256:7b537995b09bda336a22b3c139cfbef751d4361d506880ea559fb7e2180f291f\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"container\":\"aef9e475482dbd32adca70d243921a21c53b24820c3a6e764d86d65e0506cb2d\",\"container_config\":{\"Hostname\":\"c673fc810c50\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) \",\"CMD [\\\"sh\\\"]\"],\"ArgsEscaped\":true,\"Image\":\"sha256:7b537995b09bda336a22b3c139cfbef751d4361d506880ea559fb7e2180f291f\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"created\":\"2017-06-15T20:42:30.659714375Z\",\"docker_version\":\"17.03.1-ce\",\"id\":\"bbed08f07a6bccc8aca4f6053dd1b5bdf1050f830e0989738e6532dd4a703a58\",\"os\":\"linux\",\"parent\":\"a8159aa8135af6a6ffa715cae0e7f1595198a304828146231de3062811619857\",\"throwaway\":true}",
	// 			},
	// 			{
	// 				V1Compatibility: "{\"id\":\"a8159aa8135af6a6ffa715cae0e7f1595198a304828146231de3062811619857\",\"created\":\"2017-06-15T20:42:07.973730453Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) ADD file:aa56bc8f2fea9c0c81ca085bfa273ad1a3b0d46f51b8c9c61b483340c902024f in / \"]}}",
	// 			},
	// 		},
	// 	},
	// 	"alpine": {
	// 		Name: "library/alpine",
	// 		Tag:  "latest",
	// 		// Digest: "sha256:9c2c06bdf64688f49b70eea9613a83b2880d8d62f941a5b07a3053c14a64f970",
	// 		FSLayers: []schema1.FSLayer{
	// 			{
	// 				BlobSum: "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4",
	// 			},
	// 			{
	// 				BlobSum: "sha256:88286f41530e93dffd4b964e1db22ce4939fffa4a4c665dab8591fbab03d4926",
	// 			},
	// 		},
	// 		History: []schema1.History{
	// 			{
	// 				V1Compatibility: "{\"architecture\":\"amd64\",\"config\":{\"Hostname\":\"e1ede117fb1e\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\"],\"ArgsEscaped\":true,\"Image\":\"sha256:ac1fc1931356fa238379d061cb216c4bed2f150991298c20b166accf0604d3b1\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"container\":\"6a3726d15fee7d345097a26ebc1b9e5b4de25dee759a921def9059a4d0cd2261\",\"container_config\":{\"Hostname\":\"e1ede117fb1e\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) \",\"CMD [\\\"/bin/sh\\\"]\"],\"ArgsEscaped\":true,\"Image\":\"sha256:ac1fc1931356fa238379d061cb216c4bed2f150991298c20b166accf0604d3b1\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{}},\"created\":\"2017-06-27T18:42:16.849872208Z\",\"docker_version\":\"17.03.1-ce\",\"id\":\"5511ab30246fb138ba659434d8391e24d324b65ed042bef3d28ac717174b5747\",\"os\":\"linux\",\"parent\":\"472b4fc7eba3e7f6588d6965c0e14caaea23f5e61194c19f9c498ca6963ecaa8\",\"throwaway\":true}",
	// 			},
	// 			{
	// 				V1Compatibility: "{\"id\":\"472b4fc7eba3e7f6588d6965c0e14caaea23f5e61194c19f9c498ca6963ecaa8\",\"created\":\"2017-06-27T18:41:51.5382773Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) ADD file:4583e12bf5caec40b861a3409f2a1624c3f3556cc457edb99c9707f00e779e45 in / \"]}}",
	// 			},
	// 		},
	// 	},
	// }

	// ubuntu, alpine
	schema2Manifests := map[string]schema2.Manifest{
		"ubuntu": {
			Versioned: schema2.SchemaVersion,
			// Versioned: {
			// 	SchemaVersion: 2,
			// 	MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
			// },
			Config: distribution.Descriptor{
				MediaType: "application/vnd.docker.container.image.v1+json",
				Size:      1507,
				Digest:    "sha256:c30178c5239f2937c21c261b0365efcda25be4921ccb95acd63beeeb78786f27",
			},
			Layers: []distribution.Descriptor{
				{
					MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
					Size:      699243,
					Digest:    "sha256:27144aa8f1b9e066514d7f765909367584e552915d0d4bc2f5b7438ba7d1033a",
				},
			},
		},
		"alpine": {
			Versioned: schema2.SchemaVersion,
			// {
			// SchemaVersion: 2,
			// MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
			// },
			Config: distribution.Descriptor{
				MediaType: "application/vnd.docker.container.image.v1+json",
				Size:      1520,
				Digest:    "sha256:7328f6f8b41890597575cbaadc884e7386ae0acc53b747401ebce5cf0d624560",
			},
			Layers: []distribution.Descriptor{
				{
					MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
					Size:      1990402,
					Digest:    "sha256:88286f41530e93dffd4b964e1db22ce4939fffa4a4c665dab8591fbab03d4926",
				},
			},
		},
	}

	// schema2Actual := schema2.Manifest{
	// 	Versioned: manifest.Versioned{
	// 		SchemaVersion: 2,
	// 		MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
	// 	},
	// 	Config: distribution.Descriptor{
	// 		MediaType: "application/vnd.docker.container.image.v1+json",
	// 		Size:      1507,
	// 		Digest:    "c30178c5239f2937c21c261b0365efcda25be4921ccb95acd63beeeb78786f27",
	// 	},
	// }

	mockData := []PushMock{
		{
			image: "ubuntu",
		},
		{
			image: "alpine",
		},
	}

	for i, data := range mockData {
		ref, err := reference.ParseNamed(data.image)
		assert.NoError(t, err, "Could not parse reference from image %s", data.image)

		options := Options{
			Destination: os.TempDir(),
			Reference:   ref,
			Timeout:     3600 * time.Second,
			Outstream:   nil,
		}

		ic := NewImageC(options, streamformatter.NewJSONStreamFormatter(), nil)
		err = ic.PrepareManifestAndLayers()
		if err != nil {
			assert.NoError(t, err, "Manifest preparation failed:  %s", err.Error())
			continue
		}

		expectedManifest := schema2Manifests[data.image]
		actualManifest := ic.Pusher.PushManifest

		// Check version
		assert.Equal(t, expectedManifest.MediaType, actualManifest.MediaType, "Mediatype of the manifest.Versioned is incorrect for mock #%d", i)
		assert.Equal(t, expectedManifest.SchemaVersion, actualManifest.SchemaVersion, "Schema version of the manifest.Versioned is incorrect for mock #%d", i)

	}
}

func TestUpdateV2MetaData(t *testing.T) {
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

	layer2 := &ImageWithMeta{
		Image: &models.Image{
			ID:     uid.New().String(),
			Parent: layer1.ID,
		},
		V2Meta: []dmetadata.V2Metadata{{
			SourceRepository: "docker.io/library/busybox",
		}},
	}
	layerCache.Add(layer2)

	ref, err := reference.ParseNamed("busybox:latest")
	if err != nil {
		t.Error(err.Error())
	}
	imageID := uid.New()

	err = cache.RepositoryCache().AddReference(ref, imageID.String(), false, layer2.ID, false)
	if err != nil {
		t.Error(err.Error())
	}

	err = cache.RepositoryCache().AddReference(ref, imageID.String(), false, layer1.ID, false)
	if err != nil {
		t.Error(err.Error())
	}

	// try to update V2MetaData
	newSourceRepo := "docker.io/library/busybox"

	err = UpdateV2MetaData(ref, newSourceRepo)
	assert.NoError(t, err, "UpdataeV2MetaData failed: %s", err)

	l1, err := LayerCache().Get(layer1.ID)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 1, len(l1.V2Meta), "The number of source repositories should be one!")

	// try to update V2MetaData
	newSourceRepo = "docker.io/test/anotherBusybox"

	err = UpdateV2MetaData(ref, newSourceRepo)
	assert.NoError(t, err, "UpdataeV2MetaData failed: %s", err)

	l2, err := LayerCache().Get(layer2.ID)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 2, len(l2.V2Meta), "The number of source repositories should be two!")
}
