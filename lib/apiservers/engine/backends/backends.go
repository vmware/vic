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

package backends

import (
	"context"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/events"
	"github.com/docker/docker/registry"
	rc "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/engine/backends/container"
	apiclient "github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/misc"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/imagec"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/sys"
)

const (
	PortlayerName      = "Backend Engine"
	IndexServerAddress = "registry-1.docker.io"

	// RetryTimeSeconds defines how many seconds to wait between retries
	RetryTimeSeconds = 2
)

var (
	portLayerClient     *apiclient.PortLayer
	portLayerServerAddr string
	portLayerName       string
	productName         string
	productVersion      string

	vchConfig *config.VirtualContainerHostConfigSpec

	insecureRegistries []string
	RegistryCertPool   *x509.CertPool

	eventService *events.Events
)

func Init(portLayerAddr, product string, config *config.VirtualContainerHostConfigSpec, insecureRegs []url.URL) error {
	_, _, err := net.SplitHostPort(portLayerAddr)
	if err != nil {
		return err
	}

	vchConfig = config
	productName = product

	if config != nil {
		if config.Version != nil {
			productVersion = config.Version.ShortVersion()
		}
		if productVersion == "" {
			portLayerName = product + " Backend Engine"
		} else {
			portLayerName = product + " " + productVersion + " Backend Engine"
		}

		loadRegistryCACerts()
	} else {
		portLayerName = product + " Backend Engine"
	}

	t := rc.New(portLayerAddr, "/", []string{"http"})

	portLayerClient = apiclient.New(t, nil)
	portLayerServerAddr = portLayerAddr

	// block indefinitely while waiting on the portlayer to respond to pings
	// the vic-machine installer timeout will intervene if this blocks for too long
	pingPortLayer()

	if err := hydrateCaches(); err != nil {
		return err
	}

	log.Info("Creating image store")
	if err := createImageStore(); err != nil {
		log.Errorf("Failed to create image store")
		return err
	}

	serviceOptions := registry.ServiceOptions{}
	for _, r := range insecureRegs {
		// host:port was previously stored in the Path
		// component of the URL. Go 1.8 parser puts it
		// the Host component. VCH upgrade must move it
		// from Path to Host.
		insecureRegistries = append(insecureRegistries, r.Host)
	}
	if len(insecureRegistries) > 0 {
		serviceOptions.InsecureRegistries = insecureRegistries
	}

	eventService = events.New()

	return nil
}

func hydrateCaches() error {

	const waiters = 3

	wg := sync.WaitGroup{}
	wg.Add(waiters)
	errChan := make(chan error, waiters)

	go func() {
		defer wg.Done()
		if err := imagec.InitializeLayerCache(portLayerClient); err != nil {
			errChan <- fmt.Errorf("Failed to initialize layer cache: %s", err)
			return
		}
		log.Info("Layer cache initialized successfully")
		errChan <- nil
	}()

	go func() {
		defer wg.Done()
		if err := cache.InitializeImageCache(portLayerClient); err != nil {
			errChan <- fmt.Errorf("Failed to initialize image cache: %s", err)
			return
		}
		log.Info("Image cache initialized successfully")

		// container cache relies on image cache so we share a goroutine to update
		// them serially
		if err := syncContainerCache(); err != nil {
			errChan <- fmt.Errorf("Failed to update container cache: %s", err)
			return
		}
		log.Info("Container cache updated successfully")
		errChan <- nil
	}()

	go func() {
		log.Info("Refreshing repository cache")
		defer wg.Done()
		if err := cache.NewRepositoryCache(portLayerClient); err != nil {
			errChan <- fmt.Errorf("Failed to create repository cache: %s", err.Error())
			return
		}
		errChan <- nil
		log.Info("Repository cache updated successfully")
	}()

	wg.Wait()
	close(errChan)

	var errs []string
	for err := range errChan {
		if err != nil {
			// accumulate all errors into one
			errs = append(errs, err.Error())
		}
	}

	var e error
	if len(errs) > 0 {
		e = fmt.Errorf(strings.Join(errs, ", "))
	}
	return e
}

func PortLayerClient() *apiclient.PortLayer {
	return portLayerClient
}

func PortLayerServer() string {
	return portLayerServerAddr
}

