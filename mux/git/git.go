package git

import (
	"net/http"

	"github.com/AaronO/go-git-http"
	"github.com/AaronO/go-git-http/auth"
	"ztaylor.me/http/mux"
)

// Matcher returns true for http git requests
var Matcher = mux.MatcherFunc(func(r *http.Request) bool {
	if ua := r.Header["User-Agent"][0]; len(ua) > 2 && ua[:3] == "git" {
		return true
	}
	return false
})

// NewHandler creates a new default githttp Handler
func NewHandler(path string) http.Handler {
	return githttp.New(path)
}

// AuthNoPush is httpmuxgit middleware for restricting all git.Push requests
var AuthNoPush = auth.Authenticator(func(info auth.AuthInfo) (bool, error) {
	if info.Push {
		return false, nil
	}
	return true, nil
})

// NewRoute creates a default kind of http git route
//
// Uses Matcher, NewHandler, AuthNoPush
func NewRoute(path string) *mux.Route {
	return &mux.Route{
		Matcher: Matcher,
		Handler: AuthNoPush(NewHandler(path)),
	}
}
