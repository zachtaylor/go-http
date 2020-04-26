package session

import "net/http"

// Service manages sessions
type Service interface {
	// Count returns number of active Sessions
	Count() int
	// Cookie returns the Session reffered by cookies, if valid
	Cookie(*http.Request) (*T, error)
	// Get returns a Session by id, if any
	Get(string) *T
	// Find returns a Session by name, if any
	Find(string) *T
	// Grant returns a new Session granted to the username
	Grant(string) *T
}
