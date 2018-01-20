package http

import (
	"fmt"
	"testing"
	"time"
)

func TestTokenSpeed(t *testing.T) {
	Tstart := time.Now()
	token := RegisterToken("test")
	d := time.Now().Sub(Tstart)
	fmt.Println("created token:", token, " time:", d)
}
