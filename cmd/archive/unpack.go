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
)

func main() {
	ctx := context.Background()
	op := trace.NewOperation(ctx, "Unpack") // TODO op ID?
	op.Infof("XXX New unpack operation created")

	var args []string
	for _, arg := range os.Args {
		args = append(args, fmt.Sprintf("%s, ", arg))
	}
	op.Infof("%s", args)

	if len(os.Args) < 4 {
		op.Errorf("XXX Not enough arguments passed. Arguments were: %s", args)
		os.Exit(5)
	}

	root := os.Args[2]
	op.Infof("XXX root inside binary %s", root)

	filterSpec, err := archive.DecodeFilterSpec(op, &os.Args[3])
	if err != nil {
		op.Errorf("Couldn't deserialize filterspec %s", os.Args[3])
		os.Exit(11)
	}

	op.Infof("XXX filterSpec decoded %s", filterSpec)

	fi, err := os.Stat(root)
	if err != nil {
		// the target unpack path does not exist. We should not get here.
		op.Errorf("tar unpack target does not exist: %s", root)
		os.Exit(9)
	}

	if !fi.IsDir() {
		err := fmt.Errorf("unpack root target is not a directory: %s", root)
		op.Error(err)
		os.Exit(10)
	}
	op.Infof("XXX root exists: %s", root)

	err = os.Chdir(root)
	if err != nil {
		op.Errorf("XXX error while chdir outside chroot: %s", err.Error())
		os.Exit(4)
	}

	op.Infof("XXX chdir'd")
	err = syscall.Chroot(root)
	if err != nil {
		op.Errorf("XXX error while chrootin': %s", err.Error())
		os.Exit(3)
	}
	op.Infof("XXX chrooted")

	err = os.Chdir("/")
	if err != nil {
		op.Errorf("XXX error while chdir inside chroot: %s", err.Error())
		os.Exit(2)
	}
	op.Infof("XXX chdir'd")

	if err = archive.InvokeUnpack(op, os.Stdin, filterSpec, "/"); err != nil {
		op.Error(err)
		os.Exit(8)
	}

	op.Infof("XXX success!")
	os.Exit(0)
}
