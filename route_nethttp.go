package http

import "net/http"

type NetHttpRoute struct {
	Matcher
	http.Handler
}

func (route *NetHttpRoute) Match(r *Request) bool {
	return route.Matcher.Match(r)
}

func (route *NetHttpRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route.Handler.ServeHTTP(w, r)
}

func (route *NetHttpRoute) Respond(r *Request) error {
	return ErrRespondPathRaw
}
