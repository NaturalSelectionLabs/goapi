// Package apidoc contains a middleware to serve the OpenAPI document.
package apidoc

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

//go:embed swagger-ui
var swaggerFiles embed.FS

// Install the several endpoints to serve the openapi document for g.
// If config can be nil if you don't want to modify the generated doc.
func Install(g *goapi.Group, config func(doc *openapi.Document) *openapi.Document) {
	if config == nil {
		config = func(doc *openapi.Document) *openapi.Document { return doc }
	}

	op := &operation{}

	g.GET("/openapi.json", op)

	op.doc = config(g.OpenAPI())

	dir, _ := fs.Sub(swaggerFiles, "swagger-ui")

	fs := http.FileServer(http.FS(dir))

	g.Use(middlewares.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, g.Prefix())
			fs.ServeHTTP(w, r)
		})
	}))
}

type resOK struct {
	goapi.StatusOK
	Data any `response:"direct"`
}

func (resOK) Description() string {
	return "It will return the OpenAPI doc in JSON format."
}

// operation is the operation to serve the OpenAPI document.
type operation struct {
	doc *openapi.Document
}

var _ goapi.OperationOpenAPI = &operation{}
var _ goapi.OperationHandler = &operation{}

// OpenAPI implements the [goapi.OperationOpenAPI] interface.
func (*operation) OpenAPI(doc openapi.Operation) openapi.Operation {
	doc.Description = "It responds the OpenAPI doc for this service in JSON format."
	return doc
}

// Handle implements the [goapi.OperationHandler] interface.
func (op *operation) Handle() resOK {
	return resOK{Data: op.doc}
}
