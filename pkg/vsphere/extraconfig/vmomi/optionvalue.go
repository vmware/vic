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

// Package vmomi is in a separate package to avoid the transitive inclusion of govmomi
// as a fundamental dependency of the main extraconfig
package vmomi

import (
	"fmt"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

// OptionValueMap returns a map from array of OptionValues
func OptionValueMap(src []types.BaseOptionValue) map[string]string {
	// create the key/value store from the extraconfig slice for lookups
	kv := make(map[string]string)
	for i := range src {
		k := src[i].GetOptionValue().Key
		v := src[i].GetOptionValue().Value.(string)
		if v == "<nil>" {
			v = ""
		}
		kv[k] = v
	}
	return kv
}

// OptionValueSource is a convenience method to generate a MapSource source from
// and array of OptionValue's
func OptionValueSource(src []types.BaseOptionValue) extraconfig.DataSource {
	kv := OptionValueMap(src)
	return extraconfig.MapSource(kv)
}

// OptionValueFromMap is a convenience method to convert a map into a BaseOptionValue array
func OptionValueFromMap(data map[string]string) []types.BaseOptionValue {
	if len(data) == 0 {
		return nil
	}

	array := make([]types.BaseOptionValue, len(data))

	i := 0
	for k, v := range data {
		if v == "" {
			v = "<nil>"
		}
		array[i] = &types.OptionValue{Key: k, Value: v}
		i++
	}

	return array
}

// OptionValueArrayToString translates the options array in to a Go formatted structure dump
func OptionValueArrayToString(options []types.BaseOptionValue) string {
	// create the key/value store from the extraconfig slice for lookups
	kv := make(map[string]string)
	for i := range options {
		k := options[i].GetOptionValue().Key
		v := options[i].GetOptionValue().Value.(string)
		kv[k] = v
	}

	return fmt.Sprintf("%#v", kv)
}
