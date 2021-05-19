package format

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

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

// VECTOR2 ====================================================================

type Vector2Property struct {
	x float64
	y float64
}

func NewVector2Property(x, y float64) Vector2Property {
	return Vector2Property{
		x: x,
		y: y,
	}
}

func (v2p Vector2Property) Code() byte {
	return 6
}

func (v2p Vector2Property) String() string {
	return fmt.Sprintf("%f, %f", v2p.x, v2p.y)
}

func (v2p Vector2Property) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, v2p.x)
	binary.Write(buf, binary.LittleEndian, v2p.y)
	return buf.Bytes()
}

// VECTOR3 ====================================================================

type Vector3Property struct {
	x float64
	y float64
	z float64
}

func NewVector3Property(x, y, z float64) Vector3Property {
	return Vector3Property{
		x: x,
		y: y,
		z: z,
	}
}

func (v3p Vector3Property) Code() byte {
	return 7
}

func (v3p Vector3Property) String() string {
	return fmt.Sprintf("%f, %f, %f", v3p.x, v3p.y, v3p.z)
}

func (v3p Vector3Property) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, v3p.x)
	binary.Write(buf, binary.LittleEndian, v3p.y)
	binary.Write(buf, binary.LittleEndian, v3p.z)
	return buf.Bytes()
}

// QUATERNION =================================================================

type QuaternionProperty struct {
	x float64
	y float64
	z float64
	w float64
}

func NewQuaternionProperty(x, y, z, w float64) QuaternionProperty {
	return QuaternionProperty{
		x: x,
		y: y,
		z: z,
		w: w,
	}
}

func (qp QuaternionProperty) Code() byte {
	return 8
}

func (qp QuaternionProperty) String() string {
	return fmt.Sprintf("%f, %f, %f, %f", qp.x, qp.y, qp.z, qp.w)
}

func (qp QuaternionProperty) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, qp.x)
	binary.Write(buf, binary.LittleEndian, qp.y)
	binary.Write(buf, binary.LittleEndian, qp.z)
	binary.Write(buf, binary.LittleEndian, qp.w)
	return buf.Bytes()
}

// MATRIX3x3 ==================================================================

type Matrix3x3Property struct {
	r1c1, r1f2, r1c3 float64
	r2c1, r2f2, r2c3 float64
	r3c1, r3f2, r3c3 float64
}

func NewMatrix3x3Property(r1c1, r1f2, r1c3, r2c1, r2f2, r2c3, r3c1, r3f2, r3c3 float64) Matrix3x3Property {
	return Matrix3x3Property{
		r1c1: r1c1,
		r1f2: r1f2,
		r1c3: r1c3,
		r2c1: r2c1,
		r2f2: r2f2,
		r2c3: r2c3,
		r3c1: r3c1,
		r3f2: r3f2,
		r3c3: r3c3,
	}
}

func (m3p Matrix3x3Property) Code() byte {
	return 9
}

func (m3p Matrix3x3Property) String() string {
	return fmt.Sprintf(
		"[[%f, %f, %f], [%f, %f, %f], [%f, %f, %f]]",
		m3p.r1c1,
		m3p.r1f2,
		m3p.r1c3,
		m3p.r2c1,
		m3p.r2f2,
		m3p.r2c3,
		m3p.r3c1,
		m3p.r3f2,
		m3p.r3c3,
	)
}

func (m3p Matrix3x3Property) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, m3p.r1c1)
	binary.Write(buf, binary.LittleEndian, m3p.r1f2)
	binary.Write(buf, binary.LittleEndian, m3p.r1c3)
	binary.Write(buf, binary.LittleEndian, m3p.r2c1)
	binary.Write(buf, binary.LittleEndian, m3p.r2f2)
	binary.Write(buf, binary.LittleEndian, m3p.r2c3)
	binary.Write(buf, binary.LittleEndian, m3p.r3c1)
	binary.Write(buf, binary.LittleEndian, m3p.r3f2)
	binary.Write(buf, binary.LittleEndian, m3p.r3c3)
	return buf.Bytes()
}

// MATRIX4x4 ==================================================================

type Matrix4x4Property struct {
	r1c1, r1f2, r1c3, r1c4 float64
	r2c1, r2f2, r2c3, r2c4 float64
	r3c1, r3f2, r3c3, r3c4 float64
	r4c1, r4f2, r4c3, r4c4 float64
}

func NewMatrix4x4Property(r1c1, r1f2, r1c3, r1c4, r2c1, r2f2, r2c3, r2c4, r3c1, r3f2, r3c3, r3c4, r4c1, r4f2, r4c3, r4c4 float64) Matrix4x4Property {
	return Matrix4x4Property{
		r1c1: r1c1,
		r1f2: r1f2,
		r1c3: r1c3,
		r1c4: r1c4,
		r2c1: r2c1,
		r2f2: r2f2,
		r2c3: r2c3,
		r2c4: r2c4,
		r3c1: r3c1,
		r3f2: r3f2,
		r3c3: r3c3,
		r3c4: r3c4,
		r4c1: r4c1,
		r4f2: r4f2,
		r4c3: r4c3,
		r4c4: r4c4,
	}
}

func (m4p Matrix4x4Property) Code() byte {
	return 10
}

func (m4p Matrix4x4Property) String() string {
	return fmt.Sprintf(
		"[[%f, %f, %f, %f], [%f, %f, %f, %f], [%f, %f, %f, %f], [%f, %f, %f, %f]]",
		m4p.r1c1,
		m4p.r1f2,
		m4p.r1c3,
		m4p.r1c4,
		m4p.r2c1,
		m4p.r2f2,
		m4p.r2c3,
		m4p.r2c4,
		m4p.r3c1,
		m4p.r3f2,
		m4p.r3c3,
		m4p.r3c4,
		m4p.r4c1,
		m4p.r4f2,
		m4p.r4c3,
		m4p.r4c4,
	)
}

