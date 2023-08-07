package goapi

import "net/http"

type Middleware func(http.ResponseWriter, *http.Request, http.HandlerFunc)
