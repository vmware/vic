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

package vicbackends

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/go-swagger/client"

	"github.com/docker/docker/api/types/backend"
	derr "github.com/docker/docker/errors"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	dnetwork "github.com/docker/engine-api/types/network"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	plscopes "github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	plmodels "github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/metadata"
)

//***********
// Mock proxy
//***********

type CreateHandleMockData struct {
	createInputID   string
	retID           string
	retHandle       string
	retErr          error
	createErrSubstr string
}

type AddToScopeMockData struct {
	createInputID   string
	retHandle       string
	retErr          error
	createErrSubstr string
}

type AddVolumesMockData struct {
	retHandle       string
	retErr          error
	createErrSubstr string
}

type CommitHandleMockData struct {
	createInputID   string
	createErrSubstr string

	retErr error
}

type LogMockData struct {
	continaerID string
	running     bool
}

type MockContainerProxy struct {
	mockRespIndices      []int
	mockCreateHandleData []CreateHandleMockData
	mockAddToScopeData   []AddToScopeMockData
	mockAddVolumesData   []AddVolumesMockData
	mockCommitData       []CommitHandleMockData
}

const (
	SUCCESS              = 0
	dummyContainerID     = "abc123"
	dummyContainerID_tty = "tty123"
)

var dummyContainers []string = []string{dummyContainerID, dummyContainerID_tty}

func NewMockContainerProxy() *MockContainerProxy {
	return &MockContainerProxy{
		mockRespIndices:      make([]int, 4),
		mockCreateHandleData: MockCreateHandleData(),
		mockAddToScopeData:   MockAddToScopeData(),
		mockAddVolumesData:   MockAddVolumesData(),
		mockCommitData:       MockCommitData(),
	}
}

func MockCreateHandleData() []CreateHandleMockData {
	createHandleTimeoutErr := client.NewAPIError("unknown error", "context deadline exceeded", http.StatusServiceUnavailable)

	mockCreateHandleData := []CreateHandleMockData{
		{"busybox", "321cba", "handle", nil, ""},
		{"busybox", "", "", derr.NewRequestNotFoundError(fmt.Errorf("No such image: abc123")), "No such image"},
		{"busybox", "", "", derr.NewErrorWithStatusCode(createHandleTimeoutErr, http.StatusInternalServerError), "context deadline exceeded"},
	}

	return mockCreateHandleData
}

