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
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/certificate"
	"github.com/vmware/vic/pkg/trace"
)

const (
	publicNetName = "public-net"
	bridgeNetName = "bridge-net"
	imgStoreName  = "img-store"
	computeName   = "compute"
)

var (
	engineInstaller = NewEngineInstaller()
	c               config
)

type config struct {
	addr     string
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

	flag.StringVar(&conf.addr, "addr", ":1337", "Listen address - must include host and port (addr:port)")
	flag.StringVar(&conf.serveDir, "data", "/opt/vmware/engine_installer", "Directory containing vic-machine and HTML data")

	flag.Parse()

	generateCert(conf)
}

func main() {
	Init(&c)

	defer trace.End(trace.Begin(""))
	log.Infoln("starting installer-egine")

	mux := http.NewServeMux()

	// attact static asset routes
	routes := []string{"css", "images", "fonts"}
	for _, route := range routes {
		httpPath := fmt.Sprintf("/%s/", route)
		dirPath := fmt.Sprintf("%s/html/%s/", c.serveDir, route)
		mux.Handle(httpPath, http.StripPrefix(httpPath, http.FileServer(http.Dir(dirPath))))
	}

	// attach root index route
	mux.Handle("/", http.HandlerFunc(indexHandler))
	mux.Handle("/ws", http.HandlerFunc(logStream.wsServer))
	mux.Handle("/cmd", http.HandlerFunc(parseCmdArgs))

	// start the web server
	t := &tls.Config{}
	t.Certificates = []tls.Certificate{c.cert}
	s := &http.Server{
		Addr:      c.addr,
		Handler:   mux,
		TLSConfig: t,
	}

	log.Infof("Starting installer-engine server on %s", s.Addr)
	log.Fatal(s.ListenAndServeTLS("", ""))

	// start the web server
	// log.Infof("installer-engine listening on localhost:%s\n", c.addr)
	// http.ListenAndServe(c.addr, mux)
}

func generateCert(conf *config) {
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

func indexHandler(resp http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))

	if req.Method == http.MethodPost {
		// verify login

		engineInstaller.Target = req.FormValue("target")
		engineInstaller.User = req.FormValue("user")
		engineInstaller.Password = req.FormValue("password")
		if err := engineInstaller.verifyLogin(); err != nil {
			// login failed so show login form again
			renderTemplate(resp, "html/auth.html", &AuthHTML{InvalidLogin: true})
		} else {
			// vCenter login successful, set resource drop downs
			opts := engineInstaller.populateConfigOptions()
			html := &ExecHTMLOptions{}

			//add bridge select box to html template
			html.PublicNetwork = getSelectOptionHTML(opts.Networks, publicNetName)

			//add bridge select box to html template
			html.BridgeNetwork = getSelectOptionHTML(opts.Networks, bridgeNetName)

			//add datastores select box to html template
			html.ImageStore = getSelectOptionHTML(opts.Datastores, imgStoreName)

			//add compute resources select box to html template
			html.ComputeResource = getSelectOptionHTML(opts.ResourcePools, computeName)

			html.User = engineInstaller.User
			html.Password = engineInstaller.Password
			html.Target = engineInstaller.Target
			html.Name = engineInstaller.Name
			html.Thumbprint = engineInstaller.Thumbprint
			html.CreateCommand = engineInstaller.CreateCommand

			renderTemplate(resp, "html/exec.html", html)
		}
	} else {
		renderTemplate(resp, "html/auth.html", nil)
	}
}

func getSelectOptionHTML(arr []string, id string) template.HTML {
	templ := template.HTML(fmt.Sprintf("<div class=\"select\"><select name=\"%s\">", id))
	for _, option := range arr {
		optionHTML := fmt.Sprintf("<option>%s</option>", option)
		templ = template.HTML(fmt.Sprintf("%s\n%s", templ, optionHTML))
	}
	return template.HTML(fmt.Sprintf("%s\n</select></div>", templ))
}

func renderTemplate(resp http.ResponseWriter, filename string, data interface{}) {
	log.Infof("render: %s", filename)
	filename = fmt.Sprintf("%s/%s", c.serveDir, filename)
	log.Infof("render: %s", filename)
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(resp, data); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
}

func parseCmdArgs(resp http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(&engineInstaller); err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte("500"))
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %v", e)
		}
	} else {
		//build the vic create command from the installer variables
		engineInstaller.buildCreateCommand(c.serveDir)
		log.Infoln(engineInstaller)
		resp.WriteHeader(200)
		resp.Write([]byte(engineInstaller.CreateCommand))
	}

}
