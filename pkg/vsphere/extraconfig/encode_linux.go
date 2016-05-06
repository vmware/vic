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

// GuestInfoSink uses the rpcvmx mechanism to update the guestinfo key/value map as
// the datasink for encoding target structures
func GuestInfoSink() (DataSink, error) {
	guestinfo := rpcvmx.NewConfig()

	if !vmcheck.IsVirtualWorld() {
		return nil, errors.New("not in a virtual world")
	}

	return func(key, value string) error {
		if value == "" {
			value = "<nil>"
		}

		log.Debugf("GuestInfoSink: setting key: %s, value: %#v", key, value)
		err := guestinfo.SetString(key, value)
		return err
	}, nil
}
