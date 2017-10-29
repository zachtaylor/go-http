package httptrack

import (
	// "net/http"
	"sync"
	"ztaylor.me/log"
)

var lock = sync.Mutex{}

var mem = make(map[string]bool)

func Witness(addr string) {
	lock.Lock()
	defer lock.Unlock()

	if !mem[addr] {
		mem[addr] = true
		log.Add("Addr", addr).Debug("httptrack: new customer")
	}
}

var names = make(map[string]string)

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
	names[name] = addr
	names[addr] = name

}
