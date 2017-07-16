// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"path/filepath"
	"strings"
)

const (
	// CopyTo is used to indicate that the desired filter spec is for a CopyTo direction
	CopyTo = true

	// CopyFrom is used to indicate that the desired filter spec is for the CopyFrom direction
	CopyFrom = false
)

// GenerateFilterSpec will populate the appropriate relative Rebase and Strip paths based on the supplied scenarios. Inclusion/Exclusion should be constructed separately. Please also note that any mount that exists before the copy target that is not primary and comes before the primary target will have a bogus filterspec, since it would not be written or read to.
func GenerateFilterSpec(copyPath string, mountPoint string, primaryTarget bool, direction bool) FilterSpec {
	var filter FilterSpec

	// Note I know they are just booleans, if that changes then this statement will not need to.
	if direction == CopyTo {
		filter = generateCopyToFilterSpec(copyPath, mountPoint, primaryTarget)
	} else {
		filter = generateCopyFromFilterSpec(copyPath, mountPoint, primaryTarget)
	}

	filter.Exclusions = make(map[string]struct{})
	filter.Inclusions = make(map[string]struct{})
	return filter
}

func generateCopyFromFilterSpec(copyPath string, mountPoint string, primaryTarget bool) FilterSpec {
	var filter FilterSpec

	//we only need the right most element.
	_, first := filepath.Split(copyPath)

	// primary target was provided so we wil need to split the target and take the right most element for the rebase.
	// then strip set the strip as the target path.
	if primaryTarget {
		filter.RebasePath = removeLeadingSlash(first)
		filter.StripPath = removeLeadingSlash(strings.TrimPrefix(copyPath, mountPoint))
		return filter
	}

	// non primary target was provided. in this case we will need rebase to include the right most member of the target(or "/") joined to the front of the mountPath - the target path. 3
	filter.RebasePath = removeLeadingSlash(filepath.Join(first, strings.TrimPrefix(mountPoint, copyPath)))
	filter.StripPath = ""
	return filter
}

func generateCopyToFilterSpec(copyPath string, mountPoint string, primaryTarget bool) FilterSpec {
	var filter FilterSpec

	// primary target was provided so we will need to rebase header assets for this mount to have the target in front for the write.
	if primaryTarget {
		filter.RebasePath = removeLeadingSlash(strings.TrimPrefix(copyPath, mountPoint))
		filter.StripPath = ""
		return filter
	}

	// non primary target, this implies that the asset header has part of the mount point path in it. We must strip out that part since the non primary target will be mounted and be looking at the world from it's own root "/"
	filter.RebasePath = ""
	filter.StripPath = removeLeadingSlash(strings.TrimPrefix(mountPoint, copyPath))

	return filter
}

// removeLeadingSlash will remove the '/' from in front of a target path
// we use this to ensure relative pathing
func removeLeadingSlash(path string) string {
	return strings.TrimPrefix(path, "/")
}

func AddMountExclusions(currentMount string, filter *FilterSpec, mounts []string) error {
	if filter == nil {
		return fmt.Errorf("filterSpec for (%s) was nil, cannot add exclusions", currentMount)
	}

	for _, mount := range mounts {
		if strings.HasPrefix(mount, currentMount) && currentMount != mount {
			// exclusions are relative to the mount so the leading `/` should be removed unless we decide otherwise.
			filter.Exclusions[removeLeadingSlash(strings.TrimPrefix(mount, currentMount))] = struct{}{}
		}
	}
	return nil
}
