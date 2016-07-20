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

package datastore

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/test/env"
	"golang.org/x/net/context"
)

// Used in testing

func Session(ctx context.Context, t *testing.T) *session.Session {
	config := &session.Config{
		Service: env.URL(t),

		/// XXX Why does this insist on having this field populated?
		DatastorePath: env.DS(t),

		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}

	s, err := session.NewSession(config).Create(ctx)
	if err != nil {
		t.SkipNow()
	}

	return s
}

func DSsetup(t *testing.T) (context.Context, *Helper, func()) {
	ctx := context.Background()
	sess := Session(ctx, t)

	ds, err := NewHelper(ctx, sess, sess.Datastore, TestName("dstests"))
	if !assert.NoError(t, err) {
		return ctx, nil, nil
	}

	f := func() {
		log.Infof("Removing test root %s", ds.RootURL)
		err := tasks.Wait(ctx, func(context.Context) (tasks.Waiter, error) {
			return ds.fm.DeleteDatastoreFile(ctx, ds.RootURL, sess.Datacenter)
		})

		if err != nil {
			log.Errorf(err.Error())
			return
		}
	}

	return ctx, ds, f
}
