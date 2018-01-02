package http

import (
	"ztaylor.me/log"
)

type LiteralRoute struct {
	Path string
	Responder
}

func (route *LiteralRoute) Match(s string) bool {
	return route.Path == s
}

func (route *LiteralRoute) Handle(a Agent, r *Request) {
	if err := route.Respond(a, r); err != nil {
		log.Add("Quest", r.Quest).Add("Error", err).Error("request failed")
	}
}
