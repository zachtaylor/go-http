package session

import (
	"net/http"

	"ztaylor.me/http/cookies"
)

// cookieName is the string "SessionID"
const cookieName = "SessionID"

func cookieRead(r *http.Request) (string, error) {
	return cookies.Read(r, cookieName)
}

func cookieWrite(w http.ResponseWriter, val string, secure bool) {
	if secure {
		cookies.WriteSecure(w, cookieName, val)
	} else {
		cookies.WriteLax(w, cookieName, val)
	}
}

func cookieErase(w http.ResponseWriter, secure bool) {
	if secure {
		cookies.WriteSecureExpired(w, cookieName)
	} else {
		cookies.WriteLaxExpired(w, cookieName)
	}
}
