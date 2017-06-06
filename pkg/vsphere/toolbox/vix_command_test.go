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

package toolbox

import (
	"bytes"
	"context"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/vmware/vic/pkg/vsphere/toolbox/hgfs"
)

type VixCommandClient struct {
	Service *Service
	Header  *VixCommandRequestHeader
	creds   []byte
}

func NewVixCommandClient() *VixCommandClient {
	Trace = testing.Verbose()
	hgfs.Trace = Trace

	creds, _ := (&VixUserCredentialNamePassword{
		Name:     "user",
		Password: "pass",
	}).MarshalBinary()

	header := new(VixCommandRequestHeader)
	header.Magic = vixCommandMagicWord

	header.UserCredentialType = vixUserCredentialNamePassword
	header.CredentialLength = uint32(len(creds))

	in := new(mockChannelIn)
	out := new(mockChannelOut)

	return &VixCommandClient{
		creds:   creds,
		Header:  header,
		Service: NewService(in, out),
	}
}

func (c *VixCommandClient) Request(op uint32, size int, m encoding.BinaryMarshaler) []byte {
	c.Header.OpCode = op
	c.Header.BodyLength = uint32(size)

	var buf bytes.Buffer
	_, _ = buf.Write([]byte("\"reqname\"\x00"))
	_ = binary.Write(&buf, binary.LittleEndian, c.Header)

	b, err := m.MarshalBinary()
	if err != nil {
		panic(err)
	}
	_, _ = buf.Write(b)

	data := append(buf.Bytes(), c.creds...)
	reply, err := c.Service.VixCommand.Dispatch(data)
	if err != nil {
		panic(err)
	}

	return reply
}

func TestMarshalVixMsgStartProgramRequest(t *testing.T) {
	requests := []*VixMsgStartProgramRequest{
		{},
		{
			ProgramPath: "/bin/date",
		},
		{
			ProgramPath: "/bin/date",
			Arguments:   "--date=@2147483647",
		},
		{
			ProgramPath: "/bin/date",
			WorkingDir:  "/tmp",
		},
		{
			ProgramPath: "/bin/date",
			WorkingDir:  "/tmp",
			EnvVars:     []string{"FOO=bar"},
		},
		{
			ProgramPath: "/bin/date",
			WorkingDir:  "/tmp",
			EnvVars:     []string{"FOO=bar", "BAR=foo"},
		},
	}

	for i, in := range requests {
		buf, err := in.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		out := new(VixMsgStartProgramRequest)

		err = out.UnmarshalBinary(buf)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(in, out) {
			t.Errorf("%d marshal mismatch", i)
		}
	}
}

func vixRC(buf []byte) int {
	args := bytes.SplitN(buf, []byte{' '}, 2)
	rc, err := strconv.Atoi(string(args[0]))
	if err != nil {
		panic(err)
	}
	return rc
}

