package cookies // import "ztaylor.me/http/cookies"

import "net/http"

// Read a cookie value
func Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// WriteLax adds a Set-Cookie header to w, with no expire time
func WriteLax(w http.ResponseWriter, name, val string) {
	w.Header().Set("Set-Cookie", name+"="+val+"; Path=/; SameSite=Lax;")
}

// WriteLaxExpired adds a Set-Cookie header to w, with value missing, and is expired
func WriteLaxExpired(w http.ResponseWriter, name string) {
	w.Header().Set("Set-Cookie", name+"=; Path=/; Expires==Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax; ")
}

// WriteSecure adds a Set-Cookie header to w, with Secure, and with no expire time
func WriteSecure(w http.ResponseWriter, name, val string) {
	w.Header().Set("Set-Cookie", name+"="+val+"; Path=/; Secure; SameSite=Strict;")
}

// WriteSecureExpired adds a Set-Cookie header to w, with Secure, and with value missing, and is expired
func WriteSecureExpired(w http.ResponseWriter, name string) {
	w.Header().Set("Set-Cookie", name+"=; Path=/; Expires==Thu, 01 Jan 1970 00:00:00 GMT; Secure; SameSite=Strict; ")
}
