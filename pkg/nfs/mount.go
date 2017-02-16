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

// MOUNT
// RFC 1813 Section 5.0

import (
	"errors"
	"fmt"

	"github.com/vmware/vic/pkg/nfs/rpc"
	"github.com/vmware/vic/pkg/nfs/xdr"
)

const (
	MOUNT_PROG = 100005
	MOUNT_VERS = 3

	MOUNTPROC3_NULL   = 0
	MOUNTPROC3_MNT    = 1
	MOUNTPROC3_UMNT   = 3
	MOUNTPROC3_EXPORT = 5

	MNT3_OK             = 0     // no error
	MNT3ERR_PERM        = 1     // Not owner
	MNT3ERR_NOENT       = 2     // No such file or directory
	MNT3ERR_IO          = 5     // I/O error
	MNT3ERR_ACCES       = 13    // Permission denied
	MNT3ERR_NOTDIR      = 20    // Not a directory
	MNT3ERR_INVAL       = 22    // Invalid argument
	MNT3ERR_NAMETOOLONG = 63    // Filename too long
	MNT3ERR_NOTSUPP     = 10004 // Operation not supported
	MNT3ERR_SERVERFAULT = 10006 // A failure on the server
)

type Export struct {
	Dir    string
	Groups []Group
}

type Group struct {
	Name string
}

type Mount struct {
	*rpc.Client
	auth    rpc.Auth
	dirPath string
	Addr    string
}

func (m *Mount) Unmount() error {
	type umount struct {
		rpc.Header
		Dirpath string
	}

	_, err := m.Call(&umount{
		rpc.Header{
			Rpcvers: 2,
			Prog:    MOUNT_PROG,
			Vers:    MOUNT_VERS,
			Proc:    MOUNTPROC3_UMNT,
			// Weirdly, the spec calls for AUTH_UNIX or better, but AUTH_NULL
			// works here on a linux NFS kernel server.  Follow the spec
			// anyway.
			Cred: m.auth,
			Verf: rpc.AUTH_NULL,
		},
		m.dirPath,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Mount) Mount(dirpath string, auth rpc.Auth) (*Target, error) {
	type mount struct {
		rpc.Header
		Dirpath string
	}

	buf, err := m.Call(&mount{
		rpc.Header{
			Rpcvers: 2,
			Prog:    MOUNT_PROG,
			Vers:    MOUNT_VERS,
			Proc:    MOUNTPROC3_MNT,
			Cred:    auth,
			Verf:    rpc.AUTH_NULL,
		},
		dirpath,
	})
	if err != nil {
		return nil, err
	}

	mountstat3, buf := xdr.Uint32(buf)
	switch mountstat3 {
	case MNT3_OK:
		fh, buf := xdr.Opaque(buf)
		_, buf = xdr.Uint32List(buf)

		m.dirPath = dirpath
		m.auth = auth

		vol, err := NewTarget(m.Addr, auth, fh, dirpath)
		if err != nil {
			return nil, err
		}

		return vol, nil

	case MNT3ERR_PERM:
		return nil, errors.New("MNT3ERR_PERM")
	case MNT3ERR_NOENT:
		return nil, errors.New("MNT3ERR_NOENT")
	case MNT3ERR_IO:
		return nil, errors.New("MNT3ERR_IO")
	case MNT3ERR_ACCES:
		return nil, errors.New("MNT3ERR_ACCES")
	case MNT3ERR_NOTDIR:
		return nil, errors.New("MNT3ERR_NOTDIR")
	case MNT3ERR_NAMETOOLONG:
		return nil, errors.New("MNT3ERR_NAMETOOLONG")
	}
	return nil, fmt.Errorf("unknown mount stat: %d", mountstat3)
}

func DialMount(addr string) (*Mount, error) {
	// get MOUNT port
	m := rpc.Mapping{
		Prog: MOUNT_PROG,
		Vers: MOUNT_VERS,
		Prot: rpc.IPPROTO_TCP,
		Port: 0,
	}

	client, err := DialService(addr, m)
	if err != nil {
		return nil, err
	}

	return &Mount{
		Client: client,
		Addr:   addr,
	}, nil
}
