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

package management

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/simulator"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

func TestFinder(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	trace.Logger.Level = log.DebugLevel
	ctx := context.Background()

	for i, model := range []*simulator.Model{simulator.ESX(), simulator.VPX()} {
		t.Logf("%d", i)
		defer model.Remove()
		if i == 1 {
			model.Datacenter = 2
			model.Cluster = 2
			model.Host = 2
			model.Pool = 0
		}
		err := model.Create()
		if err != nil {
			t.Fatal(err)
		}

		s := model.Service.NewServer()
		defer s.Close()

		s.URL.User = url.UserPassword("user", "pass")
		s.URL.Path = ""
		t.Logf("server URL: %s", s.URL)

		var input *data.Data
		if i == 0 {
			input = getESXData(s.URL)
		} else {
			input = getVPXData(s.URL)
		}
		if err != nil {
			t.Fatal(err)
		}
		validator, err := validate.NewValidator(ctx, input)
		if err != nil {
			t.Errorf("Failed to create validator: %s", err)
		}
		if _, err = validator.ValidateTarget(ctx, input); err != nil {
			t.Logf("Got expected error to validate target: %s", err)
		}
		validator.AllowEmptyDC()
		if _, err = validator.ValidateTarget(ctx, input); err != nil {
			t.Errorf("Failed to valiate target: %s", err)
		}
		prefix := fmt.Sprintf("p%d-", i)
		if err = createTestData(ctx, validator.Session, prefix); err != nil {
			t.Errorf("Failed to create test data: %s", err)
		}
		// FIXME: ServerFaultCode: no such object: SearchIndex:ha-searchindex
		//		testSearchVCHs(t, validator)
	}
}

type testSearchDispatcher struct {
	*Dispatcher
}

func (td *testSearchDispatcher) isVCH(vm *vm.VirtualMachine) (bool, error) {
	return true, nil
}

func testSearchVCHs(t *testing.T, v *validate.Validator) {
	d := &Dispatcher{
		session: v.Session,
		ctx:     v.Context,
		isVC:    v.Session.IsVC(),
	}

	td := &testSearchDispatcher{d}
	vchs, err := td.SearchVCHs("")
	if err != nil {
		t.Errorf("Failed to search vchs: %s", err)
	}
	t.Logf("Got %d VCHs", len(vchs))
}

func createTestData(ctx context.Context, sess *session.Session, prefix string) error {
	dcs, err := sess.Finder.DatacenterList(ctx, "*")
	if err != nil {
		return err
	}
	for _, dc := range dcs {
		sess.Finder.SetDatacenter(dc)
		sess.Datacenter = dc
		resources := &Node{
			Kind: rpNode,
			Name: prefix + "Root",
			Children: []*Node{
				{
					Kind: rpNode,
					Name: prefix + "pool1",
					Children: []*Node{
						{
							Kind: vmNode,
							Name: prefix + "pool1",
						},
						{
							Kind: rpNode,
							Name: prefix + "pool1-2",
							Children: []*Node{
								{
									Kind: rpNode,
									Name: prefix + "pool1-2-1",
									Children: []*Node{
										{
											Kind: vmNode,
											Name: prefix + "vch1-2-1",
										},
									},
								},
							},
						},
					},
				},
				{
					Kind: vmNode,
					Name: prefix + "vch2",
				},
			},
		}
		if err = createResources(ctx, sess, resources); err != nil {
			return err
		}
	}
	return nil
}

type nodeKind string

const (
	vmNode   = nodeKind("VM")
	rpNode   = nodeKind("RP")
	vappNode = nodeKind("VAPP")
)

type Node struct {
	Kind     nodeKind
	Name     string
	Children []*Node
}

func createResources(ctx context.Context, sess *session.Session, node *Node) error {
	rootPools, err := sess.Finder.ResourcePoolList(ctx, "*")
	if err != nil {
		return err
	}
	for _, pool := range rootPools {
		base := path.Base(path.Dir(pool.InventoryPath))
		log.Debugf("root pool base name %q", base)
		if err = createNodes(ctx, sess, pool, node, base); err != nil {
			return err
		}
	}
	return nil
}

func createNodes(ctx context.Context, sess *session.Session, pool *object.ResourcePool, node *Node, base string) error {
	log.Debugf("create node %+v", node)
	if node == nil {
		return nil
	}
	spec := simulator.NewResourceConfigSpec()
	node.Name = fmt.Sprintf("%s-%s", base, node.Name)
	switch node.Kind {
	case rpNode:
		child, err := pool.Create(ctx, node.Name, spec)
		if err != nil {
			return err
		}
		for _, childNode := range node.Children {
			return createNodes(ctx, sess, child, childNode, base)
		}
	case vappNode:
		confSpec := types.VAppConfigSpec{
			VmConfigSpec: types.VmConfigSpec{},
		}
		vapp, err := pool.CreateVApp(ctx, node.Name, spec, confSpec, nil)
		if err != nil {
			return err
		}
		config := types.VirtualMachineConfigSpec{
			Name:    node.Name,
			GuestId: string(types.VirtualMachineGuestOsIdentifierOtherGuest),
			Files: &types.VirtualMachineFileInfo{
				VmPathName: fmt.Sprintf("[LocalDS_0] %s", node.Name),
			},
		}
		if _, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
			return vapp.CreateChildVM(ctx, config, nil)
		}); err != nil {
			return err
		}
	case vmNode:
		config := types.VirtualMachineConfigSpec{
			Name:    node.Name,
			GuestId: string(types.VirtualMachineGuestOsIdentifierOtherGuest),
			Files: &types.VirtualMachineFileInfo{
				VmPathName: fmt.Sprintf("[LocalDS_0] %s", node.Name),
			},
		}
		folder := sess.Folders(ctx).VmFolder
		if _, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
			return folder.CreateVM(ctx, config, pool, nil)
		}); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}
