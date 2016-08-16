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
	"archive/zip"
	"compress/gzip"
	"crypto/tls"
	"html/template"
	"net"
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/tlsconfig"
	"github.com/vmware/vic/lib/vicadmin"
	"github.com/vmware/vic/pkg/trace"
)

type server struct {
	auth Authenticator
	l    net.Listener
	addr string
	mux  *http.ServeMux
}

type format int

const (
	formatTGZ format = iota
	formatZip
)

func (s *server) listen(useTLS bool) error {
	defer trace.End(trace.Begin(""))

	var err error

	// FIXME: assignment copies lock value to tlsConfig: crypto/tls.Config contains sync.Once contains sync.Mutex
	tlsconfig := func(c *tls.Config) *tls.Config {
		return &tls.Config{
			Certificates:             c.Certificates,
			NameToCertificate:        c.NameToCertificate,
			GetCertificate:           c.GetCertificate,
			RootCAs:                  c.RootCAs,
			NextProtos:               c.NextProtos,
			ServerName:               c.ServerName,
			ClientAuth:               c.ClientAuth,
			ClientCAs:                c.ClientCAs,
			InsecureSkipVerify:       c.InsecureSkipVerify,
			CipherSuites:             c.CipherSuites,
			PreferServerCipherSuites: c.PreferServerCipherSuites,
			SessionTicketsDisabled:   c.SessionTicketsDisabled,
			SessionTicketKey:         c.SessionTicketKey,
			ClientSessionCache:       c.ClientSessionCache,
			MinVersion:               c.MinVersion,
			MaxVersion:               c.MaxVersion,
			CurvePreferences:         c.CurvePreferences,
		}
	}(&tlsconfig.ServerDefault)

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

	s.l = tls.NewListener(innerListener, tlsconfig)
	return nil
}

func (s *server) listenPort() int {
	return s.l.Addr().(*net.TCPAddr).Port
}

// handleFunc does preparatory work and then calls the HandleFunc method owned by the HTTP multiplexer
func (s *server) handleFunc(link string, handler func(http.ResponseWriter, *http.Request)) {
	defer trace.End(trace.Begin(""))

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
	s.mux.HandleFunc("/logs.tar.gz", s.tarDefaultLogs)
	s.handleFunc("/logs.zip", s.zipDefaultLogs)

	// tar of appliance system logs + container logs
	s.mux.HandleFunc("/container-logs.tar.gz", s.tarContainerLogs)
	s.handleFunc("/container-logs.zip", s.zipContainerLogs)

	s.mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	s.mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))
	s.mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts/"))))

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

	//s.handleFunc("/", s.index)
	s.mux.HandleFunc("/", s.index)
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

func (s *server) bundleContainerLogs(res http.ResponseWriter, req *http.Request, f format) {
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

	s.bundleLogs(res, req, readers, f)
}

func (s *server) tarDefaultLogs(res http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))

	s.bundleLogs(res, req, defaultReaders, formatTGZ)
}
func (s *server) zipDefaultLogs(res http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))

	s.bundleLogs(res, req, defaultReaders, formatZip)
}

func (s *server) bundleLogs(res http.ResponseWriter, req *http.Request, readers map[string]entryReader, f format) {
	defer trace.End(trace.Begin(""))

	var err error
	if f == formatTGZ {
		res.Header().Set("Content-Type", "application/x-gzip")
		z := gzip.NewWriter(res)
		defer z.Close()
		err = tarEntries(readers, z)
	} else if f == formatZip {
		res.Header().Set("Content-Type", "application/zip")
		z := zip.NewWriter(res)
		defer z.Close()
		err = zipEntries(readers, z)
	}

	if err != nil {
		log.Errorf("Error bundling logs: %s", err)
	}
}

func (s *server) tarContainerLogs(res http.ResponseWriter, req *http.Request) {
	s.bundleContainerLogs(res, req, formatTGZ)
}

func (s *server) zipContainerLogs(res http.ResponseWriter, req *http.Request) {
	s.bundleContainerLogs(res, req, formatZip)
}

func (s *server) tailFiles(res http.ResponseWriter, req *http.Request, names []string) {
	defer trace.End(trace.Begin(""))

	cc := res.(http.CloseNotifier).CloseNotify()

	fw := &flushWriter{
		f: res.(http.Flusher),
		w: res,
	}

	done := make(chan bool)
	for _, file := range names {
		go tailFile(fw, file, &done)
	}

	<-cc
	for _, _ = range names {
		done <- true
	}
}

func (s *server) index(res http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))
	ctx := context.Background()
	sess, err := client()
	v := vicadmin.NewValidator(ctx, &vchConfig, sess)

	tmpl, err := template.ParseFiles("dashboard.html")
	err = tmpl.ExecuteTemplate(res, "dashboard.html", v)
	if err != nil {
		log.Errorf("Error parsing template: %s", err)
	}
}
