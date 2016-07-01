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
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"
)

const (
	vixCommandMagicWord = 0xd00d0001

	vixCommandGetToolsState = 62
	vixCommandStartProgram  = 185
)

type VixMsgHeader struct {
	Magic          uint32
	MessageVersion uint16

	TotalMessageLength uint32
	HeaderLength       uint32
	BodyLength         uint32
	CredentialLength   uint32

	CommonFlags uint8
}

type VixCommandRequestHeader struct {
	VixMsgHeader

	OpCode       uint32
	RequestFlags uint32

	TimeOut uint32

	Cookie         uint64
	ClientHandleID uint32

	UserCredentialType uint32
}

type VixMsgStartProgramRequest struct {
	VixCommandRequestHeader

	header struct {
		StartMinimized    uint8
		ProgramPathLength uint32
		ArgumentsLength   uint32
		WorkingDirLength  uint32
		NumEnvVars        uint32
		EnvVarLength      uint32
	}

	ProgramPath string
	Arguments   string
	WorkingDir  string
	EnvVars     []string
}

type VixCommandHandler func(string, VixCommandRequestHeader, []byte) ([]byte, error)

type VixRelayedCommandHandler struct {
	Out Channel

	ProcessStartCommand func(*VixMsgStartProgramRequest) (int, error)

	handlers map[uint32]VixCommandHandler
}

func RegisterVixRelayedCommandHandler(service *Service) *VixRelayedCommandHandler {
	handler := &VixRelayedCommandHandler{
		Out:      service.out,
		handlers: make(map[uint32]VixCommandHandler),
	}

	service.RegisterHandler("Vix_1_Relayed_Command", handler.Dispatch)

	handler.RegisterHandler(vixCommandGetToolsState, handler.GetToolsState)

	handler.RegisterHandler(vixCommandStartProgram, handler.StartCommand)

	handler.ProcessStartCommand = handler.ExecCommandStart

	return handler
}

func (c *VixRelayedCommandHandler) Dispatch(data []byte) ([]byte, error) {
	// See ToolsDaemonTcloGetQuotedString
	if data[0] == '"' {
		data = data[1:]
	}

	var name string

	ix := bytes.IndexByte(data, '"')
	if ix > 0 {
		name = string(data[:ix])
		data = data[ix+1:]
	}
	// skip the NULL
	if data[0] == 0 {
		data = data[1:]
	}

	if Trace {
		fmt.Fprintf(os.Stderr, "vix dispatch '%s'...\n%s\n", name, hex.Dump(data))
	}

	var header VixCommandRequestHeader
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	if header.Magic != vixCommandMagicWord {
		return nil, errors.New("VIX_E_INVALID_MESSAGE_HEADER")
	}

	handler, ok := c.handlers[header.OpCode]
	if !ok {
		return nil, errors.New("VIX_E_UNRECOGNIZED_COMMAND_IN_GUEST")
	}

	creds := buf.Bytes()[header.BodyLength:]
	// TODO: ignoring credentials for now
	_ = c.authenticate(header, creds)

	// All Foundry tools commands return results that start with a foundry error
	// and a guest-OS-specific error (e.g. errno)
	var rc, errno int

	response, err := handler(name, header, buf.Bytes())
	if err != nil {
		// TODO: support the other 10 million VIX_E_* errors
		rc = 1 // VIX_E_FAIL
		// TODO: inspect err for system error, setting errno
	}

	return append([]byte(fmt.Sprintf("%d %d ", rc, errno)), response...), nil
}

func (c *VixRelayedCommandHandler) RegisterHandler(op uint32, handler VixCommandHandler) {
	c.handlers[op] = handler
}

