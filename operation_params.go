package goapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	ff "github.com/NaturalSelectionLabs/goapi/lib/flat-fields"
	"github.com/NaturalSelectionLabs/jschema"
	"github.com/iancoleman/strcase"
	"github.com/xeipuuv/gojsonschema"
)

type paramsIn int

const (
	inHeader paramsIn = iota
	inBody
	inURL
)

// Params represents the parameter of a request.
type Params interface {
	paramsIn() paramsIn
}

var tParams = reflect.TypeOf(new(Params)).Elem()

// InHeader is a flag that can be embedded into a struct to mark it
// as a container for request header parameters.
type InHeader struct{}

var _ Params = InHeader{}

func (InHeader) paramsIn() paramsIn {
	return inHeader
}

// InURL is a flag that can be embedded into a struct to mark it
// as a container for request url parameters.
type InURL struct{}

var _ Params = InURL{}

func (InURL) paramsIn() paramsIn {
	return inURL
}

// InBody is a flag that can be embedded into a struct to mark it
// as a container for request body.
type InBody struct{}

var _ Params = InBody{}

func (InBody) paramsIn() paramsIn {
	return inBody
}

var tContext = reflect.TypeOf(new(context.Context)).Elem()

var tRequest = reflect.TypeOf((*http.Request)(nil))

type parsedParam struct {
	in     paramsIn
	param  reflect.Type
	fields []*parsedField

	isContext bool
	isRequest bool

	bodyValidator *gojsonschema.Schema
}

func (p *parsedParam) loadURL(qs url.Values) (reflect.Value, error) { //nolint: gocognit
	val := reflect.New(p.param)

	for _, f := range p.fields {
		var fv reflect.Value

		if !f.InPath && f.slice { //nolint: nestif
			vs, has := qs[f.name]
			if has { //nolint: gocritic
				fv = reflect.MakeSlice(f.sliceType, len(vs), len(vs))
			} else if f.hasDefault {
				fv = f.defaultVal
			} else {
				continue
			}

			for i, v := range vs {
				val, err := toValue(f.item, v)
				if err != nil {
					return reflect.Value{}, fmt.Errorf("failed to parse url param `%s`: %w", f.name, err)
				}

				fv.Index(i).Set(val)
			}
		} else {
			vs, has := qs[f.name]
			if has { //nolint: gocritic
				var err error
				fv, err = toValue(f.item, vs[0])
				if err != nil {
					return reflect.Value{}, fmt.Errorf("failed to parse url path param `%s`: %w", f.name, err)
				}
			} else if f.required {
				if !f.InPath {
					return reflect.Value{}, fmt.Errorf("missing url query param `%s`", f.name)
				}

				return reflect.Value{}, fmt.Errorf("missing url path param `%s`", f.name)
			} else if f.hasDefault {
				fv = f.defaultVal
			}
		}

		if f.ptr && !f.slice {
			if fv.IsValid() {
				c := reflect.New(f.item)
				c.Elem().Set(fv)
				f.flatField.Set(val, c)
			}
		} else {
			f.flatField.Set(val, fv)
		}

		if err := f.validate(val); err != nil {
			return reflect.Value{}, err
		}
	}

	return val.Elem(), nil
}

func (p *parsedParam) loadHeader(h http.Header) (reflect.Value, error) {
	qs := url.Values{}

	// h to qs
	for k, vs := range h {
		for _, v := range vs {
			qs.Add(toHeaderName(k), v)
		}
	}

	return p.loadURL(qs)
}

func (p *parsedParam) loadBody(body io.Reader) (reflect.Value, error) {
	val := reflect.New(p.param)
	ref := val.Interface()

	err := json.NewDecoder(body).Decode(&ref)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("failed to parse json body: %w", err)
	}

	check, _ := p.bodyValidator.Validate(gojsonschema.NewGoLoader(ref))
	if !check.Valid() {
		return reflect.Value{}, fmt.Errorf("request body is invalid: %v", check.Errors())
	}

	return val.Elem(), nil
}

type parsedField struct {
	name       string // the normalized name of the field
	item       reflect.Type
	flatField  *ff.FlattenedField
	ptr        bool
	slice      bool
	sliceType  reflect.Type
	required   bool
	InPath     bool
	hasDefault bool
	defaultVal reflect.Value
	example    reflect.Value

	schema    *jschema.Schema
	validator *gojsonschema.Schema
}

func (f *parsedField) validate(val reflect.Value) error {
	v := f.flatField.Get(val).Interface()
	res, _ := f.validator.Validate(gojsonschema.NewGoLoader(v))

	if !res.Valid() {
		return fmt.Errorf("param `%s` is invalid: %v", f.name, res.Errors())
	}

	return nil
}

func toOperationName(name string) string {
	return strcase.ToLowerCamel(name)
}

func toHeaderName(name string) string {
	return strcase.ToKebab(name)
}

func toPathName(name string) string {
	return strcase.ToKebab(name)
}

func toQueryName(name string) string {
	return strcase.ToSnake(name)
}

var tUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

