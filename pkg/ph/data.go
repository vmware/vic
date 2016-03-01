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
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	containerEvents = []*ContainerEvent{}
	containerMutex  = &sync.Mutex{}

	imageEvents = []*ImageEvent{}
	imageMutex  = &sync.Mutex{}

	configHandler = []ConfigHandler{}
	configMutex   = &sync.Mutex{}
)

// InstallationData is designed for VCH install/uninstall information report
type InstallationData struct {
	// name: vic_deployment
	Type string `json:"@type"`
	// bonneville build version
	Version string `json:"version"`

	// command params
	Operation        string `json:"operation"`
	Name             string `json:"name"`
	OS               string `json:"os"`
	CIDR             string `json:"cidr"`
	Cluster          string `json:"cluster"`
	ContainerNetwork string `json:"container_network"`
	Datacenter       string `json:"datacenter"`
	Datastore        string `json:"datastore"`
	DNS              string `json:"dns"`
	DockerOpts       string `json:"docker_opts"`
	ExternalNetwork  string `json:"external_network"`
	Force            bool   `json:"force"`
	Host             string `json:"host"`
	IP               string `json:"ip"`
	MemoryMB         int    `json:"memory_mb"`
	NumCPUs          int    `json:"num_cpus"`
	Pool             string `json:"pool"`
	Target           string `json:"target"`
	Timeout          int    `json:"timeout"`
	Verify           bool   `json:"verify"`

	VCID                 string `json:"vc_id"`
	ContainerPortgroupID string `json:"container_pg_id"`
	ExternalPortgroupID  string `json:"external_pg_id"`
	HostID               string `json:"host_id"`
	DatastoreID          string `json:"datastore_id"`
	ApplianceID          string `json:"appliance_id"`

	Status     string `json:"status"`
	Message    string `json:"message"`
	FailedStep string `json:"failed_step"`
	// start timestamp
	StartTime string `json:"start"`
	// stop timestamp
	FinishTime string `json:"finish"`
}

