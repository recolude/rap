package metadata

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/EliCDavis/vector"
	rapbin "github.com/recolude/rap/internal/io/binary"
)

const HEX_PREFIX = "0x"

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

// func (sp StringProperty) UnmarshalProperty(data interface{}) error {
// 	err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, data)), &sp)
// 	sp.str = sp.str
// 	return err
// }

func (sp *StringProperty) UnmarshalJSON(b []byte) error {
	var data string
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	sp.str = data
	return nil
}

func (sp StringProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(sp.str)
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

func UnmarshalNewInt32Property(b []byte) (Int32Property, error) {
	var p Int32Property
	err := json.Unmarshal(b, &p)
	return p, err
}

func (ip Int32Property) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`%v`, data)), &ip)
	return err
}

func (ip *Int32Property) UnmarshalJSON(b []byte) error {
	var data int32
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	ip.i = data
	return nil
}

func (ip Int32Property) MarshalJSON() ([]byte, error) {
	return json.Marshal(ip.i)
}

// FLOAT32 ====================================================================

type Float32Property struct {
	f float32
}

func NewFloat32Property(f float32) Float32Property {
	return Float32Property{
		f: f,
	}
}

func (fp Float32Property) Code() byte {
	return 2
}

func (fp Float32Property) String() string {
	return fmt.Sprintf("%f", fp.f)
}

func (fp Float32Property) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, fp.f)
	return buf.Bytes()
}

func UnmarshalNewFloat32Property(b []byte) (Float32Property, error) {
	var p Float32Property
	err := json.Unmarshal(b, &p)
	return p, err
}

func (fp Float32Property) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`%v`, data)), &fp)
	return err
}

func (fp *Float32Property) UnmarshalJSON(b []byte) error {
	var data float32
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	fp.f = data
	return nil
}

func (fp Float32Property) MarshalJSON() ([]byte, error) {
	return json.Marshal(fp.f)
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

func UnmarshalNewBoolProperty(b []byte) (BoolProperty, error) {
	var p BoolProperty
	err := json.Unmarshal(b, &p)
	return p, err
}

func (bp BoolProperty) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`%v`, data)), &bp)
	return err
}

func (bp *BoolProperty) UnmarshalJSON(b []byte) error {
	var data bool
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	bp.b = data
	return nil
}

func (bp BoolProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(bp.b)
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

func UnmarshalNewByteProperty(b []byte) (ByteProperty, error) {
	var p ByteProperty
	err := json.Unmarshal(b, &p)
	return p, err
}

func (bp ByteProperty) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, data)), &bp)
	return err
}

func (bp *ByteProperty) UnmarshalJSON(b []byte) error {
	var data string
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	b, err := hex.DecodeString(strings.TrimPrefix(data, HEX_PREFIX))
	if err != nil {
		return err
	}
	bp.b = b[0]
	return nil
}

func (bp ByteProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(HEX_PREFIX + hex.EncodeToString(bp.Data()))
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

func UnmarshalNewVector2Property(b []byte) (Vector2Property, error) {
	var p Vector2Property
	err := json.Unmarshal(b, &p)
	return p, err
}

func (v2p Vector2Property) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, data)), &v2p)
	return err
}

func (v2p *Vector2Property) UnmarshalJSON(b []byte) error {
	var data struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	v2p.x = data.X
	v2p.y = data.Y
	return nil
}

func (v2p Vector2Property) MarshalJSON() ([]byte, error) {
	data := struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}{
		v2p.x,
		v2p.y,
	}
	return json.Marshal(data)
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

func UnmarshalNewVector3Property(b []byte) (Vector3Property, error) {
	var p Vector3Property
	err := json.Unmarshal(b, &p)
	return p, err
}

func (v3p Vector3Property) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`%v`, data)), &v3p)
	return err
}

func (v3p *Vector3Property) UnmarshalJSON(b []byte) error {
	var data struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	v3p.x = data.X
	v3p.y = data.Y
	v3p.z = data.Z
	return nil
}

func (v3p Vector3Property) MarshalJSON() ([]byte, error) {
	data := struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}{
		v3p.x,
		v3p.y,
		v3p.z,
	}
	return json.Marshal(data)
}

// METADATA ===================================================================

type MetadataProperty struct {
	block Block
}

