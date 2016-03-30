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
	"bytes"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"

	"crypto/tls"

	"github.com/vmware/vic/pkg/vsphere/session"
)

var (
	logFileDir = "/var/log/vic/"

	config struct {
		session.Config
		addr         string
		dockerHost   string
		vmPath       string
		hostCertFile string
		hostKeyFile  string
		tls          bool
	}

	defaultReaders []entryReader

	datastore types.ManagedObjectReference

	datastoreInventoryPath string
)

func init() {
	flag.StringVar(&config.addr, "l", ":2378", "Listen address")
	flag.StringVar(&config.dockerHost, "docker-host", "127.0.0.1:2376", "Docker host")
	flag.StringVar(&config.CertFile, "cert", "", "VMOMI Client certificate file")
	flag.StringVar(&config.hostCertFile, "hostcert", "", "Host certificate file")
	flag.StringVar(&config.KeyFile, "key", "", "VMOMI Client private key file")
	flag.StringVar(&config.hostKeyFile, "hostkey", "", "Host private key file")
	flag.StringVar(&config.Service, "sdk", "", "The ESXi or vCenter URL")
	flag.StringVar(&config.DatacenterPath, "dc", "", "Name of the Datacenter")
	flag.StringVar(&config.DatastorePath, "ds", "", "Name of the Datastore")
	flag.StringVar(&config.ClusterPath, "cluster", "", "Name of the cluster")
	flag.BoolVar(&config.Insecure, "insecure", false, "Allow connection when sdk certificate cannot be verified")
	flag.BoolVar(&config.tls, "tls", true, "Set to false to disable -hostcert and -hostkey and enable plain HTTP")

	// This is only applicable for containers hosted under the VCH VM folder
	// This will not function for vSAN
	flag.StringVar(&config.vmPath, "vm-path", "", "Docker vm path")
}

func logFiles() []string {
	names := []string{}
	files, _ := ioutil.ReadDir(logFileDir)
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, fmt.Sprintf("%s/%s", logFileDir, f.Name()))
		}
	}

	return names
}

