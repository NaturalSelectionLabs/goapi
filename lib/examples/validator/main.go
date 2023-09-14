package main

import (
	"log"

	"github.com/NaturalSelectionLabs/goapi"
)

type Params struct {
	goapi.InURL
	Email string `validate:"email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func main() {
	// validate := validator.New()

	r := goapi.New()
	// r.Router().Validate = func(v interface{}) *openapi.DefaultError {
	// 	err := validate.Struct(v)
	// 	if err != nil {
	// 		return &openapi.DefaultError{Code: openapi.ErrCodeInvalidParam, Message: err.Error()}
	// 	}

	// 	return nil
	// }

	r.GET("/", func(p Params) goapi.StatusOK {
		return goapi.StatusOK{}
	})

	log.Println(r.Start(":3000"))
}
