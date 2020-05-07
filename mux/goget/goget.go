package goget

import (
	"ztaylor.me/http/handler"
	"ztaylor.me/http/mux"
	"ztaylor.me/http/router"
)

// Domain creates a new *mux.Route which handles go get style challenges for the given domain name
func Domain(domain string) *mux.Route {
	return &mux.Route{
		Router:  router.UserAgent("Go-http-client"),
		Handler: handler.GoGet(domain),
	}
}
