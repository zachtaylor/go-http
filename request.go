package http

import (
	"net/http"
	"ztaylor.me/json"
)

type Request struct {
	Quest string
	Agent
	Data json.Json
}

func NewRequest() *Request {
	return &Request{
		Data: json.Json{},
	}
}

func RequestFromNet(r *http.Request, w http.ResponseWriter) *Request {
	req := NewRequest()
	req.Quest = r.RequestURI
	if session, _ := ReadRequestCookie(r); session != nil {
		req.Agent = AgentFromWriter(w)
		req.Metadata()["Username"] = session.Username
		req.Metadata()["SessionId"] = session.Id
	}
	return req
}

func RequestFromSocketMessage(msg *SocketMessage, s *Socket) *Request {
	return &Request{msg.Uri, s, msg.Data}
}
