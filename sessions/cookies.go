package sessions

import (
	"net/http"

	"ztaylor.me/http/cookies"
)

// cookieName is the key "SessionID"
const cookieName = "SessionID"

func cookieWrite(w http.ResponseWriter, val string) {
	cookies.Write(w, cookieName, val)
}

func cookieRead(r *http.Request) (string, error) {
	return cookies.Read(r, cookieName)
}
