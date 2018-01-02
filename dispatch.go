package http

import (
	"ztaylor.me/log"
)

var router = make([]Route, 0)

func Dispatch(a Agent, r *Request) {
	for _, route := range router {
		if !route.Match(r.Quest) {
		} else if err := route.Respond(a, r); err != nil {
			log.WithFields(log.Fields{
				"Error":     err,
				"Quest":     r.Quest,
				"Username":  r.Metadata()["Username"],
				"SessionId": r.Metadata()["SessionId"],
			}).Error("http.dispatch handle error")
		}
	}
}
