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

package backends

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/backend"

	"github.com/vmware/vic/lib/imagec"
)

func getMockReader() (io.ReadCloser, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new tar archive.
	tw := tar.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling license."},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}
	// Make sure to check the error on Close.
	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}

	// Open the tar archive for reading.
	r := bytes.NewReader(buf.Bytes())
	return ioutil.NopCloser(r), nil
}

func TestDownload(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	rc, err := getMockReader()
	if err != nil {
		t.Errorf("Failed to get mocked reader: %s", err)
	}

	tests := []struct {
		repo string
		tag  string
	}{
		{repo: "registry-1.docker.io", tag: ""},
		{repo: "registry-1.docker.io", tag: "mycommit"},
		{repo: "myrepo.io", tag: ""},
		{repo: "myrepo.io", tag: "mycommit"},
		{repo: "", tag: ""},
	}
	for _, test := range tests {
		config := &backend.ContainerCommitConfig{}
		config.Tag = test.tag
		config.Repo = test.repo
		ic, err := getImagec(config)
		if err != nil {
			t.Errorf("Failed to get imagec: %s", err)
			return
		}
		layer, err := downloadDiff(rc, "abcd", ic.Options)
		if err != nil {
			t.Errorf("Failed to download layer: %s", err)
			return
		}
		t.Logf("layer id: %#v", layer)
		destDir := path.Join(imagec.DestinationDirectory(ic.Options), layer.ID)
		destination := path.Join(destDir, layer.ID+".tar")
		if _, err := os.Stat(destination); err != nil {
			t.Errorf("diff file %s is not created", destination)
		}
		os.RemoveAll(destDir)
	}
}
