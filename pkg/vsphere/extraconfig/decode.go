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
	"os"
	"reflect"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/vim25/types"
)

var (
	DecodeLogLevel = log.InfoLevel
)

type decoder func(src DataSource, dest reflect.Value, prefix string) reflect.Value

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

	intfDecoders = map[reflect.Type]decoder{}
}

// decode is the generic switcher that decides which decoder to use for a field
func decode(src DataSource, dest reflect.Value, prefix string) reflect.Value {
	// obtain the handler from the map, checking for the more specific interfaces first
	dec, ok := intfDecoders[dest.Type()]
	if ok {
		return dec(src, dest, prefix)
	}

	dec, ok = kindDecoders[dest.Kind()]
	if ok {
		return dec(src, dest, prefix)
	}

	log.Debugf("Skipping unsupported field, interface: %T, kind %s", dest, dest.Kind())
	return dest
}

// decodeString is the degenerative case where what we get is what we need
func decodeString(src DataSource, dest reflect.Value, prefix string) reflect.Value {

	value, err := src(prefix)
	if err != nil {
		log.Debugf("No value found in data source for string at key \"%s\"", prefix)
	}

	return reflect.ValueOf(value)
}

// decodePrimitive wraps the fromString primitive decoding in a manner that can be called via decode
func decodePrimitive(src DataSource, dest reflect.Value, prefix string) reflect.Value {
	// see if there's a value to decode
	v, err := src(prefix)
	if err != nil {
		log.Debugf("No value available for key to primitive %s", prefix)
		return dest
	}

	t := dest.Type()
	value := fromString(reflect.Zero(t), v)
	dest.Set(value)

	return dest
}

func decodePtr(src DataSource, dest reflect.Value, prefix string) reflect.Value {
	// value representing the run-time data
	log.Debugf("Decoding struct into object: %#v", dest)

	// if the pointer is nil we need to create the destination type
	var target reflect.Value
	if dest.IsNil() {
		target = reflect.New(dest.Type().Elem())
	} else {
		target = dest.Elem()
	}

	// check to see if the resulting object is not nil
	// If it is nil, then there was nothing to decode and the pointer remains nil
	result := decode(src, target, prefix)
	log.Debugf("target is now %#v, %+q ", target, target.Type())
	if result.IsValid() && result != reflect.Zero(dest.Type()) {
		if result.CanAddr() {
			dest.Elem().Set(result)
		} else {
			// no idea why some things come out of here unaddressable while others don't.
			// The main item that hits this path is TestNetPointer
			return result
		}
	}

	return dest
}

func decodeStruct(src DataSource, dest reflect.Value, prefix string) reflect.Value {
	// value representing the run-time data
	log.Debugf("Decoding struct into object: %#v", dest)

	// if the value is the zero type, we have to create ourselves
	var this reflect.Value
	if !dest.IsValid() {
		log.Debugf("Making new struct for %s", prefix)
		this = reflect.New(dest.Type())
	} else {
		this = dest
	}

	// do we have any data for this struct at all
	var valid bool

	// iterate through every field in the struct
	for i := 0; i < this.NumField(); i++ {
		field := this.Field(i)
		key := calculateKeyFromField(this.Type().Field(i), prefix)

		// Dump what we have so far
		log.Debugf("Key: %s, Kind: %s Value: %s", key, field.Kind(), field.String())

		// check to see if the resulting object is not nil
		// If it is nil, then there was nothing to decode
		result := decode(src, field, key)
		if result.IsValid() {
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

func decodeSlice(src DataSource, dest reflect.Value, prefix string) reflect.Value {
	// value representing the run-time data
	log.Debugf("Decoding struct into object: %#v", dest)
	kind := dest.Type().Elem().Kind()

	if kind == reflect.Uint8 {
		// special []byte array handling
		log.Debugf("Converting string to []byte")
		bytes, err := src(prefix)
		if err != nil {
			log.Debugf("No value found in data source for []byte \"%s\"", prefix)
			return dest
		}

		// couldn't use this.SetBytes for some unaddressable reason
		//this.SetBytes([]byte(bytes))
		return reflect.ValueOf([]byte(bytes))
	}

	// do we have any data for this struct at all
	length := 0

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
	}

	// determine the key given the array type
	curLen := this.Len()
	if kind == reflect.Struct {
		for i := 0; i < length; i++ {
			// convert key to name|index format
			key := fmt.Sprintf("%s|%d", prefix, i)

			// if there's already a struct in the array at this index then we pass that as the current
			// value
			var cur reflect.Value
			if i < curLen {
				cur = this.Index(i)
			}

			result := decode(src, cur, key)
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

func decodeMap(src DataSource, dest reflect.Value, prefix string) reflect.Value {
	// value representing the run-time data
	log.Debugf("Decoding struct into object: %#v", dest)

	// if the value is the zero type, we have to create ourselves
	var this reflect.Value
	if !dest.IsValid() {
		log.Debugf("Making new maps for %s", prefix)
		this = reflect.MakeMap(dest.Type())
	} else {
		this = dest
	}

	mapkeys, err := src(prefix)
	if err != nil {
		log.Debugf("No value found in data source for maps keys \"%s\"", prefix)
	}

	keytype := this.Type().Key()

	// split the list of map keys and iterate
	for _, value := range strings.Split(mapkeys, "|") {
		k := fromString(reflect.Zero(keytype), value)
		target := this.MapIndex(k)

		key := fmt.Sprintf("%s|%s", prefix, value)

		// check to see if the resulting object is not nil
		// If it is nil, then there was nothing to decode and the pointer remains nil
		result := decode(src, target, key)
		if result.IsValid() {
			this.SetMapIndex(k, result)
		}
	}

	return this
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

// DataSource provides a function that, give a key will return a value
// this is to be used during extraConfig decode to obtain values
type DataSource func(string) (string, error)

// Decode convert given type to []types.BaseOptionValue
func Decode(src DataSource, dest interface{}) interface{} {
	defer log.SetLevel(log.GetLevel())
	log.SetLevel(DecodeLogLevel)

	value := decode(src, reflect.ValueOf(dest), DefaultPrefix)

	return value.Interface()
}

// MapSource takes a key/value map and uses that as the datasource for decoding into
// target structures
func MapSource(src map[string]string) DataSource {
	return func(key string) (string, error) {
		val, ok := src[key]
		if !ok {
			return "", os.ErrNotExist
		}
		return val, nil
	}
}

// OptionValueSource is a convenience method to generate a MapSource source from
// and array of OptionValue's
func OptionValueSource(src []types.BaseOptionValue) DataSource {
	// create the key/value store from the extraconfig slice for lookups
	kv := make(map[string]string)
	for i := range src {
		k := src[i].GetOptionValue().Key
		v := src[i].GetOptionValue().Value.(string)
		kv[k] = v
	}

	return MapSource(kv)
}
