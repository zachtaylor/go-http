package ws

import (
	"bytes"

	"golang.org/x/net/websocket"
	"ztaylor.me/events"
	"ztaylor.me/http/sessions"
	"ztaylor.me/js"
	"ztaylor.me/log"
)

const EVTopen = "http/mux/ws.Open"
const EVTlogin = "http/mux/ws.Login"
const EVTreceive = "http/mux/ws.Receive"
const EVTclose = "http/mux/ws.Close"

type Socket struct {
	conn    *websocket.Conn
	session *sessions.Grant
}

func Open(conn *websocket.Conn) *Socket {
	s := &Socket{
		conn: conn,
	}
	return s
}

func (socket *Socket) User() string {
	if socket == nil || socket.session == nil {
		return ""
	}
	return socket.session.Name
}

func (socket Socket) String() string {
	return "ws(" + socket.User() + ")://" + socket.conn.Request().RemoteAddr
}

func (socket *Socket) Login(session *sessions.Grant) {
	if socket.session != nil {
		log.Add("Socket", socket).Add("Session", session).Warn("http/socket: login duplicated")
		return
	}

	socket.session = session
	Service.Store(socket)
	events.Fire(EVTlogin, socket, session)
	log.Add("Socket", socket).Info("http/socket: login")
}

func (socket *Socket) Write(s []byte) {
	if socket.conn != nil {
		websocket.Message.Send(socket.conn, s)
	}
}

func (socket *Socket) WriteJson(json js.Object) {
	socket.Write([]byte(json.String()))
}

func (socket *Socket) Watch() {
	events.Fire(EVTopen, socket)
	receiver := make(chan *Message)
	for {
		go socket.Listen(receiver)
		msg := <-receiver
		if msg != nil {
			Service.Dispatch(msg)
			events.Fire(EVTreceive, socket, msg)
		} else {
			log.Add("Socket", socket).Debug("http/socket: done")
			break
		}
	}
	Service.Remove(socket.String())
	events.Fire(EVTclose, socket)
}

func (socket *Socket) Listen(receiver chan *Message) {
	s := ""
	msg := &Message{
		User: socket.User(),
		Data: js.Object{},
	}
	if socket == nil || socket.conn == nil {
		log.Add("Socket", socket).Warn("http/socket: listen socket is nil")
	} else if err := websocket.Message.Receive(socket.conn, &s); err != nil {
		if err.Error() != "EOF" {
			log.Add("Error", err).Error("http/socket: receive error")
			close(receiver)
		}
	} else if err := js.NewDecoder(bytes.NewBufferString(s)).Decode(&msg); err != nil {
		log.WithFields(log.Fields{
			"Socket": socket,
			"Val":    s,
			"Error":  err,
		}).Error("http/socket: receive decode error")
	} else {
		log.WithFields(log.Fields{
			"Socket":  socket,
			"Message": msg.Name,
		}).Debug("http/socket: receive")
		receiver <- msg
	}
}
