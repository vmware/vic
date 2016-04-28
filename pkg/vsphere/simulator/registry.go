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

var Map *Registry

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

func (r *Registry) CreateReference(item mo.Reference) types.ManagedObjectReference {
	r.counter++
	kind := reflect.TypeOf(item).Elem().Name()

	return types.ManagedObjectReference{
		Type:  kind,
		Value: fmt.Sprintf("%s-%d", strings.ToLower(kind), r.counter),
	}
}

func (r *Registry) PutEntity(parent mo.Entity, item mo.Entity) mo.Entity {
	e := item.Entity()

	if e.Self.Type == "" {
		e.Self = r.CreateReference(item)
	}

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

func (r *Registry) Put(item mo.Reference) {
	r.m.Lock()
	defer r.m.Unlock()

	r.objects[item.Reference()] = item
}

func (r *Registry) Remove(item types.ManagedObjectReference) {
	r.m.Lock()
	defer r.m.Unlock()

	delete(r.objects, item)
}