func (m4p Matrix4x4Property) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, m4p.r1c1)
	binary.Write(buf, binary.LittleEndian, m4p.r1f2)
	binary.Write(buf, binary.LittleEndian, m4p.r1c3)
	binary.Write(buf, binary.LittleEndian, m4p.r1c4)
	binary.Write(buf, binary.LittleEndian, m4p.r2c1)
	binary.Write(buf, binary.LittleEndian, m4p.r2f2)
	binary.Write(buf, binary.LittleEndian, m4p.r2c3)
	binary.Write(buf, binary.LittleEndian, m4p.r2c4)
	binary.Write(buf, binary.LittleEndian, m4p.r3c1)
	binary.Write(buf, binary.LittleEndian, m4p.r3f2)
	binary.Write(buf, binary.LittleEndian, m4p.r3c3)
	binary.Write(buf, binary.LittleEndian, m4p.r3c4)
	binary.Write(buf, binary.LittleEndian, m4p.r4c1)
	binary.Write(buf, binary.LittleEndian, m4p.r4f2)
	binary.Write(buf, binary.LittleEndian, m4p.r4c3)
	binary.Write(buf, binary.LittleEndian, m4p.r4c4)
	return buf.Bytes()
}

// METADATA ===================================================================

type MetadataProperty struct {
	block Metadata
}

func NewMetadataProperty(block Metadata) MetadataProperty {
	return MetadataProperty{block}
}

func (m4p MetadataProperty) Code() byte {
	return 11
}

func (mp MetadataProperty) String() string {
	out := strings.Builder{}
	out.WriteString("{")
	i := 0
	for key, prop := range mp.block.Mapping() {
		fmt.Fprintf(&out, " \"%s\": %s", key, prop.String())
		if i < len(mp.block.Mapping())-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("}")
	return out.String()
}

func (mp MetadataProperty) Data() []byte {
	buf := new(bytes.Buffer)

	i := 0
	mappingWithIndex := make([]string, len(mp.block.Mapping()))
	for key := range mp.block.Mapping() {
		mappingWithIndex[i] = key
		i++
	}

	buf.Write(rapbin.StringArrayToBytes(mappingWithIndex))

	for _, key := range mappingWithIndex {
		buf.WriteByte(mp.block.Mapping()[key].Code())
		buf.Write(mp.block.Mapping()[key].Data())
	}

	return buf.Bytes()
}

func readFloat64s(r io.Reader, count int) ([]float64, error) {
	outValues := make([]float64, count)
	for i := 0; i < count; i++ {
		var val float64
		err := binary.Read(r, binary.LittleEndian, &val)
		if err != nil {
			return nil, err
		}
		outValues[i] = val
	}
	return outValues, nil
}

func readNestedMetadatablock(b *bytes.Reader) (Metadata, error) {
	metadataKeys, _, err := rapbin.ReadStringArray(b)
	if err != nil {
		return EmptyMetadataBlock(), err
	}

	metadata := make(map[string]Property)

	for _, key := range metadataKeys {
		metadata[key], err = ReadProperty(b)
		if err != nil {
			return EmptyMetadataBlock(), err
		}
	}

	return NewMetadataBlock(metadata), nil
}

func ReadProperty(b *bytes.Reader) (Property, error) {
	propertyType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	switch propertyType {
	case 0:
		str, _, err := rapbin.ReadString(b)
		if err != nil {
			return nil, err
		}
		return NewStringProperty(str), nil

	case 1:
		var int32Val int32
		err := binary.Read(b, binary.LittleEndian, &int32Val)
		if err != nil {
			return nil, err
		}
		return NewIntProperty(int32Val), nil

	case 2:
		var float32Val float32
		err := binary.Read(b, binary.LittleEndian, &float32Val)
		if err != nil {
			return nil, err
		}
		return NewFloat32Property(float32Val), nil

	case 3:
		return NewBoolProperty(true), nil
	case 4:
		return NewBoolProperty(false), nil
	case 5:
		byteVal, err := b.ReadByte()
		if err != nil {
			return nil, err
		}
		return NewByteProperty(byteVal), nil

	case 6:
		vals, err := readFloat64s(b, 2)
		if err != nil {
			return nil, err
		}
		return NewVector2Property(vals[0], vals[1]), nil

	case 7:
		vals, err := readFloat64s(b, 3)
		if err != nil {
			return nil, err
		}
		return NewVector3Property(vals[0], vals[1], vals[2]), nil

	case 8:
		vals, err := readFloat64s(b, 4)
		if err != nil {
			return nil, err
		}
		return NewQuaternionProperty(vals[0], vals[1], vals[2], vals[3]), nil

	case 9:
		vals, err := readFloat64s(b, 9)
		if err != nil {
			return nil, err
		}
		return NewMatrix3x3Property(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5], vals[6], vals[7], vals[8]), nil
	case 10:
		vals, err := readFloat64s(b, 16)
		if err != nil {
			return nil, err
		}
		return NewMatrix4x4Property(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5], vals[6], vals[7], vals[8], vals[9], vals[10], vals[11], vals[12], vals[13], vals[14], vals[15]), nil
	case 11:

		metadataBlock, err := readNestedMetadatablock(b)
		if err != nil {
			return nil, err
		}
		return NewMetadataProperty(metadataBlock), nil
	}

	return nil, fmt.Errorf("unrecognized property type code: %d", int(propertyType))
}
