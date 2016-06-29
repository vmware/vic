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

package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestNewHook(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	hook := NewErrorHook(w)
	levels := hook.Levels()
	if len(levels) != 1 && levels[0] != logrus.FatalLevel {
		t.Errorf("Returned levels %#v are different than expected", levels)
	}

	w.Reset(&b)
	hook.Fire(&logrus.Entry{Message: "Fatal Test", Level: logrus.FatalLevel})
	w.Flush()

	if string(b.Bytes()) == "" {
		t.Errorf("Fatal test failed %s", string(b.Bytes()))
	}
}
