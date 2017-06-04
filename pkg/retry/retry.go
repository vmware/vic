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

package retry

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cenkalti/backoff"

	"github.com/vmware/vic/pkg/trace"
)

const (
	// Random numbers generated by fair dice roll :)
	defaultInitialInterval     = 10 * time.Millisecond
	defaultRandomizationFactor = 0.5
	defaultMultiplier          = 10
	defaultMaxInterval         = 30 * time.Second
	defaultMaxElapsedTime      = 1 * time.Minute
)

// simplified config for configuring a back off object. Callers should populate and supply this to DoWithConfig
type BackoffConfig struct {
	InitialInterval     time.Duration
	RandomizationFactor float64
	Multiplier          float64
	MaxInterval         time.Duration

	// this field will indicate the maximum amount of "sleep" time that will occur.
	MaxElapsedTime time.Duration
}

// Generate a new BackoffConfig with default values
func NewBackoffConfig() *BackoffConfig {
	return &BackoffConfig{
		InitialInterval:     defaultInitialInterval,
		RandomizationFactor: defaultRandomizationFactor,
		Multiplier:          defaultMultiplier,
		MaxInterval:         defaultMaxInterval,
		MaxElapsedTime:      defaultMaxElapsedTime,
	}
}

// Do retries the given function until defaultMaxInterval time passes, while sleeping some time between unsuccessful attempts
// if retryOnError returns true, continue retry, otherwise, return error
func Do(operation func() error, retryOnError func(err error) bool) error {
	bConf := &BackoffConfig{
		InitialInterval:     defaultInitialInterval,
		RandomizationFactor: defaultRandomizationFactor,
		Multiplier:          defaultMultiplier,
		MaxInterval:         defaultMaxInterval,
		MaxElapsedTime:      defaultMaxElapsedTime,
	}

	return DoWithConfig(operation, retryOnError, bConf)
}

// DoWithConfig will attempt an operation while retrying using an exponential back off based on the config supplied by the caller. The retry decider is the supplied function retryOnError
func DoWithConfig(operation func() error, retryOnError func(err error) bool, config *BackoffConfig) error {
	defer trace.End(trace.Begin(""))

	var err error
	var next time.Duration
	b := &backoff.ExponentialBackOff{
		InitialInterval:     config.InitialInterval,
		RandomizationFactor: config.RandomizationFactor,
		Multiplier:          config.Multiplier,
		MaxInterval:         config.MaxInterval,
		MaxElapsedTime:      config.MaxElapsedTime,
		Clock:               backoff.SystemClock,
	}
	// Reset the interval back to the initial retry interval and restart the timer
	b.Reset()
	for {
		if err = operation(); err == nil {
			log.Debugf("Will not try again. Operation succeeded")
			return nil
		}

		if next = b.NextBackOff(); next == backoff.Stop {
			log.Errorf("Will stop trying again. Operation failed with %#+v", err)
			return err
		}

		// check error
		if !retryOnError(err) {
			log.Errorf("Operation failed with %#+v", err)
			return err
		}
		// Expected error
		log.Warnf("Will try again in %s. Operation failed with detected error", next)

		// sleep and try again
		time.Sleep(next)
		continue
	}
}

// OnError is the simplest of retry deciders. If an error occurs it will indicate a retry is needed.
func OnError(err error) bool {
	return err != nil
}
