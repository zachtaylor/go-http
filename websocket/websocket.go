package websocket // import "ztaylor.me/http/websocket"

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
	"ztaylor.me/http/json"
	"ztaylor.me/http/sessions"
)

// Service provides websocket server functionality
type Service interface {
	Handler
	// New connects a websocket
	Connect(*websocket.Conn)
	// Count returns the number of open sockets
	Count() int
}

// UpgradeHandler provides a websocket handshake func
func UpgradeHandler(s Service) http.Handler {
	return websocket.Handler(func(c *websocket.Conn) {
		s.Connect(c)
	})
}

// T is a websocket connection
type T struct {
	Session *sessions.T
	conn    *websocket.Conn
	done    chan bool
}

// New creates an initialied orphan websocket
func New(conn *websocket.Conn) *T {
	return &T{
		conn: conn,
		done: make(chan bool),
	}
}

// Key is RemoteAddr
//
// Since websockets block remote ports, it is guaranteed to be unique
func (t *T) Key() string {
	return t.conn.Request().RemoteAddr
}

// GetUser returns a non-empty string value to represent this websocket
func (t *T) GetUser() string {
	if t != nil && t.Session != nil {
		return t.Session.Name()
	} else if i := strings.LastIndex(t.conn.Request().RemoteAddr, ":"); i < 0 {
		return "anon"
	} else {
		return "anon#" + t.conn.Request().RemoteAddr[i+1:]
	}
}

// Done is an observable channel that closes when the socket has been closed
func (t *T) Done() <-chan bool {
	return t.done
}

// Close closes the observable channel
func (t *T) Close() {
	if t.done != nil {
		close(t.done)
		t.done = nil
	}
}

// Write sends a buffer to the underlying websocket connection
func (t *T) Write(s []byte) {
	if t.conn != nil {
		websocket.Message.Send(t.conn, s)
	}
}

var errSocketClosed = errors.New("websocket closed")

// NextMessage reads a Message from the socket API
func (t *T) NextMessage() (*Message, error) {
	s, msg := "", &Message{}
	if t == nil || t.conn == nil {
		return nil, errSocketClosed
	} else if err := websocket.Message.Receive(t.conn, &s); err != nil {
		return nil, err
	} else if err := json.Decode(bytes.NewBufferString(s), msg); err != nil {
		return nil, err
	}
	msg.User = t.GetUser()
	return msg, nil
}

func (t *T) String() string {
	return "ws(" + t.GetUser() + ")://" + t.conn.Request().RemoteAddr
}
