package sessions

import (
	"time"
)

func Watch() {
	for now := range time.Tick(1 * time.Second) {
		revokelist := make([]uint, 0)

		for sessionId, session := range Cache {
			if session.Expire.Before(now) {
				revokelist = append(revokelist, sessionId)
			}
		}

		for _, sessionId := range revokelist {
			if session := Cache[sessionId]; session != nil {
				session.Revoke()
			}
		}
	}
}
