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

package vsphere

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

const (
	UUIDPath = "/sys/class/dmi/id/product_serial"
)

// GetSelf gets VirtualMachine reference for the VM this process is running on
func GetSelf(s *session.Session) (*object.VirtualMachine, error) {
	u, err := VMUUID()
	if err != nil {
		return nil, err
	}

	search := object.NewSearchIndex(s.Client.Client)
	ref, err := search.FindByUuid(context.Background(), s.Datacenter, u, true, nil)
	if err != nil {
		return nil, err
	}

	if ref == nil {
		return nil, fmt.Errorf("can't find the hosting vm")
	}

	vm := object.NewVirtualMachine(s.Client.Client, ref.Reference())
	return vm, nil
}

// VMUUID gets the BIOS UUID via the sys interface.  This UUID is known by vphsere
func VMUUID() (string, error) {
	f, err := os.Open(UUIDPath)
	if err != nil {
		return "", err
	}
	r := bufio.NewReader(f)
	uuidstr, err := r.ReadString('\n')
	for err != nil {
		return "", err
	}

	// check the uuid starts with "VMware-"
	if !strings.HasPrefix(uuidstr, "VMware-") {
		return "", fmt.Errorf("cannot find this VM's UUID")
	}

	uuidstr = strings.Replace(uuidstr, "VMware-", "", 1)
	uuidstr = strings.Replace(uuidstr, " ", "", -1)
	uuidstr = strings.TrimSuffix(uuidstr, "\n")

	// need to add dashes, e.g. "564d395e-d807-e18a-cb25-b79f65eb2b9f"
	uuidstr = fmt.Sprintf("%s-%s-%s-%s", uuidstr[0:8], uuidstr[8:12], uuidstr[12:21], uuidstr[21:])

	return uuidstr, nil
}
