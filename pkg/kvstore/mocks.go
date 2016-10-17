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
	"io"
	"io/ioutil"
	"os"
)

type MockBackend struct {
	buf []byte
}

// Creates path and ovewrites whatever is there
func (m *MockBackend) Save(ctx context.Context, r io.Reader, pth string) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	m.buf = buf

	return nil
}

func (m *MockBackend) Load(ctx context.Context, pth string) (io.ReadCloser, error) {
	if len(m.buf) == 0 {
		return nil, os.ErrNotExist
	}

	return ioutil.NopCloser(bytes.NewReader(m.buf)), nil
}
