package http

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/json"
	"ztaylor.me/log"
)

type Socket struct {
	conn *websocket.Conn
	mdat json.Json
	done chan interface{}
}

func Open(conn *websocket.Conn) *Socket {
	return &Socket{
		conn: conn,
		mdat: json.Json{},
		done: make(chan interface{}),
	}
}

func (socket *Socket) Metadata() json.Json {
	return socket.mdat
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
			log.Add("", "").Debug("sockets: done")
			socket.conn = nil
		case msg := <-socket.Listen():
			if msg != nil {
				go Dispatch(socket, msg)
			} else {
				close(socket.done)
			}
		}
	}
}

func (socket *Socket) Listen() chan *Request {
	receiver := make(chan *Request)
	go func() {
		msg := NewSocketMessage()
		err := websocket.JSON.Receive(socket.conn, msg)
		if err != nil {
			receiver <- nil
			if err.Error() != "EOF" {
				log.Add("Error", err).Error("sockets: receive")
			}
		} else {
			log.Debug("sockets: receive")
			receiver <- RequestFromSocketMessage(msg, socket)
		}
		close(receiver)
	}()
	return receiver
}
