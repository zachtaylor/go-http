package session

import "errors"

// ErrCookieNotFound is returned by `Cache.Cookie` when the `Request` does not contain the SessionID Cookie
var ErrCookieNotFound = errors.New("cookie not found")

// ErrCookieInvalid is returned by `Cache.Cookie` when the `Request` contains a SessionID Cookie with an invalid (or expired) SessionID
var ErrCookieInvalid = errors.New("cookie invalid")
