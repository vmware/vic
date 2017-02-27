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

package nfs

import (
	"bytes"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/nfs/rpc"
	"github.com/vmware/vic/pkg/nfs/xdr"
)

type File struct {
	*Target

	// current position
	curr   uint64
	fsinfo *FSInfo

	// filehandle to the file
	fh []byte
}

func (f *File) Read(p []byte) (int, error) {
	type ReadArgs struct {
		rpc.Header
		FH     []byte
		Offset uint64
		Count  uint32
	}

	type ReadRes struct {
		Follows uint32
		Attrs   struct {
			Attrs Fattr
		}
		Count uint32
		Eof   uint32
		Data  struct {
			Length uint32
		}
	}

	readSize := uint32(min(uint64(f.fsinfo.RTPref), uint64(len(p))))
	log.Debugf("read(%x) len=%d offset=%d", f.fh, readSize, f.curr)

	buf, err := f.call(&ReadArgs{
		Header: rpc.Header{
			Rpcvers: 2,
			Prog:    NFS3_PROG,
			Vers:    NFS3_VERS,
			Proc:    NFSPROC3_READ,
			Cred:    f.auth,
			Verf:    rpc.AUTH_NULL,
		},
		FH:     f.fh,
		Offset: uint64(f.curr),
		Count:  readSize,
	})

	if err != nil {
		log.Debugf("read(%x): %s", f.fh, err.Error())
		return 0, err
	}

	r := bytes.NewBuffer(buf)
	readres := &ReadRes{}
	if err = xdr.Read(r, readres); err != nil {
		return 0, err
	}

	f.curr = f.curr + uint64(readres.Data.Length)
	n, err := r.Read(p[:readres.Data.Length])
	if err != nil {
		return n, err
	}

	if readres.Eof != 0 {
		err = io.EOF
	}

	return n, err
}

func (f *File) Write(p []byte) (int, error) {
	type WriteArgs struct {
		rpc.Header
		FH     []byte
		Offset uint64
		Count  uint32

		// UNSTABLE(0), DATA_SYNC(1), FILE_SYNC(2) default
		How      uint32
		Contents []byte
	}

	totalToWrite := len(p)
	writeSize := uint64(min(uint64(f.fsinfo.WTPref), uint64(totalToWrite)))

	_, err := f.call(&WriteArgs{
		Header: rpc.Header{
			Rpcvers: 2,
			Prog:    NFS3_PROG,
			Vers:    NFS3_VERS,
			Proc:    NFSPROC3_WRITE,
			Cred:    f.auth,
			Verf:    rpc.AUTH_NULL,
		},
		FH:       f.fh,
		Offset:   f.curr,
		Count:    uint32(writeSize),
		How:      2,
		Contents: p[:writeSize],
	})

	if err != nil {
		log.Debugf("write(%x): %s", f.fh, err.Error())
		return int(writeSize), err
	}

	log.Debugf("write(%x) len=%d offset=%d written=%d total=%d",
		f.fh, writeSize, f.curr, writeSize, totalToWrite)

	f.curr = f.curr + writeSize

	return int(writeSize), nil
}

func (f *File) Close() error {
	type CommitArg struct {
		rpc.Header
		FH     []byte
		Offset uint64
		Count  uint32
	}

	_, err := f.call(&CommitArg{
		Header: rpc.Header{
			Rpcvers: 2,
			Prog:    NFS3_PROG,
			Vers:    NFS3_VERS,
			Proc:    NFSPROC3_COMMIT,
			Cred:    f.auth,
			Verf:    rpc.AUTH_NULL,
		},
		FH: f.fh,
	})

	if err != nil {
		log.Debugf("commit(%x): %s", f.fh, err.Error())
		return err
	}

	return nil
}

// OpenFile writes to an existing file or creates one
func (v *Target) OpenFile(path string, perm os.FileMode) (*File, error) {
	_, fh, err := v.Lookup(path)
	if err != nil {
		if os.IsNotExist(err) {
			fh, err = v.Create(path, perm)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	f := &File{
		Target: v,
		fsinfo: v.fsinfo,
		fh:     fh,
	}

	return f, nil
}

// Open opens a file for reading
func (v *Target) Open(path string) (*File, error) {
	_, fh, err := v.Lookup(path)
	if err != nil {
		return nil, err
	}

	f := &File{
		Target: v,
		fsinfo: v.fsinfo,
		fh:     fh,
	}

	return f, nil
}

func min(x, y uint64) uint64 {
	if x > y {
		return y
	}
	return x
}
