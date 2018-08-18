package http

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/log"
)

var SocketHandler = websocket.Handler(func(conn *websocket.Conn) {
	log.Add("Addr", conn.Request().RemoteAddr).Debug("http/socket_handler")

	socket := Open(conn)

	if session, _ := ReadRequestCookie(conn.Request()); session != nil {
		socket.Login(session)
	}

	socket.Watch()
})
