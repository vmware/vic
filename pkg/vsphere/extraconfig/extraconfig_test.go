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
	"net/url"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/vim25/types"
)

// [BEGIN] SLIMMED DOWNED and MODIFIED VERSION of github.com/vmware/vic/metadata
type Common struct {
	ExecutionEnvironment string

	ID string `vic:"0.1" scope:"hidden" key:"id"`

	Name string `vic:"0.1" scope:"hidden" key:"name"`

	Notes string `vic:"0.1" scope:"hidden" key:"notes"`
}

type ContainerVM struct {
	Common `vic:"0.1" scope:"read-only" key:"common"`

	Version string `vic:"0.1" scope:"hidden" key:"version"`

	Aliases map[string]string

	Interaction url.URL

	AgentKey []byte
}

type ExecutorConfig struct {
	Common `vic:"0.1" scope:"read-only" key:"common"`

	Sessions map[string]SessionConfig `vic:"0.1" scope:"hidden" key:"sessions"`

	Key []byte `json:"byte"`
}

type Cmd struct {
	Path string `vic:"0.1" scope:"hidden" key:"path"`

	Args []string `vic:"0.1" scope:"hidden" key:"args"`

	Env []string `vic:"0.1" scope:"hidden" key:"env"`

	Dir string `vic:"0.1" scope:"hidden" key:"dir"`

	Cmd *exec.Cmd `vic:"0.1" scope:"hidden" key:"cmd"`
}

type SessionConfig struct {
	Common `vic:"0.1" scope:"hidden" key:"common" json:"page"`

	Cmd Cmd `vic:"0.1" scope:"hidden" key:"cmd"`

	Tty bool `vic:"0.1" scope:"hidden" key:"tty"`
}

// [END] SLIMMED VERSION of github.com/vmware/vic/metadata

func TestBasic(t *testing.T) {
	type Type struct {
		Int    int     `vic:"0.1" scope:"read-write" key:"int"`
		Bool   bool    `vic:"0.1" scope:"read-write" key:"bool"`
		Float  float64 `vic:"0.1" scope:"read-write" key:"float"`
		String string  `vic:"0.1" scope:"read-write" key:"string"`
	}

	Struct := Type{
		42,
		true,
		3.14,
		"Grrr",
	}

	encoded := Encode(Struct)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo.int", Value: "42"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo.bool", Value: "true"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo.float", Value: "3.14E+00"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo.string", Value: "Grrr"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, Struct, decoded, "Encoded and decoded does not match")
}

func TestBasicMap(t *testing.T) {
	t.Skip("https://github.com/stretchr/testify/issues/288")

	type Type struct {
		IntMap map[string]int `vic:"0.1" scope:"read-only" key:"intmap"`
	}

	IntMap := Type{
		map[string]int{
			"1st": 12345,
			"2nd": 67890,
		},
	}

	// Encode
	encoded := Encode(IntMap)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/intmap|1st", Value: "12345"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/intmap|2nd", Value: "67890"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/intmap", Value: "1st|2nd"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	// Decode to new variable
	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, IntMap, decoded, "Encoded and decoded does not match")

	// Decode to already existing variable
	IntMapOptimusPrime := Type{
		map[string]int{
			"first":  1,
			"second": 2,
			"1st":    0,
		},
	}
	Decode(encoded, &IntMapOptimusPrime)

	// We expect a merge and over-write
	expectedOptimusPrime := Type{
		map[string]int{
			"1st":    12345,
			"2nd":    67890,
			"first":  1,
			"second": 2,
		},
	}
	assert.Equal(t, IntMapOptimusPrime, expectedOptimusPrime, "Decoded and expected does not match")

}

func TestBasicSlice(t *testing.T) {
	type Type struct {
		IntSlice []int `vic:"0.1" scope:"read-only" key:"intslice"`
	}

	IntSlice := Type{
		[]int{1, 2, 3, 4, 5},
	}

	encoded := Encode(IntSlice)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/intslice~", Value: "1|2|3|4|5"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/intslice", Value: "4"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, IntSlice, decoded, "Encoded and decoded does not match")
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

	encoded := Encode(Embedded)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/common~id", Value: "0xDEADBEEF"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/common~name", Value: "Embedded"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/common~notes", Value: ""},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, Embedded, decoded, "Encoded and decoded does not match")
}

