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

package simulator

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/xml"
)

// Method encapsulates a decoded SOAP client request
type Method struct {
	Name string
	This types.ManagedObjectReference
	Body types.AnyType
}

// Service decodes incoming requests and dispatches to a Handler
type Service struct {
	client *vim25.Client

	readAll func(io.Reader) ([]byte, error)

	TLS *tls.Config
}

// Server provides a simulator Service over HTTP
type Server struct {
	*httptest.Server
	URL *url.URL

	caFile string
}

// New returns an initialized simulator Service instance
func New(instance *ServiceInstance) *Service {
	s := &Service{
		readAll: ioutil.ReadAll,
	}

	s.client, _ = vim25.NewClient(context.Background(), s)

	return s
}

type serverFaultBody struct {
	Reason *soap.Fault `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault,omitempty"`
}

func (b *serverFaultBody) Fault() *soap.Fault { return b.Reason }

func serverFault(msg string) soap.HasFault {
	return &serverFaultBody{Reason: Fault(msg, &types.InvalidRequest{})}
}

// Fault wraps the given message and fault in a soap.Fault
func Fault(msg string, fault types.BaseMethodFault) *soap.Fault {
	f := &soap.Fault{
		Code:   "ServerFaultCode",
		String: msg,
	}

	f.Detail.Fault = fault

	return f
}

func (s *Service) call(method *Method) soap.HasFault {
	handler := Map.Get(method.This)

	if handler == nil {
		return &serverFaultBody{
			Reason: Fault(fmt.Sprintf("no such object: %s", method.This),
				&types.ManagedObjectNotFound{Obj: method.This}),
		}
	}

	name := method.Name

	if strings.HasSuffix(name, vTaskSuffix) {
		// Make golint happy renaming "Foo_Task" -> "FooTask"
		name = name[:len(name)-len(vTaskSuffix)] + sTaskSuffix
	}

	m := reflect.ValueOf(handler).MethodByName(name)
	if !m.IsValid() {
		return serverFault(fmt.Sprintf("%s does not implement: %s", method.This, method.Name))
	}

	res := m.Call([]reflect.Value{reflect.ValueOf(method.Body)})

	return res[0].Interface().(soap.HasFault)
}

// RoundTrip implements the soap.RoundTripper interface in process.
// Rather than encode/decode SOAP over HTTP, this implementation uses reflection.
func (s *Service) RoundTrip(ctx context.Context, request, response soap.HasFault) error {
	field := func(r soap.HasFault, name string) reflect.Value {
		return reflect.ValueOf(r).Elem().FieldByName(name)
	}

	// Every struct passed to soap.RoundTrip has "Req" and "Res" fields
	req := field(request, "Req")

	// Every request has a "This" field.
	this := req.Elem().FieldByName("This")

	method := &Method{
		Name: req.Elem().Type().Name(),
		This: this.Interface().(types.ManagedObjectReference),
		Body: req.Interface(),
	}

	res := s.call(method)

	if err := res.Fault(); err != nil {
		return soap.WrapSoapFault(err)
	}

	field(response, "Res").Set(field(res, "Res"))

	return nil
}

// ServeHTTP implements the http.Handler interface
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := s.readAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		log.Printf("error reading body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var res soap.HasFault

	method, err := UnmarshalBody(body)
	if err != nil {
		res = serverFault(err.Error())
	} else {
		res = s.call(method)
	}

	if res.Fault() == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	e := xml.NewEncoder(w)
	err = e.Encode(&soap.Envelope{Body: res})
	if err == nil {
		err = e.Flush()
	}

	if err != nil {
		log.Printf("error encoding %s response: %s", method.Name, err)
	}
}

func (s *Service) findDatastore(query url.Values) (*Datastore, error) {
	ctx := context.Background()

	finder := find.NewFinder(s.client, false)
	dc, err := finder.DatacenterOrDefault(ctx, query.Get("dcName"))
	if err != nil {
		return nil, err
	}

	finder.SetDatacenter(dc)

	ds, err := finder.DatastoreOrDefault(ctx, query.Get("dsName"))
	if err != nil {
		return nil, err
	}

	return Map.Get(ds.Reference()).(*Datastore), nil
}

const folderPrefix = "/folder/"

