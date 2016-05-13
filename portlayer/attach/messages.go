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

import "golang.org/x/crypto/ssh"

// All of the messages passed over the ssh channel/global mux are (or will be)
// defined here.

type Message interface {
	// Returns the message name
	RequestType() string

	// Marshalled version of the message
	Marshal() []byte

	// Unmarshal unpacks the message
	Unmarshal([]byte) error
}

// WindowChangeMsg the RFC4254 struct
const WindowChangeReq = "window-change"

type WindowChangeMsg struct {
	Columns  uint32
	Rows     uint32
	WidthPx  uint32
	HeightPx uint32
}

func (wc *WindowChangeMsg) RequestType() string {
	return WindowChangeReq
}

func (wc *WindowChangeMsg) Marshal() []byte {
	return ssh.Marshal(*wc)
}

func (wc *WindowChangeMsg) Unmarshal(payload []byte) error {
	return ssh.Unmarshal(payload, wc)
}
