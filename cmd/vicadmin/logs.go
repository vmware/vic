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

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/pprof"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

const (
	uint32max = (1 << 32) - 1
	tailLines = 8
)

type dlogReader struct {
	c    *session.Session
	name string
	host *object.HostSystem
}

func (r dlogReader) open() (entry, error) {
	defer trace.End(trace.Begin(r.name))

	name := r.name
	if r.host != nil {
		name = fmt.Sprintf("%s-%s", path.Base(r.host.InventoryPath), r.name)
	}

	m := object.NewDiagnosticManager(r.c.Vim25())
	ctx := context.Background()

	// Currently we collect the tail of diagnostic log files to avoid
	// reading the entire file into memory or writing local disk.

	// get LineEnd without any LineText
	h, err := m.BrowseLog(ctx, r.host, r.name, math.MaxInt32, 0)
	if err != nil {
		return nil, err
	}

	// DiagnosticManager::DEFAULT_MAX_LINES_PER_BROWSE = 1000
	start := h.LineEnd - 1000

	h, err = m.BrowseLog(ctx, r.host, r.name, start, 0)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	for _, line := range h.LineText {
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	return newBytesEntry(name+".log", buf.Bytes()), nil
}

func logFiles() []string {
	defer trace.End(trace.Begin(""))

	names := []string{}
	for _, f := range logFileList {
		names = append(names, fmt.Sprintf("%s/%s", logFileDir, f))
	}

	return names
}

func configureReaders() map[string]entryReader {
	defer trace.End(trace.Begin(""))

	pprofPaths := map[string]string{
		// verbose
		"verbose": "/debug/pprof/goroutine?debug=2",
		// concise
		"concise": "/debug/pprof/goroutine?debug=1",
		"block":   "/debug/pprof/block?debug=1",
		"heap":    "/debug/pprof/heap?debug=1",
	}

	pprofSources := map[string]string{
		"docker":    pprof.GetPprofEndpoint(pprof.DockerPort).String(),
		"portlayer": pprof.GetPprofEndpoint(pprof.PortlayerPort).String(),
		"vicadm":    pprof.GetPprofEndpoint(pprof.VicadminPort).String(),
		"vic-init":  pprof.GetPprofEndpoint(pprof.VCHInitPort).String(),
	}

	readers := map[string]entryReader{
		"proc-mounts": fileReader("/proc/mounts"),
		"uptime":      commandReader("uptime"),
		"df":          commandReader("df"),
		"free":        commandReader("free"),
		"netstat":     commandReader("netstat -ant"),
		"iptables":    commandReader("sudo iptables --list"),
		"ip-link":     commandReader("ip link"),
		"ip-addr":     commandReader("ip addr"),
		"ip-route":    commandReader("ip route"),
		"lsmod":       commandReader("lsmod"),
		// TODO: ls without shelling out
		"disk-by-path":  commandReader("ls -l /dev/disk/by-path"),
		"disk-by-label": commandReader("ls -l /dev/disk/by-label"),
		// To check we are not leaking any fds
		"proc-self-fd": commandReader("ls -l /proc/self/fd"),
	}

	// add the pprof collection
	for sname, source := range pprofSources {
		for pname, paths := range pprofPaths {
			rname := fmt.Sprintf("%s/%s", sname, pname)
			readers[rname] = urlReader(source + paths)
		}
	}

	for _, path := range logFiles() {
		// Strip off leading '/'
		readers[path[1:]] = fileReader(path)
	}

	if config.vmPath == "" {
		log.Info("vm-path not set, skipping datastore log collection")
	} else {
		err := findDatastore()

		if err != nil {
			log.Warning(err)
		}
	}

	return readers
}

func findDiagnosticLogs(c *session.Session) (map[string]entryReader, error) {
	defer trace.End(trace.Begin(""))

	// When connected to VC, we collect vpxd.log and hostd.log for all cluster hosts attached to the datastore.
	// When connected to ESX, we just collect hostd.log.
	const (
		vpxdKey  = "vpxd:vpxd.log"
		hostdKey = "hostd"
	)

	logs := map[string]entryReader{}
	var err error

	if c.IsVC() {
		logs[vpxdKey] = dlogReader{c, vpxdKey, nil}

		var hosts []*object.HostSystem
		if c.Cluster == nil && c.Host != nil {
			hosts = []*object.HostSystem{c.Host}
		} else {
			hosts, err = c.Datastore.AttachedClusterHosts(context.TODO(), c.Cluster)
			if err != nil {
				return nil, err
			}
		}

		for _, host := range hosts {
			lname := fmt.Sprintf("%s/%s", hostdKey, host)
			logs[lname] = dlogReader{c, hostdKey, host}
		}
	} else {
		logs[hostdKey] = dlogReader{c, hostdKey, nil}
	}

	return logs, nil
}
func tarEntries(readers map[string]entryReader, out io.Writer) error {
	defer trace.End(trace.Begin(""))

	r, w := io.Pipe()
	t := tar.NewWriter(w)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	// stream tar to out
	go func() {
		_, err := io.Copy(out, r)
		if err != nil {
			log.Errorf("error copying tar: %s", err)
		}
		wg.Done()
	}()

	for name, r := range readers {
		log.Infof("Collecting log with reader %s(%#v)", name, r)

		e, err := r.open()
		if err != nil {
			log.Warningf("error reading %s(%s): %s\n", name, r, err)
			continue
		}

		header := tar.Header{
			Name:    name,
			Size:    e.Size(),
			Mode:    0640,
			ModTime: time.Now(),
		}

		err = t.WriteHeader(&header)
		if err != nil {
			log.Errorf("Failed to write header for %s: %s", header.Name, err)
			continue
		}

		log.Infof("%s has size %d", header.Name, header.Size)

		// be explicit about the number of bytes to copy as the log files will likely
		// be written to during this exercise
		_, err = io.CopyN(t, e, e.Size())
		_ = e.Close()
		if err != nil {
			log.Errorf("Failed to write content for %s: %s", header.Name, err)
			continue
		}
	}

	_ = t.Flush()
	_ = w.Close()
	wg.Wait()
	_ = r.Close()

	return nil
}
func zipEntries(readers map[string]entryReader, out *zip.Writer) error {
	defer trace.End(trace.Begin(""))
	defer out.Close()
	defer out.Flush()

	for name, r := range readers {
		log.Infof("Collecting log with reader %s(%#v)", name, r)

		e, err := r.open()
		if err != nil {
			log.Warningf("error reading %s(%s): %s\n", name, r, err)
			continue
		}
		sz := e.Size()
		header := &zip.FileHeader{
			Name:   name,
			Method: zip.Deflate,
		}

		header.SetModTime(time.Now())
		header.SetMode(0644)
		if sz > uint32max {
			header.UncompressedSize = uint32max
		} else {
			header.UncompressedSize = uint32(e.Size())
		}

		w, err := out.CreateHeader(header)

		if err != nil {
			log.Errorf("Failed to create Zip writer for %s: %s", header.Name, err)
			continue
		}

		log.Infof("%s has size %d", header.Name, sz)

		// be explicit about the number of bytes to copy as the log files will likely
		// be written to during this exercise
		_, err = io.CopyN(w, e, sz)
		_ = e.Close()
		if err != nil {
			log.Errorf("Failed to write content for %s: %s", header.Name, err)
			continue
		}
		log.Infof("Wrote %d bytes to %s", sz, header.Name)
	}
	return nil
}

func tailFile(wr io.Writer, file string, done *chan bool) error {
	nlines := tailLines
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	pos, err := f.Seek(0, 2)
	if err != nil {
		return err
	}
	b := make([]byte, 1)
	for pos > 0 && nlines >= 0 {
		pos--
		_, err := f.ReadAt(b, pos)
		if err != nil && err != io.EOF {
			return err
		}
		if b[0] == '\n' {
			nlines--
		}
	}

	pos++
	if _, err = f.Seek(pos, 0); err != nil {
		return err
	}

	b = make([]byte, 1024)
	// We KNOW there's a data race here.
	// But it doesn't break anything, so we just trap it.
	defer func() {
		_ = recover()
	}()
	var n int
	for true {
		for err != io.EOF {
			n, err = f.Read(b)
			if err != nil && err != io.EOF {
				return err
			}
			fmt.Fprint(wr, string(b[:n]))
		}
		select {
		case _ = <-*done:
			return nil
		default:
		}
		time.Sleep(50 * time.Millisecond)
		err = nil
	}
	return nil
}
