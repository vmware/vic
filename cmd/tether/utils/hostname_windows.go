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

// +build !mock,!mock_hostname

package utils

import (
	"errors"

	"github.com/vmware/vic/pkg/trace"
)

// SetHostname sets the system hostname
func SetHostname(hostname string) error {
	defer trace.End(trace.Begin("setting hostname to " + hostname))

	return errors.New("not implemented on windows")
}
