package binary

import (
	"bytes"
	"encoding/binary"
	"io"
)

func UvarintArrayToBytes(strs []uint) []byte {
	buf := new(bytes.Buffer)

	varByte := make([]byte, binary.MaxVarintLen64)
	read := binary.PutUvarint(varByte, uint64(len(strs)))
	buf.Write(varByte[:read])

	for _, s := range strs {
		read := binary.PutUvarint(varByte, uint64(s))
		buf.Write(varByte[:read])
	}
	return buf.Bytes()
}

func ReadUvarIntArray(r io.Reader) ([]uint, int, error) {
	len, bytesRead, err := ReadUvarint(r)
	if err != nil {
		return nil, bytesRead, err
	}

	out := make([]uint, len)
	for i := range out {
		str, read, err := ReadUvarint(r)
		bytesRead += read
		if err != nil {
			return nil, bytesRead, err
		}
		out[i] = uint(str)
	}

	return out, bytesRead, nil
}
