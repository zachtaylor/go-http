package http

import (
	"errors"
	"net/http"
)

var ErrRespondPathRaw = errors.New("http path cannot respond outside http")

type NetHttpRoute struct {
	http.Handler
	Matcher
}

func NewRouteNetHttp(m Matcher, h http.Handler) Route {
	return &NetHttpRoute{
		Handler: h,
		Matcher: m,
	}
}

func (route *NetHttpRoute) Match(s string) bool {
	return route.Matcher.Match(s)
}

func (route *NetHttpRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route.Handler.ServeHTTP(w, r)
}

func (route *NetHttpRoute) Respond(r *Request) error {
	return ErrRespondPathRaw
}
