package http

import (
	"regexp"
)

type Matcher interface {
	Match(string) bool
}

type StringMatcher string

func (s StringMatcher) Match(s2 string) bool {
	return string(s) == s2
}

func RegexMatcher(s string) Matcher {
	return &regexMatcher{regexp.MustCompile(s)}
}

type regexMatcher struct {
	*regexp.Regexp
}

func (r *regexMatcher) Match(s string) bool {
	return r.MatchString(s)
}