func configureReaders() []entryReader {
	dockerServer := "http://" + config.dockerHost

	readers := []entryReader{
		fileReader("/proc/mounts"),
		urlReader(dockerServer + "/debug/pprof/goroutine?debug=1"),
		urlReader(dockerServer + "/debug/pprof/block?debug=1"),
		urlReader(dockerServer + "/debug/pprof/heap?debug=1"),
		commandReader("uptime"),
		commandReader("df"),
		commandReader("free"),
		commandReader("netstat -ant"),
		commandReader("sudo iptables --list"),
		commandReader("ip link"),
		commandReader("ip addr"),
		commandReader("ip route"),
		commandReader("lsmod"),
		// TODO: ls without shelling out
		commandReader("ls -l /dev/disk/by-path"),
		commandReader("ls -l /dev/disk/by-label"),
		// To check we are not leaking any fds
		commandReader("sudo ls -l /proc/self/fd"),
	}

	for _, path := range logFiles() {
		readers = append(readers, fileReader(path))
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

type entryReader interface {
	open() (entry, error)
}

type entry interface {
	io.ReadCloser
	Name() string
	Size() int64
}

type bytesEntry struct {
	io.ReadCloser
	name string
	size int64
}

func (e *bytesEntry) Name() string {
	return e.name
}

func (e *bytesEntry) Size() int64 {
	return e.size
}

func newBytesEntry(name string, b []byte) entry {
	r := bytes.NewReader(b)

	return &bytesEntry{
		ReadCloser: ioutil.NopCloser(r),
		size:       int64(r.Len()),
		name:       name,
	}
}

type commandReader string

func (path commandReader) open() (entry, error) {
	args := strings.Split(string(path), " ")
	cmd := exec.Command(args[0], args[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return newBytesEntry(string(path), output), nil
}

type fileReader string

type fileEntry struct {
	io.ReadCloser
	os.FileInfo
}

func (path fileReader) open() (entry, error) {
	f, err := os.Open(string(path))
	if err != nil {
		return nil, err
	}

	s, err := os.Stat(string(path))
	if err != nil {
		return nil, err
	}

	// Files in /proc always have struct stat.st_size==0, so just read it into memory.
	if s.Size() == 0 && strings.HasPrefix(f.Name(), "/proc/") {
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		if err != nil {
			return nil, err
		}

		return newBytesEntry(f.Name(), b), nil
	}

	return &fileEntry{
		ReadCloser: f,
		FileInfo:   s,
	}, nil
}

type urlReader string

func httpEntry(name string, res *http.Response) (entry, error) {
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	if res.ContentLength > 0 {
		return &bytesEntry{
			ReadCloser: res.Body,
			size:       res.ContentLength,
			name:       name,
		}, nil
	}

	// If we don't have Content-Length, read into memory for the tar.Header.Size
	body, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return newBytesEntry(name, body), nil
}

func (path urlReader) open() (entry, error) {
	res, err := http.Get(string(path))
	if err != nil {
		return nil, err
	}

	return httpEntry(string(path), res)
}

type datastoreReader struct {
	ds   *object.Datastore
	path string
}

// find datastore logs for the appliance itself and all containers
func findDatastoreLogs(c *session.Session) ([]entryReader, error) {
	ds := object.NewDatastore(c.Vim25(), datastore)
	ds.InventoryPath = datastoreInventoryPath

	var readers []entryReader
	ctx := context.Background()
	b, err := ds.Browser(ctx)
	if err != nil {
		return nil, err
	}

	spec := types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{"vmware.log", "*.debug"},
	}

	task, err := b.SearchDatastoreSubFolders(ctx, ds.Path(config.vmPath), &spec)
	if err != nil {
		return nil, err
	}

	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, errors.New("empty search results")
	}

	if info.Error != nil {
		return nil, errors.New(info.Error.LocalizedMessage)
	}

	res, ok := info.Result.(types.ArrayOfHostDatastoreBrowserSearchResults)
	if !ok {
		return nil, fmt.Errorf("unexpected search result type: %T", info.Result)
	}

	for _, r := range res.HostDatastoreBrowserSearchResults {
		folder := r.FolderPath
		if ix := strings.Index(folder, "] "); ix != -1 {
			folder = folder[ix+2:]
		}

		for _, f := range r.File {
			name := path.Join(folder, f.GetFileInfo().Path)
			readers = append(readers, &datastoreReader{ds, name})
		}
	}

	return readers, nil
}

func (r datastoreReader) open() (entry, error) {
	u, ticket, err := r.ds.ServiceTicket(context.Background(), r.path, "GET")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	req.AddCookie(ticket)

	res, err := r.ds.Client().Do(req)
	if err != nil {
		return nil, err
	}

	return httpEntry(r.path, res)
}

type dlogReader struct {
	c    *session.Session
	name string
	host *object.HostSystem
}

func findDiagnosticLogs(c *session.Session) ([]entryReader, error) {
	// When connected to VC, we collect vpxd.log and hostd.log for all cluster hosts attached to the datastore.
	// When connected to ESX, we just collect hostd.log.
	const (
		vpxdKey  = "vpxd:vpxd.log"
		hostdKey = "hostd"
	)

	var logs []entryReader
	var err error

	if c.IsVC() {
		logs = append(logs, dlogReader{c, vpxdKey, nil})

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
			logs = append(logs, dlogReader{c, hostdKey, host})
		}
	} else {
		logs = append(logs, dlogReader{c, hostdKey, nil})
	}

	return logs, nil
}

func (r dlogReader) open() (entry, error) {
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

func client() (*session.Session, error) {
	ctx := context.Background()

	session := session.NewSession(&config.Config)
	_, err := session.Connect(ctx)
	if err != nil {
		return nil, err
	}

	_, err = session.Populate(ctx)
	if err != nil {
		// no a critical error for vicadmin
		log.Warn(err)
	}

	return session, nil
}

func findDatastore() error {
	session, err := client()
	if err != nil {
		return err
	}
	defer session.Client.Logout(context.Background())

	datastore = session.Datastore.Reference()
	datastoreInventoryPath = session.Datastore.InventoryPath

	return nil
}

func tarEntries(readers []entryReader, out io.Writer) error {
	r, w := io.Pipe()
	t := tar.NewWriter(w)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	// stream tar to out
	go func() {
		_, err := io.Copy(out, r)
		if err != nil {
			log.Printf("error copying tar: %s", err)
		}
		wg.Done()
	}()

	for _, r := range readers {
		e, err := r.open()
		if err != nil {
			log.Warningf("error reading %s: %s\n", r, err)
			continue
		}

		header := tar.Header{
			Name:    url.QueryEscape(e.Name()),
			Size:    e.Size(),
			Mode:    0640,
			ModTime: time.Now(),
		}

		err = t.WriteHeader(&header)
		if err != nil {
			return err
		}

		_, err = io.Copy(t, e)
		_ = e.Close()
		if err != nil {
			return err
		}
	}

	_ = t.Flush()
	_ = w.Close()
	wg.Wait()
	_ = r.Close()

	return nil
}

type server struct {
	l     net.Listener
	addr  string
	mux   *http.ServeMux
	links []string
}

func (s *server) listen(useTLS bool) error {
	var err error
	if !useTLS {
		s.l, err = net.Listen("tcp", s.addr)
		return err
	}

	innerListener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	certificate, err := tls.LoadX509KeyPair(config.hostCertFile, config.hostKeyFile)
	if err != nil {
		log.Fatalf("Could not load either key file %s or certificate file %s",
			config.hostKeyFile, config.hostCertFile)
		return err
	}

	tlsconfig := tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	s.l = tls.NewListener(innerListener, &tlsconfig)
	return nil
}

func (s *server) listenPort() int {
	return s.l.Addr().(*net.TCPAddr).Port
}

func (s *server) handleFunc(link string, handler func(http.ResponseWriter, *http.Request)) {
	s.links = append(s.links, link)
	s.mux.HandleFunc(link, handler)
}

func (s *server) serve() error {
	s.mux = http.NewServeMux()

	// tar of appliance system logs
	s.handleFunc("/logs.tar.gz", s.tarDefaultLogs)

	// tar of appliance system logs + container logs
	s.handleFunc("/container-logs.tar.gz", s.tarContainerLogs)

	// tail all logFiles
	s.mux.HandleFunc("/logs/tail", func(w http.ResponseWriter, r *http.Request) {
		s.tailFiles(w, r, logFiles())
	})

	for _, path := range logFiles() {
		name := filepath.Base(path)
		p := path

		// get single log file (no tail)
		s.handleFunc("/logs/"+name, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, p)
		})

		// get single log file (with tail)
		s.mux.HandleFunc("/logs/tail/"+name, func(w http.ResponseWriter, r *http.Request) {
			s.tailFiles(w, r, []string{p})
		})
	}

	s.mux.HandleFunc("/", s.index)

	server := &http.Server{
		Handler: s.mux,
	}

	defaultReaders = configureReaders()

	return server.Serve(s.l)
}

