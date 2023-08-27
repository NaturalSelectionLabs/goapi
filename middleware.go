package goapi

import "net/http"

type Middleware interface {
	Handle(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

type MiddlewareFunc func(http.ResponseWriter, *http.Request, http.HandlerFunc)

func (m MiddlewareFunc) Handle(w http.ResponseWriter, rq *http.Request, next http.HandlerFunc) {
	m(w, rq, next)
}
