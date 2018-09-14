package sessions

import "time"

// Lifetime is the length of time a session is granted for
var Lifetime = 6 * time.Hour

// Grant is a record of a Session
type Grant struct {
	ID     string
	Name   string
	Time   time.Time
	Expire time.Time
}

// NewGrant creates a Grant with the specified name
func NewGrant(name string) *Grant {
	return &Grant{
		ID:     Keygen(),
		Name:   name,
		Time:   time.Now(),
		Expire: time.Now().Add(Lifetime),
	}
}
