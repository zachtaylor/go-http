package ws // import "ztaylor.me/http/ws"

// Handler handles a websocket Message
type Handler interface {
	ServeWS(*Socket, *Message)
}

// HandlerFunc wraps a func to implement Handler
type HandlerFunc func(*Socket, *Message)

// ServeWS calls the HandlerFunc with m, which provides Handler
func (h HandlerFunc) ServeWS(s *Socket, m *Message) {
	h(s, m)
}

// Router tests a Message
type Router interface {
	Route(*Message) bool
}

// Route holds pointers to Router and Handler
type Route struct {
	Router
	Handler
}

// Plugin provides Router and Handler
type Plugin interface {
	Router
	Handler
}
