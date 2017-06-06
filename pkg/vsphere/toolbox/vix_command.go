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
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/vmware/vic/pkg/vsphere/toolbox/hgfs"
)

const (
	vixCommandMagicWord = 0xd00d0001

	vixCommandGetToolsState = 62

	vixCommandStartProgram     = 185
	vixCommandListProcessesEx  = 186
	vixCommandReadEnvVariables = 187
	vixCommandTerminateProcess = 193

	vixCommandCreateDirectoryEx        = 178
	vixCommandMoveGuestFileEx          = 179
	vixCommandMoveGuestDirectory       = 180
	vixCommandCreateTemporaryFileEx    = 181
	vixCommandCreateTemporaryDirectory = 182
	vixCommandSetGuestFileAttributes   = 183
	vixCommandDeleteGuestFileEx        = 194
	vixCommandDeleteGuestDirectoryEx   = 195

	vixCommandListFiles                     = 177
	vmxiHgfsSendPacketCommand               = 84
	vixCommandInitiateFileTransferFromGuest = 188
	vixCommandInitiateFileTransferToGuest   = 189

	// VIX_USER_CREDENTIAL_NAME_PASSWORD
	vixUserCredentialNamePassword = 1

	// VIX_E_* constants from vix.h
	vixOK                 = 0
	vixFail               = 1
	vixFileNotFound       = 4
	vixFileAlreadyExists  = 12
	vixFileAccessError    = 13
	vixAuthenticationFail = 35

	vixUnrecognizedCommandInGuest = 3025
	vixInvalidMessageHeader       = 10000
	vixInvalidMessageBody         = 10001
	vixNotAFile                   = 20001
	vixNotADirectory              = 20002
	vixNoSuchProcess              = 20003
	vixDirectoryNotEmpty          = 20006

	// VIX_COMMAND_* constants from vixCommands.h
	vixCommandGuestReturnsBinary = 0x80

	// VIX_FILE_ATTRIBUTES_ constants from vix.h
	vixFileAttributesDirectory = 0x0001
	vixFileAttributesSymlink   = 0x0002
)

type VixError int

func (err VixError) Error() string {
	return fmt.Sprintf("vix error=%d", err)
}

// vixErrorCode does its best to map the given error to a VIX error code.
// See also: Vix_TranslateErrno
func vixErrorCode(err error) int {
	if xerr, ok := err.(VixError); ok {
		return int(xerr)
	}

	if xerr, ok := err.(*os.PathError); ok {
		if errno, ok := xerr.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ENOTEMPTY:
				return vixDirectoryNotEmpty
			}
		}
	}

	switch {
	case os.IsNotExist(err):
		return vixFileNotFound
	case os.IsExist(err):
		return vixFileAlreadyExists
	case os.IsPermission(err):
		return vixFileAccessError
	default:
		return vixFail
	}
}

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

type VixCommandKillProcessRequest struct {
	VixCommandRequestHeader

	header struct {
		Pid     int64
		Options uint32
	}
}

type VixMsgListProcessesExRequest struct {
	VixCommandRequestHeader

	header struct {
		Key     uint32
		Offset  uint32
		NumPids uint32
	}

	Pids []int64
}

type VixMsgReadEnvironmentVariablesRequest struct {
	VixCommandRequestHeader

	header struct {
		NumNames    uint32
		NamesLength uint32
	}

	Names []string
}

type VixMsgCreateTempFileExRequest struct {
	VixCommandRequestHeader

	header struct {
		Options             int32
		FilePrefixLength    uint32
		FileSuffixLength    uint32
		DirectoryPathLength uint32
		PropertyListLength  uint32
	}

	FilePrefix    string
	FileSuffix    string
	DirectoryPath string
}

type VixMsgFileRequest struct {
	VixCommandRequestHeader

	header struct {
		FileOptions         int32
		GuestPathNameLength uint32
	}

	GuestPathName string
}

type VixMsgDirRequest struct {
	VixCommandRequestHeader

	header struct {
		FileOptions          int32
		GuestPathNameLength  uint32
		FilePropertiesLength uint32
		Recursive            bool
	}

	GuestPathName string
}

type VixCommandRenameFileExRequest struct {
	VixCommandRequestHeader

	header struct {
		CopyFileOptions      int32
		OldPathNameLength    uint32
		NewPathNameLength    uint32
		FilePropertiesLength uint32
		Overwrite            bool
	}

	OldPathName string
	NewPathName string
}

