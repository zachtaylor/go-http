package ws

import (
	"errors"
	"sync"
)

// ErrSocketKeyExists is returned by service.Store when the key is duplicate
var ErrSocketKeyExists = errors.New("socket key exists")

type Service struct {
	cache map[string]*Socket
	mux   Mux
	sync.Mutex
}

func (service *Service) Count() int {
	return len(service.cache)
}

func (service *Service) Get(key string) *Socket {
	return service.cache[key]
}

func (service *Service) Store(key string, socket *Socket) (err error) {
	service.Lock()
	if service.cache[key] == nil {
		service.cache[key] = socket
	} else {
		err = ErrSocketKeyExists
	}
	service.Unlock()
	return
}

func (service *Service) Remove(key string) {
	service.Lock()
	delete(service.cache, key)
	service.Unlock()
}

func (service *Service) AddRoute(r Route) {
	service.mux = append(service.mux, r)
}

func (service *Service) Dispatch(socket *Socket, msg *Message) {
	for _, plugin := range service.mux {
		if plugin.Route(msg) {
			plugin.ServeWS(socket, msg)
			return
		}
	}
}
