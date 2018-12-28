package sessions

import (
	"net/http"

	"ztaylor.me/http/cookies"
)

// CookieName is the string key SessionID
const CookieName = "SessionID"

// FromRequestCookie returns the session referred by the cookie, if valid
func FromRequestCookie(r *http.Request) *T {
	if cookie, err := cookies.Read(r, CookieName); err != nil {
		return nil
	} else if session := Service.SessionID(cookie); session == nil {
		return nil
	} else {
		return session
	}
}

// WriteCookie is a convenience to write response header
func WriteCookie(w http.ResponseWriter, session *T) {
	cookies.Write(w, CookieName, session.id)
}

// EraseCookie is a convenience to write response header
func EraseCookie(w http.ResponseWriter) {
	cookies.Write(w, CookieName, "")
}
