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

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/vim25/types"
)

var (
	// EncodeLogLevel value
	EncodeLogLevel = log.InfoLevel
)

type encoder func(sink DataSink, src reflect.Value, prefix string)

var kindEncoders map[reflect.Kind]encoder
var intfEncoders map[reflect.Type]encoder

func init() {
	kindEncoders = map[reflect.Kind]encoder{
		reflect.String:  encodeString,
		reflect.Struct:  encodeStruct,
		reflect.Slice:   encodeSlice,
		reflect.Array:   encodeSlice,
		reflect.Map:     encodeMap,
		reflect.Ptr:     encodePtr,
		reflect.Int:     encodePrimitive,
		reflect.Int8:    encodePrimitive,
		reflect.Int16:   encodePrimitive,
		reflect.Int32:   encodePrimitive,
		reflect.Int64:   encodePrimitive,
		reflect.Bool:    encodePrimitive,
		reflect.Float32: encodePrimitive,
		reflect.Float64: encodePrimitive,
	}

	intfEncoders = map[reflect.Type]encoder{}
}

// decode is the generic switcher that decides which decoder to use for a field
func encode(sink DataSink, src reflect.Value, prefix string) {
	// obtain the handler from the map, checking for the more specific interfaces first
	dec, ok := intfEncoders[src.Type()]
	if ok {
		dec(sink, src, prefix)
		return
	}

	dec, ok = kindEncoders[src.Kind()]
	if ok {
		dec(sink, src, prefix)
		return
	}

	log.Debugf("Skipping unsupported field, interface: %T, kind %s", src, src.Kind())
}

// encodeString is the degenerative case where what we get is what we need
func encodeString(sink DataSink, src reflect.Value, prefix string) {
	sink(prefix, src.String())
}

// encodePrimitive wraps the toString primitive encoding in a manner that can be called via encode
func encodePrimitive(sink DataSink, src reflect.Value, prefix string) {
	sink(prefix, toString(src))
}

func encodePtr(sink DataSink, src reflect.Value, prefix string) {
	log.Debugf("Encoding object: %#v", src)

	if src.IsNil() {
		// no need to attempt anything
		return
	}

	encode(sink, src.Elem(), prefix)
}

func encodeStruct(sink DataSink, src reflect.Value, prefix string) {
	log.Debugf("Encoding object: %#v", src)

	// iterate through every field in the struct
	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		key := calculateKeyFromField(src.Type().Field(i), prefix)
		if key == "" {
			log.Debugf("Skipping field %s with empty computed key", src.Type().Field(i).Name)
			continue
		}

		// Dump what we have so far
		log.Debugf("Key: %s, Kind: %s Value: %s", key, field.Kind(), field.String())

		encode(sink, field, key)
	}
}

func encodeSlice(sink DataSink, src reflect.Value, prefix string) {
	log.Debugf("Encoding object: %#v", src)

	length := src.Len()
	if length == 0 {
		log.Debug("Skipping empty slice")
		return
	}

	// determine the key given the array type
	kind := src.Type().Elem().Kind()
	if kind == reflect.Uint8 {
		// special []byte array handling

		log.Debugf("Converting []byte to string")
		encode(sink, src.Convert(reflect.TypeOf("")), prefix)
		return

	} else if kind != reflect.Struct {
		// else assume it's primitive - we'll panic/recover and continue it not
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("unable to encode %s (slice): %s", src.Type(), err)
			}
		}()

		values := make([]string, length)
		for i := 0; i < length; i++ {
			values[i] = toString(src.Index(i))
		}

		// convert key to name|index format
		key := fmt.Sprintf("%s~", prefix)
		sink(key, strings.Join(values, "|"))
	} else {
		for i := 0; i < length; i++ {
			// convert key to name|index format
			key := fmt.Sprintf("%s|%d", prefix, i)
			encode(sink, src.Index(i), key)
		}
	}

	// prefix contains the length of the array
	// seems insane calling toString(ValueOf(..)) but it means we're using the same path for everything
	sink(prefix, toString(reflect.ValueOf(length-1)))
}

func encodeMap(sink DataSink, src reflect.Value, prefix string) {
	log.Debugf("Encoding object: %#v", src)

	// iterate over keys and recurse
	mkeys := src.MapKeys()
	length := len(mkeys)
	if length == 0 {
		log.Debug("Skipping empty map")
		return
	}

	keys := make([]string, length)
	for i, v := range mkeys {
		keys[i] = toString(v)
		key := fmt.Sprintf("%s|%s", prefix, keys[i])
		encode(sink, src.MapIndex(v), key)
	}
	// sort the keys before joining
	sort.Strings(keys)
	sink(prefix, strings.Join(keys, "|"))
}

// toString converts a basic type to its string representation
func toString(field reflect.Value) string {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(field.Bool())
	case reflect.String:
		return field.String()
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'E', -1, 64)
	default:
		panic(field.Type().String() + " is an unhandled type")
	}
}

// DataSink provides a function that, given a key/value will persist that
// in some manner suited for later retrieval
type DataSink func(string, string)

// Encode convert given type to []types.BaseOptionValue
func Encode(sink DataSink, dest interface{}) {
	defer log.SetLevel(log.GetLevel())
	log.SetLevel(EncodeLogLevel)

	encode(sink, reflect.ValueOf(dest), DefaultPrefix)
}

// MapSink takes a map and populates it with key/value pairs from the encode
func MapSink(sink map[string]string) DataSink {
	return func(key, value string) {
		sink[key] = value
	}
}

// OptionValueFromMap is a convenience method to convert a map into a BaseOptionValue array
func OptionValueFromMap(data map[string]string) []types.BaseOptionValue {
	if len(data) == 0 {
		return nil
	}

	array := make([]types.BaseOptionValue, len(data))

	i := 0
	for k, v := range data {
		array[i] = &types.OptionValue{Key: k, Value: v}
		i++
	}

	return array
}
