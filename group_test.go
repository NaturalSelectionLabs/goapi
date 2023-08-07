package goapi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

func TestGroup(t *testing.T) {
	g := got.T(t)

	r := goapi.New()

	type Params struct {
		ID string
	}

	g.Eq(g.Panic(func() {
		r.GET("/users/{userID}", func() {})
	}), "path parameter must be in kebab-case: userID")

	r.GET("/users/{id}", func(params *Params) (string, error) {
		return params.ID, nil
	})

	r.POST("/error", func() error {
		return fmt.Errorf("error")
	})

	tr := g.Serve()
	tr.Mux.Handle("/", r)

	g.Eq(g.Req("", tr.URL("/users/123")).JSON(), map[string]interface{}{
		"data": "123",
	})

	g.Eq(g.Req("", tr.URL("/error")).StatusCode, http.StatusNotFound)

	g.Eq(g.Req(http.MethodPost, tr.URL("/error")).JSON(), map[string]interface{}{
		"error": map[string]interface{} /* len=2 */ {
			"code":    `*errors.errorString`, /* len=19 */
			"message": "error",
		},
	})
}
