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

package proxy

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	vicarchive "github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/pkg/trace"
)

type VicArchiveProxy interface {
	ArchiveExportReader(ctx context.Context, store, ancestorStore, deviceID, ancestor string, data bool, filterSpec vicarchive.FilterSpec) (io.ReadCloser, error)
	ArchiveImportWriter(ctx context.Context, store, deviceID string, filterSpec vicarchive.FilterSpec) (io.WriteCloser, error)
}

//------------------------------------
// ArchiveProxy
//------------------------------------

type ArchiveProxy struct {
	client *client.PortLayer
}

func NewArchiveProxy(client *client.PortLayer) VicArchiveProxy {
	return &ArchiveProxy{client: client}
}

// ArchiveExportReader streams a tar archive from the portlayer.  Once the stream is complete,
// an io.Reader is returned and the caller can use that reader to parse the data.
func (a *ArchiveProxy) ArchiveExportReader(ctx context.Context, store, ancestorStore, deviceID, ancestor string, data bool, filterSpec vicarchive.FilterSpec) (io.ReadCloser, error) {
	defer trace.End(trace.Begin(deviceID))

	if store == "" || deviceID == "" {
		return nil, fmt.Errorf("ArchiveExportReader called with either empty store or deviceID")
	}

	var err error

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		// make sure we get out of io.Copy if context is canceled
		select {
		case <-ctx.Done():
			// Attempt to tell the portlayer to cancel the stream.  This is one way of cancelling the
			// stream.  The other way is for the caller of this function to close the returned CloseReader.
			// Callers of this function should do one but not both.
			err := pipeReader.Close()
			if err != nil {
				log.Errorf("Error closing pipereader in ArchiveExportReader: %s", err.Error())
			}
		}
	}()

	go func() {
		params := storage.NewExportArchiveParamsWithContext(ctx).
			WithStore(store).
			WithAncestorStore(&ancestorStore).
			WithDeviceID(deviceID).
			WithAncestor(&ancestor).
			WithData(data)

		// Encode the filter spec
		encodedFilter := ""
		if valueBytes, merr := json.Marshal(filterSpec); merr == nil {
			encodedFilter = base64.StdEncoding.EncodeToString(valueBytes)
			params = params.WithFilterSpec(&encodedFilter)
			log.Infof(" encodedFilter = %s", encodedFilter)
		}

		_, err = a.client.Storage.ExportArchive(params, pipeWriter)
		if err != nil {
			log.Errorf("Error from ExportArchive: %s", err.Error())
			switch err := err.(type) {
			case *storage.ExportArchiveInternalServerError:
				plErr := InternalServerError(fmt.Sprintf("Server error from archive reader for device %s", deviceID))
				log.Errorf(plErr.Error())
				pipeWriter.CloseWithError(plErr)
			case *storage.ImportArchiveLocked:
				plErr := ResourceLockedError(fmt.Sprintf("Resource locked for device %s", deviceID))
				log.Errorf(plErr.Error())
				pipeWriter.CloseWithError(plErr)
			default:
				//Check for EOF.  Since the connection, transport, and data handling are
				//encapsulated inside of Swagger, we can only detect EOF by checking the
				//error string
				if strings.Contains(err.Error(), swaggerSubstringEOF) {
					log.Debugf("swagger error %s", err.Error())
					pipeWriter.Close()
				} else {
					pipeWriter.CloseWithError(err)
				}
			}
		} else {
			pipeWriter.Close()
		}
	}()

	return pipeReader, nil
}

// ArchiveImportWriter initializes a write stream for a path.  This is usually called
// for getting a writer during docker cp TO container.
func (a *ArchiveProxy) ArchiveImportWriter(ctx context.Context, store, deviceID string, filterSpec vicarchive.FilterSpec) (io.WriteCloser, error) {
	defer trace.End(trace.Begin(deviceID))

	if store == "" || deviceID == "" {
		return nil, fmt.Errorf("ArchiveImportWriter called with either empty store or deviceID")
	}

	var err error

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		// make sure we get out of io.Copy if context is canceled
		select {
		case <-ctx.Done():
			pipeWriter.Close()
		}
	}()

	go func() {
		// encodedFilter and destination are not required (from swagge spec) because
		// they are allowed to be empty.
		params := storage.NewImportArchiveParamsWithContext(ctx).
			WithStore(store).
			WithDeviceID(deviceID).
			WithArchive(pipeReader)

		// Encode the filter spec
		encodedFilter := ""
		if valueBytes, merr := json.Marshal(filterSpec); merr == nil {
			encodedFilter = base64.StdEncoding.EncodeToString(valueBytes)
			params = params.WithFilterSpec(&encodedFilter)
		}

		_, err = a.client.Storage.ImportArchive(params)
		if err != nil {
			switch err := err.(type) {
			case *storage.ImportArchiveInternalServerError:
				plErr := InternalServerError(fmt.Sprintf("Server error from archive writer for device %s", deviceID))
				log.Errorf(plErr.Error())
				pipeReader.CloseWithError(plErr)
			case *storage.ImportArchiveLocked:
				plErr := ResourceLockedError(fmt.Sprintf("Resource locked for device %s", deviceID))
				log.Errorf(plErr.Error())
				pipeReader.CloseWithError(plErr)
			default:
				//Check for EOF.  Since the connection, transport, and data handling are
				//encapsulated inside of Swagger, we can only detect EOF by checking the
				//error string
				if strings.Contains(err.Error(), swaggerSubstringEOF) {
					log.Errorf(err.Error())
					pipeReader.Close()
				} else {
					pipeReader.CloseWithError(err)
				}
			}
		} else {
			pipeReader.Close()
		}
	}()

	return pipeWriter, nil
}
