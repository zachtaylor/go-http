package session

import (
	"net/http"

	"ztaylor.me/http/cookies"
)

// cookieName is the string "SessionID"
const cookieName = "SessionID"

func cookieWrite(w http.ResponseWriter, val string) {
	cookies.Write(w, cookieName, val)
}

func cookieRead(r *http.Request) (string, error) {
	return cookies.Read(r, cookieName)
}
