package websocket

import (
	"golang.org/x/net/websocket"
	"ztaylor.me/cast"
)

var pingTimeout = cast.Minute

var lonely = cast.BytesS(`{"uri":"/ping"}`)

// watch performs socket i/o and sends when it gets lonely
func watch(service Service, t *T) {
	for heat, pingTimer := 0, cast.NewTimer(pingTimeout); ; {
		select {
		case <-pingTimer.C:
			t.Send(lonely) // falls onto send chan
		case buff := <-t.send:
			if heat > 0 {
				heat--
			}
			if !pingTimer.Stop() {
				<-pingTimer.C
			}
			pingTimer.Reset(pingTimeout)
			if err := websocket.Message.Send(t.conn, buff); err != nil {
				go drainMessageChan(t.recv)
				t.Close()
				return
			}
		case msg := <-t.recv:
			if msg == nil {
				if !pingTimer.Stop() {
					<-pingTimer.C
				}
				pingTimer.Reset(pingTimeout)
				t.Close()
				return
			}
			if heat > heatline {
				<-cast.After(cast.Duration(heat-heatline) * 100 * cast.Millisecond)
				heat--
			}
			heat++
			go service.ServeWS(t, msg)
		}
	}
}
