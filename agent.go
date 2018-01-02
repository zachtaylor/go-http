package http

import (
	"ztaylor.me/json"
)

type Agent interface {
	Metadata() json.Json
	Write(string)
	WriteJson(json.Json)
}
