package acme

import (
	"net/http"

	"ztaylor.me/http/mux"
	"ztaylor.me/http/router"
)

// Path creates a new *mux.Route for the given file system path to use for ACME challenges on route "/.well-known/acme-challenge/"
func Path(path string) *mux.Route {
	return &mux.Route{
		Router:  router.Path("/.well-known/acme-challenge/"),
		Handler: http.FileServer(http.Dir(path)),
	}
}
