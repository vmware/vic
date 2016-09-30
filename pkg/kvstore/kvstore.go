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

package kvstore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"

	"github.com/vmware/vic/pkg/trace"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

// This package implements a very basic key/value store.  It is up to the
// caller to provision the namespace.

type KeyValueStore struct {
	b Backend

	kv map[string][]byte

	fileName string

	l sync.RWMutex
}

type Backend interface {
	// Creates path and overwrites whatever is there.
	Upload(ctx context.Context, r io.Reader, pth string) error

	// Downloads from the given path.
	Download(ctx context.Context, pth string) (io.ReadCloser, error)

	// Moves the given path.
	Mv(ctx context.Context, fromPath, toPath string) error
}

// Create a new KeyValueStore instance using the given Backend with the given
// file.  If the file exists on the Backend, it is restored.
func NewKeyValueStore(op trace.Operation, store Backend, fileName string) (*KeyValueStore, error) {
	p := &KeyValueStore{
		b:        store,
		kv:       make(map[string][]byte),
		fileName: fileName,
	}

	if err := p.restore(op); err != nil {
		return nil, err
	}

	op.Infof("KeyValueStore(%s) restored %d keys", fileName, len(p.kv))

	return p, nil
}

func (p *KeyValueStore) restore(op trace.Operation) error {
	p.l.Lock()
	defer p.l.Unlock()

	rc, err := p.b.Download(op, p.fileName)
	if err != nil {
		// We need to check for 404 vs something else here.
		op.Errorf("KeyValueStore(%s) ignoring error: %s", p.fileName, err)
		return nil
	}
	defer rc.Close()

	if err = json.NewDecoder(rc).Decode(&p.kv); err != nil {
		return err
	}

	return nil
}

// Set a key to the KeyValueStore with the given value.  If they key already
// exists, the value is overwritten.
func (p *KeyValueStore) Set(op trace.Operation, key string, value []byte) error {
	p.l.Lock()
	defer p.l.Unlock()

	// get the old value in case we need to roll back
	oldvalue, ok := p.kv[key]

	if ok && bytes.Compare(oldvalue, value) == 0 {
		// NOOP
		return nil
	}

	p.kv[key] = value

	if err := p.save(op); err != nil && ok {
		// revert if failure
		p.kv[key] = oldvalue
		return err
	}

	return nil
}

// Get retrieves a key from the KeyValueStore.
func (p *KeyValueStore) Get(op trace.Operation, key string) ([]byte, error) {
	p.l.RLock()
	defer p.l.RUnlock()

	v, ok := p.kv[key]
	if !ok {
		return []byte{}, ErrKeyNotFound
	}

	return v, nil
}

// Delete removes a key from the KeyValueStore.
func (p *KeyValueStore) Delete(op trace.Operation, key string) error {
	p.l.Lock()
	defer p.l.Unlock()

	oldvalue, ok := p.kv[key]
	if !ok {
		return ErrKeyNotFound
	}

	delete(p.kv, key)

	if err := p.save(op); err != nil {
		// restore the key
		p.kv[key] = oldvalue
		return err
	}

	return nil
}

// Save persists the KeyValueStore to the Backend.
func (p *KeyValueStore) Save(op trace.Operation) error {
	p.l.Lock()
	defer p.l.Unlock()
	return p.save(op)
}

func (p *KeyValueStore) save(op trace.Operation) error {
	buf, err := json.Marshal(p.kv)
	if err != nil {
		return err
	}

	// upload to an ephemeral file
	tmpfile := p.fileName + ".tmp"

	r := bytes.NewReader(buf)
	if err = p.b.Upload(op, r, tmpfile); err != nil {
		op.Errorf("Error uploading %s: %s", tmpfile, err)
		return err
	}

	op.Debugf("KeyValueStore(%s) Saving...", p.fileName)
	if err := p.b.Mv(op, tmpfile, p.fileName); err != nil {
		op.Errorf("Error moving %s: %s", tmpfile, err)
		return err
	}

	return nil
}
