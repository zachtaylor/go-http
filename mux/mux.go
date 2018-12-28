package mux // import "ztaylor.me/http/mux"

import (
	"net/http"

	"ztaylor.me/log"
)

// Mux is set of Routers
//
// provides http.Handler
// safely invokes each router, using Router.Matcher in the order added to choose Handler
type Mux struct {
	routers []Router
}

// NewMux creates a new Server instance
func NewMux() *Mux {
	return &Mux{
		routers: make([]Router, 0),
	}
}

// AddRouter appends a router to the Server
func (mux *Mux) AddRouter(r Router) {
	mux.routers = append(mux.routers, r)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var router Router
	for _, route := range mux.routers {
		if route.Match(r) {
			router = route
			break
		}
	}
	log.Protect(func() {
		router.ServeHTTP(w, r)
	})
}

// Map is shorthand for AddRoute
func (mux *Mux) Map(m Matcher, h http.Handler) {
	mux.AddRouter(&Route{
		Matcher: m,
		Handler: h,
	})
}

// MapLit is shorthand for Map with MatcherLit
func (mux *Mux) MapLit(path string, h http.Handler) {
	mux.Map(MatcherLit(path), h)
}

// MapRgx is shorthand for Map with MatcherRegex
func (mux *Mux) MapRgx(path string, h http.Handler) {
	mux.Map(MatcherRegex(path), h)
}

// ListenAndServe starts this Server
func (mux *Mux) ListenAndServe(port string) error {
	return http.ListenAndServe(port, mux)
}

// ListenAndServeTLS starts this Server
func (mux *Mux) ListenAndServeTLS(cert string, key string) error {
	return http.ListenAndServeTLS(":443", cert, key, mux)
}
