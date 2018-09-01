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

package handlers

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/list"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/cmd/vic-machine/create"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/trace"
)

type mockFinder struct {
	mock.Mock

	path string
}

func (mf *mockFinder) Element(ctx context.Context, t types.ManagedObjectReference) (*list.Element, error) {
	args := mf.Called(ctx, t)

	if p := args.Get(0); p != nil {
		return &list.Element{
			Path: p.(string),
		}, args.Error(1)
	}

	return nil, args.Error(1)
}

func TestFromManagedObject(t *testing.T) {
	op := trace.NewOperation(context.Background(), "TestFromManagedObject")
	var m *models.ManagedObject

	expectedType := "t"

	expected := ""
	actual, _, err := fromManagedObject(op, nil, m, expectedType)
	assert.NoError(t, err, "Expected nil error, got %#v", err)
	assert.Equal(t, expected, actual)

	name := "testManagedObject"
	path := "/folder/" + name

	m = &models.ManagedObject{
		Name: name,
	}

	mf := &mockFinder{}
	mf.On("Element", op, mock.AnythingOfType("types.ManagedObjectReference")).Return(path, nil)

	expected = name
	actual, _, err = fromManagedObject(op, mf, m, expectedType)
	assert.NoError(t, err, "Expected nil error, got %#v", err)
	assert.Equal(t, expected, actual)

	m.ID = "testID"

	expected = path
	actual, actualType, err := fromManagedObject(op, mf, m, expectedType)
	assert.NoError(t, err, "Expected nil error, got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedType, actualType)

	m.Name = ""

	expected = path
	actual, actualType, err = fromManagedObject(op, mf, m, expectedType)
	assert.NoError(t, err, "Expected nil error, got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedType, actualType)
}

func TestFromManagedObject_Negative(t *testing.T) {
	op := trace.NewOperation(context.Background(), "TestFromManagedObject")

	m := &models.ManagedObject{
		ID: "testID",
	}

	mf := &mockFinder{}
	mf.On("Element", op, mock.AnythingOfType("types.ManagedObjectReference")).Return(nil, nil)

	expected := ""
	actual, _, err := fromManagedObject(op, mf, m, "t")
	assert.Error(t, err, "Expected error when no resource found")
	assert.Equal(t, expected, actual)

	expected = ""
	actual, _, err = fromManagedObject(op, mf, m, "t", "u", "v")
	assert.Error(t, err, "Expected error when no resource found")
	assert.Equal(t, expected, actual)
}

func TestFromManagedObject_Fallback(t *testing.T) {
	op := trace.NewOperation(context.Background(), "TestFromManagedObject")

	m := &models.ManagedObject{
		ID: "testID",
	}

	mf := &mockFinder{}
	mf.On("Element", op, mock.MatchedBy(func(t types.ManagedObjectReference) bool { return t.Type == "t" })).Return(nil, &find.NotFoundError{})
	mf.On("Element", op, mock.MatchedBy(func(t types.ManagedObjectReference) bool { return t.Type == "u" })).Return(nil, nil)
	mf.On("Element", op, mock.MatchedBy(func(t types.ManagedObjectReference) bool { return t.Type == "v" })).Return("Result", nil)
	mf.On("Element", op, mock.MatchedBy(func(t types.ManagedObjectReference) bool { return t.Type == "e" })).Return(nil, fmt.Errorf("Expected"))

	expected := ""
	actual, actualType, err := fromManagedObject(op, mf, m, "t")
	assert.Error(t, err, "Expected error when NotFoundError encountered")
	assert.Equal(t, expected, actual)
	assert.Equal(t, "", actualType)

	expected = ""
	actual, actualType, err = fromManagedObject(op, mf, m, "t", "u")
	assert.Error(t, err, "Expected error when no resource found")
	assert.Equal(t, expected, actual)
	assert.Equal(t, "", actualType)

	expected = "Result"
	actual, actualType, err = fromManagedObject(op, mf, m, "t", "u", "v")
	assert.NoError(t, err, "Did not expect error when third type returns valid result, but got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, "v", actualType)

	expected = "Result"
	actual, actualType, err = fromManagedObject(op, mf, m, "t", "v")
	assert.NoError(t, err, "Did not expect error when second type returns valid result, but got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, "v", actualType)

	expected = "Result"
	actual, actualType, err = fromManagedObject(op, mf, m, "u", "v")
	assert.NoError(t, err, "Did not expect error when second type returns valid result, but got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, "v", actualType)

	expected = "Result"
	actual, actualType, err = fromManagedObject(op, mf, m, "u", "e", "v")
	assert.NoError(t, err, "Did not expect error when third type returns valid result, but got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, "v", actualType)
}

func TestFromCIDR(t *testing.T) {
	var m models.CIDR

	expected := ""
	actual := fromCIDR(&m)
	assert.Equal(t, expected, actual)

	m = "10.10.1.0/32"

	expected = string(m)
	actual = fromCIDR(&m)
	assert.Equal(t, expected, actual)
}

func TestFromGateway(t *testing.T) {
	var m *models.Gateway

	expected := ""
	actual := fromGateway(m)
	assert.Equal(t, expected, actual)

	m = &models.Gateway{
		Address: "192.168.31.37",
		RoutingDestinations: []models.CIDR{
			"192.168.1.1/24",
			"172.17.0.1/24",
		},
	}

	expected = "192.168.1.1/24,172.17.0.1/24:192.168.31.37"
	actual = fromGateway(m)
	assert.Equal(t, expected, actual)
}

func TestCreateVCH(t *testing.T) {
	vch := &models.VCH{
		Name:  "test-vch",
		Debug: 3,
		Compute: &models.VCHCompute{
			Resource: &models.ManagedObject{
				Name: "TestCluster",
			},
		},
		Storage: &models.VCHStorage{
			ImageStores: []string{
				"ds://test/datastore",
			},
		},
		Network: &models.VCHNetwork{
			Bridge: &models.VCHNetworkBridge{
				IPRange: "17.16.0.0/12",
				PortGroup: &models.ManagedObject{
					Name: "bridge",
				},
			},
			Public: &models.Network{
				PortGroup: &models.ManagedObject{
					Name: "public",
				},
			},
		},
		Registry: &models.VCHRegistry{
			ImageFetchProxy: &models.VCHRegistryImageFetchProxy{
				HTTP:    "http://example.com",
				HTTPS:   "https://example.com",
				NoProxy: []strfmt.URI{"localhost", ".example.com"},
			},
			Insecure: []string{
				"https://insecure.example.com",
			},
			Whitelist: []string{
				"10.0.0.0/8",
			},
		},
		Auth: &models.VCHAuth{
			Server: &models.VCHAuthServer{
				Generate: &models.VCHAuthServerGenerate{
					Cname: "vch.example.com",
					Organization: []string{
						"VMware, Inc.",
					},
					Size: &models.ValueBits{
						Value: models.Value{Value: 2048},
						Units: "bits",
					},
				},
			},
		},
		SyslogAddr: "tcp://syslog.example.com:4444",
		Container: &models.VCHContainer{
			NameConvention: "container-{id}",
		},
	}

	op := trace.NewOperation(context.Background(), "testing")
	defer func() {
		err := os.RemoveAll("test-vch")
		assert.NoError(t, err, "Error removing temp directory: %s", err)
	}()

	pass := "testpass"
	data := &data.Data{
		Target: &common.Target{
			URL:      &url.URL{Host: "10.10.1.2"},
			User:     "testuser",
			Password: &pass,
		},
	}

	ca := newCreate()
	ca.Data = data
	ca.DisplayName = "test-vch"
	err := ca.ProcessParams(op)
	assert.NoError(t, err, "Error while processing params: %s", err)
	op.Infof("ca EnvFile: %s", ca.Certs.EnvFile)

	mf := &mockFinder{}
	mf.On("Element", op, mock.AnythingOfType("types.ManagedObjectReference")).Return(nil, nil)

	cb, err := buildCreate(op, data, mf, vch)
	assert.NoError(t, err, "Error while processing params: %s", err)

	a := reflect.ValueOf(ca).Elem()
	b := reflect.ValueOf(cb).Elem()

	if err = compare(a, b, 0); err != nil {
		t.Fatalf("Error comparing create structs: %s", err)
	}
}

func newCreate() *create.Create {
	debug := 3
	ca := create.NewCreate()
	ca.Debug = common.Debug{Debug: &debug}
	ca.Compute = common.Compute{DisplayName: "TestCluster"}
	ca.ImageDatastorePath = "ds://test/datastore"
	ca.BridgeIPRange = "17.16.0.0/12"
	ca.BridgeNetworkName = "bridge"
	ca.PublicNetworkName = "public"
	ca.Certs.Cname = "vch.example.com"
	ca.Certs.Org = cli.StringSlice{"VMware, Inc."}
	ca.Certs.KeySize = 2048
	httpProxy := "http://example.com"
	httpsProxy := "https://example.com"
	noProxy := "localhost,.example.com"
	ca.Proxies = common.Proxies{
		HTTPProxy:  &httpProxy,
		HTTPSProxy: &httpsProxy,
		NoProxy:    &noProxy,
	}
	ca.Registries = common.Registries{
		InsecureRegistriesArg:  cli.StringSlice{"https://insecure.example.com"},
		WhitelistRegistriesArg: cli.StringSlice{"10.0.0.0/8"},
	}
	ca.SyslogAddr = "tcp://syslog.example.com:4444"
	ca.ContainerNameConvention = "container-{id}"
	ca.Certs.CertPath = "test-vch"
	ca.Certs.NoSaveToDisk = true

	return ca
}

func compare(a, b reflect.Value, index int) (err error) {
	switch a.Kind() {
	case reflect.Invalid, reflect.Uint8: // skip uint8 as generated cert data is not expected to match
		// NOP
	case reflect.Ptr:
		ae := a.Elem()
		be := b.Elem()
		if !ae.IsValid() != !be.IsValid() {
			return fmt.Errorf("Expected pointer validity to match for for %s", a.Type().Name())
		}
		return compare(ae, be, index)
	case reflect.Interface:
		return compare(a.Elem(), b.Elem(), index)
	case reflect.Struct:
		for i := 0; i < a.NumField(); i++ {
			if err = compare(a.Field(i), b.Field(i), i); err != nil {
				fmt.Printf("Field name a: %s, b: %s, index: %d\n", a.Type().Field(i).Name, b.Type().Field(i).Name, i)
				return err
			}
		}
	case reflect.Slice:
		m := min(a.Len(), b.Len())
		for i := 0; i < m; i++ {
			if err = compare(a.Index(i), b.Index(i), i); err != nil {
				return err
			}
		}
	case reflect.Map:
		keys := []string{}
		for _, key := range a.MapKeys() {
			keys = append(keys, key.String())
		}
		sort.Strings(keys)
		for i, key := range keys {
			if err = compare(a.MapIndex(reflect.ValueOf(key)), b.MapIndex(reflect.ValueOf(key)), i); err != nil {
				return err
			}
		}
	case reflect.String:
		if a.String() != b.String() {
			return fmt.Errorf("String fields not equal: %s != %s", a.String(), b.String())
		}
	default:
		if a.CanInterface() && b.CanInterface() {
			if a.Interface() != b.Interface() {
				return fmt.Errorf("Elements are not equal: %#v != %#v", a.Interface(), b.Interface())
			}
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
