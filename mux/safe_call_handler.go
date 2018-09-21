package mux

import (
	"net/http"

	"ztaylor.me/log"
)

// SafeCallHandler uses log.Protect() to protect a http.Handler call
func SafeCallHandler(f http.Handler, w http.ResponseWriter, r *http.Request) {
	log.Protect(func() {
		f.ServeHTTP(w, r)
	})
}
