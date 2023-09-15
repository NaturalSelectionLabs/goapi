package main

import (
	"context"
	"fmt"
	"log"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/goapi/lib/middlewares/apidoc"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

// To test the example, start the server
//
//	go run ./lib/examples/basic
//
// Then open http://127.0.0.1:3000 in your browser.
func main() {
	r := goapi.New()

	r.POST("/login", func(p ParamsLogin) ResLogin {
		// If the username and password are not correct, return a LoginFail response.
		if p.Username != "admin" || p.Password != "123456" {
			return goapi.StatusUnauthorized{}
		}

		// If the username and password are correct, return a LoginOK response.
		return ResLoginOK{
			// goapi will automatically use the standard case conversion,
			// Here SetCookie will be converted to Set-Cookie in http.
			// Same works for url path and query.
			Header: struct{ SetCookie string }{
				SetCookie: "token=123456",
			},
		}
	})

	// You can use multiple parameters at the same time to get url values, headers, request context, or request body.
	// The order of the parameters doesn't matter.
	r.GET("/users/{id}/posts", func(c context.Context, f ParamsPosts, h ParamsHeader) ResPosts {
		if h.Cookie != "token=123456" {
			return goapi.StatusUnauthorized{}
		}

		return ResPostsOK{
			Data: fetchPosts(c, f.ID, f.Type.String(), f.Keyword),
			Meta: 100,
		}
	})

	// Install endpoints for openapi doc.
	apidoc.Install(r, func(doc *openapi.Document) *openapi.Document {
		// Use this callback to customize the openapi document.
		doc.Info.Title = "Basic Example"
		doc.Info.Version = "0.0.1"
		return doc
	})

	log.Println(r.Start(":3000"))
}

// Simulate slow data fetching from database.
func fetchPosts(c context.Context, id int, keyword, typ string) []string {
	posts := []string{}

	for i := 0; i < 2; i++ {
		if c.Err() != nil { // abort if the request is canceled.
			return posts
		}

		posts = append(posts, fmt.Sprintf("user %d posted %s %s %d", id, typ, keyword, i))
	}

	return posts
}
