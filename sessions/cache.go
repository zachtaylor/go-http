package sessions

import (
	"net/http"
	"sync"
	"time"

	"ztaylor.me/keygen"
)

// NewCache creates a sessions caching Service with expiry goroutine
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
	lock   sync.Mutex
}

// Count returns number of active Sessions
func (s *Cache) Count() int {
	return len(s.cache)
}

// Get returns a Session by id, if any
func (s *Cache) Get(id string) *T {
	return s.cache[id]
}

// Cookie returns the Session reffered by cookies, if valid
func (s *Cache) Cookie(r *http.Request) (t *T) {
	if cookie, err := cookieRead(r); err != nil {
	} else if session := s.Get(cookie); session == nil {
	} else {
		t = session
	}
	return
}

// Grant returns a new Session granted to the username
func (s *Cache) Grant(name string) (t *T) {
	t = New(name)
	s.lock.Lock()
	t.id = s.keygen.Keygen()
	for s.cache[t.id] != nil {
		t.id = s.keygen.Keygen()
	}
	s.cache[t.id] = t
	s.lock.Unlock()
	return
}

// watch monitors the Cache forever
func (s *Cache) watch(d time.Duration) {
	tick := time.Minute
	for now := range time.Tick(tick) {
		s.lock.Lock()
		for _, t := range s.cache {
			if t.done == nil || now.Sub(t.Time()) > d {
				go s.expire(t)
			}
		}
		s.lock.Unlock()
	}
}

func (s *Cache) expire(t *T) {
	s.lock.Lock()
	delete(s.cache, t.id)
	s.lock.Unlock()
	t.Revoke()
}
