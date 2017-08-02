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

package proxy

import (
	"fmt"
	"net/http"

	derr "github.com/docker/docker/api/errors"
)

// InternalServerError returns a 500 docker error on a portlayer error.
func InternalServerError(msg string) error {
	return derr.NewErrorWithStatusCode(fmt.Errorf("Server error from portlayer: %s", msg), http.StatusInternalServerError)
}

// ResourceLockedError returns a 423 http status
func ResourceLockedError(msg string) error {
	return derr.NewErrorWithStatusCode(fmt.Errorf("Resource locked: %s", msg), http.StatusLocked)
}

// ResourceNotFoundError returns a 404 http status
func ResourceNotFoundError(msg string) error {
	return derr.NewErrorWithStatusCode(fmt.Errorf("No such %s", msg), http.StatusNotFound)
}
