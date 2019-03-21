package track // import "ztaylor.me/http/track"

import "time"

// LoginDetails is a login event record
type LoginDetails struct {
	Name string
	Addr string
	time.Time
}

// Service is a storage connection for login details
type Service interface {
	GetAccountAddrs(name string) ([]*LoginDetails, error)
	SaveAccountAddr(name string, addr string, t time.Time) error
}

// Cache provides in-heap implementation of Service
type Cache map[string][]*LoginDetails

// GetAccountAddrs retrieves all login details
func (c Cache) GetAccountAddrs(name string) ([]*LoginDetails, error) {
	if c[name] == nil {
		c[name] = make([]*LoginDetails, 0)
	}
	return c[name], nil
}

// SaveAccountAddr adds a new login detail
func (c Cache) SaveAccountAddr(name string, addr string, t time.Time) error {
	c[name] = append(c[name], &LoginDetails{name, addr, t})
	return nil
}
