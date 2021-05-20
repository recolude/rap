package metadata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	rapbin "github.com/recolude/rap/internal/io/binary"
)

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

func readNestedMetadatablock(b *bytes.Reader) (Block, error) {
	metadataKeys, _, err := rapbin.ReadStringArray(b)
	if err != nil {
		return EmptyBlock(), err
	}

	metadata := make(map[string]Property)

	for _, key := range metadataKeys {
		metadata[key], err = ReadProperty(b)
		if err != nil {
			return EmptyBlock(), err
		}
	}

	return NewBlock(metadata), nil
}

func readPropData(b *bytes.Reader, propertyType byte) (Property, error) {
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
		return NewIntProperty(int(int32Val)), nil

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

	case 12:
		var unixTime int64
		err := binary.Read(b, binary.LittleEndian, &unixTime)
		if err != nil {
			return nil, err
		}
		return NewTimeProperty(time.Unix(0, unixTime)), nil

	case 13:
		fallthrough
	case 14:
		len, _, err := rapbin.ReadUvarint(b)
		if err != nil {
			return nil, err
		}
		adjustedType := propertyType - 13
		props := make([]Property, len)
		for i := 0; i < int(len); i++ {
			props[i], err = readPropData(b, adjustedType)
			if err != nil {
				return nil, err
			}
		}
		return newArrayProperty(adjustedType, props), nil
	}
	return nil, fmt.Errorf("unrecognized property type code: %d", int(propertyType))
}

func ReadProperty(b *bytes.Reader) (Property, error) {
	propertyType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	return readPropData(b, propertyType)
}
