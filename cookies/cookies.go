package cookies // import "ztaylor.me/http/cookies"

import (
	"net/http"
)

// Read a cookie value
func Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Write a cookie value with Path=/; and no time
func Write(w http.ResponseWriter, name string, val string) {
	w.Header().Set("Set-Cookie", name+"="+val+"; Path=/;")
}
