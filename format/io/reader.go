package io

import (
	"compress/flate"
	"fmt"
	"io"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/encoding"
	"github.com/recolude/rap/format/metadata"
	"github.com/recolude/rap/internal/io/binary"
	"github.com/recolude/rap/internal/io/rapv1"
)

// https://dave.cheney.net/2019/01/27/eliminate-error-handling-by-eliminating-errors
type errReader struct {
	io.Reader
	err error
	n   int
}

func (e *errReader) Read(p []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}

	var n int
	n, e.err = io.ReadFull(e.Reader, p)
	e.n += n
	return n, e.err
}

func (e *errReader) TotalRead() int {
	return e.n
}

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

func (r Reader) readEncoders() ([]encoding.Encoder, int, error) {
	totalBytesRead := 0

	encoderSignatures, read, err := binary.ReadStringArray(r.in)
	totalBytesRead += read
	if err != nil {
		return nil, totalBytesRead, err
	}

	encoderVersions := make([]uint64, len(encoderSignatures))
	for i := range encoderSignatures {
		val, read, err := binary.ReadUvarint(r.in)
		if err != nil {
			return nil, totalBytesRead, err
		}
		totalBytesRead += read
		encoderVersions[i] = val
	}

	encoders := make([]encoding.Encoder, len(encoderSignatures))
	for i, desiredEncoderSignature := range encoderSignatures {
		found := false
		for _, registeredEncoder := range r.encoders {
			if registeredEncoder.Signature() == desiredEncoderSignature {
				if registeredEncoder.Version() >= uint(encoderVersions[i]) {
					encoders[i] = registeredEncoder
					found = true
				} else {
					return nil, totalBytesRead, fmt.Errorf(
						"registered encoder (%s) version is behind what is found in recording: %d < %d",
						desiredEncoderSignature,
						registeredEncoder.Version(),
						encoderVersions[i],
					)
				}
			}
		}
		if found == false {
			return nil, totalBytesRead, fmt.Errorf("no registered encoder has signature %s", desiredEncoderSignature)
		}
	}

	return encoders, totalBytesRead, nil
}

func readRecordingMetadataBlock(in io.Reader, metadataKeys []string) (metadata.Block, error) {
	propMapping := make(map[string]metadata.Property)

	keyIndecies, _, err := binary.ReadUvarIntArray(in)
	if err != nil {
		return metadata.EmptyBlock(), err
	}

	for _, key := range keyIndecies {
		propMapping[metadataKeys[key]], err = metadata.ReadProperty(in)
		if err != nil {
			return metadata.EmptyBlock(), err
		}
	}

	return metadata.NewBlock(propMapping), nil
}

func recursiveBuidRecordings(inStream io.Reader, metadataKeys []string, encoders []encoding.Encoder, headers [][]byte) (format.Recording, int, error) {
	// in := bytes.NewReader(recordingData)
	er := &errReader{Reader: inStream}

	// Read Recording id
	recordingID, _, err := binary.ReadString(er)

	// Read Recording name
	recordingName, _, err := binary.ReadString(er)

	// Read Recording metadata
	recordingMetadataBlock, err := readRecordingMetadataBlock(er, metadataKeys)
	if err != nil {
		return nil, er.TotalRead(), err
	}

	// read num streams
	numStreams, _, err := binary.ReadUvarint(er)

	// read streams
	allStreams := make([]format.CaptureCollection, numStreams)
	for i := 0; i < int(numStreams); i++ {
		encoderIndex, _, _ := binary.ReadUvarint(er)

		streamName, _, _ := binary.ReadString(er)

		times, _ := decodeTime(er)
		captureBody, _, _ := binary.ReadBytesArray(er)
		stream, _ := encoders[encoderIndex].Decode(streamName, headers[encoderIndex], captureBody, times)
		allStreams[i] = stream
	}

	// read binary references
	numBinaryReferences, _, err := binary.ReadUvarint(er)
	binReferences := make([]format.BinaryReference, numBinaryReferences)

	for i := 0; i < int(numBinaryReferences); i++ {
		name, _, _ := binary.ReadString(er)
		uri, _, _ := binary.ReadString(er)
		refSize, _, _ := binary.ReadUvarint(er)
		// Read Recording metadata
		block, _ := readRecordingMetadataBlock(er, metadataKeys)

		binReferences[i] = NewBinaryReference(name, uri, refSize, block)
	}

	// read binaries
	numBinaries, _, err := binary.ReadUvarint(er)
	binaries := make([]format.Binary, numBinaries)

	for i := 0; i < int(numBinaries); i++ {
		name, _, _ := binary.ReadString(er)
		refSize, _, _ := binary.ReadUvarint(er)

		// Read Recording metadata
		block, _ := readRecordingMetadataBlock(er, metadataKeys)

		allData := make([]byte, refSize)
		io.ReadFull(er, allData)

		binaries[i] = NewBinary(name, allData, block)
	}

	// read num recordings
	numRecordings, _, err := binary.ReadUvarint(er)

	allChildRecordings := make([]format.Recording, numRecordings)
	for i := 0; i < int(numRecordings); i++ {
		childRec, _, err := recursiveBuidRecordings(er, metadataKeys, encoders, headers)
		if err != nil {
			return nil, er.TotalRead(), err
		}
		allChildRecordings[i] = childRec
	}

	return format.NewRecording(recordingID, recordingName, allStreams, allChildRecordings, recordingMetadataBlock, binaries, binReferences), er.TotalRead(), er.err
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
	encodersToUse, bytesRead, err := r.readEncoders()
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	compressedFlag := []byte{0}
	bytesRead, err = r.in.Read(compressedFlag)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	var readcloser io.Reader = r.in
	if compressedFlag[0] == 1 {
		readcloser = flate.NewReader(r.in)
	}

	encoderHeaders := make([][]byte, len(encodersToUse))
	for i := range encoderHeaders {
		header, read, err := binary.ReadBytesArray(readcloser)
		totalBytesRead += read
		if err != nil {
			return nil, totalBytesRead, err
		}
		encoderHeaders[i] = header
	}

	// Read off metadata keys
	metdataKeys, bytesRead, err := binary.ReadStringArray(readcloser)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	// Read off recordings
	rec, bytesRead, err := recursiveBuidRecordings(readcloser, metdataKeys, encodersToUse, encoderHeaders)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	return rec, totalBytesRead, err
}
