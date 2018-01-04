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
