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
	"bufio"
	"bytes"
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

	"github.com/Sirupsen/logrus"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
)

const (
	OAuthToken = "Top_Secret_Token"
	Image      = "library/photon"
	Tag        = "latest"

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
	options.reference = "busybox"
	if err := ParseReference(); err != nil {
		t.Errorf(err.Error())
	}

	options.reference = "library/busybox"
	if err := ParseReference(); err != nil {
		t.Errorf(err.Error())
	}

	options.reference = "library/busybox:latest"
	if err := ParseReference(); err != nil {
		t.Errorf(err.Error())
	}

	// should fail
	options.reference = "library/busybox@invalid"
	if err := ParseReference(); err == nil {
		t.Errorf(err.Error())
	}
}

func TestLearnRegistryURL(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
			http.Error(w, "You shall not pass", http.StatusUnauthorized)
		}))
	defer s.Close()

	options.registry = s.URL[7:]
	options.image = Image
	options.tag = Tag

	// should fail
	_, err := LearnRegistryURL(options)
	if err == nil {
		t.Errorf(err.Error())
	}

	// should pass
	options.insecureAllowHTTP = true
	_, err = LearnRegistryURL(options)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestLearnAuthURL(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("www-authenticate",
				"Bearer realm=\"https://auth.docker.io/token\",service=\"registry.docker.io\",scope=\"repository:library/photon:pull\"")
			http.Error(w, "You shall not pass", http.StatusUnauthorized)
		}))
	defer s.Close()

	options.registry = s.URL
	options.image = Image
	options.tag = Tag

	url, err := LearnAuthURL(options)
	if err != nil {
		t.Errorf(err.Error())
	}

	if url.String() != "https://auth.docker.io/token?scope=repository%3Alibrary%2Fphoton%3Apull&service=registry.docker.io" {
		t.Errorf("Returned url %s is different than expected", url)
	}
}

func TestFetchToken(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			body, err := json.Marshal(&Token{Token: OAuthToken})
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

	token, err := FetchToken(url)
	if err != nil {
		t.Errorf(err.Error())
	}

	if token.Token != OAuthToken {
		t.Errorf("Returned token %s is different than expected", token.Token)
	}
}

func TestFetchImageManifest(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(DefaultManifest))

		}))
	defer s.Close()

	options.registry = s.URL
	options.image = Image
	options.tag = Tag
	options.token = &Token{Token: OAuthToken}

	// create a temporary directory
	dir, err := ioutil.TempDir("", "imagec")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer os.RemoveAll(dir)

	options.destination = dir

	manifest, err := FetchImageManifest(options)
	if err != nil {
		t.Errorf(err.Error())
	}
	if manifest.FSLayers[0].BlobSum != DigestSHA256EmptyData {
		t.Errorf("Returned manifest %#v is different than expected", manifest)
	}
}

func TestFetchImageBlob(t *testing.T) {
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

	options.registry = s.URL
	options.image = Image
	options.tag = Tag
	options.token = &Token{Token: OAuthToken}

	// create a temporary directory
	dir, err := ioutil.TempDir("", "imagec")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer os.RemoveAll(dir)

	options.destination = dir

	parent := "scratch"
	image := ImageWithMeta{
		Image: &models.Image{
			ID:     LayerID,
			Parent: &parent,
			Store:  Storename,
		},
		meta:  LayerHistory,
		layer: FSLayer{BlobSum: DigestSHA256LayerContent},
	}
	diffID, err := FetchImageBlob(options, &image)
	if err != nil {
		t.Errorf(err.Error())
	}
	if diffID == "" {
		t.Errorf("Expected a diffID, got nil.")
	}

	tarFile, err := ioutil.ReadFile(path.Join(DestinationDirectory(), LayerID, LayerID+".tar"))
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

	hist, err := ioutil.ReadFile(path.Join(DestinationDirectory(), LayerID, LayerID+".json"))
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(hist) != LayerHistory {
		t.Errorf(err.Error())
	}
}

