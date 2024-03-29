package encoding

import (
	"io"
	"strings"

	"github.com/ancalabrese/reload/internal/encoding/json"
	"github.com/ancalabrese/reload/internal/encoding/xml"
	"github.com/ancalabrese/reload/internal/encoding/yaml"
)

type Codec interface {
	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v any) error
}

// New returns the right Codec based on the file type or nil if not suppported.
func New(fileExtension string) Codec {
	if strings.Contains(fileExtension, "json") {
		return json.Codec{}
	}

	if strings.Contains(fileExtension, "yaml") ||
		strings.Contains(fileExtension, "yml") {
		return yaml.Codec{}
	}

	if strings.Contains(fileExtension, "xml") {
		return xml.Codec{}
	}
	return nil
}
