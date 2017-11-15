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
	"path"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/trace"
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
	Timestamp  string
}

type VCHLogger struct {
	Pipe *BufferedPipe
	SignalChan chan DatastoreReadySignal
}

//// pipe: the streaming readwriter pipe to hold log messages
//var pipe *BufferedPipe
//
//// signalChan: channel for signaling when datastore folder is ready
//var signalChan chan DatastoreReadySignal

func NewVCHLogger() *VCHLogger {
	return &VCHLogger{
		Pipe: NewBufferedPipe(),
		SignalChan: make(chan DatastoreReadySignal),
	}
}

//// Init initializes the logger, creates the streaming pipe and makes the singaling channel.
//func Init() {
//	pipe = NewBufferedPipe()
//	signalChan = make(chan DatastoreReadySignal)
//}

// Run waits until the signal arrives and uploads the streaming pipe to datastore
func (logger *VCHLogger) Run() {
	sig := <-logger.SignalChan
	// suffix the log file name with caller operation ID and timestamp
	logFileName := "vic-machine" + "_" + sig.Timestamp + "_" + sig.Name + "_" + sig.Operation.ID() + ".log"
	sig.Datastore.Upload(sig.Operation.Context, logger.Pipe, path.Join(sig.VMPathName, logFileName), nil)
}

//// GetPipe returns the streaming pipe of the vch logger
//func GetPipe() *BufferedPipe {
//	return pipe
//}

// Signal signals the logger that the datastore folder is ready
func (logger *VCHLogger) Signal(sig DatastoreReadySignal) {
	logger.SignalChan <- sig
}

// Close stops the logger by closing the underlying pipe
func (logger *VCHLogger) Close() {
	logger.Pipe.Close()
}
