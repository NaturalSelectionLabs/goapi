package main

import (
	"log"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/go-playground/validator/v10"
)

type Params struct {
	goapi.InURL
	Email string `validate:"email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func main() {
	validate := validator.New()

	r := goapi.New()
	r.Router().Validate = func(v interface{}) *openapi.Error {
		err := validate.Struct(v)
		if err != nil {
			return &openapi.Error{Code: openapi.CodeInvalidParam, Message: err.Error()}
		}

		return nil
	}

	r.GET("/", func(p Params) goapi.StatusOK {
		return goapi.StatusOK{}
	})

	log.Println(r.Start(":3000"))
}
