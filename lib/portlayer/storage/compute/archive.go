// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

// StatPath will use Guest Tools to stat a given path in the container
package compute

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"sync"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
)

func FileTransferToGuest(op trace.Operation, vc *exec.Container, path string, reader io.ReadCloser) error {
	defer trace.End(trace.Begin(""))

	// set up file manager
	client := vc.VIM25Reference()
	filemgr, err := guest.NewOperationsManager(client, vc.VMReference()).FileManager(op)
	if err != nil {
		return err
	}

	// authenticate client and parse container host/port
	auth := types.NamePasswordAuthentication{
		Username: vc.ExecConfig.ID,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	var buf bytes.Buffer
	var size int64
	gZipReader := zipTar(op, reader)

	go func() {
		defer reader.Close()
		defer wg.Done()

		size, err = io.Copy(&buf, gZipReader)
		if err != nil {
			op.Errorf("error gzipping tar stream: %s", err.Error())
			size = 0
		}
	}()

	wg.Wait()

	if size == 0 {
		return fmt.Errorf("archive upload size cannot be 0")
	}

	guestTransferURL, err := filemgr.InitiateFileTransferToGuest(op, &auth, path, &types.GuestPosixFileAttributes{}, size, true)
	if err != nil {
		return err
	}

	url, err := client.ParseURL(guestTransferURL)
	if err != nil {
		return err
	}

	// upload tar archive to url
	op.Debugf("Uploading: %v --- %s on %s", url, path, vc.ExecConfig.ID)
	params := soap.DefaultUpload
	params.ContentLength = size
	op.Debugf("%s", string(buf.Bytes()))
	err = client.Upload(bytes.NewReader(buf.Bytes()), url, &params)
	if err != nil {
		return err
	}
	return nil
}

func FileTransferFromGuest(op trace.Operation, vc *exec.Container, path string) (io.ReadCloser, error) {
	defer trace.End(trace.Begin(""))

	// set up file manager.
	client := vc.VIM25Reference()
	filemgr, err := guest.NewOperationsManager(client, vc.VMReference()).FileManager(op)
	if err != nil {
		return nil, err
	}
	auth := types.NamePasswordAuthentication{
		Username: vc.ExecConfig.ID,
	}

	// authenticate client and parse container host/port.
	guestInfo, err := filemgr.InitiateFileTransferFromGuest(op, &auth, path)
	if err != nil {
		return nil, err
	}

	url, err := client.ParseURL(guestInfo.Url)
	if err != nil {
		return nil, err
	}

	// download from guest. if download is a file, create a tar out of it.
	// guest tools will not tar up single files.
	op.Debugf("Downloading: %v --- %s from %d", url, path, vc.ExecConfig.ID)
	params := soap.DefaultDownload
	rc, contentLength, err := client.Download(url, &params)
	if err != nil {
		return nil, err
	}

	gc, err := createTarFromFile(op, rc, vc, path, contentLength)
	if err != nil {
		return nil, err
	}

	return unZipTar(op, gc), nil
}

//----------
// Utility Functions
//----------

func createTarFromFile(op trace.Operation, reader io.ReadCloser, vc *exec.Container, path string, size int64) (io.ReadCloser, error) {
	stat, err := StatPath(op, vc, path)
	if err != nil {
		return nil, err
	}

	if types.GuestFileType(stat.Type) != types.GuestFileTypeFile {
		return reader, nil
	}

	tarOut, tarIn := io.Pipe()
	go func() {
		gZipWriter := gzip.NewWriter(tarIn)
		tarWriter := tar.NewWriter(gZipWriter)
		defer reader.Close()
		defer tarIn.Close()
		defer gZipWriter.Close()
		defer tarWriter.Close()

		hdr := &tar.Header{
			Name:    filepath.Base(stat.Path),
			Size:    size,
			ModTime: *stat.Attributes.GetGuestFileAttributes().ModificationTime,
		}
		switch types.GuestFileType(stat.Type) {
		case types.GuestFileTypeDirectory:
			hdr.Mode = int64(os.ModeDir)
		case types.GuestFileTypeSymlink:
			hdr.Mode = int64(os.ModeSymlink)
		default:
			hdr.Mode = int64(0600)
		}

		op.Debugf("Write Header: %v", *hdr)
		if err = tarWriter.WriteHeader(hdr); err != nil {
			tarIn.CloseWithError(err)
			return
		}

		op.Debugf("Write Body: %d", size)
		// write file content body
		if _, err := io.CopyN(tarWriter, reader, size); err != nil {
			tarIn.CloseWithError(err)
			return
		}

		op.Debugf("Return")
	}()

	return tarOut, nil
}

func unZipTar(op trace.Operation, reader io.ReadCloser) io.ReadCloser {
	// create a writer for gzip compressiona nd a tar archive
	tarOut, tarIn := io.Pipe()
	go func() {
		gZipReader, err := gzip.NewReader(reader)
		if err != nil {
			op.Errorf("Error in unziptar: %s", err.Error())
			tarIn.CloseWithError(err)
			return
		}
		tarReader := tar.NewReader(gZipReader)
		tarWriter := tar.NewWriter(tarIn)
		defer reader.Close()
		defer tarIn.Close()
		defer gZipReader.Close()
		defer tarWriter.Close()

		// grab tar stream from guest tools. zip it up if there are no errors
		for {
			hdr, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				op.Errorf("Error in unziptar: %s", err.Error())
				tarIn.CloseWithError(err)
				return
			}

			op.Debugf("read/write header: %#v", *hdr)
			if err = tarWriter.WriteHeader(hdr); err != nil {
				op.Errorf("Error in unziptar: %s", err.Error())
				tarIn.CloseWithError(err)
				return
			}

			op.Debugf("read/write body")
			if _, err := io.Copy(tarWriter, tarReader); err != nil {
				op.Errorf("Error in unziptar: %s", err.Error())
				tarIn.CloseWithError(err)
				return
			}
		}

		op.Debugf("return")
	}()

	return tarOut
}

func zipTar(op trace.Operation, reader io.Reader) io.ReadCloser {
	gZipOut, gZipIn := io.Pipe()
	go func() {
		tarReader := tar.NewReader(reader)
		gZipWriter := gzip.NewWriter(gZipIn)
		tarWriter := tar.NewWriter(gZipWriter)
		defer gZipIn.Close()
		defer gZipWriter.Close()
		defer tarWriter.Close()

		for {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}
			if err != nil {
				op.Errorf("Error in zipTar: %s", err.Error())
				gZipIn.CloseWithError(err)
				return
			}

			if err := tarWriter.WriteHeader(header); err != nil {
				op.Errorf("Error in zipTar: %s", err.Error())
				gZipIn.CloseWithError(err)
				return
			}

			if _, err := io.Copy(tarWriter, tarReader); err != nil {
				op.Errorf("Error in zipTar: %s", err.Error())
				gZipIn.CloseWithError(err)
				return
			}
		}
	}()
	return gZipOut
}
