package mux

import "net/http"

// RouterAny creates a router with a consistent response
type RouterAny bool

// Route satisfies Router by returning the bool
func (router RouterAny) Route(_ *http.Request) bool {
	return bool(router)
}

// RoutersAnd creates a Router from any number of Routers that requires all return true
func RoutersAnd(routers ...Router) Router {
	return RouterFunc(func(r *http.Request) bool {
		for _, router := range routers {
			if !router.Route(r) {
				return false
			}
		}
		return true
	})
}

// RoutersOr creates a Router from any number of Routers that requires any return true
func RoutersOr(routers ...Router) Router {
	return RouterFunc(func(r *http.Request) bool {
		for _, router := range routers {
			if router.Route(r) {
				return true
			}
		}
		return false
	})
}

// RouterHost creates a Router for Request host from given string
type RouterHost string

// Route satisfies Router by matching the Request host
func (router RouterHost) Route(r *http.Request) bool {
	return string(router) == r.Host
}

// RouterFunc casts Router from a basic func
type RouterFunc func(*http.Request) bool

// Route satisfies Router by calling the func
func (f RouterFunc) Route(r *http.Request) bool {
	return f(r)
}

// RouterHTTP is a Router that returns if Request TLS is nil
var RouterHTTP = RouterFunc(func(r *http.Request) bool {
	return r.TLS == nil
})

// RouterHTTPS is a Router that returns if Request TLS is non-nil
var RouterHTTPS = RouterFunc(func(r *http.Request) bool {
	return r.TLS != nil
})

// RouterUserAgent is a Router for matching User-Agent
type RouterUserAgent string

// Route satisfies Router by matching the first chars of User-Agent
func (s RouterUserAgent) Route(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	ls := len(s)
	return len(ua) >= ls && ua[:ls] == string(s)
}

// RouterGoGet checks r.Header["User-Agent"] is "Go-http-client"
var RouterGoGet = RouterUserAgent("Go-http-client")

// RouterGit checks r.Header["User-Agent"] is "git"
var RouterGit = RouterUserAgent("git")

// RouterPath creates a Router for "== r.URL.Path"
type RouterPath string

func (s RouterPath) Route(r *http.Request) bool {
	return string(s) == r.URL.Path
}

// RouterPathStarts is a Router for Request path starting with given string
type RouterPathStarts string

// Route satisfies Router by matching the given prefix
func (router RouterPathStarts) Route(r *http.Request) bool {
	if len(r.URL.Path) < len(router) {
		return false
	}
	return string(router) == r.URL.Path[:len(router)]
}

type routerMethod string

func (router routerMethod) Route(r *http.Request) bool {
	return string(router) == r.Method
}

var (
	// RouterCONNECT is a Router that returns if Request method is CONNECT
	RouterCONNECT = routerMethod("CONNECT")

	// RouterDELETE is a Router that returns if Request method is DELETE
	RouterDELETE = routerMethod("DELETE")

	// RouterGET is a Router that returns if Request method is GET
	RouterGET = routerMethod("GET")

	// RouterHEAD is a Router that returns if Request method is HEAD
	RouterHEAD = routerMethod("HEAD")

	// RouterOPTIONS is a Router that returns if Request method is OPTIONS
	RouterOPTIONS = routerMethod("OPTIONS")

	// RouterPOST is a Router that returns if Request method is POST
	RouterPOST = routerMethod("POST")

	// RouterPUT is a Router that returns if Request method is PUT
	RouterPUT = routerMethod("PUT")

	// RouterTRACE is a Router that returns if Request method is TRACE
	RouterTRACE = routerMethod("TRACE")
)
