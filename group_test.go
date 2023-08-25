package goapi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

type MyResponse goapi.ResponseOK

func (r *MyResponse) StatusCode() int {
	return http.StatusNotModified
}

func TestGroup(t *testing.T) {
	g := got.T(t)

	r := goapi.New()

	type Params struct {
		ID string
	}

	g.Eq(g.Panic(func() {
		r.GET("/users/{userID}", func() {})
	}), "path parameter must be in kebab-case: userID")

	r.GET("/users/{id}", func(params *Params) (string, string) {
		return params.ID, params.ID
	})

	r.POST("/error", func() error {
		return fmt.Errorf("error")
	})

	r.PUT("/override-res", func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNotModified)
	})

	r.GET("/override-header", func() goapi.Response {
		return &goapi.ResponseOKHeader{
			ResHeader: http.Header{
				"x-test": []string{"ok"},
			},
		}
	})

	r.GET("/posts/{id}", func(params *Params) goapi.Response {
		return &goapi.ResponseOK{
			Data: params.ID,
		}
	})

	tr := g.Serve()
	tr.Mux.Handle("/", r)

	g.Eq(g.Req("", tr.URL("/users/123?user_filter=1")).JSON(), map[string]interface{}{
		"data": "123",
		"meta": "123",
	})

	g.Eq(g.Req("", tr.URL("/error")).StatusCode, http.StatusNotFound)

	g.Eq(g.Req(http.MethodPost, tr.URL("/error")).JSON(), map[string]interface{}{
		"error": map[string]interface{} /* len=2 */ {
			"code":    `*errors.errorString`, /* len=19 */
			"message": "error",
		},
	})

	g.Eq(g.Req(http.MethodPut, tr.URL("/override-res")).StatusCode, http.StatusNotModified)

	g.Eq(g.Req("", tr.URL("/users/123?userFilter=1")).JSON(), map[string]interface{}{
		"error": map[string]interface{} /* len=2 */ {
			"code":    `*errors.errorString`,                       /* len=19 */
			"message": `query key is not snake styled: userFilter`, /* len=41 */
		},
	})

	g.Eq(g.Req("", tr.URL("/users/123?userFilter=1")).JSON(), map[string]interface{}{
		"error": map[string]interface{} /* len=2 */ {
			"code":    `*errors.errorString`,                       /* len=19 */
			"message": `query key is not snake styled: userFilter`, /* len=41 */
		},
	})

	g.Eq(g.Req("", tr.URL("/override-header")).Header.Get("x-test"), "ok")
}
