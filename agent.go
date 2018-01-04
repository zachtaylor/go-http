package http

import (
	"ztaylor.me/json"
)

type Agent interface {
	Name() string
	Write(string)
	WriteJson(json.Json)
}
