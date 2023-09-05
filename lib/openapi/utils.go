package openapi

import (
	"encoding/json"
	"strings"
)

// JSON returns the OpenAPI doc in JSON format.
func (doc *Document) JSON() string {
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (p Path) MarshalJSON() ([]byte, error) {
	m := map[string]Operation{}

	for method, op := range p {
		m[strings.ToLower(method.String())] = op
	}

	return json.Marshal(m)
}
