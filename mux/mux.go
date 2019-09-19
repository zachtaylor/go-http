package mux // import "ztaylor.me/http/mux"

import "net/http"

// Mux is slice of *Route
type Mux []*Route

// Route contains routing, handling
type Route struct {
	Router
	http.Handler
}

// Router sorts Requests
type Router interface {
	// Route returns whether the router accepts this request
	Route(*http.Request) bool
}

// Add is shorthand for Plugin
func (mux *Mux) Add(m Router, h http.Handler) {
	mux.Route(&Route{
		Handler: h,
		Router:  m,
	})
}

// Route appends a Route to this Mux
func (mux *Mux) Route(r *Route) {
	*mux = append(*mux, r)
}

// ServeHTTP satisfies http.Handler by calling each member Plugin
func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var router *Route
	for _, route := range mux {
		if route.Route(r) {
			router = route
			break
		}
	}
	if router != nil {
		router.ServeHTTP(w, r)
	}
}