func PortLayerName() string {
	return portLayerName
}

func ProductName() string {
	return productName
}

func ProductVersion() string {
	return productVersion
}

func VchConfig() *config.VirtualContainerHostConfigSpec {
	return vchConfig
}

func pingPortLayer() {
	ticker := time.NewTicker(RetryTimeSeconds * time.Second)
	defer ticker.Stop()
	params := misc.NewPingParamsWithContext(context.TODO())

	log.Infof("Waiting for portlayer to come up")

	for range ticker.C {
		if _, err := portLayerClient.Misc.Ping(params); err == nil {
			log.Info("Portlayer is up and responding to pings")
			return
		}
	}
}

func createImageStore() error {
	// TODO(jzt): we should move this to a utility package or something
	host, err := sys.UUID()
	if err != nil {
		log.Errorf("Failed to determine host UUID")
		return err
	}

	// attempt to create the image store if it doesn't exist
	store := &models.ImageStore{Name: host}
	_, err = portLayerClient.Storage.CreateImageStore(
		storage.NewCreateImageStoreParamsWithContext(ctx).WithBody(store),
	)

	if err != nil {
		if _, ok := err.(*storage.CreateImageStoreConflict); ok {
			log.Debugf("Store already exists")
			return nil
		}
		return err
	}
	log.Infof("Image store created successfully")
	return nil
}

func InsecureRegistries() []string {
	registries := []string{}
	for _, reg := range insecureRegistries {
		registries = append(registries, reg)
	}

	return registries
}

// syncContainerCache runs once at startup to populate the container cache
func syncContainerCache() error {
	log.Debugf("Updating container cache")

	backend := NewContainerBackend()
	client := backend.containerProxy.Client()

	reqParams := containers.NewGetContainerListParamsWithContext(ctx).WithAll(swag.Bool(true))
	containme, err := client.Containers.GetContainerList(reqParams)
	if err != nil {
		return errors.Errorf("Failed to retrieve container list from portlayer: %s", err)
	}

	log.Debugf("Found %d containers", len(containme.Payload))
	cc := cache.ContainerCache()
	var errs []string
	for _, info := range containme.Payload {
		container := ContainerInfoToVicContainer(*info)
		cc.AddContainer(container)
		if err = setPortMapping(info, backend, container); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.Errorf("Failed to set port mapping: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func setPortMapping(info *models.ContainerInfo, backend *Container, container *container.VicContainer) error {
	if info.ContainerConfig.State == "" {
		log.Infof("container state is nil")
		return nil
	}
	if info.ContainerConfig.State != "Running" || len(container.HostConfig.PortBindings) == 0 {
		log.Infof("Container info state: %s", info.ContainerConfig.State)
		log.Infof("container portbinding: %+v", container.HostConfig.PortBindings)
		return nil
	}

	log.Debugf("Set port mapping for container %q, portmapping %+v", container.Name, container.HostConfig.PortBindings)
	client := backend.containerProxy.Client()
	endpointsOK, err := client.Scopes.GetContainerEndpoints(
		scopes.NewGetContainerEndpointsParamsWithContext(ctx).WithHandleOrID(container.ContainerID))
	if err != nil {
		return err
	}
	for _, e := range endpointsOK.Payload {
		if len(e.Ports) > 0 {
			if err = MapPorts(container.HostConfig, e, container.ContainerID); err != nil {
				log.Errorf(err.Error())
				return err
			}
		}
	}
	return nil
}

func loadRegistryCACerts() {
	if len(vchConfig.RegistryCertificateAuthorities) == 0 {
		return
	}

	rootCertPool, err := x509.SystemCertPool()
	log.Debugf("Loaded %d CAs for registries from system CA bundle", len(rootCertPool.Subjects()))
	if err != nil {
		log.Errorf("Unable to load system CAs")
		return
	}

	if !rootCertPool.AppendCertsFromPEM(vchConfig.RegistryCertificateAuthorities) {
		log.Errorf("Unable to load CAs for registry access in config")
		return
	}

	RegistryCertPool = rootCertPool

	log.Debugf("Loaded %d CAs for registries from config", len(rootCertPool.Subjects()))
}

func EventService() *events.Events {
	return eventService
}
