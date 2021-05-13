package binary

import (
	"bytes"
	"encoding/binary"
	"io"
)

// BytesArrayToBytes creates a new byte array, prefixed with the length of the byte array
func BytesArrayToBytes(b []byte) []byte {
	buf := new(bytes.Buffer)

	varByte := make([]byte, 4)
	read := binary.PutUvarint(varByte, uint64(len(b)))
	buf.Write(varByte[:read])

	buf.Write(b)
	return buf.Bytes()
}

// ReadBytesArray first reads the length of the byte array, then reads in a
// buffer of that length
func ReadBytesArray(r io.Reader) ([]byte, int, error) {
	len, bytesRead, err := ReadUvarint(r)
	if err != nil {
		return nil, bytesRead, err
	}

	out := make([]byte, len)
	read, err := io.ReadFull(r, out)
	if err != nil && err != io.EOF {
		return nil, read + bytesRead, err
	}

	if read != int(len) {
		return nil, read + bytesRead, io.ErrUnexpectedEOF
	}

	return out, read + bytesRead, nil
}
