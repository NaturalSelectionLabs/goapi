// go test -bench=. -benchmem -cpuprofile profile.out ./lib/bench

package bench_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ResGoapi struct {
	goapi.StatusOK
	Data string
}

type ParamsGoapi struct {
	goapi.InURL
	ID      int
	Keyword string
}

func Benchmark_goapi(b *testing.B) {
	r := goapi.New()

	r.GET("/users/{id}/posts", func(p ParamsGoapi) ResGoapi {
		return ResGoapi{Data: fmt.Sprintf("%d %s", p.ID, p.Keyword)}
	})

	go func() { _ = r.Start(":3000") }()
	b.Cleanup(func() { _ = r.Shutdown(context.Background()) })

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req("http://localhost:3000/users/123/posts?keyword=test")
	}
}

type ResEcho struct {
	Data string `json:"data"`
}

type ParamsEcho struct {
	ID      int    `param:"id"`
	Keyword string `query:"keyword"`
}

func Benchmark_echo(b *testing.B) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log.ERROR)

	e.GET("/users/:id/posts", func(c echo.Context) error {
		p := ParamsEcho{}
		_ = c.Bind(&p)
		return c.JSON(200, ResEcho{Data: fmt.Sprintf("%d %s", p.ID, p.Keyword)})
	})

	go func() { _ = e.Start(":3001") }()
	b.Cleanup(func() { _ = e.Shutdown(context.Background()) })

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req("http://localhost:3001/users/123/posts?keyword=test")
	}
}

func req(u string) {
	res, err := http.Get(u) //nolint: noctx
	if err != nil {
		panic(err)
	}

	d := ResEcho{}

	_ = json.NewDecoder(res.Body).Decode(&d)

	if d.Data != "123 test" {
		panic("invalid response")
	}

	_ = res.Body.Close()
}
