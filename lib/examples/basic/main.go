package main

import (
	"fmt"
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

	// You can use func(http.ResponseWriter, *http.Request) to override the default handler behavior.
	// Here we use it to return the openapi doc.
	router.GET("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		doc := router.OpenAPI(nil)
		doc.Info.Title = "Basic Example"
		doc.Info.Version = "0.0.1"
		_, _ = w.Write([]byte(doc.JSON()))
	})

	_ = http.ListenAndServe(":3000", router.Server())
}

type PostsParams struct {
	goapi.InURL
	ID int
	// Use default tag to mark this field as optional,
	// you can also use pointer to mark it as optional.
	// The default value should be a json string.
	Keyword string `default:"\"go\""`
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
var iLogin = goapi.Interface(new(Login))

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

var iPosts = goapi.Interface(new(Posts))

type PostsOK struct {
	goapi.StatusOK
	// Use Data to store the main response data.
	Data []string
	// Use Meta to store info like pagination.
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
