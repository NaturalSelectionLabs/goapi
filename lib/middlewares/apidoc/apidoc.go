// Package apidoc contains a middleware to serve the OpenAPI document.
package apidoc

import (
	"embed"
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

	g.GET("/", op)

	op.doc = config(g.OpenAPI())

	fs := http.FileServer(http.FS(swaggerFiles))

	g.Group("/swagger-ui").Use(middlewares.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, g.Prefix())
			fs.ServeHTTP(w, r)
		})
	}))
}

type params struct {
	goapi.InHeader

	Accept string `default:"application/json"`
}

type res interface {
	goapi.Response
}

var _ = goapi.Interface(new(res), resOK{}, resRedirect{})

type resOK struct {
	goapi.StatusOK
	Data any `response:"direct"`
}

func (resOK) Description() string {
	return "It will return the OpenAPI doc in JSON format."
}

type resRedirect struct {
	goapi.StatusFound
	Header headerRedirect
}

func (resRedirect) Description() string {
	return "It will redirect the browser to the Swagger UI to render the generated OpenAPI doc."
}

type headerRedirect struct {
	Location string `description:"Redirect to ./swagger-ui"`
}

// operation is the operation to serve the OpenAPI document.
type operation struct {
	doc *openapi.Document
}

var _ goapi.OperationOpenAPI = &operation{}

// OpenAPI implements the [goapi.OperationOpenAPI] interface.
func (*operation) OpenAPI(doc openapi.Operation) openapi.Operation {
	doc.Description = "It will auto redirect the browser to the Swagger UI to render the generated OpenAPI doc. " +
		"If you request it with `Accept: application/json` header, it will return the OpenAPI doc in JSON format."
	return doc
}

// Handle implements the [goapi.OperationHandler] interface.
func (op *operation) Handle(p params, _ *http.Request) res {
	if strings.Contains(p.Accept, "application/json") {
		return resOK{Data: op.doc}
	}

	return resRedirect{
		Header: headerRedirect{
			Location: "swagger-ui",
		},
	}
}
