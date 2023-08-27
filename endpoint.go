package goapi

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/iancoleman/strcase"
)

var (
	tHTTPResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	tHTTPRequest        = reflect.TypeOf((*http.Request)(nil))
	tParamDecoder       = reflect.TypeOf((*ParamDecoder)(nil)).Elem()
	tError              = reflect.TypeOf((*error)(nil)).Elem()
)

type Operation struct {
	method      openapi.Method
	openapiPath string
	vHandler    reflect.Value
	tHandler    reflect.Type

	pathPattern    *regexp.Regexp
	pathParams     []string
	overrideWriter bool

	meta *OperationMeta
}

type OperationMeta struct {
	// Summary is used for display in the openapi UI.
	Summary string
	// Description is used for display in the openapi UI.
	Description string
	// OperationID is a unique string used to identify an individual operation.
	// This can be used by tools and libraries to provide functionality for
	// referencing and calling the operation from different parts of your application.
	OperationID string
	// Tags are used for grouping operations together for display in the openapi UI.
	Tags []string

	ResponseDescription string
}

type ParamDecoder interface {
	DecodeParam([]string)
}

func (op *Operation) Handle(w http.ResponseWriter, r *http.Request, pathParams []string) {
	body := map[string]json.RawMessage{}

	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			writeResponse(w, &ResponseBadRequest{Error: toError(err)})
			return
		}

		if len(b) > 0 {
			err := json.Unmarshal(b, &body)
			if err != nil {
				writeResponse(w, &ResponseBadRequest{Error: toError(err)})
				return
			}
		}
	}

	args := make([]reflect.Value, op.tHandler.NumIn())

	for i := range args {
		tArg := op.tHandler.In(i)

		switch tArg {
		case tHTTPResponseWriter:
			args[i] = reflect.ValueOf(w)

		case tHTTPRequest:
			args[i] = reflect.ValueOf(r)

		default:
			var err error

			args[i], err = op.Params(tArg.Elem(), pathParams, r, body)
			if err != nil {
				writeResponse(w, &ResponseBadRequest{Error: toError(err)})
				return
			}
		}
	}

	ret := op.vHandler.Call(args)

	if op.overrideWriter {
		return
	}

	if len(ret) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if len(ret) > 0 {
		last := ret[len(ret)-1].Interface()

		if err, ok := last.(error); ok {
			if _, ok := last.(ErrorCode); ok {
				writeResponse(w, &ResponseBadRequest{Error: toError(err)})
			} else {
				writeResponse(w, &ResponseInternalServerError{Error: toError(err)})
			}

			return
		}
	}

	first := ret[0].Interface()

	if res, ok := first.(Response); ok {
		writeResponse(w, res)
		return
	}

	switch len(ret) {
	case 1:
		writeResponse(w, &ResponseOK{first, nil})

	default:
		writeResponse(w, &ResponseOK{first, ret[1].Interface()})
	}
}

func (op *Operation) Match(method, path string) []string {
	if method != op.method.String() {
		return nil
	}

	matches := op.pathPattern.FindStringSubmatch(path)

	if matches == nil {
		return nil
	}

	return matches[1:]
}

func (op *Operation) Params( //nolint: gocognit
	tArg reflect.Type,
	pathParams []string,
	r *http.Request,
	body map[string]json.RawMessage,
) (reflect.Value, error) {
	arg := reflect.New(tArg)

	params := make(map[string][]string, len(pathParams))

	for i, name := range op.pathParams {
		params[name] = []string{pathParams[i]}
	}

	for i := 0; i < tArg.NumField(); i++ {
		tField := tArg.Field(i)
		vField := arg.Elem().Field(i)

		switch inWhere(tField) {
		case InHeader:
			vs := r.Header.Values(toHeaderName(tField.Name))
			if len(vs) > 0 {
				assign(tField, vField, vs[0])
			}

			continue

		case InBody:
			if len(body) > 0 {
				if b, ok := body[toBodyName(tField.Name)]; ok {
					v := reflect.New(tField.Type)

					err := json.Unmarshal(b, v.Interface())
					if err != nil {
						return reflect.Value{}, err
					}

					setVal(vField, tField.Type.Kind() == reflect.Ptr, v.Elem())
				}

				continue
			}

		case InOthers:
		}

		qs := r.URL.Query()

		var field reflect.Value
		if tField.Type.Kind() == reflect.Ptr {
			field = reflect.New(tField.Type.Elem())
		} else {
			field = reflect.New(tField.Type)
		}

		if dec, ok := field.Interface().(ParamDecoder); ok {
			if vs, ok := params[toPathName(tField.Name)]; ok {
				dec.DecodeParam(vs)
			} else if vs, ok := qs[toQueryName(tField.Name)]; ok {
				dec.DecodeParam(vs)
			}

			if tField.Type.Kind() == reflect.Ptr {
				vField.Set(field)
			} else {
				vField.Set(field.Elem())
			}

			continue
		}

		if vs, ok := params[toPathName(tField.Name)]; ok {
			assign(tField, vField, vs[0])
		} else if vs, ok := qs[toQueryName(tField.Name)]; ok {
			if tField.Type.Kind() == reflect.Slice {
				vField.Set(reflect.MakeSlice(tField.Type, len(vs), len(vs)))
				for i, v := range vs {
					assign(tField, vField.Index(i), v)
				}
			} else {
				assign(tField, vField, vs[0])
			}
		}
	}

	return arg, nil
}

var regOpenAPIPath = regexp.MustCompile(`\{([^}]+)\}`)

// Converts OpenAPI style path to Go Regexp and returns path parameters.
func openAPIPathToRegexp(path string) (*regexp.Regexp, []string, error) {
	params := []string{}

	// Replace OpenAPI wildcards with Go RegExp named wildcards
	regexPath := regOpenAPIPath.ReplaceAllStringFunc(path, func(m string) string {
		param := m[1 : len(m)-1]          // Strip outer braces from parameter
		params = append(params, param)    // Add param to list
		return "(?P<" + param + ">[^/]+)" // Replace with Go Regexp named wildcard
	})

	// Make sure the path starts with a "^", ends with a "$", and escape slashes
	regexPath = "^" + strings.ReplaceAll(regexPath, "/", "\\/") + "$"

	// Compile the regular expression
	r, err := regexp.Compile(regexPath)
	if err != nil {
		return nil, nil, err
	}

	return r, params, nil
}

func assign(tField reflect.StructField, vField reflect.Value, val string) {
	kind := vField.Kind()
	ptr := kind == reflect.Ptr

	if ptr {
		kind = tField.Type.Elem().Kind()
	}

	switch kind { //nolint: exhaustive
	case reflect.Int:
		n, _ := strconv.ParseInt(val, 10, 64)
		setVal(vField, ptr, reflect.ValueOf(int(n)))

	case reflect.Float64:
		n, _ := strconv.ParseFloat(val, 64)
		setVal(vField, ptr, reflect.ValueOf(n))

	default:
		setVal(vField, ptr, reflect.ValueOf(val))
	}
}

func setVal(vField reflect.Value, ptr bool, val reflect.Value) {
	if ptr {
		p := reflect.New(val.Type())
		p.Elem().Set(val)
		val = p
	}

	vField.Set(val)
}

type InWhere int

const (
	InHeader InWhere = iota
	InBody
	InOthers
)

func inWhere(t reflect.StructField) InWhere {
	switch t.Tag.Get("in") {
	case "header":
		return InHeader
	case "body":
		return InBody
	default:
		return InOthers
	}
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

func toBodyName(name string) string {
	return strcase.ToLowerCamel(name)
}
