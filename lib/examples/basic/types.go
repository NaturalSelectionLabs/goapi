package main

import "github.com/NaturalSelectionLabs/goapi"

type ParamsPosts struct {
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
	// Use embedded struct to share common parameters.
	ParamsPagination
}

type ParamsPagination struct {
	Limit  int `default:"5"`
	Offset int `default:"0"`
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

type ParamsLogin struct {
	goapi.InBody
	Username string
	Password string
}

type ResLogin interface {
	goapi.Response
}

// Creates a set to store all the implementations of the Login interface.
var _ = goapi.Interface(new(ResLogin), ResLoginOK{}, goapi.StatusUnauthorized{})

type ResLoginOK struct {
	goapi.StatusNoContent
	Header struct {
		SetCookie string
	}
}

func (ResLoginOK) Description() string {
	return "Login successfully." // openapi description for the response.
}

type ResPosts interface {
	goapi.Response
}

var _ = goapi.Interface(new(ResPosts), ResPostsOK{}, goapi.StatusUnauthorized{})

type ResPostsOK struct {
	goapi.StatusOK
	// Use Data to store the main response data.
	Data []string
	// Use Meta to store info like pagination.
	// Here we use it to store the total number of posts.
	Meta int
}

type ParamsHeader struct {
	goapi.InHeader
	Cookie string
}
