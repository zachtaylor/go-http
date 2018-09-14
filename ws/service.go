package ws

import (
	"sync"
)

type CacheService struct {
	cache   map[string]*Socket
	routers []Router
	sync.Mutex
}

func (service *CacheService) Count() int {
	return len(service.cache)
}

func (service *CacheService) Get(key string) *Socket {
	return service.cache[key]
}

func (service *CacheService) Store(socket *Socket) error {
	var err error
	service.Lock()
	if service.cache[socket.String()] == nil {
		service.cache[socket.String()] = socket
	} else {
		err = ErrKeyExists
	}
	service.Unlock()
	return err
}

func (service *CacheService) Remove(key string) {
	service.Lock()
	delete(service.cache, key)
	service.Unlock()
}

func (service *CacheService) AddRoute(r Route) {
	service.routers = append(service.routers, r)
}

func (service *CacheService) Dispatch(msg *Message) {
	for _, router := range service.routers {
		if router.Match(msg) {
			router.ServeWS(msg)
			return
		}
	}
}

var Service interface {
	Count() int
	Get(string) *Socket
	Store(*Socket) error
	Remove(string)
	AddRoute(Route)
	Dispatch(*Message)
} = &CacheService{
	cache:   make(map[string]*Socket),
	routers: make([]Router, 0),
}
