package wsocks

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/events"
	"ztaylor.me/log"
)

var Handler = websocket.Handler(func(conn *websocket.Conn) {
	log.Add("RemoteAddr", conn.Request().RemoteAddr).Debug("wsocks: open")
	socket := Open(conn)
	events.Fire("WebsocketOpen", socket, conn.Request())
	socket.Watch()
	events.Fire("WebsocketClose", socket)
})
