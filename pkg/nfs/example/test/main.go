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

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/nfs"
	"github.com/vmware/vic/pkg/nfs/rpc"
)

func main() {
	log.SetLevel(log.DebugLevel)
	if len(os.Args) != 3 {
		log.Infof("%s <host>:<target path> <test directory to be created>", os.Args[0])
		os.Exit(-1)
	}

	b := strings.Split(os.Args[1], ":")

	host := b[0]
	target := b[1]
	dir := os.Args[2]

	log.Infof("host=%s target=%s dir=%s\n", host, target, dir)

	mount, err := nfs.DialMount(host)
	if err != nil {
		log.Fatalf("unable to dial MOUNT service: %v", err)
	}
	defer mount.Close()

	auth := rpc.NewAuthUnix("hasselhoff", 1001, 1001)

	v, err := mount.Mount(target, auth.Auth())
	if err != nil {
		log.Fatalf("unable to mount volume: %v", err)
	}
	defer v.Close()

	if _, err = v.Mkdir(dir, 0775); err != nil {
		log.Fatalf("mkdir error: %v", err)
	}

	if _, err = v.Mkdir(dir, 0775); err == nil {
		log.Fatalf("mkdir expected error")
	}

	// make a nested dir
	if _, err = v.Mkdir(dir+"/a", 0775); err != nil {
		log.Fatalf("mkdir error: %v", err)
	}

	// make a nested dir
	if _, err = v.Mkdir(dir+"/a/b", 0775); err != nil {
		log.Fatalf("mkdir error: %v", err)
	}

	dirs, err := ls(v, ".")
	if err != nil {
		log.Fatalf("ls: %s", err.Error())
	}

	// check the length.  There should only be 1 entry in the target (aside from . and ..)
	if len(dirs) != 3 {
		log.Fatalf("expected 3 dirs, got %d", len(dirs))
	}

	// 10 MB file
	if err = testFileRW(v, "10mb", 10*1024*1024); err != nil {
		log.Fatalf("fail")
	}

	// 7b file
	if err = testFileRW(v, "7b", 7); err != nil {
		log.Fatalf("fail")
	}

	// should return an error
	if err = v.RemoveAll("7b"); err == nil {
		log.Fatalf("expected a NOTADIR error")
	} else {
		nfserr := err.(*nfs.Error)
		if nfserr.ErrorNum != nfs.NFS3ERR_NOTDIR {
			log.Fatalf("Wrong error")
		}
	}

	if err = v.Remove("7b"); err != nil {
		log.Fatalf("rm(7b) err: %s", err.Error())
	}

	if err = v.Remove("10mb"); err != nil {
		log.Fatalf("rm(10mb) err: %s", err.Error())
	}

	_, _, err = v.Lookup(dir)
	if err != nil {
		log.Fatalf("lookup error: %s", err.Error())
	}

	if _, err = ls(v, "."); err != nil {
		log.Fatalf("ls: %s", err.Error())
	}

	if err = v.RmDir(dir); err == nil {
		log.Fatalf("expected not empty error")
	}

	for _, fname := range []string{"/one", "/two", "/a/one", "/a/two", "/a/b/one", "/a/b/two"} {
		if err = testFileRW(v, dir+fname, 10); err != nil {
			log.Fatalf("fail")
		}
	}

	if err = v.RemoveAll(dir); err != nil {
		log.Fatalf("error removing files: %s", err.Error())
	}

	outDirs, err := ls(v, ".")
	if err != nil {
		log.Fatalf("ls: %s", err.Error())
	}

	if len(outDirs) != 2 {
		log.Fatalf("directory should be empty!")
	}

	if err = mount.Unmount(); err != nil {
		log.Fatalf("unable to umount target: %v", err)
	}

	if err = mount.Close(); err != nil {
		log.Fatalf("error unmounting: %s", err.Error())
	}

	log.Infof("PASSED")
}

func testFileRW(v *nfs.Target, name string, filesize uint64) error {

	// create a temp file
	f, err := os.Open("/dev/urandom")
	if err != nil {
		log.Errorf("error openning random: %s", err.Error())
		return err
	}

	wr, err := v.OpenFile(name, 0777)
	if err != nil {
		log.Errorf("write fail: %s", err.Error())
		return err
	}

	// calculate the sha
	h := sha256.New()
	t := io.TeeReader(f, h)

	// Copy filesize
	_, err = io.CopyN(wr, t, int64(filesize))
	if err != nil {
		log.Errorf("error copying: %s", err.Error())
		return err
	}
	expectedSum := h.Sum(nil)

	if err = wr.Close(); err != nil {
		log.Errorf("error committing: %s", err.Error())
		return err
	}

	//
	// get the file we wrote and calc the sum
	rdr, err := v.Open(name)
	if err != nil {
		log.Errorf("read error: %v", err)
		return err
	}

	h = sha256.New()
	t = io.TeeReader(rdr, h)

	_, err = ioutil.ReadAll(t)
	if err != nil {
		log.Errorf("readall error: %v", err)
		return err
	}
	actualSum := h.Sum(nil)

	if bytes.Compare(actualSum, expectedSum) != 0 {
		log.Fatalf("sums didn't match. actual=%x expected=%s", actualSum, expectedSum) //  Got=0%x expected=0%x", string(buf), testdata)
	}

	log.Printf("Sums match %x %x", actualSum, expectedSum)
	return nil
}

func ls(v *nfs.Target, path string) ([]*nfs.EntryPlus, error) {
	dirs, err := v.ReadDirPlus(path)
	if err != nil {
		return nil, fmt.Errorf("readdir error: %s", err.Error())
	}

	log.Infof("dirs:")
	for _, dir := range dirs {
		log.Infof("\t%s\t%d:%d\t0%o", dir.FileName, dir.Attr.Attr.UID, dir.Attr.Attr.GID, dir.Attr.Attr.Mode)
	}

	return dirs, nil
}
