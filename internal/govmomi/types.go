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

package llpm

import "github.com/vmware/govmomi/vim25/types"

type RetrieveInternalConfigManagerRequest struct {
	This types.ManagedObjectReference `xml:"_this"`
}

type RetrieveInternalConfigManagerResponse struct {
	Returnval *InternalConfigManager `xml:"urn:vim25 returnval"`
}

// Request structure for LLPM.deleteVm
type LowLevelCreateVm_Task struct {
	This       types.ManagedObjectReference   `xml:"_this"`
	ConfigSpec types.VirtualMachineConfigSpec `xml:"configSpec"`
}

type LowLevelCreateVm_TaskResponse struct {
	Returnval types.ManagedObjectReference `xml:"returnval"`
}

// Request structure for LLPM.deleteVm
type DeleteVm_Task struct {
	This       types.ManagedObjectReference   `xml:"_this"`
	ConfigInfo types.VirtualMachineConfigInfo `xml:"configInfo"`
}

type DeleteVm_TaskResponse struct {
	Returnval types.ManagedObjectReference `xml:"returnval"`
}

// Request structure for LLPM.deleteVm
type DeleteVmExceptDisks_Task struct {
	This       types.ManagedObjectReference `xml:"_this"`
	VmPathName string                       `xml:"vmPathName"`
}

type DeleteVmExceptDisks_TaskResponse struct {
	Returnval types.ManagedObjectReference `xml:"returnval"`
}

// Request structure for LLPM.reconfigVM
type LLPMReconfigVM_Task struct {
	This       types.ManagedObjectReference   `xml:"_this"`
	ConfigSpec types.VirtualMachineConfigSpec `xml:"configSpec"`
}

// Common response type for reconfigVM operations
type ReconfigVM_TaskResponse struct {
	Returnval types.ManagedObjectReference `xml:"returnval"`
}
