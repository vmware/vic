// Copyright 2018 VMware, Inc. All Rights Reserved.
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

package nfs

import (
	"errors"
	"fmt"
	"io"

	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"
)

// Export reads the delta between child and parent volume layers, returning
// the difference as a tar archive.
//
// store - the volume store containing the two layers
// id - must inherit from ancestor if ancestor is specified
// ancestor - the volume layer up the chain against which to diff
// spec - describes filters on paths found in the data (include, exclude, strip)
// data - set to true to include file data in the tar archive, false to include headers only
// Export creates and returns a tar archive containing data found between an nfs layer one or all of its ancestors
func (v *VolumeStore) Export(op trace.Operation, id, ancestor string, spec *archive.FilterSpec, data bool) (io.ReadCloser, error) {
	return nil, fmt.Errorf("vSphere Integrated Containers does not yet implement Export for nfs volumes")
}

func (v *VolumeStore) NewDataSource(op trace.Operation, id string) (storage.DataSource, error) {
	return nil, errors.New("NFS VolumeStore does not yet implement NewDataSource")
}
