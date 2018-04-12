package http

import (
	"regexp"
)

type Matcher interface {
	Match(*Request) bool
}

type QuestMatcher string

func (quest QuestMatcher) Match(r *Request) bool {
	return string(quest) == r.Quest
}

func RegexMatcher(s string) Matcher {
	return &regexMatcher{regexp.MustCompile(s)}
}

type regexMatcher struct {
	*regexp.Regexp
}

func (regex *regexMatcher) Match(r *Request) bool {
	return regex.MatchString(r.Quest)
}
