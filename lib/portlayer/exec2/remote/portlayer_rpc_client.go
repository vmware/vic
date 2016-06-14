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

package remote

import (
	"net/rpc"

	"github.com/vmware/vic/lib/portlayer/exec2"
)

type PortLayerRPCClient struct {
}

var serverAddress = "localhost"
var client *rpc.Client

func (*PortLayerRPCClient) Connect() error {
	// Ignore Init args on the client - that is the server's responsibilty
	var err error
	client, err = rpc.DialHTTP("tcp", serverAddress+":1234")
	return err
}

type CreateArgs struct {
	Name string
}

func (*PortLayerRPCClient) CreateContainer(name string) (exec2.Handle, error) {
	args := &CreateArgs{Name: name}
	var reply exec2.Handle
	err := client.Call("PortLayerRpcServer.CreateContainer", args, &reply)
	return reply, err
}

func (*PortLayerRPCClient) GetHandle(cid exec2.ID) (exec2.Handle, error) {
	var reply exec2.Handle
	err := client.Call("PortLayerRpcServer.GetHandle", cid, &reply)
	return reply, err
}

type CopyToArgs struct {
	Handle    exec2.Handle
	TargetDir string
	Filename  string
	Data      []byte
}

func (*PortLayerRPCClient) CopyTo(handle exec2.Handle, targetDir string, fname string, data []byte) (exec2.Handle, error) {
	args := &CopyToArgs{Handle: handle, TargetDir: targetDir, Filename: fname, Data: data}
	var reply exec2.Handle
	err := client.Call("PortLayerRpcServer.CopyTo", args, &reply)
	return reply, err
}

type SetEntryPointArgs struct {
	Handle   exec2.Handle
	WorkDir  string
	ExecPath string
	Args     string
}

func (*PortLayerRPCClient) SetEntryPoint(handle exec2.Handle, workDir string, execPath string, args string) (exec2.Handle, error) {
	epArgs := &SetEntryPointArgs{Handle: handle, WorkDir: workDir, ExecPath: execPath, Args: args}
	var reply exec2.Handle
	err := client.Call("PortLayerRpcServer.SetEntryPoint", epArgs, &reply)
	return reply, err
}

type ExecProcessArgs struct {
	Handle   exec2.Handle
	ExecPath string
	ExecArgs string
}

func (*PortLayerRPCClient) ExecProcess(handle exec2.Handle, execPath string, execArgs string) (exec2.Handle, error) {
	args := &ExecProcessArgs{Handle: handle, ExecPath: execPath, ExecArgs: execArgs}
	var reply exec2.Handle
	err := client.Call("PortLayerRpcServer.ExecProcess", args, &reply)
	return reply, err
}

type SetLimitsArgs struct {
	Handle   exec2.Handle
	MemoryMb int
	CPUMhz   int
}

func (*PortLayerRPCClient) SetLimits(handle exec2.Handle, memoryMb int, cpuMhz int) (exec2.Handle, error) {
	args := &SetLimitsArgs{Handle: handle, MemoryMb: memoryMb, CPUMhz: cpuMhz}
	var reply exec2.Handle
	err := client.Call("PortLayerRpcServer.SetLimits", args, &reply)
	return reply, err
}

type SetRunStateArgs struct {
	Handle   exec2.Handle
	RunState exec2.RunState
}

func (*PortLayerRPCClient) SetRunState(handle exec2.Handle, runState exec2.RunState) (exec2.Handle, error) {
	args := &SetRunStateArgs{Handle: handle, RunState: runState}
	var reply exec2.Handle
	err := client.Call("PortLayerRpcServer.SetRunState", args, &reply)
	return reply, err
}

type CommitArgs struct {
	Handle exec2.Handle
}

func (*PortLayerRPCClient) Commit(handle exec2.Handle) (exec2.ID, error) {
	args := &CommitArgs{Handle: handle}
	var reply exec2.ID
	err := client.Call("PortLayerRpcServer.Commit", args, &reply)
	return reply, err
}

func (*PortLayerRPCClient) DestroyContainer(cid exec2.ID) error {
	/* Ignore the reply */
	var reply exec2.ID
	return client.Call("PortLayerRpcServer.DestroyContainer", cid, &reply)
}