type VixMsgListFilesRequest struct {
	VixCommandRequestHeader

	header struct {
		FileOptions         int32
		GuestPathNameLength uint32
		PatternLength       uint32
		Index               int32
		MaxResults          int32
		Offset              uint64
	}

	GuestPathName string
	Pattern       string
}

type VixMsgSetGuestFileAttributesRequest struct {
	VixCommandRequestHeader

	header struct {
		FileOptions         int32
		AccessTime          int64
		ModificationTime    int64
		OwnerID             int32
		GroupID             int32
		Permissions         int32
		Hidden              bool
		ReadOnly            bool
		GuestPathNameLength uint32
	}

	GuestPathName string
}

type VixCommandHgfsSendPacket struct {
	VixCommandRequestHeader

	header struct {
		PacketSize uint32
		Timeout    int32
	}

	Packet []byte
}

type VixCommandInitiateFileTransferToGuestRequest struct {
	VixCommandRequestHeader

	header struct {
		Options             int32
		GuestPathNameLength uint32
		Overwrite           bool
	}

	GuestPathName string
}

type VixCommandHandler func(string, VixCommandRequestHeader, []byte) ([]byte, error)

type VixRelayedCommandHandler struct {
	Out *ChannelOut

	ProcessManager *ProcessManager

	Authenticate func(VixCommandRequestHeader, []byte) error

	ProcessStartCommand func(*ProcessManager, *VixMsgStartProgramRequest) (int64, error)

	handlers map[uint32]VixCommandHandler

	FileServer *hgfs.Server
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
		Out:            service.out,
		ProcessManager: NewProcessManager(),
		handlers:       make(map[uint32]VixCommandHandler),
	}

	service.RegisterHandler("Vix_1_Relayed_Command", handler.Dispatch)

	handler.RegisterHandler(vixCommandGetToolsState, handler.GetToolsState)

	handler.RegisterHandler(vixCommandStartProgram, handler.StartCommand)
	handler.RegisterHandler(vixCommandTerminateProcess, handler.KillProcess)
	handler.RegisterHandler(vixCommandListProcessesEx, handler.ListProcessesEx)
	handler.RegisterHandler(vixCommandReadEnvVariables, handler.ReadEnvironmentVariables)

	handler.RegisterHandler(vixCommandCreateTemporaryFileEx, handler.CreateTemporaryFileEx)
	handler.RegisterHandler(vixCommandCreateTemporaryDirectory, handler.CreateTemporaryDirectory)
	handler.RegisterHandler(vixCommandDeleteGuestFileEx, handler.DeleteFile)
	handler.RegisterHandler(vixCommandCreateDirectoryEx, handler.CreateDirectory)
	handler.RegisterHandler(vixCommandDeleteGuestDirectoryEx, handler.DeleteDirectory)
	handler.RegisterHandler(vixCommandMoveGuestFileEx, handler.MoveFile)
	handler.RegisterHandler(vixCommandMoveGuestDirectory, handler.MoveDirectory)
	handler.RegisterHandler(vixCommandListFiles, handler.ListFiles)
	handler.RegisterHandler(vixCommandSetGuestFileAttributes, handler.SetGuestFileAttributes)
	handler.RegisterHandler(vixCommandInitiateFileTransferFromGuest, handler.InitiateFileTransferFromGuest)
	handler.RegisterHandler(vixCommandInitiateFileTransferToGuest, handler.InitiateFileTransferToGuest)

	handler.RegisterHandler(vmxiHgfsSendPacketCommand, handler.ProcessHgfsPacket)

	handler.ProcessStartCommand = handler.ExecCommandStart

	return handler
}

