package openapi

import (
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

	Summary     string                `json:"summary,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
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
	Description string   `json:"description"`
	Headers     Headers  `json:"headers,omitempty"`
	Content     *Content `json:"content,omitempty"`
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
	Schemas         map[string]*jschema.Schema `json:"schemas"`
	SecuritySchemes map[string]SecurityScheme  `json:"securitySchemes,omitempty"`
}

type SecurityScheme struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Scheme           string            `json:"scheme,omitempty"`
	BearerFormat     string            `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlowsObject `json:"flows,omitempty"`
	OpenIdConnectUrl string            `json:"openIdConnectUrl,omitempty"`
}

type OAuthFlowsObject struct {
	Implicit          *OAuthFlowObject `json:"implicit,omitempty"`
	Password          *OAuthFlowObject `json:"password,omitempty"`
	ClientCredentials *OAuthFlowObject `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlowObject `json:"authorizationCode,omitempty"`
}

type OAuthFlowObject struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

// Error is an error object that contains information about a failed request.
// Reference: https://github.com/microsoft/api-guidelines/blob/vNext/Guidelines.md#error--object
type Error struct {
	// Code is a machine-readable error code.
	Code string `json:"code,omitempty"`
	// Message is a human-readable error message.
	Message string `json:"message,omitempty"`
	// Target is a human-readable description of the target of the error.
	Target string `json:"target,omitempty"`
	// Details is an array of structured error details objects.
	Details []Error `json:"details,omitempty"`
	// InnerError is a generic error object that is used by the service developer for debugging.
	InnerError any `json:"innererror,omitempty"`
}

const (
	CodeNotFound      = "not_found"
	CodeInvalidParam  = "invalid_param"
	CodeInternalError = "internal_error"
)

type ResponseFormat interface {
	format()
}

type ResponseFormatErr struct {
	Error any `json:"error"`
}

func (ResponseFormatErr) format() {}

type ResponseFormatMeta struct {
	Data any `json:"data"`
	Meta any `json:"meta"`
}

func (ResponseFormatMeta) format() {}

type ResponseFormatData struct {
	Data any `json:"data"`
}

func (ResponseFormatData) format() {}
