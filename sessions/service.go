package sessions

import "net/http"

// Service manages sessions
type Service interface {
	// Count returns number of active Sessions
	Count() int
	// Get returns a Session by id, if any
	Get(string) *T
	// Cookie returns the Session reffered by cookies, if valid
	Cookie(r *http.Request) (t *T)
	// Grant returns a new Session granted to the username
	Grant(string) *T
}
