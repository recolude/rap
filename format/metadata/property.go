package metadata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/EliCDavis/vector"
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

func NewIntProperty(i int) Int32Property {
	return Int32Property{
		i: int32(i),
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

// METADATA ===================================================================

type MetadataProperty struct {
	block Block
}

func NewMetadataProperty(block Block) MetadataProperty {
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

// Time =======================================================================

type TimeProperty struct {
	nanoseconds int64
}

func NewTimeProperty(t time.Time) TimeProperty {
	return TimeProperty{
		nanoseconds: t.UnixNano(),
	}
}

func (tp TimeProperty) Code() byte {
	return 12
}

func (tp TimeProperty) String() string {
	return fmt.Sprintf("%d ns", tp.nanoseconds)
}

func (tp TimeProperty) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, tp.nanoseconds)
	return buf.Bytes()
}

// ARRAY =====================================================================

type ArrayProperty struct {
	originalBaseCode byte
	props            []Property
}

func newArrayProperty(originalBaseCode byte, props []Property) ArrayProperty {
	return ArrayProperty{
		originalBaseCode: originalBaseCode,
		props:            props,
	}
}

func (sp ArrayProperty) Code() byte {
	return 13 + sp.originalBaseCode
}

func (sp ArrayProperty) String() string {
	return fmt.Sprintf("Type %d Array of %d elements", sp.originalBaseCode, len(sp.props))
}

func (sp ArrayProperty) Data() []byte {
	buf := bytes.Buffer{}

	numBinaries := make([]byte, 4)
	read := binary.PutUvarint(numBinaries, uint64(len(sp.props)))
	buf.Write(numBinaries[:read])

	for _, prop := range sp.props {
		buf.Write(prop.Data())
	}
	return buf.Bytes()
}

func NewStringArrayProperty(strs []string) ArrayProperty {
	strProps := make([]Property, len(strs))
	for i, str := range strs {
		strProps[i] = NewStringProperty(str)
	}
	return newArrayProperty(0, strProps)
}

func NewIntArrayProperty(strs []int) ArrayProperty {
	strProps := make([]Property, len(strs))
	for i, str := range strs {
		strProps[i] = NewIntProperty(str)
	}
	return newArrayProperty(1, strProps)
}

func NewFloat32ArrayProperty(entries []float32) ArrayProperty {
	strProps := make([]Property, len(entries))
	for i, entry := range entries {
		strProps[i] = NewFloat32Property(entry)
	}
	return newArrayProperty(2, strProps)
}

func NewVector2ArrayProperty(entries []vector.Vector2) ArrayProperty {
	strProps := make([]Property, len(entries))
	for i, entry := range entries {
		strProps[i] = NewVector2Property(entry.X(), entry.Y())
	}
	return newArrayProperty(6, strProps)
}

func NewVector3ArrayProperty(entries []vector.Vector3) ArrayProperty {
	strProps := make([]Property, len(entries))
	for i, entry := range entries {
		strProps[i] = NewVector3Property(entry.X(), entry.Y(), entry.Z())
	}
	return newArrayProperty(7, strProps)
}

func NewMetadataArrayProperty(entries []Block) ArrayProperty {
	strProps := make([]Property, len(entries))
	for i, entry := range entries {
		strProps[i] = NewMetadataProperty(entry)
	}
	return newArrayProperty(11, strProps)
}

func NewTimestampArrayProperty(entries []time.Time) ArrayProperty {
	props := make([]Property, len(entries))
	for i, entry := range entries {
		props[i] = NewTimeProperty(entry)
	}
	return newArrayProperty(12, props)
}

// BINARY ARRAY ===============================================================

type ArrayPropertyRaw struct {
	originalBaseCode byte
	data             []byte
	divison          int
}

func (sp ArrayPropertyRaw) Code() byte {
	return 13 + sp.originalBaseCode
}

func (sp ArrayPropertyRaw) String() string {
	return fmt.Sprintf("Type %d Array of %d elements", sp.originalBaseCode, len(sp.data)/sp.divison)
}

func (sp ArrayPropertyRaw) Data() []byte {
	return sp.data
}

func NewBinaryArrayProperty(binarr []byte) ArrayPropertyRaw {
	return ArrayPropertyRaw{
		data:             rapbin.BytesArrayToBytes(binarr),
		originalBaseCode: 5,
		divison:          1,
	}
}

func NewBoolArrayProperty(boolArr []bool) ArrayPropertyRaw {
	data := make([]byte, len(boolArr))
	for i, b := range boolArr {
		if b {
			data[i] = 1
		} else {
			data[i] = 0
		}
	}

	return ArrayPropertyRaw{
		data:             rapbin.BytesArrayToBytes(data),
		originalBaseCode: 3,
		divison:          1,
	}
}
