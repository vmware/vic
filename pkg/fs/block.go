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

package fs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ctdk/sbinfo"
)

//this contains the fs block info that we are interested in.
type DeviceInfo struct {
	DevicePath  string
	UUID        string
	VolumeLabel string
}

type FileSystem interface {
	DevPath() string
	UUID() string
	Label() string
	Info() DeviceInfo
}

type Ext4FileSystem struct {
	FilesystemInfo DeviceInfo
}

func (fs Ext4FileSystem) DevPath() string {
	return fs.FilesystemInfo.DevicePath
}

func (fs Ext4FileSystem) UUID() string {
	return fs.FilesystemInfo.UUID
}

func (fs Ext4FileSystem) Label() string {
	return fs.FilesystemInfo.VolumeLabel
}

func (fs Ext4FileSystem) Info() DeviceInfo {
	return fs.FilesystemInfo
}

type BlockDevices struct {
	FileSystems []FileSystem
}

func NewBlockDevices(DevicesPath string) (BlockDevices, error) {
	devices := BlockDevices{}

	_, err := os.Stat(DevicesPath)
	if err != nil {
		return devices, err
	}

	deviceTypes := []SuperBlock{Ext4SuperBlock{}}

	//TODO: make this more efficient we should not be trying to read
	//      blocks that are already identified.
	for _, sbReader := range deviceTypes {
		fs, err := sbReader.Read(DevicesPath)
		if err != nil {
			devices.FileSystems = append(devices.FileSystems, fs...)
		}
	}
	return devices, nil
}

func (b *BlockDevices) DevicesByLabel() map[string]FileSystem {
	deviceLabelMap := make(map[string]FileSystem)

	for _, device := range b.FileSystems {
		deviceLabelMap[device.Label()] = device
	}
	return deviceLabelMap
}

//interface for identifying super block information
type SuperBlock interface {
	Read(devicePath string) ([]FileSystem, error)
}

//ext4 super block implementation
type Ext4SuperBlock struct{}

func (m Ext4SuperBlock) Read(devicesPath string) ([]FileSystem, error) {
	ext4BlockDevices := make([]FileSystem, 1, 1)

	//first make sure devicesPath is valid
	devicesDir, err := os.Stat(devicesPath)
	if err != nil {
		return nil, err
	}

	if !devicesDir.IsDir() {
		return nil, fmt.Errorf("Supplied device path is not a directory")
	}

	blockDevices, err := ioutil.ReadDir(devicesPath)
	if err != nil {
		return nil, err
	}

	for _, block := range blockDevices {
		blockPath := fmt.Sprintf("%s/%s", devicesPath, block.Name())
		info, berr := sbinfo.ReadExt2Superblock(blockPath)
		fs := Ext4FileSystem{
			FilesystemInfo: m.newExt4DeviceInfo(info, blockPath),
		}

		if berr != nil {
			//ignore this block, it is probably not ext4
			continue
		}

		ext4BlockDevices = append(ext4BlockDevices, fs)
	}

	return ext4BlockDevices, nil
}

func (m *Ext4SuperBlock) newExt4DeviceInfo(info *sbinfo.Ext2Sb, devicePath string) DeviceInfo {
	ext4info := DeviceInfo{
		DevicePath:  devicePath,
		UUID:        string(info.SUUID[:]),
		VolumeLabel: string(info.SVolumeName[:]),
	}
	return ext4info
}
