package sessions

import (
	"errors"
	"fmt"
	"strconv"
	"ztaylor.me/json"
)

var errArgNotFound = errors.New("msg: arg not found")
var errArgType = errors.New("msg: arg type conversion")

type WebsocketMessage struct {
	Name string
	Data json.Json
}

func NewWebsocketMessage() *WebsocketMessage {
	return &WebsocketMessage{
		Data: json.Json{},
	}
}

func (m *WebsocketMessage) Arg(k string) interface{} {
	return m.Data[k]
}

func (m *WebsocketMessage) IArg(k string) (int, error) {
	switch v := m.Arg(k).(type) {
	case int:
		return v, nil
	case string:
		i, err := strconv.ParseInt(v, 10, 0)
		return int(i), err
	case float64:
		return int(v), nil
	case nil:
		return 0, errArgNotFound
	default:
		return 0, errArgType
	}
}

func (m *WebsocketMessage) SArg(k string) (string, error) {
	switch v := m.Arg(k).(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%v", v), nil
	case int:
		return fmt.Sprintf("%v", v), nil
	default:
		return "", errArgNotFound
	}
}
