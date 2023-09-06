package main

import (
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/goapi"
)

// $ go run ./lib/examples/basic
// $ curl localhost:3000/login -id '{"username": "admin", "password": "123456"}'
// Set-Cookie: token=123456
//
// $ curl 'localhost:3000/users/3/posts?keyword=sky' -H 'Cookie: token=123456'
// {"data":["post1","post2"],"meta":"User 3 using keyword: sky"}

func main() {
	router := goapi.New()

	router.POST("/login", func(p LoginParams) Login {
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
	})

	// You can use multiple parameters at the same time to get url values, headers, or request body.
	// The order of the parameters doesn't matter.
	router.GET("/users/{id}/posts", func(f PostsParams, h Header) Posts {
		if h.Cookie != "token=123456" {
			return Unauthorized{}
		}

		return PostsOK{
			Data: []string{"post1", "post2"},
			Meta: fmt.Sprintf("User %d using keyword: %s", f.ID, f.Keyword),
		}
	})

	_ = http.ListenAndServe(":3000", router.Server())
}

type PostsParams struct {
	goapi.InURL
	ID      int
	Keyword string
}

type LoginParams struct {
	goapi.InBody
	Username string
	Password string
}

type Login interface {
	goapi.Response
}

// Creates a set to store all the implementations of the Login interface.
var iLogin = goapi.Vary(new(Login))

type LoginOK struct {
	goapi.StatusNoContent
	Header struct {
		SetCookie string
	}
}

var _ = iLogin.Add(LoginOK{})

type Posts interface {
	goapi.Response
}

var iPosts = goapi.Vary(new(Posts))

type PostsOK struct {
	goapi.StatusOK
	Data []string
	Meta string
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
