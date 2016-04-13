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

package auth

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

const (
	usernameEnvVar = "vic_username"
	passwordEnvVar = "vic_password"
)

// Authenticator provides a contract for authenticating HTTP handlers
type Authenticator interface {

	// Authenticate should wrap an http.HandlerFunc with implementation-specific authentication
	Authenticate(http.HandlerFunc) http.HandlerFunc
}

// None is an empty struct for disabling authentication
type None struct{}

// Authenticate when called on None is a noop
func (n *None) Authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return handler
}

// BasicHTTP provides metadata for HTTP Basic Authentication
type BasicHTTP struct {
	username string
	password string
}

// NewBasicHTTP acts as a constructor for a BasicHTTP authenticator
// Invoking it with empty arguments will pull values from environment variables
func NewBasicHTTP(username, password string) *BasicHTTP {
	var b *BasicHTTP
	if len(username) == 0 || len(password) == 0 {
		b = &BasicHTTP{
			username: os.Getenv(usernameEnvVar),
			password: os.Getenv(passwordEnvVar),
		}

		if len(b.username) == 0 || len(b.password) == 0 {
			log.Fatalf("You're attempting to use Basic HTTP Authentication but you do not have a username "+
				"and/or password set in your environment. Please export the variables %s and %s with a non-empty value "+
				"in order to enable Basic HTTP Authentication.",
				usernameEnvVar, passwordEnvVar)
		}

	} else {
		b = &BasicHTTP{
			username: username,
			password: password,
		}
	}
	return b
}

// Checks credentials for HTTP Basic Authentication
func (b *BasicHTTP) validateCredentials(u string, p string) bool {
	return u == b.username && p == b.password
}

// Authenticate for BasicHTTP wraps handler in Basic HTTP Authentication
func (b *BasicHTTP) Authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || !b.validateCredentials(user, password) {
			w.Header().Add("WWW-Authenticate", "Basic realm=vicadmin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}
