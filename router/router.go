package router

import (
	"net/http"
	"strings"
)

// Bool satisfies mux.Router with a consistent response
type Bool bool

// Route satisfies mux.Router by returning the bool
func (b Bool) Route(_ *http.Request) bool {
	return bool(b)
}

// True is a router that always returns true
var True Bool = true

// And creates a Router from any number of Routers that requires all return true
type And []interface {
	Route(*http.Request) bool
}

// Route satisfies Router by verifying all routers return true
func (routers And) Route(r *http.Request) bool {
	for _, router := range routers {
		if !router.Route(r) {
			return false
		}
	}
	return true
}

// Or creates a Router from any number of Routers that requires any return true
type Or []interface {
	Route(*http.Request) bool
}

// Route satisfies Router by verifying a router returns true
func (routers Or) Route(r *http.Request) bool {
	for _, router := range routers {
		if router.Route(r) {
			return true
		}
	}
	return false
}

// Host creates a Router for Request host from given string
type Host string

// Route satisfies Router by matching the Request host
func (host Host) Route(r *http.Request) bool {
	return string(host) == r.Host
}

// Func casts Router from a basic func
type Func func(*http.Request) bool

// Route satisfies Router by calling the func
func (f Func) Route(r *http.Request) bool {
	return f(r)
}

// HTTP is a Router that returns if Request TLS is nil
var HTTP = Func(func(r *http.Request) bool {
	return r.TLS == nil
})

// HTTPS is a Router that returns if Request TLS is non-nil
var HTTPS = Func(func(r *http.Request) bool {
	return r.TLS != nil
})

// SinglePageApp is a Router that checks for Single Page App response
//
// http.Request.Method is GET
// http.Request.URL.Path does not have file ext after last /
// http.Request.Header["Accept"] contains "text/html"
var SinglePageApp = Func(func(r *http.Request) bool {
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

// UserAgent is a Router for matching User-Agent
type UserAgent string

// Route satisfies   by matching the first chars of User-Agent
func (s UserAgent) Route(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	ls := len(s)
	return len(ua) >= ls && ua[:ls] == string(s)
}
