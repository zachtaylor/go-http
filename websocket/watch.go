package websocket

import "time"

// heatline tunes what watch considers excessive user input (starts at 0)
var heatline = 2

var pingTimeout = time.Minute

var lonely []byte = []byte(`{"uri":"/ping"}`)

// Watch performs socket i/o and sends when lonely
func Watch(service Service, t *T) {
	for heat, pingTimer, resetCD := 0, time.NewTimer(pingTimeout), time.Now(); ; {
		select {
		case <-pingTimer.C:
			t.Write(lonely) // falls onto send chan
			pingTimer.Reset(pingTimeout)
		case buff := <-t.send:
			if heat > 0 {
				heat--
			}
			// pingTimer.Reset on write effect has a 1 sec cooldown
			if now := time.Now(); now.Sub(resetCD) > time.Second {
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
			if now := time.Now(); now.Sub(resetCD) > time.Second {
				if !pingTimer.Stop() {
					<-pingTimer.C
				}
				pingTimer.Reset(pingTimeout)
				resetCD = now
			}
			if heat > heatline {
				<-time.After(time.Duration(heat-heatline) * 100 * time.Millisecond)
				heat--
			}
			heat++
			go service.ServeWS(t, msg)
		}
	}
}
