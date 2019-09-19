package git // import "ztaylor.me/http/handler/git"

import (
	"net/http"

	githttp "github.com/AaronO/go-git-http"
	"github.com/AaronO/go-git-http/auth"
)

// Default creates a new default githttp Handler
//
// Push is disabled
func Default(path string) http.Handler {
	return authNoPush(githttp.New(path))
}

// authNoPush is httpmuxgit middleware for restricting all git.Push requests
var authNoPush = auth.Authenticator(func(info auth.AuthInfo) (bool bool, error error) {
	return !info.Push, nil
})
