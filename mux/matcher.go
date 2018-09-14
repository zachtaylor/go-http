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

type matcherJSON struct {
}

func (_ *matcherJSON) Match(r *http.Request) bool {
	return strings.HasSuffix(r.RequestURI, ".json")
}

// MatcherJSON returns a Matcher that checks http.Request.RequestURI ends with .json
func MatcherJSON() *matcherJSON {
	return nil
}

type matcherMethod string

func (method matcherMethod) Match(r *http.Request) bool {
	return string(method) == r.Method
}

// MatcherGET returns a Matcher that checks http.Request.Method is GET
func MatcherGET() Matcher {
	return matcherMethod(http.MethodGet)
}

// MatcherDELETE returns a Matcher that checks http.Request.Method is DELETE
func MatcherDELETE() Matcher {
	return matcherMethod(http.MethodDelete)
}

// MatcherPOST returns a Matcher that checks http.Request.Method is POST
func MatcherPOST() Matcher {
	return matcherMethod(http.MethodPost)
}

// MatcherPUT returns a Matcher that checks http.Request.Method is PUT
func MatcherPUT() Matcher {
	return matcherMethod(http.MethodPut)
}
