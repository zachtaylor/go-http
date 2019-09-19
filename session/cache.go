package session

import (
	"net/http"
	"sync"
	"time"

	"ztaylor.me/keygen"
)

// NewCache creates a sessions caching Service and an expiry-monitor goroutine
func NewCache(lifetime time.Duration) Service {
	c := &Cache{
		keygen: keygen.DefaultSettings,
		cache:  make(map[string]*T),
	}
	go c.watch(lifetime)
	return c
}

// Cache implements Service in-memory
type Cache struct {
	keygen keygen.Keygener
	cache  map[string]*T
	lock   sync.Mutex // guards cache write
}

// Cookie returns the Session reffered by cookies, if valid
func (c *Cache) Cookie(r *http.Request) (t *T) {
	if cookie, err := cookieRead(r); err != nil {
		// no cookie found
	} else if session := c.Get(cookie); session == nil {
		// cookie exists but session does not
	} else {
		t = session
	}
	return
}

// Count returns number of active Sessions
func (c *Cache) Count() int {
	return len(c.cache)
}

// Find returns a Session by name, if any
func (c *Cache) Find(name string) (t *T) {
	c.lock.Lock() // guards cache write
	for _, s := range c.cache {
		if name == s.Name() {
			t = s // fill return buffer
			break // go to unlock, return
		}
	}
	c.lock.Unlock() // defer has overhead
	return
}

// Get returns a Session by id, if any
func (c *Cache) Get(id string) *T {
	return c.cache[id]
}

// Grant returns a new Session granted to the username
func (c *Cache) Grant(name string) *T {
	t := New(name)
	c.lock.Lock() // guards cache write
	for t.id = c.keygen.Keygen(); c.cache[t.id] != nil; t.id = c.keygen.Keygen() {
	}
	c.cache[t.id] = t
	c.lock.Unlock() // defer has overhead
	return t
}

// Remove removes a Session from the Cache, and closes the Session
func (c *Cache) Remove(t *T) {
	c.lock.Lock() // guards cache write
	delete(c.cache, t.id)
	c.lock.Unlock() // defer has overhead
	t.Close()
}

// watch monitors the Cache forever
func (c *Cache) watch(d time.Duration) {
	tick := time.Minute
	for now := range time.Tick(tick) {
		c.lock.Lock() // guard cache write
		for _, t := range c.cache {
			if t.done == nil || now.Sub(t.Time()) > d {
				go c.Remove(t) // awaits lock
			}
		}
		c.lock.Unlock() // avoid stack frame and defer overhead
	}
}
