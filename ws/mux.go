package ws

import "ztaylor.me/log"

// Mux is set of Routers
//
// implements Handler
type Mux struct {
	routers []Router
}

// NewMux creates a new Mux
func NewMux() *Mux {
	return &Mux{
		routers: make([]Router, 0),
	}
}

// AddRouter appends a router to the Mux
func (mux *Mux) AddRouter(r Router) {
	mux.routers = append(mux.routers, r)
}

// ServeWS routes a message
func (mux *Mux) ServeWS(socket *Socket, m *Message) {
	for _, route := range mux.routers {
		if route.Match(m) {
			go log.Protect(func() {
				route.ServeWS(socket, m)
			})
			return
		}
	}
}

// Map is shorthand for AddRoute
func (mux *Mux) Map(m Matcher, h Handler) {
	mux.AddRouter(&Route{
		Matcher: m,
		Handler: h,
	})
}

// MapLit is shorthand for Map with MatcherLit
func (mux *Mux) MapLit(path string, h Handler) {
	mux.Map(MatcherLit(path), h)
}

// MapRgx is shorthand for Map with MatcherRegex
func (mux *Mux) MapRgx(path string, h Handler) {
	mux.Map(MatcherRegex(path), h)
}
