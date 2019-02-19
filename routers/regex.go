package routers

import (
	"net/http"
	"regexp"

	"ztaylor.me/http/mux"
)

type routerRegex struct {
	*regexp.Regexp
}

func (rgx *routerRegex) Route(r *http.Request) bool {
	return rgx.MatchString(r.URL.Path)
}

// RouterRegex creates a regexp match check against http.Request.RequestURI
func RouterRegex(s string) mux.Router {
	return &routerRegex{regexp.MustCompile(s)}
}
