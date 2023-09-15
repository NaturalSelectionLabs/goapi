package goapi

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/ysmood/vary"
)

// Operation is a handler for a specific HTTP method and path.
// We use reflection to constrain the handler function signature,
// to make it follow the openapi spec.
type Operation struct {
	group  *Group
	method openapi.Method
	path   *Path

	vHandler reflect.Value
	tHandler reflect.Type
	params   []*parsedParam

	tRes reflect.Type

	override http.HandlerFunc

	openapi OperationOpenAPI
}

func (g *Group) newOperation(method openapi.Method, path string, handler OperationHandler) *Operation {
	p, err := newPath(path)
	if err != nil {
		panic(err)
	}

	if h, ok := handler.(func(http.ResponseWriter, *http.Request)); ok {
		return &Operation{
			group:    g,
			method:   method,
			path:     p,
			override: h,
		}
	}

	vHandler := reflect.ValueOf(handler)
	tHandler := vHandler.Type()

	var doc OperationOpenAPI

	if _, has := tHandler.MethodByName("Handle"); has {
		vHandler = vHandler.MethodByName("Handle")
		tHandler = vHandler.Type()

		doc, _ = handler.(OperationOpenAPI)
	} else if tHandler.Kind() != reflect.Func {
		panic("handler must be a function or a struct with Handle method")
	}

	params := []*parsedParam{}
	for i := 0; i < tHandler.NumIn(); i++ {
		params = append(params, g.parseParam(p, tHandler.In(i)))
	}

	if tHandler.NumOut() != 1 {
		panic("handler must return a single value")
	}

	tRes := tHandler.Out(0)

	return &Operation{
		group:    g,
		method:   method,
		path:     p,
		vHandler: vHandler,
		tHandler: tHandler,
		params:   params,
		tRes:     tRes,
		openapi:  doc,
	}
}

// Handler implements the [middlewares.Middleware] interface.
func (op *Operation) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != op.method.String() {
			next.ServeHTTP(w, r)
			return
		}

		params := op.path.match(r.URL.Path)
		if params == nil {
			next.ServeHTTP(w, r)
			return
		}

		if op.override != nil {
			op.override(w, r)
			return
		}

		qs := r.URL.Query()
		for k, v := range params {
			qs.Set(k, v)
		}

		op.handle(w, r, qs)
	})
}

func (op *Operation) handle(w http.ResponseWriter, r *http.Request, qs url.Values) {
	params := []reflect.Value{}

	for _, p := range op.params {
		if p.isContext {
			params = append(params, reflect.ValueOf(r.Context()))

			continue
		}

		if p.isRequest {
			params = append(params, reflect.ValueOf(r))

			continue
		}

		var param reflect.Value

		var err error

		switch p.in {
		case inHeader:
			param, err = p.loadHeader(r.Header)
		case inURL:
			param, err = p.loadURL(qs)
		case inBody:
			param, err = p.loadBody(r.Body)
		}

		if err != nil {
			middlewares.ResponseError(w, http.StatusBadRequest, &openapi.Error{
				Code:    openapi.CodeInvalidParam,
				Message: err.Error(),
			})

			return
		}

		params = append(params, param)
	}

	res := op.vHandler.Call(params)[0]

	resType := res.Type()
	if resType.Kind() == reflect.Interface {
		setType := resType
		res = res.Elem()
		resType = res.Type()

		if _, ok := Interfaces[vary.ID(setType)]; !ok {
			panic(fmt.Sprintf("handler response of path `%s` must goapi.Interface(new(%s))", op.path.path, setType.String()))
		}

		if _, ok := Interfaces[vary.ID(setType)].Implementations[vary.ID(resType)]; !ok {
			panic(fmt.Sprintf("handler response of path `%s` must goapi.Interface(new(%s), %s{})",
				op.path.path, setType.String(), resType.String()))
		}
	}

	op.parseResponse(resType).write(w, res)
}

// OperationHandler is a function to handle input and output of a http operation.
type OperationHandler any

// OperationOpenAPI allows a handler customize the OpenAPI doc of its corresponding operation.
type OperationOpenAPI interface {
	OpenAPI(doc openapi.Operation) openapi.Operation
}
