package http

import (
	"net/http"
	"strconv"
)

var SessionService interface {
	// Count returns the number of active sessions
	Count() int
	// Get returns a session with the given id, if it exists
	Get(uint) *Session
	// Find returns all sessions associated with the username
	Find(string) []*Session
	// Grant creates a new session for the username
	Grant(string) *Session
	// Revoke deletes a session, if it exists
	Revoke(uint)
}

func ReadRequestCookie(r *http.Request) (*Session, error) {
	if sessionCookie, err := r.Cookie("SessionID"); err != nil {
		return nil, err
	} else if sessionID, err := strconv.ParseUint(sessionCookie.Value, 10, 0); err != nil {
		return nil, ErrCookieFormat
	} else if session := SessionService.Get(uint(sessionID)); session == nil {
		return nil, ErrCookieSession
	} else {
		return session, nil
	}
}
