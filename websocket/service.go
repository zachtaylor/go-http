package websocket

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/cast"
)

// Service provides websocket server functionality
type Service interface {
	Handler
	// New connects a websocket
	Connect(*websocket.Conn)
	// Count returns the number of open sockets
	Count() int
	// Message is a macro for SendMessage(NewMessage)
	Message(string, cast.JSON)
	// SendMessage calls Write with cast []byte Message.JSON().String()
	SendMessage(*Message)
	// Send sends a buffer to all websocket connections
	Send(s []byte)
}
