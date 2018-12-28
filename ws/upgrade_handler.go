package ws

import (
	"net/http"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/sessions"
	"ztaylor.me/log"
)

// UpgradeHandler provides a websocket handshake func
func UpgradeHandler(h Handler) http.Handler {
	return websocket.Handler(func(conn *websocket.Conn) {
		log.Add("Addr", conn.Request().RemoteAddr).Debug("http/socket_handler")

		socket := NewSocket(conn)

		if session := sessions.FromRequestCookie(conn.Request()); session != nil {
			socket.Login(session)
		}

		socket.watch(h)
	})
}
