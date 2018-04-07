package httptrack

import (
	"time"
)

func init() {
	Service = MemoryService{}
}

type MemoryService map[string][]*LoginDetails

func (m MemoryService) GetAccountAddrs(name string) ([]*LoginDetails, error) {
	if m[name] == nil {
		m[name] = make([]*LoginDetails, 0)
	}
	return m[name], nil
}

func (m MemoryService) SaveAccountAddr(name string, addr string, t time.Time) error {
	data, _ := m.GetAccountAddrs(name)
	m[name] = append(data, &LoginDetails{name, addr, t})
	return nil
}
