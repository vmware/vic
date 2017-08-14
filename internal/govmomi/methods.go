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

import (
	"context"

	"github.com/vmware/govmomi/vim25/soap"
)

type RetrieveInternalConfigManagerBody struct {
	Req    *RetrieveInternalConfigManagerRequest  `xml:"urn:vim25 RetrieveInternalConfigManager"`
	Res    *RetrieveInternalConfigManagerResponse `xml:"urn:vim25 RetrieveInternalConfigManagerResponse"`
	Fault_ *soap.Fault
}

func (b *RetrieveInternalConfigManagerBody) Fault() *soap.Fault { return b.Fault_ }

func RetrieveInternalConfigManager(ctx context.Context, r soap.RoundTripper, req *RetrieveInternalConfigManagerRequest) (*RetrieveInternalConfigManagerResponse, error) {
	var reqBody, resBody RetrieveInternalConfigManagerBody

	reqBody.Req = req

	if err := r.RoundTrip(ctx, &reqBody, &resBody); err != nil {
		return nil, err
	}

	return resBody.Res, nil
}

type LowLevelCreateVmBody struct {
	Req    *LowLevelCreateVm_Task         `xml:"urn:internalvim25 LowLevelCreateVm_Task"`
	Res    *LowLevelCreateVm_TaskResponse `xml:"urn:vim25 LowLevelCreateVm_TaskResponse"`
	Fault_ *soap.Fault
}

func (b *LowLevelCreateVmBody) Fault() *soap.Fault { return b.Fault_ }

func LowLevelCreateVm(ctx context.Context, r soap.RoundTripper, req *LowLevelCreateVm_Task) (*LowLevelCreateVm_TaskResponse, error) {
	var reqBody, resBody LowLevelCreateVmBody

	reqBody.Req = req

	if err := r.RoundTrip(ctx, &reqBody, &resBody); err != nil {
		return nil, err
	}

	return resBody.Res, nil
}

type LowLevelDeleteVmBody struct {
	Req    *DeleteVm_Task         `xml:"urn:internalvim25 DeleteVm_Task"`
	Res    *DeleteVm_TaskResponse `xml:"urn:vim25 DeleteVm_TaskResponse"`
	Fault_ *soap.Fault
}

func (b *LowLevelDeleteVmBody) Fault() *soap.Fault { return b.Fault_ }

func LowLevelDeleteVm(ctx context.Context, r soap.RoundTripper, req *DeleteVm_Task) (*DeleteVm_TaskResponse, error) {
	var reqBody, resBody LowLevelDeleteVmBody

	reqBody.Req = req

	if err := r.RoundTrip(ctx, &reqBody, &resBody); err != nil {
		return nil, err
	}

	return resBody.Res, nil
}

type LowLevelDeleteVmExceptDisksBody struct {
	Req    *DeleteVmExceptDisks_Task         `xml:"urn:internalvim25 DeleteVmExceptDisks_Task"`
	Res    *DeleteVmExceptDisks_TaskResponse `xml:"urn:vim25 DeleteVmExceptDisks_TaskResponse"`
	Fault_ *soap.Fault
}

func (b *LowLevelDeleteVmExceptDisksBody) Fault() *soap.Fault { return b.Fault_ }

func LowLevelDeleteVmExceptDisks(ctx context.Context, r soap.RoundTripper, req *DeleteVmExceptDisks_Task) (*DeleteVmExceptDisks_TaskResponse, error) {
	var reqBody, resBody LowLevelDeleteVmExceptDisksBody

	reqBody.Req = req

	if err := r.RoundTrip(ctx, &reqBody, &resBody); err != nil {
		return nil, err
	}

	return resBody.Res, nil
}

type LowLevelReconfigVMBody struct {
	Req    *LLPMReconfigVM_Task     `xml:"urn:internalvim25 LLPMReconfigVM_Task"`
	Res    *ReconfigVM_TaskResponse `xml:"urn:vim25 LLPMReconfigVM_TaskResponse"`
	Fault_ *soap.Fault              `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault,omitempty"`
}

func (b *LowLevelReconfigVMBody) Fault() *soap.Fault { return b.Fault_ }

func LowLevelReconfigVM(ctx context.Context, r soap.RoundTripper, req *LLPMReconfigVM_Task) (*ReconfigVM_TaskResponse, error) {
	var reqBody, resBody LowLevelReconfigVMBody

	reqBody.Req = req

	if err := r.RoundTrip(ctx, &reqBody, &resBody); err != nil {
		return nil, err
	}

	return resBody.Res, nil
}
