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

// Package utils provides basic helper functions for interacting with DOS
package utils

import (
	"bufio"
	"strings"
)

// Output from a cmd call to MS-DOS consists of an echo of the command and the output, separated by CRLF
// We want just the second line without the CRLF
func StripCommandOutput(output string) string {
	reader := strings.NewReader(output)
	reader2 := bufio.NewReader(reader)
	_, _ = reader2.ReadString('\n') // throw away the echo line
	line2, err := reader2.ReadString('\n')
	if err == nil && len(line2) > 2 {
		return line2[:len(line2)-2] // Assume line ends with CRLF
	}
	return ""
}
