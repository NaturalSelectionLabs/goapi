package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/NaturalSelectionLabs/goapi"
)

// To test the example, start the server
//
//	go run ./lib/examples/basic
//
// Then run
//
//	bash ./lib/examples/basic/test.sh
func main() {
	r := goapi.New()

	r.POST("/login", func(p LoginParams) Login {
		// If the username and password are not correct, return a LoginFail response.
		if p.Username != "admin" || p.Password != "123456" {
			return Unauthorized{}
		}

		// If the username and password are correct, return a LoginOK response.
		return LoginOK{
			// goapi will automatically use the standard case conversion,
			// Here SetCookie will be converted to Set-Cookie in http.
			// Same works for url path and query.
			Header: struct{ SetCookie string }{
				SetCookie: "token=123456",
			},
		}
	}, goapi.Description("Login with username and password.")) // openapi description for the endpoint.

	// You can use multiple parameters at the same time to get url values, headers, request context, or request body.
	// The order of the parameters doesn't matter.
	r.GET("/users/{id}/posts", func(c context.Context, f PostsParams, h Header) Posts {
		if h.Cookie != "token=123456" {
			return Unauthorized{}
		}

		return PostsOK{
			Data: fetchPosts(c, f.ID, f.Type.String(), f.Keyword),
			Meta: 100,
		}
	})

	// You can use func(http.ResponseWriter, *http.Request) to override the default handler behavior.
	// Here we use it to return the openapi doc.
	r.GET("/openapi.json", func(w http.ResponseWriter, rq *http.Request) {
		doc := r.OpenAPI(nil)
		doc.Info.Title = "Basic Example"
		doc.Info.Version = "0.0.1"
		_, _ = w.Write([]byte(doc.JSON()))
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

type PostsParams struct {
	goapi.InURL
	// Use description tag to describe the openapi parameter.
	ID int `description:"User ID"`
	// Type of the posts to fetch.
	// You can use json tag to override the default parameter naming behavior.
	Type PostType `json:"t"`
	// Use default tag to mark this field as optional,
	// you can also use pointer to mark it as optional.
	// The default value should be a json string.
	Keyword string `default:"\"go\""`
}

// Type of a post.
// When using enumer with -json and -values flags, the generated openapi will respect it.
//
//go:generate go run github.com/dmarkham/enumer@latest -type=PostType -trimprefix=PostType -transform=lower -json -values
type PostType int

const (
	PostTypeAll PostType = iota
	PostTypeGame
	PostTypeMusic
)

type LoginParams struct {
	goapi.InBody
	Username string
	Password string
}

type Login interface {
	goapi.Response
}

// Creates a set to store all the implementations of the Login interface.
var iLogin = goapi.Interface(new(Login))

type LoginOK struct {
	goapi.StatusNoContent
	Header struct {
		SetCookie string
	}
}

var _ = iLogin.Add(LoginOK{})

func (LoginOK) Description() string {
	return "Login successfully." // openapi description for the response.
}

type Posts interface {
	goapi.Response
}

var iPosts = goapi.Interface(new(Posts))

type PostsOK struct {
	goapi.StatusOK
	// Use Data to store the main response data.
	Data []string
	// Use Meta to store info like pagination.
	// Here we use it to store the total number of posts.
	Meta int
}

var _ = iPosts.Add(PostsOK{})

type Header struct {
	goapi.InHeader
	Cookie string
}

type Unauthorized struct {
	goapi.StatusUnauthorized
}

var _ = iLogin.Add(Unauthorized{})
var _ = iPosts.Add(Unauthorized{})
