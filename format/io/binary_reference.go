package io

import (
	"github.com/recolude/rap/format/metadata"
)

type BinaryReference struct {
	name  string
	uri   string
	size  uint64
	block metadata.Block
}

func NewBinaryReference(name string, uri string, size uint64, block metadata.Block) BinaryReference {
	return BinaryReference{
		name:  name,
		uri:   uri,
		size:  size,
		block: block,
	}
}

func (br BinaryReference) Name() string {
	return br.name
}

func (br BinaryReference) URI() string {
	return br.uri
}

func (br BinaryReference) Size() uint64 {
	return br.size
}

func (br BinaryReference) Metadata() metadata.Block {
	return br.block
}
