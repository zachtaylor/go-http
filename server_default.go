package http

import "net/http"

var defaultServer *Router

func DefaultServer() *Router {
	if defaultServer == nil {
		defaultServer = NewRouter()
	}
	return defaultServer
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DefaultServer().ServeHTTP(w, r)
}

func Dispatch(q *Request) {
	DefaultServer().Dispatch(q)
}

func MapLit(s string, f func(*Request) error) {
	DefaultServer().AddRoute(&ServiceRoute{QuestMatcher(s), f})
}

func MapRgx(s string, f func(*Request) error) {
	DefaultServer().AddRoute(&ServiceRoute{RegexMatcher(s), f})
}

func MapRawLit(s string, h http.Handler) {
	DefaultServer().AddRoute(&NetHttpRoute{QuestMatcher(s), h})
}

func MapRawRgx(s string, h http.Handler) {
	DefaultServer().AddRoute(&NetHttpRoute{RegexMatcher(s), h})
}