func TestStruct(t *testing.T) {
	type Type struct {
		Common Common `vic:"0.1" scope:"read-only" key:"common"`
	}

	Struct := Type{
		Common: Common{
			ID:   "0xDEADBEEF",
			Name: "Struct",
		},
	}

	encoded := Encode(Struct)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/common~id", Value: "0xDEADBEEF"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/common~name", Value: "Struct"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/common~notes", Value: ""},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, Struct, decoded, "Encoded and decoded does not match")
}

func TestTime(t *testing.T) {
	type Type struct {
		Time time.Time `vic:"0.1" scope:"read-only" key:"time"`
	}

	Time := Type{
		Time: time.Date(2009, 11, 10, 23, 00, 00, 0, time.UTC),
	}

	encoded := Encode(Time)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/time", Value: "2009-11-10 23:00:00 +0000 UTC"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, Time, decoded, "Encoded and decoded does not match")
}

func TestNet(t *testing.T) {
	type Type struct {
		Net net.IPNet `vic:"0.1" scope:"read-only" key:"net"`
	}

	_, n, _ := net.ParseCIDR("127.0.0.1/8")
	Net := Type{
		Net: *n,
	}

	encoded := Encode(Net)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/net", Value: "127.0.0.0/8"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, Net, decoded, "Encoded and decoded does not match")
}

func TestNetPointer(t *testing.T) {
	type Type struct {
		Net *net.IPNet `vic:"0.1" scope:"read-only" key:"net"`
	}

	_, n, _ := net.ParseCIDR("127.0.0.1/8")
	Net := Type{
		Net: n,
	}

	encoded := Encode(Net)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/net", Value: "127.0.0.0/8"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

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

	encoded := Encode(Time)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/time", Value: "2009-11-10 23:00:00 +0000 UTC"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, Time, decoded, "Encoded and decoded does not match")
}

func TestStructMap(t *testing.T) {
	t.Skip("https://github.com/stretchr/testify/issues/288")

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

	encoded := Encode(StructMap)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key1~id", Value: "0xDEADBEEF"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key1~name", Value: "beef"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key1~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key2~id", Value: "0x8BADF00D"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key2~name", Value: "food"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key2~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key3~id", Value: "0xDEADF00D"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key3~name", Value: "dead"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|Key3~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map", Value: "Key1|Key2|Key3"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, StructMap, decoded, "Encoded and decoded does not match")
}

func TestIntStructMap(t *testing.T) {
	t.Skip("https://github.com/stretchr/testify/issues/288")

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

	encoded := Encode(StructMap)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|1~id", Value: "0xDEADBEEF"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|1~name", Value: "beef"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|1~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|2~id", Value: "0x8BADF00D"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|2~name", Value: "food"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|2~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|3~id", Value: "0xDEADF00D"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|3~name", Value: "dead"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map|3~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/map", Value: "1|2|3"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

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

	encoded := Encode(StructSlice)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/slice|0~id", Value: "0xDEADFEED"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/slice|0~name", Value: "feed"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/slice|0~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/slice|1~id", Value: "0xFACEFEED"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/slice|1~name", Value: "face"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/slice|1~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/slice", Value: "1"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, StructSlice, decoded, "Encoded and decoded does not match")
}

func TestMultipleScope(t *testing.T) {
	MultipleScope := struct {
		MultipleScope string `vic:"0.1" scope:"read-only,hidden,non-persistent" key:"multiscope"`
	}{
		"MultipleScope",
	}

	encoded := Encode(MultipleScope)
	var expected []types.BaseOptionValue
	assert.Equal(t, encoded, expected, "Not equal")
}

func TestUnknownScope(t *testing.T) {
	UnknownScope := struct {
		UnknownScope int `vic:"0.1" scope:"unknownscope" key:"unknownscope"`
	}{
		42,
	}

	encoded := Encode(UnknownScope)
	var expected []types.BaseOptionValue
	assert.Equal(t, encoded, expected, "Not equal")
}

