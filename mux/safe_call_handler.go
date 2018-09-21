package mux

import (
	"net/http"

	"ztaylor.me/log"
)

// SafeCallHandler uses defer recover() to protect a http.Handler call
func SafeCallHandler(f http.Handler, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			log.WithFields(log.Fields{
				"Host":    r.Host,
				"Path":    r.RequestURI,
				"Cookies": r.Cookies,
				"Remote":  r.RemoteAddr,
				"Error":   e,
			}).Error("mux: recover")
		}
	}()
	f.ServeHTTP(w, r)
}
