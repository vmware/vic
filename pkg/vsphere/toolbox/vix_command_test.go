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
	"time"

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

func (c *VixCommandClient) Request(op uint32, m encoding.BinaryMarshaler) []byte {
	b, err := m.MarshalBinary()
	if err != nil {
		panic(err)
	}

	c.Header.OpCode = op
	c.Header.BodyLength = uint32(len(b))

	var buf bytes.Buffer
	_, _ = buf.Write([]byte("\"reqname\"\x00"))
	_ = binary.Write(&buf, binary.LittleEndian, c.Header)

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
		request.GuestPathName = name

		reply := c.Request(vixCommandInitiateFileTransferFromGuest, request)

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
		force bool
		fail  bool
	}{
		{false, true},  // exists == OK
		{true, false},  // exists, but overwrite == OK
		{false, false}, // does not exist == FAIL
	}

	for i, test := range tests {
		request.GuestPathName = name
		request.header.Overwrite = test.force

		reply := c.Request(vixCommandInitiateFileTransferToGuest, request)

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
			if test.force {
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

	reply := c.Request(vmxiHgfsSendPacketCommand, request)

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

	reply := c.Request(vixCommandStartProgram, exec)
	rc := vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	r := bytes.Trim(bytes.Split(reply, []byte{' '})[2], "\x00")
	pid, _ := strconv.Atoi(string(r))

	exec.ProgramPath = "bar"
	reply = c.Request(vixCommandStartProgram, exec)
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

	reply = c.Request(vixCommandListProcessesEx, ps)
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

	reply = c.Request(vixCommandTerminateProcess, kill)
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	kill.header.Pid = 33333
	reply = c.Request(vixCommandTerminateProcess, kill)
	rc = vixRC(reply)
	if rc != vixNoSuchProcess {
		t.Fatalf("rc: %d", rc)
	}
}

func TestVixGetenv(t *testing.T) {
	c := NewVixCommandClient()

	env := os.Environ()
	key := strings.SplitN(env[0], "=", 2)[0]

	tests := []struct {
		names  []string
		expect int
	}{
		{nil, len(env)},              // all env
		{[]string{key, "ENOENT"}, 1}, // specific vars, 1 exists 1 does not
	}

	for i, test := range tests {
		env := &VixMsgReadEnvironmentVariablesRequest{
			Names: test.names,
		}
		reply := c.Request(vixCommandReadEnvVariables, env)
		rc := vixRC(reply)
		if rc != vixOK {
			t.Fatalf("%d) rc: %d", i, rc)
		}

		num := bytes.Count(reply, []byte("<ev>"))
		if num != test.expect {
			t.Errorf("%d) getenv(%v): %d", i, test.names, num)
		}
	}
}

func TestVixDirectories(t *testing.T) {
	c := NewVixCommandClient()

	mktemp := &VixMsgCreateTempFileExRequest{
		FilePrefix: "toolbox-",
	}

	// mktemp -d
	reply := c.Request(vixCommandCreateTemporaryDirectory, mktemp)
	rc := vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	dir := strings.TrimSuffix(string(reply[4:]), "\x00")

	mkdir := &VixMsgDirRequest{
		GuestPathName: dir,
	}

	// mkdir $dir == EEXIST
	reply = c.Request(vixCommandCreateDirectoryEx, mkdir)
	rc = vixRC(reply)
	if rc != vixFileAlreadyExists {
		t.Fatalf("rc: %d", rc)
	}

	// mkdir $dir/ok == OK
	mkdir.GuestPathName = dir + "/ok"
	reply = c.Request(vixCommandCreateDirectoryEx, mkdir)
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	// rm of a dir should fail, regardless if empty or not
	reply = c.Request(vixCommandDeleteGuestFileEx, &VixMsgFileRequest{
		GuestPathName: mkdir.GuestPathName,
	})
	rc = vixRC(reply)
	if rc != vixNotAFile {
		t.Errorf("rc: %d", rc)
	}

	// rmdir $dir/ok == OK
	reply = c.Request(vixCommandDeleteGuestDirectoryEx, mkdir)
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	// rmdir $dir/ok == ENOENT
	reply = c.Request(vixCommandDeleteGuestDirectoryEx, mkdir)
	rc = vixRC(reply)
	if rc != vixFileNotFound {
		t.Fatalf("rc: %d", rc)
	}

	// mkdir $dir/1/2 == ENOENT (parent directory does not exist)
	mkdir.GuestPathName = dir + "/1/2"
	reply = c.Request(vixCommandCreateDirectoryEx, mkdir)
	rc = vixRC(reply)
	if rc != vixFileNotFound {
		t.Fatalf("rc: %d", rc)
	}

	// mkdir -p $dir/1/2 == OK
	mkdir.header.Recursive = true
	reply = c.Request(vixCommandCreateDirectoryEx, mkdir)
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	// rmdir $dir == ENOTEMPTY
	mkdir.GuestPathName = dir
	mkdir.header.Recursive = false
	reply = c.Request(vixCommandDeleteGuestDirectoryEx, mkdir)
	rc = vixRC(reply)
	if rc != vixDirectoryNotEmpty {
		t.Fatalf("rc: %d", rc)
	}

	// rm -rf $dir == OK
	mkdir.header.Recursive = true
	reply = c.Request(vixCommandDeleteGuestDirectoryEx, mkdir)
	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}
}

