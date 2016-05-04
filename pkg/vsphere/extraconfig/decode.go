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
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/vim25/types"
)

var (
	DecodeLogLevel = log.InfoLevel
)

// fromString converts string representation of a basic type to basic type
func fromString(field reflect.Value, value string) reflect.Value {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Errorf("Failed to convert value %#v (%s) to int: %s", value, field.Kind(), err.Error())
			return field
		}
		return reflect.ValueOf(int(s))

	case reflect.Bool:
		s, err := strconv.ParseBool(value)
		if err != nil {
			log.Errorf("Failed to convert value %#v (%s) to bool: %s", value, field.Kind(), err.Error())
			return field
		}
		return reflect.ValueOf(s)

	case reflect.String:
		return reflect.ValueOf(value)

	case reflect.Float32, reflect.Float64:
		s, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Errorf("Failed to convert value %#v (%s) to float: %s", value, field.Kind(), err.Error())
			return field
		}
		return reflect.ValueOf(s)

	}
	log.Debugf("Invalid Kind: %s (%#v)", field.Kind(), value)

	return field
}

// decodeWithPrefix convert []types.BaseOptionValue to type dest
func decodeWithPrefix(kv map[string]string, dest interface{}, prefix string) interface{} {
	// value representing the run-time data
	value := reflect.ValueOf(dest).Elem()
	log.Debugf("Value: %#v", value)

	// determine the kind as it changes how we get underlying data
	switch value.Kind() {
	case reflect.Invalid:
		log.Errorf("Invalid Kind: %#v", value)
		return nil

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

			// returns the i'th field of the struct i
			field := value.Field(i)

			// process properties
			if len(properties) > 0 {
				// do not recurse into nested type
				if properties[0] == "omitnested" {
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

			// ensure we have the tag
			v, ok := kv[key]
			if ok {
				log.Debugf("Setting the %s with %s", key, v)
				// set the field with the value
				field.Set(fromString(field, v))
			}

			// all check passed, start to work on i'th field of the struct dest
			switch field.Kind() {
			case reflect.Struct:
				log.Debugf("Struct begin")

				// type cast struct to supported types (eg; time.Time, url.Url etc.)
				switch field.Interface().(type) {
				case time.Time:
					// from https://golang.org/src/time/format.go?s=12854:12883#L407
					t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", v)
					if err != nil {
						log.Errorf("Failed to convert value %#v to time: %s", v, err.Error())
					}
					field.Set(reflect.ValueOf(t))
				case net.IPNet:
					_, n, err := net.ParseCIDR(v)
					if err != nil {
						log.Errorf("Failed to convert value %#v to IPNET: %s", v, err.Error())
					}
					field.Set(reflect.ValueOf(*n))
				default:
					member := reflect.New(field.Type())
					decodeWithPrefix(kv, member.Interface(), key)
					field.Set(member.Elem())
				}
				log.Debugf("Struct end")

			case reflect.Slice, reflect.Array:
				// recurse for struct, slice and array types
				log.Debugf("Slice|Array begin")

				lengthValue := fromString(reflect.ValueOf(0), v)
				length := int(lengthValue.Int()) + 1

				// create the slice
				slice := reflect.MakeSlice(field.Type(), length, length)
				if field.Type().Elem().Kind() == reflect.Struct {
					for i := 0; i < length; i++ {
						member := reflect.New(field.Type().Elem())
						// convert key to name|index format
						mkey := fmt.Sprintf("%s|%d", key, i)
						decodeWithPrefix(kv, member.Interface(), mkey)

						// set the i'th slice item
						slice.Index(i).Set(member.Elem())
					}
				} else {
					// convert key to name~ format
					mkey := fmt.Sprintf("%s~", key)
					// lookup the key and split it
					values := strings.Split(kv[mkey], "|")
					for i := 0; i < length; i++ {
						x := values[i]
						t := field.Type().Elem()
						k := fromString(reflect.Zero(t), x)
						// set the i'th slice item
						slice.Index(i).Set(k)
					}
				}
				// set the slice
				field.Set(slice)

				log.Debugf("Slice|Array end")

			case reflect.Map:
				log.Debugf("Map begin")

				// create and set the map if neccessary
				if field.IsNil() {
					field.Set(reflect.MakeMap(field.Type()))
				}

				// split v and iterate over it
				for _, value := range strings.Split(v, "|") {
					if field.Type().Elem().Kind() == reflect.Struct {
						member := reflect.New(field.Type().Elem())

						mkey := fmt.Sprintf("%s|%s", key, value)
						decodeWithPrefix(kv, member.Interface(), mkey)

						t := field.Type().Key()
						k := fromString(reflect.Zero(t), value)
						field.SetMapIndex(k, member.Elem())
					} else {
						values := strings.Split(v, "|")
						for i := 0; i < len(values); i++ {
							v := values[i]
							// convert key to name|mapkey format
							mkey := fmt.Sprintf("%s|%s", key, v)
							t := field.Type().Elem()
							k := fromString(reflect.Zero(t), kv[mkey])

							// set the map item
							field.SetMapIndex(reflect.ValueOf(v), k)
						}
					}
				}

				log.Debugf("Map end")

			case reflect.Ptr:
				log.Debugf("Prt begin")
				// FIXME: non struct pointers
				if field.Type().Elem().Kind() == reflect.Struct {
					v, ok := kv[key]
					if ok {
						if field.Type().Elem() == reflect.TypeOf(time.Time{}) {
							// from https://golang.org/src/time/format.go?s=12854:12883#L407
							t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", v)
							if err != nil {
								log.Errorf("Failed to convert value %#v to time: %s", v, err.Error())
							}
							// set the pointer
							field.Set(reflect.ValueOf(&t))
						} else if field.Type().Elem() == reflect.TypeOf(net.IPNet{}) {
							_, n, err := net.ParseCIDR(v)
							if err != nil {
								log.Errorf("Failed to convert value %#v to IPNET: %s", v, err.Error())
							}
							field.Set(reflect.ValueOf(n))
						} else {
							member := reflect.New(field.Type().Elem())

							decodeWithPrefix(kv, member.Interface(), key)

							// set the pointer
							field.Set(member)
						}
					}
				}
				log.Debugf("Ptr end")
			}
		}
		return dest
	default:
		log.Debugf("Skipping not supported kind %s", value.Kind())
		return nil
	}
}

// Decode convert given type to []types.BaseOptionValue
func Decode(src []types.BaseOptionValue, dest interface{}) interface{} {
	log.SetLevel(DecodeLogLevel)

	// create the key/value store from the extraconfig slice for lookups
	kv := make(map[string]string)
	for i := range src {
		k := src[i].GetOptionValue().Key
		v := src[i].GetOptionValue().Value.(string)
		kv[k] = v
	}

	return decodeWithPrefix(kv, dest, DefaultPrefix)
}
