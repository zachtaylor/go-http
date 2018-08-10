package http

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
