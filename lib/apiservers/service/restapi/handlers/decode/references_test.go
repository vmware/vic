// Copyright 2018 VMware, Inc. All Rights Reserved.
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

package decode

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/list"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/lib/apiservers/service/models"
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
	actual, _, err := FromManagedObject(op, nil, m, expectedType)
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
	actual, _, err = FromManagedObject(op, mf, m, expectedType)
	assert.NoError(t, err, "Expected nil error, got %#v", err)
	assert.Equal(t, expected, actual)

	m.ID = "testID"

	expected = path
	actual, actualType, err := FromManagedObject(op, mf, m, expectedType)
	assert.NoError(t, err, "Expected nil error, got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedType, actualType)

	m.Name = ""

	expected = path
	actual, actualType, err = FromManagedObject(op, mf, m, expectedType)
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
	actual, _, err := FromManagedObject(op, mf, m, "t")
	assert.Error(t, err, "Expected error when no resource found")
	assert.Equal(t, expected, actual)

	expected = ""
	actual, _, err = FromManagedObject(op, mf, m, "t", "u", "v")
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
	actual, actualType, err := FromManagedObject(op, mf, m, "t")
	assert.Error(t, err, "Expected error when NotFoundError encountered")
	assert.Equal(t, expected, actual)
	assert.Equal(t, "", actualType)

	expected = ""
	actual, actualType, err = FromManagedObject(op, mf, m, "t", "u")
	assert.Error(t, err, "Expected error when no resource found")
	assert.Equal(t, expected, actual)
	assert.Equal(t, "", actualType)

	expected = "Result"
	actual, actualType, err = FromManagedObject(op, mf, m, "t", "u", "v")
	assert.NoError(t, err, "Did not expect error when third type returns valid result, but got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, "v", actualType)

	expected = "Result"
	actual, actualType, err = FromManagedObject(op, mf, m, "t", "v")
	assert.NoError(t, err, "Did not expect error when second type returns valid result, but got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, "v", actualType)

	expected = "Result"
	actual, actualType, err = FromManagedObject(op, mf, m, "u", "v")
	assert.NoError(t, err, "Did not expect error when second type returns valid result, but got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, "v", actualType)

	expected = "Result"
	actual, actualType, err = FromManagedObject(op, mf, m, "u", "e", "v")
	assert.NoError(t, err, "Did not expect error when third type returns valid result, but got %#v", err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, "v", actualType)
}
