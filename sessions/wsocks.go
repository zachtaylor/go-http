package sessions

import (
	"net/http"
	"ztaylor.me/events"
	"ztaylor.me/http/wsocks"
	"ztaylor.me/log"
)

func init() {
	events.On("WebsocketOpen", func(args ...interface{}) {
		s := args[0].(*wsocks.Socket)
		r := args[1].(*http.Request)
		if session, _ := ReadRequestCookie(r); session != nil {
			log.Add("SessionId", session.Id).Add("Username", session.Username).Info("sessions: socket match")
			s.SessionId = session.Id
		}
	})
}
