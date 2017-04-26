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

package syslog

import (
	"log/syslog"
	"net"

	"github.com/Sirupsen/logrus"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	Writer        *syslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewSyslogHook(network, raddr string, priority syslog.Priority, tag string) (*SyslogHook, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &SyslogHook{
		Writer:        w,
		SyslogNetwork: network,
		SyslogRaddr:   raddr,
	}, err
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	// just use the message since the timestamp
	// is added by the syslog package
	line := entry.Message

	err := func() error {
		switch entry.Level {
		case logrus.PanicLevel:
			return hook.Writer.Crit(line)
		case logrus.FatalLevel:
			return hook.Writer.Crit(line)
		case logrus.ErrorLevel:
			return hook.Writer.Err(line)
		case logrus.WarnLevel:
			return hook.Writer.Warning(line)
		case logrus.InfoLevel:
			return hook.Writer.Info(line)
		case logrus.DebugLevel:
			return hook.Writer.Debug(line)
		default:
			return nil
		}
	}()

	if err != nil {
		// this will force a re-connection
		hook.Writer.Close()
	}

	return err
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func AddSyslog(proto, raddr, tag string) error {
	hook, err := NewSyslogHook(proto, raddr, syslog.LOG_INFO, tag)
	if err := err.(net.Error); err != nil {
		if !err.Temporary() && !err.Timeout() {
			return err
		}
	}

	logrus.AddHook(hook)
	return nil
}
