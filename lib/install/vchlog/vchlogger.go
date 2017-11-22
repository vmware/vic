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

package vchlog

import (
	"io/ioutil"
	"path"
	"time"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/vic/pkg/trace"
)

// Reminder: When changing these, consider the impact to backwards compatibility.
const (
	LogFilePrefix = "vic-machine" // logFilePrefix is the prefix for file names of all vic-machine log files
	LogFileSuffix = ".log"        // logFileSuffix is the suffix for file names of all vic-machine log files
)

// DatastoreReadySignal serves as a signal struct indicating datastore folder path is available
// Datastore: the govmomi datastore object
// Name: the name of the vic-machine process that sends the signal (e.g. "create", "inspect")
// LogFileName: the filename of the destination path on datastore
// Context: the caller context when sending the signal
// VMPathName: the datastore path
type DatastoreReadySignal struct {
	Datastore  *object.Datastore
	Name       string
	Operation  trace.Operation
	VMPathName string
	Timestamp  time.Time
}

type VCHLogger struct {
	// pipe: the streaming readwriter pipe to hold log messages
	pipe *BufferedPipe

	// signalChan: channel for signaling when datastore folder is ready
	signalChan chan DatastoreReadySignal
}

type Receiver interface {
	Signal(sig DatastoreReadySignal)
}

// New creates the logger, with the streaming pipe and singaling channel.
func New() *VCHLogger {
	return &VCHLogger{
		pipe:       NewBufferedPipe(),
		signalChan: make(chan DatastoreReadySignal),
	}
}

// Run waits until the signal arrives and uploads the streaming pipe to datastore
func (l *VCHLogger) Run() {
	sig := <-l.signalChan
	// suffix the log file name with caller operation ID and timestamp
	logFileName := LogFilePrefix + "_" + sig.Timestamp.UTC().Format(time.RFC3339) + "_" + sig.Name + "_" + sig.Operation.ID() + LogFileSuffix
	param := soap.DefaultUpload
	param.ContentLength = -1
	sig.Datastore.Upload(sig.Operation.Context, ioutil.NopCloser(l.pipe), path.Join(sig.VMPathName, logFileName), &param)
}

// GetPipe returns the streaming pipe of the vch logger
func (l *VCHLogger) GetPipe() *BufferedPipe {
	return l.pipe
}

// Signal signals the logger that the datastore folder is ready
func (l *VCHLogger) Signal(sig DatastoreReadySignal) {
	l.signalChan <- sig
}

// Close stops the logger by closing the underlying pipe
func (l *VCHLogger) Close() error {
	return l.pipe.Close()
}
