package main

import (
	"github.com/NaturalSelectionLabs/goapi"
)

// ParamsHeader is the parameters for fetching posts.
type ParamsHeader struct {
	// Use goapi.InHeader to get headers values into rest of the fields.
	goapi.InHeader

	Cookie string
}

// ParamsPosts is the parameters for fetching posts.
type ParamsPosts struct {
	// Use goapi.InURL to get url parameters in path or query into rest of the fields.
	goapi.InURL

	// Use description tag to describe the openapi parameter.
	ID int `description:"User ID" min:"1" examples:"[1, 2]"`
	// Type of the posts to fetch.
	// You can use json tag to override the default parameter naming behavior.
	Type PostType `json:"t"`
	// Use default tag to mark this field as optional,
	// you can also use pointer to mark it as optional.
	// Supported tags: min, max, format, pattern.
	// You can also use [goapi.Router.AddFormatChecker] to add custom format checker.
	// Please use json string for default and examples tags if possible, if you forget to quote strings
	// goapi will help to quote it, but it only works for plain string.
	Keyword string `default:"go" minLen:"1" pattern:"^[a-z]+$" examples:"[\"a\", \"b\"]"`
	// Use embedded struct to share common parameters.
	ParamsPagination
}

// ParamsPagination is the parameters for pagination.
type ParamsPagination struct {
	Limit  int `default:"5"`
	Offset int `default:"0"`
}

// ParamsLogin is the parameters for login.
// If we don't embed goapi.InURL or goapi.InHeader to the struct,
// It will be treated as the request body json.
// It should be treated as a common json struct of golang, goapi won't do any special handling for it,
// such as default field tag won't work.
type ParamsLogin struct {
	Username string `json:"username" format:"email"`
	// Here format:"password" is a custom format checker added by [goapi.Router.AddFormatChecker].
	Password string `json:"password" format:"password"`
}

// PostType of a post.
// When using enumer with -json and -values flags, the generated openapi will respect it.
//
//go:generate go run github.com/ysmood/enumer@v0.1.0 -type=PostType -trimprefix=PostType -transform=lower -json -values
type PostType int

const (
	// PostTypeAll is the default value.
	PostTypeAll PostType = iota
	// PostTypeGame .
	PostTypeGame
	// PostTypeMusic .
	PostTypeMusic
)

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
