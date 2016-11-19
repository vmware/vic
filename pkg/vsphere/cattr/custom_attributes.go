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

package cattr

import (
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/session"

	"context"
)

const (
	InventoryCategory = "ContainersInventory"
	ManagedEntityType = "ManagedEntity"
	ImageCategory     = "ImageUpload"
	PowerCategory     = "PoweringOn"
	ProbingCategory   = "Probing"

	VCH        = "Virtual Container Host"
	InProgress = "In progress"
	Completed  = "Completed"
	Failed     = "Failed"
)

type Manager struct {
	session *session.Session
	ctx     context.Context
}

func NewManager(session *session.Session, ctx context.Context) *Manager {
	return &Manager{
		session,
		ctx,
	}
}

func (m *Manager) isDuplicateName(err error) bool {
	if soap.IsSoapFault(err) {
		fault := soap.ToSoapFault(err)
		if _, ok := fault.VimFault().(types.DuplicateName); ok {
			return true
		}
	}
	return false
}

func (m *Manager) CustomField(managedObject types.ManagedObjectReference, managedType string, key string, value string) error {
	mgr, err := object.GetCustomFieldsManager(m.session.Client.Client)
	if err != nil {
		return err
	}

	if managedType == "" {
		_, err = mgr.Add(m.ctx, key, managedObject.Type, nil, nil)
	} else {
		_, err = mgr.Add(m.ctx, key, managedType, nil, nil)
	}

	if err != nil && !m.isDuplicateName(err) {
		return err
	}

	ID, err := mgr.FindKey(m.ctx, key)
	if err != nil {
		return err
	}

	err = mgr.Set(m.ctx, managedObject, ID, value)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) MarkAsVCH(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, ManagedEntityType, InventoryCategory, VCH)
}

func (m *Manager) PoweringOnInProgress(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", PowerCategory, InProgress)
}

func (m *Manager) PoweringOnCompleted(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", PowerCategory, Completed)
}

func (m *Manager) PoweringOnFailed(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", PowerCategory, Failed)
}

func (m *Manager) ImageUploadInProgress(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", ImageCategory, InProgress)
}

func (m *Manager) ImageUploadCompleted(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", ImageCategory, Completed)
}

func (m *Manager) ImageUploadFailed(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", ImageCategory, Failed)
}

func (m *Manager) ProbingInProgress(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", ProbingCategory, InProgress)
}

func (m *Manager) ProbingCompleted(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", ProbingCategory, Completed)
}

func (m *Manager) ProbingFailed(managedObject types.ManagedObjectReference) error {
	return m.CustomField(managedObject, "", ProbingCategory, Failed)
}
