package http

import (
	"net/http"
	"ztaylor.me/log"
)

type Route interface {
	Matcher
	http.Handler
	Respond(*Request) error
}

type route struct {
	Matcher
	ResponderFunc func(*Request) error
}

func (route *route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := route.ResponderFunc(RequestFromNet(r, w)); err != nil {
		log.Error(err)
	}
}

func (route *route) Respond(r *Request) error {
	return route.ResponderFunc(r)
}
