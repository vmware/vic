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

package toolbox

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Property type enum as defined in open-vm-tools/lib/include/vix.h
const (
	_ = iota // ANY type not supported
	vixPropertyTypeInt32
	vixPropertyTypeString
	vixPropertyTypeBool
	_ // HANDLE type not supported
	vixPropertyTypeInt64
	vixPropertyTypeBlob
)

// Property ID enum as defined in open-vm-tools/lib/include/vixOpenSource.h
const (
	VixPropertyGuestToolsAPIOptions     = 4501
	VixPropertyGuestOsFamily            = 4502
	VixPropertyGuestOsVersion           = 4503
	VixPropertyGuestToolsProductNam     = 4511
	VixPropertyGuestToolsVersion        = 4500
	VixPropertyGuestName                = 4505
	VixPropertyGuestOsVersionShort      = 4520
	VixPropertyGuestStartProgramEnabled = 4540
)

type VixProperty struct {
	header struct {
		ID     int32
		Kind   int32
		Length int32
	}

	data struct {
		Int32  int32
		String string
		Bool   uint8
		Int64  int64
		Blob   []byte
	}
}

var int32Size int32

func init() {
	var i int32
	int32Size = int32(binary.Size(&i))
}

type VixPropertyList []*VixProperty

func NewInt32Property(ID int32, val int32) *VixProperty {
	p := new(VixProperty)
	p.header.ID = ID
	p.header.Kind = vixPropertyTypeInt32
	p.header.Length = int32Size
	p.data.Int32 = val
	return p
}

func NewStringProperty(ID int32, val string) *VixProperty {
	p := new(VixProperty)
	p.header.ID = ID
	p.header.Kind = vixPropertyTypeString
	p.header.Length = int32(len(val) + 1)
	p.data.String = val
	return p
}

func NewBoolProperty(ID int32, val bool) *VixProperty {
	p := new(VixProperty)
	p.header.ID = ID
	p.header.Kind = vixPropertyTypeBool
	p.header.Length = 1
	if val {
		p.data.Bool = 1
	}
	return p
}

func NewInt64Property(ID int32, val int64) *VixProperty {
	p := new(VixProperty)
	p.header.ID = ID
	p.header.Kind = vixPropertyTypeInt64
	p.header.Length = int32Size * 2
	p.data.Int64 = val
	return p
}

func NewBlobProperty(ID int32, val []byte) *VixProperty {
	p := new(VixProperty)
	p.header.ID = ID
	p.header.Kind = vixPropertyTypeBlob
	p.header.Length = int32(len(val))
	p.data.Blob = val
	return p
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (p *VixProperty) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, &p.header)
	if err != nil {
		return nil, err
	}

	switch p.header.Kind {
	case vixPropertyTypeBool:
		err = binary.Write(buf, binary.LittleEndian, p.data.Bool)
	case vixPropertyTypeInt32:
		err = binary.Write(buf, binary.LittleEndian, p.data.Int32)
	case vixPropertyTypeInt64:
		err = binary.Write(buf, binary.LittleEndian, p.data.Int64)
	case vixPropertyTypeString:
		_, err = buf.WriteString(p.data.String)
		if err == nil {
			err = buf.WriteByte(0)
		}
	case vixPropertyTypeBlob:
		_, err = buf.Write(p.data.Blob)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (p *VixProperty) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &p.header)
	if err != nil {
		return err
	}

	switch p.header.Kind {
	case vixPropertyTypeBool:
		return binary.Read(buf, binary.LittleEndian, &p.data.Bool)
	case vixPropertyTypeInt32:
		return binary.Read(buf, binary.LittleEndian, &p.data.Int32)
	case vixPropertyTypeInt64:
		return binary.Read(buf, binary.LittleEndian, &p.data.Int64)
	case vixPropertyTypeString:
		s := make([]byte, p.header.Length-1)
		if _, err := buf.Read(s); err != nil {
			return err
		}
		if _, err := buf.ReadByte(); err != nil { // discard \0
			return err
		}
		p.data.String = string(s)
	case vixPropertyTypeBlob:
		p.data.Blob = make([]byte, p.header.Length)
		if _, err := buf.Read(p.data.Blob); err != nil {
			return err
		}
	default:
		return errors.New("VIX_E_UNRECOGNIZED_PROPERTY")
	}

	return nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (l *VixPropertyList) UnmarshalBinary(data []byte) error {
	headerSize := int32Size * 3

	for {
		p := new(VixProperty)

		err := p.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		*l = append(*l, p)

		offset := headerSize + p.header.Length
		data = data[offset:]

		if len(data) == 0 {
			return nil
		}
	}
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (l *VixPropertyList) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	for _, p := range *l {
		b, err := p.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if _, err = buf.Write(b); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
