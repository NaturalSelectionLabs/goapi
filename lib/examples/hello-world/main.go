// Package main .
package main

import (
	"log"

	"github.com/NaturalSelectionLabs/goapi"
)

// Hello is the response for hello world.
type Hello struct {
	goapi.StatusOK // response http status code 200

	Data string
}

func main() {
	g := goapi.New()

	g.GET("/", func() Hello {
		return Hello{Data: "Hello World!"}
	})

	log.Println(g.Start(":3000"))
}
