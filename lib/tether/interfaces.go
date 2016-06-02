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

package tether

import (
	"io"

	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/dio"
	"golang.org/x/net/context"
)

type Operations interface {
	Setup() error
	Cleanup() error
	// Log returns the tether debug log writer
	Log() (io.Writer, error)

	SetHostname(hostname string) error
	Apply(endpoint *metadata.NetworkEndpoint) error
	MountLabel(label, target string, ctx context.Context) error
	Fork() error

	SessionLog(session *SessionConfig) (dio.DynamicMultiWriter, error)
	HandleSessionExit(config *ExecutorConfig, session *SessionConfig) bool
	ProcessEnv(env []string) []string
}

type Tether interface {
	Start() error
	Stop() error
	Reload()
	Register(name string, ext Extension)
}

type Extension interface {
	Start() error
	Reload(config *ExecutorConfig) error
	Stop() error
}
