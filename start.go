package http

import (
	"net/http"

	"ztaylor.me/log"
)

func Start(port string) {
	StartServer(port, DefaultServer())
}

func StartServer(port string, h http.Handler) {
	log.Error(http.ListenAndServe(port, h))
}

func StartTLS(cert string, key string, h http.Handler) {
	go StartServer(":80", http.HandlerFunc(redirectHttps))
	log.Error(http.ListenAndServeTLS(":443", cert, key, h))
}

func redirectHttps(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req,
		"https://"+req.Host+req.URL.String(),
		http.StatusMovedPermanently,
	)
}
