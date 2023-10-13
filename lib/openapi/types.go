package openapi

import (
	"github.com/NaturalSelectionLabs/jschema"
)

// Version represents the version of an OpenAPI document.
type Version string

// MarshalJSON implements the [json.Marshaler] interface.
func (v Version) MarshalJSON() ([]byte, error) {
	// We only support openapi 3.1.0 for now
	return []byte(`"3.1.0"`), nil
}

// Document represents an OpenAPI document.
type Document struct {
	OpenAPI      Version         `json:"openapi"`
	Info         Info            `json:"info"`
	Servers      []Server        `json:"servers,omitempty"`
	Paths        map[string]Path `json:"paths"`
	Components   Components      `json:"components"`
	Tags         []Tag           `json:"tags,omitempty"`
	ExternalDocs *ExternalDocs   `json:"externalDocs,omitempty"`
	Extension    Extension       `json:"x-extension,omitempty"`
}

// Info represents the info section of an OpenAPI document.
type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

// Path represents a path in an OpenAPI document.
type Path map[Method]Operation

// Server represents a server in an OpenAPI document.
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// Operation represents an operation in an OpenAPI document.
type Operation struct {
	Parameters  []Parameter             `json:"parameters,omitempty"`
	RequestBody *RequestBody            `json:"requestBody,omitempty"`
	Responses   map[StatusCode]Response `json:"responses"`

	Summary     string                `json:"summary,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Extension   Extension             `json:"x-extension,omitempty"`
}

// Parameter represents a parameter in an OpenAPI document.
type Parameter struct {
	Name        string             `json:"name"`
	In          ParamIn            `json:"in"`
	Schema      *jschema.Schema    `json:"schema"`
	Description string             `json:"description,omitempty"`
	Required    bool               `json:"required,omitempty"`
	Examples    map[string]Example `json:"examples,omitempty"`
}

// Example represents an example in an OpenAPI document.
type Example struct {
	Value         jschema.JVal `json:"value"`
	Summary       string       `json:"summary,omitempty"`
	Description   string       `json:"description,omitempty"`
	ExternalValue string       `json:"externalValue,omitempty"`
}

// RequestBody represents a request body in an OpenAPI document.
type RequestBody struct {
	Content  *Content `json:"content,omitempty"`
	Required bool     `json:"required,omitempty"`
}

// Response represents a response in an OpenAPI document.
type Response struct {
	Description string   `json:"description"`
	Headers     Headers  `json:"headers,omitempty"`
	Content     *Content `json:"content,omitempty"`
}

// Headers represents a headers in an OpenAPI document.
type Headers map[string]Header

// Header represents a header in an OpenAPI document.
type Header struct {
	Description string          `json:"description,omitempty"`
	Schema      *jschema.Schema `json:"schema"`
}

const (
	// ContentTypeJSON represents the JSON http content type.
	ContentTypeJSON = "application/json"
	// ContentTypeBin represents the binary http content type.
	ContentTypeBin = "application/octet-stream"
)

// Content represents a content in an OpenAPI document.
type Content map[string]*Schema

// Schema represents a schema in an OpenAPI document.
type Schema struct {
	Schema *jschema.Schema `json:"schema"`
}

// Components represents the components section of an OpenAPI document.
type Components struct {
	Schemas         map[string]*jschema.Schema `json:"schemas"`
	SecuritySchemes map[string]SecurityScheme  `json:"securitySchemes,omitempty"`
}

// SecurityScheme represents a security scheme in an OpenAPI document.
type SecurityScheme struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Scheme           string            `json:"scheme,omitempty"`
	BearerFormat     string            `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlowsObject `json:"flows,omitempty"`
	OpenIDConnectURL string            `json:"openIdConnectUrl,omitempty"`
}

// OAuthFlowsObject represents a OAuthFlowsObject in an OpenAPI document.
type OAuthFlowsObject struct {
	Implicit          *OAuthFlowObject `json:"implicit,omitempty"`
	Password          *OAuthFlowObject `json:"password,omitempty"`
	ClientCredentials *OAuthFlowObject `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlowObject `json:"authorizationCode,omitempty"`
}

// OAuthFlowObject represents a OAuthFlowObject in an OpenAPI document.
type OAuthFlowObject struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

// Extension represents an extension in an OpenAPI document.
type Extension any

// Tag represents a tag in an OpenAPI document.
type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// ExternalDocs represents an externalDocs in an OpenAPI document.
type ExternalDocs struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}
