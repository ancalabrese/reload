package json

import (
	"encoding/json"
	"io"
)

type Codec struct{}

func (Codec) Encode(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", " ")
}

func (Codec) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
