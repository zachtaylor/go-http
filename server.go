package http

import (
	"net/http"
)

type server int

var Server = server(1)

func (_ server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router {
		if route.Match(r.RequestURI) {
			route.ServeHTTP(w, r)
			return
		}
	}
}
