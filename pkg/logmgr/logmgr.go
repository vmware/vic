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

package logmgr

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/trace"
)

type RotateInterval uint32

const LogRotateBinary = "/usr/sbin/logrotate"

const (
	Daily RotateInterval = iota
	Hourly
	Weekly
	Monthly
)

type logRotateConfig struct {
	rotateInterval RotateInterval
	logFilePath    string
	logFileName    string
	maxLogSize     int64
	maxLogFiles    int64
	compress       bool
}

// ConfigFileContent formats log configuration according to logrotate requirements.
func (lrc *logRotateConfig) ConfigFileContent() string {
	b := make([]string, 0, 32)
	if lrc.compress {
		b = append(b, "compress")
	}

	switch lrc.rotateInterval {
	case Hourly:
		b = append(b, string("hourly"))
	case Daily:
		b = append(b, string("daily"))
	case Weekly:
		b = append(b, string("weekly"))
	case Monthly:
		fallthrough
	default:
		b = append(b, string("monthly"))
	}

	b = append(b, fmt.Sprintf("rotate %d", lrc.maxLogFiles))
	if lrc.maxLogSize > 0 {
		b = append(b, fmt.Sprintf("size %d", lrc.maxLogSize))
	}

	if lrc.maxLogSize > 2 {
		b = append(b, fmt.Sprintf("minsize %d", lrc.maxLogSize-1))
	}

	b = append(b, "copytruncate")
	b = append(b, "dateext")
	b = append(b, "dateformat -%Y%m%d-%s")

	for i, v := range b {
		b[i] = "    " + v
	}

	return fmt.Sprintf("%s {\n %s\n}\n", lrc.logFilePath, strings.Join(b, "\n"))
}

// LogManager runs logrotate for specified log files.
// TODO: Upload compressed logs into vSphere storage.
// TODO: Upload all logs into vSphere storage during graceful shutdown.
type logManager struct {
	// Frequency of running log rotate.
	runInterval time.Duration

	// list of log files and their log rotate parameters to rotate.
	logFiles []*logRotateConfig

	// channel gets closed on stop.
	closed chan struct{}
	op     trace.Operation

	// used to wait until logrotate goroutine stops.
	wg sync.WaitGroup
	// just to make sure start is not called twice accidentaly.
	once sync.Once

	logConfig string
}

// NewLogManager creates a new log manager instance.
func NewLogManager(runInterval time.Duration) (*logManager, error) {
	lm := &logManager{
		runInterval: runInterval,
		op:          trace.NewOperation(context.Background(), "logrotate"),
	}
	if s, err := os.Stat(LogRotateBinary); err != nil || s.IsDir() {
		return nil, fmt.Errorf("logrotate is not available at %s, without it logs will not be rotated", LogRotateBinary)
	}
	return lm, nil
}

// AddLogRotate adds a log to rotate.
func (lm *logManager) AddLogRotate(logFilePath string, ri RotateInterval, maxSize, maxLogFiles int64, compress bool) {
	lm.logFiles = append(lm.logFiles, &logRotateConfig{
		rotateInterval: ri,
		logFilePath:    logFilePath,
		logFileName:    filepath.Base(logFilePath),
		maxLogSize:     maxSize,
		maxLogFiles:    maxLogFiles,
		compress:       compress,
	})
}

// Reload - just to satisfy Tether interface.
func (lm *logManager) Reload(*tether.ExecutorConfig) error { return nil }

// Start log rotate loop.
func (lm *logManager) Start() error {
	if len(lm.logFiles) == 0 {
		lm.op.Errorf("Attempt to start logrotate with no log files configured.")
		return nil
	}
	lm.once.Do(func() {
		lm.wg.Add(1)
		lm.logConfig = lm.buildConfig()
		lm.op.Debugf("logrotate config: %s", lm.logConfig)

		go func() {
			for {
				lm.rotateLogs()
				select {
				case <-time.After(lm.runInterval):
				case <-lm.closed:
					lm.rotateLogs()
					lm.wg.Done()
					return
				}
			}
		}()
	})
	return nil
}

// Stop loop.
func (lm *logManager) Stop() error {
	select {
	case <-lm.closed:
	default:
		close(lm.closed)
	}
	lm.wg.Wait()
	return nil
}

func (lm *logManager) saveConfig(logConf string) string {
	tf, err := ioutil.TempFile("", "vic-logrotate-conf-")
	if err != nil {
		lm.op.Errorf("Failed to create temp file for logrotate: %v", err)
		return ""
	}

	tempFilePath := tf.Name()
	if _, err = tf.Write([]byte(logConf)); err != nil {
		lm.op.Errorf("Failed to store logrotate config %s: %v", tempFilePath, err)
		if err = tf.Close(); err != nil {
			lm.op.Errorf("Failed to close temp file %s: %v", tempFilePath, err)
		}
		if err = os.Remove(tempFilePath); err != nil {
			lm.op.Errorf("Failed to remove temp file %s: %v", tempFilePath, err)
		}
		return ""
	}

	if err = tf.Close(); err != nil {
		lm.op.Errorf("Failed to close logrotate config file %s: %v", tempFilePath, err)
		return ""
	}
	return tempFilePath
}

func (lm *logManager) buildConfig() string {
	c := make([]string, 0, len(lm.logFiles))
	for _, v := range lm.logFiles {
		c = append(c, v.ConfigFileContent())
	}
	return strings.Join(c, "\n")
}

func (lm *logManager) rotateLogs() {
	// Check if logrotate config exists, create one
	configFile := lm.saveConfig(lm.logConfig)

	if configFile == "" {
		lm.op.Errorf("Can not run logrotate dues to missing logrotate config")
		return
	}
	// remove config file as soon as logrotate finishes its work.
	defer os.Remove(configFile)

	lm.op.Debugf("Running logrotate: %s %s", LogRotateBinary, configFile)

	if err := exec.Command(LogRotateBinary, configFile).Run(); err == nil {
		lm.op.Debugf("logrotate finished succesfully")
	} else {
		lm.op.Errorf("Failed to run logrotate: %v", err)
	}
}
