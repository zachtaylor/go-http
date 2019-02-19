package ws

import "ztaylor.me/log"

// Mux is slice of Plugins
//
// implements Handler
type Mux []Plugin

// ServeWS routes a message
func (mux Mux) ServeWS(socket *Socket, m *Message) {
	for _, route := range mux {
		if route.Route(m) {
			go log.Protect(func() {
				route.ServeWS(socket, m)
			})
			return
		}
	}
}
