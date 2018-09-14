package git

import (
	"net/http"

	"github.com/AaronO/go-git-http"
	"github.com/AaronO/go-git-http/auth"
	"ztaylor.me/http/mux"
)

// Matcher provides http/mux.Matcher for http git requests
type Matcher struct {
}

// Match provides http/mux.Matcher for http git requests
func (_ *Matcher) Match(r *http.Request) bool {
	if ua := r.Header["User-Agent"][0]; len(ua) > 2 && ua[:3] == "git" {
		return true
	}
	return false
}

// NewMatcher creates a new Matcher that returns true for http git requests
func NewMatcher() *Matcher {
	return nil
}

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
// Uses NewMatcher, NewHandler, AuthNoPush
func NewRoute(path string) *mux.Route {
	return &mux.Route{
		Matcher: NewMatcher(),
		Handler: AuthNoPush(NewHandler(path)),
	}
}
