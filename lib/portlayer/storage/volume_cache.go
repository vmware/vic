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

package storage

import (
	"net/url"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// VolumeLookupCache caches Volume references to volumes in the system.
type VolumeLookupCache struct {

	// Maps IDs to Volumes.
	//
	// id -> Volume
	vlc     map[string]Volume
	vlcLock sync.Mutex

	// The underlying data storage implementation
	volumeStore VolumeStorer
}

func NewVolumeLookupCache(ctx context.Context, vs VolumeStorer) (*VolumeLookupCache, error) {
	v := &VolumeLookupCache{
		vlc:         make(map[string]Volume),
		volumeStore: vs,
	}

	return v, v.rebuildCache(ctx)
}

func (v *VolumeLookupCache) VolumeCreate(ctx context.Context, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*Volume, error) {
	v.vlcLock.Lock()
	defer v.vlcLock.Unlock()

	// check if it exists
	_, ok := v.vlc[ID]
	if ok {
		return nil, os.ErrExist
	}

	vol, err := v.volumeStore.VolumeCreate(ctx, ID, store, capacityKB, info)
	if err != nil {
		return nil, err
	}
	// Add it to the cache.
	v.vlc[vol.ID] = *vol

	return vol, nil
}

func (v *VolumeLookupCache) VolumeDestroy(ctx context.Context, ID string) error {
	v.vlcLock.Lock()
	defer v.vlcLock.Unlock()

	// Check if it exists
	vol, ok := v.vlc[ID]
	if !ok {
		return os.ErrNotExist
	}

	// remove it from the volumestore
	if err := v.volumeStore.VolumeDestroy(ctx, ID); err != nil {
		return err
	}
	delete(v.vlc, vol.ID)

	return nil
}

func (v *VolumeLookupCache) VolumeGet(ctx context.Context, ID string) (*Volume, error) {
	v.vlcLock.Lock()
	defer v.vlcLock.Unlock()

	// look in the cache
	vol, ok := v.vlc[ID]
	if !ok {
		return nil, os.ErrNotExist
	}

	return &vol, nil
}

func (v *VolumeLookupCache) VolumesList(ctx context.Context) ([]*Volume, error) {
	v.vlcLock.Lock()
	defer v.vlcLock.Unlock()

	// look in the cache, return the list
	l := make([]*Volume, 0, len(v.vlc))
	for _, vol := range v.vlc {
		// this is idiotic
		var e Volume
		e = vol
		l = append(l, &e)
	}

	return l, nil
}

// goto the volume store and repopulate the cache.
func (v *VolumeLookupCache) rebuildCache(ctx context.Context) error {

	// Lock everything because we're rewriting the whole cache
	v.vlcLock.Lock()
	defer v.vlcLock.Unlock()

	log.Info("Refreshing volume cache.")
	// if it's not in the cache, check the volumeStore, cache the result, and return the list.
	vols, err := v.volumeStore.VolumesList(ctx)
	if err != nil {
		return err
	}

	for _, vol := range vols {
		log.Infof("Volumestore: Found vol %s on store %s.", vol.ID, vol.Store)
		// Add it to the cache.
		v.vlc[vol.ID] = *vol
	}

	return nil
}
