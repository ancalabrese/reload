package reload

import (
	"encoding/json"
)

type jsonCodec struct{}

func (jsonCodec) Encode(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", " ")
}

func (jsonCodec) Decode(b []byte, v any) error {
	return json.Unmarshal(b, v)
}
