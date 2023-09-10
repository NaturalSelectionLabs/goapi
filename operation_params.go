package goapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	ff "github.com/NaturalSelectionLabs/goapi/lib/flat-fields"
	"github.com/NaturalSelectionLabs/jschema"
	"github.com/iancoleman/strcase"
)

type paramsIn int

const (
	inHeader paramsIn = iota
	inBody
	inURL
)

type Params interface {
	paramsIn() paramsIn
}

var tParams = reflect.TypeOf(new(Params)).Elem()

type InHeader struct{}

var _ Params = InHeader{}

func (InHeader) paramsIn() paramsIn {
	return inHeader
}

type InURL struct{}

var _ Params = InURL{}

func (InURL) paramsIn() paramsIn {
	return inURL
}

type InBody struct{}

var _ Params = InBody{}

func (InBody) paramsIn() paramsIn {
	return inBody
}

var tContext = reflect.TypeOf(new(context.Context)).Elem()

type parsedParam struct {
	in     paramsIn
	param  reflect.Type
	fields []*parsedField

	isContext bool
}

func parseParam(path *Path, p reflect.Type) *parsedParam {
	if p == tContext {
		return &parsedParam{isContext: true}
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
			fields = append(fields, parseHeaderField(f))
		}

	case inURL:
		parsed.in = inURL

		for _, f := range flat.Fields {
			fields = append(fields, parseURLField(path, f))
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
		parsed.in = inBody
	}

	parsed.fields = fields

	return parsed
}

var ErrMissingParam = errors.New("missing parameter in request")

func (p *parsedParam) loadURL(qs url.Values) (reflect.Value, error) {
	val := reflect.New(p.param)

	for _, f := range p.fields {
		var fv reflect.Value

		if !f.InPath && f.slice { //nolint: nestif
			vs, ok := qs[f.name]
			if ok { //nolint: gocritic
				fv = reflect.MakeSlice(f.sliceType, len(vs), len(vs))
			} else if f.hasDefault {
				fv = f.defaultVal
			} else {
				continue
			}

			for i, v := range vs {
				val, err := toValue(f.item, v)
				if err != nil {
					return reflect.Value{}, fmt.Errorf("failed to parse param `%s`: %w", f.name, err)
				}

				fv.Index(i).Set(val)
			}
		} else {
			vs, has := qs[f.name]
			if has { //nolint: gocritic
				var err error
				fv, err = toValue(f.item, vs[0])
				if err != nil {
					return reflect.Value{}, fmt.Errorf("failed to parse path param `%s`: %w", f.name, err)
				}
			} else if f.required {
				return reflect.Value{}, fmt.Errorf("%w, param: %s", ErrMissingParam, f.name)
			} else if f.hasDefault {
				fv = f.defaultVal
			}
		}

		if f.ptr && !f.slice {
			c := reflect.New(f.item)
			c.Elem().Set(fv)
			f.flatField.Set(val, c)
		} else {
			f.flatField.Set(val, fv)
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

	return val.Elem(), nil
}

type parsedField struct {
	name        string
	item        reflect.Type
	flatField   *ff.FlattenedField
	ptr         bool
	slice       bool
	sliceType   reflect.Type
	required    bool
	InPath      bool
	hasDefault  bool
	defaultVal  reflect.Value
	example     reflect.Value
	description string
}

func parseHeaderField(flatField *ff.FlattenedField) *parsedField {
	f := flatField.Field
	parsed := parseField(flatField)
	parsed.name = toHeaderName(f.Name)
	parsed.name = tagName(f.Tag, parsed.name)

	return parsed
}

func parseURLField(path *Path, flatField *ff.FlattenedField) *parsedField {
	f := parseField(flatField)

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

func parseField(flatField *ff.FlattenedField) *parsedField {
	t := flatField.Field
	f := &parsedField{flatField: flatField, required: true}
	tf := t.Type

	switch t.Type.Kind() { //nolint: exhaustive
	case reflect.Ptr, reflect.Slice:
		f.ptr = true
	default:
		f.ptr = false
	}

	if tf.Kind() == reflect.Ptr {
		tf = tf.Elem()
		f.required = false
	}

	if tf.Kind() == reflect.Slice {
		f.slice = true
		f.sliceType = tf
		f.item = tf.Elem()
		f.required = false
	} else {
		f.item = tf
	}

	if d, ok := t.Tag.Lookup("default"); ok {
		var v any
		if f.slice {
			v = reflect.New(f.sliceType).Interface()
		} else {
			v = reflect.New(f.item).Interface()
		}

		err := json.Unmarshal([]byte(d), &v)
		if err != nil {
			panic("failed to parse tag `default` of `" + t.Name + "`: " + err.Error())
		}

		f.required = false
		f.hasDefault = true
		f.defaultVal = reflect.Indirect(reflect.ValueOf(v))
	}

	if d, ok := t.Tag.Lookup("example"); ok {
		var v any

		err := json.Unmarshal([]byte(d), &v)
		if err != nil {
			panic("failed to parse tag `example` of `" + t.Name + "`: " + err.Error())
		}

		f.example = reflect.Indirect(reflect.ValueOf(v))
	}

	if d, ok := t.Tag.Lookup("description"); ok {
		f.description = d
	}

	return f
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