func MockAddToScopeData() []AddToScopeMockData {
	addToScopeNotFound := plscopes.AddContainerNotFound{
		Payload: &plmodels.Error{
			Message: "Scope not found",
		},
	}

	addToScopeNotFoundErr := fmt.Errorf("ContainerProxy.AddContainerToScope: Scopes error: %s", addToScopeNotFound.Error())

	addToScopeTimeout := plscopes.AddContainerInternalServerError{
		Payload: &plmodels.Error{
			Message: "context deadline exceeded",
		},
	}

	addToScopeTimeoutErr := fmt.Errorf("ContainerProxy.AddContainerToScope: Scopes error: %s", addToScopeTimeout.Error())

	mockAddToScopeData := []AddToScopeMockData{
		{"busybox", "handle", nil, ""},
		{"busybox", "handle", derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"), http.StatusInternalServerError), "failed to create a portlayer"},
		{"busybox", "handle", derr.NewErrorWithStatusCode(addToScopeNotFoundErr, http.StatusInternalServerError), "Scope not found"},
		{"busybox", "handle", derr.NewErrorWithStatusCode(addToScopeTimeoutErr, http.StatusInternalServerError), "context deadline exceeded"},
	}

	return mockAddToScopeData
}

func MockAddVolumesData() []AddVolumesMockData {
	return nil
}

func MockCommitData() []CommitHandleMockData {
	noSuchImageErr := fmt.Errorf("No such image: busybox")

	mockCommitData := []CommitHandleMockData{
		{"buxybox", "", nil},
		{"busybox", "failed to create a portlayer", derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"), http.StatusInternalServerError)},
		{"busybox", "No such image", derr.NewRequestNotFoundError(noSuchImageErr)},
	}

	return mockCommitData
}

func (m *MockContainerProxy) GetMockDataCount() (int, int, int, int) {
	return len(m.mockCreateHandleData), len(m.mockAddToScopeData), len(m.mockAddVolumesData), len(m.mockCommitData)
}

func (m *MockContainerProxy) SetMockDataResponse(createHandleResp int, addToScopeResp int, addVolumeResp int, commitContainerResp int) {
	m.mockRespIndices[0] = createHandleResp
	m.mockRespIndices[1] = addToScopeResp
	m.mockRespIndices[2] = addVolumeResp
	m.mockRespIndices[3] = commitContainerResp
}

func (m *MockContainerProxy) CreateContainerHandle(imageID string, config types.ContainerCreateConfig) (string, string, error) {
	respIdx := m.mockRespIndices[0]

	if respIdx >= len(m.mockCreateHandleData) {
		return "", "", nil
	}
	return m.mockCreateHandleData[respIdx].retID, m.mockCreateHandleData[respIdx].retHandle, m.mockCreateHandleData[respIdx].retErr
}

func (m *MockContainerProxy) AddContainerToScope(handle string, config types.ContainerCreateConfig) (string, error) {
	respIdx := m.mockRespIndices[1]

	if respIdx >= len(m.mockAddToScopeData) {
		return "", nil
	}

	return m.mockAddToScopeData[respIdx].retHandle, m.mockAddToScopeData[respIdx].retErr
}

func (m *MockContainerProxy) AddVolumesToContainer(handle string, config types.ContainerCreateConfig) (string, error) {
	respIdx := m.mockRespIndices[2]

	if respIdx >= len(m.mockAddVolumesData) {
		return "", nil
	}

	return m.mockAddVolumesData[respIdx].retHandle, m.mockAddVolumesData[respIdx].retErr
}

func (m *MockContainerProxy) CommitContainerHandle(handle, imageID string) error {
	respIdx := m.mockRespIndices[3]

	if respIdx >= len(m.mockCommitData) {
		return nil
	}

	return m.mockCommitData[respIdx].retErr
}

func (m *MockContainerProxy) StreamContainerLogs(name string, out io.Writer, started chan struct{}, showTimestamps bool, followLogs bool, since int64, tailLines int64) error {
	var lineCount int64 = 10

	close(started)

	for i := int64(0); i < lineCount; i++ {
		if !followLogs && i > tailLines {
			break
		}
		if followLogs && i > tailLines {
			time.Sleep(500 * time.Millisecond)
		}

		fmt.Fprintf(out, "line %d\n", i)
	}

	return nil
}

func (m *MockContainerProxy) ContainerRunning(vc *viccontainer.VicContainer) (bool, error) {
	// Assume container is running if container in cache.  If we need other conditions
	// in the future, we can add it, but for now, just assume running.
	container := cache.ContainerCache().GetContainer(vc.ContainerID)

	if container == nil {
		return false, nil
	}

	return true, nil
}

func AddMockImageToCache() {
	mockImage := &metadata.ImageConfig{
		ImageID: "e732471cb81a564575aad46b9510161c5945deaf18e9be3db344333d72f0b4b2",
		Name:    "busybox",
		Tags:    []string{"latest"},
	}
	mockImage.Config = &container.Config{
		Hostname:     "55cd1f8f6e5b",
		Domainname:   "",
		User:         "",
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Env:          []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		Cmd:          []string{"sh"},
		Image:        "sha256:e732471cb81a564575aad46b9510161c5945deaf18e9be3db344333d72f0b4b2",
		Volumes:      nil,
		WorkingDir:   "",
		Entrypoint:   nil,
		OnBuild:      nil,
	}

	cache.ImageCache().AddImage(mockImage)
}

func AddMockContainerToCache() {
	AddMockImageToCache()

	image, err := cache.ImageCache().GetImage("e732471cb81a564575aad46b9510161c5945deaf18e9be3db344333d72f0b4b2")
	if err == nil {
		vc := viccontainer.NewVicContainer()
		vc.ImageID = image.ID
		vc.Config = image.Config //Set defaults.  Overrides will get copied below.
		vc.Config.Tty = false
		vc.ContainerID = dummyContainerID

		cache.ContainerCache().AddContainer(vc)

		vc = viccontainer.NewVicContainer()
		vc.ImageID = image.ID
		vc.Config = image.Config
		vc.Config.Tty = true
		vc.ContainerID = dummyContainerID_tty

		cache.ContainerCache().AddContainer(vc)
	}
}

//***********
// Tests
//***********

// TestContainerCreateEmptyImageCache() attempts a ContainerCreate() with an empty image
// cache
func TestContainerCreateEmptyImageCache(t *testing.T) {
	mockContainerProxy := NewMockContainerProxy()

	// Create our personality Container backend
	cb := &Container{
		containerProxy: mockContainerProxy,
	}

	// mock a container create config
	var config types.ContainerCreateConfig

	config.HostConfig = &container.HostConfig{}
	config.Config = &container.Config{}
	config.NetworkingConfig = &dnetwork.NetworkingConfig{}
	config.Config.Image = "busybox"

	_, err := cb.ContainerCreate(config)

	assert.Contains(t, err.Error(), "No such image", "Error (%s) should have 'No such image' for an empty image cache", err.Error())
}

// TestCreateHandle() cycles through all possible input/outputs for creating a handle
// and calls vicbackends.ContainerCreate().  The idea is that if creating handle fails
// then vicbackends.ContainerCreate() should return errors from that.
func TestCreateHandle(t *testing.T) {
	mockContainerProxy := NewMockContainerProxy()

	// Create our personality Container backend
	cb := &Container{
		containerProxy: mockContainerProxy,
	}

	AddMockImageToCache()

	// mock a container create config
	var config types.ContainerCreateConfig

	config.HostConfig = &container.HostConfig{}
	config.Config = &container.Config{}
	config.NetworkingConfig = &dnetwork.NetworkingConfig{}

	mockCreateHandleData := MockCreateHandleData()

	// Iterate over create handler responses and see what the composite ContainerCreate()
	// returns.  Since the handle is the first operation, we expect to receive a create handle
	// error.
	count, _, _, _ := mockContainerProxy.GetMockDataCount()

	for i := 0; i < count; i++ {
		if i == SUCCESS { //skip success case
			continue
		}

		mockContainerProxy.SetMockDataResponse(i, 0, 0, 0)
		config.Config.Image = mockCreateHandleData[i].createInputID
		_, err := cb.ContainerCreate(config)

		assert.Contains(t, err.Error(), mockCreateHandleData[i].createErrSubstr)
	}
}

// TestContainerAddToScope() assumes container handle create succeeded and cycles through all
// possible input/outputs for adding container to scope and calls vicbackends.ContainerCreate()
func TestContainerAddToScope(t *testing.T) {
	mockContainerProxy := NewMockContainerProxy()

	// Create our personality Container backend
	cb := &Container{
		containerProxy: mockContainerProxy,
	}

	AddMockImageToCache()

	// mock a container create config
	var config types.ContainerCreateConfig

	config.HostConfig = &container.HostConfig{}
	config.Config = &container.Config{}
	config.NetworkingConfig = &dnetwork.NetworkingConfig{}

	mockAddToScopeData := MockAddToScopeData()

	// Iterate over create handler responses and see what the composite ContainerCreate()
	// returns.  Since the handle is the first operation, we expect to receive a create handle
	// error.
	_, count, _, _ := mockContainerProxy.GetMockDataCount()

	for i := 0; i < count; i++ {
		if i == SUCCESS { //skip success case
			continue
		}

		mockContainerProxy.SetMockDataResponse(0, i, 0, 0)
		config.Config.Image = mockAddToScopeData[i].createInputID
		_, err := cb.ContainerCreate(config)

		assert.Contains(t, err.Error(), mockAddToScopeData[i].createErrSubstr)
	}
}

// TestContainerAddVolumes() assumes container handle create succeeded and cycles through all
// possible input/outputs for committing the handle and calls vicbackends.ContainerCreate()
func TestCommitHandle(t *testing.T) {
	mockContainerProxy := NewMockContainerProxy()

	// Create our personality Container backend
	cb := &Container{
		containerProxy: mockContainerProxy,
	}

	AddMockImageToCache()

	// mock a container create config
	var config types.ContainerCreateConfig

	config.HostConfig = &container.HostConfig{}
	config.Config = &container.Config{}
	config.NetworkingConfig = &dnetwork.NetworkingConfig{}

	mockCommitHandleData := MockCommitData()

	// Iterate over create handler responses and see what the composite ContainerCreate()
	// returns.  Since the handle is the first operation, we expect to receive a create handle
	// error.
	_, _, _, count := mockContainerProxy.GetMockDataCount()

	for i := 0; i < count; i++ {
		if i == SUCCESS { //skip success case
			continue
		}

		mockContainerProxy.SetMockDataResponse(0, 0, 0, i)
		config.Config.Image = mockCommitHandleData[i].createInputID
		_, err := cb.ContainerCreate(config)

		assert.Contains(t, err.Error(), mockCommitHandleData[i].createErrSubstr)
	}

}

// TestContainerLogs() tests the docker logs api when user asks for entire log
func TestContainerLogs(t *testing.T) {
	mockContainerProxy := NewMockContainerProxy()

	// Create our personality Container backend
	cb := &Container{
		containerProxy: mockContainerProxy,
	}

	// Prepopulate our image and container cache with dummy data
	AddMockContainerToCache()

	// Create a buffer io.writer
	var writer bytes.Buffer

	successDuration := 1 * time.Second

	// Create our mock table
	mockData := []struct {
		Config          backend.ContainerLogsConfig
		ExpectedSuccess bool
		ExpectedFollow  bool
	}{
		{
			Config: backend.ContainerLogsConfig{
				ContainerLogsOptions: types.ContainerLogsOptions{
					ShowStdout: true,
					ShowStderr: true,
					Tail:       "all",
				},
				OutStream: &writer,
			},
			ExpectedSuccess: true,
			ExpectedFollow:  false,
		},
		{
			Config: backend.ContainerLogsConfig{
				ContainerLogsOptions: types.ContainerLogsOptions{
					ShowStdout: false,
					ShowStderr: false,
				},
				OutStream: &writer,
			},
			ExpectedSuccess: false,
			ExpectedFollow:  false,
		},
		{
			Config: backend.ContainerLogsConfig{
				ContainerLogsOptions: types.ContainerLogsOptions{
					ShowStdout: true,
					ShowStderr: true,
					Follow:     true,
				},
				OutStream: &writer,
			},
			ExpectedSuccess: true,
			ExpectedFollow:  true,
		},
	}

	for _, containerID := range dummyContainers {
		for _, data := range mockData {
			started := make(chan struct{})

			start := time.Now()
			err := cb.ContainerLogs(containerID, &data.Config, started)
			end := time.Now()

			select {
			case <-started:
			default:
				close(started)
			}

			if data.ExpectedSuccess {
				assert.Nil(t, err, "Expected success, but got error, config: %#v", data.Config)

			} else {
				assert.NotEqual(t, err, nil, "Expected error but received nil, config: %#v", data.Config)
			}

			immediate := start.Add(successDuration)

			didFollow := immediate.Before(end) //determines if logs continued to stream

			if data.ExpectedFollow {
				assert.True(t, didFollow, "Expected logs to follow but didn't (%s, %s), config: %#v", start.String(), end.String(), data.Config)
			} else {
				assert.False(t, didFollow, "Expected logs to NOT follow but it did, config: %#v", data.Config)
			}
		}
	}
}
