package websocket

import "regexp"

// Router tests a Message
type Router interface {
	Route(*Message) bool
}

// Handler handles a websocket Message
type Handler interface {
	ServeWS(*T, *Message)
}

// Plugin provides Router and Handler
type Plugin interface {
	Router
	Handler
}

// Route is a pointer pair that implements Plugin
type Route struct {
	Router
	Handler
}

// Mux is slice of Plugins
type Mux []Plugin

// ServeWS routes a message
func (mux Mux) ServeWS(t *T, m *Message) {
	for _, p := range mux {
		if p.Route(m) {
			go mux.callPlugin(p, t, m)
			return
		}
	}
}

// callPlugin uses recover() to guard call to ServeWS
func (mux Mux) callPlugin(p Plugin, t *T, m *Message) {
	defer recover()
	p.ServeWS(t, m)
}

// RouterSet is used to combine Routers into a single Router
type RouterSet []Router

// Route checks that all included Routers return true
func (set RouterSet) Route(m *Message) bool {
	for _, router := range set {
		if router == nil || !router.Route(m) {
			return false
		}
	}
	return true
}

// RouterFunc turns a func into a Router
type RouterFunc func(*Message) bool

// Route implements Router by calling the underlying func
func (f RouterFunc) Route(m *Message) bool {
	return f(m)
}

// HandlerFunc turns a func into a Handler
type HandlerFunc func(*T, *Message)

// ServeWS calls the HandlerFunc with m, which provides Handler
func (h HandlerFunc) ServeWS(t *T, m *Message) {
	h(t, m)
}

type routerRegex struct {
	*regexp.Regexp
}

func (rgx *routerRegex) Route(m *Message) bool {
	return rgx.MatchString(m.URI)
}

// RouterRegex creates a regexp match check against Message.Name
func RouterRegex(s string) Router {
	return &routerRegex{regexp.MustCompile(s)}
}

// RouterLit creates a literal match check against Message.Name
type RouterLit string

// Route implements Router by literally matching the Message.URI
func (s RouterLit) Route(m *Message) bool {
	return string(s) == m.URI
}
