package ws

import "regexp"

// RouterSet is used to combine Routers into a single Router
type RouterSet []Router

// Match checks that all included Routers return true
func (set RouterSet) Route(m *Message) bool {
	for _, router := range set {
		if router == nil || !router.Route(m) {
			return false
		}
	}
	return true
}

// RouterFunc turns a func into a Router
type RouterFunc func(*Message) bool

// Route implements Router by calling the underlying func
func (f RouterFunc) Route(m *Message) bool {
	return f(m)
}

type routerRegex struct {
	*regexp.Regexp
}

func (rgx *routerRegex) Route(m *Message) bool {
	return rgx.MatchString(m.URI)
}

// RouterRegex creates a regexp match check against Message.Name
func RouterRegex(s string) Router {
	return &routerRegex{regexp.MustCompile(s)}
}

// RouterLit creates a literal match check against Message.Name
type RouterLit string

// Route implements Router by literally matching the Message.URI
func (s RouterLit) Route(m *Message) bool {
	return string(s) == m.URI
}
