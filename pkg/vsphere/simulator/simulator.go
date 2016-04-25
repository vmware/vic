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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"

	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/xml"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"github.com/vmware/vic/pkg/vsphere/simulator/vc"
)

// Method encapsulates a decoded SOAP client request
type Method struct {
	Name string
	Body interface{}
}

// Handler implements a response handler for a Method request
type Handler interface {
	Name() string
	Call(*Method) soap.HasFault
}

// Service decodes incoming requests and dispatches to a Handler
type Service struct {
	handlers map[string]Handler

	readAll func(io.Reader) ([]byte, error)
}

// Server provides a simulator Service over HTTP
type Server struct {
	*httptest.Server
	URL *url.URL
}

// New returns an initialized simulator Service instance
func New() *Service {
	s := &Service{
		handlers: make(map[string]Handler),

		readAll: ioutil.ReadAll,
	}

	return s
}

// Register method handlers with a Service
func (s *Service) Register(handlers ...Handler) {
	for _, handler := range handlers {
		s.handlers[handler.Name()] = handler
	}
}

var handlers = []Handler{
	&Login{},
	&CurrentTime{},
}

// AsESX registers service handlers with an ESX personality
func (s *Service) AsESX() *Service {
	h := []Handler{
		&RetrieveServiceContent{esx.ServiceContent},
	}
	s.Register(append(handlers, h...)...)
	return s
}

// AsVC registers service handlers with a vCenter personality
func (s *Service) AsVC() *Service {
	h := []Handler{
		&RetrieveServiceContent{vc.ServiceContent},
	}
	s.Register(append(handlers, h...)...)
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
		if handler, ok := s.handlers[method.Name]; ok {
			res = handler.Call(method)
		} else {
			res = serverFault(fmt.Sprintf("no handler for %s", method.Name))
		}
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

// NewServer returns an http Server instance for the given service
func (s *Service) NewServer() *Server {
	mux := http.NewServeMux()
	path := "/sdk"
	mux.Handle(path, s)

	ts := httptest.NewServer(mux)

	u, _ := url.Parse(ts.URL)
	u.Path = path

	return &Server{
		Server: ts,
		URL:    u,
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

	return &Method{Name: kind, Body: val}, nil
}
