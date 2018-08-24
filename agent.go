package http

import "ztaylor.me/js"

type Agent interface {
	Name() string
	Write([]byte)
	WriteJson(js.Object)
}
