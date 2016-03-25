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

package ph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

func TestImageAnonymization(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	phClient := NewClient("line001")
	image := NewImageEvent()
	image.ApplianceID = "vch1"
	image.Operation = "deployment"
	image.ParentID = "030120"
	image.ImageID = "030121"
	image.Registry = "localhost:5000/mytest"
	image.VMDKPath = "vm001/vm001.vmdk"
	phClient.AnonymizeImageEvent(image)
	if image.Registry != "http://334389048b872a533002b34d73f8c29fd09efc50:5000/mytest" {
		t.Errorf("Failed to anonymize image registry, was %s", image.Registry)
	}
	image.Registry = "https://localhost/mytest"
	phClient.AnonymizeImageEvent(image)
	if image.Registry != "https://334389048b872a533002b34d73f8c29fd09efc50/mytest" {
		t.Errorf("Failed to anonymize image registry, was %s", image.Registry)
	}
	image.Registry = "127.0.0.1:5000/mytest"
	image.VMDKPath = "vm001/vm001.vmdk"
	phClient.AnonymizeImageEvent(image)
	if image.Registry != "http://4b84b15bff6ee5796152495a230e45e3d7e947d9:5000/mytest" {
		t.Errorf("Failed to anonymize image registry, was %s", image.Registry)
	}
	image.Registry = "https://127.0.0.1/mytest"
	phClient.AnonymizeImageEvent(image)
	if image.Registry != "https://4b84b15bff6ee5796152495a230e45e3d7e947d9/mytest" {
		t.Errorf("Failed to anonymize image registry, was %s", image.Registry)
	}
}

func TestMultiPost(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	phClient := NewClient("line001")
	install := NewInstallationData()
	install.ApplianceID = "line001"
	install.ContainerPortgroupID = ""
	install.DatastoreID = ""
	install.DockerOpts = ""
	install.ExternalPortgroupID = ""
	install.FailedStep = "init"
	install.FinishTime = "10:10:10"
	install.Force = false
	install.Message = "failed to wait initialization status"
	install.Operation = "install"
	install.OS = "windows"
	install.Status = "failed"
	install.VCID = "vc123"

	array := make([]*InstallationData, 2)
	array[0] = install
	array[1] = install
	log.Debugf("object: %s", array)

	err := phClient.POST(*install)
	if err != nil {
		t.Errorf("failed to post json %s", err)
	}
}

func TestTimer(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	// Test server that always responds with 201 code, and specific payload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var postData []*ContainerEvent
		if err := json.Unmarshal(body, &postData); err != nil {
			t.Errorf("Posted data format is wrong")
		}
		log.Debugf("Posted data: %s", postData)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[{"test":"pass"}]`)
	}))
	defer server.Close()

	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	phClient := NewClient("line001")
	phClient.SetPHAddress("http://testaddress")
	phClient.Transport = transport
	container := NewContainerEvent()
	container.ApplianceID = "line001"
	container.ContainerID = "container1"
	container.CPU = 1
	container.EventTime = time.Now().UTC().String()
	container.ImageID = "image1"
	container.Memory = 2048
	container.Operation = "create"
	container.StartSeconds = 8
	phClient.AddContainerEvent(container)
	phClient.StartPOST(time.Millisecond*100, time.Millisecond*300)

	ticker := time.NewTicker(time.Millisecond * 50)
	go func() {
		for _ = range ticker.C {
			phClient.AddContainerEvent(container)
		}
	}()
	time.Sleep(time.Millisecond * 1000)
	ticker.Stop()
	phClient.stopPOST()
}

type MyConfigData struct {
	Type string `json:"@type"`
	CPU  int    `json:"cpu"`
}

type MyConfigHandler1 struct {
}

func (h *MyConfigHandler1) Name() string {
	return "Test1"
}

func (h *MyConfigHandler1) ConfigData() interface{} {
	data := &MyConfigData{
		Type: "myTable",
		CPU:  2,
	}
	return data
}

type MyConfigHandler2 struct {
}

func (h *MyConfigHandler2) Name() string {
	return "Test2"
}

func (h *MyConfigHandler2) ConfigData() interface{} {
	data := NewCRDProductInstance()
	data.Name = "VC1"
	return data
}

func TestConfigHandler(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	// Test server that always responds with 201 code, and specific payload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		log.Debugf("Posted data: %s", body)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[{"test":"pass"}]`)
	}))
	defer server.Close()

	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	phClient := NewClient("line001")
	phClient.SetPHAddress("http://testaddress")
	phClient.Transport = transport
	phClient.AddConfigHandler(&MyConfigHandler1{})
	phClient.AddConfigHandler(&MyConfigHandler2{})
	phClient.StartPOST(time.Millisecond*100, time.Millisecond*300)
	time.Sleep(time.Millisecond * 1000)
	phClient.stopPOST()
}
