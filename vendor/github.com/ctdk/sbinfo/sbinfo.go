//Package sbinfo reads information from a filesystem's superblock and returns a
// struct representing that superblock. At the moment only works with ext2/3/4.
package sbinfo

import (
	"fmt"
	"os"
	"unsafe"
)

const ext2Magic uint16 = 0xef53

// ReadExt2Superblock reads the superblock from the given device or file and
// returns a struct representing the superblock, or an error on failure.
func ReadExt2Superblock(devPath string) (*Ext2Sb, error) {
	fp, err := os.Open(devPath)
	if err != nil {
		return nil, err
	}
	// read in the first 2k bytes
	buf := make([]byte, 2048)
	n, err := fp.Read(buf)
	if n != 2048 {
		nerr := fmt.Errorf("Expected 2048 bytes, only read %d from %s", n, devPath)
		return nil, nerr
	}
	if err != nil {
		return nil, err
	}
	sbRaw := buf[1024:]
	sb := (*Ext2Sb)(unsafe.Pointer(&sbRaw[0]))
	// This needs to be able to check if it's running on a bigendian system
	// and switch up the sb values as needed.
	if sb.SMagic != ext2Magic {
		sbErr := fmt.Errorf("Bad magic number for %s", devPath)
		return nil, sbErr
	}
	return sb, nil
}
