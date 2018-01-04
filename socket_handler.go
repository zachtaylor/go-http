package http

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/events"
	"ztaylor.me/log"
)

var SocketHandler = websocket.Handler(func(conn *websocket.Conn) {
	log.Add("RemoteAddr", conn.Request().RemoteAddr).Debug("http.sockets: open")
	socket := Open(conn)

	if session, _ := ReadRequestCookie(conn.Request()); session != nil {
		socket.Session = session
		log.Add("SessionId", session.Id).Add("Username", session.Username).Info("http.sockets: session match")
	}

	events.Fire("WebsocketOpen", socket)

	socket.Watch()

	events.Fire("WebsocketClose", socket)
})
