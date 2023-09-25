// Package main .
package main

import (
	"net/http"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()

	g := goapi.New()

	g.GET("/hello", func() Res {
		return Res{Data: "World"}
	})

	e.Use(func(ctx *gin.Context) {
		g.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Next()
		})).ServeHTTP(ctx.Writer, ctx.Request)
	})

	_ = e.Run(":3000")
}

// Res .
type Res struct {
	goapi.StatusOK
	Data string
}