func TestPingPortLayer(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("\"OK\""))
		}))
	defer s.Close()

	options.host = s.URL[7:]

	b, err := PingPortLayer()
	if err != nil || b != true {
		t.Errorf(err.Error())
	}
}

func TestListImages(t *testing.T) {
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

	options.host = s.URL[7:]

	m, err := ListImages(Storename, nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	if m["7bd023c8937ded982c1b98da453b1a5afec86f390ffad8fa0f4fba244a6155f1"].ID != "7bd023c8937ded982c1b98da453b1a5afec86f390ffad8fa0f4fba244a6155f1" {
		t.Errorf("Returned list %#v is different than expected", m)
	}
}

func TestNewHook(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	hook := NewErrorHook(w)
	levels := hook.Levels()
	if len(levels) != 1 && levels[0] != logrus.FatalLevel {
		t.Errorf("Returned levels %#v are different than expected", levels)
	}

	w.Reset(&b)
	hook.Fire(&logrus.Entry{Message: "Fatal Test", Level: logrus.FatalLevel})
	w.Flush()

	if string(b.Bytes()) == "" {
		t.Errorf("Fatal test failed %s", string(b.Bytes()))
	}
}

func TestArrayUpdate(t *testing.T) {
	attributes := []string{"green", "blue", "red"}

	// remove an existing value
	attributes, chg := arrayUpdate("green", attributes, Remove)
	if !chg && len(attributes) == 2 {
		t.Error("arrayUpdate failed to remove existing value")
	}

	// add a new value
	attributes, chg = arrayUpdate("green", attributes, Add)
	if !chg && len(attributes) == 3 {
		t.Error("arrayUpdate failed to add new value")
	}

	// add existing value
	attributes, chg = arrayUpdate("green", attributes, Add)
	if chg && len(attributes) == 3 {
		t.Error("arrayUpdate failed while adding existing value")
	}

	// attempt to remove non-existent value
	attributes, chg = arrayUpdate("foobar", attributes, Remove)
	if chg && len(attributes) == 3 {
		t.Error("arrayUpdate failed while removing non-existent value")
	}

}

func TestFetchScenarios(t *testing.T) {
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

	options.registry = s.URL
	options.image = Image
	options.tag = Tag

	// create a temporary directory
	dir, err := ioutil.TempDir("", "imagec")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer os.RemoveAll(dir)

	options.destination = dir

	options.token = nil
	// try without token, response will miss www-authenticate so we will retry (and fail eventually)
	_, err = FetchImageManifest(options)
	if err == nil {
		t.Errorf("Condition didn't fail (testing failure)")
	}

	// set the www-authenticate
	wwwAuthenticate = true
	// try without token, response will carry a valid www-authenticate so we won't retry
	_, err = FetchImageManifest(options)
	if err != nil {
		// we should get a DNR error
		if _, isDNR := err.(DoNotRetry); !isDNR {
			t.Errorf(err.Error())
		}
	}
	wwwAuthenticate = false

	// set a valid token
	options.token = &Token{Token: OAuthToken}

	// enable invalid token test
	invalidToken = true
	// valid token but faulty hub, we should retry but eventually fail
	_, err = FetchImageManifest(options)
	if err == nil {
		t.Errorf("Condition didn't fail (testing failure)")
	}
	invalidToken = false

	// enable insufficient_scope test
	insufficientScope = true
	// valid token but image is missing we shouldn't retry
	_, err = FetchImageManifest(options)
	if err != nil {
		// we should get a ImageNotFoundError
		if _, imageErr := err.(ImageNotFoundError); !imageErr {
			t.Errorf(err.Error())
		}
	}
	insufficientScope = false

	// valid token, existing image, correct header. We should succeed
	_, err = FetchImageManifest(options)
	if err != nil {
		t.Errorf(err.Error())
	}
}
