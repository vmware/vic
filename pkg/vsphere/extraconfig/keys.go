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

package extraconfig

import (
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	// DefaultTagName value
	DefaultTagName = "vic"
	// DefaultPrefix value
	DefaultPrefix = ""
	// DefaultGuestInfoPrefix value
	DefaultGuestInfoPrefix = "guestinfo"
)

const (
	// Invalid value
	Invalid = 1 << iota
	// Hidden value
	Hidden
	// ReadOnly value
	ReadOnly
	// ReadWrite value
	ReadWrite
	// NonPersistent value
	NonPersistent
	// Volatile value
	Volatile
)

// calculateScope returns the uint representation of scope tag
func calculateScope(scopes []string) uint {
	var scope uint
	for _, v := range scopes {
		switch v {
		case "hidden":
			scope |= Hidden
		case "read-only":
			scope |= ReadOnly
		case "read-write":
			scope |= ReadWrite
		case "non-persistent":
			scope |= NonPersistent
		case "volatile":
			scope |= Volatile
		default:
			return Invalid
		}
	}
	return scope
}

func calculateScopeFromKey(key string) []string {
	scopes := []string{}

	if !strings.HasPrefix(key, DefaultGuestInfoPrefix) {
		scopes = append(scopes, "hidden")
	}

	if strings.Contains(key, "/") {
		scopes = append(scopes, "read-only")
	} else {
		scopes = append(scopes, "read-write")
	}

	return scopes
}

func calculateKeyFromField(field reflect.StructField, prefix string) string {
	//skip unexported fields
	if field.PkgPath != "" {
		log.Debugf("Skipping %s (not exported)", field.Name)
		return ""
	}

	// get the annotations
	tags := field.Tag
	log.Debugf("Tags: %#v", tags)

	var key string
	var scopes []string

	// do we have DefaultTagName?
	if tags.Get(DefaultTagName) == "" {
		log.Debugf("%s not tagged - inheriting parent scope", field.Name)

		key = field.Name
		scopes = calculateScopeFromKey(prefix)
	} else {
		// get the scopes
		scopes = strings.Split(tags.Get("scope"), ",")
		log.Debugf("Scopes: %#v", scopes)

		// get the keys and split properties from it
		keys := strings.Split(tags.Get("key"), ",")
		key = keys[0]
		properties := keys[1:]
		log.Debugf("Keys: %#v Properties: %#v", key, properties)

		// process properties
		if len(properties) > 0 {
			// do not recurse into nested type
			if properties[0] == "omitnested" {
				log.Debugf("Skipping %s (omitnested)", key)
				return ""
			}

			log.Debugf("Unknown property %s (%s)", key, properties[0])
			return ""
		}
	}

	// re-calculate the key based on the scope and prefix
	if key = calculateKey(scopes, prefix, key); key == "" {
		log.Debugf("Skipping %s (unknown scope %s)", key, scopes)
		return ""
	}

	return key
}

// calculateKey calculates the key based on the scope and prefix
func calculateKey(scopes []string, prefix string, key string) string {
	scope := calculateScope(scopes)
	if scope&Invalid != 0 || scope&NonPersistent != 0 {
		return ""
	}

	var guestinfo string
	var seperator string
	var oseperator string

	// Trim the whitespaces
	key = strings.TrimSpace(key)

	if scope&Hidden != 0 {
		guestinfo = ""
		seperator = "/"
		oseperator = "."
	}
	if scope&ReadOnly != 0 {
		guestinfo = DefaultGuestInfoPrefix
		seperator = "/"
		oseperator = "."
	}
	if scope&ReadWrite != 0 {
		guestinfo = DefaultGuestInfoPrefix
		seperator = "."
		oseperator = "/"
	}

	// set up the correct base prefix first
	base := strings.Replace(prefix, oseperator, seperator, -1)

	if strings.HasPrefix(base, DefaultGuestInfoPrefix) {
		if guestinfo == "" {
			// this key is hidden - strip the prefix and separator
			base = base[len(DefaultGuestInfoPrefix)+1:]
		} else {
			// this key is already exposed
			guestinfo = ""
		}
	}

	// the add detail to the base
	if guestinfo == "" && prefix == "" {
		return key
	}

	if guestinfo == "" {
		return strings.Join([]string{base, key}, seperator)
	}

	if prefix == "" {
		return strings.Join([]string{guestinfo, key}, seperator)
	}

	return strings.Join([]string{guestinfo, prefix, key}, seperator)
}
