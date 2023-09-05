package goapi

import (
	"reflect"

	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/NaturalSelectionLabs/jschema"
	"github.com/ysmood/vary"
)

type Description interface {
	Description() string
}

var tDescription = reflect.TypeOf((*Description)(nil)).Elem()

func (r *Group) OpenAPI(schemas *jschema.Schemas) *openapi.Document {
	return r.router.OpenAPI(schemas)
}

// OpenAPI returns the OpenAPI doc of the router.
// You can use [json.Marshal] to convert it to a JSON string.
func (r *Router) OpenAPI(schemas *jschema.Schemas) *openapi.Document {
	doc := &openapi.Document{
		Paths: map[string]openapi.Path{},
	}

	if schemas == nil {
		s := jschema.New("#/components/schemas")
		schemas = &s
	}

	for _, m := range r.middlewares {
		op, ok := m.(*Operation)
		if !ok {
			continue
		}

		if _, has := doc.Paths[op.path.path]; !has {
			doc.Paths[op.path.path] = openapi.Path{}
		}

		doc.Paths[op.path.path][op.method] = operationDoc(*schemas, op)
	}

	doc.Components.Schemas = schemas.JSON()

	return doc
}

func operationDoc(s jschema.Schemas, op *Operation) openapi.Operation {
	doc := openapi.Operation{
		Parameters: []openapi.Parameter{},
		Responses:  map[openapi.StatusCode]openapi.Response{},
	}

	if op.meta != nil {
		doc.Summary = op.meta.Summary
		doc.Description = op.meta.Description
		doc.OperationID = op.meta.OperationID
		doc.Tags = op.meta.Tags
	}

	for _, p := range op.params {
		var params []openapi.Parameter

		switch p.in {
		case inHeader:
			params = append(params, headerParamDoc(s, p)...)

		case inURL:
			params = append(params, urlParamDoc(s, p)...)

		case inBody:
			doc.RequestBody = &openapi.RequestBody{
				Content: &openapi.Content{
					JSON: &openapi.Schema{
						Schema: s.DefineT(p.param),
					},
				},
			}
		}

		doc.Parameters = append(doc.Parameters, params...)
	}

	doc.Responses = resDoc(s, op)

	return doc
}

func urlParamDoc(s jschema.Schemas, p *parsedParam) []openapi.Parameter {
	arr := []openapi.Parameter{}

	for _, f := range p.fields {
		if f.skip {
			continue
		}

		in := openapi.QUERY

		if f.InPath {
			in = openapi.PATH
		}

		arr = append(arr, openapi.Parameter{
			Name:        f.name,
			In:          in,
			Schema:      s.DefineT(f.item),
			Description: f.description,
			Required:    f.required,
		})
	}

	return arr
}

func headerParamDoc(s jschema.Schemas, p *parsedParam) []openapi.Parameter {
	arr := []openapi.Parameter{}

	for _, f := range p.fields {
		if f.skip {
			continue
		}

		arr = append(arr, openapi.Parameter{
			Name:        f.name,
			In:          openapi.HEADER,
			Schema:      s.DefineT(f.item),
			Description: f.description,
			Required:    f.required,
		})
	}

	return arr
}

func resDoc(s jschema.Schemas, op *Operation) map[openapi.StatusCode]openapi.Response {
	list := map[openapi.StatusCode]openapi.Response{}

	add := func(t reflect.Type) {
		parsedRes := parseResponse(t)
		res := openapi.Response{
			Description: getDescription(t),
			Headers:     resHeaderDoc(s, parsedRes.header),
			Content: &openapi.Content{
				JSON: &openapi.Schema{
					Schema: s.DefineT(t),
				},
			},
		}

		list[openapi.StatusCode(parsedRes.statusCode)] = res
	}

	if it := vary.Get(vary.NewID(op.tRes)); it != nil {
		for _, t := range it.Implementations {
			add(t)
		}
	} else {
		add(op.tRes)
	}

	return list
}

func resHeaderDoc(s jschema.Schemas, t reflect.Type) openapi.Headers {
	if t == nil {
		return nil
	}

	parsed := parseParam(nil, t)

	headers := openapi.Headers{}

	for _, f := range parsed.fields {
		if f.skip {
			continue
		}

		headers[f.name] = openapi.Header{
			Description: f.description,
			Schema:      s.DefineT(f.item),
		}
	}

	return headers
}

func getDescription(t reflect.Type) string {
	if t.Implements(tDescription) {
		return reflect.New(t).Elem().Interface().(Description).Description()
	}

	return ""
}
