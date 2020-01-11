package websocket

import "ztaylor.me/cast"

// Message is websocket message data
type Message struct {
	URI  string
	Data cast.JSON
}

// NewMessage creates a Message
func NewMessage(uri string, json cast.JSON) *Message {
	return &Message{
		URI:  uri,
		Data: json,
	}
}

// JSON implements conversion to JSON data
func (m *Message) JSON() cast.JSON {
	if m == nil {
		return nil
	}
	return cast.JSON{
		"uri":  m.URI,
		"data": m.Data,
	}
}

func (m *Message) String() string {
	return "websocket.Message{URI:" + m.URI + "}"
}
