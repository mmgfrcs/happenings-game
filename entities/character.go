package entities

type Character struct {
  ID         string `json:"id"`
	Name       string `json:"name"`
	Description       string `json:"desc"`
	Attributes map[string]int16 `json:"attributes"`
}