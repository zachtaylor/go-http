package session

import (
	"net/http"
	"sync"
	"time"

	"ztaylor.me/keygen"
)

// Cache implements Service in-memory
type Cache struct {
	life   time.Duration
	keygen keygen.Keygener
	cache  map[string]*T
	lock   sync.Mutex // guards cache write
}

// NewCache creates a sessions caching Service and an expiry-monitor goroutine
func NewCache(lifetime time.Duration) Service {
	c := &Cache{
		life:   lifetime,
		keygen: keygen.DefaultSettings,
		cache:  make(map[string]*T),
	}
	return c
}

// Cookie returns the Session reffered by cookies, if valid
func (c *Cache) Cookie(r *http.Request) (t *T, err error) {
	if cookie, _ := cookieRead(r); len(cookie) < 1 {
		err = ErrCookieNotFound
	} else if session := c.Get(cookie); session == nil {
		err = ErrCookieInvalid
	} else {
		t = session
	}
	return
}

// Count returns number of active Sessions
func (c *Cache) Count() int {
	return len(c.cache)
}

// Get returns a Session by id, if any
func (c *Cache) Get(id string) *T {
	return c.cache[id]
}

// Find returns a Session by name, if any
func (c *Cache) Find(name string) (t *T) {
	c.lock.Lock() // guards cache write
	for _, s := range c.cache {
		if name == s.name {
			t = s
			break
		}
	}
	c.lock.Unlock()
	return
}

// Grant returns a new Session granted to the username
func (c *Cache) Grant(name string) (t *T) {
	var id string
	c.lock.Lock() // guards cache write
	for {
		id = c.keygen.Keygen()
		if c.cache[id] == nil {
			break
		}
	}
	t = New(id, name, c.life)
	c.cache[id] = t
	c.lock.Unlock()
	return
}

// Remove removes a Session from the Cache, and closes the Session
func (c *Cache) Remove(t *T) {
	c.lock.Lock() // guards cache write
	delete(c.cache, t.id)
	c.lock.Unlock()
	t.Close()
}
