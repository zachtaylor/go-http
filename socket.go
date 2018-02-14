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
	events.Bus
}

type SocketSlice []*Socket

func Open(conn *websocket.Conn) *Socket {
	s := &Socket{
		name: "ws://" + conn.Request().RemoteAddr,
		conn: conn,
		Bus:  events.Bus{},
	}
	s.Bus.On(events.EVTfire, events.FireGlobal)
	return s
}

func (socket *Socket) Name() string {
	return socket.name
}

func (socket *Socket) Login(session *Session) {
	if socket.Session != nil {
		log.Add("SessionId", socket.Id).Add("Username", socket.Username).Warn("http/socket: login duplicated")
		return
	}

	socket.Session = session
	StoreSocket(socket)
	socket.Fire(EVTsocket_login, session, socket)
	log.Add("SessionId", session.Id).Add("Username", socket.Username).Info("http/socket: login")
}

func (socket *Socket) Write(s string) {
	socket.WriteJson(js.Object{
		"error": s,
	})
}

func (slice SocketSlice) WriteAll(s string) {
	for _, socket := range slice {
		socket.Write(s)
	}
}

func (socket *Socket) WriteJson(json js.Object) {
	if socket.conn != nil {
		websocket.Message.Send(socket.conn, json.String())
	}
}

func (slice SocketSlice) WriteAllJson(json js.Object) {
	for _, socket := range slice {
		socket.WriteJson(json)
	}
}

func (socket *Socket) Watch() {
	socket.Fire(EVTsocket_open, socket)
	for {
		req := <-socket.Listener()
		if req != nil {
			Dispatch(req)
			socket.Fire(EVTsocket_receive, req)
		} else {
			log.Add("Name", socket.name).Debug("http/socket: done")
			return
		}
	}
	RemoveSocket(socket)
	socket.Fire(EVTsocket_close, socket)
}

func (socket *Socket) Listener() chan *Request {
	receiver := make(chan *Request)
	go func() {
		s := ""
		msg := SocketMessage{"", js.Object{}}
		log := log.Add("Session", socket.Session)
		if socket == nil || socket.conn == nil {
			log.Warn("http/socket: listen to nil socket")
		} else if err := websocket.Message.Receive(socket.conn, &s); err != nil {
			if err.Error() != "EOF" {
				log.Add("Error", err).Error("http/socket: receive error")
			}
		} else if err := js.NewDecoder(bytes.NewBufferString(s)).Decode(&msg); err != nil {
			log.Add("Error", err).Add("Val", s).Error("http/socket: receive decode error")
		} else {
			log.Add("Uri", msg.Uri).Debug("http/socket: receive")
			receiver <- RequestFromSocketMessage(&msg, socket)
		}

		close(receiver)
	}()
	return receiver
}
