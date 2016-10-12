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
