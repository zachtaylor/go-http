package http

// Error is a string, used by const ErrX
type Error string

func (e Error) String() string {
	return string(e)
}

func (e Error) Error() string {
	return string(e)
}

// ErrRespondPathRaw is used as a return by a ServiceFunc
const ErrRespondPathRaw = Error("http: route cannot respond outside http")

// ErrCookieFormat is used as a return by a ReadRequestCookie
const ErrCookieFormat = Error("http: cookie format must parse uint")

// ErrCookieSession is used as a return by a ReadRequestCookie
const ErrCookieSession = Error("http: cookie session missing")
