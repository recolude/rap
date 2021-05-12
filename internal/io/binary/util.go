package binary

import (
	"encoding/binary"
	"errors"
	"io"
)

var overflow = errors.New("binary: varint overflows a 64-bit integer")

// ReadUvarint reads an encoded unsigned integer from r and returns it as a
// uint64.
func ReadUvarint(r io.Reader) (uint64, int, error) {
	b := []byte{0}

	var x uint64
	var s uint
	var i int
	var totalRead int
	for i = 0; i < binary.MaxVarintLen64; i++ {
		read, err := r.Read(b)
		totalRead += read
		if err != nil {
			return x, totalRead, err
		}
		if b[0] < 0x80 {
			if i == 9 && b[0] > 1 {
				return x, totalRead, overflow
			}
			return x | uint64(b[0])<<s, i + 1, nil
		}
		x |= uint64(b[0]&0x7f) << s
		s += 7
	}
	return x, totalRead, overflow
}

func UnsignedFloatBSTToBytes(value, start, duration float64, out []byte) {
	curValue := start + (duration / 2.0)
	increment := duration / 4.0

	for byteIndex := 0; byteIndex < len(out); byteIndex++ {

		// Clear whatever byte might be there
		out[byteIndex] = 0

		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			if value < curValue {
				curValue -= increment
			} else {
				out[byteIndex] = out[byteIndex] | (1 << bitIndex)
				curValue += increment
			}
			increment /= 2.0
		}
	}
}

func BytesToUnisngedFloatBST(start, duration float64, in []byte) float64 {
	curValue := start + (duration / 2.0)
	increment := duration / 4.0

	for byteIndex := 0; byteIndex < len(in); byteIndex++ {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			if (in[byteIndex]>>byte(bitIndex))&1 == 1 {
				curValue += increment
			} else {
				curValue -= increment
			}

			increment /= 2.0
		}
	}

	return curValue
}
