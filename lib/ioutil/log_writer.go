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

package ioutil

import (
	"encoding/json"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
)

// LogEntry attaches a timestamp to a log entry
type LogEntry struct {
	Log  string    `json:"log"`
	Time time.Time `json:"time"`
}

// LogWriter wraps entries in a LogEntry struct with the
// timestamp of the entry.
type LogWriter struct {
	w io.Writer
}

// NewLogWriter wraps an io.WriteCloser in a LogWriter
func NewLogWriter(w io.WriteCloser) *LogWriter {
	return &LogWriter{w: w}
}

// Write takes an incoming byte slice and wraps it in a LogEntry,
// adding a timestamp. It then serializes the LogEntry and writes the
// serialized bytes to the underlying stream, adding a newline at the end for
// human readability
func (lw *LogWriter) Write(p []byte) (n int, err error) {

	bytes, err := json.Marshal(LogEntry{
		Log:  string(p),
		Time: time.Now(),
	})

	if err != nil {
		return 0, err
	}

	w := 0 // running total of bytes written
	for w != len(bytes) {
		n, err := lw.w.Write(bytes[w:])
		if err != nil {
			log.Errorf("Error writing JSON entry: %s", err.Error())
			return n, err
		}
		w += n
	}

	// add a newline to the end of a JSON log entry
	lw.w.Write([]byte{'\n'})

	// inform the caller that all received bytes were successfully written
	return len(p), err
}
