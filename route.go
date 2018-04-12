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
	var err error
	route.try(r, &err)
	return err
}

func (route *route) try(r *Request, err *error) {
	defer func() {
		if e := recover(); e != nil {
			log.WithFields(log.Fields{
				"Error":   e,
				"Quest":   r.Quest,
				"Agent":   r.Agent,
				"Session": r.Session,
			}).Error("route panic")
		}
	}()
	*err = route.ResponderFunc(r)
}
