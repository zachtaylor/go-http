package sessions // import "ztaylor.me/http/sessions"

import (
	"sync"
	"time"
)

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

func (t *T) UpdateTime() {
	t.Lock()
	t.time = time.Now()
	t.Unlock()
}

func (t *T) Revoke() {
	t.Lock()
	if t.done != nil {
		close(t.done)
		t.done = nil
	}
	t.Unlock()
}

func (t *T) Done() <-chan bool {
	return t.done
}

func (t T) String() string {
	return "sessions.Grant#" + t.id
}
