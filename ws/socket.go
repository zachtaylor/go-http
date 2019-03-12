package ws

import (
	"bytes"
	"errors"
	"strings"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/sessions"
	"ztaylor.me/js"
	"ztaylor.me/log"
)

var errSocketClosed = errors.New("socket closed")

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
	for {
		msg, err := socket.nextMessage()
		if err != nil {
			if err.Error() != "EOF" {
				log.Add("Error", err).Error("http/ws: receive error")
			}
			service.Remove(socket.String())
			close(socket.done)
			return
		}
		h.ServeWS(socket, msg)
	}
}

func (socket *Socket) nextMessage() (*Message, error) {
	s := ""
	msg := &Message{}
	if socket == nil || socket.conn == nil {
		return nil, errSocketClosed
	} else if err := websocket.Message.Receive(socket.conn, &s); err != nil {
		return nil, err
	} else if err := js.NewDecoder(bytes.NewBufferString(s)).Decode(&msg); err != nil {
		return nil, err
	}
	msg.User = socket.GetUser()
	return msg, nil
}
