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

// +build !windows,!darwin

package tether

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/vmware/govmomi/toolbox"
	"github.com/vmware/govmomi/toolbox/hgfs"
	"github.com/vmware/govmomi/toolbox/vix"
	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/portlayer/storage/vsphere"
	"github.com/vmware/vic/pkg/trace"
)

// Toolbox is a tether extension that wraps toolbox.Service
type Toolbox struct {
	*toolbox.Service

	sess struct {
		sync.Mutex
		session *SessionConfig
	}

	// IDs that can be used to authenticate
	authIDs map[string]struct{}

	stop chan struct{}
}

var (
	defaultArchiveHandler hgfs.ArchiveHandler
)

func init() {
	fileHandler := hgfs.NewArchiveHandler()
	if afh, ok := fileHandler.(*hgfs.ArchiveHandler); ok {
		defaultArchiveHandler = *afh
	}
}

// NewToolbox returns a tether.Extension that wraps the vsphere/toolbox service
func NewToolbox() *Toolbox {
	in := toolbox.NewBackdoorChannelIn()
	out := toolbox.NewBackdoorChannelOut()

	service := toolbox.NewService(in, out)
	service.PrimaryIP = toolbox.DefaultIP

	// We need to see the output for archive read and write. We need to gate these on the
	// debug variable on vch creation, same as the personality and persona.
	hgfs.Trace = true
	toolbox.Trace = true

	return &Toolbox{
		Service: service,
		authIDs: make(map[string]struct{}),
	}
}

// Start implementation of the tether.Extension interface
func (t *Toolbox) Start() error {
	t.stop = make(chan struct{})
	on := make(chan struct{})

	t.Service.Power.PowerOn.Handler = func() error {
		log.Info("toolbox: service is ready (power on event received)")
		close(on)
		return nil
	}

	err := t.Service.Start()
	if err != nil {
		return err
	}

	// Wait for the vmx to send the OS_PowerOn message,
	// at which point it will be ready to service vix command requests.
	log.Info("toolbox: waiting for initialization")

	select {
	case <-on:
	case <-time.After(time.Second):
		log.Warn("toolbox: timeout waiting for power on event")
	}

	return nil
}

// Stop implementation of the tether.Extension interface
func (t *Toolbox) Stop() error {
	t.Service.Stop()

	t.Service.Wait()

	close(t.stop)

	return nil
}

// Reload implementation of the tether.Extension interface
func (t *Toolbox) Reload(config *ExecutorConfig) error {
	if config != nil && config.Sessions != nil {
		t.sess.Lock()
		defer t.sess.Unlock()
		t.sess.session = config.Sessions[config.ID]
	}

	// we allow the primary session
	t.authIDs[config.ID] = struct{}{}
	// we also allow any device IDs that are attached
	for _, mspec := range config.Mounts {
		t.authIDs[mspec.Source.Host] = struct{}{}
	}

	return nil
}

// InContainer configures the toolbox to run within a container VM
func (t *Toolbox) InContainer() *Toolbox {
	t.Power.Halt.Handler = t.halt

	cmd := t.Service.Command
	cmd.Authenticate = t.containerAuthenticate
	cmd.ProcessStartCommand = t.containerStartCommand

	cmd.FileServer.RegisterFileHandler(hgfs.ArchiveScheme, &hgfs.ArchiveHandler{
		Read:  toolboxOverrideArchiveRead,
		Write: toolboxOverrideArchiveWrite,
	})

	return t
}

func (t *Toolbox) session() *SessionConfig {
	t.sess.Lock()
	defer t.sess.Unlock()
	return t.sess.session
}

func (t *Toolbox) kill(name string) error {
	session := t.session()
	if session == nil {
		return fmt.Errorf("failed to kill container: process not found")
	}

	session.Lock()
	defer session.Unlock()
	return t.killHelper(session, name)
}

func (t *Toolbox) killHelper(session *SessionConfig, name string) error {
	if name == "" {
		name = string(ssh.SIGTERM)
	}

	if session.Cmd.Process == nil {
		return fmt.Errorf("the session %s hasn't launched yet", session.ID)
	}

	sig := new(msgs.SignalMsg)
	err := sig.FromString(name)
	if err != nil {
		return err
	}

	num := syscall.Signal(sig.Signum())

	log.Infof("sending signal %s (%d) to %s", sig.Signal, num, session.ID)

	if err := session.Cmd.Process.Signal(num); err != nil {
		return fmt.Errorf("failed to signal %s: %s", session.ID, err)
	}

	return nil
}

func (t *Toolbox) containerAuthenticate(_ vix.CommandRequestHeader, data []byte) error {
	var c vix.UserCredentialNamePassword
	if err := c.UnmarshalBinary(data); err != nil {
		return err
	}

	session := t.session()
	if session == nil {
		return errors.New("not yet initialized")
	}

	session.Lock()
	defer session.Unlock()

	// no authentication yet, just using container ID and device IDs as a sanity check for now
	if _, ok := t.authIDs[c.Name]; !ok {
		return errors.New("failed to verify authentication ID")
	}

	return nil
}

func (t *Toolbox) containerStartCommand(m *toolbox.ProcessManager, r *vix.StartProgramRequest) (int64, error) {
	switch r.ProgramPath {
	case "kill":
		return -1, t.kill(r.Arguments)
	case "reload":
		return -1, ReloadConfig()
	default:
		return -1, fmt.Errorf("unknown command %q", r.ProgramPath)
	}
}

