package http

import (
	"net/http"
	"ztaylor.me/log"
)

var router = make([]Route, 0)

func Router(r Route) {
	router = append(router, r)
}
func MapLit(s string, f func(*Request) error) {
	Router(&route{StringMatcher(s), f})
}
func MapRgx(s string, f func(*Request) error) {
	Router(&route{RegexMatcher(s), f})
}
func MapRawLit(s string, h http.Handler) {
	Router(NewRouteNetHttp(StringMatcher(s), h))
}
func MapRawRgx(s string, h http.Handler) {
	Router(NewRouteNetHttp(RegexMatcher(s), h))
}

func Dispatch(r *Request) {
	for _, route := range router {
		if !route.Match(r.Quest) {
		} else if err := route.Respond(r); err != nil {
			log.WithFields(log.Fields{
				"Error":   err,
				"Quest":   r.Quest,
				"Session": r.Session,
			}).Error("dispatch respond error")
		} else {
			return
		}
	}
}
