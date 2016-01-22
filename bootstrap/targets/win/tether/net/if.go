// +build linux,!generic

package net

import (
	"os"
	"syscall"
)

type sockaddr struct {
	family uint16
	port   uint16
	addr   [20]byte
}

type ifreq struct {
	ifr_name [syscall.IFNAMSIZ]byte
	flags    uint16
}

type ifconf struct {
	rlen uint32  // lenght (in bytes) of ifr array
	pad  uint32  //
	ifr  uintptr // ptr to ifr array
}

// Convert a null (0x00) terminated C string into a Go string
func cStrToString(b []byte) string {
	clen := 0
	for ; clen < len(b); clen++ {
		if b[clen] == 0 {
			break
		}
	}
	return string(b[:clen])
}

// cribbed from https://github.com/davecheney/pcap/blob/master/bpf.go#L43 for now
func ioctl(fd int, request, argp uintptr) error {
	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, argp)
	return os.NewSyscallError("ioctl", int(errorp))
}

func ifup(name string) error {
	fd, fderr := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)

	// get currently flags
	request := ifreq{ifr_name: name}
	unsafe := unsafe.Pointer(&request)
	err := ioctl(fd, syscall.SIOCGIFFLAGS, uintptr(unsafe))
	if err != nil {
		return err
	}

	flags := (*uint16)(unsafe.Pointer((uintptr(unsafe) + syscall.IFNAMSIZ)))

	// add the activation flags
	(*flags) |= IFF_UP & IFF_RUNNING

	// bring the interface up
	err := ioctl(fd, syscall.SIOCSIFFLAGS, uintptr(unsafe))
	if err != nil {
		return err
	}

	return nil
}
