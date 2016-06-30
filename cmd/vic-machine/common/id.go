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

package common

import (
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
)

type VCHID struct {
	// VCH id
	ID string
}

func (i *VCHID) IDFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "id",
			Value:       "",
			Usage:       "The ID of the Virtual Container Host - not supported until vic-machine ls is ready",
			Destination: &i.ID,
		},
	}
}

func (i *VCHID) ProcessID() error {
	if i.ID != "" {
		log.Warnf("ID of Virtual Container Host is not supported until vic-machine ls is ready. For details, please refer github issue #810")
		return cli.NewExitError("Can specify --compute-resource and --name", 1)
	}
	return nil
}
