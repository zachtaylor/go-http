package http

import (
	"net/http"
	"ztaylor.me/log"
)

type Route interface {
	Match(string) bool
	http.Handler
	Respond(*Request) error
}

type route struct {
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
