package io

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/recolude/rap/internal/io/binary"
	"github.com/recolude/rap/internal/io/rapv1"
	"github.com/recolude/rap/pkg/data"
	"github.com/recolude/rap/pkg/encoding"
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

func recursiveBuidRecordings(recordingData []byte, encoders []encoding.Encoder, headers [][]byte) (data.Recording, error) {
	in := bytes.NewReader(recordingData)

	// Read Recording name
	recordingName, _, err := binary.ReadString(in)
	if err != nil {
		return nil, err
	}

	// Read Recording metadata
	meatadataBlock, _, err := binary.ReadStringArray(in)
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]string)
	for i := 0; i < len(meatadataBlock); i += 2 {
		metadata[meatadataBlock[i]] = meatadataBlock[i+1]
	}

	// read num streams
	numStreams, _, err := binary.ReadUvarint(in)
	if err != nil {
		return nil, err
	}

	// read streams
	allStreams := make([]data.CaptureStream, numStreams)
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

	allChildRecordings := make([]data.Recording, numRecordings)
	for i := 0; i < int(numRecordings); i++ {
		childRecData, _, err := binary.ReadBytesArray(in)
		if err != nil {
			return nil, err
		}
		childRec, err := recursiveBuidRecordings(childRecData, encoders, headers)
		allChildRecordings[i] = childRec
	}

	return data.NewRecording(recordingName, allStreams, allChildRecordings, metadata, nil), nil
}

func (r Reader) Read() (data.Recording, int, error) {
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

	compressedRecData, bytesRead, err := binary.ReadBytesArray(r.in)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	deflateReader := flate.NewReader(bytes.NewReader(compressedRecData))

	uncompresseRecordingData, err := ioutil.ReadAll(deflateReader)
	if err != nil {
		return nil, totalBytesRead, err
	}

	rec, err := recursiveBuidRecordings(uncompresseRecordingData, encodersToUse, encoderHeaders)

	return rec, totalBytesRead, err
}
