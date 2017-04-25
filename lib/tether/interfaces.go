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

package tether

import (
	"context"
	"io"
	"net/url"

	"github.com/vmware/vic/pkg/dio"
)

// Operations defines the set of operations that Tether depends upon. These are split out for:
// * portability
// * dependency injection (primarily for testing)
// * behavioural control (e.g. what behaviour is required when a session exits)
type Operations interface {
	Setup(Config) error
	Cleanup() error
	// Log returns the tether debug log writer
	Log() (io.Writer, error)

	SetHostname(hostname string, aliases ...string) error
	SetupFirewall(config *ExecutorConfig) error
	Apply(endpoint *NetworkEndpoint) error
	MountLabel(ctx context.Context, label, target string) error
	MountTarget(ctx context.Context, source url.URL, target string, mountOptions string) error
	CopyExistingContent(source string) error
	Fork() error
	// Returns two DynamicMultiWriters for stdout and stderr
	SessionLog(session *SessionConfig) (dio.DynamicMultiWriter, dio.DynamicMultiWriter, error)
	// Returns a function to invoke after the session state has been persisted
	HandleSessionExit(config *ExecutorConfig, session *SessionConfig) func()
	ProcessEnv(env []string) []string
}

// Tether presents the consumption interface for code needing to run a tether
type Tether interface {
	Start() error
	Stop() error
	Reload()
	Register(name string, ext Extension)
}

// Extension is a very simple extension interface for supporting code that need to be
// notified when the configuration is reloaded.
type Extension interface {
	Start() error
	Reload(config *ExecutorConfig) error
	Stop() error
}

type Config interface {
	UpdateNetworkEndpoint(e *NetworkEndpoint) error
	Flush() error
}
