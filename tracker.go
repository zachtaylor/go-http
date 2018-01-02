package http

import (
	"sync"
	"ztaylor.me/log"
)

var tracker = NewTracker()
var Track = tracker.Track
var GetTrackRecord = tracker.Record
var TrackPair = tracker.Pair

type Tracker struct {
	sync.Mutex
	data map[string]string
}

func NewTracker() *Tracker {
	return &Tracker{sync.Mutex{}, make(map[string]string)}
}

func (t *Tracker) Track(addr string) {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.data[addr]; !ok {
		t.data[addr] = ""
		log.Add("Addr", addr).Debug("http.track: address recorded")
	}
}

func (t *Tracker) Record(name string) []string {
	t.Lock()
	defer t.Unlock()

	s := make([]string, 0)
	for k, v := range t.data {
		if name == v {
			s = append(s, k)
		}
	}

	return s
}

func (t *Tracker) Pair(name string, addr string) {
	t.Lock()
	defer t.Unlock()
	if t.data[addr] == "" || t.data[name] == "" || t.data[name] != addr {
		t.data[name] = addr
		t.data[addr] = name
		log.Add("Addr", addr).Add("Username", name).Debug("http.track: account address paired")
	}
}
