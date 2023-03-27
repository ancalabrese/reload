package yaml

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Codec struct{}

func (Codec) Encode(w io.Writer, v any) error {
	b, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("yaml encoder: encoding error: %w", err)
	}
	_, err = w.Write(b)
	return err
}

func (Codec) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
