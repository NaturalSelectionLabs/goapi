// Package main ...
package main

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/goapi/lib/middlewares/apidoc"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

// This example shows 2 simple endpoints for typical blog website, login and fetch posts.
// To test the example, start the server
//
//	go run ./lib/examples/basic
//
// Then open http://127.0.0.1:3000 in your browser.
// You can also run the test curl command in cli:
//
//	bash ./lib/examples/basic/test.sh
func main() {
	g := goapi.New()

	g.Router().AddFormatChecker("password", passwordChecker{})

	g.POST("/login", func(p ParamsLogin) ResLogin {
		// If the username and password are not correct, return a LoginFail response.
		if p.Username != "a@a.com" || p.Password != "123456" {
			return goapi.StatusUnauthorized{}
		}

		// If the username and password are correct, return a LoginOK response.
		return ResLoginOK{
			// goapi will automatically use the standard case conversion,
			// Here SetCookie will be converted to Set-Cookie in http.
			// Same works for url path and query.
			Header: ResLoginHeader{
				SetCookie: "token=123456",
			},
		}
	})

	// You can use multiple parameters at the same time to get url values, headers, request context, or request body.
	// The order of the parameters doesn't matter.
	g.GET("/users/{id}/posts", GetPosts{})

	// Install endpoints for openapi doc.
	apidoc.Install(g, func(doc *openapi.Document) *openapi.Document {
		// Use this callback to customize the openapi document.
		doc.Info.Title = "Basic Example"
		doc.Info.Version = "0.0.1"
		return doc
	})

	log.Println(g.Start(":3000"))
}

// GetPosts is the handler for fetching posts of a user.
type GetPosts struct{}

// Handle implements [goapi.OperationHandler] which let us to handle the request.
func (GetPosts) Handle(c context.Context, f ParamsPosts, h ParamsHeader) ResPosts {
	if h.Cookie != "token=123456" {
		return goapi.StatusUnauthorized{}
	}

	return ResPostsOK{
		Data: fetchPosts(c, f.ID, f.Type.String(), f.Keyword),
		Meta: 100,
	}
}

// OpenAPI implements [goapi.OperationOpenAPI] which let us to customize the generated openapi document
// for the current handler.
func (GetPosts) OpenAPI(doc openapi.Operation) openapi.Operation {
	doc.OperationID = "GetPosts"
	doc.Description = "Fetch posts of a user."
	doc.Tags = []string{"posts"}

	return doc
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

type passwordChecker struct {
}

func (passwordChecker) IsFormat(input interface{}) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(input.(string))
}
