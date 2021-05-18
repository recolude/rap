package io

import (
	"bytes"
	"compress/flate"
	encbin "encoding/binary"
	"fmt"
	"io"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/encoding"
	"github.com/recolude/rap/internal/io/binary"
	"github.com/recolude/rap/internal/io/rapv1"
)

type Reader struct {
	encoders []encoding.Encoder
	in       io.Reader
}

func NewReader(encoders []encoding.Encoder, r io.Reader) Reader {
	return Reader{
		encoders: encoders,
		in:       r,
	}
}

func (r Reader) readEncoders() ([]encoding.Encoder, [][]byte, int, error) {
	totalBytesRead := 0

	encoderSignatures, read, err := binary.ReadStringArray(r.in)
	totalBytesRead += read
	if err != nil {
		return nil, nil, totalBytesRead, err
	}

	encoderVersions, read, err := binary.ReadUintArray(r.in)
	totalBytesRead += read
	if err != nil {
		return nil, nil, totalBytesRead, err
	}

	encoders := make([]encoding.Encoder, len(encoderSignatures))
	for i, desiredEncoderSignature := range encoderSignatures {
		found := false
		for _, registeredEncoder := range r.encoders {
			if registeredEncoder.Signature() == desiredEncoderSignature {
				if registeredEncoder.Version() >= encoderVersions[i] {
					encoders[i] = registeredEncoder
					found = true
				} else {
					return nil, nil, totalBytesRead, fmt.Errorf(
						"registered encoder (%s) version is behind what is found in recording: %d < %d",
						desiredEncoderSignature,
						registeredEncoder.Version(),
						encoderVersions[i],
					)
				}
			}
		}
		if found == false {
			return nil, nil, totalBytesRead, fmt.Errorf("no registered encoder has signature %s", desiredEncoderSignature)
		}
	}

	encoderHeaders := make([][]byte, len(encoders))
	for i := range encoderHeaders {
		header, read, err := binary.ReadBytesArray(r.in)
		totalBytesRead += read
		if err != nil {
			return nil, nil, totalBytesRead, err
		}
		encoderHeaders[i] = header
	}

	return encoders, encoderHeaders, totalBytesRead, nil
}

func readFloat64s(r io.Reader, count int) ([]float64, error) {
	outValues := make([]float64, count)
	for i := 0; i < count; i++ {
		var val float64
		err := encbin.Read(r, encbin.LittleEndian, &val)
		if err != nil {
			return nil, err
		}
		outValues[i] = val
	}
	return outValues, nil
}

func readNestedMetadatablock(b *bytes.Reader) (format.Metadata, error) {
	metadataKeys, _, err := binary.ReadStringArray(b)
	if err != nil {
		return format.EmptyMetadataBlock(), err
	}

	metadata := make(map[string]format.Property)

	// valuesBlock, _, err := binary.ReadBytesArray(in)
	// if err != nil {
	// 	return format.Metadata{}, err
	// }
	// valuesBlockBuffer := bytes.NewReader(valuesBlock)

	for _, key := range metadataKeys {
		metadata[key], err = readProperty(b)
		if err != nil {
			return format.EmptyMetadataBlock(), err
		}
	}

	return format.NewMetadataBlock(metadata), nil
}

func readProperty(b *bytes.Reader) (format.Property, error) {
	propertyType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	switch propertyType {
	case 0:
		str, _, err := binary.ReadString(b)
		if err != nil {
			return nil, err
		}
		return format.NewStringProperty(str), nil

	case 1:
		var int32Val int32
		err := encbin.Read(b, encbin.LittleEndian, &int32Val)
		if err != nil {
			return nil, err
		}
		return format.NewIntProperty(int32Val), nil

	case 2:
		var float32Val float32
		err := encbin.Read(b, encbin.LittleEndian, &float32Val)
		if err != nil {
			return nil, err
		}
		return format.NewFloat32Property(float32Val), nil

	case 3:
		return format.NewBoolProperty(true), nil
	case 4:
		return format.NewBoolProperty(false), nil
	case 5:
		byteVal, err := b.ReadByte()
		if err != nil {
			return nil, err
		}
		return format.NewByteProperty(byteVal), nil

	case 6:
		vals, err := readFloat64s(b, 2)
		if err != nil {
			return nil, err
		}
		return format.NewVector2Property(vals[0], vals[1]), nil

	case 7:
		vals, err := readFloat64s(b, 3)
		if err != nil {
			return nil, err
		}
		return format.NewVector3Property(vals[0], vals[1], vals[2]), nil

	case 8:
		vals, err := readFloat64s(b, 4)
		if err != nil {
			return nil, err
		}
		return format.NewQuaternionProperty(vals[0], vals[1], vals[2], vals[3]), nil

	case 9:
		vals, err := readFloat64s(b, 9)
		if err != nil {
			return nil, err
		}
		return format.NewMatrix3x3Property(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5], vals[6], vals[7], vals[8]), nil
	case 10:
		vals, err := readFloat64s(b, 16)
		if err != nil {
			return nil, err
		}
		return format.NewMatrix4x4Property(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5], vals[6], vals[7], vals[8], vals[9], vals[10], vals[11], vals[12], vals[13], vals[14], vals[15]), nil
	case 11:

		metadataBlock, err := readNestedMetadatablock(b)
		if err != nil {
			return nil, err
		}
		return format.NewMetadataProperty(metadataBlock), nil
	}

	return nil, fmt.Errorf("unrecognized property type code: %d", int(propertyType))
}

