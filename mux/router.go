package mux

import "net/http"

// Router provides both Matcher and http.Handler
type Router interface {
	Matcher
	http.Handler
}

// Route holds pointers to Matcher, and http.Handler
//
// provides Router
type Route struct {
	Matcher
	http.Handler
}

// Routes is a branching node in a routing tree
//
// provides Router
type Routes struct {
	Matcher
	*Mux
}

// NewRoutes creates a Routes
func NewRoutes(m Matcher) *Routes {
	return &Routes{
		Matcher: m,
		Mux:     NewMux(),
	}
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route.Handler.ServeHTTP(w, r)
}

func (routes *Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes.Mux.ServeHTTP(w, r)
}
