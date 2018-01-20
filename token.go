package http

import (
	"math/rand"
	"sync"
	"time"
)

var tokens = make(map[string]string)
var tokenRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var tokenRing = []int{
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
}
var tokenLock sync.Mutex
var TokenLength = 24

func CheckToken(v string) string {
	tokenLock.Lock()
	k := tokens[v]
	tokenLock.Unlock()
	return k
}

func RemoveToken(v string) {
	tokenLock.Lock()
	delete(tokens, v)
	tokenLock.Unlock()
}

func RegisterToken(k string) string {
	tokenLock.Lock()
	val := newTokenVal()
	for tokens[val] != "" {
		val = newTokenVal()
	}
	tokens[val] = k
	tokenLock.Unlock()
	return val
}

func newTokenVal() string {
	v := ""
	ringp := 0
	for i := 0; i < TokenLength; i++ {
		ringp = (ringp + (tokenRand.Int() % len(tokenRing))) % len(tokenRing)
		v = v + string(tokenRing[ringp])
	}
	return v
}
