package ws

import (
	"bytes"
	"strings"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/sessions"
	"ztaylor.me/js"
	"ztaylor.me/log"
)

type Socket struct {
	Session *sessions.Grant
	conn    *websocket.Conn
	done    chan bool
}

func Open(conn *websocket.Conn) *Socket {
	s := &Socket{
		conn: conn,
		done: make(chan bool),
	}
	return s
}

// GetUser returns a nonempty string value to represent this socket
//
// If the underlying Session is not available, GetUser returns a best-guess "anon" name
func (socket *Socket) GetUser() string {
	if socket != nil && socket.Session != nil {
		return socket.Session.Name
	} else if i := strings.LastIndex(socket.conn.Request().RemoteAddr, ":"); i < 0 {
		return "anon"
	} else {
		return "anon#" + socket.conn.Request().RemoteAddr[i+1:]
	}
}

func (socket Socket) String() string {
	return "ws(" + socket.GetUser() + ")://" + socket.conn.Request().RemoteAddr
}

func (socket *Socket) Done() chan bool {
	done := make(chan bool)
	go func() {
		defer close(done)

		if socket.done == nil {
			return
		}
		<-socket.done
	}()
	return done
}

func (socket *Socket) Login(session *sessions.Grant) {
	if socket.Session != nil {
		log.Add("Socket", socket).Add("Session", session).Warn("http/ws: login duplicated")
		return
	}

	socket.Session = session
	Service.Store(socket)
	log.Add("Socket", socket).Info("http/ws: login")
}

func (socket *Socket) Write(s []byte) {
	if socket.conn != nil {
		websocket.Message.Send(socket.conn, s)
	}
}

func (socket *Socket) WriteJson(json js.Object) {
	socket.Write([]byte(json.String()))
}

func (socket *Socket) Watch(h Handler) {
	receiver := make(chan *Message, 1)
	go log.Protect(func() {
		for {
			msg := socket.Listen()
			if msg == nil {
				close(receiver)
				return
			}
			receiver <- msg
		}
	})

	for {
		select {
		case msg := <-receiver:
			if msg != nil {
				h.ServeWS(socket, msg)
			}
		case _, ok := <-socket.done:
			if !ok {
				log.Add("Socket", socket).Info("http/ws: closed")
				Service.Remove(socket.String())
				return
			}
		}
	}
}

func (socket *Socket) Listen() *Message {
	s := ""
	msg := &Message{}
	if socket == nil || socket.conn == nil {
		log.Add("Socket", socket).Warn("http/ws: listen socket is nil")
		return nil
	} else if err := websocket.Message.Receive(socket.conn, &s); err != nil {
		if err.Error() != "EOF" {
			log.Add("Error", err).Error("http/ws: receive error")
		}
		return nil
	} else if err := js.NewDecoder(bytes.NewBufferString(s)).Decode(&msg); err != nil {
		log.WithFields(log.Fields{
			"Socket": socket,
			"Val":    s,
			"Error":  err,
		}).Error("http/ws: receive decode error")
		return nil
	}
	msg.User = socket.GetUser()
	log.WithFields(log.Fields{
		"Socket":  socket,
		"Message": msg,
	}).Debug("http/ws: receive")
	return msg
}
