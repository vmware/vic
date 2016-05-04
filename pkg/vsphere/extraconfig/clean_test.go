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
	"net"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// make it verbose during testing
func init() {
	DecodeLogLevel = log.DebugLevel
	EncodeLogLevel = log.DebugLevel
}

func TestEmbedded(t *testing.T) {

	type Type struct {
		Common `vic:"0.1" scope:"read-only" key:"common"`
	}

	Embedded := Type{
		Common: Common{
			ID:   "0xDEADBEEF",
			Name: "Embedded",
		},
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), Embedded)

	expected := map[string]string{
		"guestinfo/common/id":    "0xDEADBEEF",
		"guestinfo/common/name":  "Embedded",
		"guestinfo/common/notes": "",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Type
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, Embedded, decoded, "Encoded and decoded does not match")
}

func TestNetPointer(t *testing.T) {
	type Type struct {
		Net *net.IPNet `vic:"0.1" scope:"read-only" key:"net"`
	}

	// 127.0.0.1/8
	n := net.IPNet{IP: net.IP{0x7f, 0x0, 0x0, 0x1}, Mask: net.IPMask{0xff, 0x0, 0x0, 0x0}}
	Net := Type{
		Net: &n,
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), Net)

	expected := map[string]string{
		"guestinfo/net/IP":   "\u007f\x00\x00\x01",
		"guestinfo/net/Mask": "\xff\x00\x00\x00",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Type
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, Net, decoded, "Encoded and decoded does not match")
}

func TestTimePointer(t *testing.T) {
	d := time.Date(2009, 11, 10, 23, 00, 00, 0, time.UTC)

	type Type struct {
		Time *time.Time `vic:"0.1" scope:"read-only" key:"time"`
	}

	Time := Type{
		Time: &d,
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), Time)

	expected := map[string]string{
		"guestinfo/time": "2009-11-10 23:00:00 +0000 UTC",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Type
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, Time, decoded, "Encoded and decoded does not match")
}

func TestStructMap(t *testing.T) {
	type Type struct {
		StructMap map[string]Common `vic:"0.1" scope:"read-only" key:"map"`
	}

	StructMap := Type{
		map[string]Common{
			"Key1": Common{
				ID:   "0xDEADBEEF",
				Name: "beef",
			},
			"Key2": Common{
				ID:   "0x8BADF00D",
				Name: "food",
			},
			"Key3": Common{
				ID:   "0xDEADF00D",
				Name: "dead",
			},
		},
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), StructMap)

	expected := map[string]string{
		"guestinfo/map|Key1/id":    "0xDEADBEEF",
		"guestinfo/map|Key1/name":  "beef",
		"guestinfo/map|Key1/notes": "",
		"guestinfo/map|Key2/id":    "0x8BADF00D",
		"guestinfo/map|Key2/name":  "food",
		"guestinfo/map|Key2/notes": "",
		"guestinfo/map|Key3/id":    "0xDEADF00D",
		"guestinfo/map|Key3/name":  "dead",
		"guestinfo/map|Key3/notes": "",
		"guestinfo/map":            "Key1|Key2|Key3",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Type
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, StructMap, decoded, "Encoded and decoded does not match")
}

func TestIntStructMap(t *testing.T) {
	type Type struct {
		StructMap map[int]Common `vic:"0.1" scope:"read-only" key:"map"`
	}

	StructMap := Type{
		map[int]Common{
			1: Common{
				ID:   "0xDEADBEEF",
				Name: "beef",
			},
			2: Common{
				ID:   "0x8BADF00D",
				Name: "food",
			},
			3: Common{
				ID:   "0xDEADF00D",
				Name: "dead",
			},
		},
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), StructMap)

	expected := map[string]string{
		"guestinfo/map|1/id":    "0xDEADBEEF",
		"guestinfo/map|1/name":  "beef",
		"guestinfo/map|1/notes": "",
		"guestinfo/map|2/id":    "0x8BADF00D",
		"guestinfo/map|2/name":  "food",
		"guestinfo/map|2/notes": "",
		"guestinfo/map|3/id":    "0xDEADF00D",
		"guestinfo/map|3/name":  "dead",
		"guestinfo/map|3/notes": "",
		"guestinfo/map":         "1|2|3",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Type
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, StructMap, decoded, "Encoded and decoded does not match")
}

