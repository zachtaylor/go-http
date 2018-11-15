package sessions

import (
	"net/http"

	"ztaylor.me/http/cookies"
)

func ReadCookie(r *http.Request) *T {
	if cookie, err := cookies.Read(r, "SessionID"); err != nil {
		return nil
	} else if session := Service.SessionID(cookie); session == nil {
		return nil
	} else {
		return session
	}
}

func WriteCookie(w http.ResponseWriter, session *T) {
	cookies.Write(w, "SessionID", session.id)
}
