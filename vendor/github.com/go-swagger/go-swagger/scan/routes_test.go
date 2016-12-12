// Copyright 2015 go-swagger maintainers
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

package scan

import (
	goparser "go/parser"
	"log"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestRouteExpression(t *testing.T) {
	assert.Regexp(t, rxRoute, "swagger:route DELETE /orders/{id} deleteOrder")
	assert.Regexp(t, rxRoute, "swagger:route GET /v1.2/something deleteOrder")
}

func TestRoutesParser(t *testing.T) {
	docFile := "../fixtures/goparsing/classification/operations/todo_operation.go"
	fileTree, err := goparser.ParseFile(classificationProg.Fset, docFile, nil, goparser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	rp := newRoutesParser(classificationProg)
	var ops spec.Paths
	err = rp.Parse(fileTree, &ops)
	assert.NoError(t, err)

	assert.Len(t, ops.Paths, 3)

	po, ok := ops.Paths["/pets"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"listPets",
		"Lists pets filtered by some parameters.",
		"This will show all available pets by default.\nYou can get the pets that are out of stock",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)
	assertOperation(t,
		po.Post,
		"createPet",
		"Create a pet based on the parameters.",
		"",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"listOrders",
		"lists orders filtered by some parameters.",
		"",
		[]string{"orders"},
		[]string{"orders:read", "https://www.googleapis.com/auth/userinfo.email"},
	)
	assertOperation(t,
		po.Post,
		"createOrder",
		"create an order based on the parameters.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders/{id}"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"orderDetails",
		"gets the details for an order.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperation(t,
		po.Put,
		"updateOrder",
		"Update the details for an order.",
		"When the order doesn't exist this will return an error.",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperation(t,
		po.Delete,
		"deleteOrder",
		"delete a particular order.",
		"",
		nil,
		[]string{"read", "write"},
	)
}

func TestRoutesParserBody(t *testing.T) {
	docFile := "../fixtures/goparsing/classification/operations_body/todo_operation_body.go"
	fileTree, err := goparser.ParseFile(classificationProg.Fset, docFile, nil, goparser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	rp := newRoutesParser(classificationProg)
	var ops spec.Paths
	err = rp.Parse(fileTree, &ops)
	assert.NoError(t, err)

	assert.Len(t, ops.Paths, 3)

	po, ok := ops.Paths["/pets"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"listPets",
		"Lists pets filtered by some parameters.",
		"This will show all available pets by default.\nYou can get the pets that are out of stock",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)
	assertOperationBody(t,
		po.Post,
		"createPet",
		"Create a pet based on the parameters.",
		"",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"listOrders",
		"lists orders filtered by some parameters.",
		"",
		[]string{"orders"},
		[]string{"orders:read", "https://www.googleapis.com/auth/userinfo.email"},
	)
	assertOperationBody(t,
		po.Post,
		"createOrder",
		"create an order based on the parameters.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders/{id}"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"orderDetails",
		"gets the details for an order.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperationBody(t,
		po.Put,
		"updateOrder",
		"Update the details for an order.",
		"When the order doesn't exist this will return an error.",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperationBody(t,
		po.Delete,
		"deleteOrder",
		"delete a particular order.",
		"",
		nil,
		[]string{"read", "write"},
	)
}

func assertOperation(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string) {
	assert.NotNil(t, op)
	assert.Equal(t, summary, op.Summary)
	assert.Equal(t, description, op.Description)
	assert.Equal(t, id, op.ID)
	assert.EqualValues(t, tags, op.Tags)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Consumes)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Produces)
	assert.EqualValues(t, []string{"http", "https", "ws", "wss"}, op.Schemes)
	assert.Len(t, op.Security, 2)
	akv, ok := op.Security[0]["api_key"]
	assert.True(t, ok)
	// akv must be defined & not empty
	assert.NotNil(t, akv)
	assert.Empty(t, akv)

	vv, ok := op.Security[1]["oauth"]
	assert.True(t, ok)
	assert.EqualValues(t, scopes, vv)

	assert.NotNil(t, op.Responses.Default)
	assert.Equal(t, "#/responses/genericError", op.Responses.Default.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.Equal(t, "#/responses/someResponse", rsp.Ref.String())
	rsp, ok = op.Responses.StatusCodeResponses[422]
	assert.True(t, ok)
	assert.Equal(t, "#/responses/validationError", rsp.Ref.String())
}

func assertOperationBody(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string) {
	assert.NotNil(t, op)
	assert.Equal(t, summary, op.Summary)
	assert.Equal(t, description, op.Description)
	assert.Equal(t, id, op.ID)
	assert.EqualValues(t, tags, op.Tags)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Consumes)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Produces)
	assert.EqualValues(t, []string{"http", "https", "ws", "wss"}, op.Schemes)
	assert.Len(t, op.Security, 2)
	akv, ok := op.Security[0]["api_key"]
	assert.True(t, ok)
	// akv must be defined & not empty
	assert.NotNil(t, akv)
	assert.Empty(t, akv)

	vv, ok := op.Security[1]["oauth"]
	assert.True(t, ok)
	assert.EqualValues(t, scopes, vv)

	assert.NotNil(t, op.Responses.Default)
	assert.Equal(t, "", op.Responses.Default.Ref.String())
	assert.Equal(t, "#/definitions/genericError", op.Responses.Default.Schema.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.Equal(t, "", rsp.Ref.String())
	assert.Equal(t, "#/definitions/someResponse", rsp.Schema.Ref.String())
	rsp, ok = op.Responses.StatusCodeResponses[422]
	assert.True(t, ok)
	assert.Equal(t, "", rsp.Ref.String())
	assert.Equal(t, "#/definitions/validationError", rsp.Schema.Ref.String())
}
