package http

import (
	"net/http"

	"ztaylor.me/log"
)

type Router struct {
	routes []Route
}

func (r *Router) AddRoute(route Route) {
	r.routes = append(r.routes, route)
}

func (r *Router) MapLit(s string, f func(*Request) error) {
	r.AddRoute(&ServiceRoute{QuestMatcher(s), f})
}

func (r *Router) MapRgx(s string, f func(*Request) error) {
	r.AddRoute(&ServiceRoute{RegexMatcher(s), f})
}

func (r *Router) MapRawLit(s string, h http.Handler) {
	r.AddRoute(&NetHttpRoute{QuestMatcher(s), h})
}

func (r *Router) MapRawRgx(s string, h http.Handler) {
	r.AddRoute(&NetHttpRoute{RegexMatcher(s), h})
}

func NewRouter() *Router {
	return &Router{
		routes: make([]Route, 0),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, q *http.Request) {
	req := RequestFromNet(q, w)
	for _, route := range r.routes {
		if route.Match(req) {
			route.ServeHTTP(w, q)
			return
		}
	}
}

func (r *Router) Dispatch(q *Request) {
	for _, route := range r.routes {
		if !route.Match(q) {
			continue
		}
		if err := route.Respond(q); err != nil && err != ErrRespondPathRaw {
			log.WithFields(log.Fields{
				"Error":   err,
				"Quest":   q.Quest,
				"Agent":   q.Agent,
				"Session": q.Session,
			}).Error("http/router: dispatch")
		}
		return
	}
}
