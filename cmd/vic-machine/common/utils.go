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
	"fmt"
	"net/url"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/pkg/errors"
)

// https://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2046088
const unsuppCharsRegex = `%|&|\*|\$|#|@|!|\\|/|:|\?|"|<|>|;|'|\|`

// Same as unsuppCharsRegex but allows / and : for datastore paths
const unsuppCharsDatastoreRegex = `%|&|\*|\$|#|@|!|\\|\?|"|<|>|;|'|\|`

var reUnsupp = regexp.MustCompile(unsuppCharsRegex)
var reUnsuppDatastore = regexp.MustCompile(unsuppCharsDatastoreRegex)

func LogErrorIfAny(clic *cli.Context, err error) error {
	if err == nil {
		return nil
	}

	log.Errorf("--------------------")
	log.Errorf("%s %s failed: %s\n", clic.App.Name, clic.Command.Name, errors.ErrorStack(err))
	return cli.NewExitError("", 1)
}

// CheckUnsupportedChars returns an error if string contains special characters
func CheckUnsupportedChars(s string) error {
	return checkUnsupportedChars(s, reUnsuppDatastore)
}

// CheckUnsupportedCharsDatastore returns an error if a datastore string contains special characters
func CheckUnsupportedCharsDatastore(s string) error {
	return checkUnsupportedChars(s, reUnsuppDatastore)
}

func checkUnsupportedChars(s string, re *regexp.Regexp) error {
	st := []byte(s)
	var v []int
	//this is validation step for characters in a datastore URI
	if v = re.FindIndex(st); v == nil {
		return nil
	}
	return fmt.Errorf("unsupported character %q in %q", s[v[0]:v[1]], s)
}

// CheckNFSURlValidation is used to validate propernfs arguments to the --volume-store flag and avoid checking inputs for datastore restricted characters if a url is provided
func CheckURLValidation(input string) *url.URL {
	parts := strings.Split(input, ":")
	rawURL := strings.Join(parts[0:len(parts)-1], ":")

	if len(parts) < 2 {
		return nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}

	// at the very lease, these fields must exist for a good input to exist for an nfs share
	if u.Host == "" || u.Scheme == "" || u.Path == "" {
		return nil
	}
	return u
}
