package httptrack

import (
	"sync"
	"ztaylor.me/log"
)

var lock = sync.Mutex{}
var witness = make(map[string]bool)
var names = make(map[string]string)

func Witness(addr string) {
	lock.Lock()
	defer lock.Unlock()

	if !witness[addr] {
		witness[addr] = true
		log.Add("Addr", addr).Debug("httptrack: new customer")
	}
}

func Match(name string) []string {
	lock.Lock()
	defer lock.Unlock()

	s := make([]string, 0)
	for k, v := range names {
		if name == v {
			s = append(s, k)
		}
	}

	return s
}

func Link(name string, addr string) {
	lock.Lock()
	defer lock.Unlock()
	if names[addr] == "" || names[name] == "" || names[name] != addr {
		names[name] = addr
		names[addr] = name
		log.Add("Addr", addr).Add("Username", name).Debug("httptrack: new ip account pair")
	}
}
