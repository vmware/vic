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

package simulator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

var Map = NewRegistry()

type Registry struct {
	m       sync.Mutex
	objects map[types.ManagedObjectReference]mo.Reference
	counter int
}

func NewRegistry() *Registry {
	r := &Registry{
		objects: make(map[types.ManagedObjectReference]mo.Reference),
	}

	return r
}

func TypeName(item mo.Reference) string {
	return reflect.TypeOf(item).Elem().Name()
}

// newReference returns a new MOR, where Type defaults to type of the given item
// and Value defaults to a unique id for the given type.
func (r *Registry) newReference(item mo.Reference) types.ManagedObjectReference {
	ref := item.Reference()

	if ref.Type == "" {
		ref.Type = TypeName(item)
	}

	if ref.Value == "" {
		r.counter++
		ref.Value = fmt.Sprintf("%s-%d", strings.ToLower(ref.Type), r.counter)
	}

	return ref
}

// NewEntity sets Entity().Self with a new, unique Value.
// Useful for creating object instances from templates.
func (r *Registry) NewEntity(item mo.Entity) mo.Entity {
	e := item.Entity()
	e.Self.Value = ""
	e.Self = r.newReference(item)
	return item
}

func (r *Registry) PutEntity(parent mo.Entity, item mo.Entity) mo.Entity {
	e := item.Entity()

	if parent != nil {
		e.Parent = &parent.Entity().Self
	}

	r.Put(item)

	return item
}

func (r *Registry) Get(ref types.ManagedObjectReference) mo.Reference {
	r.m.Lock()
	defer r.m.Unlock()

	return r.objects[ref]
}

func (r *Registry) Put(item mo.Reference) mo.Reference {
	r.m.Lock()
	defer r.m.Unlock()

	ref := item.Reference()
	if ref.Type == "" || ref.Value == "" {
		ref = r.newReference(item)
		// mo.Reference() returns a value, not a pointer so use reflect to set the Self field
		reflect.ValueOf(item).Elem().FieldByName("Self").Set(reflect.ValueOf(ref))
	}

	r.objects[ref] = item

	return item
}

func (r *Registry) Remove(item types.ManagedObjectReference) {
	r.m.Lock()
	defer r.m.Unlock()

	delete(r.objects, item)
}

// getEntityParent traverses up the inventory and returns the first object of type kind.
// If no object of type kind is found, the method will panic when it reaches the
// inventory root Folder where the Parent field is nil.
func (r *Registry) getEntityParent(item mo.Entity, kind string) mo.Entity {
	for {
		parent := item.Entity().Parent

		item = Map.Get(*parent).(mo.Entity)

		if item.Reference().Type == kind {
			return item
		}
	}
}

// getEntityDatacenter returns the Datacenter containing the given item
func (r *Registry) getEntityDatacenter(item mo.Entity) *mo.Datacenter {
	return r.getEntityParent(item, "Datacenter").(*mo.Datacenter)
}

// FindByName returns the first mo.Entity of the given refs whose Name field is equal to the given name.
// If there is no match, nil is returned.
// This method is useful for cases where objects are required to have a unique name, such as Datastore with
// a HostStorageSystem or HostSystem within a ClusterComputeResource.
func (r *Registry) FindByName(name string, refs []types.ManagedObjectReference) mo.Entity {
	for _, ref := range refs {
		e := Map.Get(ref).(mo.Entity)

		if name == e.Entity().Name {
			return e
		}
	}

	return nil
}
