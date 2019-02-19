package routers

import (
	"net/http"
	"strings"

	"ztaylor.me/http/mux"
)

// SinglePageApp is a Router that checks for Single Page App response
//
// http.Request.Method is GET
// http.Request.URL.Path does not have file ext after last /
// http.Request.Header["Accept"] contains "text/html"
var SinglePageApp = mux.RouterFunc(func(r *http.Request) bool {
	if r.Method != http.MethodGet || r.URL.Path == "/" {
		return false
	}
	path := r.URL.Path
	if i := strings.LastIndex(path, "/"); i > 1 {
		path = path[i:]
	}
	if strings.Contains(path, ".") {
		return false
	}
	return strings.Contains(r.Header.Get("Accept"), "text/html")
})
