package fs

import (
	"fmt"
	"github.com/ctdk/sbinfo"
	"io/ioutil"
	"os"
)

//this contains the fs block info that we are interested in.
type Fsinfo struct {
	DevicePath  string
	UUID        string
	VolumeLabel string
}

type fsDeviceManager interface {
	FindMountableBlockDevices(devicesPath string) ([]Fsinfo, error)
	GetDevicesByLabel(devicesPath string) (map[string]Fsinfo, error)
}

type Ext4DeviceManager struct{}

func (m *Ext4DeviceManager) FindMountableBlockDevices(devicesPath string) ([]Fsinfo, error) {
	ext4BlockDevices := make([]Fsinfo, 1, 1)

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

		if berr != nil {
			//ignore this block, it is probably not ext4
			continue
		}

		ext4BlockDevices = append(ext4BlockDevices, m.newExt4Fsinfo(info, blockPath))
	}

	return ext4BlockDevices, nil
}

func (m *Ext4DeviceManager) GetDeviceByLbel(devicesPath string) (map[string]Fsinfo, error) {
	deviceLabelMap := make(map[string]Fsinfo)

	devices, err := m.FindMountableBlockDevices(devicesPath)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		deviceLabelMap[device.VolumeLabel] = device
	}
	return deviceLabelMap, nil
}

func (m *Ext4DeviceManager) newExt4Fsinfo(info *sbinfo.Ext2Sb, devicePath string) Fsinfo {
	ext4info := Fsinfo{
		DevicePath:  devicePath,
		UUID:        string(info.SUUID[:]),
		VolumeLabel: string(info.SVolumeName[:]),
	}
	return ext4info
}
