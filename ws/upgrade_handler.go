package ws

import (
	"net/http"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/sessions"
	"ztaylor.me/log"
)

// UpgradeHandler provides a websocket handshake func
func UpgradeHandler(sessions *sessions.Service, h Handler) http.Handler {
	socks := &Service{
		cache: make(map[string]*Socket),
		mux:   make(Mux, 0),
	}
	return websocket.Handler(func(conn *websocket.Conn) {
		log.Add("Addr", conn.Request().RemoteAddr).Debug("http/socket_handler")

		socket := NewSocket(conn)

		if session := sessions.ReadRequestCookie(conn.Request()); session != nil {
			socket.Login(socks, session)
		}

		socket.watch(socks, h)
	})
}