func (t *Toolbox) halt() error {
	session := t.session()
	if session == nil {
		return fmt.Errorf("failed to halt container: not initialized yet")
	}

	session.Lock()
	defer session.Unlock()

	if session.Cmd.Process == nil {
		return fmt.Errorf("the session %s hasn't launched yet", session.ID)
	}

	log.Infof("stopping %s", session.ID)

	if err := t.killHelper(session, session.StopSignal); err != nil {
		return err
	}

	// Killing the executor session in the container VM will stop the tether and its extensions.
	// If that doesn't happen within the timeout, send a SIGKILL.
	select {
	case <-t.stop:
		log.Infof("%s has stopped", session.ID)
		return nil
	case <-time.After(time.Second * 10):
	}

	log.Warnf("killing %s", session.ID)

	return session.Cmd.Process.Kill()
}

// ----------
// Online DataSource and DataSink Override Handlers
// ----------

func toolboxOverrideArchiveRead(u *url.URL, tr *tar.Reader) error {

	// special behavior when using disk-labels and filterspec
	diskLabel := u.Query().Get(vsphere.DiskLabelQueryName)
	filterSpec := u.Query().Get(vsphere.FilterSpecQueryName)
	if diskLabel != "" && filterSpec != "" {
		op := trace.NewOperation(context.Background(), "ToolboxOnlineDataSink: %s", u.String())
		op.Debugf("Reading from tar archive to path %s: %s", u.Path, u.String())
		spec, err := archive.DecodeFilterSpec(op, &filterSpec)
		if err != nil {
			op.Debugf(err.Error())
			return err
		}
		diskPath, err := mountDiskLabel(diskLabel)
		if err != nil {
			op.Debugf(err.Error())
			return err
		}
		defer unmount(op, diskPath)

		// no need to join on u.Path here. u.Path == spec.Rebase, but
		// Unpack will rebase tar headers for us. :thumbsup:
		err = archive.Unpack(op, tr, spec, diskPath)
		if err != nil {
			op.Debugf(err.Error())
		}
		op.Debugf("Finished reading from tar archive to path %s: %s", u.Path, u.String())
		return err
	}
	return defaultArchiveHandler.Read(u, tr)

}

func toolboxOverrideArchiveWrite(u *url.URL, tw *tar.Writer) error {

	// special behavior when using disk-labels and filterspec
	diskLabel := u.Query().Get(vsphere.DiskLabelQueryName)
	filterSpec := u.Query().Get(vsphere.FilterSpecQueryName)
	if diskLabel != "" && filterSpec != "" {
		op := trace.NewOperation(context.Background(), "ToolboxOnlineDataSource: %s", u.String())
		op.Debugf("Writing to archive from %s: %s", u.Path, u.String())

		spec, err := archive.DecodeFilterSpec(op, &filterSpec)
		if err != nil {
			op.Debugf(err.Error())
			return err
		}

		// get the container fs mount
		diskPath, err := mountDiskLabel(diskLabel)
		if err != nil {
			op.Debugf(err.Error())
			return err
		}
		defer unmount(op, diskPath)

		rc, err := archive.Diff(op, diskPath, "", spec, true, false)
		if err != nil {
			op.Debugf(err.Error())
			return err
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer rc.Close()
			defer wg.Done()
			tr := tar.NewReader(rc)
			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					err = nil
					op.Debugf("Finished writing to archive from %s: %s with error %#v", u.Path, u.String(), err)
					break
				}
				if err != nil {
					op.Debugf("error writing tar: %s", err.Error())
					break
				}
				op.Debugf("Writing header: %#s", *hdr)
				err = tw.WriteHeader(hdr)
				if err != nil {
					op.Debugf("error writing tar header: %s", err.Error())
					break
				}
				_, err = io.Copy(tw, tr)
				if err != nil {
					op.Debugf("error writing tar contents: %s", err.Error())
					break
				}
			}
		}()
		wg.Wait()
		return err
	}
	return defaultArchiveHandler.Write(u, tw)
}

func mountDiskLabel(label string) (string, error) {
	// We don't really need to bind mount the vmdk, we know it will
	// always be attached at '/'
	if label == "containerfs" {
		return "/", nil
	}

	// otherwise, label represents a volume that needs to be bind mounted
	path := path.Join("/dev/disk/by-label/", label)
	_, err := os.Lstat(path)
	// label does not exist
	if err != nil {
		return "", err
	}

	// label exists, bind mount it to a tmp directory
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("toolbox-%s", label))
	if err := Sys.Syscall.Mount(path, tmpDir, "bind", syscall.MS_BIND, ""); err != nil {
		return "", fmt.Errorf("faild to mount %s to %s: %s", path, tmpDir, err)
	}

	return tmpDir, nil
}

func unmount(op trace.Operation, unmountPath string) {
	// don't unmount the root vmdk
	if unmountPath == "/" {
		return
	}

	// unmount the disk from the temporary directory
	if err := Sys.Syscall.Unmount(unmountPath, syscall.MNT_DETACH); err != nil {
		op.Debugf("failed to unmount %s: %s", unmountPath, err.Error())
	}
	// finally, remove the temporary directory
	os.RemoveAll(unmountPath)
}
