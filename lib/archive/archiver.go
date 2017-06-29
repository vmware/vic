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

package archive

import (
	"fmt"
	"io"
	"net/url"

	"github.com/vmware/vic/pkg/trace"
)

// FilterType specifies what type of filter to apply in a FilterSpec
type FilterType int

const (
	// Include specifies a path inclusion.
	// Note: An included path that sits under an excluded path should be included.
	// Example: if / is excluded, and /files is included, /files should  be added to the archive.
	Include = iota
	// Exclude specifies a path exclusion.
	Exclude
	// Rebase specifies a path rebase.
	// The rebase path is prepended to the path in the archive header.
	Rebase
	// Strip specifies a path strip.
	// The inverse of a rebase, the path is stripped from the header path
	// before writing to disk.
	Strip
)

// FilterSpec describes rules for handling specified paths during archival
type FilterSpec struct {
	Inclusions map[string]struct{}
	Exclusions map[string]struct{}
	RebasePath string
	StripPath  string
}

// Archiver defines an API for creating archives consisting of data that
// exist between a child and parent layer, as well as unpacking archives
// to container filesystems
type Archiver interface {

	// Export reads the delta between child and ancestor layers, returning
	// the difference as a tar archive.
	//
	// store - the store containing the two layers
	// id - must inherit from ancestor if ancestor is specified
	// ancestor - the old layer against which to diff. If omitted, the entire filesystem will be included
	// spec - describes filters on paths found in the data (include, exclude, rebase, strip)
	// data - include file data in the tar archive if true, headers only otherwise
	Export(op trace.Operation, store *url.URL, id, ancestor string, spec *FilterSpec, data bool) (io.ReadCloser, error)
}

// CreateFilterSpec creates a FilterSpec from a supplied map
func CreateFilterSpec(op trace.Operation, spec map[string]FilterType) (*FilterSpec, error) {
	var rebase, strip string

	fs := &FilterSpec{
		Inclusions: make(map[string]struct{}),
		Exclusions: make(map[string]struct{}),
	}

	if spec == nil {
		return fs, nil
	}

	for k, v := range spec {
		switch v {
		case Include:
			fs.Inclusions[k] = struct{}{}
		case Exclude:
			fs.Exclusions[k] = struct{}{}
		case Rebase:
			if rebase != "" {
				return nil, fmt.Errorf("Error creating filter spec: only one rebase path allowed")
			}
			rebase = k
		case Strip:
			if strip != "" {
				return nil, fmt.Errorf("Error creating filter spec: only one strip path allowed")
			}
			strip = k
		default:
			return nil, fmt.Errorf("Invalid filter specification: %d", v)
		}
	}

	fs.RebasePath = rebase
	fs.StripPath = strip

	return fs, nil
}
