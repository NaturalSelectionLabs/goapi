// Package main .
package main

import (
	"github.com/NaturalSelectionLabs/goapi"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	g := goapi.New()

	g.GET("/hello", func() Res {
		return Res{Data: "World"}
	})

	e.Use(echo.WrapMiddleware(g.Handler))

	_ = e.Start(":3000")
}

// Res .
type Res struct {
	goapi.StatusOK
	Data string
}
