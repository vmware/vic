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

package metadata

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type guestInfoConfig struct {
	// guestinfo *rpcvmx.Config
}

// New generates a handle to a ConfigLoader
func New() ConfigLoader {
	config := guestInfoConfig{}
	// config.guestinfo = rpcvmx.NewConfig()

	return config
}

// LoadConfig will do so from the VMs GuestInfo
func (c guestInfoConfig) LoadConfig(configblob string) (*ExecutorConfig, error) {
	// commented out until we move rpc tools to pure Go - we need a static binary
	// if !vmcheck.IsVirtualWorld() {
	// 	return nil, errors.New("not in a virtual world")
	// }

	// configblob, err := c.guestinfo.String(key, "")
	// if err != nil {
	// 	return nil, err
	// }
	if len(configblob) == 0 {
		return nil, errors.New("failed to retrieve populated config blob from guestinfo." + key)
	}

	config := &ExecutorConfig{}
	// err := gob.NewDecoder(bytes.NewBuffer([]byte(configblob))).Decode(config)
	data, err := base64.StdEncoding.DecodeString(configblob)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (c guestInfoConfig) StoreConfig(config *ExecutorConfig) (string, error) {
	// var metadataBuf bytes.Buffer
	// err := gob.NewEncoder(&metadataBuf).Encode(config)
	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	// return metadataBuf.String(), nil
	return base64.StdEncoding.EncodeToString(data), nil
}
