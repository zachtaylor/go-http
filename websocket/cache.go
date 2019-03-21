package websocket

import (
	"errors"
	"sync"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/sessions"
)

// ErrSocketKeyExists is returned by Cache.Store when the key is duplicate
var ErrSocketKeyExists = errors.New("socket key exists")

// Cache manages the websockets
type Cache struct {
	sessions sessions.Service
	cache    map[string]*T
	mux      Mux
	sync.Mutex
}

// NewCache creates a Cache with internals initialized
func NewCache(sessions sessions.Service) *Cache {
	return &Cache{
		sessions: sessions,
		cache:    make(map[string]*T),
		mux:      make(Mux, 0),
	}
}

// New creates a Socket
func (cache *Cache) New(conn *websocket.Conn) *T {
	t := New(conn)
	if cache.sessions != nil {
		if s := cache.sessions.Cookie(conn.Request()); s != nil {
			t.Session = s
		}
		cache.Store(t.String(), t)
	}
	return t
}

// Count returns the number of open sockets
func (cache *Cache) Count() int {
	return len(cache.cache)
}

// Get returns a socket with the given id, if available
func (cache *Cache) Get(key string) *T {
	return cache.cache[key]
}

// Store saves a socket to the cache
func (cache *Cache) Store(key string, socket *T) (err error) {
	cache.Lock()
	if cache.cache[key] == nil {
		cache.cache[key] = socket
	} else {
		err = ErrSocketKeyExists
	}
	cache.Unlock()
	return
}

// Remove deletes a socket from the cache
func (cache *Cache) Remove(key string) {
	cache.Lock()
	delete(cache.cache, key)
	cache.Unlock()
}

// Plugin adds a plugin to the mux
func (cache *Cache) Plugin(p Plugin) {
	cache.mux = append(cache.mux, p)
}

// ServeWS calls the mux ServeWS
func (cache *Cache) ServeWS(t *T, m *Message) {
	cache.mux.ServeWS(t, m)
}
