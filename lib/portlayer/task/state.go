// Copyright 2018 VMware, Inc. All Rights Reserved.
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

package task

import (
	"github.com/vmware/vic/lib/config/executor"
)

const (
	runningState = "running"
	stoppedState = "stopped"
	createdState = "created"
	failedState  = "failed"
	unknownState = "unknown"
)

// State takes the given snapshot of a Task and determines the state of the Task.
func State(e *executor.SessionConfig) (string, error) {
	// TODO: Should probably be made a function on SessionConfig...

	switch {
	case e.Started == "" && e.Detail.StartTime == 0 && e.Detail.StopTime == 0:
		return createdState, nil
	case e.Started == "true" && e.Detail.StartTime > e.StopTime:
		return runningState, nil
	case e.Started == "true" && e.Detail.StartTime <= e.Detail.StopTime:
		return stoppedState, nil
	case e.Started != "" && e.StartTime > e.Detail.StartTime:
		return failedState, nil
	default:
		return unknownState, nil
	}
}
