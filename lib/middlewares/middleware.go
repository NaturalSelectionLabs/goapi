package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

type Middleware interface {
	// A Middleware https://cs.opensource.google/go/x/pkgsite/+/68be0dd1:internal/middleware/middleware.go
	Handler(next http.Handler) http.Handler
}

type Func func(next http.Handler) http.Handler

func (fn Func) Handler(next http.Handler) http.Handler {
	return fn(next)
}

// Identity is a middleware that does nothing.
var Identity = Func(func(next http.Handler) http.Handler {
	return next
})

func ResponseError(w http.ResponseWriter, code int, err *openapi.Error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(openapi.ResponseFormatErr{Error: err})
}

func Chain(ms ...Middleware) Middleware {
	return Func(func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			next = ms[i].Handler(next)
		}

		return next
	})
}
