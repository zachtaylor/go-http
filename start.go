package http

import (
	"net/http"
	"ztaylor.me/log"
)

func Start(port string) {
	log.Error(http.ListenAndServe(port, Server))
}

func StartTLS(cert string, key string) {
	go http.ListenAndServe(":80", http.HandlerFunc(redirectHttps))
	log.Error(http.ListenAndServeTLS(":443", cert, key, Server))
}

func redirectHttps(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req,
		"https://"+req.Host+req.URL.String(),
		http.StatusMovedPermanently)
}
