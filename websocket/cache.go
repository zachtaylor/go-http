package websocket

import (
	"ztaylor.me/cast"
	"ztaylor.me/cast/charset"
	"ztaylor.me/http/session"
	"ztaylor.me/keygen"
)

// heatline tunes what cache considers excessive input (starts at 0)
var heatline = 2

// Cache implements Service
type Cache struct {
	sessions session.Service
	cache    map[string]*T
	mux      Mux
	lock     cast.Mutex // guard cache
}

// NewCache builds a Cache, required for internals
func NewCache(sessions session.Service) Service {
	return &Cache{
		sessions: sessions,
		cache:    make(map[string]*T),
		mux:      make(Mux, 0),
	}
}

// Connect connects a websocket
//
// Adds 3 hooks
//
// 1) Server Request Message.URI = "/connect"
//
// 2) Server Request Message.URI = "/disconnect"
//
// 3) Client Request Message.URI = "/ping"
//
func (c *Cache) Connect(conn *Conn) {
	t := New(conn)
	if c.sessions == nil {
	} else if s, _ := c.sessions.Cookie(conn.Request()); s != nil {
		t.Session = s
	}
	c.lock.Lock()
	for t.ID == "" || c.cache[t.ID] != nil {
		t.ID = keygen.New(16, charset.AlphaCapitalNumeric, keygen.DefaultSettings.Rand)
	}
	c.cache[t.ID] = t
	c.lock.Unlock()
	c.ServeWS(t, NewMessage("/connect", nil))
	watch(c, t)
	c.ServeWS(t, NewMessage("/disconnect", nil))
	c.lock.Lock()
	delete(c.cache, t.ID)
	c.lock.Unlock()
}

// Message is a macro for SendMessage(NewMessage)
func (c *Cache) Message(uri string, json cast.JSON) {
	c.SendMessage(NewMessage(uri, json))
}

// SendMessage calls Write with cast []byte Message.JSON().String()
func (c *Cache) SendMessage(m *Message) {
	c.Send(cast.BytesS(m.JSON().String()))
}

// Send sends a buffer all websocket connections
func (c *Cache) Send(s []byte) {
	c.lock.Lock()
	for _, socket := range c.cache {
		socket.Send(s)
	}
	c.lock.Unlock()
}

// Count returns the number of open sockets
func (c *Cache) Count() int {
	return len(c.cache)
}

// Add alloc new route, call Route
func (c *Cache) Add(router Router, handler Handler) {
	c.Route(&Route{router, handler})
}

// Route adds a route to the mux
func (c *Cache) Route(r *Route) {
	c.mux = append(c.mux, r)
}

// ServeWS calls the mux ServeWS
func (c *Cache) ServeWS(t *T, m *Message) {
	c.mux.ServeWS(t, m)
}
