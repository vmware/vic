// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package management

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/simulator"
)

func TestGuestInfoSecret(t *testing.T) {
	ctx := context.Background()

	for i, m := range []*simulator.Model{simulator.ESX(), simulator.VPX()} {

		defer m.Remove()
		err := m.Create()
		if err != nil {
			t.Fatal(err)
		}

		server := m.Service.NewServer()
		defer server.Close()

		var s *session.Session
		if i == 0 {
			s, err = getESXSession(ctx, server.URL.String())
		} else {
			s, err = getVPXSession(ctx, server.URL.String())
		}
		if err != nil {
			t.Fatal(err)
		}

		if s, err = s.Populate(ctx); err != nil {
			t.Fatal(err)
		}

		name := "my-vm"
		vmx := fmt.Sprintf("%s/%s.vmx", name, name)
		ds := s.Datastore
		secretKey, err := extraconfig.NewSecretKey()
		if err != nil {
			t.Fatal(err)
		}

		spec := types.VirtualMachineConfigSpec{
			Name:    name,
			GuestId: string(types.VirtualMachineGuestOsIdentifierOtherGuest),
			Files: &types.VirtualMachineFileInfo{
				VmPathName: fmt.Sprintf("[%s] %s", ds.Name(), vmx),
			},
			ExtraConfig: []types.BaseOptionValue{
				&types.OptionValue{
					Key:   extraconfig.GuestInfoSecretKey,
					Value: secretKey.String(),
				},
			},
		}

		folder := s.Folders(ctx).VmFolder
		task, err := folder.CreateVM(ctx, spec, s.Pool, nil)
		if err != nil {
			t.Fatal(err)
		}
		err = task.Wait(ctx)
		if err != nil {
			t.Fatal(err)
		}

		d := &Dispatcher{
			session:    s,
			ctx:        ctx,
			vmPathName: name,
		}

		// Attempt to extract the secret without setting the session's datastore
		d.session.Datastore = nil
		secret, err := d.GuestInfoSecret(name)
		assert.Nil(t, secret)
		assert.Equal(t, err, errNilDatastore)

		d.session.Datastore = ds

		// Attempt to extract the secret with an empty .vmx file
		// TODO: simulator should write ExtraConfig to the .vmx file
		secret, err = d.GuestInfoSecret(name)
		assert.Nil(t, secret)
		assert.Equal(t, err, errSecretKeyNotFound)

		// Write malformed key-value pairs
		dir := simulator.Map.Get(ds.Reference()).(*simulator.Datastore).Info.(*types.LocalDatastoreInfo).Path
		text := fmt.Sprintf("foo.bar = \"baz\"\n%s \"%s\"\n", extraconfig.GuestInfoSecretKey, secretKey.String())
		if err = ioutil.WriteFile(path.Join(dir, vmx), []byte(text), 0); err != nil {
			t.Fatal(err)
		}

		// Attempt to extract the secret from an incorrectly populated .vmx file
		secret, err = d.GuestInfoSecret(name)
		assert.Nil(t, secret)
		assert.Error(t, err)

		// Write an invalid key that only prefix-matches the secret key
		text = fmt.Sprintf("%s = \"%s\"\n", extraconfig.GuestInfoSecretKey+"1", secretKey.String())
		if err = ioutil.WriteFile(path.Join(dir, vmx), []byte(text), 0); err != nil {
			t.Fatal(err)
		}

		// Attempt to extract the secret from an incorrectly populated .vmx file
		secret, err = d.GuestInfoSecret(name)
		assert.Nil(t, secret)
		assert.Equal(t, err, errSecretKeyNotFound)

		// Write valid key-value pairs
		text = fmt.Sprintf("foo.bar = \"baz\"\n%s = \"%s\"\n", extraconfig.GuestInfoSecretKey, secretKey.String())
		if err = ioutil.WriteFile(path.Join(dir, vmx), []byte(text), 0); err != nil {
			t.Fatal(err)
		}

		// Extract the secret from a correctly populated .vmx file
		secret, err = d.GuestInfoSecret(name)
		assert.NoError(t, err)
		assert.Equal(t, secret.String(), secretKey.String())
	}
}
