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

package tether

import (
	"errors"

	"github.com/vmware/vic/pkg/trace"
)

const (
	pidFilePath = "Temp"
)

func (t *tether) childReaper() error {
	// TODO: windows child process notifications
	return errors.New("Child reaping unimplemented on windows")
}

func (t *tether) stopReaper() {
	defer trace.End(trace.Begin("Shutting down child reaping"))
}

func lookPath(file string, env []string, dir string) (string, error) {
	return "", errors.New("unimplemented on windows")
}

func establishPty(session *SessionConfig) error {
	return errors.New("unimplemented on windows")
}

func establishNonPty(session *SessionConfig) error {
	return errors.New("unimplemented on windows")
}
