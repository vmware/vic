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

package vsphere

import (
	"math/rand"
	"net/url"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/test/env"
	"golang.org/x/net/context"
)

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

// test if we can get a Datastore via the rooturl
func TestDatastoreGetDatastores(t *testing.T) {
	ctx, ds, cleanupfunc := dSsetup(t)
	if t.Failed() {
		return
	}
	defer cleanupfunc()

	firstSummary, err := ds.Summary(ctx)
	if !assert.NoError(t, err) {
		return
	}

	t.Logf("Name:\t%s\n", firstSummary.Name)
	t.Logf("  Path:\t%s\n", ds.ds.InventoryPath)
	t.Logf("  Type:\t%s\n", firstSummary.Type)
	t.Logf("  URL:\t%s\n", firstSummary.Url)
	t.Logf("  Capacity:\t%.1f GB\n", float64(firstSummary.Capacity)/(1<<30))
	t.Logf("  Free:\t%.1f GB\n", float64(firstSummary.FreeSpace)/(1<<30))

	inMap := make(map[string]*url.URL)

	p, err := url.Parse(ds.RootURL)
	if !assert.NoError(t, err) {
		return
	}

	inMap["foo"] = p

	dstores, err := GetDatastores(context.TODO(), ds.s, inMap)
	if !assert.NoError(t, err) || !assert.NotNil(t, dstores) {
		return
	}

	secondSummary, err := ds.Summary(ctx)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Equal(t, firstSummary, secondSummary) {
		return
	}
}

func TestDatastoreCreateDir(t *testing.T) {
	ctx, ds, cleanupfunc := dSsetup(t)
	if t.Failed() {
		return
	}
	defer cleanupfunc()

	_, err := ds.Ls(ctx, "")
	if !assert.NoError(t, err) {
		return
	}
}

func TestDatastoreMkdirAndLs(t *testing.T) {
	ctx, ds, cleanupfunc := dSsetup(t)
	if t.Failed() {
		return
	}
	defer cleanupfunc()

	dirs := []string{"dir1", "dir1/child1"}

	// create the dir then test it exists by calling ls
	for _, dir := range dirs {
		_, err := ds.Mkdir(ctx, true, dir)
		if !assert.NoError(t, err) {
			return
		}

		_, err = ds.Ls(ctx, dir)
		if !assert.NoError(t, err) {
			return
		}
	}
}

// From https://siongui.github.io/2015/04/13/go-generate-random-string/
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func dSsetup(t *testing.T) (context.Context, *Datastore, func()) {
	ctx := context.Background()
	sess := Session(ctx, t)
	log.SetLevel(log.DebugLevel)

	ds, err := NewDatastore(ctx, sess, sess.Datastore, RandomString(10)+"dstests")
	if !assert.NoError(t, err) {
		return ctx, nil, nil
	}

	f := func() {
		log.Debugf("Removing test root %s", ds.RootURL)
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
