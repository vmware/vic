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

package main

import (
	"errors"
	"io"
	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/dio"
)

type operations struct {
	tether.BaseOperations

	logging bool
}

func (t *operations) Log() (io.Writer, error) {
	return nil, errors.New("not implemented on OSX")
}

// sessionLogWriter returns a writer that will persist the session output
func (t *operations) SessionLog(session *tether.SessionConfig) (dio.DynamicMultiWriter, error) {
	return nil, errors.New("not implemented on OSX")
}

func (t *operations) Setup(sink tether.Config) error {

	if err := t.BaseOperations.Setup(sink); err != nil {
		return err
	}

	// TODO: enabled for initial dev debugging only
	log.Info("Launching pprof server on port 6060")
	go func() {
		log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	return nil
}
