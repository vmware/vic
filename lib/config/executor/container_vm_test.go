// Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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

package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnv(t *testing.T) {
	sess := SessionConfig{
		Cmd: Cmd{
			Env: []string{
				"hello=world",
				"goodbye=",
			},
		},
	}

	assert.Equal(t, "world", *sess.GetEnv("hello"), "Get on set variable should return set value")
	assert.Equal(t, "", *sess.GetEnv("goodbye"), "Get on set variable with empty value should return empty string")
	assert.Nil(t, sess.GetEnv("nope"), "Get on an unset variable should return nil")
}

func TestSetEnvUpdateValue(t *testing.T) {
	sess := SessionConfig{
		Cmd: Cmd{
			Env: []string{
				"hello=world",
				"goodbye=",
			},
		},
	}

	require.Equal(t, "world", *sess.GetEnv("hello"), "Get on set variable should return set value")
	require.Equal(t, "", *sess.GetEnv("goodbye"), "Get on set variable with empty value should return empty string")

	newVal := "sapients"
	old := *sess.SetEnv("hello", newVal)

	assert.Equal(t, "world", old, "Expected old value to be return on update")
	assert.Equal(t, newVal, *sess.GetEnv("hello"), "Expected new value to be updated value")

	assert.Equal(t, "", *sess.GetEnv("goodbye"), "Expected unmodified value not to have changed")
}

func TestSetEnvEmptyValue(t *testing.T) {
	sess := SessionConfig{
		Cmd: Cmd{
			Env: []string{
				"hello=world",
				"goodbye=",
			},
		},
	}

	require.Equal(t, "world", *sess.GetEnv("hello"), "Get on set variable should return set value")
	require.Equal(t, "", *sess.GetEnv("goodbye"), "Get on set variable with empty value should return empty string")

	newVal := "sapients"
	old := *sess.SetEnv("goodbye", newVal)

	assert.Equal(t, "", old, "Expected old value to be return on update")
	assert.Equal(t, newVal, *sess.GetEnv("goodbye"), "Expected new value to be updated value")

	assert.Equal(t, "world", *sess.GetEnv("hello"), "Expected unmodified value not to have changed")
}

func TestSetEnvNewEnv(t *testing.T) {
	sess := SessionConfig{
		Cmd: Cmd{
			Env: []string{
				"hello=world",
				"goodbye=",
			},
		},
	}

	require.Equal(t, "world", *sess.GetEnv("hello"), "Get on set variable should return set value")
	require.Equal(t, "", *sess.GetEnv("goodbye"), "Get on set variable with empty value should return empty string")

	newKey := "solo"
	require.Nil(t, sess.GetEnv(newKey), "Expected nil return for unset value")

	// checking we can set a new env with value of empty string
	newVal := ""
	_ = sess.SetEnv(newKey, newVal)
	assert.Equal(t, newVal, *sess.GetEnv(newKey), "Expected new value to be updated value")

	// checking we can set a new env with value specified
	newKey = "absence"
	newVal = "fonder"

	_ = sess.SetEnv(newKey, newVal)
	assert.Equal(t, newVal, *sess.GetEnv(newKey), "Expected new value to be updated value")

	assert.Equal(t, "world", *sess.GetEnv("hello"), "Expected unmodified value not to have changed")
	assert.Equal(t, "", *sess.GetEnv("goodbye"), "Expected unmodified value not to have changed")
}
