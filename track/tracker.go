package httptrack

import (
	"time"
	"ztaylor.me/log"
)

func Addr(addr string) {
	Pair("", addr)
}

func Pair(name string, addr string) {
	log := log.Add("Name", name).Add("Addr", addr)
	if Service == nil {
		log.Warn("httptrack: service is nil")
		return
	}

	if err := Service.SaveAccountAddr(name, addr, time.Now()); err == nil {
		log.Debug("httptrack: address recorded")
	} else {
		log.Add("Error", err).Error("httptrack: address recording")
	}
}
