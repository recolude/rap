package io

import "github.com/recolude/rap/format"

type BinaryReference struct {
	name     string
	uri      string
	size     uint64
	metadata format.Metadata
}

func NewBinaryReference(name string, uri string, size uint64, metadata format.Metadata) BinaryReference {
	return BinaryReference{
		name:     name,
		uri:      uri,
		size:     size,
		metadata: metadata,
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

func (br BinaryReference) Metadata() format.Metadata {
	return br.metadata
}
