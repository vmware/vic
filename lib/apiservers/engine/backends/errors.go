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

package backends

import (
	"fmt"
	derr "github.com/docker/docker/errors"
	"net/http"
)

// InvalidVolumeError is returned when the user specifies a client directory as a volume.
type InvalidVolumeError struct {
}

func (e InvalidVolumeError) Error() string {
	return fmt.Sprintf("%s does not support mounting directories as a data volume.", ProductName())
}

// InvalidBindError is returned when create/run -v has more params than allowed.
type InvalidBindError struct {
	volume string
}

func (e InvalidBindError) Error() string {
	return fmt.Sprintf("volume bind input is invalid: -v %s", e.volume)
}

// VolumeJoinNotFoundError returns a 404 docker error for a volume join request.
func VolumeJoinNotFoundError(msg string) error {
	return derr.NewRequestNotFoundError(fmt.Errorf(msg))
}

// VolumeCreateNotFoundError returns a 404 docker error for a volume create request.
func VolumeCreateNotFoundError(msg string) error {
	return derr.NewErrorWithStatusCode(fmt.Errorf("No volume store named (%s) exists", msg), http.StatusInternalServerError)
}
