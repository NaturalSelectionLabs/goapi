package openapi

// CommonError is an error object that contains information about a failed request.
// Reference: https://github.com/microsoft/api-guidelines/blob/vNext/Guidelines.md#error--object
type CommonError[C any] struct {
	// Code is a machine-readable error code.
	Code C `json:"code"`
	// Message is a human-readable error message.
	Message string `json:"message,omitempty"`
	// Target is a human-readable description of the target of the error.
	Target string `json:"target,omitempty"`
	// Details is an array of structured error details objects.
	Details []CommonError[C] `json:"details,omitempty"`
	// InnerError is a generic error object that is used by the service developer for debugging.
	InnerError any `json:"innererror,omitempty"`
}

var _ error = &CommonError[Code]{}

func (e CommonError[E]) Error() string {
	return e.Message
}

// Error ...
type Error CommonError[Code]

// ResponseFormat for the json response body.
type ResponseFormat interface {
	format()
}

// ResponseFormatErr is the error response format.
type ResponseFormatErr struct {
	Error any `json:"error"`
}

func (ResponseFormatErr) format() {}

// ResponseFormatMeta is the data and meta response format.
type ResponseFormatMeta struct {
	Data any `json:"data"`
	Meta any `json:"meta"`
}

func (ResponseFormatMeta) format() {}

// ResponseFormatData is the data response format.
type ResponseFormatData struct {
	Data any `json:"data"`
}

func (ResponseFormatData) format() {}
