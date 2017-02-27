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

package rpc

import (
	"bytes"
	"fmt"

	"github.com/vmware/vic/pkg/nfs/xdr"
)

// PORTMAP
// RFC 1057 Section A.1

const (
	PMAP_PORT = 111
	PMAP_PROG = 100000
	PMAP_VERS = 2

	PMAPPROC_GETPORT = 3
	PMAPPROC_DUMP    = 4

	IPPROTO_TCP = 6
	IPPROTO_UDP = 17
)

type Mapping struct {
	Prog uint32
	Vers uint32
	Prot uint32
	Port uint32
}

type Portmapper struct {
	*Client
	host string
}

func (p *Portmapper) Getport(mapping Mapping) (int, error) {
	type getport struct {
		Header
		Mapping
	}
	msg := &getport{
		Header{
			Rpcvers: 2,
			Prog:    PMAP_PROG,
			Vers:    PMAP_VERS,
			Proc:    PMAPPROC_GETPORT,
			Cred:    AUTH_NULL,
			Verf:    AUTH_NULL,
		},
		mapping,
	}
	buf, err := p.Call(msg)
	if err != nil {
		return 0, err
	}
	port, _ := xdr.Uint32(buf)
	return int(port), nil
}

func (p *Portmapper) Dump() ([]Mapping, error) {
	type dump struct {
		Header
	}
	msg := &dump{
		Header{
			Rpcvers: 2,
			Prog:    PMAP_PROG,
			Vers:    PMAP_VERS,
			Proc:    PMAPPROC_DUMP,
			Cred:    AUTH_NULL,
			Verf:    AUTH_NULL,
		},
	}
	buf, err := p.Call(msg)
	if err != nil {
		return nil, err
	}
	var mm []Mapping
	for r := bytes.NewBuffer(buf); r.Len() > 0; {
		type index struct {
			Index uint32
			M     Mapping
		}
		var i index
		if err := xdr.Read(r, &i); err != nil {
			return nil, err
		}
		mm = append(mm, i.M)
	}
	return mm, nil
}

func DialPortmapper(net, host string) (*Portmapper, error) {
	client, err := DialTCP(net, nil, fmt.Sprintf("%s:%d", host, PMAP_PORT))
	if err != nil {
		return nil, err
	}
	return &Portmapper{client, host}, nil
}
