package sessions // import "ztaylor.me/http/sessions"

import (
	"net/http"
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

// NewService creates a sessions Manager with expiry goroutine
func NewService(lifetime time.Duration) (s *Service) {
	s = &Service{
		keygen: Keygen,
		cache:  make(map[string]*T),
	}
	go s.watch(lifetime)
	return
}

// EraseSessionID writes a SessionID header that is empty value
func EraseSessionID(w http.ResponseWriter) {
	cookieWrite(w, "")
}
