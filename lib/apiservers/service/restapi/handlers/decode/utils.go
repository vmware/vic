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

package decode

import (
	"regexp"
	"fmt"
)

// TODO [AngieCris]: complete duplicate of common/utils.go
// TODO [AngieCris]: move the regex validation to swagger (if possible)

// https://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2046088
const unsuppCharsRegex = `%|&|\*|\$|#|@|!|\\|/|:|\?|"|<|>|;|'|\|`

// Same as unsuppCharsRegex but allows / and : for datastore paths
const unsuppCharsDatastoreRegex = `%|&|\*|\$|#|@|!|\\|\?|"|<|>|;|'|\|`

var reUnsupp = regexp.MustCompile(unsuppCharsRegex)
var reUnsuppDatastore = regexp.MustCompile(unsuppCharsDatastoreRegex)


// CheckUnsupportedChars returns an error if string contains special characters
func CheckUnsupportedChars(s string) error {
	return checkUnsupportedChars(s, reUnsupp)
}

// CheckUnsupportedCharsDatastore returns an error if a datastore string contains special characters
func CheckUnsupportedCharsDatastore(s string) error {
	return checkUnsupportedChars(s, reUnsuppDatastore)
}

func checkUnsupportedChars(s string, re *regexp.Regexp) error {
	st := []byte(s)
	var v []int
	// this is validation step for characters in a datastore URI
	if v = re.FindIndex(st); v == nil {
		return nil
	}
	return fmt.Errorf("unsupported character %q in %q", s[v[0]:v[1]], s)
}