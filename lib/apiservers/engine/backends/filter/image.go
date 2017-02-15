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

package filter

import (
	"fmt"
	"path"

	"github.com/docker/docker/api/types/filters"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
)

/*
* ValidateImageFilters will validate the image filters are
* valid docker filters / values and supported by vic.
*
* The function will reuse dockers filter validation
*
 */
func ValidateImageFilters(cmdFilters filters.Args, acceptedFilters map[string]bool, unSupportedFilters map[string]bool) (*FilterContext, error) {

	// ensure filter options are valid and supported by vic
	if err := ValidateFilters(cmdFilters, acceptedFilters, unSupportedFilters); err != nil {
		return nil, err
	}

	// return value
	imgFilterContext := &FilterContext{}

	err := cmdFilters.WalkValues("before", func(value string) error {
		before, err := cache.ImageCache().Get(value)
		if before == nil {
			err = fmt.Errorf("No such image: %s", value)
		} else {
			imgFilterContext.BeforeID = &before.ImageID
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	err = cmdFilters.WalkValues("since", func(value string) error {
		since, err := cache.ImageCache().Get(value)
		if since == nil {
			err = fmt.Errorf("No such image: %s", value)
		} else {
			imgFilterContext.SinceID = &since.ImageID
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return imgFilterContext, nil

}

/*
*	IncludeImage will evaluate the filter criteria in filterContext against the provided
* 	image and determine what action to take.  There are three options:
*		* IncludeAction
*		* ExcludeAction
*		* StopAction
*
 */
func IncludeImage(imgFilters filters.Args, listContext *FilterContext) FilterAction {

	// filter common requirements
	act := filterCommon(listContext, imgFilters)
	if act != IncludeAction {
		return act
	}

	// filter on image reference
	if imgFilters.Include("reference") {
		// references for this imageID
		refs := cache.RepositoryCache().References(listContext.ID)
		for _, ref := range refs {
			err := imgFilters.WalkValues("reference", func(value string) error {
				// match on complete ref ie. busybox:latest
				matchRef, _ := path.Match(value, ref.String())
				// match on repo only ie. busybox
				matchName, _ := path.Match(value, ref.Name())
				if !matchRef && !matchName {
					return fmt.Errorf("reference not found")
				}
				return nil
			})
			if err != nil {
				return ExcludeAction
			}
		}
	}
	return IncludeAction
}