func TestVixFiles(t *testing.T) {
	c := NewVixCommandClient()

	mktemp := &VixMsgCreateTempFileExRequest{
		FilePrefix: "toolbox-",
	}

	// mktemp -d
	reply := c.Request(vixCommandCreateTemporaryDirectory, mktemp)
	rc := vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	dir := strings.TrimSuffix(string(reply[4:]), "\x00")

	max := 12
	var total int

	// mktemp
	for i := 0; i <= max; i++ {
		mktemp = &VixMsgCreateTempFileExRequest{
			DirectoryPath: dir,
		}

		reply = c.Request(vixCommandCreateTemporaryFileEx, mktemp)
		rc = vixRC(reply)
		if rc != vixOK {
			t.Fatalf("rc: %d", rc)
		}
	}

	// name of the last file temp file we created, we'll mess around with it then delete it
	name := strings.TrimSuffix(string(reply[4:]), "\x00")

	// test ls of a single file
	ls := &VixMsgListFilesRequest{
		GuestPathName: name,
	}

	reply = c.Request(vixCommandListFiles, ls)

	rc = vixRC(reply)
	if rc != vixOK {
		t.Fatalf("rc: %d", rc)
	}

	num := bytes.Count(reply, []byte("<fxi>"))
	if num != 1 {
		t.Errorf("ls %s: %d", name, num)
	}

	num = bytes.Count(reply, []byte("<rem>0</rem>"))
	if num != 1 {
		t.Errorf("ls %s: %d", name, num)
	}

	mv := &VixCommandRenameFileExRequest{
		OldPathName: name,
		NewPathName: name + "-new",
	}

	for _, expect := range []int{vixOK, vixFileNotFound} {
		reply = c.Request(vixCommandMoveGuestFileEx, mv)
		rc = vixRC(reply)
		if rc != expect {
			t.Errorf("rc: %d", rc)
		}

		if expect == vixOK {
			// test file type is properly checked
			reply = c.Request(vixCommandMoveGuestDirectory, &VixCommandRenameFileExRequest{
				OldPathName: mv.NewPathName,
				NewPathName: name,
			})
			rc = vixRC(reply)
			if rc != vixNotADirectory {
				t.Errorf("rc: %d", rc)
			}

			// test Overwrite flag is properly checked
			reply = c.Request(vixCommandMoveGuestFileEx, &VixCommandRenameFileExRequest{
				OldPathName: mv.NewPathName,
				NewPathName: mv.NewPathName,
			})
			rc = vixRC(reply)
			if rc != vixFileAlreadyExists {
				t.Errorf("rc: %d", rc)
			}
		}
	}

	// rmdir of a file should fail
	reply = c.Request(vixCommandDeleteGuestDirectoryEx, &VixMsgDirRequest{
		GuestPathName: mv.NewPathName,
	})

	rc = vixRC(reply)
	if rc != vixNotADirectory {
		t.Errorf("rc: %d", rc)
	}

	file := &VixMsgFileRequest{
		GuestPathName: mv.NewPathName,
	}

	for _, expect := range []int{vixOK, vixFileNotFound} {
		reply = c.Request(vixCommandDeleteGuestFileEx, file)
		rc = vixRC(reply)
		if rc != expect {
			t.Errorf("rc: %d", rc)
		}
	}

	// ls again now that file is gone
	reply = c.Request(vixCommandListFiles, ls)

	rc = vixRC(reply)
	if rc != vixFileNotFound {
		t.Errorf("rc: %d", rc)
	}

	// ls
	ls = &VixMsgListFilesRequest{
		GuestPathName: dir,
	}
	ls.header.MaxResults = 5 // default is 50

	for i := 0; i < 5; i++ {
		reply = c.Request(vixCommandListFiles, ls)

		if Trace {
			fmt.Fprintf(os.Stderr, "%s: %q\n", dir, string(reply[4:]))
		}

		var rem int
		_, err := fmt.Fscanf(bytes.NewReader(reply[4:]), "<rem>%d</rem>", &rem)
		if err != nil {
			t.Fatal(err)
		}

		num := bytes.Count(reply, []byte("<fxi>"))
		total += num
		ls.header.Offset += uint64(num)

		if rem == 0 {
			break
		}
	}

	if total != max {
		t.Errorf("expected %d, got %d", max, total)
	}

	// mv $dir ${dir}-old
	mv = &VixCommandRenameFileExRequest{
		OldPathName: dir,
		NewPathName: dir + "-old",
	}

	for _, expect := range []int{vixOK, vixFileNotFound} {
		reply = c.Request(vixCommandMoveGuestDirectory, mv)
		rc = vixRC(reply)
		if rc != expect {
			t.Errorf("rc: %d", rc)
		}

		if expect == vixOK {
			// test file type is properly checked
			reply = c.Request(vixCommandMoveGuestFileEx, &VixCommandRenameFileExRequest{
				OldPathName: mv.NewPathName,
				NewPathName: dir,
			})
			rc = vixRC(reply)
			if rc != vixNotAFile {
				t.Errorf("rc: %d", rc)
			}

			// test Overwrite flag is properly checked
			reply = c.Request(vixCommandMoveGuestDirectory, &VixCommandRenameFileExRequest{
				OldPathName: mv.NewPathName,
				NewPathName: mv.NewPathName,
			})
			rc = vixRC(reply)
			if rc != vixFileAlreadyExists {
				t.Errorf("rc: %d", rc)
			}
		}
	}

	rmdir := &VixMsgDirRequest{
		GuestPathName: mv.NewPathName,
	}

	// rm -rm $dir
	for _, rmr := range []bool{false, true} {
		rmdir.header.Recursive = rmr

		reply = c.Request(vixCommandDeleteGuestDirectoryEx, rmdir)
		rc = vixRC(reply)
		if rmr {
			if rc != vixOK {
				t.Fatalf("rc: %d", rc)
			}
		} else {
			if rc != vixDirectoryNotEmpty {
				t.Fatalf("rc: %d", rc)
			}
		}
	}
}

