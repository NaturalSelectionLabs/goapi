package main

import (
	"github.com/NaturalSelectionLabs/goapi"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	router := goapi.New()

	router.GET("/hello", func() Res {
		return Res{Data: "World"}
	})

	e.Use(echo.WrapMiddleware(router.Handler))

	_ = e.Start(":3000")
}

type Res struct {
	goapi.StatusOK
	Data string
}