// converts the val to the kind of value.
func toValue(t reflect.Type, val string) (reflect.Value, error) {
	if t.Kind() == reflect.String || t.Implements(tUnmarshaler) || reflect.New(t).Type().Implements(tUnmarshaler) {
		val = strconv.Quote(val)
	}

	v := reflect.New(t)

	err := json.Unmarshal([]byte(val), v.Interface())
	if err != nil {
		return reflect.Value{}, fmt.Errorf("can't parse `%s` to expected value, %w", val, err)
	}

	return v.Elem(), nil
}

func tagName(t reflect.StructTag, name string) string {
	tag := jschema.ParseJSONTag(t)

	if tag != nil && tag.Name != "" {
		return tag.Name
	}

	return name
}

func firstProp(s *jschema.Schema) (p *jschema.Schema) { //nolint: nonamedreturns
	for _, p = range s.Properties {
		break
	}

	return p
}

func parseParam(s jschema.Schemas, path *Path, p reflect.Type) *parsedParam {
	if p == tContext {
		return &parsedParam{isContext: true}
	}

	if p == tRequest {
		return &parsedParam{isRequest: true}
	}

	if p.Kind() != reflect.Struct || !p.Implements(tParams) {
		panic("expect parameter to be a struct and embedded with goapi.InHeader, goapi.InURL, or goapi.InBody," +
			" but got: " + p.String())
	}

	parsed := &parsedParam{param: p}
	fields := []*parsedField{}
	flat := ff.Parse(p)

	switch reflect.New(p).Interface().(Params).paramsIn() {
	case inHeader:
		parsed.in = inHeader

		for _, f := range flat.Fields {
			fields = append(fields, parseHeaderField(s, f))
		}

	case inURL:
		parsed.in = inURL

		for _, f := range flat.Fields {
			fields = append(fields, parseURLField(s, path, f))
		}

		for _, n := range path.names {
			has := false

			for _, f := range fields {
				if f.name == n {
					has = true
				}
			}

			if !has {
				panic("expect to have path parameter for {" + n + "} in " + p.String())
			}
		}

	case inBody:
		validateBodyParam(p)

		parsed.in = inBody

		scm := s.ToStandAlone(s.DefineT(p))

		validator, _ := gojsonschema.NewSchema(gojsonschema.NewGoLoader(scm))
		parsed.bodyValidator = validator
	}

	parsed.fields = fields

	return parsed
}

func parseHeaderField(s jschema.Schemas, flatField *ff.FlattenedField) *parsedField {
	f := flatField.Field
	parsed := parseField(s, flatField)
	parsed.name = toHeaderName(f.Name)
	parsed.name = tagName(f.Tag, parsed.name)

	return parsed
}

func parseURLField(s jschema.Schemas, path *Path, flatField *ff.FlattenedField) *parsedField {
	f := parseField(s, flatField)

	t := flatField.Field

	f.name = toPathName(t.Name)
	if path.contains(f.name) {
		if f.hasDefault {
			panic("path parameter cannot have tag `default`, param: " + t.Name)
		}

		if f.slice {
			panic("path parameter cannot be an slice, param: " + t.Name)
		}

		if !f.required {
			panic("path parameter cannot be optional, param: " + t.Name)
		}

		f.InPath = true
	} else {
		f.name = toQueryName(t.Name)
	}

	f.name = tagName(t.Tag, f.name)

	return f
}

func parseField(s jschema.Schemas, flatField *ff.FlattenedField) *parsedField {
	f := flatField.Field
	parsed := &parsedField{flatField: flatField, required: true}
	t := f.Type

	switch t.Kind() { //nolint: exhaustive
	case reflect.Ptr, reflect.Slice:
		parsed.ptr = true
	default:
		parsed.ptr = false
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		parsed.required = false
	}

	if t.Kind() == reflect.Slice {
		parsed.slice = true
		parsed.sliceType = t
		parsed.item = t.Elem()
		parsed.required = false
	} else {
		parsed.item = t
	}

	parsed.schema = s.ToStandAlone(fieldSchema(s, f))

	if _, ok := f.Tag.Lookup("default"); ok {
		parsed.required = false
		parsed.hasDefault = true
		parsed.defaultVal = reflect.ValueOf(parsed.schema.Default)
	}

	if _, ok := f.Tag.Lookup("example"); ok {
		parsed.example = reflect.ValueOf(parsed.schema.Example)
	}

	scm := parsed.schema

	if !parsed.required {
		s := &jschema.Schema{Defs: parsed.schema.Defs}
		s.AnyOf = []*jschema.Schema{parsed.schema, {Type: jschema.TypeNull}}
		scm = s
	}

	validator, _ := gojsonschema.NewSchema(gojsonschema.NewGoLoader(scm))
	parsed.validator = validator

	return parsed
}

func validateBodyParam(p reflect.Type) {
	p = indirectType(p)

	if p.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < p.NumField(); i++ {
		f := p.Field(i)
		if _, has := f.Tag.Lookup(string(jschema.JTagDefault)); has {
			panic(fmt.Sprintf("goapi.InBody field `%s.%s` don't support default field tag, "+
				"it should be treated as common json field, "+
				"if you want to make this field as optional you should use omitempty or pointer type",
				p.String(),
				f.Name,
			))
		}
	}
}

func indirectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}

	return t
}
