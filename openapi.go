package goapi

import (
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/NaturalSelectionLabs/jschema"
)

func (r *Group) OpenAPI() *openapi.Document {
	return r.router.OpenAPI()
}

// OpenAPI returns the OpenAPI doc of the router.
// You can use [json.Marshal] to convert it to a JSON string.
func (r *Router) OpenAPI() *openapi.Document {
	doc := &openapi.Document{
		Paths: map[string]openapi.Path{},
	}
	schemas := jschema.New("#/components/schemas")

	doc.Components.Schemas = schemas.JSON()

	return doc
}
