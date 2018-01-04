package http

import (
	"net/http"
	"ztaylor.me/json"
)

type NetHttpAgent struct {
	name string
	http.ResponseWriter
}

func AgentFromNetHttp(r *http.Request, w http.ResponseWriter) Agent {
	return &NetHttpAgent{r.UserAgent(), w}
}

func (a *NetHttpAgent) Name() string {
	return a.name
}

func (a *NetHttpAgent) Write(s string) {
	a.ResponseWriter.Write([]byte(s))
}

func (a *NetHttpAgent) WriteJson(json json.Json) {
	a.Write(json.String())
}
