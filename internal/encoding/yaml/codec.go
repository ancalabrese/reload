package yaml

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Codec struct{}

func (Codec) Encode(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (Codec) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
