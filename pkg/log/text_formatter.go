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

// level strings padded to match the length of the longest level,
// which is "UNKNOWN" currently. Indexed according to levels in
// logrus, e.g. levelStrs[logrus.InfoLevel] == "INFO ".
var levelStrs = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN ",
	"INFO ",
	"DEBUG",
}

const unknownLevel = "UNKWN"

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

func levelToString(level logrus.Level) string {
	if level <= logrus.DebugLevel {
		return levelStrs[level]
	}

	return unknownLevel
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	t := f.timeStamp(entry)
	l := levelToString(entry.Level)

	if entry.Message == "" {
		return []byte(t + " " + l + "\n"), nil
	}

	// prefix each line of the message with timestamp
	// level information
	s := bufio.NewScanner(strings.NewReader(entry.Message))
	// Define a split function that separates on newlines
	onNewLine := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == '\n' {
				return i + 1, data[:i], nil
			}
		}
		// There is one final token to be delivered, which may be the empty string.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
		// but does not trigger an error to be returned from Scan itself.
		return 0, data, bufio.ErrFinalToken
	}
	s.Split(onNewLine)

	var b bytes.Buffer
	for s.Scan() {
		b.WriteString(t + " " + l + " " + s.Text() + "\n")
	}

	if s.Err() != nil {
		return nil, s.Err()
	}

	return b.Bytes(), nil
}

func (f *TextFormatter) timeStamp(entry *logrus.Entry) string {
	return entry.Time.Format(f.TimestampFormat)
}
