package http

import (
	"net/http"

	"ztaylor.me/log"
)

type ServiceRoute struct {
	Matcher
	ServiceFunc func(*Request) error
}

func (route *ServiceRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := route.ServiceFunc(RequestFromNet(r, w)); err != nil {
		log.Add("Error", err).Error("http/route_service")
	}
}

func (route *ServiceRoute) Respond(r *Request) error {
	var err error
	route.try(r, &err)
	return err
}

func (route *ServiceRoute) try(r *Request, err *error) {
	defer func() {
		if e := recover(); e != nil {
			log.WithFields(log.Fields{
				"Error":   e,
				"Quest":   r.Quest,
				"Agent":   r.Agent,
				"Session": r.Session,
			}).Error("http/route_service: panic")
		}
	}()
	*err = route.ServiceFunc(r)
}
