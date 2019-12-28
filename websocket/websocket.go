package websocket // import "ztaylor.me/http/websocket"

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/cast"
	"ztaylor.me/http/session"
)

// T is a websocket connection
type T struct {
	ID      string
	Session *session.T
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

// Message is a macro for WriteMessage(NewMessage)
func (t *T) Message(uri string, json cast.JSON) {
	t.WriteMessage(NewMessage(uri, json))
}

// WriteMessage calls Write with cast []byte Message.JSON().String()
func (t *T) WriteMessage(m *Message) {
	t.Write(cast.BytesS(m.JSON().String()))
}

// Write sends a buffer to the underlying websocket connection
func (t *T) Write(s []byte) {
	if conn := t.conn; conn != nil {
		go websocket.Message.Send(conn, s)
	}
}

// NextMessage reads a Message from the socket API
func (t *T) NextMessage() (*Message, error) {
	s, msg := "", &Message{}
	if t == nil || t.conn == nil {
		return nil, cast.EOF
	} else if err := websocket.Message.Receive(t.conn, &s); err != nil {
		return nil, err
	} else if err := cast.DecodeJSON(cast.NewBuffer(s), msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// NextChan creates a chan of *Message using NextMessage
func (t *T) NextChan() chan *Message {
	msgs := make(chan *Message)
	go func() {
		for { // loop
			if msg, err := t.NextMessage(); err == nil {
				msgs <- msg
			} else if err == cast.EOF {
				t.Close()
				break
			}
		} // loop
		close(msgs)
	}()
	return msgs
}

func (t *T) String() string {
	return "websocket.T{" + t.conn.Request().RemoteAddr + "}"
}
