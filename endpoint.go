package goapi

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"

	"github.com/iancoleman/strcase"
)

var (
	tHTTPResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	tHTTPRequest        = reflect.TypeOf((*http.Request)(nil))
	tError              = reflect.TypeOf((*error)(nil)).Elem()
)

type Endpoint struct {
	method      string
	openapiPath string
	vHandler    reflect.Value
	tHandler    reflect.Type
	tags        []Tag

	pathPattern *regexp.Regexp
	pathParams  []string
}

func (e *Endpoint) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	pathParams := e.Match(r.Method, r.URL.Path)

	if pathParams == nil {
		next(w, r)
		return
	}

	args := make([]reflect.Value, e.tHandler.NumIn())

	for i := range args {
		tArg := e.tHandler.In(i)

		switch tArg {
		case tHTTPResponseWriter:
			args[i] = reflect.ValueOf(w)

		case tHTTPRequest:
			args[i] = reflect.ValueOf(r)

		default:
			if tArg.Kind() != reflect.Ptr || tArg.Elem().Kind() != reflect.Struct {
				panic("Params argument must be a pointer to a struct: " + tArg.String())
			}

			args[i] = e.Params(tArg.Elem(), pathParams, r.URL.Query())
		}
	}

	ret := e.vHandler.Call(args)

	switch len(ret) {
	case 0:
		w.WriteHeader(http.StatusNoContent)

	case 1:
		if ret[0].IsNil() {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		writeResponse(w, nil, nil, ret[0].Interface().(error))

	case 2:
		if !ret[1].IsNil() {
			writeResponse(w, nil, nil, ret[0].Interface().(error))
			return
		}

		writeResponse(w, ret[0].Interface(), nil, nil)

	default:
		if !ret[2].IsNil() {
			writeResponse(w, nil, nil, ret[0].Interface().(error))
			return
		}

		writeResponse(w, ret[0].Interface(), ret[1].Interface(), nil)
	}
}

func (e *Endpoint) Match(method, path string) []string {
	if method != e.method {
		return nil
	}

	matches := e.pathPattern.FindStringSubmatch(path)

	if matches == nil {
		return nil
	}

	return matches[1:]
}

func (e *Endpoint) Params(tArg reflect.Type, pathParams []string, qs url.Values) reflect.Value {
	arg := reflect.New(tArg)

	params := make(map[string][]string, len(pathParams))

	for i, name := range e.pathParams {
		params[name] = []string{pathParams[i]}
	}

	for i := 0; i < tArg.NumField(); i++ {
		tField := tArg.Field(i)
		vField := arg.Elem().Field(i)

		if v, ok := params[strcase.ToKebab(tField.Name)]; ok {
			vField.SetString(v[0])
		} else if v, ok := qs[strcase.ToSnake(tField.Name)]; ok {
			if tField.Type.Kind() == reflect.Slice {
				vField.Set(reflect.ValueOf(v))
			} else {
				vField.SetString(v[0])
			}
		}
	}

	return arg
}

func (e *Endpoint) SetTags(tags ...Tag) *Endpoint {
	e.tags = tags
	return e
}

type Tag struct {
	name string
}
