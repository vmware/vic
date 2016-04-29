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
	"os"
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

type decoder func(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value

var kindDecoders map[reflect.Kind]decoder
var intfDecoders map[reflect.Type]decoder

func init() {
	kindDecoders = map[reflect.Kind]decoder{
		reflect.String:  decodeString,
		reflect.Struct:  decodeStruct,
		reflect.Slice:   decodeSlice,
		reflect.Array:   decodeSlice,
		reflect.Map:     decodeMap,
		reflect.Ptr:     decodePtr,
		reflect.Int:     decodePrimitive,
		reflect.Int8:    decodePrimitive,
		reflect.Int16:   decodePrimitive,
		reflect.Int32:   decodePrimitive,
		reflect.Int64:   decodePrimitive,
		reflect.Bool:    decodePrimitive,
		reflect.Float32: decodePrimitive,
		reflect.Float64: decodePrimitive,
	}

	intfDecoders = map[reflect.Type]decoder{
		reflect.TypeOf(time.Time{}): decodeTime,
	}
}

// decode is the generic switcher that decides which decoder to use for a field
func decode(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {
	// if depth has reached zero, we skip decoding entirely
	if depth.depth == 0 {
		return dest
	}
	depth.depth--

	// obtain the handler from the map, checking for the more specific interfaces first
	dec, ok := intfDecoders[dest.Type()]
	if ok {
		return dec(src, dest, prefix, depth)
	}

	dec, ok = kindDecoders[dest.Kind()]
	if ok {
		return dec(src, dest, prefix, depth)
	}

	log.Debugf("Skipping unsupported field, interface: %T, kind %s", dest, dest.Kind())
	return dest
}

// decodeString is the degenerative case where what we get is what we need
func decodeString(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {

	value, err := src(prefix)
	if err != nil {
		log.Debugf("No value found in data source for string at key \"%s\"", prefix)
	}

	return reflect.ValueOf(value)
}

// decodePrimitive wraps the fromString primitive decoding in a manner that can be called via decode
func decodePrimitive(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {
	var this reflect.Value
	if !dest.CanAddr() {
		log.Debugf("Making new primitive for %s", prefix)
		ptr := reflect.New(dest.Type())
		this = ptr.Elem()
	} else {
		log.Debugf("Reusing existing struct for %s", prefix)
		this = dest
	}

	// see if there's a value to decode
	v, err := src(prefix)
	if err != nil {
		log.Debugf("No value available for key to primitive %s", prefix)
		return dest
	}

	t := this.Type()
	value := fromString(reflect.Zero(t), v)
	this.Set(value)

	return this
}

func decodePtr(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {
	// if we're not following pointers, then return immediately
	if !depth.follow {
		return dest
	}

	// value representing the run-time data
	log.Debugf("Decoding pointer into object: %#v", dest)

	// if the pointer is nil we need to create the destination type
	var target reflect.Value
	if dest.IsNil() {
		target = reflect.New(dest.Type().Elem())
	} else {
		target = dest.Elem()
	}

	// check to see if the resulting object is not nil
	// If it is nil, then there was nothing to decode and the pointer remains nil
	result := decode(src, target, prefix, depth)
	log.Debugf("target is now %#v, %+q ", target, target.Type())
	kind := result.Kind()
	if kind == reflect.Ptr {
		return result
	}

	if result == reflect.Zero(dest.Type()) || !result.IsValid() {
		// leave the pointer as nil if the result is zero type or invalid
		return dest
	}

	// neither pointer, nor zero
	// NOTE: if the returned result is not addressable this can panic - that generally
	// indicates an incorrect implementation of a decodeX method... those should always
	// return addressable Values. See decodeByteSlice as an example - this uses make([]byte)
	// rather than built in string(bytes) conversion specifically to get an addressable return
	dest.Elem().Set(result)

	return dest
}

var typeType = reflect.TypeOf((*reflect.Type)(nil)).Elem()

func decodeStruct(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {
	// value representing the run-time data
	log.Debugf("Decoding struct into object: %#v, type: %s", dest, dest.Type().Name())

	var this reflect.Value
	if !dest.CanAddr() {
		log.Debugf("Making new struct for %s", prefix)
		ptr := reflect.New(dest.Type())
		this = ptr.Elem()
	} else {
		log.Debugf("Reusing existing struct for %s", prefix)
		this = dest
	}

	// do we have any data for this struct at all
	var valid bool

	// iterate through every field in the struct
	for i := 0; i < this.NumField(); i++ {
		field := this.Field(i)
		key, fdepth := calculateKeyFromField(this.Type().Field(i), prefix, depth)
		if key == "" {
			// this is either a malformed key or explicitly skipped
			continue
		}

		// Dump what we have so far
		log.Debugf("Key: %s, Kind: %s Value: %s", key, field.Kind(), field.String())

		// check to see if the resulting object is not nil
		// If it is nil, then there was nothing to decode
		result := decode(src, field, key, fdepth)
		if result.IsValid() {
			log.Debugf("Setting field %s to %#v", this.Type().Field(i).Name, result)
			field.Set(result)
			valid = true
		} else {
			log.Debugf("Invalid result for field %s", this.Type().Field(i).Name)
		}
	}

	if !valid {
		return reflect.Value{}
	}

	log.Debugf("Return decoded structure for %s: %#v", prefix, this)
	return this
}

func decodeByteSlice(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {
	log.Debugf("Converting string to []byte")
	bytes, err := src(prefix)
	if err != nil {
		log.Debugf("No value found in data source for []byte \"%s\"", prefix)
		return dest
	}

	length := len(bytes)

	// we don't even try to merge byte arrays - no idea how to get append behaviour
	// correct with reflection
	// use make([]byte) rather than built in string(bytes) conversion to get an addressable return value
	log.Debugf("Making new slice for %s", prefix)
	this := make([]byte, length, length)

	copy(this, bytes)

	return reflect.ValueOf(this)
}

func decodeSlice(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {
	// value representing the run-time data
	log.Debugf("Decoding struct into object: %#v", dest)
	kind := dest.Type().Elem().Kind()

	if kind == reflect.Uint8 {
		return decodeByteSlice(src, dest, prefix, depth)
	}

	// do we have any data for this struct at all
	length := 0
	curLen := 0

	// get the length of the array
	len, err := src(prefix)
	if err != nil || len == "" {
		log.Debugf("No value available for key %s - will create empty array if needed", prefix)
	} else {
		// if there's any data at all then we can assume we need to be extant
		lengthValue := fromString(reflect.ValueOf(0), len)
		length = int(lengthValue.Int()) + 1
	}

	var this reflect.Value
	if !dest.IsValid() || dest.IsNil() {
		log.Debugf("Making new slice for %s", prefix)
		this = reflect.MakeSlice(dest.Type(), length, length)
	} else {
		this = dest
		this.SetLen(length)
		curLen = this.Len()
	}

	// determine the key given the array type
	if kind == reflect.Struct {
		for i := 0; i < length; i++ {
			// convert key to name|index format
			key := fmt.Sprintf("%s|%d", prefix, i)

			// if there's already a struct in the array at this index then we pass that as the current
			// value
			var cur reflect.Value
			if i < curLen {
				cur = this.Index(i)
			} else {
				cur = reflect.Zero(dest.Type().Elem())
			}

			result := decode(src, cur, key, depth)
			if result.IsValid() {
				this.Index(i).Set(result)
			}

			continue
		}

		return this
	}

	// convert key to name|index format
	key := fmt.Sprintf("%s~", prefix)
	kval, err := src(key)
	if err != nil {
		log.Debugf("No value found in data source for key \"%s\"", key)
		return this
	}

	// lookup the key and split it
	values := strings.Split(kval, "|")
	for i := 0; i < length; i++ {
		v := values[i]
		t := this.Type().Elem()
		k := fromString(reflect.Zero(t), v)
		// set the i'th slice item
		this.Index(i).Set(k)
	}

	return this
}

func decodeMap(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {
	// value representing the run-time data
	log.Debugf("Decoding struct into object: %#v", dest)

	// if the value is the zero type, we have to create ourselves
	var this reflect.Value
	if !dest.IsValid() || dest.IsNil() {
		log.Debugf("Making new maps for %s", prefix)
		this = reflect.MakeMap(dest.Type())
	} else {
		this = dest
	}

	mapkeys, err := src(prefix)
	if err != nil {
		log.Debugf("No value found in data source for maps keys \"%s\"", prefix)
		return this
	}

	keytype := this.Type().Key()
	valtype := this.Type().Elem()

	// split the list of map keys and iterate
	for _, value := range strings.Split(mapkeys, "|") {
		k := fromString(reflect.Zero(keytype), value)
		target := this.MapIndex(k)
		if !target.IsValid() {
			target = reflect.Zero(valtype)
		}

		key := fmt.Sprintf("%s|%s", prefix, value)

		// check to see if the resulting object is not nil
		// If it is nil, then there was nothing to decode and the pointer remains nil
		result := decode(src, target, key, depth)
		if result.IsValid() {
			this.SetMapIndex(k, result)
		}
	}

	return this
}

func decodeTime(src DataSource, dest reflect.Value, prefix string, depth recursion) reflect.Value {
	v, err := src(prefix)
	if err != nil {
		log.Debugf("No value found in data source for time \"%s\"", prefix)
	}

	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", v)
	if err != nil {
		log.Debugf("Failed to convert value \"%s\" to time", v)
	}

	return reflect.ValueOf(t)
}

// fromString converts string representation of a basic type to basic type
func fromString(field reflect.Value, value string) reflect.Value {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Errorf("Failed to convert value %#v (%s) to int: %s", value, field.Kind(), err.Error())
			return field
		}
		return reflect.ValueOf(s).Convert(field.Type())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			log.Errorf("Failed to convert value %#v (%s) to uint: %s", value, field.Kind(), err.Error())
			return field
		}
		return reflect.ValueOf(s).Convert(field.Type())

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
func decodeWithPrefix(src func(string) (string, error), dest interface{}, prefix string) interface{} {
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

// Decode convert given type to []types.BaseOptionValue
func Decode(src DataSource, dest interface{}) interface{} {
	defer log.SetLevel(log.GetLevel())
	log.SetLevel(DecodeLogLevel)

	value := decode(src, reflect.ValueOf(dest), DefaultPrefix, Unbounded)

			// ensure we have the tag
			v, err := src(key)
			if err == nil {
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
					decodeWithPrefix(src, member.Interface(), key)
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
						key := fmt.Sprintf("%s|%d", key, i)
						decodeWithPrefix(src, member.Interface(), key)

						// set the i'th slice item
						slice.Index(i).Set(member.Elem())
					}
				} else {
					// convert key to name~ format
					key := fmt.Sprintf("%s~", key)
					kval, err := src(key)
					if err != nil {
						log.Debugf("No value found in data source for key \"%s\"", key)
						continue
					}
					// lookup the key and split it
					values := strings.Split(kval, "|")
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
					// if the map value is a pointer to a struct
					if field.Type().Elem().Kind() == reflect.Ptr && field.Type().Elem().Elem().Kind() == reflect.Struct {
						member := reflect.New(field.Type().Elem().Elem())

						key := fmt.Sprintf("%s|%s", key, value)
						decodeWithPrefix(src, member.Interface(), key)

						t := field.Type().Key()
						k := fromString(reflect.Zero(t), value)
						field.SetMapIndex(k, member.Elem().Addr())
					} else if field.Type().Elem().Kind() == reflect.Struct {
						member := reflect.New(field.Type().Elem())

						key := fmt.Sprintf("%s|%s", key, value)
						decodeWithPrefix(src, member.Interface(), key)

						t := field.Type().Key()
						k := fromString(reflect.Zero(t), value)
						field.SetMapIndex(k, member.Elem())
					} else {
						values := strings.Split(v, "|")
						for i := 0; i < len(values); i++ {
							v := values[i]
							// convert key to name|mapkey format
							key := fmt.Sprintf("%s|%s", key, v)
							kval, err := src(key)
							if err != nil {
								log.Debugf("No value found in data source for key \"%s\"", key)
								continue
							}

							t := field.Type().Elem()
							k := fromString(reflect.Zero(t), kval)

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
					kval, err := src(key)
					if err != nil {
						log.Debugf("No value found in data source for key \"%s\"", key)
						continue
					}
					if field.Type().Elem() == reflect.TypeOf(time.Time{}) {
						// from https://golang.org/src/time/format.go?s=12854:12883#L407
						t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", kval)
						if err != nil {
							log.Errorf("Failed to convert value %#v to time: %s", kval, err.Error())
						}
						// set the pointer
						field.Set(reflect.ValueOf(&t))
					} else if field.Type().Elem() == reflect.TypeOf(net.IPNet{}) {
						_, n, err := net.ParseCIDR(kval)
						if err != nil {
							log.Errorf("Failed to convert value %#v to IPNET: %s", kval, err.Error())
						}
						field.Set(reflect.ValueOf(n))
					} else {
						member := reflect.New(field.Type().Elem())

						decodeWithPrefix(src, member.Interface(), key)

						// set the pointer
						field.Set(member)
					}
				}
				log.Debugf("Ptr end")
			}
		}
		return val, nil
	}
}

// Decode convert given type to []types.BaseOptionValue
func Decode(src func(string) (string, error), dest interface{}) interface{} {
	log.SetLevel(DecodeLogLevel)

	return decodeWithPrefix(src, dest, DefaultPrefix)
}

// MapSource takes a key/value map and uses that as the datasource for decoding into
// target structures
func MapSource(src map[string]string) func(string) (string, error) {
	return func(key string) (val string, err error) {
		val, ok := src[key]
		if !ok {
			err = os.ErrNotExist
		}
		return
	}
}

// OptionValueSource is a convenience method to generate a MapSource source from
// and array of OptionValue's
func OptionValueSource(src []types.BaseOptionValue) func(string) (string, error) {
	// create the key/value store from the extraconfig slice for lookups
	kv := make(map[string]string)
	for i := range src {
		k := src[i].GetOptionValue().Key
		v := src[i].GetOptionValue().Value.(string)
		kv[k] = v
	}

	return MapSource(kv)
}
