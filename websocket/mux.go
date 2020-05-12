package websocket

// Router tests a Message
type Router interface {
	Route(*Message) bool
}

// Handler handles a websocket Message
type Handler interface {
	ServeWS(*T, *Message)
}

// Route points to a pair of routing and handling logic
type Route struct {
	Router
	Handler
}

// Mux is slice of Route
type Mux []*Route

// Add is shorthand for Plugin
func (mux *Mux) Add(m Router, h Handler) {
	mux.Route(&Route{
		Handler: h,
		Router:  m,
	})
}

// Route appends a Route to this Mux
func (mux *Mux) Route(r *Route) {
	*mux = append(*mux, r)
}

// ServeWS routes a message
func (mux Mux) ServeWS(t *T, m *Message) {
	for _, r := range mux {
		if r.Route(m) {
			go mux.callRoute(r, t, m)
			return
		}
	}
}

// callRoute uses recover() to guard call to ServeWS
func (mux Mux) callRoute(r *Route, t *T, m *Message) {
	defer recover()
	r.ServeWS(t, m)
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

// RouterLit creates a literal match check against Message.Name
type RouterLit string

// Route implements Router by literally matching the Message.URI
func (s RouterLit) Route(m *Message) bool {
	return string(s) == m.URI
}
