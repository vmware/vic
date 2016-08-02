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

	// VIX_USER_CREDENTIAL_NAME_PASSWORD
	vixUserCredentialNamePassword = 1

	// VIX_E_* constants from vix.h
	vixOK                         = 0
	vixFail                       = 1
	vixAuthenticationFail         = 35
	vixUnrecognizedCommandInGuest = 3025
	vixInvalidMessageHeader       = 10000
	vixInvalidMessageBody         = 10001
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
	Out *ChannelOut

	Authenticate func(VixCommandRequestHeader, []byte) error

	ProcessStartCommand func(*VixMsgStartProgramRequest) (int, error)

	handlers map[uint32]VixCommandHandler
}

type VixUserCredentialNamePassword struct {
	header struct {
		NameLength     uint32
		PasswordLength uint32
	}

	Name     string
	Password string
}

func registerVixRelayedCommandHandler(service *Service) *VixRelayedCommandHandler {
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

func vixCommandResult(rc int, err error, response []byte) []byte {
	// All Foundry tools commands return results that start with a foundry error
	// and a guest-OS-specific error (e.g. errno)
	errno := 0

	if err != nil {
		// TODO: inspect err for system error, setting errno

		response = []byte(err.Error())
	}

	return append([]byte(fmt.Sprintf("%d %d ", rc, errno)), response...)
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
		fmt.Fprintf(os.Stderr, "vix dispatch %q...\n%s\n", name, hex.Dump(data))
	}

	var header VixCommandRequestHeader
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	if header.Magic != vixCommandMagicWord {
		return vixCommandResult(vixInvalidMessageHeader, nil, nil), nil
	}

	handler, ok := c.handlers[header.OpCode]
	if !ok {
		return vixCommandResult(vixUnrecognizedCommandInGuest, nil, nil), nil
	}

	if header.OpCode != vixCommandGetToolsState {
		// Every command expect GetToolsState requires authentication
		creds := buf.Bytes()[header.BodyLength:]

		err = c.authenticate(header, creds[:header.CredentialLength])
		if err != nil {
			return vixCommandResult(vixAuthenticationFail, err, nil), nil
		}
	}

	rc := vixOK

	response, err := handler(name, header, buf.Bytes())
	if err != nil {
		rc = vixFail
	}

	return vixCommandResult(rc, err, response), nil
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

	src, _ := props.MarshalBinary()
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

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

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

		x := buf.Next(int(field.len))
		*field.val = string(bytes.TrimRight(x, "\x00"))
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
	if c.Authenticate != nil {
		return c.Authenticate(r, data)
	}

	switch r.UserCredentialType {
	case vixUserCredentialNamePassword:
		var c VixUserCredentialNamePassword

		if err := c.UnmarshalBinary(data); err != nil {
			return err
		}

		if Trace {
			fmt.Fprintf(traceLog, "ignoring credentials: %q:%q\n", c.Name, c.Password)
		}

		return nil
	default:
		return fmt.Errorf("unsupported UserCredentialType=%d", r.UserCredentialType)
	}
}

func (c *VixUserCredentialNamePassword) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(bytes.TrimRight(data, "\x00"))

	err := binary.Read(buf, binary.LittleEndian, &c.header)
	if err != nil {
		return err
	}

	str, err := base64.StdEncoding.DecodeString(string(buf.Bytes()))
	if err != nil {
		return err
	}

	c.Name = string(str[0:c.header.NameLength])
	c.Password = string(str[c.header.NameLength+1 : len(str)-1])

	return nil
}

func (c *VixUserCredentialNamePassword) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	c.header.NameLength = uint32(len(c.Name))
	c.header.PasswordLength = uint32(len(c.Password))

	_ = binary.Write(buf, binary.LittleEndian, &c.header)

	src := append([]byte(c.Name+"\x00"), []byte(c.Password+"\x00")...)

	enc := base64.StdEncoding
	pwd := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(pwd, src)
	_, _ = buf.Write(pwd)
	_ = buf.WriteByte(0)

	return buf.Bytes(), nil
}