func vixCommandResult(header VixCommandRequestHeader, rc int, err error, response []byte) []byte {
	// All Foundry tools commands return results that start with a foundry error
	// and a guest-OS-specific error (e.g. errno)
	errno := 0

	if err != nil {
		// TODO: inspect err for system error, setting errno

		response = []byte(err.Error())
	}

	buf := bytes.NewBufferString(fmt.Sprintf("%d %d ", rc, errno))

	if header.CommonFlags&vixCommandGuestReturnsBinary != 0 {
		// '#' delimits end of ascii and the start of the binary data (see ToolsDaemonTcloReceiveVixCommand)
		_ = buf.WriteByte('#')
	}

	_, _ = buf.Write(response)

	if header.CommonFlags&vixCommandGuestReturnsBinary == 0 {
		// this is not binary data, so it should be a NULL terminated string (see ToolsDaemonTcloReceiveVixCommand)
		_ = buf.WriteByte(0)
	}

	return buf.Bytes()
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
		return vixCommandResult(header, vixInvalidMessageHeader, nil, nil), nil
	}

	handler, ok := c.handlers[header.OpCode]
	if !ok {
		return vixCommandResult(header, vixUnrecognizedCommandInGuest, nil, nil), nil
	}

	if header.OpCode != vixCommandGetToolsState {
		// Every command expect GetToolsState requires authentication
		creds := buf.Bytes()[header.BodyLength:]

		err = c.authenticate(header, creds[:header.CredentialLength])
		if err != nil {
			return vixCommandResult(header, vixAuthenticationFail, err, nil), nil
		}
	}

	rc := vixOK

	response, err := handler(name, header, buf.Bytes())
	if err != nil {
		rc = vixErrorCode(err)
	}

	return vixCommandResult(header, rc, err, response), nil
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
		NewBoolProperty(VixPropertyGuestTerminateProcessEnabled, true),
		NewBoolProperty(VixPropertyGuestListProcessesEnabled, true),
		NewBoolProperty(VixPropertyGuestReadEnvironmentVariableEnabled, true),
		NewBoolProperty(VixPropertyGuestMakeDirectoryEnabled, true),
		NewBoolProperty(VixPropertyGuestDeleteFileEnabled, true),
		NewBoolProperty(VixPropertyGuestDeleteDirectoryEnabled, true),
		NewBoolProperty(VixPropertyGuestMoveDirectoryEnabled, true),
		NewBoolProperty(VixPropertyGuestMoveFileEnabled, true),
		NewBoolProperty(VixPropertyGuestCreateTempFileEnabled, true),
		NewBoolProperty(VixPropertyGuestCreateTempDirectoryEnabled, true),
		NewBoolProperty(VixPropertyGuestListFilesEnabled, true),
		NewBoolProperty(VixPropertyGuestChangeFileAttributesEnabled, true),
		NewBoolProperty(VixPropertyGuestInitiateFileTransferFromGuestEnabled, true),
		NewBoolProperty(VixPropertyGuestInitiateFileTransferToGuestEnabled, true),
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

	pid, err := c.ProcessStartCommand(c.ProcessManager, r)
	if err != nil {
		return nil, err
	}

	return append([]byte(fmt.Sprintf("%d", pid)), 0), nil
}

func (c *VixRelayedCommandHandler) ExecCommandStart(m *ProcessManager, r *VixMsgStartProgramRequest) (int64, error) {
	return m.Start(r, NewProcess())
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixCommandKillProcessRequest) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixCommandKillProcessRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	return binary.Read(buf, binary.LittleEndian, &r.header)
}

