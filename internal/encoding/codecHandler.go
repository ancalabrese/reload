package encoding

type CodecHandler interface {
	Encode(v any) ([]byte, error)
	Decode(b []byte, v any) error
}
