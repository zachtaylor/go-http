package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"ztaylor.me/events"
)

var SessionCache = make(map[uint]*Session)

func init() {
	go watch()
}

func Get(username string) *Session {
	for _, session := range SessionCache {
		if username == session.Username {
			return session
		}
	}
	return nil
}

func Grant(username string, lifetime time.Duration) *Session {
	session := NewSession(lifetime)
	session.Username = username
	SessionCache[session.Id] = session
	events.Fire("SessionGrant", session)
	return session
}

func Revoke(username string) {
	if session := Get(username); session != nil {
		session.Revoke()
	}
}

func ReadRequestCookie(r *http.Request) (*Session, error) {
	if sessionCookie, err := r.Cookie("SessionId"); err == nil {
		if sessionId, err := strconv.ParseInt(sessionCookie.Value, 10, 0); err == nil {
			if session := SessionCache[uint(sessionId)]; session != nil {
				return session, nil
			} else if sessionId == 0 {
				return nil, nil
			} else {
				return nil, errors.New("invalid cookie")
			}
		} else {
			return nil, errors.New("cookie format")
		}
	} else {
		return nil, errors.New("session missing")
	}
}

func EraseSessionId(w http.ResponseWriter) {
	w.Header().Set("Set-Cookie", "SessionId=0; Path=/;")
}

func watch() {
	for now := range time.Tick(1 * time.Second) {
		revokelist := make([]uint, 0)

		for sessionId, session := range SessionCache {
			if session.Expire.Before(now) {
				revokelist = append(revokelist, sessionId)
			}
		}

		for _, sessionId := range revokelist {
			if session := SessionCache[sessionId]; session != nil {
				session.Revoke()
			}
		}
	}
}
