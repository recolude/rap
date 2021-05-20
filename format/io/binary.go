package io

import (
	"bytes"
	"io"

	"github.com/recolude/rap/format/metadata"
)

type Binary struct {
	name  string
	data  []byte
	block metadata.Block
}

func NewBinary(name string, data []byte, block metadata.Block) Binary {
	return Binary{
		name:  name,
		data:  data,
		block: block,
	}
}

func (br Binary) Name() string {
	return br.name
}

func (br Binary) Size() uint64 {
	return uint64(len(br.data))
}

func (br Binary) Metadata() metadata.Block {
	return br.block
}

func (br Binary) Data() io.Reader {
	return bytes.NewReader(br.data)
}
