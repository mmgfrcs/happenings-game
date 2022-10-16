package entities

type Attribute struct {
  	ShortName    string `json:"short"`
	Name     string `json:"name"`
	Description     string `json:"desc"`
	Range    int16    `json:"range"`
	Critical bool   `json:"critical"`
}