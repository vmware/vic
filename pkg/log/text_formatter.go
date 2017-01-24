// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

package log

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/Sirupsen/logrus"
)

type TextFormatter struct {
	// TimestampFormat is the format used to print the timestamp.  By default
	// an RFC3339 timestamp is used.
	TimestampFormat string
}

// NewTextFormatter returns a text formatter
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		TimestampFormat: "Jan _2 2006 15:04:05.000Z07:00",
	}
}

func trimOrPadToSize(s string, size int) string {
	if len(s) > size {
		return s[:size]
	}

	return s + strings.Repeat(" ", size-len(s))
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	t := f.timeStamp(entry)
	l := strings.ToUpper(trimOrPadToSize(entry.Level.String(), 5))

	var b bytes.Buffer

	if entry.Message == "" {
		b.Grow(len(t) + len(" ") + len(l) + len("\n"))
		b.WriteString(t)
		b.WriteString(" ")
		b.WriteString(l)
		b.WriteString("\n")
		return b.Bytes(), nil
	}

	s := bufio.NewScanner(strings.NewReader(entry.Message))
	for s.Scan() {
		b.Grow(len(t) + len(" ") + len(l) + len(" ") + len(s.Text()) + len("\n"))
		b.WriteString(t)
		b.WriteString(" ")
		b.WriteString(l)
		b.WriteString(" ")
		b.WriteString(s.Text())
		b.WriteString("\n")
	}

	if s.Err() != nil {
		return nil, s.Err()
	}

	return b.Bytes(), nil
}

func (f *TextFormatter) timeStamp(entry *logrus.Entry) string {
	return entry.Time.Format(f.TimestampFormat)
}
