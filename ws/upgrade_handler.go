package ws

import (
	"net/http"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/sessions"
	"ztaylor.me/log"
)

func UpgradeHandler(h Handler) http.Handler {
	return websocket.Handler(func(conn *websocket.Conn) {
		log.Add("Addr", conn.Request().RemoteAddr).Debug("http/socket_handler")

		socket := Open(conn)

		if session := sessions.ReadCookie(conn.Request()); session != nil {
			socket.Login(session)
		}

		socket.Watch(h)
	})
}
