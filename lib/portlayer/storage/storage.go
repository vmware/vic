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

package storage

import (
	"context"
	"io"
	"net/url"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

var (
	once sync.Once

	importers map[string]Importer
	exporters map[string]Exporter
)

func create(ctx context.Context, session *session.Session, pool *object.ResourcePool) error {
	var err error

	mngr := view.NewManager(session.Vim25())

	// Create view of VirtualMachine objects under the VCH's resource pool
	Config.ContainerView, err = mngr.CreateContainerView(ctx, pool.Reference(), []string{"VirtualMachine"}, true)
	if err != nil {
		return err
	}
	return nil
}

func Init(ctx context.Context, session *session.Session, pool *object.ResourcePool, source extraconfig.DataSource, _ extraconfig.DataSink) error {
	defer trace.End(trace.Begin(""))

	var err error

	once.Do(func() {
		// Grab the storage layer config blobs from extra config
		extraconfig.Decode(source, &Config)
		log.Debugf("Decoded VCH config for storage: %#v", Config)

		err = create(ctx, session, pool)

		importers = make(map[string]Importer)
		exporters = make(map[string]Exporter)
	})
	return err
}

func RegisterImporter(op trace.Operation, store string, i Importer) {
	op.Infof("Registering importer: %s => %T", store, i)

	importers[store] = i
}

func RegisterExporter(op trace.Operation, store string, e Exporter) {
	op.Infof("Registering exporter: %s => %T", store, e)

	exporters[store] = e
}

func GetImporter(store string) (Importer, bool) {
	i, ok := importers[store]
	return i, ok
}

func GetExporter(store string) (Exporter, bool) {
	e, ok := exporters[store]
	return e, ok
}

// Resolver defines methods for mapping ids to URLS, and urls to owners of that device
type Resolver interface {
	// URL returns a url to the data source representing `id`
	URL(op trace.Operation, id string) (*url.URL, error)
	// Owners returns a list of VMs that are using the resource specified by `url`
	Owners(op trace.Operation, url *url.URL, filter func(vm *mo.VirtualMachine) bool) ([]*vm.VirtualMachine, error)
}

// DataSource defines the methods for exporting data to/from a specific data source
// as a tar archive
type DataSource interface {
	io.Closer

	Export(op trace.Operation, spec *archive.FilterSpec, data bool) (io.ReadCloser, error)

	// Source returns the mechanism by which the data source is accessed
	// Examples:
	//     vmdk mounted locally: *os.File
	//     nfs volume:  XDR-client
	//     via guesttools:  tar stream
	Source() interface{}
}

type DataSink interface {
	io.Closer

	Import(op trace.Operation, spec *archive.FilterSpec, tarStream io.ReadCloser) error

	// Sink returns the mechanism by which the data sink is accessed
	// Examples:
	//     vmdk mounted locally: *os.File
	//     nfs volume:  XDR-client
	//     via guesttools:  tar stream
	Sink() interface{}
}

// Importer defines the methods needed to write data into a storage element
type Importer interface {
	Import(op trace.Operation, ID string, spec *archive.FilterSpec, tarStream io.ReadCloser) error
	NewDataSink(op trace.Operation, id string) (DataSink, error)
}

// Exporter defines the methods needed to read data from a storage element, optionally diff with an ancestor
type Exporter interface {
	Export(op trace.Operation, id, ancestor string, spec *archive.FilterSpec, data bool) (io.ReadCloser, error)
	NewDataSource(op trace.Operation, id string) (DataSource, error)
}
