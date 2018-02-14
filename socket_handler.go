package http

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/log"
)

var SocketHandler = websocket.Handler(func(conn *websocket.Conn) {
	log.Add("RemoteAddr", conn.Request().RemoteAddr).Debug("http/socket: open")

	socket := Open(conn)

	if session, _ := ReadRequestCookie(conn.Request()); session != nil {
		socket.Login(session)
	}

	socket.Watch()
})
