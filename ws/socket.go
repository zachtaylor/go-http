package ws

import (
	"bytes"
	"strings"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/sessions"
	"ztaylor.me/js"
	"ztaylor.me/log"
)

// Socket is a websocket connection
type Socket struct {
	Session *sessions.T
	conn    *websocket.Conn
	done    chan bool
}

// NewSocket creates a Socket
func NewSocket(conn *websocket.Conn) *Socket {
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
		return socket.Session.Name()
	} else if i := strings.LastIndex(socket.conn.Request().RemoteAddr, ":"); i < 0 {
		return "anon"
	} else {
		return "anon#" + socket.conn.Request().RemoteAddr[i+1:]
	}
}

func (socket Socket) String() string {
	return "ws(" + socket.GetUser() + ")://" + socket.conn.Request().RemoteAddr
}

// Done is an observable channel that indicates the socket has been closed
func (socket *Socket) Done() <-chan bool {
	return socket.done
}

// Login saves this socket in the Service
func (socket *Socket) Login(service *Service, session *sessions.T) {
	if socket.Session != nil {
		service.Remove(socket.String())
	}
	socket.Session = session
	service.Store(socket.String(), socket)
}

// Write sends a buffer to the underlying websocket connection
func (socket *Socket) Write(s []byte) {
	if socket.conn != nil {
		websocket.Message.Send(socket.conn, s)
	}
}

// WriteJson uses Socket.Write and js.Object.String
func (socket *Socket) WriteJson(json js.Object) {
	socket.Write([]byte(json.String()))
}

func (socket *Socket) watch(service *Service, h Handler) {
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
		case msg, ok := <-receiver:
			if ok {
				h.ServeWS(socket, msg)
			} else {
				log.Protect(func() {
					close(socket.done)
				})
			}
		case _, ok := <-socket.done:
			if !ok {
				service.Remove(socket.String())
				return
			}
		}
	}
}

// Listen creates the next Message from the Socket by waiting forever
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
