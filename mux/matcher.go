package mux

import (
	"net/http"
	"regexp"
	"strings"
)

// Matcher tests a Request
type Matcher interface {
	Match(*http.Request) bool
}

// MatcherSet is used to combine Matchers into a single Matcher
type MatcherSet []Matcher

// Matches creates a new MatcherSet with the provided Matchers
//
// This is syntactic sugar for MatcherSet{...}
func Matches(matches ...Matcher) Matcher {
	return MatcherSet(matches)
}

// Match checks that all included Matchers return true
func (set MatcherSet) Match(r *http.Request) bool {
	for _, matcher := range set {
		if matcher == nil || !matcher.Match(r) {
			return false
		}
	}
	return true
}

type matcherFunc func(*http.Request) bool

func (f matcherFunc) Match(r *http.Request) bool {
	return f(r)
}

// MatcherFunc turns a func into a Matcher
func MatcherFunc(f func(*http.Request) bool) Matcher {
	return matcherFunc(f)
}

type matcherRegex struct {
	*regexp.Regexp
}

func (rgx *matcherRegex) Match(r *http.Request) bool {
	return rgx.MatchString(r.RequestURI)
}

// MatcherRegex creates a regexp match check against http.Request.RequestURI
func MatcherRegex(s string) Matcher {
	return &matcherRegex{regexp.MustCompile(s)}
}

type matcherLit string

func (s matcherLit) Match(r *http.Request) bool {
	return string(s) == r.RequestURI
}

// MatcherLit creates a literal match check against http.Request.RequestURI
func MatcherLit(s string) Matcher {
	return matcherLit(s)
}

type matcherHost string

func (host matcherHost) Match(r *http.Request) bool {
	return string(host) == r.Host
}

// MatcherHost creates a literal match check against http.Request.Host
func MatcherHost(s string) Matcher {
	return matcherHost(s)
}

// MatcherJSON checks http.Request.RequestURI ends with .json
var MatcherJSON = MatcherFunc(func(r *http.Request) bool {
	return strings.HasSuffix(r.RequestURI, ".json")
})

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
	if r.Method != http.MethodGet {
		return false
	}
	if r.URL.Path == "/" {
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

// MatcherGit checks r.Header["User-Agent"] is "git"
var MatcherGit = MatcherFunc(func(r *http.Request) bool {
	if ua := r.Header["User-Agent"][0]; len(ua) > 2 && ua[:3] == "git" {
		return true
	}
	return false
})

// MatcherGoGet checks r.Header["User-Agent"] is "Go-http-client"
var MatcherGoGet = matcherFunc(func(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	return len(ua) > 13 && ua[:14] == "Go-http-client"
})
