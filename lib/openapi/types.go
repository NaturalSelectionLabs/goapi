package openapi

import (
	"net/http"

	"github.com/NaturalSelectionLabs/jschema"
)

type OpenAPIVersion string

func (v OpenAPIVersion) MarshalJSON() ([]byte, error) {
	// We only support openapi 3.1.0 for now
	return []byte(`"3.1.0"`), nil
}

// Document represents an OpenAPI document.
type Document struct {
	OpenAPI    OpenAPIVersion  `json:"openapi"`
	Info       Info            `json:"info"`
	Servers    []Server        `json:"servers,omitempty"`
	Paths      map[string]Path `json:"paths"`
	Components Components      `json:"components"`
}

type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

type Path map[Method]Operation

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Operation struct {
	Parameters  []Parameter             `json:"parameters,omitempty"`
	RequestBody *RequestBody            `json:"requestBody,omitempty"`
	Responses   map[StatusCode]Response `json:"responses"`

	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	OperationID string   `json:"operationId,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type Parameter struct {
	Name        string          `json:"name"`
	In          ParamIn         `json:"in"`
	Schema      *jschema.Schema `json:"schema"`
	Description string          `json:"description,omitempty"`
	Required    bool            `json:"required,omitempty"`
}

type RequestBody struct {
	Content  *Content `json:"content,omitempty"`
	Required bool     `json:"required,omitempty"`
}

type Response struct {
	Description string      `json:"description,omitempty"`
	Headers     http.Header `json:"headers,omitempty"`
	Content     *Content    `json:"content,omitempty"`
}

type Headers map[string]Header

type Header struct {
	Description string          `json:"description,omitempty"`
	Schema      *jschema.Schema `json:"schema"`
}

type Content struct {
	JSON *Schema `json:"application/json"`
}

type Schema struct {
	Schema *jschema.Schema `json:"schema"`
}

type Components struct {
	Schemas map[string]*jschema.Schema `json:"schemas"`
}
