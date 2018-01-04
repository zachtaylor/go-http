package http

import (
	"net/http"
	"strings"
	"ztaylor.me/json"
)

type Request struct {
	Quest    string
	Remote   string
	Language string
	Data     json.Json
	*Session
	Agent
}

func NewRequest() *Request {
	return &Request{
		Data: json.Json{},
	}
}

func RequestFromNet(r *http.Request, w http.ResponseWriter) *Request {
	req := NewRequest()
	req.Quest = r.RequestURI
	req.Remote = r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")]
	for k, v := range r.Form {
		req.Data[k] = v
	}
	req.Agent = AgentFromNetHttp(r, w)
	if session, _ := ReadRequestCookie(r); session != nil {
		req.Session = session
	}
	return req
}

func RequestFromSocketMessage(msg *SocketMessage, s *Socket) *Request {
	req := NewRequest()
	req.Quest = msg.Uri
	req.Remote = s.Name()
	req.Data = msg.Data
	req.Agent = s
	req.Session = s.Session
	return req
}
