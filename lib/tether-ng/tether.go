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

package tetherng

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/vmware/vic/lib/tether-ng/types"
)

type Signaler interface {
	Running(ctx context.Context, sessionID string) bool
	Kill(ctx context.Context, sessionID string) error
}

// Interaction calls this to release Waiter
type Releaser interface {
	Release(ctx context.Context, out chan<- chan struct{})
}

// Process calls this to wait Releaser to release
type Waiter interface {
	Wait(ctx context.Context, in <-chan chan struct{})
}

// Process calls this to mutate writers/readers
type Interactor interface {
	PseudoTerminal(ctx context.Context, in <-chan *types.Session) <-chan struct{}
	NonInteract(ctx context.Context, in <-chan *types.Session) <-chan struct{}

	Close(ctx context.Context, in <-chan *types.Session) <-chan struct{}
}

//
type Reaper interface {
	Reap(ctx context.Context) error
}

type Plugin interface {
	Configure(ctx context.Context, config *types.ExecutorConfig) error

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	UUID(ctx context.Context) uuid.UUID
}

type PluginRegistrar interface {
	Register(ctx context.Context, plugin Plugin) error
	Unregister(ctx context.Context, plugin Plugin) error
	Plugins(ctx context.Context) []Plugin
}

type Tether struct {
	PluginRegistrar

	ctx context.Context

	m       sync.RWMutex
	plugins map[uuid.UUID]Plugin
}

func NewTether(ctx context.Context) PluginRegistrar {
	return &Tether{
		ctx:     ctx,
		plugins: make(map[uuid.UUID]Plugin),
	}
}

func (t *Tether) Register(ctx context.Context, plugin Plugin) error {
	t.m.Lock()
	defer t.m.Unlock()

	uuid := plugin.UUID(ctx)

	if plugin == nil {
		return fmt.Errorf("tether: Driver is nil")
	}
	if _, dup := t.plugins[uuid]; dup {
		return fmt.Errorf("tether: Plugin %s already registered", uuid)
	}
	t.plugins[uuid] = plugin

	return nil
}

func (t *Tether) Unregister(ctx context.Context, plugin Plugin) error {
	t.m.Lock()
	defer t.m.Unlock()

	uuid := plugin.UUID(ctx)
	_, ok := t.plugins[uuid]
	if !ok {
		return fmt.Errorf("tether: No such plugin %s", uuid)
	}
	delete(t.plugins, uuid)

	return nil
}

func (t *Tether) Plugins(ctx context.Context) []Plugin {
	t.m.RLock()
	defer t.m.RUnlock()

	var list []Plugin
	for uuid := range t.plugins {
		list = append(list, t.plugins[uuid])
	}
	return list
}
