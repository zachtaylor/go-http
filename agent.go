package http

import (
	"ztaylor.me/js"
)

type Agent interface {
	Name() string
	Write(string)
	WriteJson(js.Object)
}
