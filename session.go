package http

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"ztaylor.me/events"
	"ztaylor.me/log"
)

var SessionLifetime = 1 * time.Hour

type Session struct {
	Id       uint
	Username string
	Expire   time.Time
	Done     chan error
	sync.Mutex
}

func (session *Session) Refresh() {
	session.Expire = time.Now().Add(SessionLifetime)
	log.Add("Session", session).Debug("http/session: refresh")
}

func (session *Session) Close() {
	session.Expire = time.Now()
	close(session.Done)
	log.Add("Session", session).Debug("http/session: close")
	events.Fire("SessionClose", session.Username)
}

func (session *Session) WriteCookie(w http.ResponseWriter) {
	w.Header().Set("Set-Cookie", "SessionId="+strconv.Itoa(int(session.Id))+"; Path=/;")
}

func (session Session) String() string {
	return fmt.Sprintf("#%d:%s", session.Id, session.Username)
}
