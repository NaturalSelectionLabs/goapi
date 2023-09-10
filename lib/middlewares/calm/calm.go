package calm

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

type Calm struct {
	PrintStack bool
	Logger     *slog.Logger
}

var _ middlewares.Middleware = (*Calm)(nil)

func New() *Calm {
	return &Calm{
		PrintStack: true,
		Logger:     slog.Default(),
	}
}

// A middleware to recover from panic.
func (c *Calm) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprint(err)
				if c.PrintStack {
					c.Logger.Error(msg, "stack", string(debug.Stack()))
				}
				middlewares.ResponseError(w, http.StatusInternalServerError, &openapi.Error{
					Code:    openapi.CodeInternalError,
					Message: msg,
				})
			}
		}()

		h.ServeHTTP(w, rq)
	})
}
