package http

import (
	"io"
	"ztaylor.me/json"
)

type WriterAgent struct {
	w    io.Writer
	mdat json.Json
}

func AgentFromWriter(w io.Writer) Agent {
	return &WriterAgent{
		w:    w,
		mdat: json.Json{},
	}
}

func (a *WriterAgent) Metadata() json.Json {
	return a.mdat
}

func (a *WriterAgent) Write(s string) {
	a.w.Write([]byte(s))
}

func (a *WriterAgent) WriteJson(json json.Json) {
	a.Write(json.String())
}
