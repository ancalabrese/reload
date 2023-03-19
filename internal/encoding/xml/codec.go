package xml

import (
	"encoding/xml"
	"io"
)

type Codec struct{}

func (c Codec) Encode(v any) ([]byte, error) {
	return xml.Marshal(v)
}

func (c Codec) Decode(r io.Reader, v any) error {
	return xml.NewDecoder(r).Decode(v)
}
