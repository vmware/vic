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

package fs

import "github.com/vmware/vic/pkg/trace"

// FilesystemType represents the filesystem in use by a virtual disk
type FilesystemType string

const (
	// Ext4 represents the ext4 file system
	TypeExt4 FilesystemType = "ext4"

	// Xfs represents the XFS file system
	TypeXfs FilesystemType = "xfs"

	// Ntfs represents the NTFS file system
	TypeNtfs FilesystemType = "ntfs"

	// Raw represents no filesystem - raw device
	TypeRaw FilesystemType = "raw"
)

var Filesystems = map[FilesystemType]Filesystem{
	TypeRaw:  NewRaw(),
	TypeExt4: NewExt4(),
	TypeXfs:  NewXFS(),
}

// Filesystem defines the interface for handling an attached virtual disk
type Filesystem interface {
	Mkfs(op trace.Operation, devPath, label string) error
	SetLabel(op trace.Operation, devPath, labelName string) error
	Mount(op trace.Operation, devPath, targetPath string, options []string) error
	Unmount(op trace.Operation, path string) error
}

func GetFilesystem(name FilesystemType) Filesystem {
	return Filesystems[name]
}
