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

package vsphere

import (
	"io"
)

const (
	// Persistent is used as a constant input for disk persistence configuration
	persistent = true
)

// cleanReader wraps an io.ReadCloser, calling the supplied cleanup function on Close()
type cleanReader struct {
	io.ReadCloser
	clean func()
}

func (c *cleanReader) Read(p []byte) (int, error) {
	return c.ReadCloser.Read(p)
}

func (c *cleanReader) Close() error {
	defer c.clean()
	return c.ReadCloser.Close()
}
