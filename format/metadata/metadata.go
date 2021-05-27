package metadata

type Block struct {
	mapping map[string]Property
}

func NewBlock(mapping map[string]Property) Block {
	return Block{mapping: mapping}
}

func EmptyBlock() Block {
	return NewBlock(make(map[string]Property))
}

func (m Block) Mapping() map[string]Property {
	return m.mapping
}
