package wsocks

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/events"
	"ztaylor.me/json"
	"ztaylor.me/log"
)

type Socket struct {
	SessionId uint
	Done      chan interface{}
	conn      *websocket.Conn
	events.Bus
}

func Open(conn *websocket.Conn) *Socket {
	return &Socket{
		Done: make(chan interface{}),
		conn: conn,
		Bus:  events.Bus{},
	}
}

func (socket *Socket) Send(name string, data json.Json) {
	if socket.conn == nil {
		log.Add("SessionId", socket.SessionId).Add("Name", name).Warn("wsocks: conn is nil")
		return
	}
	websocket.Message.Send(socket.conn, json.Json{
		"name": name,
		"data": data,
	}.String())
}

func (socket *Socket) Watch() {
	for socket.conn != nil {
		msgIn, msgErr := socket.receivers()
		select {
		case <-socket.Done:
			socket.conn = nil
			socket.Fire("Done")
		case err := <-msgErr:
			if err.Error() != "EOF" {
				log.Add("Error", err).Add("SessionId", socket.SessionId).Error("wsocks: receive")
			}
			close(socket.Done)
		case msg := <-msgIn:
			if msg != nil {
				socket.Fire("Message", msg)
			} else {
				close(socket.Done)
			}
		}
	}
}

func (socket *Socket) receivers() (chan *Message, chan error) {
	receiver := make(chan *Message)
	errors := make(chan error)
	go func() {
		msg := NewMessage()
		err := websocket.JSON.Receive(socket.conn, msg)
		if err != nil {
			errors <- err
		} else {
			receiver <- msg
		}
		close(receiver)
		close(errors)
	}()
	return receiver, errors
}
