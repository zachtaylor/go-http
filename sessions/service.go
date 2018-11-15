package sessions

import (
	"sync"
	"time"
)

var Service Provider = &MemoryService{
	cache: make(map[string]*T),
}

type MemoryService struct {
	cache map[string]*T
	sync.Mutex
}

func (service *MemoryService) Count() int {
	return len(service.cache)
}

func (service *MemoryService) SessionID(id string) *T {
	return service.cache[id]
}

func (service *MemoryService) NewGrant(name string) *T {
	t := New(name)
	service.Lock()
	t.id = Keygen()
	for service.cache[t.id] != nil {
		t.id = Keygen()
	}
	service.cache[t.id] = t
	service.Unlock()
	return t
}

func (service *MemoryService) expire(t *T) {
	service.Lock()
	delete(service.cache, t.id)
	service.Unlock()
	t.Revoke()
}

// Watch requires goroutine to monitor session expiry
//
// as in `go service.Watch(time.Duration)` before `server.Start()`
func (service *MemoryService) Watch(d time.Duration) {
	for now := range time.Tick(1 * time.Second) {
		service.Lock()
		for _, t := range service.cache {
			if now.Sub(t.Time()) > d {
				go service.expire(t)
			}
		}
		service.Unlock()
	}
}
