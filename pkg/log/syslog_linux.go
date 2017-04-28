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
	"log/syslog"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cenkalti/backoff"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	Writer   *syslog.Writer
	Proto    string
	RAddr    string
	Priority syslog.Priority
	Tag      string

	entries chan *logrus.Entry
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewSyslogHook(network, raddr string, priority syslog.Priority, tag string) *SyslogHook {
	return &SyslogHook{
		Proto:    network,
		RAddr:    raddr,
		Priority: priority,
		Tag:      tag,
		entries:  make(chan *logrus.Entry, 100),
	}
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	select {
	case hook.entries <- entry:
	default:
		// make room
		<-hook.entries
		hook.entries <- entry
	}

	return nil
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func AddSyslog(proto, raddr, tag string) error {
	hook := NewSyslogHook(proto, raddr, syslog.LOG_INFO, tag)
	w, err := syslog.Dial(hook.Proto, hook.RAddr, hook.Priority, hook.Tag)
	if err != nil {
		switch err.(type) {
		case *net.AddrError, *net.ParseError:
			return err
		}
	} else {
		hook.Writer = w
	}

	go func() {
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 0

		for {
			b.Reset()
			for hook.Writer == nil {
				w, err := syslog.Dial(hook.Proto, hook.RAddr, hook.Priority, hook.Tag)
				if err != nil {
					// ignore error and retry
					time.Sleep(b.NextBackOff())

					continue
				}

				hook.Writer = w
			}

			for entry := range hook.entries {
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
					hook.Writer.Close()
					hook.Writer = nil
					break
				}
			}
		}
	}()

	logrus.AddHook(hook)
	return nil
}
