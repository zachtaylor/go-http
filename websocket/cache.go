package websocket

import (
	"ztaylor.me/cast"
	"ztaylor.me/cast/charset"
	"ztaylor.me/http/session"
	"ztaylor.me/keygen"
)

// Cache implements Service
type Cache struct {
	cache map[string]*T
	lock  cast.Mutex // guard cache
}

// NewCache builds a Cache, required for internals
func NewCache(sessions session.Service) *Cache {
	return &Cache{
		cache: make(map[string]*T),
	}
}

// Count returns the number of open sockets
func (c *Cache) Count() int {
	return len(c.cache)
}

// Get returns the websocket for the given key
func (c *Cache) Get(id string) *T {
	return c.cache[id]
}

// New creates a socket connection
func (c *Cache) New(conn *Conn) *T {
	c.lock.Lock()
	t := New("", conn)
	for t.id == "" || c.cache[t.id] != nil {
		t.id = keygen.New(16, charset.AlphaCapitalNumeric, keygen.DefaultSettings.Rand)
	}
	c.cache[t.id] = t
	c.lock.Unlock()
	return t
}

// Remove deletes a socket ID from the cache
func (c *Cache) Remove(t *T) {
	c.lock.Lock()
	delete(c.cache, t.id)
	c.lock.Unlock()
}

// Keys creates a new slice of the active connected keys
func (c *Cache) Keys() []string {
	c.lock.Lock()
	keys := make([]string, len(c.cache))
	i := 0
	for k := range c.cache {
		keys[i] = k
		i++
	}
	c.lock.Unlock()
	return keys
}
