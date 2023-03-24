package entities

import "encoding/json"

type Action struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"desc"`
	Tags        []string          `json:"tags"`
	Mods        []map[string]interface{} `json:"mod"`
}
