package messages

type Message interface {
	Type() int
	Parse([]any) error
	Marshal() []any
}

type BinaryPayload interface {
	PayloadIsBinary() bool
	Payload() []byte
	PayloadSerializer() uint64
}
