package json

import (
	"encoding/json"
	"io"
)

// Encode wraps encoding/json.Marshal
func Encode(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{\"json-error\": " + err.Error() + "}"
	}
	return string(data)
}

// Decode wraps encoding/json.Decoder.Decode
func Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
