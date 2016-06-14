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

package extraconfig

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vmw-guestinfo/rpcvmx"
	"github.com/vmware/vmw-guestinfo/vmcheck"
)

// GuestInfoSource uses the rpcvmx mechanism to access the guestinfo key/value map as
// the datasource for decoding into target structures
func GuestInfoSource() (DataSource, error) {
	return GuestInfoSourceWithPrefix("")
}

// GuestInfoSourceWithPrefix adds a prefix to all keys accessed. The key must not have leading
// or trailing separator characters, but may have separators in other positions. The separator
// (either . or /) will be replaced with the appropriate value for the key in question.
func GuestInfoSourceWithPrefix(prefix string) (DataSource, error) {
	guestinfo := rpcvmx.NewConfig()

	if !vmcheck.IsVirtualWorld() {
		return nil, errors.New("not in a virtual world")
	}

	return func(key string) (string, error) {
		key = addPrefixToKey(DefaultGuestInfoPrefix, prefix, key)

		value, err := guestinfo.String(key, "")
		if value == "" {
			// if we don't assume this then we never return error for missing keys
			// which is problematic for decoding
			err = errors.New("no value for key")
		} else if value == "<nil>" {
			value = ""
		}

		log.Debugf("GuestInfoSource: key: %s, value: %#v, error: %s", key, value, err)
		return value, err
	}, nil
}
