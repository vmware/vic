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
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/certificate"
)

type config struct {
	addr     string
	certPath string
	keyPath  string
	cert     tls.Certificate
	serveDir string
}

func Init(conf *config) {
	ud := syscall.Getuid()
	gd := syscall.Getgid()
	log.Info(fmt.Sprintf("Current UID/GID = %d/%d", ud, gd))
	/* TODO FIXME
	if ud == 0 {
		log.Error("Error: must not run as root.")
		os.Exit(1)
	}
	*/

	flag.StringVar(&conf.addr, "addr", ":8443", "Listen address - must include host and port (addr:port)")
	flag.StringVar(&conf.certPath, "cert", "", "Path to server certificate in PEM format")
	flag.StringVar(&conf.keyPath, "key", "", "Path to server certificate key in PEM format")
	flag.StringVar(&conf.serveDir, "dir", "/data/fileserver", "Directory to serve")

	flag.Parse()

	if (conf.certPath == "" && conf.keyPath != "") || (conf.certPath != "" && conf.keyPath == "") {
		log.Errorf("Both certificate and key must be specified")
	}

	var err error
	if conf.certPath != "" {
		log.Infof("Loading certificate %s and key %s", conf.certPath, conf.keyPath)
		conf.cert, err = tls.LoadX509KeyPair(conf.certPath, conf.keyPath)
		if err != nil {
			log.Fatalf("Failed to load certificate %s and key %s: %s", conf.certPath, conf.keyPath, err)
		}
	} else {
		log.Info("Generating self signed certificate")
		c, k, err := certificate.CreateSelfSigned(conf.addr, []string{"VMware, Inc."}, 2048)
		if err != nil {
			log.Errorf("Failed to generate a self-signed certificate: %s. Exiting.", err.Error())
			os.Exit(1)
		}
		conf.cert, err = tls.X509KeyPair(c.Bytes(), k.Bytes())
		if err != nil {
			log.Errorf("Failed to load generated self-signed certificate: %s. Exiting.", err.Error())
			os.Exit(1)
		}
	}
	log.Infof("Loaded certificate")
}

func main() {
	var c config
	Init(&c)

	t := &tls.Config{}
	t.Certificates = []tls.Certificate{c.cert}
	s := &http.Server{
		Addr:      c.addr,
		Handler:   http.FileServer(http.Dir(c.serveDir)),
		TLSConfig: t,
	}

	log.Infof("Starting server on %s", s.Addr)
	log.Fatal(s.ListenAndServeTLS("", ""))
}
