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

package management

import (
	"github.com/vmware/govmomi/vim25/types"
)

var rolePrefix = "vic-vch-"

var RoleVCenter = types.AuthorizationRole{
	Name: "vcenter",
	Privilege: []string{
		"Datastore.Config",
	},
}

var RoleDataCenter = types.AuthorizationRole{
	Name: "datacenter",
	Privilege: []string{
		"Datastore.Config",
		"Datastore.FileManagement",
		"VirtualMachine.Config.AddNewDisk",
		"VirtualMachine.Config.AdvancedConfig",
		"VirtualMachine.Config.RemoveDisk",
		"VirtualMachine.Inventory.Create",
		"VirtualMachine.Inventory.Delete",
	},
}

var RoleCluster = types.AuthorizationRole{
	Name: "cluster",
	Privilege: []string{
		"Datastore.AllocateSpace",
		"Datastore.Browse",
		"Datastore.Config",
		"Datastore.DeleteFile",
		"Datastore.FileManagement",
		"Host.Config.SystemManagement",
	},
}

var RoleDataStore = types.AuthorizationRole{
	Name: "datastore",
	Privilege: []string{
		"Datastore.AllocateSpace",
		"Datastore.Browse",
		"Datastore.Config",
		"Datastore.DeleteFile",
		"Datastore.FileManagement",
		"Host.Config.SystemManagement",
	},
}

var RoleNetwork = types.AuthorizationRole{
	Name: "network",
	Privilege: []string{
		"Network.Assign",
	},
}

var RoleEndpoint = types.AuthorizationRole{
	Name: "endpoint",
	Privilege: []string{
		"DVPortgroup.Modify",
		"DVPortgroup.PolicyOp",
		"DVPortgroup.ScopeOp",
		"Resource.AssignVMToPool",
		"VirtualMachine.Config.AddExistingDisk",
		"VirtualMachine.Config.AddNewDisk",
		"VirtualMachine.Config.AddRemoveDevice",
		"VirtualMachine.Config.AdvancedConfig",
		"VirtualMachine.Config.EditDevice",
		"VirtualMachine.Config.RemoveDisk",
		"VirtualMachine.Config.Rename",
		"VirtualMachine.GuestOperations.Execute",
		"VirtualMachine.Interact.DeviceConnection",
		"VirtualMachine.Interact.PowerOff",
		"VirtualMachine.Interact.PowerOn",
		"VirtualMachine.Inventory.Create",
		"VirtualMachine.Inventory.Delete",
		"VirtualMachine.Inventory.Register",
		"VirtualMachine.Inventory.Unregister",
	},
}

// Configuration for the ops-user
var OpsUserRBACConf = RBACConfig{
	Resources: []RBACResource{
		{
			Type:      VCenter,
			Propagate: false,
			Role:      RoleVCenter,
		},
		{
			Type:      Datacenter,
			Propagate: true,
			Role:      RoleDataCenter,
		},
		{
			Type:      Cluster,
			Propagate: true,
			Role:      RoleDataStore,
		},
		{
			Type:      DatastoreFolder,
			Propagate: true,
			Role:      RoleDataStore,
		},
		{
			Type:      Datastore,
			Propagate: false,
			Role:      RoleDataStore,
		},
		{
			Type:      VSANDatastore,
			Propagate: false,
			Role:      RoleDataStore,
		},
		{
			Type:      Network,
			Propagate: true,
			Role:      RoleNetwork,
		},
		{
			Type:      Endpoint,
			Propagate: true,
			Role:      RoleEndpoint,
		},
	},
}
