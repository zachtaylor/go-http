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
	done chan interface{}
	*Session
}

func Open(conn *websocket.Conn) *Socket {
	return &Socket{
		name: "ws://" + conn.Request().RemoteAddr,
		conn: conn,
		done: make(chan interface{}),
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
	for socket.conn != nil {
		select {
		case <-socket.done:
			log.Debug("socket done")
			socket.conn = nil
		case req := <-socket.Listen():
			if req != nil {
				go Dispatch(req)
				events.Fire("WebsocketRequest", req)
			} else {
				close(socket.done)
			}
		}
	}
}

func (socket *Socket) Listen() chan *Request {
	receiver := make(chan *Request)
	go func() {
		s := bytes.NewBufferString("")
		msg := SocketMessage{"", json.Json{}}
		if err := websocket.Message.Receive(socket.conn, &s); err != nil {
			if err.Error() != "EOF" {
				log.Add("Error", err).Error("socket receive")
			}
			receiver <- nil
		} else if err := json.NewDecoder(s).Decode(&msg); err != nil {
			log.Add("Error", err).Add("Val", s).Error("socket receive decode")
			receiver <- nil
		} else {
			log.Add("Uri", msg.Uri).Add("Username", socket.Username).Debug("socket receive")
			receiver <- RequestFromSocketMessage(&msg, socket)
		}

		close(receiver)
	}()
	return receiver
}
