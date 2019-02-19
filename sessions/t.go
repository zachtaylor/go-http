package sessions // import "ztaylor.me/http/sessions"

import (
	"net/http"
	"sync"
	"time"
)

// T is a Session
type T struct {
	id   string
	name string
	time time.Time
	done chan bool
	lock sync.Mutex
}

// ID returns the SessionID
func (t *T) ID() string {
	return t.id
}

// Name returns the original name of the session
func (t *T) Name() string {
	return t.name
}

// Time returns the auth time
func (t *T) Time() time.Time {
	return t.time
}

// UpdateTime is used to reset the auth time
func (t *T) UpdateTime() {
	t.lock.Lock()
	t.time = time.Now()
	t.lock.Unlock()
}

// Revoke breaks the session voluntarily and immediately
func (t *T) Revoke() {
	t.lock.Lock()
	if t.done != nil {
		close(t.done)
		t.done = nil
	}
	t.lock.Unlock()
}

// Done returns the observe channel, or nil if the session is already closed
func (t *T) Done() <-chan bool {
	return t.done
}

// String returns a string representation of the session
func (t T) String() string {
	return "SessionID#" + t.id
}

// WriteCookie is a convenience to write response header with this SessionID
func (t *T) WriteCookie(w http.ResponseWriter) {
	if t != nil {
		cookieWrite(w, t.id)
	} else {
		cookieWrite(w, "")
	}
}
