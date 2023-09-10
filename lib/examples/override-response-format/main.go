package main

import (
	"log"

	"github.com/NaturalSelectionLabs/goapi"
)

type Res struct {
	goapi.StatusOK
	Data string
}

func main() {
	r := goapi.NewRouter()

	r.FormatResponse = func(format goapi.ResponseFormat) any {
		switch f := format.(type) {
		// Return the data directly without nested "data" field.
		case goapi.ResponseFormatData:
			return f.Data

		// Return the error directly without nested "error" field.
		case goapi.ResponseFormatErr:
			return f.Error

		default:
			panic("unknown format")
		}
	}

	g := r.Group("")

	g.GET("/", func() Res {
		return Res{Data: "ok"}
	})

	log.Println(r.Start(":3000"))
}
