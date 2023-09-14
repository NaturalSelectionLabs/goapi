package apidoc

import (
	"embed"
	"net/http"
	"strings"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/NaturalSelectionLabs/jschema"
)

//go:embed swagger-ui
var swaggerFiles embed.FS

func Install(g *goapi.Group, schemas *jschema.Schemas, config func(doc *openapi.Document) *openapi.Document) {
	op := &Operation{}

	g.GET("/", op)

	op.doc = config(g.OpenAPI())

	fs := http.FileServer(http.FS(swaggerFiles))

	g.Group("/swagger-ui").Use(middlewares.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.TrimLeft(r.URL.Path, g.Prefix())
			fs.ServeHTTP(w, r)
		})
	}))
}

type params struct {
	goapi.InHeader

	Accept string
}

type res interface {
	goapi.Response
}

var _ = goapi.Interface(new(res), resOK{}, resRedirect{})

type resOK struct {
	goapi.StatusOK
	Data any `response:"direct"`
}

type resRedirect struct {
	goapi.StatusFound
	Header headerRedirect
}

type headerRedirect struct {
	Location string
}

type Operation struct {
	doc *openapi.Document
}

var _ goapi.OperationOpenAPI = &Operation{}

func (*Operation) OpenAPI(doc openapi.Operation) openapi.Operation {
	doc.Description = "It will auto redirect the browser to the Swagger UI to render the generated OpenAPI doc. " +
		"If you request it with `Accept: application/json` header, it will return the OpenAPI doc in JSON format."
	return doc
}

func (op *Operation) Handle(p params, r *http.Request) res {
	if strings.Contains(p.Accept, "application/json") {
		return resOK{Data: op.doc}
	}

	return resRedirect{
		Header: headerRedirect{
			Location: "swagger-ui",
		},
	}
}
