package git // import "ztaylor.me/http/mux/git"

import (
	"ztaylor.me/http/handler/git"
	"ztaylor.me/http/mux"
	"ztaylor.me/http/router"
)

// Path creates a new git route
func Path(path string) *mux.Route {
	return &mux.Route{
		Router:  router.UserAgent("git"),
		Handler: git.Default(path),
	}
}