func (c *VixRelayedCommandHandler) KillProcess(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixCommandKillProcessRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	if c.ProcessManager.Kill(r.header.Pid) {
		return nil, err
	}

	// TODO: could kill process started outside of toolbox

	return nil, VixError(vixNoSuchProcess)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixMsgListProcessesExRequest) MarshalBinary() ([]byte, error) {
	r.header.NumPids = uint32(len(r.Pids))

	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	for _, pid := range r.Pids {
		_ = binary.Write(buf, binary.LittleEndian, &pid)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixMsgListProcessesExRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	r.Pids = make([]int64, r.header.NumPids)

	for i := uint32(0); i < r.header.NumPids; i++ {
		err := binary.Read(buf, binary.LittleEndian, &r.Pids[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *VixRelayedCommandHandler) ListProcessesEx(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgListProcessesExRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	state := c.ProcessManager.ListProcesses(r.Pids)

	return state, nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixMsgCreateTempFileExRequest) MarshalBinary() ([]byte, error) {
	var fields []string

	add := func(s string, l *uint32) {
		*l = uint32(len(s)) // NOTE: NULL byte is not included in the length fields on the wire
		fields = append(fields, s)
	}

	add(r.FilePrefix, &r.header.FilePrefixLength)
	add(r.FileSuffix, &r.header.FileSuffixLength)
	add(r.DirectoryPath, &r.header.DirectoryPathLength)

	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	for _, val := range fields {
		_, _ = buf.Write([]byte(val))
		_ = buf.WriteByte(0)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixMsgCreateTempFileExRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	fields := []struct {
		len uint32
		val *string
	}{
		{r.header.FilePrefixLength, &r.FilePrefix},
		{r.header.FileSuffixLength, &r.FileSuffix},
		{r.header.DirectoryPathLength, &r.DirectoryPath},
	}

	for _, field := range fields {
		field.len++ // NOTE: NULL byte is not included in the length fields on the wire

		x := buf.Next(int(field.len))
		*field.val = string(bytes.TrimRight(x, "\x00"))
	}

	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixMsgReadEnvironmentVariablesRequest) MarshalBinary() ([]byte, error) {
	var env bytes.Buffer

	if n := len(r.Names); n != 0 {
		for _, e := range r.Names {
			_, _ = env.Write([]byte(e))
			_ = env.WriteByte(0)
		}
		r.header.NumNames = uint32(n)
		r.header.NamesLength = uint32(env.Len())
	}

	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	if r.header.NamesLength != 0 {
		_, _ = buf.Write(env.Bytes())
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixMsgReadEnvironmentVariablesRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	for i := 0; i < int(r.header.NumNames); i++ {
		env, rerr := buf.ReadString(0)
		if rerr != nil {
			return rerr
		}

		env = env[:len(env)-1] // discard NULL terminator
		r.Names = append(r.Names, env)
	}

	return nil
}

func (c *VixRelayedCommandHandler) ReadEnvironmentVariables(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgReadEnvironmentVariablesRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if len(r.Names) == 0 {
		for _, e := range os.Environ() {
			_, _ = buf.WriteString(fmt.Sprintf("<ev>%s</ev>", xmlEscape.Replace(e)))
		}
	} else {
		for _, key := range r.Names {
			val := os.Getenv(key)
			if val == "" {
				continue
			}
			_, _ = buf.WriteString(fmt.Sprintf("<ev>%s=%s</ev>", xmlEscape.Replace(key), xmlEscape.Replace(val)))
		}
	}

	return buf.Bytes(), nil
}

func (c *VixRelayedCommandHandler) CreateTemporaryFileEx(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgCreateTempFileExRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	f, err := ioutil.TempFile(r.DirectoryPath, r.FilePrefix+"vmware")
	if err != nil {
		return nil, err
	}

	_ = f.Close()

	return []byte(f.Name()), nil
}

func (c *VixRelayedCommandHandler) CreateTemporaryDirectory(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgCreateTempFileExRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	name, err := ioutil.TempDir(r.DirectoryPath, r.FilePrefix+"vmware")
	if err != nil {
		return nil, err
	}

	return []byte(name), nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixMsgFileRequest) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	r.header.GuestPathNameLength = uint32(len(r.GuestPathName))

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	_, _ = buf.WriteString(r.GuestPathName)
	_ = buf.WriteByte(0)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixMsgFileRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	name := buf.Next(int(r.header.GuestPathNameLength))
	r.GuestPathName = string(bytes.TrimRight(name, "\x00"))

	return nil
}

func (c *VixRelayedCommandHandler) DeleteFile(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgFileRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(r.GuestPathName)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, VixError(vixNotAFile)
	}

	err = os.Remove(r.GuestPathName)

	return nil, err
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixMsgDirRequest) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	r.header.GuestPathNameLength = uint32(len(r.GuestPathName))

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	_, _ = buf.WriteString(r.GuestPathName)
	_ = buf.WriteByte(0)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixMsgDirRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	name := buf.Next(int(r.header.GuestPathNameLength))
	r.GuestPathName = string(bytes.TrimRight(name, "\x00"))

	return nil
}

func (c *VixRelayedCommandHandler) DeleteDirectory(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgDirRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(r.GuestPathName)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, VixError(vixNotADirectory)
	}

	if r.header.Recursive {
		err = os.RemoveAll(r.GuestPathName)
	} else {
		err = os.Remove(r.GuestPathName)
	}

	return nil, err
}

func (c *VixRelayedCommandHandler) CreateDirectory(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgDirRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	mkdir := os.Mkdir

	if r.header.Recursive {
		mkdir = os.MkdirAll
	}

	err = mkdir(r.GuestPathName, 0700)

	return nil, err
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixCommandRenameFileExRequest) MarshalBinary() ([]byte, error) {
	var fields []string

	add := func(s string, l *uint32) {
		*l = uint32(len(s)) // NOTE: NULL byte is not included in the length fields on the wire
		fields = append(fields, s)
	}

	add(r.OldPathName, &r.header.OldPathNameLength)
	add(r.NewPathName, &r.header.NewPathNameLength)

	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	for _, val := range fields {
		_, _ = buf.Write([]byte(val))
		_ = buf.WriteByte(0)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixCommandRenameFileExRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	fields := []struct {
		len uint32
		val *string
	}{
		{r.header.OldPathNameLength, &r.OldPathName},
		{r.header.NewPathNameLength, &r.NewPathName},
	}

	for _, field := range fields {
		field.len++ // NOTE: NULL byte is not included in the length fields on the wire

		x := buf.Next(int(field.len))
		*field.val = string(bytes.TrimRight(x, "\x00"))
	}

	return nil
}

func (c *VixRelayedCommandHandler) MoveDirectory(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixCommandRenameFileExRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(r.OldPathName)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, VixError(vixNotADirectory)
	}

	if !r.header.Overwrite {
		info, err = os.Stat(r.NewPathName)
		if err == nil {
			return nil, VixError(vixFileAlreadyExists)
		}
	}

	return nil, os.Rename(r.OldPathName, r.NewPathName)
}

func (c *VixRelayedCommandHandler) MoveFile(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixCommandRenameFileExRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(r.OldPathName)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, VixError(vixNotAFile)
	}

	if !r.header.Overwrite {
		info, err = os.Stat(r.NewPathName)
		if err == nil {
			return nil, VixError(vixFileAlreadyExists)
		}
	}

	return nil, os.Rename(r.OldPathName, r.NewPathName)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixMsgListFilesRequest) MarshalBinary() ([]byte, error) {
	var fields []string

	add := func(s string, l *uint32) {
		if n := len(s); n != 0 {
			*l = uint32(n) + 1
			fields = append(fields, s)
		}
	}

	add(r.GuestPathName, &r.header.GuestPathNameLength)
	add(r.Pattern, &r.header.PatternLength)

	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	for _, val := range fields {
		_, _ = buf.Write([]byte(val))
		_ = buf.WriteByte(0)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixMsgListFilesRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	fields := []struct {
		len uint32
		val *string
	}{
		{r.header.GuestPathNameLength, &r.GuestPathName},
		{r.header.PatternLength, &r.Pattern},
	}

	for _, field := range fields {
		if field.len == 0 {
			continue
		}

		x := buf.Next(int(field.len))
		*field.val = string(bytes.TrimRight(x, "\x00"))
	}

	return nil
}

func (c *VixRelayedCommandHandler) ListFiles(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgListFilesRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(r.GuestPathName)
	if err != nil {
		return nil, err
	}

	var files []os.FileInfo

	if info.IsDir() {
		files, err = ioutil.ReadDir(r.GuestPathName)
		if err != nil {
			return nil, err
		}
	} else {
		files = append(files, info)
	}

	offset := r.header.Offset + uint64(r.header.Index)
	total := uint64(len(files)) - offset

	files = files[offset:]

	var remaining uint64

	if r.header.MaxResults > 0 && total > uint64(r.header.MaxResults) {
		remaining = total - uint64(r.header.MaxResults)
		files = files[:r.header.MaxResults]
	}

	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("<rem>%d</rem>", remaining))

	for _, info = range files {
		buf.WriteString(fileExtendedInfoFormat(info))
	}

	return buf.Bytes(), nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixMsgSetGuestFileAttributesRequest) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	r.header.GuestPathNameLength = uint32(len(r.GuestPathName))

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	_, _ = buf.WriteString(r.GuestPathName)
	_ = buf.WriteByte(0)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixMsgSetGuestFileAttributesRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	name := buf.Next(int(r.header.GuestPathNameLength))
	r.GuestPathName = string(bytes.TrimRight(name, "\x00"))

	return nil
}

