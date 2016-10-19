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

package log

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

var (
	baseTimestamp time.Time
	isTerminal    bool
)

func init() {
	baseTimestamp = time.Now()
	isTerminal = logrus.IsTerminal()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

type TextFormatter struct {

	// Enable colors
	Colors bool

	// Enable timestamp logging. It's useful to disable timestamps when output
	// is redirected to logging system that already adds timestamps.  When
	// disabled, a delta is used for each log line from the start of the
	// process.  Enabling will apply a timestamp derived from TimestampFormat.
	Timestamp bool

	// TimestampFormat is the format used to print the timestamp.  By default
	// an RFC3339 timestamp is used.
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// Escape noncharacter ascii strings.
	EscapeNonCharacters bool
}

// NewTextFormatter returns a text formatter with defaults appropriate for the
// TTY or file being written to.
func NewTextFormatter() *TextFormatter {

	isColorTerminal := isTerminal && (runtime.GOOS != "windows")

	formatter := &TextFormatter{
		Colors:          isColorTerminal,
		Timestamp:       false,
		TimestampFormat: time.RFC3339,
		DisableSorting:  false,
	}

	return formatter
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}

	b = &bytes.Buffer{}

	// get the time string
	ts := f.timeStamp(entry)

	if f.Colors {
		f.printColored(b, entry, keys, ts)
	} else {
		if f.Timestamp {
			f.appendKeyValue(b, "time", ts)
		}

		f.appendKeyValue(b, "level", entry.Level.String())
		if entry.Message != "" {
			f.appendKeyValue(b, "msg", entry.Message)
		}

		for _, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key])
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) timeStamp(entry *logrus.Entry) string {
	if !f.Timestamp {
		return fmt.Sprintf("%04d", miniTS())
	}

	timestampFormat := f.TimestampFormat

	if timestampFormat == "" {
		timestampFormat = logrus.DefaultTimestampFormat
	}

	return entry.Time.Format(timestampFormat)
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestamp string) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]

	fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %-44s ", levelColor, levelText, timestamp, entry.Message)

	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%+v", levelColor, k, v)
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if !f.EscapeNonCharacters {
		return false
	}

	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {

	b.WriteString(key)
	b.WriteByte('=')

	switch value := value.(type) {
	case string:
		if !f.needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	case error:
		errmsg := value.Error()
		if !f.needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	default:
		fmt.Fprint(b, value)
	}

	b.WriteByte(' ')
}
