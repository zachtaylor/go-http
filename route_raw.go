package http

import (
	"errors"
	"net/http"
)

var ErrRespondPathRaw = errors.New("http path cannot respond outside http")

type RawRoute struct {
	http.Handler
	Path string
}

func NewRouteRaw(path string, h http.Handler) Route {
	return &RawRoute{
		Handler: h,
		Path:    path,
	}
}

func (route *RawRoute) Match(s string) bool {
	return route.Path == s
}

func (route *RawRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route.Handler.ServeHTTP(w, r)
}

func (route *RawRoute) Respond(r *Request) error {
	return ErrRespondPathRaw
}
