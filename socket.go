package http

import (
	"bytes"
	"golang.org/x/net/websocket"
	"ztaylor.me/events"
	"ztaylor.me/json"
	"ztaylor.me/log"
)

type Socket struct {
	name string
	conn *websocket.Conn
	*Session
}

func Open(conn *websocket.Conn) *Socket {
	return &Socket{
		name: "ws://" + conn.Request().RemoteAddr,
		conn: conn,
	}
}

func (socket *Socket) Name() string {
	return socket.name
}

func (socket *Socket) Write(s string) {
	socket.WriteJson(json.Json{
		"error": s,
	})
}

func (socket *Socket) WriteJson(json json.Json) {
	if socket.conn != nil {
		websocket.Message.Send(socket.conn, json.String())
	}
}

func (socket *Socket) Watch() {
	for {
		req := <-socket.Listen()
		if req != nil {
			Dispatch(req)
			events.Fire("WebsocketRequest", req)
		} else {
			log.Add("Name", socket.name).Debug("http/socket: done")
			return
		}
	}
}

func (socket *Socket) Listen() chan *Request {
	receiver := make(chan *Request)
	go func() {
		s := ""
		msg := SocketMessage{"", json.Json{}}
		log := log.Add("Session", socket.Session)
		if socket == nil {
			log.Warn("listen to nil socket")
		} else if err := websocket.Message.Receive(socket.conn, &s); err != nil {
			if err.Error() != "EOF" {
				log.Add("Error", err).Error("http/socket: receive error")
			}
			receiver <- nil
		} else if err := json.NewDecoder(bytes.NewBufferString(s)).Decode(&msg); err != nil {
			log.Add("Error", err).Add("Val", s).Error("http/socket: receive decode error")
			receiver <- nil
		} else {
			log.Add("Uri", msg.Uri).Debug("http/socket: receive")
			receiver <- RequestFromSocketMessage(&msg, socket)
		}

		close(receiver)
	}()
	return receiver
}
