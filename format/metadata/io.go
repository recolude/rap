package metadata

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	rapbin "github.com/recolude/rap/internal/io/binary"
)

func WriteProprty(writer io.Writer, property Property) (int, error) {
	count, err := writer.Write([]byte{property.Code()})
	if err != nil {
		return count, err
	}
	other, err := writer.Write(property.Data())
	return count + other, err
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

func readNestedMetadatablock(b io.Reader) (Block, error) {
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

func readPropData(b io.Reader, propertyType byte) (Property, error) {
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
		byteVal := make([]byte, 1)
		_, err := b.Read(byteVal)
		if err != nil {
			return nil, err
		}
		return NewByteProperty(byteVal[0]), nil

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

	case 16:
		allbytes, _, err := rapbin.ReadBytesArray(b)
		if err != nil {
			return nil, err
		}
		return ArrayPropertyRaw{
			data:             rapbin.BytesArrayToBytes(allbytes),
			originalBaseCode: 3,
			divison:          1,
		}, err

	case 18:
		allbytes, _, err := rapbin.ReadBytesArray(b)
		if err != nil {
			return nil, err
		}
		return NewBinaryArrayProperty(allbytes), nil

	case 13:
		fallthrough
	case 15:
		fallthrough
	case 19:
		fallthrough
	case 20:
		fallthrough
	case 24:
		fallthrough
	case 25:
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

func ReadProperty(b io.Reader) (Property, error) {
	propByte := make([]byte, 1)
	_, err := b.Read(propByte)
	if err != nil {
		return nil, err
	}
	return readPropData(b, propByte[0])
}
