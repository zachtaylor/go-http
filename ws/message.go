package ws

import (
	"fmt"

	"ztaylor.me/js"
)

// Message is websocket message data
type Message struct {
	URI  string
	User string
	Data js.Object
}

func (m Message) String() string {
	return fmt.Sprintf("ws.Message{URI:\"%s\",User:\"%s\"}", m.URI, m.User)
}
