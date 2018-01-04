package http

import (
	"net/http"
	"ztaylor.me/log"
)

var router = make([]Route, 0)

func Router(r Route) {
	router = append(router, r)
}
func MapRoute(s string, f ResponderFunc) {
	Router(NewRouteLiteral(s, f))
}
func MapRegex(s string, f ResponderFunc) {
	Router(NewRouteRegex(s, f))
}
func MapRaw(s string, h http.Handler) {
	Router(NewRouteRaw(s, h))
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
		}
	}
}
