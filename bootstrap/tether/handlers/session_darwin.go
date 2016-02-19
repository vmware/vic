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

package handlers

import (
	"github.com/vmware/vic/bootstrap/tether"

	"golang.org/x/crypto/ssh"
)

func (ch *SessionHandler) AssignPty() {
}

func (ch *SessionHandler) ResizePty(winSize *tether.WindowChangeMsg) error {
	return nil
}

func (ch *SessionHandler) Signal(sig ssh.Signal) error {
	return nil
}

func (ch *SessionHandler) Exec(command string, args []string, config map[string]string) (ok bool, payload []byte) {
	return false, nil
}
