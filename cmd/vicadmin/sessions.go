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

package main

import (
	"fmt"
	// webcontext "github.com/gorilla/context"
	vchconfig "github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/vsphere/session"
	"time"
)

type UserSessionManager interface {
	New()
	Delete()
	Get()
}

type AuthType int

const (
	vSphere AuthType = iota
	ESX
	cert
)

// UserSession acts as a container for sessions
type UserSession struct {
	vsphereSess *session.Session
	created     time.Time
	authType    AuthType
	conf        *vchconfig.VirtualContainerHostConfigSpec
}

// create a new UserSession
func New(a AuthType, conf *vchconfig.VirtualContainerHostConfigSpec) *UserSession {

	var u *UserSession
	u.created = time.Now()
	u.authType = a
	switch u.authType {
	case vSphere:
	case ESX:
	case cert:
	default:
	}

	return u
}

// Delete an UserSession
func (u *UserSession) Delete() {

}

func (u *UserSession) Get() {

}

// WatchDog is responsible for expiring vSphere sessions
func WatchDog() {
	var sessions []*UserSession

	fmt.Println(sessions)
}
