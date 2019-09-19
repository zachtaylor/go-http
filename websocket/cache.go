package websocket

import (
	"errors"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/session"
)

// ErrSocketKeyExists is returned by Cache.Store when the key is duplicate
var ErrSocketKeyExists = errors.New("socket key exists")

// Cache implements Service
//
// Adds 2 hooks for connection
// 1) Message.URI = "/connect"
// 2) Message.URI = "/disconnect"
type Cache struct {
	sessions session.Service
	cache    map[string]*T
	mux      Mux
}

// NewCache builds a Cache, required for internals
func NewCache(sessions session.Service) *Cache {
	return &Cache{
		sessions: sessions,
		cache:    make(map[string]*T),
		mux:      make(Mux, 0),
	}
}
func _cacheIsService(c *Cache) Service {
	return c
}

// Connect connects a websocket
func (c *Cache) Connect(conn *websocket.Conn) {
	t := New(conn)
	if c.sessions != nil {
		if s := c.sessions.Cookie(conn.Request()); s != nil {
			t.Session = s
		}
	}
	c.cache[t.Key()] = t
	c.ServeWS(t, &Message{
		URI:  "/connect",
		User: t.GetUser(),
	})
	for msg := range t.NextChan() {
		c.ServeWS(t, msg)
	}
	c.ServeWS(t, &Message{
		URI:  "/disconnect",
		User: t.GetUser(),
	})
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
