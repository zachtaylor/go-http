package mux // import "ztaylor.me/http/mux"

import "net/http"

// Mux is slice of Plugin
type Mux []Plugin

// Plugin appends a Plugin to this Mux
func (mux *Mux) Plugin(p Plugin) {
	*mux = append(*mux, p)
}

// ServeHTTP satisfies http.Handler by calling each member Plugin
func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var router Plugin
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

// Router sorts Requests
type Router interface {
	Route(*http.Request) bool
}

// Routers is used to combine Routers into a single Router
type Routers []Router

// Add is shorthand for Plugin
func (mux *Mux) Add(m Router, h http.Handler) {
	mux.Plugin(&Route{
		Handler: h,
		Router:  m,
	})
}

// Route implements Plugin; points to Router, http.Handler
type Route struct {
	Router
	http.Handler
}

// Plugin is a union type of Router, http.Handler
type Plugin interface {
	Router
	http.Handler
}
