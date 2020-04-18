package session // import "ztaylor.me/http/session"

import (
	"net/http"
	"time"
)

// T is a Session
type T struct {
	id   string
	name string
	in   chan bool
	done chan bool
}

// New creates an initialized orphan Session
func New(id, name string, d time.Duration) (t *T) {
	t = &T{
		id:   id,
		name: name,
		in:   make(chan bool),
		done: make(chan bool),
	}
	go t.watch(d)
	return
}

// ID returns the Session ID
func (t *T) ID() string {
	return t.id
}

// Name returns the name of this Session
func (t *T) Name() string {
	return t.name
}

// Refresh sends a refresh signal
func (t *T) Refresh() {
	go t.send(true)
}

// Close sends a close signal
func (t *T) Close() {
	go t.send(false)
}

// Done returns the observe channel, or nil if the Session is already closed
func (t *T) Done() <-chan bool {
	return t.done
}

// String returns a string representation of this Session
func (t *T) String() string {
	if t == nil {
		return "nil"
	}
	return "Session#" + t.id
}

func (t *T) send(ok bool) {
	t.in <- ok
}
func (t *T) close() {
	close(t.in)
	close(t.done)
}
func (t *T) watch(d time.Duration) {
	defer t.close()
	timer := time.NewTimer(d)
	for {
		select {
		case ok := <-t.in:
			if !timer.Stop() { // fail to stop
				<-timer.C // drain the channel
			}
			if !ok { // signal close
				return
			} // signal refresh
			timer.Reset(d)
		case <-timer.C:
			return
		}
	}
}

// WriteCookie is a convenience to write SessionID header cookie
func (t *T) WriteCookie(w http.ResponseWriter) {
	if t != nil {
		cookieWrite(w, t.id)
	} else {
		EraseSessionID(w)
	}
}

// EraseSessionID writes a SessionID header cookie that is empty value
func EraseSessionID(w http.ResponseWriter) {
	cookieWrite(w, "")
}
