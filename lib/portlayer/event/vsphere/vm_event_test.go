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
	"testing"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/event"

	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	vm := newVMMO()
	k := 1
	msg := "jojo the idiot circus boy"
	vmwEve := &types.VmPoweredOnEvent{VmEvent: types.VmEvent{Event: types.Event{FullFormattedMessage: msg, Key: int32(k), Vm: &types.VmEventArgument{Vm: *vm}}}}
	vme := NewVMEvent(vmwEve)
	assert.NotNil(t, vme)
	assert.Equal(t, event.ContainerPoweredOn, vme.String())
	assert.Equal(t, vm.String(), vme.Reference())
	assert.Equal(t, k, vme.EventID())
	assert.Equal(t, msg, vme.Message())
}
