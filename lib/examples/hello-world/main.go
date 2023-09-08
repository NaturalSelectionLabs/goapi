package main

import (
	"log"

	"github.com/NaturalSelectionLabs/goapi"
)

type Hello struct {
	goapi.StatusOK // response http status code 200

	Data string
}

func main() {
	r := goapi.New()

	r.GET("/", func() Hello {
		return Hello{Data: "Hello World!"}
	})

	log.Println(r.Start(":3000"))
}
