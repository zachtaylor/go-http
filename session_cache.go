package http

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
	"ztaylor.me/events"
)

var sessionCache = make(map[uint]*Session)
var sessionCacheLock sync.Mutex

func init() {
	go watch()
}

func SessionCount() int {
	return len(sessionCache)
}

func GetSession(username string) *Session {
	sessionCacheLock.Lock()
	for _, session := range sessionCache {
		if username == session.Username {
			sessionCacheLock.Unlock()
			return session
		}
	}
	sessionCacheLock.Unlock()
	return nil
}

func GrantSession(username string) *Session {
	session := NewSession()
	session.Username = username
	sessionCache[session.Id] = session
	events.Fire("SessionGrant", session)
	return session
}

func RevokeSession(username string) {
	if session := GetSession(username); session != nil {
		sessionCacheLock.Lock()
		delete(sessionCache, session.Id)
		sessionCacheLock.Unlock()
		if session.Expire.After(time.Now()) {
			session.Close()
		}
	}
}

func ReadRequestCookie(r *http.Request) (*Session, error) {
	if sessionCookie, err := r.Cookie("SessionId"); err == nil {
		if sessionId, err := strconv.ParseInt(sessionCookie.Value, 10, 0); err == nil {
			if session := sessionCache[uint(sessionId)]; session != nil {
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
		revokelist := make([]string, 0)

		for _, session := range sessionCache {
			if session.Expire.Before(now) {
				revokelist = append(revokelist, session.Username)
			}
		}

		for _, username := range revokelist {
			RevokeSession(username)
		}
	}
}
