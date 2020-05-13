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

package plugina

type ExecutorConfig struct {
	Common `vic:"0.1" scope:"read-only" key:"common"`
}

type Common struct {
	// A reference to the components hosting execution environment, if any
	ExecutionEnvironment string

	// Unambiguous ID with meaning in the context of its hosting execution environment
	ID string `vic:"0.1" scope:"read-only" key:"id"`

	// Convenience field to record a human readable name
	Name string `vic:"0.1" scope:"hidden" key:"name"`

	// Freeform notes related to the entity
	Notes string `vic:"0.1" scope:"hidden" key:"notes"`
}

type UpdatedCommon struct {
	// A reference to the components hosting execution environment, if any
	ExecutionEnvironment string

	// Unambiguous ID with meaning in the context of its hosting execution environment
	ID string `vic:"0.1" scope:"read-only" key:"id"`

	// Convenience field to record a human readable name
	Name string `vic:"0.1" scope:"read-only" key:"name"`

	// Freeform notes related to the entity
	Notes string `vic:"0.1" scope:"hidden" key:"notes"`
}

type VirtualContainerHostConfigSpec struct {
	ExecutorConfig `vic:"0.1" scope:"read-only" key:"init"`
}

type UpdatedVCHConfigSpec struct {
	UpdatedExecutorConfig `vic:"0.1" scope:"read-only" key:"init"`
}

type UpdatedExecutorConfig struct {
	UpdatedCommon `vic:"0.1" scope:"read-only" key:"common"`
}
