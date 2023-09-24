// Package main .
package main

import (
	"log"

	"github.com/NaturalSelectionLabs/goapi"
)

// To test it:
//
//	curl 'localhost:3000/double' -d 3
func main() {
	r := goapi.New()

	goapi.Add(r, Double)

	log.Println(r.Start(":3000"))
}

// Double for "POST /double" which doubles the input to response.
func Double(num int) int {
	return num * 2
}