func TestVixFileChangeAttributes(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("running as root")
	}

	c := NewVixCommandClient()

	f, err := ioutil.TempFile("", "toolbox-")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	name := f.Name()

	// touch,chown,chmod
	chattr := &VixMsgSetGuestFileAttributesRequest{
		GuestPathName: name,
	}

	h := &chattr.header

	tests := []struct {
		expect int
		f      func()
	}{
		{
			vixOK, func() {},
		},
		{
			vixOK, func() {
				h.FileOptions = vixFileAttributeSetModifyDate
				h.ModificationTime = time.Now().Unix()
			},
		},
		{
			vixOK, func() {
				h.FileOptions = vixFileAttributeSetAccessDate
				h.AccessTime = time.Now().Unix()
			},
		},
		{
			vixFileAccessError, func() {
				h.FileOptions = vixFileAttributeSetUnixOwnerid
				h.OwnerID = 0 // fails as we are not root
			},
		},
		{
			vixFileAccessError, func() {
				h.FileOptions = vixFileAttributeSetUnixGroupid
				h.GroupID = 0 // fails as we are not root
			},
		},
		{
			vixOK, func() {
				h.FileOptions = vixFileAttributeSetUnixOwnerid
				h.OwnerID = int32(os.Getuid())
			},
		},
		{
			vixOK, func() {
				h.FileOptions = vixFileAttributeSetUnixGroupid
				h.GroupID = int32(os.Getgid())
			},
		},
		{
			vixOK, func() {
				h.FileOptions = vixFileAttributeSetUnixPermissions
				h.Permissions = int32(os.FileMode(0755).Perm())
			},
		},
		{
			vixFileNotFound, func() {
				_ = os.Remove(name)
			},
		},
	}

	for i, test := range tests {
		test.f()
		reply := c.Request(vixCommandSetGuestFileAttributes, chattr)
		rc := vixRC(reply)

		if rc != test.expect {
			t.Errorf("%d: rc=%d", i, rc)
		}
	}
}