func NewMetadataProperty(block Block) MetadataProperty {
	return MetadataProperty{block}
}

func (mp MetadataProperty) Code() byte {
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

func UnmarshalNewMetadataProperty(b []byte) (MetadataProperty, error) {
	var p MetadataProperty
	err := json.Unmarshal(b, &p)
	return p, err
}

func (mp MetadataProperty) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, data)), &mp)
	return err
}

func (mp *MetadataProperty) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	prop, err := getProperties(data)
	if err != nil {
		return err
	}
	newMp, ok := prop.(MetadataProperty)
	if !ok {
		return errors.New("unable to unmarshal MetadataProperty")
	}
	mp.block.mapping = newMp.block.mapping
	return nil
}

func (mp MetadataProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(mp.block.Mapping())
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

func UnmarshalNewTimeProperty(b []byte) (TimeProperty, error) {
	var p TimeProperty
	err := json.Unmarshal(b, &p)
	return p, err
}

func (tp TimeProperty) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, data)), &tp)
	return err
}

func (tp *TimeProperty) UnmarshalJSON(b []byte) error {
	var data string
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339Nano, data)
	if err != nil {
		return err
	}
	tp.nanoseconds = t.UnixNano()
	return nil
}

func (tp TimeProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Unix(0, tp.nanoseconds).Format(time.RFC3339Nano))
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

func (ap ArrayProperty) Code() byte {
	return 13 + ap.originalBaseCode
}

func (ap ArrayProperty) String() string {
	return fmt.Sprintf("Type %d Array of %d elements", ap.originalBaseCode, len(ap.props))
}

func (ap ArrayProperty) Data() []byte {
	buf := bytes.Buffer{}

	numBinaries := make([]byte, binary.MaxVarintLen64)
	read := binary.PutUvarint(numBinaries, uint64(len(ap.props)))
	buf.Write(numBinaries[:read])

	for _, prop := range ap.props {
		buf.Write(prop.Data())
	}
	return buf.Bytes()
}

func (ap ArrayProperty) UnmarshalProperty(data interface{}) error {
	err := json.Unmarshal([]byte(fmt.Sprintf(`%v`, data)), &ap)
	return err
}

func (ap *ArrayProperty) UnmarshalJSON(b []byte) error {
	var data []interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	prop, err := getProperties(data)
	if err != nil {
		return err
	}
	newAp, ok := prop.(ArrayProperty)
	if !ok {
		return errors.New("unable to unmarshal MetadataProperty")
	}
	ap.props = newAp.props
	ap.originalBaseCode = newAp.originalBaseCode
	return nil
}

func (ap ArrayProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(ap.props)
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

func (apr ArrayPropertyRaw) Code() byte {
	return 13 + apr.originalBaseCode
}

func (apr ArrayPropertyRaw) String() string {
	return fmt.Sprintf("Type %d Array of %d elements", apr.originalBaseCode, len(apr.data)/apr.divison)
}

func (apr ArrayPropertyRaw) Data() []byte {
	return apr.data
}

func (apr *ArrayPropertyRaw) UnmarshalJSON(b []byte) error {
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	prop, err := getProperties(data)
	if err != nil {
		return err
	}
	newApr, ok := prop.(ArrayPropertyRaw)
	if !ok {
		return errors.New("unable to unmarshal MetadataProperty")
	}
	apr.data = newApr.data
	apr.originalBaseCode = newApr.originalBaseCode
	apr.divison = newApr.divison
	return nil
}

func (apr ArrayPropertyRaw) MarshalJSON() ([]byte, error) {
	switch apr.originalBaseCode {
	case 5:
		br, _, err := rapbin.ReadBytesArray(bytes.NewBuffer(apr.data))
		if err != nil {
			return nil, err
		}
		return json.Marshal(HEX_PREFIX + hex.EncodeToString(br))
	case 3:
		br, _, err := rapbin.ReadBytesArray(bytes.NewBuffer(apr.data))
		if err != nil {
			return nil, err
		}
		bools := make([]bool, 0, len(br)/2)
		for _, b := range br {
			if b&1 == 1 {
				bools = append(bools, true)
				continue
			}
			bools = append(bools, false)
		}
		return json.Marshal(bools)
	}
	return nil, nil
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