// SetGuestFileAttributes flags as defined in vixOpenSource.h
const (
	vixFileAttributeSetAccessDate      = 0x0001
	vixFileAttributeSetModifyDate      = 0x0002
	vixFileAttributeSetReadonly        = 0x0004
	vixFileAttributeSetHidden          = 0x0008
	vixFileAttributeSetUnixOwnerid     = 0x0010
	vixFileAttributeSetUnixGroupid     = 0x0020
	vixFileAttributeSetUnixPermissions = 0x0040
)

func (r *VixMsgSetGuestFileAttributesRequest) isSet(opt int32) bool {
	return r.header.FileOptions&opt == opt
}

func (r *VixMsgSetGuestFileAttributesRequest) chtimes() error {
	var mtime, atime *time.Time

	if r.isSet(vixFileAttributeSetModifyDate) {
		t := time.Unix(r.header.ModificationTime, 0)
		mtime = &t
	}

	if r.isSet(vixFileAttributeSetAccessDate) {
		t := time.Unix(r.header.AccessTime, 0)
		atime = &t
	}

	if mtime == nil && atime == nil {
		return nil
	}

	info, err := os.Stat(r.GuestPathName)
	if err != nil {
		return err
	}

	if mtime == nil {
		t := info.ModTime()
		mtime = &t
	}

	if atime == nil {
		t := info.ModTime()
		atime = &t
	}

	return os.Chtimes(r.GuestPathName, *atime, *mtime)
}

