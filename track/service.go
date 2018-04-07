package httptrack

import (
	"time"
)

var Service interface {
	GetAccountAddrs(name string) ([]*LoginDetails, error)
	SaveAccountAddr(name string, addr string, t time.Time) error
}

type LoginDetails struct {
	Name string
	Addr string
	time.Time
}
