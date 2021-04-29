package io

import (
	"fmt"
	"io"

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

func (r Reader) readEncoders() ([]encoding.Encoder, int, error) {
	totalBytesRead := 0

	encoderSignatures, read, err := binary.ReadStringArray(r.in)
	totalBytesRead += read
	if err != nil {
		return nil, totalBytesRead, err
	}

	encoderVersions, read, err := binary.ReadUintArray(r.in)
	totalBytesRead += read
	if err != nil {
		return nil, totalBytesRead, err
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

	// Read name
	recordingName, bytesRead, err := binary.ReadString(r.in)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	// Read encoders
	encodersToUse, bytesRead, err := r.readEncoders()
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	// Read metadata
	meatadataBlock, bytesRead, err := binary.ReadStringArray(r.in)
	totalBytesRead += bytesRead
	if err != nil {
		return nil, totalBytesRead, err
	}

	metadata := make(map[string]string)
	for i := 0; i < len(meatadataBlock); i += 2 {
		metadata[meatadataBlock[i]] = meatadataBlock[i+1]
	}

	allStreams := make([]data.CaptureStream, 0)
	// Read streams
	for _, encoder := range encodersToUse {
		// read header
		header, bytesRead, err := binary.ReadBytesArray(r.in)
		totalBytesRead += bytesRead
		if err != nil {
			return nil, totalBytesRead, err
		}

		// read num streams
		len, bytesRead, err := binary.ReadUvarint(r.in)
		totalBytesRead += bytesRead
		if err != nil {
			return nil, totalBytesRead, err
		}

		// read streams
		allBodies := make([][]byte, len)
		for i := 0; i < int(len); i++ {
			captureBody, bytesRead, err := binary.ReadBytesArray(r.in)
			totalBytesRead += bytesRead
			if err != nil {
				return nil, totalBytesRead, err
			}
			allBodies[i] = captureBody
		}

		// rebuild streams
		rebuiltStreams, err := encoder.Decode(header, allBodies)
		if err != nil {
			return nil, totalBytesRead, err
		}
		allStreams = append(allStreams, rebuiltStreams...)
	}

	return data.NewRecording(recordingName, allStreams, nil, metadata, nil), totalBytesRead, nil
}
