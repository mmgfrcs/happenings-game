package entities

type Character struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"desc"`
	Attributes  []Attribute `json:"attributes"`
}

// GetAttributeStrict returns a single attribute by id. Returns an empty attribute if the attribute is not found.
func (c Character) GetAttribute(id string) Attribute {
	for i := range c.Attributes {
		if c.Attributes[i].ShortName == id {
			return c.Attributes[i]
		}
	}
	return Attribute{}
}

// GetAttributeStrict returns a single attribute by id. Same as GetAttribute, but panics if the attribute is not found.
func (c Character) GetAttributeStrict(id string) Attribute {
	attr := c.GetAttribute(id)
	if attr.ShortName == "" {
		panic("Attribute not found")
	}
	return attr
}
