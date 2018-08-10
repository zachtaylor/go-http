package http

import "net/http"

type Route interface {
	Matcher
	http.Handler
	Respond(*Request) error
}