func TestUnknownProperty(t *testing.T) {
	UnknownProperty := struct {
		UnknownProperty int `vic:"0.1" scope:"hidden" key:"unknownproperty,unknownproperty"`
	}{
		42,
	}

	encoded := Encode(UnknownProperty)
	var expected []types.BaseOptionValue
	assert.Equal(t, encoded, expected, "Not equal")
}

func TestOmitNested(t *testing.T) {
	OmitNested := struct {
		Time        time.Time `vic:"0.1" scope:"volatile" key:"time,omitnested"`
		CurrentTime time.Time `vic:"0.1" scope:"volatile" key:"time"`
	}{
		Time:        time.Date(2009, 11, 10, 23, 00, 00, 0, time.UTC),
		CurrentTime: time.Date(2009, 11, 10, 23, 00, 00, 0, time.UTC),
	}

	encoded := Encode(OmitNested)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "time", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "time", Value: "2009-11-10 23:00:00 +0000 UTC"},
	}
	assert.Equal(t, encoded, expected, "Encoded and decoded does not match")

}

func TestPointer(t *testing.T) {
	type Type struct {
		Pointer           *ContainerVM `vic:"0.1" scope:"hidden" key:"pointer"`
		PointerOmitnested *ContainerVM `vic:"0.1" scope:"non-persistent" key:"pointeromitnested,omitnested"`
	}

	Pointer := Type{
		Pointer: &ContainerVM{Version: "0.1"},
	}

	encoded := Encode(Pointer)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/pointer/common~id", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/pointer/common~name", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/pointer/common~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "pointer~version", Value: "0.1"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "pointer", Value: "pointer~version"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "pointeromitnested", Value: ""},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, Pointer, decoded, "Encoded and decoded does not match")
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
						Args: []string{"-standalone", "/bin/imagec"},
						Env:  []string{"PATH=/bin", "USER=imagec"},
						Dir:  "/",
					},
				},
			},
		},
	}
	encoded := Encode(ExecutorConfig)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/executorconfig/common~id", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/executorconfig/common~name", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/executorconfig/common~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~common~id", Value: "SessionID"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~common~name", Value: "SessionName"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~common~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~path", Value: "/vmware"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~args~", Value: "-standalone|/bin/imagec"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~args", Value: "1"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~env~", Value: "PATH=/bin|USER=imagec"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~env", Value: "1"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~dir", Value: "/"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~tty", Value: "true"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions", Value: "Session1"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

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
						Args: []string{"-standalone", "/bin/imagec"},
						Env:  []string{"PATH=/bin", "USER=imagec"},
						Dir:  "/",
					},
				},
			},
		},
	}

	encoded := Encode(ExecutorConfig)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/executorconfig/common~id", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/executorconfig/common~name", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/executorconfig/common~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~common~id", Value: "SessionID"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~common~name", Value: "SessionName"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~common~notes", Value: ""},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~path", Value: "/vmware"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~args~", Value: "-standalone|/bin/imagec"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~args", Value: "1"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~env~", Value: "PATH=/bin|USER=imagec"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~env", Value: "1"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~cmd~dir", Value: "/"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions|Session1~tty", Value: "true"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig~sessions", Value: "Session1"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "executorconfig", Value: "executorconfig~sessions"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Type
	Decode(encoded, &decoded)

	assert.Equal(t, ExecutorConfig, decoded, "Encoded and decoded does not match")
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

	encoded := Encode(outside)
	expected := []types.BaseOptionValue{
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/inside.id", Value: "inside"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo/inside.name", Value: "Inside"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo.id", Value: "outside"},
		&types.OptionValue{DynamicData: types.DynamicData{}, Key: "guestinfo.name", Value: "Outside"},
	}
	assert.Equal(t, encoded, expected, "Encoded and expected does not match")

	var decoded Outside
	Decode(encoded, &decoded)

	assert.Equal(t, outside, decoded, "Encoded and decoded does not match")

}
