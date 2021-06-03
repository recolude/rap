package binary

import (
	"bytes"
	"encoding/binary"
	"io"
)

func StringToBytes(str string) []byte {
	strBytes := []byte(str)

	varByte := make([]byte, binary.MaxVarintLen64)
	read := binary.PutUvarint(varByte, uint64(len(strBytes)))

	return append(varByte[:read], strBytes...)
}

func ReadString(r io.Reader) (string, int, error) {
	len, bytesRead, err := ReadUvarint(r)
	if err != nil {
		return "", bytesRead, err
	}

	strBuffer := make([]byte, len)
	moreBytes, err := io.ReadFull(r, strBuffer)
	if err != nil {
		return "", bytesRead + moreBytes, err
	}

	return string(strBuffer), bytesRead + moreBytes, nil
}

func StringArrayToBytes(strs []string) []byte {
	varByte := make([]byte, binary.MaxVarintLen64)
	read := binary.PutUvarint(varByte, uint64(len(strs)))

	buf := new(bytes.Buffer)
	buf.Write(varByte[:read])

	for _, s := range strs {
		buf.Write(StringToBytes(s))
	}
	return buf.Bytes()
}

func ReadStringArray(r io.Reader) ([]string, int, error) {
	len, bytesRead, err := ReadUvarint(r)
	if err != nil {
		return nil, bytesRead, err
	}

	out := make([]string, len)
	for i := range out {
		str, read, err := ReadString(r)
		bytesRead += read
		if err != nil {
			return nil, bytesRead, err
		}
		out[i] = str
	}

	return out, bytesRead, nil
}
