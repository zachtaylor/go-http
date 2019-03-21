package websocket

import "ztaylor.me/cast"

// Message is websocket message data
type Message struct {
	URI  string
	User string
	Data cast.Fields
}

func (m Message) String() string {
	return "websocket.Message{URI:" + m.URI + " User:" + m.User + " Fields:" + cast.StringI(len(m.Data)) + "}"
}
