package goapi_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

func TestOpenAPI(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	tr := g.Serve()
	tr.Mux.Handle("/", r)

	r.GET("/test", func() string { return "test" })

	g.Eq(r.OpenAPI().JSON(), nil)

	g.Eq(1, 1)
}
