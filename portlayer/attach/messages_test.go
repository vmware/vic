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

package attach

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindowChange(t *testing.T) {
	s := &WindowChangeMsg{1, 2, 3, 4}

	assert.Equal(t, s.RequestType(), WindowChangeReq)

	tmp := s.Marshal()
	out := &WindowChangeMsg{}
	out.Unmarshal(tmp)

	assert.Equal(t, s, out)
}

func TestSignal(t *testing.T) {
	s := &SignalMsg{"HUP"}

	assert.Equal(t, s.RequestType(), SignalReq)

	tmp := s.Marshal()
	out := &SignalMsg{}
	out.Unmarshal(tmp)

	assert.Equal(t, s, out)
}
