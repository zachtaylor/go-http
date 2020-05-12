package acme

import (
	"net/http"

	"ztaylor.me/http/mux"
	"ztaylor.me/http/router"
)

// Thumbprint creates a new *mux.Route for the given file system path to use for stateless ACME challenges on route "/.well-known/acme-challenge/"
func Thumbprint(thumbprint string) *mux.Route {
	lencut := 28 // len("/.well-known/acme-challenge/")
	addthumb := "." + thumbprint
	return &mux.Route{
		Router: router.PathStarts("/.well-known/acme-challenge"),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.URL.Path) < lencut {
				w.Write([]byte("error: path too short"))
				return
			}
			match := r.URL.Path[lencut:]
			w.Write([]byte(match + addthumb))
		}),
	}
}
