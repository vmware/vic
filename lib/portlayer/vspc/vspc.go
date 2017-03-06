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

package vspc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/telnet"
)

/* Vsphere telnet extension constants */
const (
	VmwareExt             byte          = 232
	KnownSuboptions1      byte          = 0
	KnownSuboptions2      byte          = 1
	UnknownSuboptionRcvd1 byte          = 2
	UnknownSuboptionRcvd2 byte          = 3
	VmotionBegin          byte          = 40
	VmotionGoahead        byte          = 41
	VmotionNotNow         byte          = 43
	VmotionPeer           byte          = 44
	VmotionPeerOK         byte          = 45
	VmotionComplete       byte          = 46
	VmotionAbort          byte          = 48
	DoProxy               byte          = 70
	WillProxy             byte          = 71
	WontProxy             byte          = 73
	VMVCUUID              byte          = 80
	GetVMVCUUID           byte          = 81
	VMName                byte          = 82
	GetVMName             byte          = 83
	VMBiosUUID            byte          = 84
	GetVMBiosUUID         byte          = 85
	VMLocationUUID        byte          = 86
	GetsVMLocationUUID    byte          = 87
	VMHealthCheckPeriod   time.Duration = 10 * time.Second
	vspcPort              int           = 2377
)
const remoteConnReadDeadline = 1 * time.Second

// VM is a struct that represents a VM
type VM struct {
	sync.Mutex
	// this is the current connection between the vspc and the containerVM serial port
	containerConn *telnet.Conn
	// in case of vmotion, at some point we will have two telnet connections to the containerVM
	// prevContainerConn is the connection of the source
	// containerConn will be the connection to the destination
	prevContainerConn *telnet.Conn
	// inVmotion is a boolean denoting whether this VM is in a state of vmotion or not
	inVmotion bool

	// remoteConn is the remote-system connection
	// this is the connection between the vspc and the attach-server
	remoteConn net.Conn

	vmUUID           string
	vmotionStarted   chan chan struct{}
	vmotionCompleted chan chan struct{}
}

// NewVM is the constructor of the VM
func NewVM(tc *telnet.Conn) *VM {
	return &VM{
		containerConn:    tc,
		inVmotion:        false,
		vmotionStarted:   make(chan chan struct{}),
		vmotionCompleted: make(chan chan struct{}),
	}
}

// Vspc is all the vspc singletons
type Vspc struct {
	vmManagerMu *sync.Mutex
	vmManager   map[string]*VM

	*telnet.Server

	attachSrvrAddr string
	attachSrvrPort uint

	doneCh chan bool
}

// NewVspc is the constructor
func NewVspc(address string, port uint, attachSrvrAddr string, attachSrvrPort uint, ch chan bool) *Vspc {
	vspc := &Vspc{
		vmManagerMu: new(sync.Mutex),
		vmManager:   make(map[string]*VM),

		attachSrvrAddr: attachSrvrAddr,
		attachSrvrPort: attachSrvrPort,

		doneCh: ch,
	}
	hdlr := Handler{vspc}
	opts := telnet.ServerOpts{
		Addr:        fmt.Sprintf("%s:%d", address, port),
		ServerOpts:  []byte{telnet.Binary, telnet.Sga, telnet.Echo},
		ClientOpts:  []byte{telnet.Binary, telnet.Sga, VmwareExt},
		DataHandler: hdlr.DataHdlr,
		CmdHandler:  hdlr.CmdHdlr,
	}
	vspc.Server = telnet.NewServer(opts)
	go vspc.monitorVMConnections()
	return vspc
}

// getVM returns the VM struct from its uuid
func (vspc *Vspc) getVM(uuid string) (*VM, bool) {
	vspc.vmManagerMu.Lock()
	defer vspc.vmManagerMu.Unlock()
	if vm, ok := vspc.vmManager[uuid]; ok {
		return vm, true
	}
	return nil, false
}

// addVM adds a VM to the map
func (vspc *Vspc) addVM(uuid string, vm *VM) {
	vspc.vmManagerMu.Lock()
	defer vspc.vmManagerMu.Unlock()
	vspc.vmManager[uuid] = vm
}

// relayReads reads from the AttachServer connection and relays the data to the telnet connection
func (vspc *Vspc) relayReads(containervm *VM, conn net.Conn) {
	vmotion := false
	var tmpBuf bytes.Buffer
	for {
		select {
		case ch := <-containervm.vmotionStarted:
			vmotion = true
			ch <- struct{}{}
			log.Infof("vspc started to buffer data coming from the remote system")
		case ch := <-containervm.vmotionCompleted:
			vmotion = false
			ch <- struct{}{}
			log.Infof("vspc stopped buffering data coming from the remote system")
		default:
			b := make([]byte, 4096)
			conn.SetReadDeadline(time.Now().Add(remoteConnReadDeadline))
			n, err := conn.Read(b)
			if n > 0 {
				log.Debugf("vspc read %d bytes from the  remote system connection", n)
				if !vmotion {
					if tmpBuf.Len() > 0 {
						buf, err := ioutil.ReadAll(&tmpBuf)
						if err != nil {
							log.Errorf("read error from vspc temporary buffer: %v", err)
						}
						log.Infof("vspc writing buffered data during vmotion to the containerVM")
						if n, err := containervm.containerConn.WriteData(buf); n == -1 {
							log.Errorf("vspc: RelayReads: %v", err)
							return
						}
					}
					if n, err := containervm.containerConn.WriteData(b[:n]); n == -1 {
						log.Errorf("vspc: RelayReads: %v", err)
						return
					}
					log.Infof("vspc relayed the read data to the containerVM")
				} else {
					tmpBuf.Write(b[:n])
				}
			}
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Info("remote system connection closed")
				return
			}
		}
	}
}

func (vspc *Vspc) vmFromTelnetConn(tc *telnet.Conn) (*VM, bool) {
	vspc.vmManagerMu.Lock()
	defer vspc.vmManagerMu.Unlock()
	for _, v := range vspc.vmManager {
		if v.containerConn == tc {
			return v, true
		}
	}
	return nil, false
}

// monitorVMConnections is to monitor connections and delete VM recoreds from the vmManager map when the containerVM exits
func (vspc *Vspc) monitorVMConnections() {
	ticker := time.NewTicker(VMHealthCheckPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				vspc.vmManagerMu.Lock()

				for k, vm := range vspc.vmManager {
					vm.Lock()
					if !vm.inVmotion && vm.containerConn.IsClosed() { // vm just shut down
						log.Debugf("(vspc) detected closed connection for VM %s", k)
						log.Debugf("(vspc) deleting vm records from the vm manager %s", k)
						delete(vspc.vmManager, k)
					}
					vm.Unlock()
				}
				vspc.vmManagerMu.Unlock()
			case <-vspc.doneCh:
				ticker.Stop()
				return
			}
		}
	}()
}
