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

	meta OperationMeta
}

type ConfigOperation func(op *Operation)

// Description for the operation.
func Description(d string) ConfigOperation {
	return func(op *Operation) { op.meta.Description = d }
}

// OperationID for the operation.
func OperationID(id string) ConfigOperation {
	return func(op *Operation) { op.meta.OperationID = id }
}

// Summary for the operation.
func Summary(s string) ConfigOperation {
	return func(op *Operation) { op.meta.Summary = s }
}

// Tags for the operation.
func Tags(tags ...string) ConfigOperation {
	return func(op *Operation) { op.meta.Tags = tags }
}

// Security for the operation.
func Security(security ...map[string][]string) ConfigOperation {
	return func(op *Operation) { op.meta.Security = security }
}

func (g *Group) newOperation(method openapi.Method, path string, handler any) *Operation {
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

	if tHandler.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	params := []*parsedParam{}
	for i := 0; i < tHandler.NumIn(); i++ {
		params = append(params, parseParam(p, tHandler.In(i)))
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
	}
}

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

		if _, ok := interfaces[vary.ID(setType)].Implementations[vary.ID(resType)]; !ok {
			panic(fmt.Sprintf("%s should goapi.Interface(new(%s), %s{})", op.path.path, setType.String(), resType.String()))
		}
	}

	op.parseResponse(resType).write(w, res)
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
	// Security is a declaration of which security mechanisms can be used for this operation.
	Security []map[string][]string
}
