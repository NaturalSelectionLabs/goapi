package bench_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
)

type Res interface {
	goapi.Response
}

var iRes = goapi.Interface(new(Res))

type ResOK struct {
	goapi.StatusOK
	Data string
}

var _ = iRes.Add(ResOK{})

func BenchmarkRes(b *testing.B) {
	r := goapi.New()

	r.GET("/hello", func() ResOK {
		return ResOK{Data: "World"}
	})

	go http.ListenAndServe(":3000", r.Server()) //nolint: errcheck

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res, err := http.Get("http://localhost:3000/hello") //nolint: noctx
		if err != nil {
			panic(err)
		}

		_, _ = io.ReadAll(res.Body)
		_ = res.Body.Close()
	}
}

func BenchmarkResInterface(b *testing.B) {
	r := goapi.New()

	r.GET("/hello", func() Res {
		return ResOK{Data: "World"}
	})

	go http.ListenAndServe(":3000", r.Server()) //nolint: errcheck

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res, err := http.Get("http://localhost:3000/hello") //nolint: noctx
		if err != nil {
			panic(err)
		}

		_, _ = io.ReadAll(res.Body)
		_ = res.Body.Close()
	}
}
