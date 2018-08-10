package http

import "github.com/cznic/mathutil"

var sessionIdGen, _ = mathutil.NewFC32(0, 999999999, true)

func NewSessionId() uint {
	return uint(sessionIdGen.Next())
}
