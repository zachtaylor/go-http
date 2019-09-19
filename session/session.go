package session // import "ztaylor.me/http/session"

import (
	"net/http"
	"sync"
	"time"
)

// New creates an initialized orphan Session
func New(name string) *T {
	return &T{
		name: name,
		time: time.Now(),
		done: make(chan bool),
	}
}

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

// Done returns the observe channel, or nil if the session is already closed
func (t *T) Done() <-chan bool {
	return t.done
}

// UpdateTime is used to reset the auth time
func (t *T) UpdateTime() {
	t.lock.Lock()
	t.time = time.Now()
	t.lock.Unlock()
}

// Close closes the observe channel, and sets it to nil
func (t *T) Close() {
	t.lock.Lock()
	if t.done != nil {
		close(t.done)
		t.done = nil
	}
	t.lock.Unlock()
}

// String returns a string representation of the session
func (t *T) String() string {
	if t == nil {
		return "(nil)"
	}
	return "SessionID#" + t.id
}

// WriteCookie is a convenience to write response header with this SessionID
func (t *T) WriteCookie(w http.ResponseWriter) {
	if t != nil {
		cookieWrite(w, t.id)
	} else {
		EraseSessionID(w)
	}
}

// EraseSessionID writes a SessionID header that is empty value
func EraseSessionID(w http.ResponseWriter) {
	cookieWrite(w, "")
}
