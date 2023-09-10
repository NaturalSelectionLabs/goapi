package goapi

import (
	"net/http"
	"reflect"

	ff "github.com/NaturalSelectionLabs/goapi/lib/flat-fields"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/NaturalSelectionLabs/jschema"
	"github.com/ysmood/vary"
)

var interfaces = vary.NewInterfaces()

func Interface(i any, ts ...any) *vary.Interface {
	return interfaces.New(i, ts...)
}

type Descriptioner interface {
	Description() string
}

var tDescriptioner = reflect.TypeOf((*Descriptioner)(nil)).Elem()

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
		s := jschema.NewWithInterfaces("#/components/schemas", interfaces)
		schemas = &s
	}

	for _, op := range r.operations {
		if op.override != nil {
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

	doc.Summary = op.meta.Summary
	doc.Description = op.meta.Description
	doc.OperationID = op.meta.OperationID
	doc.Tags = op.meta.Tags
	doc.Security = op.meta.Security

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
		in := openapi.QUERY

		if f.InPath {
			in = openapi.PATH
		}

		arr = append(arr, openapi.Parameter{
			Name:        f.name,
			In:          in,
			Schema:      fieldSchema(s, f),
			Description: f.description,
			Required:    f.required,
		})
	}

	return arr
}

func headerParamDoc(s jschema.Schemas, p *parsedParam) []openapi.Parameter {
	arr := []openapi.Parameter{}

	for _, f := range p.fields {
		arr = append(arr, openapi.Parameter{
			Name:        f.name,
			In:          openapi.HEADER,
			Schema:      fieldSchema(s, f),
			Description: f.description,
			Required:    f.required,
		})
	}

	return arr
}

func fieldSchema(s jschema.Schemas, f *parsedField) *jschema.Schema {
	scm := s.DefineT(f.item)

	raw := s.PeakSchema(scm)

	if f.defaultVal.IsValid() {
		raw.Default = f.defaultVal.Interface()
	}

	if f.example.IsValid() {
		raw.Example = f.example.Interface()
	}

	return scm
}

func resDoc(s jschema.Schemas, op *Operation) map[openapi.StatusCode]openapi.Response {
	list := map[openapi.StatusCode]openapi.Response{}

	add := func(t reflect.Type) {
		parsedRes := parseResponse(t)

		var content *openapi.Content

		if parsedRes.hasData || parsedRes.hasErr {
			scm := &jschema.Schema{
				Type:                 jschema.TypeObject,
				AdditionalProperties: ptr(false),
				Properties:           jschema.Properties{},
			}

			if parsedRes.hasErr { //nolint: gocritic
				scm.Properties["error"] = s.DefineT(parsedRes.err)
				scm.Required = []string{"error"}
			} else if parsedRes.hasMeta {
				scm.Properties["data"] = s.DefineT(parsedRes.data)
				scm.Properties["meta"] = s.DefineT(parsedRes.meta)
				scm.Required = []string{"data", "meta"}
			} else {
				scm.Properties["data"] = s.DefineT(parsedRes.data)
				scm.Required = []string{"data"}
			}

			content = &openapi.Content{
				JSON: &openapi.Schema{
					Schema: scm,
				},
			}
		}

		code := openapi.StatusCode(parsedRes.statusCode)

		res := openapi.Response{
			Description: getDescription(t, code),
			Headers:     resHeaderDoc(s, parsedRes.header),
			Content:     content,
		}

		list[code] = res
	}

	if it, has := interfaces[vary.ID(op.tRes)]; has {
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

	headers := openapi.Headers{}

	for _, flat := range ff.Parse(t).Fields {
		f := parseHeaderField(flat)
		headers[f.name] = openapi.Header{
			Description: f.description,
			Schema:      s.DefineT(f.item),
		}
	}

	return headers
}

func getDescription(t reflect.Type, code openapi.StatusCode) string {
	if t.Implements(tDescriptioner) {
		return reflect.New(t).Elem().Interface().(Descriptioner).Description()
	}

	return http.StatusText(int(code))
}

func ptr[T any](v T) *T {
	return &v
}
