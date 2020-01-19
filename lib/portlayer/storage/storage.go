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

// Package storage provides abstractions for interacting with storage locations, both management and access.
// There are two primary type of storage handled by this package at this time:
// a. VMDKs on a vSphere datastore
// b. NFS shares
//
// We have more than two sub-packages however as we use distinct packages to provide different lifecycle and
// usage semantics for types of data using those mechanisms. Examples include:
// * application logic (images)
// * application transient data (image read/write layer)
// * application persistent data (volumes)
//
// As such there are a set of concepts, with all but `store` captured as interfaces:
// * store
// * resolver
// * exporter
// * importer
// * data source
// * data sink
//
// The odd one out for these concepts is the `store`. This is because we did not want to be prescriptive about
// what a store supports, but it may well distill into an interface at a later date. Currently the various
// store implementations consist of a composition of `resolver`, `importer`, and `exporter` interfaces.
//
// The `data source` and `data sink` interfaces are present to allow abstracted read and write access respectively
// to a specific storage mechanism. The data source and data sink instances should be considered single use and not
// reused for multiple access. This semantic was chosen as it's the least prescriptive on storage availability and
// connectivity - the precise means used to read or write from storage may change, for example, based on whether
// the storage is in use (has an owner).
// The `resolver` interface is present to translate a potentially store specific ID into an unambiguous URI, and to
// allow discovery of which virtual machines, if any, currently claim that ID. This returns virtual machines instead
// of a further abstracted interface because it's not clear at time of writing what that abstraction should or would
// be.
// The `importer` and `exporter` interfaces are the factories used to acquire `data source` and `data sink` instances.
// They also provide convenience methods for directly creating and consuming the source/sink in one operation.
//
// It is not required that all storage elements implement both Importer and Exporter mechanisms - having write only or
// read only storage is viable.
//
// The data source and sink interfaces also specify a basic accessor method for the underlying mechanism. This was
// added to allow for some future-proofing, with the primary extension expected to be a generic FileWalker interface
// so that directory walk and listings can be done without needing to care about the access mechanism (e.g. common
// between a local filesystem mount, an XDR NFS client, or a VM toolbox client).
//
// The general structure of a specific storage implementation is as follows and is imposed so that it is easier to
// locate a specific functional path for a given usage/storage type combination. As of this writing there are two
// additional folders in this structure, nfs and vsphere. These contain implementation specific logic used across
// store types and must be at this level to avoid cyclic package dependencies:
// ```
// storage/                  - this package
// ├── [type]/               - the high level type of use, implies semantics (i.e., container, image, volume)
// │   ├── [type].go         - the common store implementation aspects across implementations
// │   ├── errors.go         - [type] specific error types
// │   ├── cache             - a `store` compliant cache implementation, if needed, with appropriate semantics for the type of use
// │   │   └── cache.go
// │   └── [implementation]/ - a specific implementation type (e.g., vsphere)
// │       ├── bind.go       - modifies a portlayer handle to configure active use of a specific `joined` storage
// │       ├── export.go     - implements the `read` side of the interfaces (Exporter and DataSource)
// │       ├── import.go     - a specific implementation type
// │       ├── join.go       - modifies a portlayer handle to configure basic access to a specific instance of the storage type
// │       └── store.go      - constructor and implementation for specific store type and implementation
// ├── storage.go            - interface definitions, portlayer lifecycle management functions, and the registration/lookup
// |                           mechanisms for store instances.
// └── [purpose].go          - common functions used by the various type/implementation combinations
// ```
//
// This structure is not completely consistent at the time of writing. Most notable is that the bind.go functions
// have all been rolled into the Join calls (the portlayer uses the following common verbs across components, again
// with some inconsistencies at this time: Join, Bind, Inspect, Unbind, Remove, Wait).
// In the case of the `container` store type there is only one implementation at this time (vsphere VMDKs) and the
// `implementation` subdirectory has been omitted.
//
// This file implements lookup of store implementations via the `RegisterImporter` and `RegisterExporter` functions.
// This is provided to allow a common pattern for implementing, finding, and accessing store implementations without
// the caller requiring specific knowledge about the type of a specific store. All that's required is knowledge of
// the set of `store` interfaces and the identifier with which a given store was implemented.
package storage

