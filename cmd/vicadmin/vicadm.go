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
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"

	vchconfig "github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/pprof"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
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
		"output.log",
		"vmware.log",
		"tether.debug",
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

	resources vchconfig.Resources

	vchConfig vchconfig.VirtualContainerHostConfigSpec

	defaultReaders map[string]entryReader

	datastore types.ManagedObjectReference

	datastoreInventoryPath string
)

type logfile struct {
	URL    url.URL
	VMName string
}

func init() {
	defer trace.End(trace.Begin(""))
	trace.Logger.Level = log.DebugLevel
	_ = pprof.StartPprof("vicadmin", pprof.VicadminPort)

	// We don't want to run this as root.
	ud := syscall.Getuid()
	gd := syscall.Getgid()
	log.Info(fmt.Sprintf("Current UID/GID = %d/%d", ud, gd))
	// TODO: Enable this after we figure out to NOT break the test suite with it.
	// if ud == 0 {
	// log.Errorf("Error: vicadmin must not run as root.")
	// time.Sleep(60 * time.Second)
	// os.Exit(1)
	// }

	flag.StringVar(&config.addr, "l", ":2378", "Listen address")
	flag.StringVar(&config.dockerHost, "docker-host", "127.0.0.1:2376", "Docker host")
	flag.StringVar(&config.hostCertFile, "hostcert", "", "Host certificate file")
	flag.StringVar(&config.hostKeyFile, "hostkey", "", "Host private key file")
	flag.StringVar(&config.DatacenterPath, "dc", "", "Name of the Datacenter")
	flag.StringVar(&config.DatastorePath, "ds", "", "Name of the Datastore")
	flag.StringVar(&config.ClusterPath, "cluster", "", "Path of the cluster")
	flag.StringVar(&config.PoolPath, "pool", "", "Path of the resource pool")
	flag.BoolVar(&config.Insecure, "insecure", false, "Allow connection when sdk certificate cannot be verified")
	flag.BoolVar(&config.tls, "tls", true, "Set to false to disable -hostcert and -hostkey and enable plain HTTP")

	// This is only applicable for containers hosted under the VCH VM folder
	// This will not function for vSAN
	flag.StringVar(&config.vmPath, "vm-path", "", "Docker vm path")

	flag.Parse()

	// load the vch config
	src, err := extraconfig.GuestInfoSource()
	if err != nil {
		log.Errorf("Unable to load configuration from guestinfo")
		return
	}

	extraconfig.Decode(src, &vchConfig)
}

type Authenticator interface {
	// Validate will validate a user and password combo and return a bool.
	Validate(string, string) bool
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

// listVMPaths returns an array of datastore paths for VMs associated with the
// VCH - this includes containerVMs and the appliance
func listVMPaths(ctx context.Context, s *session.Session) ([]logfile, error) {
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

	logfiles := []logfile{}
	for _, child := range children {
		path, err := child.DSPath(ctx)

		if err != nil {
			log.Errorf("Unable to get datastore path for child VM %s: %s", child.Reference(), err)
			// we need to get as many logs as possible
			continue
		}

		logname, err := child.Name(ctx)
		if err != nil {
			log.Errorf("Unable to get the vm name for %s: %s", child.Reference(), err)
			continue
		}

		log.Debugf("Adding VM for log collection: %s", path.String())

		log := logfile{
			URL:    path,
			VMName: logname,
		}

		logfiles = append(logfiles, log)
	}

	log.Infof("Collecting logs from %d VMs", len(logfiles))
	log.Infof("Found VM paths are : %#v", logfiles)
	return logfiles, nil
}

// find datastore logs for the appliance itself and all containers
func findDatastoreLogs(c *session.Session) (map[string]entryReader, error) {
	defer trace.End(trace.Begin(""))

	// Create an empty reader as opposed to a nil reader...
	readers := map[string]entryReader{}
	ctx := context.Background()

	logfiles, err := listVMPaths(ctx, c)
	if err != nil {
		detail := fmt.Sprintf("unable to perform datastore log collection due to failure looking up paths: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	for _, logfile := range logfiles {
		log.Debugf("Assembling datastore readers for %s", logfile.URL.String())
		// obtain datastore object
		ds, err := c.Finder.Datastore(ctx, logfile.URL.Host)
		if err != nil {
			log.Errorf("Failed to acquire reference to datastore %s: %s", logfile.URL.Host, err)
			continue
		}

		// generate the full paths to collect
		for _, file := range vmFiles {
			wpath := fmt.Sprintf("%s/%s", logfile.VMName, file)
			rpath := fmt.Sprintf("%s/%s", logfile.URL.Path, file)
			log.Infof("Processed File read Path : %s", rpath)
			log.Infof("Processed File write Path : %s", wpath)
			readers[wpath] = datastoreReader{
				ds:   ds,
				path: rpath,
			}

			log.Debugf("Added log file for collection: %s", logfile.URL.String())
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

func client() (*session.Session, error) {
	defer trace.End(trace.Begin(""))

	ctx := context.Background()

	session := session.NewSession(&config.Config)
	_, err := session.Connect(ctx)
	if err != nil {
		log.Warnf("Unable to connect: %s", err)
		return nil, err
	}

	_, err = session.Populate(ctx)
	if err != nil {
		// no a critical error for vicadmin
		log.Warnf("Unable to populate session: %s", err)
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

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (int, error) {
	n, err := fw.w.Write(p)

	fw.f.Flush()

	return n, err
}

func main() {
	if version.Show() {
		fmt.Fprintf(os.Stdout, "%s\n", version.String())
		return
	}

	flag.Parse()

	// If we're in an ESXi environment, then we need
	// to extract the userid/password from UserPassword
	if vchConfig.UserPassword != "" {
		newurl, _ := url.Parse(fmt.Sprintf("%s://%s@%s%s",
			vchConfig.Target.Scheme,
			vchConfig.UserPassword,
			vchConfig.Target.Host,
			vchConfig.Target.Path))
		vchConfig.Target = *newurl
	}

	config.Service = vchConfig.Target.String()
	config.ExtensionCert = vchConfig.ExtensionCert
	config.ExtensionKey = vchConfig.ExtensionKey
	config.ExtensionName = vchConfig.ExtensionName

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
