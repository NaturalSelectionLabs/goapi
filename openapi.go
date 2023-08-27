package goapi

import (
	"reflect"

	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/NaturalSelectionLabs/jschema"
)

func (r *Group) OpenAPI() *openapi.Document {
	return r.router.OpenAPI()
}

// OpenAPI returns the OpenAPI doc of the router.
// You can use [json.Marshal] to convert it to a JSON string.
func (r *Router) OpenAPI() *openapi.Document { //nolint:gocognit
	doc := &openapi.Document{
		Paths: map[string]openapi.Path{},
	}
	schemas := jschema.New("#/components/schemas")

	for _, m := range r.middlewares {
		g, ok := m.(*Group)
		if !ok {
			continue
		}

		for _, op := range g.operations {
			var opDoc openapi.Operation

			var body *jschema.Schema

			if op.meta != nil {
				opDoc = openapi.Operation{
					Summary:     op.meta.Summary,
					Description: op.meta.Description,
					OperationID: op.meta.OperationID,
					Tags:        op.meta.Tags,
				}
			}

			numIn := op.tHandler.NumIn()
			for i := 0; i < numIn; i++ {
				tArg := op.tHandler.In(i)

				switch tArg {
				case tHTTPResponseWriter, tHTTPRequest:
				default:
					tArg = tArg.Elem()
					params := []openapi.Parameter{}

					for i := 0; i < tArg.NumField(); i++ {
						tField := tArg.Field(i)

						var param openapi.Parameter

						switch inWhere(tField) {
						case InHeader:
							param = openapi.Parameter{
								Name:   toHeaderName(tField.Name),
								In:     openapi.HEADER,
								Schema: schemas.DefineT(tField.Type),
							}

							continue

						case InBody:
							if opDoc.RequestBody == nil {
								body = &jschema.Schema{
									Type:       jschema.TypeObject,
									Properties: map[string]*jschema.Schema{},
								}

								opDoc.RequestBody = &openapi.RequestBody{
									Content: &openapi.Content{
										JSON: &openapi.Schema{
											Schema: body,
										},
									},
								}
							}

							body.Properties[toBodyName(tField.Name)] = schemas.DefineT(tField.Type)

						case InOthers:
						}

						_, pathParams, _ := openAPIPathToRegexp(op.openapiPath)

						if has(pathParams, toPathName(tField.Name)) {
							param = openapi.Parameter{
								Name:   toPathName(tField.Name),
								In:     openapi.PATH,
								Schema: schemas.DefineT(tField.Type),
							}
						} else {
							param = openapi.Parameter{
								Name:   toQueryName(tField.Name),
								In:     openapi.QUERY,
								Schema: schemas.DefineT(tField.Type),
							}
						}

						param.Required = tField.Type.Kind() != reflect.Ptr
						param.Description = getDesc(tField)

						params = append(params, param)
					}

					opDoc.Parameters = params
				}
			}

			opDoc.Responses = map[openapi.StatusCode]openapi.Response{}

			numOut := op.tHandler.NumOut()
			if op.tHandler.NumOut() > 0 {
				last := op.tHandler.Out(numOut - 1)

				if last.Implements(tError) {

				}
			}

			switch op.tHandler.NumOut() {
			case 0:
				opDoc.Responses[openapi.Status204] = openapi.Response{
					Description: op.meta.ResponseDescription,
				}
			case 1:
			case 2:
			case 3:
			}

			if doc.Paths[op.openapiPath] == nil {
				doc.Paths[op.openapiPath] = openapi.Path{}
			}

			doc.Paths[op.openapiPath][op.method] = opDoc
		}
	}

	doc.Components.Schemas = schemas.JSON()

	return doc
}

func has[T comparable](s []T, v T) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}

	return false
}

func getDesc(t reflect.StructField) string {
	return t.Tag.Get("description")
}