func TestStructSlice(t *testing.T) {
	type Type struct {
		StructSlice []Common `vic:"0.1" scope:"read-only" key:"slice"`
	}

	StructSlice := Type{
		[]Common{
			Common{
				ID:   "0xDEADFEED",
				Name: "feed",
			},
			Common{
				ID:   "0xFACEFEED",
				Name: "face",
			},
		},
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), StructSlice)

	expected := map[string]string{
		"guestinfo/slice":         "1",
		"guestinfo/slice|0/id":    "0xDEADFEED",
		"guestinfo/slice|0/name":  "feed",
		"guestinfo/slice|0/notes": "",
		"guestinfo/slice|1/id":    "0xFACEFEED",
		"guestinfo/slice|1/name":  "face",
		"guestinfo/slice|1/notes": "",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Type
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, StructSlice, decoded, "Encoded and decoded does not match")
}

func TestMultipleScope(t *testing.T) {
	MultipleScope := struct {
		MultipleScope string `vic:"0.1" scope:"read-only,hidden,non-persistent" key:"multiscope"`
	}{
		"MultipleScope",
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), MultipleScope)

	expected := map[string]string{}
	assert.Equal(t, expected, encoded, "Not equal")
}

func TestUnknownScope(t *testing.T) {
	UnknownScope := struct {
		UnknownScope int `vic:"0.1" scope:"unknownscope" key:"unknownscope"`
	}{
		42,
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), UnknownScope)

	expected := map[string]string{}
	assert.Equal(t, encoded, expected, "Not equal")
}

func TestUnknownProperty(t *testing.T) {
	UnknownProperty := struct {
		UnknownProperty int `vic:"0.1" scope:"hidden" key:"unknownproperty" recurse:"unknownproperty"`
	}{
		42,
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), UnknownProperty)

	expected := map[string]string{
		"unknownproperty": "42",
	}
	assert.Equal(t, expected, encoded, "Not equal")
}

func TestOmitNested(t *testing.T) {
	OmitNested := struct {
		Time        time.Time `vic:"0.1" scope:"volatile" key:"time" recurse:"depth=0"`
		CurrentTime time.Time `vic:"0.1" scope:"volatile" key:"time"`
	}{
		Time:        time.Date(2009, 11, 10, 23, 00, 00, 0, time.UTC),
		CurrentTime: time.Date(2009, 11, 10, 23, 00, 00, 0, time.UTC),
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), OmitNested)

	expected := map[string]string{
		"time": "2009-11-10 23:00:00 +0000 UTC",
	}
	assert.Equal(t, expected, encoded, "Encoded and decoded does not match")

}

func TestComplex(t *testing.T) {
	type Type struct {
		ExecutorConfig ExecutorConfig `vic:"0.1" scope:"hidden" key:"executorconfig"`
	}

	ExecutorConfig := Type{
		ExecutorConfig{
			Sessions: map[string]SessionConfig{
				"Session1": SessionConfig{
					Common: Common{
						ID:   "SessionID",
						Name: "SessionName",
					},
					Tty: true,
					Cmd: Cmd{
						Path: "/vmware",
						Args: []string{"/bin/imagec", "-standalone"},
						Env:  []string{"PATH=/bin", "USER=imagec"},
						Dir:  "/",
					},
				},
			},
		},
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), ExecutorConfig)

	expected := map[string]string{
		"guestinfo/executorconfig/common/id":                      "",
		"guestinfo/executorconfig/common/name":                    "",
		"guestinfo/executorconfig/common/notes":                   "",
		"guestinfo/executorconfig/sessions|Session1/common/id":    "SessionID",
		"guestinfo/executorconfig/sessions|Session1/common/name":  "SessionName",
		"guestinfo/executorconfig/sessions|Session1/common/notes": "",
		"executorconfig/sessions|Session1/cmd/path":               "/vmware",
		"executorconfig/sessions|Session1/cmd/args~":              "/bin/imagec|-standalone",
		"executorconfig/sessions|Session1/cmd/args":               "1",
		"executorconfig/sessions|Session1/cmd/env~":               "PATH=/bin|USER=imagec",
		"executorconfig/sessions|Session1/cmd/env":                "1",
		"executorconfig/sessions|Session1/cmd/dir":                "/",
		"executorconfig/sessions|Session1/tty":                    "true",
		"executorconfig/sessions":                                 "Session1",
		"guestinfo.executorconfig.Key":                            "",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Type
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, ExecutorConfig, decoded, "Encoded and decoded does not match")
}

