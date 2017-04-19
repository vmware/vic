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
package client

import (
	"errors"
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/d2g/dhcp4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/system"
)

var dummyHWAddr net.HardwareAddr = []byte{0x6, 0x0, 0x0, 0x0, 0x0, 0x0}

func TestMain(m *testing.M) {
	Sys = system.System{
		UUID: uuid.New().String(),
	}

	os.Exit(m.Run())
}

func TestSetOptions(t *testing.T) {
	id, err := NewID(0, dummyHWAddr)
	assert.NoError(t, err)
	c := &client{
		id: id,
	}
	p := dhcp4.NewPacket(dhcp4.BootRequest)

	var tests = []struct {
		prl []byte
	}{
		{
			prl: []byte{
				byte(dhcp4.OptionSubnetMask),
				byte(dhcp4.OptionRouter),
			},
		},
		{
			prl: []byte{
				byte(dhcp4.OptionSubnetMask),
				byte(dhcp4.OptionRouter),
				byte(dhcp4.OptionNameServer),
			},
		},
	}

	for _, te := range tests {
		c.SetParameterRequestList(te.prl...)
		p, err := c.setOptions(p)
		assert.NoError(t, err)
		assert.NotEmpty(t, p)

		opts := p.ParseOptions()
		prl := opts[dhcp4.OptionParameterRequestList]
		assert.NotNil(t, prl)
		assert.EqualValues(t, te.prl, prl)

		// packet should have client id set
		cid := opts[dhcp4.OptionClientIdentifier]
		assert.NotNil(t, cid)
		b, _ := id.MarshalBinary()
		assert.EqualValues(t, cid, b)
	}
}

func TestWithRetry(t *testing.T) {
	errors := []error{
		syscall.Errno(syscall.EAGAIN),
		syscall.Errno(syscall.EINTR),
		errors.New("fail"),
	}

	i := 0

	err := withRetry("test fail", func() error {
		e := errors[i]
		i++
		return e
	})

	if err != errors[len(errors)-1] {
		t.Errorf("err=%s", err)
	}

	err = withRetry("test ok", func() error {
		return nil
	})

	if err != nil {
		t.Errorf("err=%s", err)
	}
}
