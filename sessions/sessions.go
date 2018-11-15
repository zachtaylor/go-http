package sessions // import "ztaylor.me/http/sessions"

import (
	"sync"
	"time"
)

// Provider is a generic sessions container
type Provider interface {
	// Count returns len active Sessions
	Count() int
	// Get finds a Session with id
	SessionID(id string) *T
	// NewGrant creates a Grant with name
	NewGrant(name string) *T
	// Watch needs a goroutine to expire Sessions older than d
	Watch(d time.Duration)
}

// T is a Session
type T struct {
	id   string
	name string
	time time.Time
	done chan bool
	sync.Mutex
}

// New creates an initialized orphan Session
func New(name string) *T {
	return &T{
		name: name,
		time: time.Now(),
		done: make(chan bool),
	}
}

func (t *T) ID() string {
	return t.id
}
func (t *T) Name() string {
	return t.name
}
func (t *T) Time() time.Time {
	return t.time
}
func (t *T) UpdateTime() time.Time {
	return t.time
}
func (t *T) Revoke() {
	t.Lock()
	if t.done != nil {
		close(t.done)
		t.done = nil
	}
	t.Unlock()
}
func (t *T) Done() chan bool {
	return t.done
}
func (t T) String() string {
	return "sessions.Grant#" + t.id
}
