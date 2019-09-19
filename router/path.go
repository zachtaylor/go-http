package router

import "net/http"

// Path satisfies mux.Router by matching r.URL.Path
type Path string

// Route satisfies Router by matching the request path literally
func (path Path) Route(r *http.Request) bool {
	return string(path) == r.URL.Path
}

// PathStarts satisfies mux.Router by matching path starting with given prefix
type PathStarts string

// Route satisfies Router by matching the path prefix
func (prefix PathStarts) Route(r *http.Request) bool {
	lp := len(prefix)
	if len(r.URL.Path) < lp {
		return false
	}
	return string(prefix) == r.URL.Path[:lp]
}