func TestComplexPointer(t *testing.T) {
	type Type struct {
		ExecutorConfig *ExecutorConfig `vic:"0.1" scope:"hidden" key:"executorconfig"`
	}

	ExecutorConfig := Type{
		&ExecutorConfig{
			Sessions: map[string]SessionConfig{
				"Session1": SessionConfig{
					Common: Common{
						ID:   "SessionID",
						Name: "SessionName",
					},
					Tty: true,
					Cmd: Cmd{
						Path: "/vmware",
						Args: []string{"/bin/imagec", "-standalone"},
						Env:  []string{"PATH=/bin", "USER=imagec"},
						Dir:  "/",
					},
				},
			},
		},
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), ExecutorConfig)

	expected := map[string]string{
		"guestinfo/executorconfig/common/id":                      "",
		"guestinfo/executorconfig/common/name":                    "",
		"guestinfo/executorconfig/common/notes":                   "",
		"guestinfo/executorconfig/sessions|Session1/common/id":    "SessionID",
		"guestinfo/executorconfig/sessions|Session1/common/name":  "SessionName",
		"guestinfo/executorconfig/sessions|Session1/common/notes": "",
		"executorconfig/sessions|Session1/cmd/path":               "/vmware",
		"executorconfig/sessions|Session1/cmd/args~":              "/bin/imagec|-standalone",
		"executorconfig/sessions|Session1/cmd/args":               "1",
		"executorconfig/sessions|Session1/cmd/env~":               "PATH=/bin|USER=imagec",
		"executorconfig/sessions|Session1/cmd/env":                "1",
		"executorconfig/sessions|Session1/cmd/dir":                "/",
		"executorconfig/sessions|Session1/tty":                    "true",
		"executorconfig/sessions":                                 "Session1",
		"guestinfo.executorconfig.Key":                            "",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Type
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, ExecutorConfig, decoded, "Encoded and decoded does not match")
}

// TestPointerDecode tests the translation from a type where the sessions are direct values to
// one where they are pointers
func TestPointerDecode(t *testing.T) {
	reference := ExecutorConfig{
		Sessions: map[string]SessionConfig{
			"Session1": SessionConfig{
				Common: Common{
					ID:   "SessionID",
					Name: "SessionName",
				},
				Tty: true,
				Cmd: Cmd{
					Path: "/vmware",
					Args: []string{"/bin/imagec", "-standalone"},
					Env:  []string{"PATH=/bin", "USER=imagec"},
					Dir:  "/",
				},
			},
		},
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), reference)

	expected := map[string]string{
		"guestinfo/common/id":                      "",
		"guestinfo/common/name":                    "",
		"guestinfo/common/notes":                   "",
		"guestinfo/sessions|Session1/common/id":    "SessionID",
		"guestinfo/sessions|Session1/common/name":  "SessionName",
		"guestinfo/sessions|Session1/common/notes": "",
		"sessions|Session1/cmd/path":               "/vmware",
		"sessions|Session1/cmd/args~":              "/bin/imagec|-standalone",
		"sessions|Session1/cmd/args":               "1",
		"sessions|Session1/cmd/env~":               "PATH=/bin|USER=imagec",
		"sessions|Session1/cmd/env":                "1",
		"sessions|Session1/cmd/dir":                "/",
		"sessions|Session1/tty":                    "true",
		"sessions":                                 "Session1",
		"guestinfo.Key":                            "",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded ExecutorConfigPointers
	Decode(MapSource(encoded), &decoded)

	// cannot assert equality at a high level because of the different structure types, but we can test the
	// common structure fragments
	assert.Equal(t, reference.Sessions["Session1"], *decoded.Sessions["Session1"], "Encoded and decoded sessions do not match")

}

func TestInsideOutside(t *testing.T) {
	type Inside struct {
		ID   string `vic:"0.1" scope:"read-write" key:"id"`
		Name string `vic:"0.1" scope:"read-write" key:"name"`
	}

	type Outside struct {
		Inside Inside `vic:"0.1" scope:"read-only" key:"inside"`
		ID     string `vic:"0.1" scope:"read-write" key:"id"`
		Name   string `vic:"0.1" scope:"read-write" key:"name"`
	}
	outside := Outside{
		Inside: Inside{
			ID:   "inside",
			Name: "Inside",
		},
		ID:   "outside",
		Name: "Outside",
	}

	encoded := map[string]string{}
	Encode(MapSink(encoded), outside)

	expected := map[string]string{
		"guestinfo.inside.id":   "inside",
		"guestinfo.inside.name": "Inside",
		"guestinfo.id":          "outside",
		"guestinfo.name":        "Outside",
	}
	assert.Equal(t, expected, encoded, "Encoded and expected does not match")

	var decoded Outside
	Decode(MapSource(encoded), &decoded)

	assert.Equal(t, outside, decoded, "Encoded and decoded does not match")

}
