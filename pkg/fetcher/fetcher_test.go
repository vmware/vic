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

	"fmt"
	"github.com/stretchr/testify/assert"
)

const (
	errJsonStr401 = `{
	"errors":
		[{"code":"UNAUTHORIZED",
		  "message":"authentication required",
		  "detail":[{"Type":"repository","Class":"","Name":"library/jiowengew","Action":"pull"}]
		}]
	}
	`
	multipleErrJsonStr = `{
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
	unexpectedStr                = `random`
	unexpectedJsonStr            = `{"nope":"nope"}`
	errJsonWithEmptyErrorsField  = `{"errors":[]}`
	errJsonWithNoMessageField    = `{"errors":[{"code":"nope","detail":"nope"}]}`
	errJsonWithEmptyMessageField = `{"errors":[{"code":"nope","message":""},{"message":""}]}`
)

func TestExtractErrResponseMessage(t *testing.T) {
	// Test set up: create the io streams for testing purposes
	// multiple streams needed: these streams only have read ends
	singleErrTestStream := ioutil.NopCloser(bytes.NewReader([]byte(errJsonStr401)))
	multipleErrTestStream := ioutil.NopCloser(bytes.NewReader([]byte(multipleErrJsonStr)))
	unexpectedStrTestStream := ioutil.NopCloser(bytes.NewReader([]byte(unexpectedStr)))
	malformedJSONTestStream := ioutil.NopCloser(bytes.NewReader([]byte(unexpectedJsonStr)))
	emptyErrorsJSONTestStream := ioutil.NopCloser(bytes.NewReader([]byte(errJsonWithEmptyErrorsField)))
	noMessageJSONTestStream := ioutil.NopCloser(bytes.NewReader([]byte(errJsonWithNoMessageField)))
	emptyMessageJSONTestStream := ioutil.NopCloser(bytes.NewReader([]byte(errJsonWithEmptyMessageField)))

	// expected error message upon malformed json to check against
	malformedJsonErrFormat := fmt.Errorf("error response json has unconventional format")

	// Test 1: single error message extraction
	msg, err := extractErrResponseMessage(singleErrTestStream)
	assert.Nil(t, err, "test: (single error message) extraction should success for well-formatted error json")
	assert.Equal(t, "authentication required", msg,
		"test: (single error message) extracted message: %s; expected: authentication required", msg)

	// Test 2: multiple error message extraction
	msg, err = extractErrResponseMessage(multipleErrTestStream)
	assert.Nil(t, err, "test: (multiple error messages) extraction should success for well-formatted error json")
	assert.Equal(t, "authentication required, image not found", msg,
		"test: (multiple error messages) extracted message: %s; expected: authentication required, image not found", msg)

	// Test 3: random string in the stream that is not a json
	msg, err = extractErrResponseMessage(unexpectedStrTestStream)
	assert.Equal(t, "", msg, "test: (non-json string) no message should be extracted")
	assert.NotNil(t, err, "test: (non-json string) extraction should fail")

	// Test 4: malformed json string
	msg, err = extractErrResponseMessage(malformedJSONTestStream)
	assert.Equal(t, "", msg, "test: (malformed json string) no message should be extracted")
	assert.Equal(t, malformedJsonErrFormat, err,
		"test: (malformed json string) error: %s; expected error: %s", err, malformedJsonErrFormat)

	// Test 5: malformed json with empty `errors` field
	msg, err = extractErrResponseMessage(emptyErrorsJSONTestStream)
	assert.Equal(t, "", msg, "test: (malformed json string, empty errors field) no message should be extracted")
	assert.Equal(t, malformedJsonErrFormat, err,
		"test: (malformed json string, empty errors field) error: %s; expected error: %s", err, malformedJsonErrFormat)

	// Test 6: malformed json with no `message` field
	msg, err = extractErrResponseMessage(noMessageJSONTestStream)
	assert.Equal(t, "", msg, "test: (malformed json string, no message field) no message should be extracted")
	assert.Equal(t, malformedJsonErrFormat, err,
		"test: (malformed json string, no message field) error: %s; expected error: %s", err, malformedJsonErrFormat)

	// Test 7: malformed json with empty string in `message` field
	msg, err = extractErrResponseMessage(emptyMessageJSONTestStream)
	assert.Equal(t, "", msg, "test: (malformed json string, empty message field) no message should be extracted")
	assert.Equal(t, malformedJsonErrFormat, err,
		"test: (malformed json string, empty message field) error: %s; expected error: %s", err, malformedJsonErrFormat)
}
