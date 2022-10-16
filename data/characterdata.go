package data

type Character struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"desc"`
	Attributes  []Attribute `json:"attributes"`
}