import (
	"context"
	"io"
	"net/url"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

var (
	once sync.Once

	importers map[string]Importer
	exporters map[string]Exporter
)

// Resolver defines methods for mapping ids to URLS, and urls to owners of that device
type Resolver interface {
	// URL returns a url to the data source representing `id`
	// For historic reasons this is not the same URL that other parts of the storage component use, but an actual
	// URL suited for locating the storage element without having additional precursor knowledge.
	URL(op trace.Operation, id string) (*url.URL, error)

	// Owners returns a list of VMs that are using the resource specified by `url`
	Owners(op trace.Operation, url *url.URL, filter func(vm *mo.VirtualMachine) bool) ([]*vm.VirtualMachine, error)
}

// DataSource defines the methods for exporting data from a specific storage element as a tar stream
type DataSource interface {
	// Close releases all resources associated with this source. Shared resources should be reference counted.
	io.Closer

	// Export performs an export of the specified files, returning the data as a tar stream. This is single use; once
	// the export has completed it should not be assumed that the source remains functional.
	//
	// spec: specifies which files will be included/excluded in the export and allows for path rebasing/stripping
	// data: if true the actual file data is included, if false only the file headers are present
	Export(op trace.Operation, spec *archive.FilterSpec, data bool) (io.ReadCloser, error)

	// Source returns the mechanism by which the data source is accessed
	// Examples:
	//     vmdk mounted locally: *os.File
	//     nfs volume:  		 XDR-client
	//     via guesttools:  	 toolbox client
	Source() interface{}

	// Stat stats the filesystem target indicated by the last entry in the given Filterspecs inclusion map
	Stat(op trace.Operation, spec *archive.FilterSpec) (*FileStat, error)
}

// DataSink defines the methods for importing data to a specific storage element from a tar stream
type DataSink interface {
	// Close releases all resources associated with this sink. Shared resources should be reference counted.
	io.Closer

	// Import performs an import of the tar stream to the source held by this DataSink.  This is single use; once
	// the export has completed it should not be assumed that the sink remains functional.
	//
	// spec: specifies which files will be included/excluded in the import and allows for path rebasing/stripping
	// tarStream: the tar stream to from which to import data
	Import(op trace.Operation, spec *archive.FilterSpec, tarStream io.ReadCloser) error

	// Sink returns the mechanism by which the data sink is accessed
	// Examples:
	//     vmdk mounted locally: *os.File
	//     nfs volume:  		 XDR-client
	//     via guesttools:  	 toolbox client
	Sink() interface{}
}

// Importer defines the methods needed to write data into a storage element. This should be implemented by the various
// store types.
type Importer interface {
	// Import allows direct construction and invocation of a data sink for the specified ID.
	Import(op trace.Operation, id string, spec *archive.FilterSpec, tarStream io.ReadCloser) error

	// NewDataSink constructs a data sink for the specified ID within the context of the Importer. This is a single
	// use sink which may hold resources until Closed.
	NewDataSink(op trace.Operation, id string) (DataSink, error)
}

// Exporter defines the methods needed to read data from a storage element, optionally diff with an ancestor. This
// shoiuld be implemented by the various store types.
type Exporter interface {
	// Export allows direct construction and invocation of a data source for the specified ID.
	Export(op trace.Operation, id, ancestor string, spec *archive.FilterSpec, data bool) (io.ReadCloser, error)

	// NewDataSource constructs a data source for the specified ID within the context of the Exporter. This is a single
	// use source which may hold resources until Closed.
	NewDataSource(op trace.Operation, id string) (DataSource, error)
}

type FileStat struct {
	LinkTarget string
	Mode       uint32
	Name       string
	Size       int64
	ModTime    time.Time
}

func init() {
	importers = make(map[string]Importer)
	exporters = make(map[string]Exporter)
}

func create(ctx context.Context, session *session.Session, pool *object.ResourcePool) error {
	var err error

	mngr := view.NewManager(session.Vim25())

	op := trace.FromContext(ctx, "storage component initialization")

	// Create view of VirtualMachine objects under the VCH's resource pool
	Config.ContainerView, err = mngr.CreateContainerView(op, pool.Reference(), []string{"VirtualMachine"}, true)
	if err != nil {
		return err
	}

	Config.DiskManager, err = disk.NewDiskManager(op, session, Config.ContainerView)
	if err != nil {
		return err
	}

	return nil
}

// Init performs basic initialization, including population of storage.Config
func Init(ctx context.Context, session *session.Session, pool *object.ResourcePool, source extraconfig.DataSource, _ extraconfig.DataSink) error {
	defer trace.End(trace.Begin(""))

	var err error

	once.Do(func() {
		// Grab the storage layer config blobs from extra config
		extraconfig.Decode(source, &Config)
		log.Debugf("Decoded VCH config for storage: %#v", Config)

		err = create(ctx, session, pool)
	})
	return err
}

// TODO: figure out why the Init calls are wrapped in once.Do - implies it can be called
// multiple times, but once Finalize is called things will not be functional.
func Finalize(ctx context.Context) error {
	if Config.ContainerView != nil {
		Config.ContainerView.Destroy(ctx)
	}

	return nil
}

// RegisterImporter registers the specified importer against the provided store for later retrieval.
func RegisterImporter(op trace.Operation, store string, i Importer) {
	op.Infof("Registering importer: %s => %T", store, i)

	importers[store] = i
}

// RegisterExporter registers the specified exporter against the provided store for later retrieval.
func RegisterExporter(op trace.Operation, store string, e Exporter) {
	op.Infof("Registering exporter: %s => %T", store, e)

	exporters[store] = e
}

// GetImporter retrieves an importer registered with the provided store.
// Will return nil, false if the store is not found.
func GetImporter(store string) (Importer, bool) {
	i, ok := importers[store]
	return i, ok
}

// GetExporter retrieves an exporter registered with the provided store.
// Will return nil, false if the store is not found.
func GetExporter(store string) (Exporter, bool) {
	e, ok := exporters[store]
	return e, ok
}

// GetImporters returns the set of known importers.
func GetImporters() []string {
	keys := make([]string, 0, len(importers))
	for key := range importers {
		keys = append(keys, key)
	}

	return keys
}

// GetExporters returns the set of known importers.
func GetExporters() []string {
	keys := make([]string, 0, len(exporters))
	for key := range exporters {
		keys = append(keys, key)
	}

	return keys
}
