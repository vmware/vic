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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
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
	// sha256 sum of LayerContent
	DigestSHA256LayerContent = "sha256:18adac3bcad6124ed2e0d8dcc3beef8d540786ef8ef52c1f9fd71fdbfe36aa8e"

	//DigestSHA256EmptyTar is the canonical sha256 digest of empty data
	DigestSHA256EmptyTar = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

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
	options.digest = Tag

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

			manifest := &Manifest{
				Name:     Image,
				Tag:      Tag,
				FSLayers: []FSLayer{FSLayer{BlobSum: DigestSHA256EmptyTar}},
			}

			body, err := json.Marshal(manifest)
			if err != nil {
				t.Errorf(err.Error())
			}
			w.Write(body)

		}))
	defer s.Close()

	options.registry = s.URL
	options.image = Image
	options.digest = Tag
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
	if manifest.FSLayers[0].BlobSum != DigestSHA256EmptyTar {
		t.Errorf("Returned manifest %#v is different than expected", manifest)
	}
}

func TestFetchImageBlob(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-gzip")

			w.Write([]byte(LayerContent))
		}))
	defer s.Close()

	options.registry = s.URL
	options.image = Image
	options.digest = Tag
	options.token = &Token{Token: OAuthToken}

	// create a temporary directory
	dir, err := ioutil.TempDir("", "imagec")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer os.RemoveAll(dir)

	options.destination = dir

	err = FetchImageBlob(options, DigestSHA256LayerContent, LayerHistory)
	if err != nil {
		t.Errorf(err.Error())
	}

	tar, err := ioutil.ReadFile(path.Join(options.destination, Image, Tag, LayerID, LayerID+".tar"))
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(tar) != LayerContent {
		t.Errorf(err.Error())
	}

	hist, err := ioutil.ReadFile(path.Join(options.destination, Image, Tag, LayerID, LayerID+".json"))
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(hist) != LayerHistory {
		t.Errorf(err.Error())
	}
}
