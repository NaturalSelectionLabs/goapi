package goapi

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"

	"github.com/iancoleman/strcase"
)

var (
	tHTTPResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	tHTTPRequest        = reflect.TypeOf((*http.Request)(nil))
	tResponse           = reflect.TypeOf((*Response)(nil)).Elem()
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
	overrideWriter := false

	for i := range args {
		tArg := e.tHandler.In(i)

		switch tArg {
		case tHTTPResponseWriter:
			overrideWriter = true
			args[i] = reflect.ValueOf(w)

		case tHTTPRequest:
			args[i] = reflect.ValueOf(r)

		default:
			if tArg.Kind() != reflect.Ptr || tArg.Elem().Kind() != reflect.Struct {
				panic("Params argument must be a pointer to a struct: " + tArg.String())
			}

			qs := r.URL.Query()

			if err := e.guardQuery(qs); err != nil {
				writeResponse(w, &ResponseError{Error: toError(err)})
				return
			}

			args[i] = e.Params(tArg.Elem(), pathParams, qs)
		}
	}

	ret := e.vHandler.Call(args)

	if overrideWriter {
		return
	}

	if len(ret) > 0 {
		last := ret[len(ret)-1]

		if last.Type() == tError {
			if !last.IsNil() {
				writeResponse(w, &ResponseError{Error: toError(ret[0].Interface().(error))})
				return
			}

			ret = ret[:len(ret)-1]
		}
	}

	if ret[0].Type() == tResponse {
		writeResponse(w, ret[0].Interface().(Response))
		return
	}

	switch len(ret) {
	case 0:
		w.WriteHeader(http.StatusNoContent)

	case 1:
		writeResponse(w, &ResponseOK{ret[0].Interface(), nil})

	default:
		writeResponse(w, &ResponseOK{ret[0].Interface(), ret[1].Interface()})
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

// Ensure query key are camelCased.
func (e *Endpoint) guardQuery(qs url.Values) error {
	for k := range qs {
		if k != strcase.ToSnake(k) {
			return fmt.Errorf("query key is not snake styled: %s", k)
		}
	}

	return nil
}

func (e *Endpoint) SetTags(tags ...Tag) *Endpoint {
	e.tags = tags
	return e
}

type Tag struct {
	name string
}