func (c *VixRelayedCommandHandler) GetToolsState(_ string, _ VixCommandRequestHeader, _ []byte) ([]byte, error) {
	hostname, _ := os.Hostname()
	osname := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	// Note that vmtoolsd sends back 40 or so of these properties, sticking with the minimal set for now.
	props := VixPropertyList{
		NewStringProperty(VixPropertyGuestOsVersion, osname),
		NewStringProperty(VixPropertyGuestOsVersionShort, osname),
		NewStringProperty(VixPropertyGuestToolsProductNam, "VMware Tools (Go)"),
		NewStringProperty(VixPropertyGuestToolsVersion, "10.0.5 build-3227872 (Compatible)"),
		NewStringProperty(VixPropertyGuestName, hostname),
		NewInt32Property(VixPropertyGuestToolsAPIOptions, 0x0001), // TODO: const VIX_TOOLSFEATURE_SUPPORT_GET_HANDLE_STATE
		NewInt32Property(VixPropertyGuestOsFamily, 1),             // TODO: const GUEST_OS_FAMILY_*
		NewBoolProperty(VixPropertyGuestStartProgramEnabled, true),
	}

	src, err := props.MarshalBinary()
	if err != nil {
		return nil, err
	}

	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)

	return buf, nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixMsgStartProgramRequest) MarshalBinary() ([]byte, error) {
	var env bytes.Buffer

	if n := len(r.EnvVars); n != 0 {
		for _, e := range r.EnvVars {
			_, _ = env.Write([]byte(e))
			_ = env.WriteByte(0)
		}
		r.header.NumEnvVars = uint32(n)
		r.header.EnvVarLength = uint32(env.Len())
	}

	var fields []string

	add := func(s string, l *uint32) {
		if n := len(s); n != 0 {
			*l = uint32(n) + 1
			fields = append(fields, s)
		}
	}

	add(r.ProgramPath, &r.header.ProgramPathLength)
	add(r.Arguments, &r.header.ArgumentsLength)
	add(r.WorkingDir, &r.header.WorkingDirLength)

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return nil, err
	}

	for _, val := range fields {
		_, _ = buf.Write([]byte(val))
		_ = buf.WriteByte(0)
	}

	if r.header.EnvVarLength != 0 {
		_, _ = buf.Write(env.Bytes())
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixMsgStartProgramRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	fields := []struct {
		len uint32
		val *string
	}{
		{r.header.ProgramPathLength, &r.ProgramPath},
		{r.header.ArgumentsLength, &r.Arguments},
		{r.header.WorkingDirLength, &r.WorkingDir},
	}

	for _, field := range fields {
		if field.len == 0 {
			continue
		}

		x := buf.Next(int(field.len) - 1)
		*field.val = string(x)

		_, err = buf.ReadByte() // discard NULL terminator
		if err != nil {
			return err
		}
	}

	for i := 0; i < int(r.header.NumEnvVars); i++ {
		env, rerr := buf.ReadString(0)
		if rerr != nil {
			return rerr
		}

		env = env[:len(env)-1] // discard NULL terminator
		r.EnvVars = append(r.EnvVars, env)
	}

	return nil
}

func (c *VixRelayedCommandHandler) StartCommand(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgStartProgramRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	pid, err := c.ProcessStartCommand(r)
	if err != nil {
		return nil, err
	}

	return append([]byte(fmt.Sprintf("%d", pid)), 0), nil
}

func (c *VixRelayedCommandHandler) ExecCommandStart(r *VixMsgStartProgramRequest) (int, error) {
	// TODO: we could map to exec.Command(...).Start here
	return -1, errors.New("not implemented")
}

func (c *VixRelayedCommandHandler) authenticate(r VixCommandRequestHeader, data []byte) error {
	buf := bytes.NewBuffer(data)

	// TODO: ignoring credentials for now, but this is how-to decode...
	if r.UserCredentialType == 1 { // VIX_USER_CREDENTIAL_NAME_PASSWORD
		pw := struct {
			NameLength     uint32
			PasswordLength uint32
		}{}

		err := binary.Read(buf, binary.LittleEndian, &pw)
		if err != nil {
			return err
		}

		length := int(r.CredentialLength - uint32(int32Size*2) - 1)

		str, err := base64.StdEncoding.DecodeString(string(buf.Next(length)))
		if err != nil {
			return err
		}

		creds := struct {
			Name, Password string
		}{
			string(str[0:pw.NameLength]),
			string(str[pw.NameLength+1 : len(str)-1]),
		}

		if Trace {
			fmt.Fprintf(os.Stderr, "ignoring credentials: '%s:%s'\n", creds.Name, creds.Password)
		}
	}

	return nil
}
