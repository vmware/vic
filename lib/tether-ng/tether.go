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
	"reflect"
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

//
type Reporter interface {
	Report(ctx context.Context, err chan<- error)
}

//
type Collector interface {
	Collect(ctx context.Context)
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
	Collector

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

func (t *Tether) Collect(ctx context.Context) {
	var chans = []chan error{}

	t.m.RLock()
	for uuid := range t.plugins {
		plugin := t.plugins[uuid]
		if reporter, ok := plugin.(Reporter); ok {
			ch := make(chan error)
			chans = append(chans, ch)
			go reporter.Report(ctx, ch)
		}
	}
	t.m.RUnlock()

	//FIXME(caglar10ur): Should unregister also remove it from here?
	cases := make([]reflect.SelectCase, len(chans))
	for i := range chans {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(chans[i])}
	}

	remaining := len(cases)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remaining--
			continue
		}

		fmt.Printf("Read from channel %#v and received %s\n", chans[chosen], value)
	}
}
