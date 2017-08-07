// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"archive/tar"
	"context"
	"io"
	"net/url"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"

	dar "github.com/docker/docker/pkg/archive"

	"github.com/vmware/govmomi/toolbox/hgfs"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/portlayer/storage/vsphere"
	"github.com/vmware/vic/lib/tether"
	viclog "github.com/vmware/vic/pkg/log"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

var (
	tthr                  tether.Tether
	defaultArchiveHandler = hgfs.NewArchiveHandler().(*hgfs.ArchiveHandler)
)

func init() {
	// Initialize logger with default TextFormatter
	log.SetFormatter(viclog.NewTextFormatter())

	// use the same logger for trace and other logging
	trace.Logger.Level = log.DebugLevel
	log.SetLevel(log.DebugLevel)

	// init and start the HUP handler
	startSignalHandler()

	pathPrefix = "/dev"
}

func main() {
	if strings.HasSuffix(os.Args[0], "-debug") {
		// very, very verbose
		extraconfig.SetLogLevel(log.DebugLevel)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("run time panic: %s : %s", r, debug.Stack())
		}
		halt()
	}()

	logFile, err := os.OpenFile(path.Join(pathPrefix, "ttyS1"), os.O_WRONLY|os.O_SYNC, 0)
	if err != nil {
		log.Errorf("Could not open serial port for debugging info. Some debug info may be lost! Error reported was %s", err)
	}

	if err = syscall.Dup3(int(logFile.Fd()), int(os.Stderr.Fd()), 0); err != nil {
		log.Errorf("Could not pipe logfile to standard error due to error %s", err)
	}

	if _, err = os.Stderr.WriteString("all stderr redirected to debug log"); err != nil {
		log.Errorf("Could not write to Stderr due to error %s", err)
	}

	// where to look for the various devices and files related to tether

	// TODO: hard code executor initialization status reporting via guestinfo here
	sshserver := NewAttachServerSSH()
	src, err := extraconfig.GuestInfoSource()
	if err != nil {
		log.Error(err)
		return
	}

	sink, err := extraconfig.GuestInfoSink()
	if err != nil {
		log.Error(err)
		return
	}

	// create the tether
	tthr = tether.New(src, sink, &operations{})

	// register the attach extension
	tthr.Register("Attach", sshserver)

	// register the toolbox extension
	toolbox := tether.NewToolbox().InContainer()
	cmd := toolbox.Service.Command
	cmd.FileServer.RegisterFileHandler(hgfs.ArchiveScheme, &hgfs.ArchiveHandler{
		Read:  toolboxOverrideArchiveRead,
		Write: toolboxOverrideArchiveWrite,
	})
	tthr.Register("Toolbox", toolbox)

	err = tthr.Start()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Clean exit from tether")
}

// exit cleanly shuts down the system
func halt() {
	log.Infof("Powering off the system")
	if strings.HasSuffix(os.Args[0], "-debug") {
		log.Info("Squashing power off for debug tether")
		return
	}

	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func startSignalHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for s := range sigs {
			switch s {
			case syscall.SIGHUP:
				log.Infof("Reloading tether configuration")
				tthr.Reload()
			default:
				log.Infof("%s signal not defined", s.String())
			}
		}
	}()
}

// toolboxOverrideArchiveRead is the online DataSink Override Handler
func toolboxOverrideArchiveRead(u *url.URL, tr *tar.Reader) error {

	// special behavior when using disk-labels and filterspec
	diskLabel := u.Query().Get(vsphere.DiskLabelQueryName)
	filterSpec := u.Query().Get(vsphere.FilterSpecQueryName)
	if diskLabel != "" && filterSpec != "" {
		op := trace.NewOperation(context.Background(), "ToolboxOnlineDataSink: %s", u.String())
		op.Debugf("Reading from tar archive to path %s: %s", u.Path, u.String())
		spec, err := archive.DecodeFilterSpec(op, &filterSpec)
		if err != nil {
			op.Errorf(err.Error())
			return err
		}
		diskPath, err := MountDiskLabel(diskLabel)
		if err != nil {
			op.Errorf(err.Error())
			return err
		}
		defer UnmountDiskLabel(op, diskPath)

		// no need to join on u.Path here. u.Path == spec.Rebase, but
		// Unpack will rebase tar headers for us. :thumbsup:
		err = archive.Unpack(op, tr, spec, diskPath)
		if err != nil {
			op.Errorf(err.Error())
		}
		op.Debugf("Finished reading from tar archive to path %s: %s", u.Path, u.String())
		return err
	}
	return defaultArchiveHandler.Read(u, tr)

}

// toolboxOverrideArchiveWrite is the Online DataSource Override Handler
func toolboxOverrideArchiveWrite(u *url.URL, tw *tar.Writer) error {

	// special behavior when using disk-labels and filterspec
	diskLabel := u.Query().Get(vsphere.DiskLabelQueryName)
	filterSpec := u.Query().Get(vsphere.FilterSpecQueryName)

	skiprecurse, _ := strconv.ParseBool(u.Query().Get(vsphere.SkipRecurseQueryName))
	skipdata, _ := strconv.ParseBool(u.Query().Get(vsphere.SkipDataQueryName))

	if diskLabel != "" && filterSpec != "" {
		op := trace.NewOperation(context.Background(), "ToolboxOnlineDataSource: %s", u.String())
		op.Debugf("Writing to archive from %s: %s", u.Path, u.String())

		spec, err := archive.DecodeFilterSpec(op, &filterSpec)
		if err != nil {
			op.Errorf(err.Error())
			return err
		}

		// get the container fs mount
		diskPath, err := MountDiskLabel(diskLabel)
		if err != nil {
			op.Errorf(err.Error())
			return err
		}
		defer UnmountDiskLabel(op, diskPath)

		var rc io.ReadCloser
		if skiprecurse {
			// we only want a single file - this is a hack while we're abusing Diff, but
			// accomplish this by generating a single entry ChangeSet
			changes := []dar.Change{
				{
					Kind: dar.ChangeModify,
					Path: u.Path,
				},
			}

			rc, err = archive.Tar(op, diskPath, changes, spec, !skipdata, false)
		} else {
			rc, err = archive.Diff(op, diskPath, "", spec, !skipdata, false)
		}

		if err != nil {
			op.Errorf(err.Error())
			return err
		}

		tr := tar.NewReader(rc)
		defer rc.Close()
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				op.Debugf("Finished writing to archive from %s: %s with error %#v", u.Path, u.String(), err)
				break
			}
			if err != nil {
				op.Errorf("error writing tar: %s", err.Error())
				return err
			}
			op.Debugf("Writing header: %#s", *hdr)
			err = tw.WriteHeader(hdr)
			if err != nil {
				op.Errorf("error writing tar header: %s", err.Error())
				return err
			}
			_, err = io.Copy(tw, tr)
			if err != nil {
				op.Errorf("error writing tar contents: %s", err.Error())
				return err
			}
		}

		return nil
	}
	return defaultArchiveHandler.Write(u, tw)
}
