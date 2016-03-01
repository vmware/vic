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

func TestMultiPost(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	phClient := NewClient("line001")
	install := NewInstallationData()
	install.ApplianceID = "line001"
	install.CIDR = "127.0.0.0/24"
	install.Cluster = "/Data Center/host/Cluster1"
	install.ContainerNetwork = "net1"
	install.ContainerPortgroupID = ""
	install.Datacenter = "DataCenter"
	install.Datastore = "datastore1"
	install.DatastoreID = ""
	install.DNS = ""
	install.DockerOpts = ""
	install.ExternalNetwork = "net2"
	install.ExternalPortgroupID = ""
	install.FailedStep = "init"
	install.FinishTime = "10:10:10"
	install.Force = false
	install.Host = "/DataCenter/host/Cluster1/host1"
	install.IP = "1.1.1.1"
	install.Message = "failed to wait initialization status"
	install.Name = "line-test1"
	install.Operation = "install"
	install.OS = "windows"
	install.Status = "failed"
	install.Target = "testvc"
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
	container.IP = ""
	container.Memory = 2048
	container.Name = "random"
	container.Operation = "create"
	container.PortMapping = "8080:8080"
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
