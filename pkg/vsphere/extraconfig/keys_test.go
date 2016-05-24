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

package extraconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func visibleRO(key string) string {
	return calculateKey([]string{"read-only"}, "", key)
}

func visibleRW(key string) string {
	return calculateKey([]string{"read-write"}, "", key)
}

func hidden(key string) string {
	return calculateKey([]string{"hidden"}, "", key)
}

func TestHidden(t *testing.T) {
	scopes := []string{"hidden"}

	key := calculateKey(scopes, "a/b", "c")

	assert.Equal(t, "a/b/c", key, "Key should remain hidden")
}

func TestHide(t *testing.T) {
	scopes := []string{"hidden"}

	key := calculateKey(scopes, DefaultGuestInfoPrefix+"/a/b", "c")

	assert.Equal(t, "a/b/c", key, "Key should be hidden")
}

func TestReveal(t *testing.T) {
	scopes := []string{"read-only"}

	key := calculateKey(scopes, "a/b", "c")

	assert.Equal(t, DefaultGuestInfoPrefix+"/a/b/c", key, "Key should be exposed")
}

func TestVisibleReadOnly(t *testing.T) {
	scopes := []string{"read-only"}

	key := calculateKey(scopes, DefaultGuestInfoPrefix+"/a/b", "c")

	assert.Equal(t, DefaultGuestInfoPrefix+"/a/b/c", key, "Key should be remain visible and read-only")
}

func TestVisibleReadWrite(t *testing.T) {
	scopes := []string{"read-write"}

	key := calculateKey(scopes, DefaultGuestInfoPrefix+".a.b", "c")

	assert.Equal(t, DefaultGuestInfoPrefix+".a.b.c", key, "Key should be remain visible and read-write")
}

func TestTopLevelReadOnly(t *testing.T) {
	scopes := []string{"read-only"}

	key := calculateKey(scopes, "", "a")

	assert.Equal(t, DefaultGuestInfoPrefix+"/a", key, "Key should be visible and read-only")
}

func TestReadOnlyToReadWrite(t *testing.T) {
	scopes := []string{"read-write"}

	key := calculateKey(scopes, DefaultGuestInfoPrefix+"/a/b", "c")

	assert.Equal(t, DefaultGuestInfoPrefix+".a.b.c", key, "Key should be visible and change to read-write")
}

func TestReadWriteToReadOnly(t *testing.T) {
	scopes := []string{"read-only"}

	key := calculateKey(scopes, DefaultGuestInfoPrefix+".a.b", "c")

	assert.Equal(t, DefaultGuestInfoPrefix+"/a/b/c", key, "Key should be visible and change to read-only")
}

func TestCompoundKey(t *testing.T) {
	scopes := []string{"read-write"}

	key := calculateKey(scopes, DefaultGuestInfoPrefix+".a", "b/c")

	assert.Equal(t, DefaultGuestInfoPrefix+".a.b.c", key, "Key should be visible and read-write")
}

func TestNoScopes(t *testing.T) {
	scopes := []string{}

	key := calculateKey(scopes, DefaultGuestInfoPrefix+".a/b", "c")
	assert.Equal(t, "a/b/c", key, "Key should be completely proscriptive")

	key = calculateKey(scopes, DefaultGuestInfoPrefix+".a.b", "c")
	assert.Equal(t, "a.b/c", key, "Key should be hidden")

	key = calculateKey(scopes, "a.b", "c")
	assert.Equal(t, "a.b/c", key, "Key should remain hidden")
}
