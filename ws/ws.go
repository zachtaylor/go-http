package ws // import "ztaylor.me/http/ws"

// Handler handles a websocket Message
type Handler interface {
	ServeWS(*Socket, *Message)
}

// HandlerFunc wraps a func to implement Handler
type HandlerFunc func(*Socket, *Message)

// Manager is a generic Socket container
type Manager interface {
	Count() int
	Get(string) *Socket
	Store(string, *Socket) error
	Remove(string)
	AddRoute(Route)
	Dispatch(*Socket, *Message)
}

// Matcher tests a Message
type Matcher interface {
	Match(*Message) bool
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

// ServeWS calls the HandlerFunc with m, which provides Handler
func (h HandlerFunc) ServeWS(s *Socket, m *Message) {
	h(s, m)
}
