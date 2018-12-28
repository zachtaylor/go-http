package handlers

import "net/http"

// RedirectHost is a Handler that uses hostname rewrite redirect http.StatusMovedPermanently
func RedirectHost(host string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proto := "http"
		if r.TLS != nil && r.TLS.ServerName != "" {
			proto += "s"
		}
		http.Redirect(w, r, proto+"://"+host+r.URL.String(), http.StatusMovedPermanently)
	})
}
