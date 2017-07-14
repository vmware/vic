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
	"path/filepath"
	"strings"
)

const (
	// CopyTo is used to indicate that the desired filter spec is for a CopyTo direction
	CopyTo = true

	// CopyFrom is used to indicate that the desired filter spec is for the CopyFrom direction
	CopyFrom = false
)

// GenerateFilterSpec will populate the appropriate relative Rebase and Strip paths based on the supplied scenarios. Inclusion/Exclusion should be constructed separately.
func GenerateFilterSpec(targetPath string, mountPoint string, primaryTarget bool, direction bool) FilterSpec {
	var filter FilterSpec

	// Note I know they are just booleans, if that changes then this statement will not need to.
	if direction == CopyTo {
		filter = generateCopyToFilterSpec(targetPath, mountPoint, primaryTarget)
	} else {
		filter = generateCopyFromFilterSpec(targetPath, mountPoint, primaryTarget)
	}

	filter.Exclusions = make(map[string]struct{})
	filter.Inclusions = make(map[string]struct{})
	return filter
}

func generateCopyFromFilterSpec(targetPath string, mountPoint string, primaryTarget bool) FilterSpec {
	var filter FilterSpec

	//we only need the right most element.
	_, first := filepath.Split(targetPath)

	// primary target was provided so we wil need to split the target and take the right most element for the rebase.
	// then strip set the strip as the target path.
	if primaryTarget {
		filter.RebasePath = first
		filter.StripPath = removeLeadingSlash(targetPath)
		return filter
	}

	// non primary target was provided. in this case we will need rebase to include the right most member of the target(or "/") joined to the front of the mountPath - the target path. 3
	filter.RebasePath = removeLeadingSlash(filepath.Join(first, strings.TrimPrefix(mountPoint, targetPath)))
	filter.StripPath = ""
	return filter
}

func generateCopyToFilterSpec(targetPath string, mountPoint string, primaryTarget bool) FilterSpec {
	var filter FilterSpec

	// primary target was provided so we will need to rebase header assets for this mount to have the target in front for the write.
	if primaryTarget {
		filter.RebasePath = removeLeadingSlash(targetPath)
		filter.StripPath = ""
		return filter
	}

	// non primary target, this implies that the asset header has part of the mount point path in it. We must strip out that part since the non primary target will be mounted and be looking at the world from it's own root "/"
	filter.RebasePath = ""
	filter.StripPath = removeLeadingSlash(strings.TrimPrefix(mountPoint, targetPath))

	return filter
}

// removeLeadingSlash will remove the '/' form in front of a target path if it is not "/"
// we use this to ensure relative pathing, except for when we assign '/'
func removeLeadingSlash(path string) string {
	if strings.HasPrefix(path, "/") && path != "/" {
		return strings.TrimPrefix(path, "/")
	}
	return path
}
