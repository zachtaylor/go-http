package http

import (
	"regexp"
	"ztaylor.me/log"
)

type RegexRoute struct {
	*regexp.Regexp
	Responder
}

func (route *RegexRoute) Match(s string) bool {
	return route.MatchString(s)
}

func (route *RegexRoute) Handle(a Agent, r *Request) {
	if err := route.Respond(a, r); err != nil {
		log.Add("Quest", r.Quest).Add("Agent", r.Agent).Add("Error", err).Error("request failed")
	}
}
