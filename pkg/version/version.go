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

package version

import (
	"flag"
	"fmt"
	"runtime"
)

// These fields are set by the compiler using the linker flags upon build via Makefile.
var (
	Version     string
	GitCommit   string
	BuildDate   string
	BuildNumber string
	State       string

	v bool
)

type Build struct {
	Version     string
	GitCommit   string
	BuildDate   string
	BuildNumber string
	State       string
}

func init() {
	flag.BoolVar(&v, "version", false, "Show version info")
}

// Show returns whether -version flag is set
func Show() bool {
	return v
}

// String returns a string representation of the version
func String() string {
	return GetBuild().String()
}

func GetBuild() *Build {
	return &Build{
		Version:     Version,
		GitCommit:   GitCommit,
		BuildDate:   BuildDate,
		BuildNumber: BuildNumber,
		State:       State,
	}
}

func (v *Build) String() string {
	if v.State == "" {
		v.State = "clean"
	}

	if v.BuildNumber == "" {
		v.BuildNumber = "N/A"
	}
	return fmt.Sprintf("%s git:%s-%s build:%s id:%s runtime:%s", v.Version, v.GitCommit, v.State, v.BuildDate, v.BuildNumber, runtime.Version())
}

func (v *Build) ShortVersion() string {
	return fmt.Sprintf("%s-%s-%s", v.Version, v.BuildNumber, v.GitCommit)
}
