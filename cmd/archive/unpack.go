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

// +build !windows

package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/pkg/trace"

	docker "github.com/docker/docker/pkg/archive"
)

const (
	Success = iota
	WrongArgumentCount
	InvalidFilterSpec
	NoTarUnpackTarget
	UnpackTargetNotDirectory
	FailedChdirBeforeChroot
	FailedChroot
	FailedChdirAfterChroot
	FailedInvokeUnpack
	FailedApplyLayer
)

// Calls our custom unpack method inside of the chroot
func customUnpack(op trace.Operation) {
	setupChroot(op)
	op.Debugf("New custom unpack operation created")

	filterSpec, err := archive.DecodeFilterSpec(op, &os.Args[3])
	if err != nil {
		op.Errorf("Couldn't deserialize filterspec %s", os.Args[3])
		os.Exit(InvalidFilterSpec)
	}

	if err = archive.UnpackNoChroot(op, os.Stdin, filterSpec, "/"); err != nil {
		op.Error(err)
		os.Exit(FailedInvokeUnpack)
	}

	os.Exit(Success)
}

// Calls docker's ApplyLayer inside of the chroot
func dockerUnpack(op trace.Operation) {
	setupChroot(op)
	op.Debugf("New docker-based unpack operation created")

	// Untar the archive
	if _, err := docker.ApplyLayer("/", os.Stdin); err != nil {
		op.Errorf("Error applying layer: %s", err)
		os.Exit(FailedApplyLayer)
	}

	os.Exit(Success)

}

// Performs sanity checking (makes sure the unpack directory exists and has permissions, etc), changes directory to the unpack directory, and calls chroot on that directory. Will exit the executable if an error occurs.
func setupChroot(op trace.Operation) {
	root := os.Args[2]

	fi, err := os.Stat(root)
	if err != nil {
		// the target unpack path does not exist. We should not get here.
		op.Errorf("tar unpack target does not exist: %s", root)
		os.Exit(NoTarUnpackTarget)
	}

	if !fi.IsDir() {
		err := fmt.Errorf("unpack root target is not a directory: %s", root)
		op.Error(err)
		os.Exit(UnpackTargetNotDirectory)
	}
	op.Debugf("root exists: %s", root)

	err = os.Chdir(root)
	if err != nil {
		op.Errorf("error while chdir outside chroot: %s", err.Error())
		os.Exit(FailedChdirBeforeChroot)
	}

	err = syscall.Chroot(root)
	if err != nil {
		op.Errorf("error while chrootin': %s", err.Error())
		os.Exit(FailedChroot)
	}

	// this seems like a no-op but it is necessary to complete the chroot
	err = os.Chdir("/")
	if err != nil {
		op.Errorf("error while chdir inside chroot: %s", err.Error())
		os.Exit(FailedChdirAfterChroot)
	}
}

func main() {
	ctx := context.Background()
	op := trace.NewOperation(ctx, "Unpack") // TODO op ID? It's os.Args[1] if we need it..

	switch len(os.Args) {
	case 4:
		// When performing an unpack via docker cp, we use a custom unpack routine, which requires a FilterSpec. Thus, if we have 4 arguments (command, op ID, unpack location, filterspec), we should perform a custom unpack
		customUnpack(op)
	case 3:
		// In the case of docker pull, we use docker's ApplyLayer function and therefore do not need a filterSpec, so if we only have command, opID, and unpack location, we can proceed with using ApplyLayer
		dockerUnpack(op)
	default:
		os.Exit(WrongArgumentCount)
	}
}
