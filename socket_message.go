package http

import (
	"ztaylor.me/json"
)

type SocketMessage struct {
	Uri  string
	Data json.Json
}

func NewSocketMessage() *SocketMessage {
	return &SocketMessage{"", json.Json{}}
}

func (msg *SocketMessage) Request(user string) *Request {
	return &Request{
		Quest: msg.Uri,
		Data:  msg.Data,
	}
}