func TestVixRelayedCommandHandler(t *testing.T) {
	Trace = true
	if !testing.Verbose() {
		// cover Trace paths but discard output
		traceLog = ioutil.Discard
	}

	in := new(mockChannelIn)
	out := new(mockChannelOut)

	service := NewService(in, out)

	vix := service.VixCommand

	msg := []byte("\"reqname\"\x00")

	_, err := vix.Dispatch(msg) // io.EOF
	if err == nil {
		t.Fatal("expected error")
	}

	header := new(VixCommandRequestHeader)

	marshal := func(m ...encoding.BinaryMarshaler) []byte {
		var buf bytes.Buffer
		_, _ = buf.Write(msg)
		_ = binary.Write(&buf, binary.LittleEndian, header)

		for _, e := range m {
			b, err := e.MarshalBinary()
			if err != nil {
				panic(err)
			}
			_, _ = buf.Write(b)
		}

		return buf.Bytes()
	}

	// header.Magic not set
	reply, _ := vix.Dispatch(marshal())
	rc := vixRC(reply)
	if rc != vixInvalidMessageHeader {
		t.Fatalf("%q", reply)
	}

	// header.OpCode not set
	header.Magic = vixCommandMagicWord
	reply, _ = vix.Dispatch(marshal())
	rc = vixRC(reply)
	if rc != vixUnrecognizedCommandInGuest {
		t.Fatalf("%q", reply)
	}

	// valid request for GetToolsState
	header.OpCode = vixCommandGetToolsState
	reply, _ = vix.Dispatch(marshal())
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("%q", reply)
	}

	// header.UserCredentialType not set
	header.OpCode = vixCommandStartProgram
	request := new(VixMsgStartProgramRequest)
	buf := marshal(request)
	reply, _ = vix.Dispatch(marshal())
	rc = vixRC(reply)
	if rc != vixAuthenticationFail {
		t.Fatalf("%q", reply)
	}

	creds, _ := (&VixUserCredentialNamePassword{
		Name:     "user",
		Password: "pass",
	}).MarshalBinary()

	header.BodyLength = uint32(binary.Size(request.header))
	header.UserCredentialType = vixUserCredentialNamePassword
	header.CredentialLength = uint32(len(creds))

	// ProgramPath not set
	buf = append(marshal(request), creds...)
	reply, _ = vix.Dispatch(buf)
	rc = vixRC(reply)
	if rc != vixFail {
		t.Fatalf("%q", reply)
	}

	vix.ProcessStartCommand = func(pm *ProcessManager, r *VixMsgStartProgramRequest) (int64, error) {
		return -1, nil
	}

	// valid request for StartProgram
	buf = append(marshal(request), creds...)
	reply, _ = vix.Dispatch(buf)
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("%q", reply)
	}

	vix.Authenticate = func(_ VixCommandRequestHeader, data []byte) error {
		var c VixUserCredentialNamePassword
		if err := c.UnmarshalBinary(data); err != nil {
			panic(err)
		}

		return errors.New("you shall not pass")
	}

	// fail auth with our own handler
	buf = append(marshal(request), creds...)
	reply, _ = vix.Dispatch(buf)
	rc = vixRC(reply)
	if rc != vixAuthenticationFail {
		t.Fatalf("%q", reply)
	}

	vix.Authenticate = nil

	// cause VixUserCredentialNamePassword.UnmarshalBinary to error
	// first by EOF reading header, second in base64 decode
	for _, l := range []uint32{1, 10} {
		header.CredentialLength = l
		buf = append(marshal(request), creds...)
		reply, _ = vix.Dispatch(buf)
		rc = vixRC(reply)
		if rc != vixAuthenticationFail {
			t.Fatalf("%q", reply)
		}
	}
}

