package sessions

import (
	"sync"
)

type MemoryService struct {
	cache map[string]*Grant
	sync.Mutex
}

func (service *MemoryService) Count() int {
	return len(service.cache)
}

func (service *MemoryService) Get(key string) *Grant {
	return service.cache[key]
}

func (service *MemoryService) Grant(name string) *Grant {
	grant := NewGrant(name)
	service.Lock()
	for service.cache[grant.ID] != nil {
		grant.ID = Keygen()
	}
	service.cache[grant.ID] = grant
	service.Unlock()
	return grant
}

func (service *MemoryService) Revoke(key string) {
	service.Lock()
	delete(service.cache, key)
	service.Unlock()
}

var Service interface {
	Count() int
	Get(string) *Grant
	Grant(string) *Grant
	Revoke(string)
} = &MemoryService{
	cache: make(map[string]*Grant),
}
