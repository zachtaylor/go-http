package websocket

import "ztaylor.me/http/session"

// Runtime implements Service
type Runtime struct {
	Mux      *Mux
	Cache    *Cache
	Sessions session.Service
}

// NewRuntime creates a Runtime, which implements Service
func NewRuntime(mux *Mux, cache *Cache, sessions session.Service) *Runtime {
	return &Runtime{
		Mux:      mux,
		Cache:    cache,
		Sessions: sessions,
	}
}

// ServeWS calls the mux ServeWS
func (rt *Runtime) ServeWS(t *T, m *Message) {
	rt.Mux.ServeWS(t, m)
}

// Count returns the number of open websocket connections
func (rt *Runtime) Count() int {
	return rt.Cache.Count()
}

// Get returns the socket for the given socket id
func (rt *Runtime) Get(id string) *T {
	return rt.Cache.Get(id)
}

// Conn merges the current thread (go routine) to monitor the connection
func (rt *Runtime) Conn(conn *Conn) {
	t := rt.Cache.New(conn)
	t.Session, _ = rt.Sessions.Cookie(conn.Request())
	go Watch(rt, t)
	rt.ServeWS(t, NewMessage("/connect", nil))
	<-t.DoneChan() // await socket close
	rt.ServeWS(t, NewMessage("/disconnect", nil))
	t.Session = nil
	rt.Cache.Remove(t)
}

// Broadcast sends a message to all connected sockets
func (rt *Runtime) Broadcast(m *Message) {
	keys := rt.Cache.Keys()
	for _, k := range keys {
		if socket := rt.Get(k); socket != nil {
			socket.Message(m)
		}
	}
}
