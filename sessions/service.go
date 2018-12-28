package sessions

import (
	"sync"
	"time"

	"ztaylor.me/keygen"
)

// Manager is a generic sessions container
type Manager interface {
	// Count returns len active Sessions
	Count() int
	// Get finds a Session with id
	SessionID(id string) *T
	// NewGrant creates a Grant with name
	NewGrant(name string) *T
	// Watch needs a goroutine to expire Sessions older than d
	Watch(d time.Duration)
}

// Service is the global sessions Manager
var Service Manager = &service{
	Keygener: keygen.DefaultSettings,
	cache:    make(map[string]*T),
}

type service struct {
	keygen.Keygener
	cache map[string]*T
	sync.Mutex
}

func (service *service) Count() int {
	return len(service.cache)
}

func (service *service) SessionID(id string) *T {
	return service.cache[id]
}

func (service *service) NewGrant(name string) *T {
	t := New(name)
	service.Lock()
	t.id = service.Keygen()
	for service.cache[t.id] != nil {
		t.id = service.Keygen()
	}
	service.cache[t.id] = t
	service.Unlock()
	return t
}

func (service *service) expire(t *T) {
	service.Lock()
	delete(service.cache, t.id)
	service.Unlock()
	t.Revoke()
}

// Watch requires goroutine to monitor session expiry. Defines session lifetime.
//
// as in `go service.Watch(time.Duration)` before `server.Start()`
func (service *service) Watch(d time.Duration) {
	tick := time.Minute
	if d < time.Minute {
		tick = time.Second
	}
	for now := range time.Tick(tick) {
		service.Lock()
		for _, t := range service.cache {
			if t.done == nil || now.Sub(t.Time()) > d {
				go service.expire(t)
			}
		}
		service.Unlock()
	}
}