func readRecordingMetadataBlock(in *bytes.Reader, metadataKeys []string) (format.Metadata, error) {
	metadata := make(map[string]format.Property)

	keyIndecies, _, err := binary.ReadUintArray(in)
	if err != nil {
		return format.Metadata{}, err
	}

	valuesBlock, _, err := binary.ReadBytesArray(in)
	if err != nil {
		return format.Metadata{}, err
	}
	valuesBlockBuffer := bytes.NewReader(valuesBlock)

	for _, key := range keyIndecies {
		metadata[metadataKeys[key]], err = readProperty(valuesBlockBuffer)
		if err != nil {
			return format.EmptyMetadataBlock(), err
		}
	}

	return format.NewMetadataBlock(metadata), nil
}

func recursiveBuidRecordings(recordingData []byte, metadataKeys []string, encoders []encoding.Encoder, headers [][]byte) (format.Recording, error) {
	in := bytes.NewReader(recordingData)

	// Read Recording name
	recordingName, _, err := binary.ReadString(in)
	if err != nil {
		return nil, err
	}

	// Read Recording metadata
	metadata, err := readRecordingMetadataBlock(in, metadataKeys)
	if err != nil {
		return nil, err
	}

	// read num streams
	numStreams, _, err := binary.ReadUvarint(in)
	if err != nil {
		return nil, err
	}

	// read streams
	allStreams := make([]format.CaptureCollection, numStreams)
	for i := 0; i < int(numStreams); i++ {

		encoderIndex, _, err := binary.ReadUvarint(in)
		if err != nil {
			return nil, err
		}

		captureBody, _, err := binary.ReadBytesArray(in)
		if err != nil {
			return nil, err
		}

		stream, err := encoders[encoderIndex].Decode(headers[encoderIndex], captureBody)
		if err != nil {
			return nil, err
		}

		allStreams[i] = stream
	}

	// read num recordings
	numRecordings, _, err := binary.ReadUvarint(in)
	if err != nil {
		return nil, err
	}

	allChildRecordings := make([]format.Recording, numRecordings)
	for i := 0; i < int(numRecordings); i++ {
		childRecData, _, err := binary.ReadBytesArray(in)
		if err != nil {
			return nil, err
		}
		childRec, err := recursiveBuidRecordings(childRecData, metadataKeys, encoders, headers)
		if err != nil {
			return nil, err
		}
		allChildRecordings[i] = childRec
	}

	return format.NewRecording("", recordingName, allStreams, allChildRecordings, metadata, nil), nil
}

func (r Reader) Read() (format.Recording, int, error) {
	if r.in == nil {
		panic("Attempting to load recording from nil reader")
	}

	totalBytesRead := 0

	// read version
	version, bytesRead, err := GetRecoringVersion(r.in)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, bytesRead, err
	}

	if version == 1 {
		rec, read, err := rapv1.ReadRecording(r.in)
		return rec, read + totalBytesRead, err
	}

	if version != 2 {
		return nil, totalBytesRead, fmt.Errorf("Unrecognized file version: %d", version)
	}

	// Read encoders
	encodersToUse, encoderHeaders, bytesRead, err := r.readEncoders()
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	deflateReader := flate.NewReader(r.in)

	// Read off metadata keys
	metdataKeys, bytesRead, err := binary.ReadStringArray(deflateReader)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	// Read off recordings
	uncompresseRecordingData, bytesRead, err := binary.ReadBytesArray(deflateReader)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	rec, err := recursiveBuidRecordings(uncompresseRecordingData, metdataKeys, encodersToUse, encoderHeaders)

	return rec, totalBytesRead, err
}
