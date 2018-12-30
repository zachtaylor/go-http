package mux // import "ztaylor.me/http/mux"

import (
	"net/http"
	"regexp"
	"strings"
)

// Mux is set of Routers
//
// provides http.Handler
// safely invokes each router, using Router.Matcher in the order added to choose Handler
type Mux struct {
	Routes []Router
}

// NewServer creates a Server container
func NewMux() *Mux {
	return &Mux{
		Routes: make([]Router, 0),
	}
}

// Router appends a Router to this Mux
func (mux *Mux) Router(r Router) {
	mux.Routes = append(mux.Routes, r)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var router Router
	for _, route := range mux.Routes {
		if route.Route(r) {
			router = route
			break
		}
	}
	if router != nil {
		router.ServeHTTP(w, r)
	}
}

// Matcher sorts Requests
type Matcher interface {
	Route(*http.Request) bool
}

// Matches is used to combine Matchers into a single Matcher
type Matches []Matcher

// Append alters this Branch with a new Matcher
func (s *Matches) Append(r Matcher) {
	*s = append(*s, r)
}

// Route checks that all included Matchers return true
func (s Matches) Route(r *http.Request) bool {
	for _, matcher := range s {
		if matcher == nil || !matcher.Route(r) {
			return false
		}
	}
	return true
}

type matcherFunc func(*http.Request) bool

func (f matcherFunc) Route(r *http.Request) bool {
	return f(r)
}

// MatcherFunc turns a func into a Matcher
func MatcherFunc(f func(*http.Request) bool) Matcher {
	return matcherFunc(f)
}

type matcherRegex struct {
	*regexp.Regexp
}

func (rgx *matcherRegex) Route(r *http.Request) bool {
	return rgx.MatchString(r.URL.Path)
}

// MatcherRegex creates a regexp match check against http.Request.RequestURI
func MatcherRegex(s string) Matcher {
	return &matcherRegex{regexp.MustCompile(s)}
}

type matcherLit string

func (s matcherLit) Route(r *http.Request) bool {
	return string(s) == r.URL.Path
}

// MatcherLit creates a literal match check against http.Request.RequestURI
func MatcherLit(s string) Matcher {
	return matcherLit(s)
}

type matcherHost string

func (host matcherHost) Route(r *http.Request) bool {
	return string(host) == r.Host
}

// MatcherHost creates a literal match check against http.Request.Host
func MatcherHost(s string) Matcher {
	return matcherHost(s)
}

// MatcherHTTP checks http.Request.TLS is nil or has a ServerName
var MatcherHTTP = MatcherFunc(func(r *http.Request) bool {
	return r.TLS == nil || r.TLS.ServerName == ""
})

// MatcherHTTPS checks http.Request.TLS has a ServerName
var MatcherHTTPS = MatcherFunc(func(r *http.Request) bool {
	return r.TLS != nil && r.TLS.ServerName != ""
})

// Map is shorthand for AddRoute
func (mux *Mux) Map(m Matcher, h http.Handler) {
	mux.Router(&Route{
		Handler: h,
		Matcher: m,
	})
}

// MapLit is shorthand for Map with MatcherLit
func (mux *Mux) MapLit(path string, h http.Handler) {
	mux.Map(MatcherLit(path), h)
}

type matcherMethod string

func (method matcherMethod) Match(r *http.Request) bool {
	return string(method) == r.Method
}

// MatcherGET is a Matcher that checks http.Request.Method is GET
var MatcherGET = matcherMethod(http.MethodGet)

// MatcherDELETE is a Matcher that checks http.Request.Method is DELETE
var MatcherDELETE = matcherMethod(http.MethodDelete)

// MatcherPOST is a Matcher that checks http.Request.Method is POST
var MatcherPOST = matcherMethod(http.MethodPost)

// MatcherPUT is a Matcher that checks http.Request.Method is PUT
var MatcherPUT = matcherMethod(http.MethodPut)

// MatcherSPA is a Matcher that checks for Single Page App response
//
// http.Request.Method is GET
// http.Request.URL.Path does not have file ext after last /
// http.Request.Header["Accept"] contains "text/html"
var MatcherSPA = matcherFunc(func(r *http.Request) bool {
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

// MapRgx is shorthand for Map with MatcherRegex
func (mux *Mux) MapRgx(path string, h http.Handler) {
	mux.Map(MatcherRegex(path), h)
}

// MatcherGit checks r.Header["User-Agent"] is "git"
var MatcherGit = MatcherFunc(func(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	return len(ua) > 2 && ua[:3] == "git"
})

// MatcherGoGet checks r.Header["User-Agent"] is "Go-http-client"
var MatcherGoGet = matcherFunc(func(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	return len(ua) > 13 && ua[:14] == "Go-http-client"
})

// Route points to http.Handler, Matcher
type Route struct {
	http.Handler
	Matcher
}

// Router is a union type of http.Handler, Matcher
type Router interface {
	http.Handler
	Matcher
}
