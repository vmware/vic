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
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"github.com/vmware/vic/pkg/vsphere/simulator/vc"
)

func TestUnmarshal(t *testing.T) {
	requests := []struct {
		body interface{}
		data string
	}{
		{
			&types.RetrieveServiceContent{
				This: types.ManagedObjectReference{
					Type: "ServiceInstance", Value: "ServiceInstance",
				},
			},
			`<?xml version="1.0" encoding="UTF-8"?>
                         <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
                           <Body>
                             <RetrieveServiceContent xmlns="urn:vim25">
                               <_this type="ServiceInstance">ServiceInstance</_this>
                             </RetrieveServiceContent>
                           </Body>
                         </Envelope>`,
		},
		{
			&types.Login{
				This: types.ManagedObjectReference{
					Type:  "SessionManager",
					Value: "SessionManager",
				},
				UserName: "root",
				Password: "secret",
			},
			`<?xml version="1.0" encoding="UTF-8"?>
                         <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
                           <Body>
                             <Login xmlns="urn:vim25">
                               <_this type="SessionManager">SessionManager</_this>
                               <userName>root</userName>
                               <password>secret</password>
                             </Login>
                           </Body>
                         </Envelope>`,
		},
	}

	for i, req := range requests {
		method, err := UnmarshalBody([]byte(req.data))
		if err != nil {
			t.Errorf("failed to decode %d (%s):", i, req, err)
		}
		if !reflect.DeepEqual(method.Body, req.body) {
			t.Errorf("malformed body %d (%#v):", i, method.Body)
		}
	}
}

func TestUnmarshalError(t *testing.T) {
	requests := []string{
		"", // io.EOF
		`<?xml version="1.0" encoding="UTF-8"?>
                 <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
                   <Body>
                   </MissingEndTag
                 </Envelope>`,
		`<?xml version="1.0" encoding="UTF-8"?>
                 <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
                   <Body>
                     <UnknownType xmlns="urn:vim25">
                       <_this type="ServiceInstance">ServiceInstance</_this>
                     </UnknownType>
                   </Body>
                 </Envelope>`,
		`<?xml version="1.0" encoding="UTF-8"?>
                 <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
                   <Body>
                   <!-- no start tag -->
                   </Body>
                 </Envelope>`,
		`<?xml version="1.0" encoding="UTF-8"?>
                 <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
                   <Body>
                     <RetrieveServiceContent xmlns="urn:vim25">
                       <_this type="ServiceInstance">ServiceInstance</_this>
                     </RetrieveServiceContent>
                   </Body>
                 </Envelope>`,
	}

	defer func() {
		typeFunc = types.TypeFunc() // reset
	}()

	ttypes := map[string]reflect.Type{
		// triggers xml.Decoder.DecodeElement error
		"RetrieveServiceContent": reflect.TypeOf(nil),
	}
	typeFunc = func(name string) (reflect.Type, bool) {
		typ, ok := ttypes[name]
		return typ, ok
	}

	for i, data := range requests {
		_, err := UnmarshalBody([]byte(data))
		if err != nil {
			continue
		}
		t.Errorf("expected %d (%s) to return an error", i, data)
	}
}

func TestServeHTTP(t *testing.T) {
	services := []*Service{
		New(NewServiceInstance(esx.ServiceContent)),
		New(NewServiceInstance(vc.ServiceContent)),
	}

	for _, s := range services {
		ts := s.NewServer()
		defer ts.Close()

		ctx := context.Background()
		client, err := govmomi.NewClient(ctx, ts.URL, true)
		if err != nil {
			t.Fatal(err)
		}

		err = client.Login(ctx, nil)
		if err == nil {
			t.Fatal("expected invalid login error")
		}

		err = client.Login(ctx, url.UserPassword("user", "pass"))
		if err != nil {
			t.Fatal(err)
		}

		now, err := methods.GetCurrentTime(ctx, client)
		if err != nil {
			t.Fatal(err)
		}

		if now.After(time.Now()) {
			t.Error()
		}
	}
}

type errorMarshal struct {
	mo.ServiceInstance
}

func (*errorMarshal) Fault() *soap.Fault {
	return nil
}

func (*errorMarshal) MarshalText() ([]byte, error) {
	return nil, errors.New("time has stopped")
}

func (h *errorMarshal) CurrentTime(types.AnyType) soap.HasFault {
	return h
}

type errorNoSuchMethod struct {
	mo.ServiceInstance
}

type errorInvalidMethod struct {
	mo.ServiceInstance
}

func (h *errorInvalidMethod) CurrentTime() soap.HasFault {
	return serverFault("notreached")
}

func TestServeHTTPErrors(t *testing.T) {
	s := New(NewServiceInstance(esx.ServiceContent))

	ts := s.NewServer()
	defer ts.Close()

	ctx := context.Background()
	client, err := govmomi.NewClient(ctx, ts.URL, true)
	if err != nil {
		t.Fatal(err)
	}

	// unregister type, covering the ServeHTTP UnmarshalBody error path
	typeFunc = func(name string) (reflect.Type, bool) {
		return nil, false
	}

	_, err = methods.GetCurrentTime(ctx, client)
	if err == nil {
		t.Error("expected error")
	}

	typeFunc = types.TypeFunc() // reset

	// cover the does not implement method error path
	s.handlers[serviceInstance] = &errorNoSuchMethod{}
	_, err = methods.GetCurrentTime(ctx, client)
	if err == nil {
		t.Error("expected error")
	}

	// cover the invalid method error path
	s.handlers[serviceInstance] = &errorInvalidMethod{}
	_, err = methods.GetCurrentTime(ctx, client)
	if err == nil {
		t.Error("expected error")
	}

	// cover the xml encode error path
	s.handlers[serviceInstance] = &errorMarshal{}
	_, err = methods.GetCurrentTime(ctx, client)
	if err == nil {
		t.Error("expected error")
	}

	// cover the no such object path
	delete(s.handlers, serviceInstance)
	_, err = methods.GetCurrentTime(ctx, client)
	if err == nil {
		t.Error("expected error")
	}

	// cover the method not supported path
	res, err := http.Get(ts.URL.String())
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %s", http.StatusMethodNotAllowed, res.Status)
	}

	// cover the ioutil.ReadAll error path
	s.readAll = func(io.Reader) ([]byte, error) {
		return nil, io.ErrShortBuffer
	}
	res, err = http.Post(ts.URL.String(), "none", nil)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %s", http.StatusBadRequest, res.Status)
	}
}
