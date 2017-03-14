package nfs

import (
	"io"
	"os"

	"github.com/fdawg4l/go-nfs-client/nfs/rpc"
	"github.com/fdawg4l/go-nfs-client/nfs/util"
	"github.com/fdawg4l/go-nfs-client/nfs/xdr"
)

type File struct {
	*Target

	// current position
	curr   uint64
	fsinfo *FSInfo

	// filehandle to the file
	fh []byte
}

func (f *File) Read(p []byte) (int, error) {
	type ReadArgs struct {
		rpc.Header
		FH     []byte
		Offset uint64
		Count  uint32
	}

	type ReadRes struct {
		Follows uint32
		Attrs   struct {
			Attrs Fattr
		}
		Count uint32
		Eof   uint32
		Data  struct {
			Length uint32
		}
	}

	readSize := uint32(min(uint64(f.fsinfo.RTPref), uint64(len(p))))
	util.Debugf("read(%x) len=%d offset=%d", f.fh, readSize, f.curr)

	r, err := f.call(&ReadArgs{
		Header: rpc.Header{
			Rpcvers: 2,
			Prog:    Nfs3Prog,
			Vers:    Nfs3Vers,
			Proc:    NFSProc3Read,
			Cred:    f.auth,
			Verf:    rpc.AuthNull,
		},
		FH:     f.fh,
		Offset: uint64(f.curr),
		Count:  readSize,
	})

	if err != nil {
		util.Debugf("read(%x): %s", f.fh, err.Error())
		return 0, err
	}

	readres := &ReadRes{}
	if err = xdr.Read(r, readres); err != nil {
		return 0, err
	}

	f.curr = f.curr + uint64(readres.Data.Length)
	n, err := r.Read(p[:readres.Data.Length])
	if err != nil {
		return n, err
	}

	if readres.Eof != 0 {
		err = io.EOF
	}

	return n, err
}

func (f *File) Write(p []byte) (int, error) {
	type WriteArgs struct {
		rpc.Header
		FH     []byte
		Offset uint64
		Count  uint32

		// UNSTABLE(0), DATA_SYNC(1), FILE_SYNC(2) default
		How      uint32
		Contents []byte
	}

	totalToWrite := len(p)
	writeSize := uint64(min(uint64(f.fsinfo.WTPref), uint64(totalToWrite)))

	_, err := f.call(&WriteArgs{
		Header: rpc.Header{
			Rpcvers: 2,
			Prog:    Nfs3Prog,
			Vers:    Nfs3Vers,
			Proc:    NFSProc3Write,
			Cred:    f.auth,
			Verf:    rpc.AuthNull,
		},
		FH:       f.fh,
		Offset:   f.curr,
		Count:    uint32(writeSize),
		How:      2,
		Contents: p[:writeSize],
	})

	if err != nil {
		util.Debugf("write(%x): %s", f.fh, err.Error())
		return int(writeSize), err
	}

	util.Debugf("write(%x) len=%d offset=%d written=%d total=%d",
		f.fh, writeSize, f.curr, writeSize, totalToWrite)

	f.curr = f.curr + writeSize

	return int(writeSize), nil
}

func (f *File) Close() error {
	type CommitArg struct {
		rpc.Header
		FH     []byte
		Offset uint64
		Count  uint32
	}

	_, err := f.call(&CommitArg{
		Header: rpc.Header{
			Rpcvers: 2,
			Prog:    Nfs3Prog,
			Vers:    Nfs3Vers,
			Proc:    NFSProc3Commit,
			Cred:    f.auth,
			Verf:    rpc.AuthNull,
		},
		FH: f.fh,
	})

	if err != nil {
		util.Debugf("commit(%x): %s", f.fh, err.Error())
		return err
	}

	return nil
}

// OpenFile writes to an existing file or creates one
func (v *Target) OpenFile(path string, perm os.FileMode) (io.ReadWriteCloser, error) {
	_, fh, err := v.Lookup(path)
	if err != nil {
		if os.IsNotExist(err) {
			fh, err = v.Create(path, perm)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	f := &File{
		Target: v,
		fsinfo: v.fsinfo,
		fh:     fh,
	}

	return f, nil
}

// Open opens a file for reading
func (v *Target) Open(path string) (io.ReadCloser, error) {
	_, fh, err := v.Lookup(path)
	if err != nil {
		return nil, err
	}

	f := &File{
		Target: v,
		fsinfo: v.fsinfo,
		fh:     fh,
	}

	return f, nil
}

func min(x, y uint64) uint64 {
	if x > y {
		return y
	}
	return x
}
