package http

import "net/http"

func Responder(f func(*Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(RequestFromNet(r, w))
	})
}
