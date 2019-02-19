package sessions

import (
	"net/http"
	"sync"
	"time"
)

// Service deals with sessions, create using Service
type Service struct {
	keygen Keygener
	cache  map[string]*T
	lock   sync.Mutex
}

// Count returns len active Sessions
func (s *Service) Count() int {
	return len(s.cache)
}

// SessionID finds a Session with id
func (s *Service) SessionID(id string) *T {
	return s.cache[id]
}

// ReadRequestCookie returns the session referred by the cookie, if valid
func (s *Service) ReadRequestCookie(r *http.Request) (t *T) {
	if cookie, err := cookieRead(r); err != nil {
	} else if session := s.SessionID(cookie); session == nil {
	} else {
		t = session
	}
	return
}

// NewGrant creates a Grant with name
func (s *Service) NewGrant(name string) (t *T) {
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

func (s *Service) expire(t *T) {
	s.lock.Lock()
	delete(s.cache, t.id)
	s.lock.Unlock()
	t.Revoke()
}

func (s *Service) watch(d time.Duration) {
	tick := time.Minute
	if d < time.Minute {
		tick = time.Second
	}
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