// CRDProductInstance information is common data required by Phone Home
type CRDProductInstance struct {
	// table name: crd_product_instance
	Type    string `json:"@type"`
	ID      string `json:"@id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Edition string `json:"edition"`
}

// InitData is designed for VCH service initialization information report.
// The details of service initialization might not be available from installer
type InitData struct {
	// name: vic_appliance
	Type                 string `json:"@type"`
	ApplianceID          string `json:"appliance_id"`
	VCID                 string `json:"vc_id"`
	ContainerPortgroupID string `json:"container_pg_id"`
	ExternalPortgroupID  string `json:"external_pg_id"`
	HostID               string `json:"host_id"`
	DatastoreID          string `json:"datastore_id"`
	MemoryMB             int    `json:"memory_mb"`
	NumCPUs              int    `json:"num_cpus"`
	IP                   string `json:"ip"`
	InitStatus           string `json:"init_status"`
	FailureMessage       string `json:"failure_message"`
	ReportTime           string `json:"report_time"`
}

// ContainerEvent is to record container configuration and operations
type ContainerEvent struct {
	// name: vic_container
	Type        string `json:"@type"`
	ApplianceID string `json:"appliance_id"`
	ContainerID string `json:"container_id"`
	ImageID     string `json:"image_id"`
	EventTime   string `json:"event_time"`
	// run, exec, stop, rm, etc.
	Operation   string `json:"operation"`
	Name        string `json:"name"`
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	PortMapping string `json:"port_mapping"`
	IP          string `json:"ip"`
	// the seconds to start a container if available
	StartSeconds int `json:"start_seconds"`
}

// ImageEvent is to record image related operations and configuraiton
type ImageEvent struct {
	// name: vic_image
	Type        string `json:"@type"`
	ApplianceID string `json:"appliance_id"`
	ImageID     string `json:"image_id"`
	Name        string `json:"name"`
	Registry    string `json:"registry"`
	// pull/push/commit, etc
	Operation string `json:"operation"`
	VMDKPath  string `json:"vmdk_path"`
	ParentID  string `json:"parent_id"`
}

type ConfigHandler interface {
	Name() string
	ConfigData() interface{}
}

func NewInstallationData() *InstallationData {
	return &InstallationData{
		Type: "vic_deployment",
	}
}

func NewCRDProductInstance() *CRDProductInstance {
	return &CRDProductInstance{
		Type: "crd_product_instance",
	}
}

func NewInitData() *InitData {
	return &InitData{
		Type: "vic_appliance",
	}
}

func NewContainerEvent() *ContainerEvent {
	return &ContainerEvent{
		Type: "vic_container",
	}
}

func NewImageEvent() *ImageEvent {
	return &ImageEvent{
		Type: "vic_image",
	}
}

func (phc *Client) postImages() error {
	imageMutex.Lock()
	if len(imageEvents) == 0 {
		log.Debugf("No image events.")
		imageMutex.Unlock()
		return nil
	}

	images := imageEvents
	imageEvents = []*ImageEvent{}
	imageMutex.Unlock()

	if err := phc.POST(images); err != nil {
		log.Debugf("Failed report image events to VMware PhoneHome for %s.", err)
		imageMutex.Lock()
		imagesNew := []*ImageEvent{}
		imagesNew = append(imagesNew, images...)
		imageEvents = append(imagesNew, imageEvents...)
		imageMutex.Unlock()
		return err
	}
	log.Debugf("Posted images %s", images)
	return nil
}

func (phc *Client) postContainers() error {
	containerMutex.Lock()
	if len(containerEvents) == 0 {
		log.Debugf("No container events.")
		containerMutex.Unlock()
		return nil
	}

	containers := containerEvents
	containerEvents = []*ContainerEvent{}
	containerMutex.Unlock()

	if err := phc.POST(containers); err != nil {
		log.Debugf("Failed report container events to VMware PhoneHome for %s", err)
		containerMutex.Lock()
		containersNew := []*ContainerEvent{}
		containersNew = append(containersNew, containers...)
		containerEvents = append(containersNew, containerEvents...)
		containerMutex.Unlock()
		return err
	}
	log.Debugf("Posted container %s", containers)
	return nil
}

func (phc *Client) postConfigurations() error {
	configMutex.Lock()
	handlers := configHandler
	configHandler = []ConfigHandler{}
	configMutex.Unlock()

	for i := range handlers {
		data := handlers[i].ConfigData()
		if err := phc.POST(data); err != nil {
			log.Debugf("Failed report configuration data for handler %s", handlers[i].Name())
		}
		log.Debugf("Posted config %s", data)
	}
	configMutex.Lock()
	handlers = append(handlers, configHandler...)
	configHandler = handlers
	configMutex.Unlock()
	return nil
}

// StartPOST periodically post events and configurations to Phone Home
func (phc *Client) StartPOST(event time.Duration, conf time.Duration) {
	go func() {
		eventChan := time.NewTicker(event).C
		confChan := time.NewTicker(conf).C
		phc.doneChan = make(chan bool)

		for {
			select {
			case <-eventChan:
				log.Debugf("Report events to Phone Home.")
				go phc.postImages()
				go phc.postContainers()
			case <-confChan:
				log.Debugf("Report configuration to Phone Home.")
				go phc.postConfigurations()
			case <-phc.doneChan:
				log.Debugf("Stop post.")
				return
			}
		}
	}()
}

func (phc *Client) stopPOST() {
	phc.doneChan <- true
}

func (phc *Client) AddImageEvent(image *ImageEvent) {
	imageMutex.Lock()
	imageEvents = append(imageEvents, image)
	imageMutex.Unlock()
}

func (phc *Client) AddContainerEvent(container *ContainerEvent) {
	containerMutex.Lock()
	containerEvents = append(containerEvents, container)
	containerMutex.Unlock()
}

func (phc *Client) AddConfigHandler(handler ConfigHandler) {
	configMutex.Lock()
	configHandler = append(configHandler, handler)
	configMutex.Unlock()
}
