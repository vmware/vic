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
	"crypto/tls"
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
	"github.com/docker/docker/pkg/tlsconfig"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/lib/pprof"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	timeout = time.Duration(2 * time.Second)
)

var (
	logFileDir  = "/var/log/vic"
	logFileList = []string{
		"docker-personality.log",
		"imagec.log",
		"port-layer.log",
		"vicadmin.log",
		"init.log",
	}

	// VMFiles is the set of files to collect per VM associated with the VCH
	vmFiles = []string{
		"vmware.log",
		string(metadata.VM + ".debug"),
	}

	config struct {
		session.Config
		addr         string
		dockerHost   string
		vmPath       string
		hostCertFile string
		hostKeyFile  string
		authType     string
		timeout      time.Time
		tls          bool
	}

	vchConfig metadata.VirtualContainerHostConfigSpec

	defaultReaders map[string]entryReader

	datastore types.ManagedObjectReference

	datastoreInventoryPath string
)

func init() {
	trace.Logger.Level = log.DebugLevel
	defer trace.End(trace.Begin(""))

	flag.StringVar(&config.addr, "l", ":2378", "Listen address")
	flag.StringVar(&config.dockerHost, "docker-host", "127.0.0.1:2376", "Docker host")
	flag.StringVar(&config.ExtensionCert, "cert", "", "VMOMI Client certificate file")
	flag.StringVar(&config.hostCertFile, "hostcert", "", "Host certificate file")
	flag.StringVar(&config.ExtensionKey, "key", "", "VMOMI Client private key file")
	flag.StringVar(&config.hostKeyFile, "hostkey", "", "Host private key file")
	flag.StringVar(&config.Service, "sdk", "", "The ESXi or vCenter URL")
	flag.StringVar(&config.DatacenterPath, "dc", "", "Name of the Datacenter")
	flag.StringVar(&config.DatastorePath, "ds", "", "Name of the Datastore")
	flag.StringVar(&config.ClusterPath, "cluster", "", "Path of the cluster")
	flag.StringVar(&config.PoolPath, "pool", "", "Path of the resource pool")
	flag.BoolVar(&config.Insecure, "insecure", false, "Allow connection when sdk certificate cannot be verified")
	flag.BoolVar(&config.tls, "tls", true, "Set to false to disable -hostcert and -hostkey and enable plain HTTP")

	// This is only applicable for containers hosted under the VCH VM folder
	// This will not function for vSAN
	flag.StringVar(&config.vmPath, "vm-path", "", "Docker vm path")

	// load the vch config
	src, err := extraconfig.GuestInfoSource()
	if err != nil {
		log.Errorf("Unable to load configuration from guestinfo")
		return
	}

	extraconfig.Decode(src, &vchConfig)
	pprof.StartPprof("vicadmin", pprof.VicadminPort)
}

type Authenticator interface {
	// Validate will validate a user and password combo and return a bool.
	Validate(string, string) bool
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
		"proc-mounts":  fileReader("/proc/mounts"),
		"uptime":       commandReader("uptime"),
		"df":           commandReader("df"),
		"free":         commandReader("free"),
		"netstat":      commandReader("netstat -ant"),
		"iptables":     commandReader("sudo iptables -L -v"),
		"iptables-nat": commandReader("sudo iptables -L -v -t nat"),
		"ip-link":      commandReader("ip link"),
		"ip-addr":      commandReader("ip addr"),
		"ip-route":     commandReader("ip route"),
		"lsmod":        commandReader("lsmod"),
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
	defer trace.End(trace.Begin(string(path)))

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
	defer trace.End(trace.Begin(string(path)))

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
	defer trace.End(trace.Begin(name))

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
	defer trace.End(trace.Begin(string(path)))
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Get(string(path))
	if err != nil {
		return nil, err
	}

	return httpEntry(string(path), res)
}

type datastoreReader struct {
	ds   *object.Datastore
	path string
}

// listVMPaths returns an array of datastore paths for VMs assocaited with the
// VCH - this includes containerVMs and the appliance
func listVMPaths(ctx context.Context, s *session.Session) ([]url.URL, error) {
	defer trace.End(trace.Begin(""))

	var err error
	var children []*vm.VirtualMachine

	if len(vchConfig.ComputeResources) == 0 {
		return nil, errors.New("compute resources is empty")
	}

	ref := vchConfig.ComputeResources[0]
	rp := compute.NewResourcePool(ctx, s, ref)
	if children, err = rp.GetChildrenVMs(ctx, s); err != nil {
		return nil, err
	}

	log.Infof("Found %d candidate VMs in resource pool %s for log collection", len(children), ref.String())

	directories := []url.URL{}
	for _, child := range children {
		path, err := child.DSPath(ctx)
		if err != nil {
			log.Errorf("Unable to get datastore path for child VM %s: %s", child.Reference(), err)
			// we need to get as many logs as possible
			continue
		}

		log.Debugf("Adding VM for log collection: %s", path.String())
		directories = append(directories, path)
	}

	log.Infof("Collecting logs from %d VMs", len(directories))
	return directories, nil
}

