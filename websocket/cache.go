package websocket

import (
	"time"

	"golang.org/x/net/websocket"
	"ztaylor.me/cast/charset"
	"ztaylor.me/http/session"
	"ztaylor.me/keygen"
)

// Cache implements Service
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
//
// Adds 3 hooks
//
// 1) Server Request Message.URI = "/connect"
//
// 2) Server Request Message.URI = "/disconnect"
//
// 3) Client Request Message.URI = "/ping"
//
func (c *Cache) Connect(conn *websocket.Conn) {
	t := New(conn)
	if c.sessions == nil {
	} else if s := c.sessions.Cookie(conn.Request()); s != nil {
		t.Session = s
	}
	for t.ID == "" || c.cache[t.ID] != nil {
		t.ID = keygen.New(16, charset.AlphaCapitalNumeric, keygen.DefaultSettings.Rand)
	}
	c.cache[t.ID] = t
	c.ServeWS(t, NewMessage("/connect", nil))
	c.watch(t)
	c.ServeWS(t, NewMessage("/disconnect", nil))
}

// watch monitors *T, and sends "/ping" when it gets lonely
func (c *Cache) watch(t *T) {
	pingTimeout := time.Minute
	for next, pingTimer := t.NextChan(), time.NewTimer(pingTimeout); ; {
		select {
		case <-pingTimer.C:
			if t.conn != nil {
				return
			}
			t.Message("/ping", nil)
		case msg := <-next:
			if msg == nil {
				pingTimer.Stop()
				return
			}
			pingTimer.Reset(pingTimeout)
			go c.ServeWS(t, msg)
		}
	}
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
