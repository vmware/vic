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
	"bufio"
	"context"
	"os"
	"syscall"

	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/pkg/trace"
)

func main() {
	ctx := context.Background()
	op := trace.NewOperation(ctx, "Unpack") // TODO op ID?
	op.Infof("XXX New unpack operation created")

	if len(os.Args) < 3 {
		op.Errorf("XXX Not enough arguments passed")
		os.Exit(5)
	}

	root := os.Args[2]
	op.Infof("XXX root inside binary %s", root)

	err := os.Chdir(root)
	if err != nil {
		op.Errorf("XXX error while chdir outside chroot: %s", err.Error())
		os.Exit(4)
	}

	err = syscall.Chroot(root)
	if err != nil {
		op.Errorf("XXX error while chrootin': %s", err.Error())
		os.Exit(3)
	}

	err = os.Chdir("/")
	if err != nil {
		op.Errorf("XXX error while chdir inside chroot: %s", err.Error())
		os.Exit(2)
	}

	r := bufio.NewReader(os.Stdin)

	filterSpecBytes := []byte{}
FilterSpec:
	f, isPrefix, err := r.ReadLine()
	if err != nil {
		op.Error(err)
		os.Exit(6)
	}

	filterSpecBytes = append(filterSpecBytes, f...)
	if isPrefix {
		op.Infof("XXX Continuing to read encoded FilterSpec..")
		// this is a goto because it was easier to reason about than a loop. fight me.
		goto FilterSpec
	}

	// everything after the FilterSpec is the tarstream
	filterSpecString := string(filterSpecBytes)
	filterSpec, err := archive.DecodeFilterSpec(op, &filterSpecString)
	if err != nil {
		op.Error(err)
		os.Exit(7)
	}

	if err = archive.InvokeUnpack(op, os.Stdin, filterSpec); err != nil {
		op.Error(err)
		os.Exit(8)
	}
	os.Exit(0)
}
