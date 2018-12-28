package ws

import (
	"errors"
	"sync"
)

// ErrSocketKeyExists is returned by service.Store when the key is duplicate
var ErrSocketKeyExists = errors.New("socket key exists")

// Service is the global websocket Manager
var Service Manager = &service{
	cache:   make(map[string]*Socket),
	routers: make([]Router, 0),
}

type service struct {
	cache   map[string]*Socket
	routers []Router
	sync.Mutex
}

func (service *service) Count() int {
	return len(service.cache)
}

func (service *service) Get(key string) *Socket {
	return service.cache[key]
}

func (service *service) Store(key string, socket *Socket) error {
	var err error
	service.Lock()
	if service.cache[key] == nil {
		service.cache[key] = socket
	} else {
		err = ErrSocketKeyExists
	}
	service.Unlock()
	return err
}

func (service *service) Remove(key string) {
	service.Lock()
	delete(service.cache, key)
	service.Unlock()
}

func (service *service) AddRoute(r Route) {
	service.routers = append(service.routers, r)
}

func (service *service) Dispatch(socket *Socket, msg *Message) {
	for _, router := range service.routers {
		if router.Match(msg) {
			router.ServeWS(socket, msg)
			return
		}
	}
}
