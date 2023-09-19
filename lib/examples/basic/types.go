package main

import (
	"github.com/NaturalSelectionLabs/goapi"
)

// ParamsPosts is the parameters for fetching posts.
type ParamsPosts struct {
	goapi.InURL
	// Use description tag to describe the openapi parameter.
	ID int `description:"User ID" min:"1"`
	// Type of the posts to fetch.
	// You can use json tag to override the default parameter naming behavior.
	Type PostType `json:"t"`
	// Use default tag to mark this field as optional,
	// you can also use pointer to mark it as optional.
	// Supported tags: min, max, format, pattern.
	// You can also use [goapi.Router.AddFormatChecker] to add custom format checker.
	Keyword string `default:"go" min:"1" pattern:"^[a-z]+$"`
	// Use embedded struct to share common parameters.
	ParamsPagination
}

// ParamsPagination is the parameters for pagination.
type ParamsPagination struct {
	Limit  int `default:"5"`
	Offset int `default:"0"`
}

// PostType of a post.
// When using enumer with -json and -values flags, the generated openapi will respect it.
//
//go:generate go run github.com/dmarkham/enumer@latest -type=PostType -trimprefix=PostType -transform=lower -json -values
type PostType int

const (
	// PostTypeAll is the default value.
	PostTypeAll PostType = iota
	// PostTypeGame .
	PostTypeGame
	// PostTypeMusic .
	PostTypeMusic
)

// ParamsLogin is the parameters for login.
type ParamsLogin struct {
	goapi.InBody
	Username string
	Password string
}

// ResLogin is the response for login.
type ResLogin interface {
	goapi.Response
}

// Creates a set to store all the implementations of the Login interface.
var _ = goapi.Interface(new(ResLogin), ResLoginOK{}, goapi.StatusUnauthorized{})

// ResLoginOK is the response for successful login.
type ResLoginOK struct {
	goapi.StatusNoContent
	Header struct {
		SetCookie string
	}
}

// Description implements [goapi.Descriptioner] which let us to customize the description of the response.
func (ResLoginOK) Description() string {
	return "Login successfully."
}

// ResLoginHeader is the header for successful login.
type ResLoginHeader struct {
	SetCookie string
}

// ResPosts is the response for fetching posts.
type ResPosts interface {
	goapi.Response
}

var _ = goapi.Interface(new(ResPosts), ResPostsOK{}, goapi.StatusUnauthorized{})

// ResPostsOK is the response for successful fetching posts.
type ResPostsOK struct {
	goapi.StatusOK
	// Use Data to store the main response data.
	Data []string
	// Use Meta to store info like pagination.
	// Here we use it to store the total number of posts.
	Meta int
}

// ParamsHeader is the parameters for fetching posts.
type ParamsHeader struct {
	goapi.InHeader
	Cookie string
}