func (s *server) stop() error {
	if s.l != nil {
		err := s.l.Close()
		s.l = nil
		return err
	}

	return nil
}

func (s *server) tarContainerLogs(res http.ResponseWriter, req *http.Request) {
	readers := append(defaultReaders, commandReader("sudo du -sh /var/lib/docker"))

	if config.Service != "" {
		c, err := client()
		if err != nil {
			log.Errorf("failed to connect: %s", err)
		} else {
			// Note: we don't want to Logout() until tarEntries() completes below
			defer c.Client.Logout(context.Background())

			logs, err := findDatastoreLogs(c)
			if err != nil {
				log.Warningf("error searching datastore: %s", err)
			} else {
				readers = append(readers, logs...)
			}

			logs, err = findDiagnosticLogs(c)
			if err != nil {
				log.Warningf("error collecting diagnostic logs: %s", err)
			} else {
				readers = append(readers, logs...)
			}
		}
	}

	s.tarLogs(res, req, readers)
}

func (s *server) tarDefaultLogs(res http.ResponseWriter, req *http.Request) {
	s.tarLogs(res, req, defaultReaders)
}

func (s *server) tarLogs(res http.ResponseWriter, req *http.Request, readers []entryReader) {
	res.Header().Set("Content-Type", "application/x-gzip")

	z := gzip.NewWriter(res)

	err := tarEntries(readers, z)
	if err != nil {
		log.Printf("error taring logs: %s", err)
	}

	_ = z.Close()
}

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (int, error) {
	n, err := fw.w.Write(p)

	fw.f.Flush()

	return n, err
}

func (s *server) tailFiles(res http.ResponseWriter, req *http.Request, names []string) {
	cc := res.(http.CloseNotifier).CloseNotify()

	fw := &flushWriter{
		f: res.(http.Flusher),
		w: res,
	}

	// TODO: tail in Go, rather than shelling out
	cmd := exec.Command("tail", append([]string{"-F"}, names...)...)
	cmd.Stdout = fw
	cmd.Stderr = fw

	if err := cmd.Start(); err != nil {
		log.Printf("error tailing logs: %s", err)
		return
	}

	defer cmd.Process.Kill()

	go cmd.Wait()

	<-cc
}

func (s *server) index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "<html><head><title>VIC Admin</title></head><body><pre>")
	for _, link := range s.links {
		fmt.Fprintf(res, "<a href=\"%s\">%s</a><br/>\n", link, link)
	}
	fmt.Fprintln(res, "</pre></body></html>")
}

func main() {
	flag.Parse()

	s := &server{
		addr: config.addr,
	}

	err := s.listen(config.tls)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s", s.addr)

	signals := []syscall.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
	}

	sigchan := make(chan os.Signal, 1)
	for _, signum := range signals {
		signal.Notify(sigchan, signum)
	}

	go func() {
		signal := <-sigchan
		log.Infof("received %s", signal)
		s.stop()
	}()

	s.serve()
}
