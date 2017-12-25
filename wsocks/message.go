package wsocks

import (
	"ztaylor.me/json"
)

type Message struct {
	Username string
	Data     json.Json
}

func NewMessage() *Message {
	return &Message{
		Data: json.Json{},
	}
}
