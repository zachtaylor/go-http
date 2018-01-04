package http

import "net/http"

type ResponderFunc func(*Request) error

func Responder(f ResponderFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(RequestFromNet(r, w))
	})
}