func (r *VixMsgSetGuestFileAttributesRequest) chown() error {
	uid := -1
	gid := -1

	if r.isSet(vixFileAttributeSetUnixOwnerid) {
		uid = int(r.header.OwnerID)
	}

	if r.isSet(vixFileAttributeSetUnixGroupid) {
		gid = int(r.header.GroupID)
	}

	if uid == -1 && gid == -1 {
		return nil
	}

	return os.Chown(r.GuestPathName, uid, gid)
}

func (r *VixMsgSetGuestFileAttributesRequest) chmod() error {
	if r.isSet(vixFileAttributeSetUnixPermissions) {
		return os.Chmod(r.GuestPathName, os.FileMode(r.header.Permissions).Perm())
	}

	return nil
}

func (c *VixRelayedCommandHandler) SetGuestFileAttributes(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgSetGuestFileAttributesRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	for _, set := range []func() error{r.chtimes, r.chown, r.chmod} {
		err = set()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixCommandHgfsSendPacket) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	_, _ = buf.Write(r.Packet)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixCommandHgfsSendPacket) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	r.Packet = buf.Next(int(r.header.PacketSize))

	return nil
}

func (c *VixRelayedCommandHandler) InitiateFileTransferFromGuest(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixMsgListFilesRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(r.GuestPathName)
	if err != nil {
		return nil, err
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		return nil, errors.New("VIX_E_INVALID_ARG")
	}

	if info.IsDir() {
		return nil, errors.New("VIX_E_NOT_A_FILE")
	}

	return []byte(fileExtendedInfoFormat(info)), nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (r *VixCommandInitiateFileTransferToGuestRequest) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	r.header.GuestPathNameLength = uint32(len(r.GuestPathName))

	_ = binary.Write(buf, binary.LittleEndian, &r.header)

	_, _ = buf.WriteString(r.GuestPathName)
	_ = buf.WriteByte(0)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (r *VixCommandInitiateFileTransferToGuestRequest) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &r.header)
	if err != nil {
		return err
	}

	name := buf.Next(int(r.header.GuestPathNameLength))
	r.GuestPathName = string(bytes.TrimRight(name, "\x00"))

	return nil
}

func (c *VixRelayedCommandHandler) InitiateFileTransferToGuest(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixCommandInitiateFileTransferToGuestRequest{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(r.GuestPathName)
	if err == nil {
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil, errors.New("VIX_E_INVALID_ARG")
		}

		if info.IsDir() {
			return nil, errors.New("VIX_E_NOT_A_FILE")
		}

		if !r.header.Overwrite {
			return nil, errors.New("VIX_E_FILE_ALREADY_EXISTS")
		}
	} else {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return nil, nil
}

func (c *VixRelayedCommandHandler) ProcessHgfsPacket(_ string, header VixCommandRequestHeader, data []byte) ([]byte, error) {
	r := &VixCommandHgfsSendPacket{
		VixCommandRequestHeader: header,
	}

	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return c.FileServer.Dispatch(r.Packet)
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
