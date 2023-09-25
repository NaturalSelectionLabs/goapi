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

	goapi.Add(g, hello)

	e.Use(func(ctx *gin.Context) {
		g.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Next()
		})).ServeHTTP(ctx.Writer, ctx.Request)
	})

	_ = e.Run(":3000")
}

func hello(any) string {
	return "World"
}
