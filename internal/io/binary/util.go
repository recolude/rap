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
