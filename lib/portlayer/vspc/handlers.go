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
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/telnet"
)

// Handler is the handler struct for the vspc
type Handler struct {
	vspc *Vspc
}

// DataHdlr is the telnet data handler
func (h *Handler) DataHdlr(w io.Writer, b []byte, tc *telnet.Conn) {
	vm, exists := h.vspc.vmFromTelnetConn(tc)
	if !exists {
		log.Errorf("vspc datahandler cannot find the vm with the provided data")
	}
	vm.remoteConn.Write(b)
}

// CmdHdlr is the telnet command handler
func (h *Handler) CmdHdlr(w io.Writer, b []byte, tc *telnet.Conn) {
	if isKnownSuboptions(b) {
		log.Infof("vspc received KNOWN-SUBOPTIONS command")
		h.HandleKnownSuboptions(w, b)
	} else if isDoProxy(b) {
		log.Infof("vspc received DO-PROXY command")
		h.HandleDoProxy(w, b)
	} else if isVmotionBegin(b) {
		log.Infof("vspc received VMOTION-BEGIN command")
		h.HandleVmotionBegin(w, tc, b)
	} else if isVmotionPeer(b) {
		log.Infof("vspc received VMOTION-PEER command")
		h.HandleVmotionPeer(w, b)
	} else if isVMUUID(b) {
		log.Infof("vspc received VMUUID command")
		h.HandleVMUUID(w, tc, b)
	} else if isVmotionComplete(b) {
		log.Infof("vspc received VMOTION-COMPLETE command")
		h.HandleVmotionComplete(w, tc, b)
	} else if isVmotionAbort(b) {
		log.Infof("vspc received VMOTION-ABORT command")
		h.HandleVmotionAbort(w, b)
	}
}

// HandleVMUUID handles the telnet vm-uuid response
func (h *Handler) HandleVMUUID(w io.Writer, tc *telnet.Conn, b []byte) {
	vmuuid := strings.Replace(string(b[3:len(b)-1]), " ", "", -1)
	log.Infof("vmuuid of the connected containerVM: %s", vmuuid)
	// check if there exists another vm with the same vmuuid
	vm, exists := h.vspc.getVM(vmuuid)
	if !exists {
		// create a new vm associated with this telnet connection
		vm = NewVM(tc)
		log.Infof("attempting to connect to the attach server")
		remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h.vspc.attachSrvrAddr, h.vspc.attachSrvrPort))
		if err != nil {
			log.Errorf("cannot connect to the attach-server: %v", err)
		}
		vm.remoteConn = remoteConn
		vm.vmUUID = vmuuid
		h.vspc.addVM(vmuuid, vm)
		// relay Reads from the remote system connection to the telnet connection associated with this vm
		go h.vspc.relayReads(vm, remoteConn)
	} else { //the vm existed before and was shut down or vmotioned
		log.Infof("established second serial-port telnet connection with vm (vmuuid: %s)", vm.vmUUID)
		vm.Lock()
		defer vm.Unlock()
		vm.prevContainerConn = vm.containerConn
		vm.containerConn = tc
	}
}

// HandleKnownSuboptions handles the known suboptions telnet command
func (h *Handler) HandleKnownSuboptions(w io.Writer, b []byte) {
	log.Infof("handling KNOWN-SUBOPTIONS")
	var resp []byte
	suboptions := b[3 : len(b)-1]
	resp = append(resp, []byte{telnet.Iac, telnet.Sb, VmwareExt, KnownSuboptions2}...)
	resp = append(resp, suboptions...)
	resp = append(resp, telnet.Iac, telnet.Se)
	log.Debugf("response to KNOWN-SUBOPTIONS: %v", resp)

	if bytes.IndexByte(suboptions, GetVMVCUUID) != -1 && bytes.IndexByte(suboptions, VMVCUUID) != -1 {
		resp = append(resp, getVMUUID()...)
	}
	w.Write(resp)
}

// HandleDoProxy handles the DO-PROXY telnet command
func (h *Handler) HandleDoProxy(w io.Writer, b []byte) {
	log.Infof("handling DO-PROXY")
	var resp []byte
	resp = append(resp, []byte{telnet.Iac, telnet.Sb, VmwareExt, WillProxy, telnet.Iac, telnet.Se}...)
	log.Debugf("response to DO-PROXY: %v", resp)
	w.Write(resp)
}

// HandleVmotionBegin handles the VMOTION-BEGIN telnet command
func (h *Handler) HandleVmotionBegin(w io.Writer, tc *telnet.Conn, b []byte) {
	if vm, exists := h.vspc.vmFromTelnetConn(tc); exists {
		vm.Lock()
		vm.inVmotion = true
		vm.Unlock()
		ch := make(chan struct{})
		vm.vmotionStarted <- ch
		<-ch
	}
	log.Infof("handling VMOTION-BEGIN")
	seq := b[3 : len(b)-1]
	var escapedSeq []byte
	for _, v := range seq {
		if v == telnet.Iac {
			escapedSeq = append(escapedSeq, telnet.Iac)
		}
		escapedSeq = append(escapedSeq, v)
	}
	secret := make([]byte, 4)
	var escapedSecret []byte
	rand.Read(secret)
	// escaping Iac
	for _, v := range secret {
		if v == telnet.Iac {
			escapedSecret = append(escapedSecret, telnet.Iac)
		}
		escapedSecret = append(escapedSecret, v)
	}
	var resp []byte
	resp = append(resp, []byte{telnet.Iac, telnet.Sb, VmwareExt, VmotionGoahead}...)
	resp = append(resp, escapedSeq...)
	resp = append(resp, escapedSecret...)
	resp = append(resp, telnet.Iac, telnet.Se)
	log.Debugf("response to VMOTION-BEGIN: %v", resp)
	w.Write(resp)
}

// HandleVmotionPeer handles the VMOTION-PEER telnet command
func (h *Handler) HandleVmotionPeer(w io.Writer, b []byte) {
	log.Infof("Handling VMOTION-PEER")
	// cookie is the sequence + secret
	cookie := b[3 : len(b)-1]
	var resp []byte
	resp = append(resp, []byte{telnet.Iac, telnet.Sb, VmwareExt, VmotionPeerOK}...)
	resp = append(resp, cookie...)
	resp = append(resp, telnet.Iac, telnet.Se)
	log.Debugf("response to VMOTION-PEER: %v", resp)
	w.Write(resp)
}

// HandleVmotionComplete handles the VMOTION-Complete telnet command
func (h *Handler) HandleVmotionComplete(w io.Writer, tc *telnet.Conn, b []byte) {
	log.Infof("handling VMOTION-COMPLETE")
	if vm, exists := h.vspc.vmFromTelnetConn(tc); exists {
		vm.Lock()
		vm.prevContainerConn = nil
		vm.inVmotion = false
		vm.Unlock()
		ch := make(chan struct{})
		vm.vmotionCompleted <- ch
		<-ch
		log.Info("vMotion completed successfully")
	} else {
		log.Errorf("couldnt find previous information of vm after vmotion (vmuuid: %s)", vm.vmUUID)
	}

}

// HandleVmotionAbort handles the VMOTION-abort telnet command
func (h *Handler) HandleVmotionAbort(w io.Writer, b []byte) {
	log.Errorf("vMotion failed")
}
