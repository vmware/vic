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

// Import takes a tar archive stream and extracts it into the target volume
func (v *VolumeStore) Import(op trace.Operation, id string, spec *archive.FilterSpec, tarStream io.ReadCloser) error {
	return fmt.Errorf("Write for nfs volumes is not Implemented")
}

func (v *VolumeStore) NewDataSink(op trace.Operation, id string) (storage.DataSink, error) {
	return nil, errors.New("NFS VolumeStore does not yet implement NewDataSink")
}
