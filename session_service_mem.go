package http

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"ztaylor.me/events"
	"ztaylor.me/log"
)

func init() {
	SessionService = NewMemSessionService()
}

type MemSessionService struct {
	Sessions map[uint]*Session
	sync.Mutex
}

func NewMemSessionService() *MemSessionService {
	s := &MemSessionService{
		Sessions: make(map[uint]*Session),
	}
	go s.watch()
	return s
}

func (mem *MemSessionService) Count() int {
	return len(mem.Sessions)
}

func (mem *MemSessionService) Get(id uint) *Session {
	return mem.Sessions[id]
}

func (mem *MemSessionService) Find(username string) []*Session {
	sessions := make([]*Session, 0)
	mem.Lock()
	for _, session := range mem.Sessions {
		if username == session.Username {
			sessions = append(sessions, session)
		}
	}
	mem.Unlock()
	return sessions
}

func (mem *MemSessionService) Grant(username string) *Session {
	session := &Session{
		ID:       NewSessionID(),
		Username: username,
		Expire:   time.Now().Add(SessionLifetime),
		Done:     make(chan error),
	}
	mem.Lock()
	mem.Sessions[session.ID] = session
	mem.Unlock()
	log.Add("Session", session).Info("http/session_service_mem: grant")
	events.Fire("SessionGrant", session)
	return session
}

func (mem *MemSessionService) Revoke(id uint) {
	mem.Lock()
	if session := mem.Sessions[id]; session != nil {
		log.Add("SessionID", session).Info("http/session_service_mem: revoke")
		session.Close()
	}
	delete(mem.Sessions, id)
	mem.Unlock()
}

func ReadRequestCookie(r *http.Request) (*Session, error) {
	if sessionCookie, err := r.Cookie("SessionID"); err == nil {
		if sessionID, err := strconv.ParseUint(sessionCookie.Value, 10, 0); err == nil {
			if session := SessionService.Get(uint(sessionID)); session != nil {
				return session, nil
			} else if sessionID == 0 {
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

func EraseSessionID(w http.ResponseWriter) {
	w.Header().Set("Set-Cookie", "SessionID=0; Path=/;")
}

func (mem *MemSessionService) watch() {
	for now := range time.Tick(1 * time.Second) {
		revokelist := make([]uint, 0)

		for _, session := range mem.Sessions {
			if session.Expire.Before(now) {
				revokelist = append(revokelist, session.ID)
			}
		}

		for _, sessionID := range revokelist {
			mem.Revoke(sessionID)
		}
	}
}
