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

package common

import (
	"testing"
)

func TestCheckUnsupportedchars(t *testing.T) {
	tests := []struct {
		S     string
		valid bool
	}{
		{"anjunabeats", true},
		{"tony-1", true},
		{"paavo_1", true},
		{"jono(1)", true},
		{"oceanlab (1)", true},
		{"test%", false},
		{"test&", false},
		{"test*", false},
		{"test$", false},
		{"test#", false},
		{"test@", false},
		{"test!", false},
		{`test\`, false},
		{"test/", false},      // U+002F
		{"test\u002f", false}, // U+002F
		{`testЯ`, true},       // U+042F
		{"test\u042f", true},  // U+042F
		{`testį`, true},       // U+012F
		{"test\u012f", true},  // U+012F
		{"test:", false},
		{"test?", false},
		{`test"`, false},
		{"test<", false},
		{"test>", false},
		{"test;", false},
		{"test'", false},
		{"test|", false},
		{"test|", false},
	}

	for _, test := range tests {
		err := CheckUnsupportedChars(test.S)
		if err != nil {
			if test.valid {
				t.Errorf("got %q, expected pass for %q", err, test.S)
			}
			t.Logf("test case %q passed", test.S)
			continue
		}
		if test.valid {
			t.Logf("test case %q passed", test.S)
			continue
		}
		t.Errorf("got %q, expected error for %q", err, test.S)
	}
}

func TestCheckUnsupportedCharsDatastore(t *testing.T) {
	tests := []struct {
		S     string
		valid bool
	}{
		{"anjunabeats", true},
		{"tony-1", true},
		{"paavo_1", true},
		{"jono(1)", true},
		{"oceanlab (1)", true},
		{"waawn/", true},
		{"tristate:", true},
		{"test%", false},
		{"test&", false},
		{"test*", false},
		{"test$", false},
		{"test#", false},
		{"test@", false},
		{"test!", false},      // U+0021
		{"test\u0021", false}, // U+0021
		{`testġ`, true},       // U+0121
		{"test\u0121", true},  // U+0121
		{`test\`, false},
		{"test?", false},
		{`test"`, false},
		{"test<", false},
		{"test>", false},
		{"test;", false},
		{"test'", false},
		{"test|", false},
	}

	for _, test := range tests {
		err := CheckUnsupportedCharsDatastore(test.S)

		if err != nil {
			if test.valid {
				t.Errorf("got %q, expected pass for %q", err, test.S)
			}
			t.Logf("test case %q passed", test.S)
			continue
		}
		if test.valid {
			t.Logf("test case %q passed", test.S)
			continue
		}
		t.Errorf("got %q, expected error for %q", err, test.S)
	}
}

func TestCheckURLValidationOutputParity(t *testing.T) {
	tests := []struct {
		S     string
		valid bool
	}{
		{"nfs://127.0.0.1:/my/data/here:labels", true},
		{"NotAGoodScheme://127.0.0.1:/not/my/server:moreLabels", true},
		{"not a URL", false},
		{"127.0.0.1:/schemeless/url:importantTarget", false},
		{"//127.0.0.1/schemeless/url:importantTarget", false},
		{"Blargh://127.0.0.1/schemeless/url:importantTarget", true},
		{"Blargh://bob:pass@127.0.0.1/schemed/url:importantTarget", true},
		{"Blargh://jim:@127.0.0.1/schemed/url:importantTarget", true},
	}

	for _, test := range tests {
		url := CheckURLValidation(test.S)

		if test.valid {
			if url == nil {
				t.Errorf("Test URL (%s) was parsed as NIL when it should not have", test.S)
				t.Fail()
			}
		} else {
			if url != nil {
				t.Errorf("Test URL (%s) was parsed as a URL when NIL should have been returned. output: (%#v)", test.S, url)
				t.Fail()
			}
		}
	}
}

func TestCheckURLValidationOutputCorrectness(t *testing.T) {

	testRawURL := "blargh://jim:@127.0.0.1/great/path:importantTarget"
	u := CheckURLValidation(testRawURL)

	correctScheme := "blargh"
	correctUser := "jim"
	correctHost := "127.0.0.1"
	correctPath := "/great/path"

	if u.Scheme != correctScheme {
		t.Errorf("expected scheme (%s), but got scheme (%s)", correctScheme, u.Scheme)
		t.Fail()
	}

	if u.User.Username() != correctUser {
		t.Errorf("expected user (%s), but got user (%s)", correctUser, u.User.Username())
		t.Fail()
	}

	if u.Hostname() != correctHost {
		t.Errorf("expected host (%s), but got host (%s)", correctHost, u.Hostname())
		t.Fail()
	}

	if u.Path != correctPath {
		t.Errorf("expected path (%s), but got path (%s)", correctPath, u.Path)
		t.Fail()
	}
}
