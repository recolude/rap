package event

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/metadata"
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

type Encoder struct{}

func NewEncoder() Encoder {
	return Encoder{}
}

func (p Encoder) Accepts(stream format.CaptureCollection) bool {
	return stream.Signature() == "recolude.event"
}

func (p Encoder) Signature() string {
	return "recolude.event"
}

func (p Encoder) Version() uint {
	return 0
}

func (p Encoder) Encode(streams []format.CaptureCollection) ([]byte, [][]byte, error) {
	eventNamesSet := make(map[string]int)
	eventKeysSet := make(map[string]int)

	streamDataBuffers := make([]bytes.Buffer, len(streams))
	for bufferIndex, stream := range streams {
		for _, c := range stream.Captures() {
			eventCapture, ok := c.(event.Capture)
			if !ok {
				return nil, nil, errors.New("capture is not of type event")
			}

			if _, ok := eventNamesSet[eventCapture.Name()]; !ok {
				eventNamesSet[eventCapture.Name()] = len(eventNamesSet)
			}

			nameIndex := make([]byte, binary.MaxVarintLen64)
			read := binary.PutUvarint(nameIndex, uint64(eventNamesSet[eventCapture.Name()]))
			streamDataBuffers[bufferIndex].Write(nameIndex[:read])

			allKeyIndxes := make([]uint, len(eventCapture.Metadata().Mapping()))
			allValueDataBuffer := bytes.Buffer{}
			keyCount := 0
			for key, val := range eventCapture.Metadata().Mapping() {
				if _, ok := eventKeysSet[key]; !ok {
					eventKeysSet[key] = len(eventKeysSet)
				}
				allKeyIndxes[keyCount] = uint(eventKeysSet[key])
				allValueDataBuffer.WriteByte(val.Code())
				allValueDataBuffer.Write(val.Data())
				keyCount++
			}

			streamDataBuffers[bufferIndex].Write(rapbinary.UvarintArrayToBytes(allKeyIndxes))
			streamDataBuffers[bufferIndex].Write(allValueDataBuffer.Bytes())
		}
	}

	streamData := make([][]byte, len(streams))
	for i, buffer := range streamDataBuffers {
		streamData[i] = buffer.Bytes()
	}

	header := bytes.Buffer{}

	allNames := make([]string, len(eventNamesSet))
	for key, index := range eventNamesSet {
		allNames[index] = key
	}
	header.Write(rapbinary.StringArrayToBytes(allNames))

	allKeys := make([]string, len(eventKeysSet))
	for key, index := range eventKeysSet {
		allKeys[index] = key
	}
	header.Write(rapbinary.StringArrayToBytes(allKeys))

	return header.Bytes(), streamData, nil
}

func readHeader(header []byte) (names []string, metadataKeys []string, err error) {
	headerReader := bytes.NewBuffer(header)
	names, _, err = rapbinary.ReadStringArray(headerReader)
	metadataKeys, _, err = rapbinary.ReadStringArray(headerReader)
	return
}

func (p Encoder) Decode(name string, header []byte, streamData []byte, times []float64) (format.CaptureCollection, error) {
	buf := bufio.NewReader(bytes.NewReader(streamData))

	eventNames, metadataKeys, err := readHeader(header)
	if err != nil {
		return nil, err
	}

	captures := make([]event.Capture, len(times))
	for i := 0; i < len(times); i++ {
		eventNameIndex, err := binary.ReadUvarint(buf)
		if err != nil {
			return nil, err
		}

		metadataIndeces, _, err := rapbinary.ReadUvarIntArray(buf)
		if err != nil {
			return nil, err
		}

		block := make(map[string]metadata.Property)
		for metadataIndex := 0; metadataIndex < len(metadataIndeces); metadataIndex++ {
			prop, err := metadata.ReadProperty(buf)
			if err != nil {
				return nil, err
			}
			block[metadataKeys[metadataIndeces[metadataIndex]]] = prop
		}
		captures[i] = event.NewCapture(times[i], eventNames[int(eventNameIndex)], metadata.NewBlock(block))
	}

	return event.NewCollection(name, captures), nil
}
