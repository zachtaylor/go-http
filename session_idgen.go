package http

import "github.com/cznic/mathutil"

var sessionIDGen, _ = mathutil.NewFC32(0, 999999999, true)

func NewSessionID() uint {
	return uint(sessionIDGen.Next())
}
