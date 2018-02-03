package http

import (
	"fmt"
	"github.com/cznic/mathutil"
	"net/http"
	"strconv"
	"sync"
	"time"
	"ztaylor.me/events"
)

var sessionIdGen, _ = mathutil.NewFC32(0, 999999999, true)
var SessionLifetime = 1 * time.Hour

type Session struct {
	Id       uint
	Username string
	Expire   time.Time
	Done     chan error
	sync.Mutex
}

func NewSession() *Session {
	return &Session{
		Id:     uint(sessionIdGen.Next()),
		Expire: time.Now().Add(SessionLifetime),
		Done:   make(chan error),
	}
}

func (session *Session) Refresh() {
	session.Expire = time.Now().Add(SessionLifetime)
}

func (session *Session) Close() {
	go func() {
		session.Expire = time.Now()
		close(session.Done)
		events.Fire("SessionClose", session.Username)
	}()
}

func (session *Session) WriteCookie(w http.ResponseWriter) {
	w.Header().Set("Set-Cookie", "SessionId="+strconv.Itoa(int(session.Id))+"; Path=/;")
}

func (session Session) String() string {
	return fmt.Sprintf("#%d:%s", session.Id, session.Username)
}