// find datastore logs for the appliance itself and all containers
func findDatastoreLogs(c *session.Session) (map[string]entryReader, error) {
	defer trace.End(trace.Begin(""))

	// Create an empty reader as opposed to a nil reader...
	readers := map[string]entryReader{}
	ctx := context.Background()

	paths, err := listVMPaths(ctx, c)
	if err != nil {
		detail := fmt.Sprintf("unable to perform datastore log collection due to failure looking up paths: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	for _, vmpath := range paths {
		log.Debugf("Assembling datastore readers for %s", vmpath.String())
		// obtain datastore object
		ds, err := c.Finder.Datastore(ctx, vmpath.Host)
		if err != nil {
			log.Errorf("Failed to acquire reference to datastore %s: %s", vmpath.Host, err)
			continue
		}

		// generate the full paths to collect
		for _, file := range vmFiles {
			// replace the VM token in file name with the VM name
			processed := strings.Replace(file, string(metadata.VM), path.Base(vmpath.Path), -1)
			rpath := fmt.Sprintf("%s/%s", vmpath.Path, processed)
			readers[rpath] = datastoreReader{
				ds:   ds,
				path: rpath,
			}

			log.Debugf("Added log file for collection: %s", vmpath.String())
		}
	}

	return readers, nil
}

func (r datastoreReader) open() (entry, error) {
	defer trace.End(trace.Begin(r.path))

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

func client() (*session.Session, error) {
	defer trace.End(trace.Begin(""))

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
	defer trace.End(trace.Begin(""))

	session, err := client()
	if err != nil {
		return err
	}
	defer session.Client.Logout(context.Background())

	datastore = session.Datastore.Reference()
	datastoreInventoryPath = session.Datastore.InventoryPath

	return nil
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

type server struct {
	auth  Authenticator
	l     net.Listener
	addr  string
	mux   *http.ServeMux
	links []string
}

func (s *server) listen(useTLS bool) error {
	defer trace.End(trace.Begin(""))

	var err error

	// Set options for TLS
	tlsconfig := tlsconfig.ServerDefault
	certificate, err := vchConfig.HostCertificate.Certificate()
	if err != nil {
		log.Errorf("Could not load certificate from config - running without TLS: %s", err)
		// TODO: add static web page with the vic
	} else {
		tlsconfig.Certificates = []tls.Certificate{*certificate}
	}

	if !useTLS || err != nil {
		s.l, err = net.Listen("tcp", s.addr)
		return err
	}

	innerListener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	s.l = tls.NewListener(innerListener, &tlsconfig)
	return nil
}

func (s *server) listenPort() int {
	return s.l.Addr().(*net.TCPAddr).Port
}

// handleFunc does preparatory work and then calls the HandleFunc method owned by the HTTP multiplexer
func (s *server) handleFunc(link string, handler func(http.ResponseWriter, *http.Request)) {
	defer trace.End(trace.Begin(""))

	s.links = append(s.links, link)

	if s.auth != nil {
		authHandler := func(w http.ResponseWriter, r *http.Request) {
			user, password, ok := r.BasicAuth()
			if !ok || !s.auth.Validate(user, password) {
				w.Header().Add("WWW-Authenticate", "Basic realm=vicadmin")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			handler(w, r)
		}
		s.mux.HandleFunc(link, authHandler)
		return
	}

	s.mux.HandleFunc(link, handler)
}

func (s *server) serve() error {
	defer trace.End(trace.Begin(""))

	s.mux = http.NewServeMux()

	// tar of appliance system logs
	s.handleFunc("/logs.tar.gz", s.tarDefaultLogs)

	// tar of appliance system logs + container logs
	s.handleFunc("/container-logs.tar.gz", s.tarContainerLogs)

	// tail all logFiles
	s.handleFunc("/logs/tail", func(w http.ResponseWriter, r *http.Request) {
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
		s.handleFunc("/logs/tail/"+name, func(w http.ResponseWriter, r *http.Request) {
			s.tailFiles(w, r, []string{p})
		})
	}

	s.handleFunc("/", s.index)
	server := &http.Server{
		Handler: s.mux,
	}

	defaultReaders = configureReaders()

	return server.Serve(s.l)
}

func (s *server) stop() error {
	defer trace.End(trace.Begin(""))

	if s.l != nil {
		err := s.l.Close()
		s.l = nil
		return err
	}

	return nil
}

func (s *server) tarContainerLogs(res http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))

	readers := defaultReaders

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
				for key, rdr := range logs {
					readers[key] = rdr
				}
			}

			logs, err = findDiagnosticLogs(c)
			if err != nil {
				log.Warningf("error collecting diagnostic logs: %s", err)
			} else {
				for key, rdr := range logs {
					readers[key] = rdr
				}
			}
		}
	}

	s.tarLogs(res, req, readers)
}

func (s *server) tarDefaultLogs(res http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))

	s.tarLogs(res, req, defaultReaders)
}

func (s *server) tarLogs(res http.ResponseWriter, req *http.Request, readers map[string]entryReader) {
	defer trace.End(trace.Begin(""))

	res.Header().Set("Content-Type", "application/x-gzip")

	z := gzip.NewWriter(res)

	err := tarEntries(readers, z)
	if err != nil {
		log.Errorf("error taring logs: %s", err)
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
	defer trace.End(trace.Begin(""))

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
		log.Errorf("error tailing logs: %s", err)
		return
	}

	defer cmd.Process.Kill()

	go cmd.Wait()

	<-cc
}

func (s *server) index(res http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))

	fmt.Fprintln(res, "<html><head><title>VIC Admin</title></head><body><pre>")
	for _, link := range s.links {
		fmt.Fprintf(res, "<a href=\"%s\">%s</a><br/>\n", link, link)
	}
	fmt.Fprintln(res, "</pre></body></html>")
}

func main() {
	defer trace.End(trace.Begin(""))

	flag.Parse()

	s := &server{
		addr: config.addr,
	}

	err := s.listen(config.tls)

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("listening on %s", s.addr)
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
