// Copyright 2017 VMware, Inc. All Rights Reserved.
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

import "os"

const (
	NFS3_OK             = 0
	NFS3ERR_PERM        = 1
	NFS3ERR_NOENT       = 2
	NFS3ERR_IO          = 5
	NFS3ERR_NXIO        = 6
	NFS3ERR_ACCES       = 13
	NFS3ERR_EXIST       = 17
	NFS3ERR_XDEV        = 18
	NFS3ERR_NODEV       = 19
	NFS3ERR_NOTDIR      = 20
	NFS3ERR_ISDIR       = 21
	NFS3ERR_INVAL       = 22
	NFS3ERR_FBIG        = 27
	NFS3ERR_NOSPC       = 28
	NFS3ERR_ROFS        = 30
	NFS3ERR_MLINK       = 31
	NFS3ERR_NAMETOOLONG = 63
	NFS3ERR_NOTEMPTY    = 66
	NFS3ERR_DQUOT       = 69
	NFS3ERR_STALE       = 70
	NFS3ERR_REMOTE      = 71
	NFS3ERR_BADHANDLE   = 10001
	NFS3ERR_NOT_SYNC    = 10002
	NFS3ERR_BAD_COOKIE  = 10003
	NFS3ERR_NOTSUPP     = 10004
	NFS3ERR_TOOSMALL    = 10005
	NFS3ERR_SERVERFAULT = 10006
	NFS3ERR_BADTYPE     = 10007
)

var errToName = map[uint32]string{
	0:     "NFS3_OK",
	1:     "NFS3ERR_PERM",
	2:     "NFS3ERR_NOENT",
	5:     "NFS3ERR_IO",
	6:     "NFS3ERR_NXIO",
	13:    "NFS3ERR_ACCES",
	17:    "NFS3ERR_EXIST",
	18:    "NFS3ERR_XDEV",
	19:    "NFS3ERR_NODEV",
	20:    "NFS3ERR_NOTDIR",
	21:    "NFS3ERR_ISDIR",
	22:    "NFS3ERR_INVAL",
	27:    "NFS3ERR_FBIG",
	28:    "NFS3ERR_NOSPC",
	30:    "NFS3ERR_ROFS",
	31:    "NFS3ERR_MLINK",
	63:    "NFS3ERR_NAMETOOLONG",
	66:    "NFS3ERR_NOTEMPTY",
	69:    "NFS3ERR_DQUOT",
	70:    "NFS3ERR_STALE",
	71:    "NFS3ERR_REMOTE",
	10001: "NFS3ERR_BADHANDLE",
	10002: "NFS3ERR_NOT_SYNC",
	10003: "NFS3ERR_BAD_COOKIE",
	10004: "NFS3ERR_NOTSUPP",
	10005: "NFS3ERR_TOOSMALL",
	10006: "NFS3ERR_SERVERFAULT",
	10007: "NFS3ERR_BADTYPE",
}

func NFS3Error(errnum uint32) error {
	switch errnum {
	case NFS3_OK:
		return nil
	case NFS3ERR_PERM:
		return os.ErrPermission
	case NFS3ERR_EXIST:
		return os.ErrExist
	case NFS3ERR_NOENT:
		return os.ErrNotExist
	default:
		if errStr, ok := errToName[errnum]; ok {
			return &Error{
				ErrorNum:    errnum,
				ErrorString: errStr,
			}
		}

		return os.ErrInvalid
	}
}

// Error represents an unexpected I/O behavior.
type Error struct {
	ErrorNum    uint32
	ErrorString string
}

func (err *Error) Error() string { return err.ErrorString }

func IsNotEmptyError(err error) bool {
	nfsErr, ok := err.(*Error)
	if !ok {
		return false
	}

	if nfsErr.ErrorNum == NFS3ERR_NOTEMPTY {
		return true
	}

	return false
}

func IsNotDirError(err error) bool {
	nfsErr, ok := err.(*Error)
	if !ok {
		return false
	}

	if nfsErr.ErrorNum == NFS3ERR_NOTDIR {
		return true
	}

	return false
}
