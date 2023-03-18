package encoding

import (
	"io"
	"strings"

	"github.com/ancalabrese/reload/internal/encoding/json"
	"github.com/ancalabrese/reload/internal/encoding/yaml"
)

type Codec interface {
	Encode(v any) ([]byte, error)
	Decode(r io.Reader, v any) error
}

// New returns the right Codec based on the file type or nil if not suppported.
func New(mimeType string) Codec {
	if strings.Contains(mimeType, "json") {
		return json.Codec{}
	}

	if strings.Contains(mimeType, "yaml") || strings.Contains(mimeType, "yml") {
		return yaml.Codec{}
	}

	return nil
}
