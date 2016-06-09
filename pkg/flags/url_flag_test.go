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

package flags

import (
	"flag"
	"net/url"
	"testing"
)

func TestURLFlag(t *testing.T) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	var val *url.URL

	fs.Var(NewURLFlag(&val), "url", "url flag")

	u := fs.Lookup("url")

	if u.DefValue != "<nil>" {
		t.Errorf("DefValue: %s", u.DefValue)
	}

	if u.Value.String() != "<nil>" {
		t.Errorf("Value: %s", u.Value)
	}

	u.Value.Set("127.0.0.1")

	if u.Value.String() != "https://:@127.0.0.1/sdk" {
		t.Errorf("Value after set: %s", u.Value)
	}

	if val == nil {
		t.Errorf("val is not set")
	}
	if val.String() != "https://:@127.0.0.1/sdk" {
		t.Errorf("val is not set correctly: %s", val.String())
	}
}
