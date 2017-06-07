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

package vsan

import (
	"context"
	"encoding/json"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	//	"github.com/vmware/govmomi/vim25/mo"
	//	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	//	"github.com/vmware/vic/pkg/trace"
)

const (
	DomObject  = "DOM_OBJECT"
	LsomObject = "LSOM_OBJECT"
	Policy     = "POLICY"
	Disk       = "DISK"
)

type VSANExtAttrs struct {
	Path string `json:"Object path"`
}

// TODO: HostVsanInternalSystem should be moved to govmomi
type HostVsanInternalSystem struct {
	object.Common
	cl *vim25.Client
}

func NewHostVsanInternalSystem(c *vim25.Client, ref types.ManagedObjectReference) *HostVsanInternalSystem {
	m := HostVsanInternalSystem{
		Common: object.NewCommon(c, ref),
		cl:     c,
	}

	return &m
}

// QueryVsanObjectUuidsByFilter returns vsan DOM object uuids by filter
func (m HostVsanInternalSystem) QueryVsanObjectUuidsByFilter(ctx context.Context, uuids []string, limit int32, version int32) ([]string, error) {
	req := types.QueryVsanObjectUuidsByFilter{
		This:    m.Reference(),
		Uuids:   uuids,
		Limit:   limit,
		Version: version,
	}

	res, err := methods.QueryVsanObjectUuidsByFilter(ctx, m.cl, &req)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}
	return res.Returnval, nil
}

// QueryVsanObjects returns vsan detail object information
func (m HostVsanInternalSystem) QueryVsanObjects(ctx context.Context, uuids []string) (string, error) {
	req := types.QueryVsanObjects{
		This:  m.Reference(),
		Uuids: uuids,
	}

	res, err := methods.QueryVsanObjects(ctx, m.cl, &req)
	if err != nil {
		return "", err
	}

	if res == nil {
		return "", nil
	}
	return res.Returnval, nil
}

// GetVsanObjExtAttrs returns vsan extended vsan attributes, including object path. This API can be slow based on API doc
func (m HostVsanInternalSystem) GetVsanObjExtAttrs(ctx context.Context, uuids []string) (map[string]VSANExtAttrs, error) {
	var extAttrs map[string]VSANExtAttrs

	req := types.GetVsanObjExtAttrs{
		This:  m.Reference(),
		Uuids: uuids,
	}

	res, err := methods.GetVsanObjExtAttrs(ctx, m.cl, &req)
	if err != nil {
		return extAttrs, err
	}

	if res == nil {
		return extAttrs, nil
	}
	log.Debugf("GetVsanObjExtAttrs: returned attributes: %s", res.Returnval)
	extAttrs = make(map[string]VSANExtAttrs)
	err = json.Unmarshal([]byte(res.Returnval), &extAttrs)
	return extAttrs, err
}

// QueryCmmds query vsan CMMDS directly
func (m HostVsanInternalSystem) QueryCmmds(ctx context.Context, owner string, cmmsType string, uuid string) (string, error) {
	query := types.HostVsanInternalSystemCmmdsQuery{
		Owner: owner,
		Type:  cmmsType,
		Uuid:  uuid,
	}
	req := types.QueryCmmds{
		This:    m.Reference(),
		Queries: []types.HostVsanInternalSystemCmmdsQuery{query},
	}

	res, err := methods.QueryCmmds(ctx, m.cl, &req)
	if err != nil {
		return "", err
	}

	if res == nil {
		return "", nil
	}
	return res.Returnval, nil
}

// DeleteVsanObjects delete vsan DOM object directly
func (m HostVsanInternalSystem) DeleteVsanObjects(ctx context.Context, uuids []string, force *bool) ([]types.HostVsanInternalSystemDeleteVsanObjectsResult, error) {
	req := types.DeleteVsanObjects{
		This:  m.Reference(),
		Uuids: uuids,
		Force: force,
	}

	res, err := methods.DeleteVsanObjects(ctx, m.cl, &req)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}
	return res.Returnval, nil
}
