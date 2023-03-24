package data

type Character struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"desc"`
	Attributes  []Attribute `json:"attributes"`
}

func (c Character) GetAttribute(id string) Attribute {
  for i := range c.Attributes {
    if c.Attributes[i].ShortName == id {
        return c.Attributes[i]
    }
  }
  return Attribute{}
}

func (c Character) GetAttributeStrict(id string) Attribute {
  attr := c.GetAttribute(id)
  if attr.ShortName == "" {
    panic("Attribute not found")
  }
  return attr
}