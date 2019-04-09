package json

import (
	"encoding/json"
	"io"
)

// Encode wraps encoding/json.Marshal
func Encode(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return []byte("{\"json-error\": \"" + err.Error() + "\"}")
	}
	return data
}

// Decode wraps encoding/json.Decoder.Decode
func Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
