package ws

// Matcher tests a Message
type Matcher interface {
	Match(*Message) bool
}

// Handler handles a Message
type Handler interface {
	ServeWS(*Socket, *Message)
}

// HandlerFunc wraps a func to implement Handler
type HandlerFunc func(*Socket, *Message)

// ServeWS calls the HandlerFunc with m, which provides Handler
func (h HandlerFunc) ServeWS(s *Socket, m *Message) {
	h(s, m)
}

// Route holds pointers to Matcher and Handler
//
// Route provides Router
type Route struct {
	Matcher
	Handler
}

// Router provides Matcher and Handler
type Router interface {
	Matcher
	Handler
}
