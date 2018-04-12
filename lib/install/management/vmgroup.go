// Copyright 2018 VMware, Inc. All Rights Reserved.
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

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/mo"
)

func (d *Dispatcher) createVMGroup(conf *config.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin("", d.op))

	if !conf.UseVMGroup {
		return nil
	}

	d.op.Debugf("Creating DRS VM Group %q on %q", conf.VMGroupName, d.appliance.Cluster)

	spec := &types.ClusterConfigSpecEx{
		GroupSpec: []types.ClusterGroupSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationAdd,
				},
				Info: &types.ClusterVmGroup{
					ClusterGroupInfo: types.ClusterGroupInfo{
						Name: conf.VMGroupName,
					},
					Vm: []types.ManagedObjectReference{d.appliance.Reference()},
				},
			},
		},
	}

	_, err := tasks.WaitForResultAndRetryIf(d.op, func(op context.Context) (tasks.Task, error) {
		return d.appliance.Cluster.Reconfigure(op, spec, true)
	}, tasks.IsTransientError)

	return err
}

func (d *Dispatcher) destroyVMGroup(conf *config.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin("", d.op))

	if !conf.UseVMGroup {
		return nil
	}

	d.op.Debugf("Checking for existence of DRS VM Group %s on %s", conf.VMGroupName, d.session.Cluster)

	var clusterConfig mo.ClusterComputeResource
	err := d.session.Cluster.Properties(d.op, d.session.Cluster.Reference(), []string{"configurationEx"}, &clusterConfig)
	if err != nil {
		d.op.Warnf("Unable to obtain cluster config: %s", err)
		return nil
	}

	groupExists := false
	clusterConfigEx := clusterConfig.ConfigurationEx.(*types.ClusterConfigInfoEx)
	for _, g := range clusterConfigEx.Group {
		if g.GetClusterGroupInfo().Name == conf.VMGroupName {
			groupExists = true
			break
		}
	}

	if !groupExists {
		d.op.Debugf("Expected VM Group cannot be found; skipping removal.")
		return nil
	}

	d.op.Infof("Removing VM Group %q", conf.VMGroupName)


	spec := &types.ClusterConfigSpecEx{
		GroupSpec: []types.ClusterGroupSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationRemove,
					RemoveKey: conf.VMGroupName,
				},
			},
		},
	}

	_, err = tasks.WaitForResultAndRetryIf(d.op, func(op context.Context) (tasks.Task, error) {
		return d.appliance.Cluster.Reconfigure(op, spec, true)
	}, tasks.IsTransientError)

	return err
}
