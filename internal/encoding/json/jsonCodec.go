package json

import (
	"encoding/json"
	"fmt"
	"io"
)

type Codec struct{}

func (Codec) Encode(w io.Writer, v any) error {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Errorf("json encoder: encoding error: %w", err)
	}
	_, err = w.Write(b)
	return err
}

func (Codec) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
