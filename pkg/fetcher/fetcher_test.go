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

package fetcher

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ErrJsonStr401 = `{
	"errors":
		[{"code":"UNAUTHORIZED",
		  "message":"authentication required",
		  "detail":[{"Type":"repository","Class":"","Name":"library/jiowengew","Action":"pull"}]
		}]
	}
	`
	MultipleErrJsonStr = `{
	"errors":
		[{"code":"UNAUTHORIZED",
		  "message":"authentication required",
		  "detail":[{"Type":"repository","Class":"","Name":"library/jiowengew","Action":"pull"}]
		},
		{"code":"NOTFOUND",
		 "message": "image not found",
		 "detail": "image not found"
		}]
	}
	`
	RandomStr                    = `random`
	RandomJsonStr                = `{"nope":"nope"}`
	ErrJsonWithEmptyErrorsField  = `{"errors":[]}`
	ErrJsonWithNoMessageField    = `{"errors":[{"code":"nope","detail":"nope"}]}`
	ErrJsonWithEmptyMessageField = `{"errors":[{"code":"nope","message":""},{"message":""}]}`
)

func TestExtractErrResponseMessage(t *testing.T) {
	// Test set up: create the io streams for testing purposes
	// multiple streams needed: these streams only have read ends
	singleErrTestStream := ioutil.NopCloser(bytes.NewReader([]byte(ErrJsonStr401)))
	multipleErrTestStream := ioutil.NopCloser(bytes.NewReader([]byte(MultipleErrJsonStr)))
	randomStrTestStream := ioutil.NopCloser(bytes.NewReader([]byte(RandomStr)))
	malformedJSONTestStream := ioutil.NopCloser(bytes.NewReader([]byte(RandomJsonStr)))
	emptyErrorsJSONTestStream := ioutil.NopCloser(bytes.NewReader([]byte(ErrJsonWithEmptyErrorsField)))
	noMessageJSONTestStream := ioutil.NopCloser(bytes.NewReader([]byte(ErrJsonWithNoMessageField)))
	emptyMessageJSONTestStream := ioutil.NopCloser(bytes.NewReader([]byte(ErrJsonWithEmptyMessageField)))

	// Test 1: single error message extraction
	msg, err := extractErrResponseMessage(singleErrTestStream)
	assert.Nil(t, err, "test: (single error message) extraction should success for well-formatted error json")
	assert.NotNil(t, msg, "test: (single error message) message extracted for well-formatted error json")
	assert.Equal(t, "authentication required", msg,
		"test: (single error message) extracted message: %s; expected: authentication required", msg)

	// Test 2: multiple error message extraction
	msg, err = extractErrResponseMessage(multipleErrTestStream)
	assert.Nil(t, err, "test: (multiple error messages) extraction should success for well-formatted error json")
	assert.NotNil(t, msg, "test: (multiple error messages) message extracted for well-formatted error json")
	assert.Equal(t, "authentication required, image not found", msg,
		"test: (multiple error messages) extracted message: %s; expected: authentication required, image not found", msg)

	// Test 3: random string in the stream that is not a json
	msg, err = extractErrResponseMessage(randomStrTestStream)
	assert.Equal(t, "", msg, "test: (non-json string) no message should be extracted")
	assert.NotNil(t, err, "test: (non-json string) extraction should fail")

	// Test 4: malformed json string
	msg, err = extractErrResponseMessage(malformedJSONTestStream)
	assert.Equal(t, "", msg, "test: (malformed json string) no message should be extracted")
	assert.NotNil(t, err, "test: (malformed json string) extraction should fail")
	assert.Equal(t, "error response json has unconventional format", err.Error(),
		"test: (malformed json string) error: error response json has unconventional format; expected error: %s", err.Error())

	// Test 5: malformed json with empty `errors` field
	msg, err = extractErrResponseMessage(emptyErrorsJSONTestStream)
	assert.Equal(t, "", msg, "test: (malformed json string, empty errors field) no message should be extracted")
	assert.NotNil(t, err, "test: (malformed json string, empty errors field) extraction should fail")
	assert.Equal(t, "error response json has unconventional format", err.Error(),
		"test: (malformed json string, empty errors field) error: %s; expected error: error response json has unconventional format", err.Error())

	// Test 6: malformed json with no `message` field
	msg, err = extractErrResponseMessage(noMessageJSONTestStream)
	assert.Equal(t, "", msg, "test: (malformed json string, no message field) no message should be extracted")
	assert.NotNil(t, err, "test: (malformed json string, no message field) extraction should fail")
	assert.Equal(t, "error response json has unconventional format", err.Error(),
		"test: (malformed json string, no message field) error: %s; expected error: error response json has unconventional format", err.Error())

	// Test 7: malformed json with empty string in `message` field
	msg, err = extractErrResponseMessage(emptyMessageJSONTestStream)
	assert.Equal(t, "", msg, "test: (malformed json string, empty message field) no message should be extracted")
	assert.NotNil(t, err, "test: (malformed json string, empty message field) extraction should fail")
	assert.Equal(t, "error response json has unconventional format", err.Error(),
		"test: (malformed json string, empty message field) error: %s; expected error: error response json has unconventional format", err.Error())
}
