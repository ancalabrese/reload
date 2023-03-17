package yaml

import "gopkg.in/yaml.v3"

type Codec struct{}

func (Codec) Encode(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (Codec) Decode(b []byte, v any) error {
	return yaml.Unmarshal(b, v)
}
