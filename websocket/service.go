package websocket

import "golang.org/x/net/websocket"

// Service provides websocket server functionality
type Service interface {
	Handler
	// New connects a websocket
	Connect(*websocket.Conn)
	// Count returns the number of open sockets
	Count() int
}
