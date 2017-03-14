package nfs

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/user"
	"syscall"
	"time"

	"github.com/fdawg4l/go-nfs-client/nfs/rpc"
	"github.com/fdawg4l/go-nfs-client/nfs/util"
)

const (
	Nfs3Prog = 100003
	Nfs3Vers = 3

	// program methods
	NFSProc3Lookup      = 3
	NFSProc3Read        = 6
	NFSProc3Write       = 7
	NFSProc3Create      = 8
	NFSProc3Mkdir       = 9
	NFSProc3Remove      = 12
	NFSProc3RmDir       = 13
	NFSProc3ReadDirPlus = 17
	NFSProc3FSInfo      = 19
	NFSProc3Commit      = 21

	// file types
	NF3Reg  = 1
	NF3Dir  = 2
	NF3Blk  = 3
	NF3Chr  = 4
	NF3Lnk  = 5
	NF3Sock = 6
	NF3FIFO = 7
)

type Diropargs3 struct {
	FH       []byte
	Filename string
}

type Sattr3 struct {
	Mode  SetMode
	UID   SetUID
	GID   SetUID
	Size  uint64
	Atime NFS3Time
	Mtime NFS3Time
}

type SetMode struct {
	Set  uint32
	Mode uint32
}

type SetUID struct {
	Set uint32
	UID uint32
}

type NFS3Time struct {
	Seconds  uint32
	Nseconds uint32
}

type Fattr struct {
	Type                uint32
	FileMode            uint32
	Nlink               uint32
	UID                 uint32
	GID                 uint32
	Filesize            uint64
	Used                uint64
	SpecData            [2]uint32
	FSID                uint64
	Fileid              uint64
	Atime, Mtime, Ctime NFS3Time
}

func (f *Fattr) Name() string {
	return ""
}

func (f *Fattr) Size() int64 {
	return int64(f.Filesize)
}

func (f *Fattr) Mode() os.FileMode {
	return os.FileMode(f.FileMode)
}

func (f *Fattr) ModTime() time.Time {
	return time.Unix(int64(f.Mtime.Seconds), int64(f.Mtime.Nseconds))
}

func (f *Fattr) IsDir() bool {
	return f.Type == NF3Dir
}

func (f *Fattr) Sys() interface{} {
	return nil
}

type EntryPlus struct {
	FileId   uint64
	FileName string
	Cookie   uint64
	Attr     struct {
		Follows uint32
		Attr    Fattr
	}
	FHSet        uint32
	FH           []byte
	ValueFollows uint32
}

func (e *EntryPlus) Name() string {
	return e.FileName
}

func (e *EntryPlus) Size() int64 {
	if e.Attr.Follows == 0 {
		return 0
	}

	return e.Attr.Attr.Size()
}

func (e *EntryPlus) Mode() os.FileMode {
	if e.Attr.Follows == 0 {
		return 0
	}

	return e.Attr.Attr.Mode()
}

func (e *EntryPlus) ModTime() time.Time {
	if e.Attr.Follows == 0 {
		return time.Time{}
	}

	return e.Attr.Attr.ModTime()
}

func (e *EntryPlus) IsDir() bool {
	if e.Attr.Follows == 0 {
		return false
	}

	return e.Attr.Attr.IsDir()
}

func (e *EntryPlus) Sys() interface{} {
	if e.Attr.Follows == 0 {
		return 0
	}

	return e.FileId
}

type FSInfo struct {
	Follows    uint32
	RTMax      uint32
	RTPref     uint32
	RTMult     uint32
	WTMax      uint32
	WTPref     uint32
	WTMult     uint32
	DTPref     uint32
	Size       uint64
	TimeDelta  NFS3Time
	Properties uint32
}

// Dial an RPC svc after getting the port from the portmapper
func DialService(addr string, prog rpc.Mapping) (*rpc.Client, error) {
	pm, err := rpc.DialPortmapper("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer pm.Close()

	port, err := pm.Getport(prog)
	if err != nil {
		return nil, err
	}

	client, err := dialService(addr, port)

	return client, nil
}

func dialService(addr string, port int) (*rpc.Client, error) {
	var (
		ldr    *net.TCPAddr
		client *rpc.Client
	)

	usr, err := user.Current()

	// Unless explicitly configured, the target will likely reject connections
	// from non-privileged ports.
	if err == nil && usr.Uid == "0" {
		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

		var p int
		for {
			p = r1.Intn(1024)
			if p < 0 {
				continue
			}

			ldr = &net.TCPAddr{
				Port: p,
			}

			client, err = rpc.DialTCP("tcp", ldr, fmt.Sprintf("%s:%d", addr, port))
			if err == nil {
				break
			}
			// bind error, try again
			if isAddrInUse(err) {
				continue
			}

			return nil, err
		}

		util.Debugf("using random port %d -> %d", p, port)
	} else {

		client, err = rpc.DialTCP("tcp", ldr, fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func isAddrInUse(err error) bool {
	if er, ok := (err.(*net.OpError)); ok {
		if syser, ok := er.Err.(*os.SyscallError); ok {
			return syser.Err == syscall.EADDRINUSE
		}
	}

	return false
}
