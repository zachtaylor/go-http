package git // import "ztaylor.me/http/mux/git"

import (
	"net/http"

	"github.com/AaronO/go-git-http"
	"github.com/AaronO/go-git-http/auth"
	"ztaylor.me/http/mux"
)

// NewHandler creates a new default githttp Handler
func NewHandler(path string) http.Handler {
	return githttp.New(path)
}

// AuthNoPush is httpmuxgit middleware for restricting all git.Push requests
var AuthNoPush = auth.Authenticator(func(info auth.AuthInfo) (bool bool, error error) {
	bool = !info.Push
	return
})

// NewRoute creates a default kind of http git route
//
// Uses Matcher, NewHandler, AuthNoPush
func NewRoute(path string) *mux.Route {
	return &mux.Route{
		Router:  mux.RouterGit,
		Handler: AuthNoPush(NewHandler(path)),
	}
}