func (s *Service) ServeDatastore(w http.ResponseWriter, r *http.Request) {
	ds, ferr := s.findDatastore(r.URL.Query())
	if ferr != nil {
		log.Printf("failed to locate datastore with query params: %s", r.URL.RawQuery)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file := strings.TrimPrefix(r.URL.Path, folderPrefix)
	p := path.Join(ds.Info.GetDatastoreInfo().Url, file)

	switch r.Method {
	case "GET":
		f, err := os.Open(p)
		if err != nil {
			log.Printf("failed to %s '%s': %s", r.Method, p, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer f.Close()

		_, _ = io.Copy(w, f)
	case "POST":
		_, err := os.Stat(p)
		if err == nil {
			// File exists
			w.WriteHeader(http.StatusConflict)
			return
		}

		// File does not exist, fallthrough to create via PUT logic
		fallthrough
	case "PUT":
		f, err := os.Create(p)
		if err != nil {
			log.Printf("failed to %s '%s': %s", r.Method, p, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, _ = io.Copy(f, r.Body)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// NewServer returns an http Server instance for the given service
func (s *Service) NewServer() *Server {
	mux := http.NewServeMux()
	path := "/sdk"
	mux.Handle(path, s)

	mux.HandleFunc(folderPrefix, s.ServeDatastore)

	// Using NewUnstartedServer() instead of NewServer(),
	// for use in main.go, where Start() blocks, we can still set ServiceHostName
	ts := httptest.NewUnstartedServer(mux)

	u := &url.URL{
		Scheme: "http",
		Host:   ts.Listener.Addr().String(),
		Path:   path,
		User:   url.UserPassword("user", "pass"),
	}

	// Enable use of SessionManagerGenericServiceTicket.HostName in govmomi, disabled by default.
	opts := u.Query()
	opts.Set("GOVMOMI_USE_SERVICE_TICKET_HOSTNAME", "true")
	u.RawQuery = opts.Encode()

	// Redirect clients to this http server, rather than HostSystem.Name
	Map.Get(*s.client.ServiceContent.SessionManager).(*SessionManager).ServiceHostName = u.Host

	if f := flag.Lookup("httptest.serve"); f != nil {
		// Avoid the blocking behaviour of httptest.Server.Start() when this flag is set
		_ = f.Value.Set("")
	}

	if s.TLS == nil {
		ts.Start()
	} else {
		ts.TLS = s.TLS
		ts.StartTLS()
		u.Scheme += "s"
	}

	return &Server{
		Server: ts,
		URL:    u,
	}
}

// Certificate returns the TLS certificate for the Server if started with TLS enabled.
// This method will panic if TLS is not enabled for the server.
func (s *Server) Certificate() *x509.Certificate {
	// By default httptest.StartTLS uses http/internal.LocalhostCert, which we can access here:
	cert, _ := x509.ParseCertificate(s.TLS.Certificates[0].Certificate[0])
	return cert
}

// CertificateInfo returns Server.Certificate() as object.HostCertificateInfo
func (s *Server) CertificateInfo() *object.HostCertificateInfo {
	info := new(object.HostCertificateInfo)
	info.FromCertificate(s.Certificate())
	return info
}

// CertificateFile returns a file name, where the file contains the PEM encoded Server.Certificate.
// The temporary file is removed when Server.Close() is called.
func (s *Server) CertificateFile() (string, error) {
	if s.caFile != "" {
		return s.caFile, nil
	}

	f, err := ioutil.TempFile("", "vcsim-")
	if err != nil {
		return "", err
	}
	defer f.Close()

	s.caFile = f.Name()
	cert := s.Certificate()
	return s.caFile, pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
}

// Close shuts down the server and blocks until all outstanding
// requests on this server have completed.
func (s *Server) Close() {
	s.Server.Close()
	if s.caFile != "" {
		_ = os.Remove(s.caFile)
	}
}

var typeFunc = types.TypeFunc()

// UnmarshalBody extracts the Body from a soap.Envelope and unmarshals to the corresponding govmomi type
func UnmarshalBody(data []byte) (*Method, error) {
	body := struct {
		Content string `xml:",innerxml"`
	}{}

	req := soap.Envelope{
		Body: &body,
	}

	err := xml.Unmarshal(data, &req)
	if err != nil {
		return nil, fmt.Errorf("xml.Unmarshal: %s", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader([]byte(body.Content)))
	decoder.TypeFunc = typeFunc // required to decode interface types

	var start *xml.StartElement

	for {
		tok, derr := decoder.Token()
		if derr != nil {
			return nil, fmt.Errorf("decoding body: %s", err)
		}
		if t, ok := tok.(xml.StartElement); ok {
			start = &t
			break
		}
	}

	kind := start.Name.Local

	rtype, ok := typeFunc(kind)
	if !ok {
		return nil, fmt.Errorf("no vmomi type defined for '%s'", kind)
	}

	var val interface{}
	if rtype != nil {
		val = reflect.New(rtype).Interface()
	}

	err = decoder.DecodeElement(val, start)
	if err != nil {
		return nil, fmt.Errorf("decoding %s: %s", kind, err)
	}

	method := &Method{Name: kind, Body: val}

	field := reflect.ValueOf(val).Elem().FieldByName("This")

	method.This = field.Interface().(types.ManagedObjectReference)

	return method, nil
}
