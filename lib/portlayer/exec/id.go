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

package exec

import "github.com/docker/docker/pkg/stringid"

// ID is a container's unique id
type ID string

// NilID is a placeholder for an empty container ID
const NilID ID = ID("")

// ParseID converts a string to ID
func ParseID(id string) ID {
	return ID(id)
}

// Truncate returns the truncated ID
func (id ID) TruncateID() ID {
	return ID(stringid.TruncateID(string(id)))
}

func (id ID) String() string {
	return string(id)
}

// GenerateID generates a new container ID
func GenerateID() ID {
	return ParseID(stringid.GenerateNonCryptoID())
}
