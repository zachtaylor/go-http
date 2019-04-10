package websocket

import (
	"errors"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/sessions"
)

// ErrSocketKeyExists is returned by Cache.Store when the key is duplicate
var ErrSocketKeyExists = errors.New("socket key exists")

// Cache implements Service
type Cache struct {
	sessions sessions.Service
	cache    map[string]*T
	mux      Mux
}

// cacheIsService is a type check
func cacheIsService(c *Cache) Service {
	return c
}

// NewCache builds a Cache, required for internals
func NewCache(sessions sessions.Service) *Cache {
	return &Cache{
		sessions: sessions,
		cache:    make(map[string]*T),
		mux:      make(Mux, 0),
	}
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
	c.watch(t) // start monitor
}
func (c *Cache) watch(t *T) {
	for {
		msg, err := t.NextMessage()
		if err != nil {
			delete(c.cache, t.Key())
			t.Close()
			return
		}
		go c.ServeWS(t, msg)
	}
}

// Count returns the number of open sockets
func (c *Cache) Count() int {
	return len(c.cache)
}

// Plugin adds a plugin to the mux
func (c *Cache) Plugin(p Plugin) {
	c.mux = append(c.mux, p)
}

// ServeWS calls the mux ServeWS
func (c *Cache) ServeWS(t *T, m *Message) {
	c.mux.ServeWS(t, m)
}
