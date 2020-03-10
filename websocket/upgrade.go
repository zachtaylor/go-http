package websocket

import (
	"net/http"

	"golang.org/x/net/websocket"
)

// UpgradeHandler provides a websocket handshake func
func UpgradeHandler(s Service) http.Handler {
	return websocket.Handler(func(c *Conn) {
		s.Connect(c)
	})
}
