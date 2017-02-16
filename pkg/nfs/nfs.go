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
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/user"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/nfs/rpc"
)

// NFS version 3
// RFC 1813

const (
	NFS3_PROG = 100003
	NFS3_VERS = 3

	// program methods
	NFSPROC3_LOOKUP      = 3
	NFSPROC3_READ        = 6
	NFSPROC3_WRITE       = 7
	NFSPROC3_CREATE      = 8
	NFSPROC3_MKDIR       = 9
	NFSPROC3_REMOVE      = 12
	NFSPROC3_RMDIR       = 13
	NFSPROC3_READDIRPLUS = 17
	NFSPROC3_FSINFO      = 19
	NFSPROC3_COMMIT      = 21

	// file types
	NF3REG  = 1
	NF3DIR  = 2
	NF3BLK  = 3
	NF3CHR  = 4
	NF3LNK  = 5
	NF3SOCK = 6
	NF3FIFO = 7
)

type Diropargs3 struct {
	FH       []byte
	Filename string
}

type Sattr3 struct {
	Mode  SetMode
	UID   SetUID
	GID   SetUID
	Size  uint64
	Atime NFS3Time
	Mtime NFS3Time
}

type SetMode struct {
	Set  uint32
	Mode uint32
}

type SetUID struct {
	Set uint32
	UID uint32
}

type NFS3Time struct {
	Seconds  uint32
	Nseconds uint32
}

type Fattr struct {
	Type                uint32
	Mode                uint32
	Nlink               uint32
	UID                 uint32
	GID                 uint32
	Size                uint64
	Used                uint64
	SpecData            [2]uint32
	FSID                uint64
	Fileid              uint64
	Atime, Mtime, Ctime NFS3Time
}

type EntryPlus struct {
	FileId   uint64
	FileName string
	Cookie   uint64
	Attr     struct {
		Follows uint32
		Attr    Fattr
	}
	FHSet        uint32
	FH           []byte
	ValueFollows uint32
}

type FSInfo struct {
	Follows    uint32
	RTMax      uint32
	RTPref     uint32
	RTMult     uint32
	WTMax      uint32
	WTPref     uint32
	WTMult     uint32
	DTPref     uint32
	Size       uint64
	TimeDelta  NFS3Time
	Properties uint32
}

// Dial an RPC svc after getting the port from the portmapper
func DialService(addr string, prog rpc.Mapping) (*rpc.Client, error) {
	pm, err := rpc.DialPortmapper("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer pm.Close()

	port, err := pm.Getport(prog)
	if err != nil {
		return nil, err
	}

	client, err := dialService(addr, port)

	return client, nil
}

func dialService(addr string, port int) (*rpc.Client, error) {
	var (
		ldr    *net.TCPAddr
		client *rpc.Client
	)

	usr, err := user.Current()

	// Unless explicitly configured, the target will likely reject connections
	// from non-privileged ports.
	if err == nil && usr.Uid == "0" {
		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

		var p int
		for {
			p = r1.Intn(1024)
			if p < 0 {
				continue
			}

			ldr = &net.TCPAddr{
				Port: p,
			}

			client, err = rpc.DialTCP("tcp", ldr, fmt.Sprintf("%s:%d", addr, port))
			if err == nil {
				break
			}
			// bind error, try again
			if isAddrInUse(err) {
				continue
			}

			return nil, err
		}

		log.Debugf("using random port %d -> %d", p, port)
	} else {

		client, err = rpc.DialTCP("tcp", ldr, fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func isAddrInUse(err error) bool {
	if er, ok := (err.(*net.OpError)); ok {
		if syser, ok := er.Err.(*os.SyscallError); ok {
			return syser.Err == syscall.EADDRINUSE
		}
	}

	return false
}