// cover misc error paths
func TestVixCommandErrors(t *testing.T) {
	r := new(VixMsgStartProgramRequest)
	err := r.UnmarshalBinary(nil)
	if err == nil {
		t.Error("expected error")
	}

	r.header.NumEnvVars = 1
	buf, _ := r.MarshalBinary()
	err = r.UnmarshalBinary(buf)
	if err == nil {
		t.Error("expected error")
	}

	c := new(VixRelayedCommandHandler)
	_, err = c.StartCommand("", r.VixCommandRequestHeader, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestVixInitiateFileTransfer(t *testing.T) {
	c := NewVixCommandClient()

	request := new(VixMsgListFilesRequest)

	f, err := ioutil.TempFile("", "toolbox")
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range []string{"a", "b", "c", "d", "e"} {
		_, _ = f.WriteString(strings.Repeat(s, 40))
	}

	_ = f.Close()

	name := f.Name()

	// 1st pass file exists == OK, 2nd pass does not exist == FAIL
	for _, fail := range []bool{false, true} {
		size := binary.Size(request.header) + len(name) + 1
		request.GuestPathName = name

		reply := c.Request(vixCommandInitiateFileTransferFromGuest, size, request)

		rc := vixRC(reply)

		if Trace {
			fmt.Fprintf(os.Stderr, "%s: %s\n", name, string(reply))
		}

		if fail {
			if rc == vixOK {
				t.Errorf("%s: %d", name, rc)
			}
		} else {
			if rc != vixOK {
				t.Errorf("%s: %d", name, rc)
			}

			err = os.Remove(name)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestVixInitiateFileTransferWrite(t *testing.T) {
	c := NewVixCommandClient()

	request := new(VixCommandInitiateFileTransferToGuestRequest)

	f, err := ioutil.TempFile("", "toolbox")
	if err != nil {
		t.Fatal(err)
	}

	_ = f.Close()

	name := f.Name()

	tests := []struct {
		force uint8
		fail  bool
	}{
		{0, true},  // exists == OK
		{1, false}, // exists, but overwrite == OK
		{0, false}, // does not exist == FAIL
	}

	for i, test := range tests {
		size := binary.Size(request.header) + len(name) + 1
		request.GuestPathName = name
		request.header.Overwrite = test.force

		reply := c.Request(vixCommandInitiateFileTransferToGuest, size, request)

		rc := vixRC(reply)

		if Trace {
			fmt.Fprintf(os.Stderr, "%s: %s\n", name, string(reply))
		}

		if test.fail {
			if rc == vixOK {
				t.Errorf("%d: %d", i, rc)
			}
		} else {
			if rc != vixOK {
				t.Errorf("%d: %d", i, rc)
			}
			if test.force != 0 {
				_ = os.Remove(name)
			}
		}
	}
}

func TestVixProcessHgfsPacket(t *testing.T) {
	c := NewVixCommandClient()

	c.Header.CommonFlags = vixCommandGuestReturnsBinary

	request := new(VixCommandHgfsSendPacket)

	op := new(hgfs.RequestCreateSessionV4)
	packet := new(hgfs.Packet)
	packet.Payload, _ = op.MarshalBinary()
	packet.Header.Version = hgfs.HeaderVersion
	packet.Header.Dummy = hgfs.OpNewHeader
	packet.Header.HeaderSize = uint32(binary.Size(&packet.Header))
	packet.Header.PacketSize = packet.Header.HeaderSize + uint32(len(packet.Payload))
	packet.Header.Op = hgfs.OpCreateSessionV4

	request.Packet, _ = packet.MarshalBinary()
	request.header.PacketSize = uint32(len(request.Packet))

	size := binary.Size(request.header) + int(request.header.PacketSize)

	reply := c.Request(vmxiHgfsSendPacketCommand, size, request)

	rc := vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	ix := bytes.IndexByte(reply, '#')
	reply = reply[ix+1:]
	err := packet.UnmarshalBinary(reply)
	if err != nil {
		t.Fatal(err)
	}

	if packet.Status != hgfs.StatusSuccess {
		t.Errorf("status=%d", packet.Status)
	}

	if packet.Dummy != hgfs.OpNewHeader {
		t.Errorf("dummy=%d", packet.Dummy)
	}

	session := new(hgfs.ReplyCreateSessionV4)
	err = session.UnmarshalBinary(packet.Payload)
	if err != nil {
		t.Fatal(err)
	}

	if session.NumCapabilities == 0 || int(session.NumCapabilities) != len(session.Capabilities) {
		t.Errorf("NumCapabilities=%d", session.NumCapabilities)
	}
}

func TestVixListProcessesEx(t *testing.T) {
	c := NewVixCommandClient()
	pm := c.Service.VixCommand.ProcessManager

	c.Service.VixCommand.ProcessStartCommand = func(pm *ProcessManager, r *VixMsgStartProgramRequest) (int64, error) {
		var p *Process
		switch r.ProgramPath {
		case "foo":
			p = NewProcessFunc(func(ctx context.Context, arg string) error {
				return nil
			})
		default:
			return -1, os.ErrNotExist
		}

		return pm.Start(r, p)
	}

	exec := &VixMsgStartProgramRequest{
		ProgramPath: "foo",
	}

	b, err := exec.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	size := len(b)
	reply := c.Request(vixCommandStartProgram, size, exec)
	rc := vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	r := bytes.Trim(bytes.Split(reply, []byte{' '})[2], "\x00")
	pid, _ := strconv.Atoi(string(r))

	exec.ProgramPath = "bar"
	reply = c.Request(vixCommandStartProgram, size, exec)
	rc = vixRC(reply)
	t.Log(VixError(rc).Error())
	if rc != vixFileNotFound {
		t.Fatalf("rc: %d", rc)
	}
	if vixErrorCode(os.ErrNotExist) != rc {
		t.Fatalf("rc: %d", rc)
	}

	pm.wg.Wait()

	ps := new(VixMsgListProcessesExRequest)

	ps.Pids = []int64{int64(pid)}

	b, err = ps.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	size = len(b)
	reply = c.Request(vixCommandListProcessesEx, size, ps)
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	n := bytes.Count(reply, []byte("<proc>"))
	if n != len(ps.Pids) {
		t.Errorf("ps -p %d=%d", pid, n)
	}

	kill := new(VixCommandKillProcessRequest)
	kill.header.Pid = ps.Pids[0]

	b, err = kill.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	size = len(b)
	reply = c.Request(vixCommandTerminateProcess, size, kill)
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	kill.header.Pid = 33333
	reply = c.Request(vixCommandTerminateProcess, size, kill)
	rc = vixRC(reply)
	if rc != vixNoSuchProcess {
		t.Fatalf("rc: %d", rc)
	}
}
