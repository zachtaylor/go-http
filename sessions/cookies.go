package sessions

import (
	"net/http"

	"ztaylor.me/http/cookies"
)

// CookieName is the string key SessionID
const CookieName = "SessionID"

func cookieWrite(w http.ResponseWriter, val string) {
	cookies.Write(w, CookieName, val)
}

func cookieRead(r *http.Request) (string, error) {
	return cookies.Read(r, CookieName)
}
