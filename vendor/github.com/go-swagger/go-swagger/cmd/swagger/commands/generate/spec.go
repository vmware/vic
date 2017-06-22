// Copyright 2015 go-swagger maintainers
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

package generate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-swagger/go-swagger/scan"
	"github.com/jessevdk/go-flags"
)

// SpecFile command to generate a swagger spec from a go application
type SpecFile struct {
	BasePath   string         `long:"base-path" short:"b" description:"the base path to use" default:"."`
	BuildTags  string         `long:"tags" short:"t" description:"build tags" default:""`
	ScanModels bool           `long:"scan-models" short:"m" description:"includes models that were annotated with 'swagger:model'"`
	Compact    bool           `long:"compact" description:"when present, doesn't prettify the the json"`
	Output     flags.Filename `long:"output" short:"o" description:"the file to write to"`
	Input      flags.Filename `long:"input" short:"i" description:"the file to use as input"`
}

// Execute runs this command
func (s *SpecFile) Execute(args []string) error {
	input, err := loadSpec(string(s.Input))
	if err != nil {
		return err
	}

	var opts scan.Opts
	opts.BasePath = s.BasePath
	opts.Input = input
	opts.ScanModels = s.ScanModels
	opts.BuildTags = s.BuildTags
	swspec, err := scan.Application(opts)
	if err != nil {
		return err
	}

	return writeToFile(swspec, !s.Compact, string(s.Output))
}

func loadSpec(input string) (*spec.Swagger, error) {
	if fi, err := os.Stat(input); err == nil {
		if fi.IsDir() {
			return nil, fmt.Errorf("expected %q to be a file not a directory", input)
		}
		sp, err := loads.Spec(input)
		if err != nil {
			return nil, err
		}
		return sp.Spec(), nil
	}
	return nil, nil
}

func writeToFile(swspec *spec.Swagger, pretty bool, output string) error {
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(swspec, "", "  ")
	} else {
		b, err = json.Marshal(swspec)
	}
	if err != nil {
		return err
	}
	if output == "" {
		fmt.Println(string(b))
		return nil
	}
	return ioutil.WriteFile(output, b, 0644)
}
