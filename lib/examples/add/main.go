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
	g := goapi.New()

	goapi.Add(g, double)

	log.Println(g.Start(":3000"))
}

// Handler for "POST /double" which doubles the input to response.
func double(num int) int {
	return num * 2
}
