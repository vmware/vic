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
		{
			&types.RetrieveProperties{
				This: types.ManagedObjectReference{Type: "PropertyCollector", Value: "ha-property-collector"},
				SpecSet: []types.PropertyFilterSpec{
					{
						DynamicData: types.DynamicData{},
						PropSet: []types.PropertySpec{
							{
								DynamicData: types.DynamicData{},
								Type:        "ManagedEntity",
								All:         (*bool)(nil),
								PathSet:     []string{"name", "parent"},
							},
						},
						ObjectSet: []types.ObjectSpec{
							{
								DynamicData: types.DynamicData{},
								Obj:         types.ManagedObjectReference{Type: "Folder", Value: "ha-folder-root"},
								Skip:        types.NewBool(false),
								SelectSet: []types.BaseSelectionSpec{ // test decode of interface
									&types.TraversalSpec{
										SelectionSpec: types.SelectionSpec{
											DynamicData: types.DynamicData{},
											Name:        "traverseParent",
										},
										Type: "ManagedEntity",
										Path: "parent",
										Skip: types.NewBool(false),
										SelectSet: []types.BaseSelectionSpec{
											&types.SelectionSpec{
												DynamicData: types.DynamicData{},
												Name:        "traverseParent",
											},
										},
									},
								},
							},
						},
						ReportMissingObjectsInResults: (*bool)(nil),
					},
				}},
			`<?xml version="1.0" encoding="UTF-8"?>
                         <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
                          <Body>
                           <RetrieveProperties xmlns="urn:vim25">
                            <_this type="PropertyCollector">ha-property-collector</_this>
                            <specSet>
                             <propSet>
                              <type>ManagedEntity</type>
                              <pathSet>name</pathSet>
                              <pathSet>parent</pathSet>
                             </propSet>
                             <objectSet>
                              <obj type="Folder">ha-folder-root</obj>
                              <skip>false</skip>
                              <selectSet xmlns:XMLSchema-instance="http://www.w3.org/2001/XMLSchema-instance" XMLSchema-instance:type="TraversalSpec">
                               <name>traverseParent</name>
                               <type>ManagedEntity</type>
                               <path>parent</path>
                               <skip>false</skip>
                               <selectSet XMLSchema-instance:type="SelectionSpec">
                                <name>traverseParent</name>
                               </selectSet>
                              </selectSet>
                             </objectSet>
                            </specSet>
                           </RetrieveProperties>
                          </Body>
                         </Envelope>`,
		},
	}

	for i, req := range requests {
		method, err := UnmarshalBody([]byte(req.data))
		if err != nil {
			t.Errorf("failed to decode %d (%s): %s", i, req, err)
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
		New(NewServiceInstance(esx.ServiceContent, esx.RootFolder)),
		New(NewServiceInstance(vc.ServiceContent, vc.RootFolder)),
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
			t.Fail()
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

func TestServeHTTPErrors(t *testing.T) {
	s := New(NewServiceInstance(esx.ServiceContent, esx.RootFolder))

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
	Map.objects[serviceInstance] = &errorNoSuchMethod{}
	_, err = methods.GetCurrentTime(ctx, client)
	if err == nil {
		t.Error("expected error")
	}

	// cover the xml encode error path
	Map.objects[serviceInstance] = &errorMarshal{}
	_, err = methods.GetCurrentTime(ctx, client)
	if err == nil {
		t.Error("expected error")
	}

	// cover the no such object path
	Map.Remove(serviceInstance)
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
