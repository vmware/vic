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
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/vim25/types"
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

var (
	// EncodeLogLevel value
	EncodeLogLevel = log.InfoLevel
)

// toString converts a basic type to its string representation
func toString(field reflect.Value) string {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(field.Bool())
	case reflect.String:
		return field.String()
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'E', -1, 64)
	default:
		return ""
	}
}

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

// calculateKey calculates the key based on the scope and prefix
func calculateKey(scopes []string, prefix string, key string) string {
	scope := calculateScope(scopes)
	if scope&Invalid != 0 || scope&NonPersistent != 0 {
		return ""
	}

	var guestinfo string
	var seperator string

	// Trim the whitespaces
	key = strings.TrimSpace(key)

	if scope&Hidden != 0 {
		guestinfo = ""
		seperator = "~"
	}
	if scope&ReadOnly != 0 {
		guestinfo = DefaultGuestInfoPrefix
		seperator = "/"
	}
	if scope&ReadWrite != 0 {
		guestinfo = DefaultGuestInfoPrefix
		seperator = "."
	}

	// no need to add another DefaultGuestInfoPrefix as a prefix if we already have one
	if strings.Contains(prefix, DefaultGuestInfoPrefix) {
		guestinfo = ""
	}

	if guestinfo == "" && prefix == "" {
		return key
	}

	if guestinfo == "" {
		return strings.Join([]string{prefix, key}, seperator)
	}

	if prefix == "" {
		return strings.Join([]string{guestinfo, key}, seperator)
	}

	return strings.Join([]string{guestinfo, prefix, key}, seperator)
}

func encodeWithPrefix(src interface{}, prefix string) []types.BaseOptionValue {
	var config []types.BaseOptionValue

	// value representing the run-time data
	value := reflect.ValueOf(src)
	log.Debugf("Value: %#v", value)

	// determine the kind as it changes how we get underlying data
	switch value.Kind() {
	case reflect.Invalid:
		log.Errorf("Invalid Kind: %#v", value)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.String, reflect.Float32, reflect.Float64:
		// handle basic types directly
		log.Debugf("Basic type: %s", prefix)

		return append(config, &types.OptionValue{Key: prefix, Value: toString(value)})
	case reflect.Slice, reflect.Array:
		log.Debugf("Slice BEGIN: %s", prefix)

		var values []string
		// iterate over slice to find what's underneath
		for i := 0; i < value.Len(); i++ {
			// recurse if we have a struct type
			// otherwise add its value to values
			if value.Index(i).Kind() == reflect.Struct {
				key := fmt.Sprintf("%s|%d", prefix, i)
				config = append(config, encodeWithPrefix(value.Index(i).Interface(), key)...)
			} else {
				values = append(values, toString(value.Index(i)))
			}
		}
		// set the slice key with the values seperated by |
		if len(values) > 0 {
			// sort the values before joining
			sort.Strings(values)
			// prefix~ contains the items
			config = append(config, &types.OptionValue{Key: fmt.Sprintf("%s~", prefix), Value: strings.Join(values, "|")})
		}
		// prefix contains the count
		config = append(config, &types.OptionValue{Key: prefix, Value: fmt.Sprintf("%d", value.Len()-1)})

		log.Debugf("Slice END %s %s", prefix, values)

		return config

	case reflect.Struct:
		// iterate through every struct type field
		for i := 0; i < value.NumField(); i++ {
			// get the underlying struct field
			structField := value.Type().Field(i)
			//skip unexported fields
			if structField.PkgPath != "" {
				log.Debugf("Skipping %s (not exported)", structField.Name)
				continue
			}

			// get the annotations
			tags := structField.Tag
			log.Debugf("Tags: %#v", tags)

			// do we have DefaultTagName?
			if tags.Get(DefaultTagName) == "" {
				log.Debugf("Skipping %s (not vic)", structField.Name)
				continue
			}

			// get the scopes
			scopes := strings.Split(tags.Get("scope"), ",")
			log.Debugf("Scopes: %#v", scopes)

			// get the keys and split properties from it
			keys := strings.Split(tags.Get("key"), ",")
			key, properties := keys[0], keys[1:]
			log.Debugf("Keys: %#v Properties: %#v", key, properties)

			// returns the i'th field of the struct src
			field := value.Field(i)

			// process properties
			if len(properties) > 0 {
				// do not recurse into nested type
				if properties[0] == "omitnested" {
					config = append(config, &types.OptionValue{Key: key, Value: toString(field)})
					log.Debugf("Skipping %s (omitnested)", key)
					continue
				}

				log.Debugf("Unknown property %s (%s)", key, properties[0])
				continue
			}

			// re-calculate the key based on the scope and prefix
			if key = calculateKey(scopes, prefix, key); key == "" {
				log.Debugf("Skipping %s (unknown scope %s)", key, scopes)
				continue
			}

			// Dump what we have so far
			log.Debugf("Key: %s, Properties: %q Kind: %s Value: %s", key, properties, field.Kind(), field.String())

			// all check passed, start to work on i'th field of the struct src
			switch field.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array:
				// recurse for struct, slice and array types
				log.Debugf("Struct|Slice|Array begin")

				// type cast struct to supported types (eg; time.Time, url.Url etc.)
				switch field.Interface().(type) {
				case time.Time:
					config = append(config, &types.OptionValue{Key: key, Value: field.Interface().(time.Time).String()})
				default:
					config = append(config, encodeWithPrefix(field.Interface(), key)...)
				}
				log.Debugf("Struct|Slice|Array end")
				continue

			case reflect.Map:
				log.Debugf("Map begin")

				var keys []string
				// iterate over keys and recurse
				for _, v := range field.MapKeys() {
					keys = append(keys, toString(v))
					key := fmt.Sprintf("%s|%s", key, toString(v))
					config = append(config, encodeWithPrefix(field.MapIndex(v).Interface(), key)...)
				}
				// sort the keys before joining
				sort.Strings(keys)
				config = append(config, &types.OptionValue{Key: key, Value: strings.Join(keys, "|")})

				log.Debugf("Map end")
				continue

			case reflect.Ptr:
				log.Debugf("Prt begin")

				if field.IsNil() {
					log.Debugf("Skipping nil pointer")
					continue
				}

				// type cast struct to supported types (eg; time.Time, url.Url etc.) or recurse
				switch field.Elem().Interface().(type) {
				case time.Time:
					config = append(config, &types.OptionValue{Key: key, Value: field.Interface().(*time.Time).String()})
				default:
					// follow the pointer and recurse
					config = append(config, encodeWithPrefix(field.Elem().Interface(), key)...)

					if len(config) > 0 {
						// add the ptr itself
						config = append(config, &types.OptionValue{Key: key, Value: config[len(config)-1].GetOptionValue().Key})
					}
				}

				log.Debugf("Ptr end")
				continue
			}
			// otherwise add it directly
			config = append(config, &types.OptionValue{Key: key, Value: toString(field)})
		}
		return config

	default:
		log.Debugf("Skipping not supported kind %s", value.Kind())
		return nil
	}
}

// Encode converts given type to []types.BaseOptionValue slice
func Encode(src interface{}) []types.BaseOptionValue {
	log.SetLevel(EncodeLogLevel)

	return encodeWithPrefix(src, DefaultPrefix)
}
