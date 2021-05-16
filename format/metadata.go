package format

type Metadata struct {
	mapping map[string]Property
}

func NewMetadataBlock(mapping map[string]Property) Metadata {
	return Metadata{mapping: mapping}
}

func EmptyMetadataBlock() Metadata {
	return Metadata{}
}

func (m Metadata) Mapping() map[string]Property {
	return m.mapping
}
