package format

import (
	"bytes"
	"encoding/binary"
	"fmt"

	rapbin "github.com/recolude/rap/internal/io/binary"
)

type Property interface {
	Code() byte
	String() string
	Data() []byte
}

// STRING =====================================================================

type StringProperty struct {
	str string
}

func NewStringProperty(str string) StringProperty {
	return StringProperty{
		str: str,
	}
}

func (sp StringProperty) Code() byte {
	return 0
}

func (sp StringProperty) String() string {
	return sp.str
}

func (sp StringProperty) Data() []byte {
	return rapbin.StringToBytes(sp.str)
}

// INT32 ======================================================================

type Int32Property struct {
	i int32
}

func NewIntProperty(i int32) Int32Property {
	return Int32Property{
		i: i,
	}
}

func (ip Int32Property) Code() byte {
	return 1
}

func (ip Int32Property) String() string {
	return fmt.Sprintf("%d", ip.i)
}

func (ip Int32Property) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ip.i)
	return buf.Bytes()
}

// FLOAT32 ====================================================================

type Float32Property struct {
	i float32
}

func NewFloat32Property(i float32) Float32Property {
	return Float32Property{
		i: i,
	}
}

func (fp Float32Property) Code() byte {
	return 2
}

func (fp Float32Property) String() string {
	return fmt.Sprintf("%f", fp.i)
}

func (fp Float32Property) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, fp.i)
	return buf.Bytes()
}

// BOOL =======================================================================

type BoolProperty struct {
	b bool
}

func NewBoolProperty(b bool) BoolProperty {
	return BoolProperty{
		b: b,
	}
}

func (bp BoolProperty) Code() byte {
	if bp.b {
		return 3
	}
	return 4
}

func (bp BoolProperty) String() string {
	if bp.b {
		return "true"
	}
	return "false"
}

func (bp BoolProperty) Data() []byte {
	return nil
}

func (bp BoolProperty) Value() bool {
	return bp.b
}

// BYTE =======================================================================

type ByteProperty struct {
	b byte
}

func NewByteProperty(b byte) ByteProperty {
	return ByteProperty{
		b: b,
	}
}

func (bp ByteProperty) Code() byte {
	return 5
}

func (bp ByteProperty) String() string {
	return fmt.Sprintf("%d", int(bp.b))
}

func (bp ByteProperty) Data() []byte {
	return []byte{bp.b}
}

func (bp ByteProperty) Value() byte {
	return bp.b
}
