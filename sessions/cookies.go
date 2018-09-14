package sessions

import (
	"net/http"

	"ztaylor.me/http/cookies"
)

func ReadCookie(r *http.Request) *Grant {
	if cookie, err := cookies.Read(r, "SessionID"); err != nil {
		return nil
	} else if session := Service.Get(cookie); session == nil {
		return nil
	} else {
		return session
	}
}

func WriteCookie(w http.ResponseWriter, session *Grant) {
	cookies.Write(w, "SessionID", session.ID)
}
