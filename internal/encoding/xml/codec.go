package xml

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Codec struct{}

func (Codec) Encode(w io.Writer, v any) error {
	b, err := xml.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Errorf("xml encoder: encoding error: %w", err)
	}
	_, err = w.Write(b)
	return err
}

func (c Codec) Decode(r io.Reader, v any) error {
	return xml.NewDecoder(r).Decode(v)
}
