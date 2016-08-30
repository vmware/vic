// Copyright 2016 VMware, Inc. All Rights Reserved.
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
	"encoding"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"
	"testing"
)

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

	// handler not set
	buf = append(marshal(request), creds...)
	reply, _ = vix.Dispatch(buf)
	rc = vixRC(reply)
	if rc != vixFail {
		t.Fatalf("%q", reply)
	}

	vix.ProcessStartCommand = func(r *VixMsgStartProgramRequest) (int, error) {
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
