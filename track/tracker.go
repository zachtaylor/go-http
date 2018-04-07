package httptrack

import (
	"time"
	"ztaylor.me/log"
)

func Match(name string, addr string) bool {
	if Service == nil {
		return false
	}

	addrs, err := Service.GetAccountAddrs(name)
	if err != nil {
		return false
	}
	for _, login := range addrs {
		if login.Addr == addr {
			return true
		}
	}
	return false
}

func Addr(addr string) {
	Pair("", addr)
}

func Pair(name string, addr string) {
	log := log.Add("Name", name).Add("Addr", addr)
	if Service == nil {
		log.Warn("httptrack: service is nil")
		return
	}
	if Match(name, addr) {
		return
	}

	if err := Service.SaveAccountAddr(name, addr, time.Now()); err == nil {
		log.Debug("httptrack: address recorded")
	} else {
		log.Add("Error", err).Error("httptrack: address recording")
	}
}
