package websocket

import (
	"ztaylor.me/cast"
)

var pingTimeout = cast.Minute

var lonely = cast.BytesS(`{"uri":"/ping"}`)

// watch performs socket i/o and sends when it gets lonely
func watch(service Service, t *T) {
	for heat, pingTimer, resetCD := 0, cast.NewTimer(pingTimeout), cast.Now(); ; {
		select {
		case <-pingTimer.C:
			t.Send(lonely) // falls onto send chan
			pingTimer.Reset(pingTimeout)
		case buff := <-t.send:
			if heat > 0 {
				heat--
			}
			if now := cast.Now(); now.Sub(resetCD) > cast.Second {
				if !pingTimer.Stop() {
					<-pingTimer.C
				}
				pingTimer.Reset(pingTimeout)
				resetCD = now
			}
			if err := Codec.Send(t.conn, buff); err != nil {
				if !pingTimer.Stop() {
					<-pingTimer.C
				}
				go drainMessageChan(t.recv)
				t.Close()
				return
			}
		case msg := <-t.recv:
			if msg == nil {
				if !pingTimer.Stop() {
					<-pingTimer.C
				}
				t.Close()
				return
			}
			if now := cast.Now(); now.Sub(resetCD) > cast.Second {
				if !pingTimer.Stop() {
					<-pingTimer.C
				}
				pingTimer.Reset(pingTimeout)
				resetCD = now
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
