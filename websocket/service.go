package websocket

// Service provides websocket server functionality
type Service interface {
	Handler
	// Count returns the number of open sockets
	Count() int
	// Get returns the socket for the given socket id
	Get(string) *T
	// Conn merges the current thread (go routine) to monitor the connection
	Conn(*Conn)
	// Broadcast sends a message to all connected sockets
	Broadcast(*Message)
}
