package http

import (
	"bytes"

	"golang.org/x/net/websocket"
	"ztaylor.me/events"
	"ztaylor.me/js"
	"ztaylor.me/log"
)

const EVTsocket_open = "SocketOpen"
const EVTsocket_login = "SocketLogin"
const EVTsocket_receive = "SocketReceive"
const EVTsocket_close = "SocketClose"

type Socket struct {
	name string
	conn *websocket.Conn
	*Session
}

func Open(conn *websocket.Conn) *Socket {
	s := &Socket{
		name: "ws://" + conn.Request().RemoteAddr,
		conn: conn,
	}
	return s
}

func (socket *Socket) Name() string {
	return socket.String()
}

func (socket Socket) String() string {
	return socket.name
}

func (socket *Socket) Login(session *Session) {
	if socket.Session != nil {
		log.Add("Name", socket.Name()).Add("SessionID", socket.ID).Add("Username", socket.Username).Warn("http/socket: login duplicated")
		return
	}

	socket.Session = session
	StoreSocket(socket)
	events.Fire(EVTsocket_login, socket, session)
	log.Add("Name", socket.Name()).Add("SessionID", session.ID).Add("Username", socket.Username).Info("http/socket: login")
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
	events.Fire(EVTsocket_open, socket)
	for {
		req := <-socket.Listener()
		if req != nil {
			Dispatch(req)
			events.Fire(EVTsocket_receive, socket, req)
		} else {
			log.Add("Name", socket.Name()).Debug("http/socket: done")
			break
		}
	}
	RemoveSocket(socket)
	events.Fire(EVTsocket_close, socket)
}

func (socket *Socket) Listener() chan *Request {
	receiver := make(chan *Request)
	go func() {
		s := ""
		msg := SocketMessage{"", js.Object{}}
		if socket == nil || socket.conn == nil {
			log.Add("Name", socket.Name()).Warn("http/socket: listen socket is nil")
		} else if err := websocket.Message.Receive(socket.conn, &s); err != nil {
			if err.Error() != "EOF" {
				log.Add("Error", err).Error("http/socket: receive error")
			}
		} else if err := js.NewDecoder(bytes.NewBufferString(s)).Decode(&msg); err != nil {
			log.WithFields(log.Fields{
				"Name":    socket.Name(),
				"Session": socket.Session,
				"Val":     s,
				"Error":   err,
			}).Error("http/socket: receive decode error")
		} else {
			log.WithFields(log.Fields{
				"Name":    socket.Name(),
				"Session": socket.Session,
				"Uri":     msg.Uri,
			}).Debug("http/socket: receive")
			receiver <- RequestFromSocketMessage(&msg, socket)
		}

		close(receiver)
	}()
	return receiver
}
