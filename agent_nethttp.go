package http

import (
	"net/http"
	"ztaylor.me/js"
)

type NetHttpAgent struct {
	name string
	http.ResponseWriter
}

func AgentFromNetHttp(r *http.Request, w http.ResponseWriter) Agent {
	return &NetHttpAgent{
		name:           r.Method + "://" + r.RemoteAddr,
		ResponseWriter: w,
	}
}

func (a *NetHttpAgent) Name() string {
	return a.name
}

func (a *NetHttpAgent) Write(s string) {
	a.ResponseWriter.Write([]byte(s))
}

func (a *NetHttpAgent) WriteJson(json js.Object) {
	a.Write(json.String())
}
