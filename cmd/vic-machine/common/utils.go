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
	log "github.com/Sirupsen/logrus"

	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/pkg/errors"
)

func LogErrorIfAny(clic *cli.Context, err error) error {
	if err == nil {
		return nil
	}

	log.Errorf("--------------------")
	log.Errorf("%s %s failed: %s\n", clic.App.Name, clic.Command.Name, errors.ErrorStack(err))
	return cli.NewExitError("", 1)
}
