package ws

import "regexp"

// MatcherSet is used to combine Matchers into a single Matcher
type MatcherSet []Matcher

// Matches creates a new MatcherSet with the provided Matchers
//
// This is syntactic sugar for MatcherSet{...}
func Matches(matches ...Matcher) Matcher {
	return MatcherSet(matches)
}

// Match checks that all included Matchers return true
func (set MatcherSet) Match(m *Message) bool {
	for _, matcher := range set {
		if matcher == nil || !matcher.Match(m) {
			return false
		}
	}
	return true
}

type matcherFunc func(*Message) bool

func (f matcherFunc) Match(m *Message) bool {
	return f(m)
}

// MatcherFunc turns a func into a Matcher
func MatcherFunc(f func(*Message) bool) Matcher {
	return matcherFunc(f)
}

type matcherRegex struct {
	*regexp.Regexp
}

func (rgx *matcherRegex) Match(m *Message) bool {
	return rgx.MatchString(m.URI)
}

// MatcherRegex creates a regexp match check against Message.Name
func MatcherRegex(s string) Matcher {
	return &matcherRegex{regexp.MustCompile(s)}
}

type matcherLit string

func (s matcherLit) Match(m *Message) bool {
	return string(s) == m.URI
}

// MatcherLit creates a literal match check against Message.Name
func MatcherLit(s string) Matcher {
	return matcherLit(s)
}
