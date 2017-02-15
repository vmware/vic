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

// Signaler provides the process signal related methods
type Signaler interface {
	Running(ctx context.Context, sessionID string) bool
	Kill(ctx context.Context, sessionID string) error
}

// Releaser is called by Interactor to release Waiter
type Releaser interface {
	Release(ctx context.Context, out chan<- chan struct{})
}

//Waiter is called by process plugin to wait Releaser to release
type Waiter interface {
	Wait(ctx context.Context, in <-chan chan struct{})
}

// Interactor is called by process plugin to mutate its writers/readers
type Interactor interface {
	PseudoTerminal(ctx context.Context, in <-chan *types.Session) <-chan struct{}
	NonInteract(ctx context.Context, in <-chan *types.Session) <-chan struct{}

	Close(ctx context.Context, in <-chan *types.Session) <-chan struct{}
}

// Reaper implements the reaper to reap processes
type Reaper interface {
	Reap(ctx context.Context) error
}

// Reporter implements the error reporting mechanism
type Reporter interface {
	Report(ctx context.Context, err chan<- error)
}

// Collector implements the error collecting mechanism
type Collector interface {
	Collect(ctx context.Context)
}

// Plugin implements the plugins
type Plugin interface {
	Configure(ctx context.Context, config *types.ExecutorConfig) error

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	UUID(ctx context.Context) uuid.UUID
}

// PluginRegistrar is the registry of Plugins
type PluginRegistrar interface {
	Register(ctx context.Context, plugin Plugin) error
	Unregister(ctx context.Context, plugin Plugin) error
	Plugins(ctx context.Context) []Plugin
}

// Tether implements PluginRegistrar and Collector
type Tether struct {
	PluginRegistrar
	Collector

	ctx context.Context

	m       sync.RWMutex
	plugins map[uuid.UUID]Plugin
}

// NewTether returns a new tether instance
func NewTether(ctx context.Context) PluginRegistrar {
	return &Tether{
		ctx:     ctx,
		plugins: make(map[uuid.UUID]Plugin),
	}
}

// Register registers the plugin
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

// Unregister unregisters the plugin
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

// Plugins returns the plugins
func (t *Tether) Plugins(ctx context.Context) []Plugin {
	t.m.RLock()
	defer t.m.RUnlock()

	var list []Plugin
	for uuid := range t.plugins {
		list = append(list, t.plugins[uuid])
	}
	return list
}

// Collect creates an error channel for each Reporter and sends it to them. Then creates a dynamic select statement using those
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
