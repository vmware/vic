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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/pkg/trace"
)

type MockBackend struct {
	buf []byte
}

// Creates path and ovewrites whatever is there
func (m *MockBackend) Upload(ctx context.Context, r io.Reader, pth string) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	m.buf = buf

	return nil
}

func (m *MockBackend) Download(ctx context.Context, pth string) (io.ReadCloser, error) {
	if len(m.buf) == 0 {
		return nil, os.ErrNotExist
	}

	return ioutil.NopCloser(bytes.NewReader(m.buf)), nil
}

func (m *MockBackend) Mv(ctx context.Context, fromPath, toPath string) error {
	return nil
}

func save(t *testing.T, kv *KeyValueStore, key string, expectedvalue []byte) {
	op := trace.NewOperation(context.Background(), "save")

	if !assert.NoError(t, kv.Set(op, key, expectedvalue)) {
		return
	}
}

func get(t *testing.T, kv *KeyValueStore, key string, expectedval []byte) {
	op := trace.NewOperation(context.Background(), "get")

	// get the value we added
	v, err := kv.Get(op, key)
	if !assert.NoError(t, err) || !assert.NotNil(t, v) || !assert.Equal(t, expectedval, v) {
		return
	}
}

func TestAddAndGet(t *testing.T) {
	mb := &MockBackend{}

	op := trace.NewOperation(context.Background(), "testaddsaveget")

	// Save some entries in parallel
	entries := 500
	wg := sync.WaitGroup{}
	wg.Add(entries)

	expected := make(map[string][]byte)

	firstkv, err := NewKeyValueStore(op, mb, "datfile")
	if !assert.NoError(t, err) || !assert.NotNil(t, firstkv) {
		return
	}

	for i := 0; i < entries; i++ {
		k := fmt.Sprintf("key-%d", i)
		v := []byte(strconv.Itoa(i))

		expected[k] = v
		go func() {
			defer wg.Done()
			save(t, firstkv, k, v)
			get(t, firstkv, k, v)
		}()
	}
	wg.Wait()

	if t.Failed() {
		return
	}

	// Restart the kv store by creating a new one and attempt to get the same
	// entries.
	secondkv, err := NewKeyValueStore(op, mb, "datfile")
	if !assert.NoError(t, err) || !assert.NotNil(t, secondkv) {
		return
	}

	wg.Add(entries)
	for k, v := range expected {
		go func(key string, value []byte) {
			defer wg.Done()
			get(t, secondkv, key, value)
		}(k, v)
	}
	wg.Wait()

	if t.Failed() {
		return
	}

	// Ovewrite all of the values and verify again
	wg.Add(entries)
	for k := range expected {
		newval := []byte("ddddd")

		expected[k] = newval
		go func(key string, value []byte) {
			defer wg.Done()
			save(t, secondkv, key, value)
			get(t, secondkv, key, value)
		}(k, newval)
	}
	wg.Wait()

	if t.Failed() {
		return
	}

	// Restart and verify the overwritten values match the expected
	thirdkv, err := NewKeyValueStore(op, mb, "datfile")
	if !assert.NoError(t, err) || !assert.NotNil(t, thirdkv) {
		return
	}

	wg.Add(entries)
	for k, v := range expected {
		go func(key string, value []byte) {
			defer wg.Done()
			get(t, thirdkv, key, value)
		}(k, v)
	}
	wg.Wait()

	if t.Failed() {
		return
	}

	// Remove all of the entries and assert nothing can be found
	wg.Add(entries)
	for k := range expected {
		go func(key string) {
			defer wg.Done()
			if !assert.NoError(t, thirdkv.Delete(op, key)) {
				return
			}

			_, err := thirdkv.Get(op, key)
			if !assert.Error(t, err) {
				return
			}
		}(k)
	}
	wg.Wait()

	if t.Failed() {
		return
	}

	// Check the kv is empty after restart
	fourthkv, err := NewKeyValueStore(op, mb, "datfile")
	wg.Add(entries)
	for k := range expected {
		go func(key string) {
			defer wg.Done()
			_, err := fourthkv.Get(op, key)
			if !assert.Error(t, err) {
				return
			}
		}(k)
	}
	wg.Wait()
}
